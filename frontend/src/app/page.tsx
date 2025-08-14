"use client"

import { useState, useEffect } from "react"
import { MainHeader } from "@/components/main-header"
import { ProfileTab } from "@/components/account-tab"
import { VideoTab } from "@/components/video-tab"
import { SettingsTab } from "@/components/settings-tab"
import "./globals.css" // Import global CSS

export default function Home() {
  const [activeTab, setActiveTab] = useState("account")
  const [isDarkMode, setIsDarkMode] = useState(false)

  useEffect(() => {
    const savedTheme = localStorage.getItem("theme")
    if (savedTheme === "dark") {
      setIsDarkMode(true)
      document.body.classList.add("dark-mode")
    } else {
      setIsDarkMode(false)
      document.body.classList.remove("dark-mode")
    }
  }, [])

  const toggleTheme = () => {
    setIsDarkMode((prevMode) => {
      const newMode = !prevMode
      if (newMode) {
        document.body.classList.add("dark-mode")
        localStorage.setItem("theme", "dark")
      } else {
        document.body.classList.remove("dark-mode")
        localStorage.setItem("theme", "light")
      }
      return newMode
    })
  }

  return (
    <div className="app-container">
      <MainHeader toggleTheme={toggleTheme} isDarkMode={isDarkMode} />
      <div className="main-content">
        <div className="tabs-root">
          <div className="tabs-list">
            <button
              className={`tabs-trigger ${activeTab === "account" ? "active" : ""}`}
              onClick={() => setActiveTab("account")}
            >
              Quản lý Account
            </button>
            <button
              className={`tabs-trigger ${activeTab === "video" ? "active" : ""}`}
              onClick={() => setActiveTab("video")}
            >
              Quản lý Video
            </button>
            <button
              className={`tabs-trigger ${activeTab === "settings" ? "active" : ""}`}
              onClick={() => setActiveTab("settings")}
            >
              Cài đặt
            </button>
          </div>
          <div className="tabs-content">
            {activeTab === "account" && <ProfileTab />}
            {activeTab === "video" && <VideoTab />}
            {activeTab === "settings" && <SettingsTab />}
          </div>
        </div>
      </div>
      <div className="bottom-status-bar">
        {/* <div className="status-left">
          <span className="font-medium">Đỗ Hữu Ben (Vĩnh viễn)</span>
        </div> */}
        <div className="status-right">
          <div className="checkbox-group">
            <input type="checkbox" id="auto-delete" className="custom-checkbox" />
            <label htmlFor="auto-delete">Auto Delete</label>
          </div>
          <div className="checkbox-group">
            <input type="checkbox" id="auto-shutdown" className="custom-checkbox" />
            <label htmlFor="auto-shutdown">Auto Shutdown</label>
          </div>
        </div>
      </div>
    </div>
  )
}
