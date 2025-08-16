package schedule

import (
	"log"
	"strconv"
	"tiktok-wails/backend/global"
	"time"
)

func InitSchedule() {
	parsedHour, err := strconv.Atoi(global.ScheduleSetting.RunAtTime)
	if err != nil {
		log.Printf("Failed to parse VALUE_RUN_AT_TIME: %v", err)
		return
	}
	hour := parsedHour
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
						// TODO: Thay thế bằng công việc thực tế
						println("Đã đến thời gian, thực hiện task!")
					}()
				}
			}
		}
	}()
}
