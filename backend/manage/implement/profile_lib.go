package implement

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"tiktok-wails/backend/global"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"gocv.io/x/gocv"
)

type ListVideos struct {
	AwemeLists []struct {
		Desc      string `json:"desc"`
		Duration  int    `json:"duration"`
		Statistic struct {
			LikeCount int `json:"digg_count"`
		} `json:"statistics"`
		Video struct {
			PlayAddr struct {
				URLList []string `json:"url_list"`
			} `json:"play_addr"`
			Cover struct {
				URLList []string `json:"url_list"`
			} `json:"cover"`
		} `json:"video"`
	} `json:"aweme_list"`
}

// PuzzleCaptchaSolver holds paths to images and provides methods to solve the CAPTCHA.
type PuzzleCaptchaSolver struct {
	GapImagePath    string
	BgImagePath     string
	OutputImagePath string
}

// NewPuzzleCaptchaSolver creates and returns a new PuzzleCaptchaSolver instance.
func NewPuzzleCaptchaSolver(gapImagePath, bgImagePath, outputImagePath string) *PuzzleCaptchaSolver {
	return &PuzzleCaptchaSolver{
		GapImagePath:    gapImagePath,
		BgImagePath:     bgImagePath,
		OutputImagePath: outputImagePath,
	}
}

// RemoveWhitespace crops an image to the area containing non-whitespace pixels.
func (s *PuzzleCaptchaSolver) RemoveWhitespace(imagePath string) (gocv.Mat, error) {
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return gocv.Mat{}, fmt.Errorf("failed to read image: %s", imagePath)
	}
	defer img.Close()

	rows := img.Rows()
	cols := img.Cols()

	minX, minY := rows, cols
	maxX, maxY := 0, 0

	// Iterate through pixels to find the bounding box of non-whitespace
	// In GoCV, you access pixel data directly if performance is critical,
	// but for simplicity, we'll iterate with GetVecb, which can be slower.
	// For production, consider direct pointer access.
	for x := 0; x < rows; x++ {
		for y := 0; y < cols; y++ {
			pixel := img.GetVecbAt(x, y) // Returns a []byte for the pixel (BGR)
			// Check if pixel is not "pure" white (e.g., all channels 255)
			// Or if there's significant variation, indicating non-whitespace
			// This check is a simplification; a more robust check might involve
			// checking for a minimum difference between channels or absolute values.
			if pixel[0] < 250 || pixel[1] < 250 || pixel[2] < 250 { // If any channel is not almost white
				minX = min(x, minX)
				minY = min(y, minY)
				maxX = max(x, maxX)
				maxY = max(y, maxY)
			}
		}
	}

	// If no non-whitespace pixels found, return empty Mat or an error
	if minX > maxX || minY > maxY {
		return gocv.Mat{}, fmt.Errorf("no non-whitespace pixels found in image: %s", imagePath)
	}

	// Define the rectangle for cropping
	rect := image.Rect(minY, minX, maxY, maxX) // gocv.Mat takes rect as (col1, row1, col2, row2)

	// Crop the image
	croppedImg := img.Region(rect)
	return croppedImg, nil
}

// ApplyEdgeDetection applies Canny edge detection to the given image.
func (s *PuzzleCaptchaSolver) ApplyEdgeDetection(img gocv.Mat) (gocv.Mat, error) {
	if img.Empty() {
		return gocv.Mat{}, fmt.Errorf("input image for edge detection is empty")
	}

	// Convert to grayscale
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorRGBToGray)

	// Apply Canny edge detection
	edges := gocv.NewMat()
	defer edges.Close()
	gocv.Canny(gray, &edges, 100, 200)

	// Convert edges back to RGB for consistency if needed, otherwise keep grayscale
	edgesRGB := gocv.NewMat()
	gocv.CvtColor(edges, &edgesRGB, gocv.ColorGrayToRGB)
	return edgesRGB, nil
}

// FindPositionOfSlide finds the x-coordinate of the top-left corner of the slide in the background picture.
func (s *PuzzleCaptchaSolver) FindPositionOfSlide(slidePic, backgroundPic gocv.Mat) (int, error) {
	if slidePic.Empty() || backgroundPic.Empty() {
		return 0, fmt.Errorf("slide or background image for template matching is empty")
	}

	// Perform template matching
	result := gocv.NewMat()
	defer result.Close()
	gocv.MatchTemplate(backgroundPic, slidePic, &result, gocv.TmCcoeffNormed, gocv.NewMat())

	// Find the minimum and maximum values and their locations
	_, _, _, maxLoc := gocv.MinMaxLoc(result)

	// The top-left corner of the best match is maxLoc
	tl := maxLoc

	// Draw rectangle around the found slide for visualization
	tplHeight := slidePic.Rows()
	tplWidth := slidePic.Cols()
	br := image.Point{X: tl.X + tplWidth, Y: tl.Y + tplHeight}
	gocv.Rectangle(&backgroundPic, image.Rectangle{Min: tl, Max: br}, color.RGBA{R: 255, A: 255}, 2) // Red color, 2px thickness

	// Save the output image
	if ok := gocv.IMWrite(s.OutputImagePath, backgroundPic); !ok {
		return 0, fmt.Errorf("failed to write output image: %s", s.OutputImagePath)
	}

	return tl.X, nil
}

// Discern performs the discernment process to find the position of the slide.
func (s *PuzzleCaptchaSolver) Discern() (int, error) {
	// Remove whitespace from the gap image
	gapImage, err := s.RemoveWhitespace(s.GapImagePath)
	if err != nil {
		return 0, fmt.Errorf("error removing whitespace from gap image: %w", err)
	}
	defer gapImage.Close()

	// Apply edge detection to the gap image
	edgeDetectedGap, err := s.ApplyEdgeDetection(gapImage)
	if err != nil {
		return 0, fmt.Errorf("error applying edge detection to gap image: %w", err)
	}
	defer edgeDetectedGap.Close()

	// Read the background image
	bgImage := gocv.IMRead(s.BgImagePath, gocv.IMReadColor)
	if bgImage.Empty() {
		return 0, fmt.Errorf("failed to read background image: %s", s.BgImagePath)
	}
	defer bgImage.Close()

	// Apply edge detection to the background image
	edgeDetectedBg, err := s.ApplyEdgeDetection(bgImage)
	if err != nil {
		return 0, fmt.Errorf("error applying edge detection to background image: %w", err)
	}
	defer edgeDetectedBg.Close()

	// Find the position of the slide
	slidePosition, err := s.FindPositionOfSlide(edgeDetectedGap, edgeDetectedBg)
	if err != nil {
		return 0, fmt.Errorf("error finding position of slide: %w", err)
	}

	return slidePosition, nil
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func DownloadImage(url, filePath string) error {
	// Bước 1: Gửi HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("lỗi khi gửi request HTTP: %w", err)
	}
	defer resp.Body.Close() // Đảm bảo đóng body sau khi hoàn thành

	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("không thể tải ảnh, mã trạng thái: %d %s", resp.StatusCode, resp.Status)
	}

	// Bước 2: Tạo file để lưu ảnh
	// Đảm bảo thư mục đích tồn tại
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("không thể tạo thư mục '%s': %w", dir, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("lỗi khi tạo file '%s': %w", filePath, err)
	}
	defer file.Close() // Đảm bảo đóng file sau khi hoàn thành

	// Bước 3: Sao chép dữ liệu từ response body vào file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("lỗi khi sao chép dữ liệu vào file: %w", err)
	}

	fmt.Printf("Đã tải ảnh thành công: %s\n", filePath)
	return nil
}

func handleCaptchaAsync(ctx context.Context, profileName string) {
	// Kiểm tra CAPTCHA với timeout dài hơn
	iframeSelector := `iframe[src*="rmc.bytedance.com/verifycenter/captcha/v2"]`

	// Tạo context riêng cho việc kiểm tra CAPTCHA
	captchaCtx, captchaCancel := context.WithTimeout(ctx, 30*time.Second)
	defer captchaCancel()

	log.Println("Bắt đầu kiểm tra CAPTCHA trong background...")

	// Kiểm tra CAPTCHA mỗi 3 giây
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-captchaCtx.Done():
			log.Println("Dừng kiểm tra CAPTCHA do timeout hoặc context cancel")
			return
		case <-ticker.C:
			// Kiểm tra xem có iframe CAPTCHA không
			var iframeExists bool
			err := chromedp.Run(captchaCtx,
				chromedp.EvaluateAsDevTools(`
                    document.querySelector('iframe[src*="rmc.bytedance.com/verifycenter/captcha/v2"]') !== null
                `, &iframeExists),
			)

			if err != nil {
				log.Printf("Lỗi khi kiểm tra CAPTCHA: %v", err)
				continue
			}

			if iframeExists {
				log.Println("Phát hiện CAPTCHA, bắt đầu xử lý...")
				if err := handleCaptcha(captchaCtx, profileName, iframeSelector); err != nil {
					log.Printf("Lỗi khi xử lý CAPTCHA: %v", err)
					// Tiếp tục kiểm tra nếu CAPTCHA fail
					continue
				}
				log.Println("Đã xử lý CAPTCHA thành công")
				return // Thoát sau khi xử lý thành công
			}
		}
	}
}

func handleCaptcha(ctx context.Context, profileName, iframeSelector string) error {
	var iframeNodeIDs []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Nodes(iframeSelector, &iframeNodeIDs, chromedp.ByQuery),
	)
	if err != nil {
		return fmt.Errorf("không thể lấy iframe nodes: %w", err)
	}

	if len(iframeNodeIDs) == 0 {
		return fmt.Errorf("không tìm thấy NodeID cho iframe CAPTCHA")
	}

	iframeNode := iframeNodeIDs[0]
	log.Printf("Đã tìm thấy iframe CAPTCHA với NodeID: %d", iframeNode.NodeID)

	if iframeNode.FrameID == "" {
		return fmt.Errorf("không tìm thấy FrameID cho iframe")
	}

	iframeTargetID := string(iframeNode.FrameID)
	iframeCtx, cancelIframe := chromedp.NewContext(ctx, chromedp.WithTargetID(target.ID(iframeTargetID)))
	defer cancelIframe()

	var btnX, btnY float64
	var rectData []float64
	var bgImage, slideImage string

	sliderBtnSelector := `.captcha-slider-btn`

	err = chromedp.Run(iframeCtx,
		network.Enable(),
		chromedp.WaitVisible(sliderBtnSelector, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
		chromedp.EvaluateAsDevTools(fmt.Sprintf(`
            (function() {
                const el = document.querySelector('%s');
                if (el) {
                    const rect = el.getBoundingClientRect();
                    return [rect.x, rect.y, rect.width, rect.height];
                }
                return null;
            })();
        `, sliderBtnSelector), &rectData),

		chromedp.EvaluateAsDevTools(`
            (function() {
                const el = document.querySelector('#captcha_verify_image');
                if (el) {
                    return el.getAttribute('src');
                }
                return null;
            })();
        `, &bgImage),

		chromedp.EvaluateAsDevTools(`
            (function() {
                const el = document.querySelector('#captcha-verify_img_slide');
                if (el) {
                    return el.getAttribute('src');
                }
                return null;
            })();
        `, &slideImage),
	)

	if err != nil {
		return fmt.Errorf("lỗi trong ngữ cảnh iframe: %w", err)
	}

	if rectData == nil || len(rectData) < 2 {
		return fmt.Errorf("không thể lấy tọa độ hợp lệ của nút CAPTCHA")
	}

	// Tải và xử lý ảnh CAPTCHA
	fmt.Printf("Đang tải ảnh nền từ: %s\n", bgImage)
	if err := DownloadImage(bgImage, global.PathHandleCaptcha+`/bg-`+profileName+`.png`); err != nil {
		return fmt.Errorf("lỗi khi tải ảnh nền: %w", err)
	}

	fmt.Printf("Đang tải ảnh trượt từ: %s\n", slideImage)
	if err := DownloadImage(slideImage, global.PathHandleCaptcha+`/slide-`+profileName+`.png`); err != nil {
		return fmt.Errorf("lỗi khi tải ảnh trượt: %w", err)
	}

	time.Sleep(2 * time.Second)

	solver := NewPuzzleCaptchaSolver(
		global.PathHandleCaptcha+`/slide-`+profileName+`.png`,
		global.PathHandleCaptcha+`/bg-`+profileName+`.png`,
		global.PathHandleCaptcha+`/result-`+profileName+`.png`,
	)

	position, err := solver.Discern()
	if err != nil {
		return fmt.Errorf("lỗi khi giải CAPTCHA: %w", err)
	}
	fmt.Printf("Vị trí của slide: %d\n", position)

	btnX = rectData[0]
	btnY = rectData[1]
	targetX := float64(position*340/552) + btnX
	targetY := btnY

	// Thực hiện kéo thả
	jsDragScript := fmt.Sprintf(`
        (async function() {
            const sliderBtn = document.querySelector('%s');
            if (!sliderBtn) {
                console.error("Nút trượt không tìm thấy trong iframe.");
                return false;
            }

            const startX = %f;
            const startY = %f;
            const endX = %f;
            const endY = %f;

            // Gửi sự kiện mousedown
            const mouseDownEvent = new MouseEvent('mousedown', {
                bubbles: true, cancelable: true, view: window,
                clientX: startX, clientY: startY, buttons: 1
            });
            sliderBtn.dispatchEvent(mouseDownEvent);

            // Tạo đường cong Bezier để mô phỏng chuyển động tự nhiên
            const steps = Math.floor(Math.random() * 20) + 25;
            const totalDistance = endX - startX;
            
            const controlX = startX + totalDistance * (0.3 + Math.random() * 0.4);
            const controlY = startY + (Math.random() - 0.5) * 10;

            for (let i = 1; i <= steps; i++) {
                const t = i / steps;
                
                const easedT = t < 0.5 
                    ? 2 * t * t 
                    : 1 - Math.pow(-2 * t + 2, 3) / 2;
                
                const currentX = Math.pow(1 - easedT, 2) * startX + 
                               2 * (1 - easedT) * easedT * controlX + 
                               Math.pow(easedT, 2) * endX;
                
                const currentY = Math.pow(1 - easedT, 2) * startY + 
                               2 * (1 - easedT) * easedT * controlY + 
                               Math.pow(easedT, 2) * endY;

                const mouseMoveEvent = new MouseEvent('mousemove', {
                    bubbles: true, cancelable: true, view: window,
                    clientX: currentX,
                    clientY: currentY,
                    buttons: 1
                });
                sliderBtn.dispatchEvent(mouseMoveEvent);

                const baseDelay = 8;
                const randomDelay = Math.random() * 12;
                const pauseChance = Math.random();
                
                if (pauseChance < 0.1 && i > steps * 0.3) {
                    await new Promise(r => setTimeout(r, 20 + Math.random() * 30));
                }
                
                await new Promise(r => setTimeout(r, baseDelay + randomDelay));
            }

            // Thêm một chút dao động cuối
            const finalAdjustments = Math.floor(Math.random() * 3) + 1;
            for (let j = 0; j < finalAdjustments; j++) {
                const adjustX = endX + (Math.random() - 0.5) * 4;
                const adjustEvent = new MouseEvent('mousemove', {
                    bubbles: true, cancelable: true, view: window,
                    clientX: adjustX, clientY: endY, buttons: 1
                });
                sliderBtn.dispatchEvent(adjustEvent);
                await new Promise(r => setTimeout(r, 15 + Math.random() * 10));
            }

            // Gửi sự kiện mouseup
            const mouseUpEvent = new MouseEvent('mouseup', {
                bubbles: true, cancelable: true, view: window,
                clientX: endX, clientY: endY, buttons: 0
            });
            sliderBtn.dispatchEvent(mouseUpEvent);

            return true;
        })();
    `, sliderBtnSelector, btnX, btnY, targetX, targetY)

	err = chromedp.Run(iframeCtx,
		chromedp.EvaluateAsDevTools(jsDragScript, nil),
		chromedp.Sleep(5*time.Second),
	)
	if err != nil {
		return fmt.Errorf("không thể kéo nút CAPTCHA: %w", err)
	}

	return nil
}
