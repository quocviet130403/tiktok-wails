package initialize

import (
	"log"
	"strconv"
	"tiktok-wails/backend/global"
	"tiktok-wails/backend/manage/service"
	"time"
)

func InitSchedule() {
	parsedHour, err := strconv.Atoi(global.ScheduleSetting.RunAtTime)
	if err != nil {
		log.Printf("Failed to parse VALUE_RUN_AT_TIME: %v", err)
		return
	}
	hour := parsedHour

	// Get Video From Douyin
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for now := range ticker.C {
			run := false
			switch global.ScheduleSetting.Time {
			case "daily":
				if now.Hour() == hour {
					run = true
				}
				if run {
					go func() {
						profileDouyins, err := service.ProfileDouyinManager().GetAllProfiles()
						if err != nil {
							log.Printf("Lỗi khi lấy danh sách profile Douyin: %v", err)
							return
						}

						for _, profile := range profileDouyins {
							go func(profile service.ProfileDouyin) {
								retry := 0
								for {
									err := service.ProfileDouyinManager().GetVideoFromProfile(profile)
									if err != nil {
										log.Printf("Lỗi khi lấy video từ profile %s: %v", profile.Nickname, err)
										retry++
										if retry >= 3 {
											log.Printf("Đã thử lại %d lần, bỏ qua profile %s", retry, profile.Nickname)
											break
										}
										time.Sleep(2 * time.Second)
										continue
									}
									break
								}
							}(profile)
						}
					}()
				}
			}
		}
	}()

	// Reup Video Auto
	tickerReup := time.NewTicker(1 * time.Second)
	go func() {
		for now := range tickerReup.C {
			run := false
			switch global.ScheduleSetting.Time {
			case "daily":
				if now.Hour() == hour+2 {
					run = true
				}
				if run {
					go func() {
						// logic code
						profiles, err := service.ProfileManager().GetAllProfiles()
						if err != nil {
							log.Printf("Lỗi khi lấy danh sách profile: %v", err)
							return
						}

						for _, profile := range profiles {
							go func(profile service.Profiles) {
								// Logic code
								allVideos, err := service.VideoManager().GetVideoReup(profile.ID)
								if err != nil {
									log.Printf("Lỗi khi lấy danh sách video từ profile: %v", err)
									return
								}

								for _, video := range allVideos {
									err := service.VideoManager().UploadVideo(profile.Name, video.Title, video.Title+" "+profile.Hashtag)
									if err != nil {
										log.Printf("Lỗi khi upload video: %v", err)
										return
									}

									err = service.VideoManager().UpdateStatusReup(video.ID, profile.ID, "done")
									if err != nil {
										log.Printf("Lỗi khi cập nhật trạng thái video: %v", err)
										return
									}
								}

							}(profile)
						}
					}()
				}
			}
		}
	}()

	// Delete uploaded videos
	tickerDelete := time.NewTicker(1 * time.Second)
	go func() {
		for now := range tickerDelete.C {
			run := false
			switch global.ScheduleSetting.Time {
			case "daily":
				if now.Hour() == hour+4 {
					run = true
				}
				if run {
					go func() {
						getAllVideos, err := service.VideoManager().GetAllVideosNP()
						if err != nil {
							log.Printf("Lỗi khi lấy danh sách video: %v", err)
							return
						}

						for _, video := range getAllVideos {
							count, err := service.VideoManager().GetCompleteProfileVideos(video.ID)
							if err != nil {
								log.Printf("Lỗi khi đếm video đã được upload: %v", err)
								return
							}

							if count == 0 {
								// Xóa video
								err := service.VideoManager().DeleteVideo(video)
								if err != nil {
									log.Printf("Lỗi khi xóa video: %v", err)
									return
								}
							}
						}
					}()
				}
			}
		}
	}()

	// Check Authen Profile
	tickerCheckAuthen := time.NewTicker(1 * time.Second)
	go func() {
		for now := range tickerCheckAuthen.C {
			run := false
			switch global.ScheduleSetting.Time {
			case "daily":
				if now.Hour() == hour+6 {
					run = true
				}
				if run {
					go func() {
						profiles, err := service.ProfileManager().GetAllProfileCheckAuthenticated()
						if err != nil {
							log.Printf("Lỗi khi lấy danh sách profile: %v", err)
							return
						}

						for _, profile := range profiles {
							go func(profile service.Profiles) {
								// Logic code
								err := service.VideoManager().LoginTiktok(profile.Name)
								if err != nil {
									log.Printf("Lỗi khi kiểm tra xác thực profile: %v", err)
									err := service.ProfileManager().UpdateAuthenticatedStatus(profile.ID, false)
									if err != nil {
										log.Printf("Lỗi khi cập nhật trạng thái xác thực profile: %v", err)
									}
									return
								}
							}(profile)
						}
					}()
				}
			}
		}
	}()
}
