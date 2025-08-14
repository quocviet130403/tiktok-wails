"use client"

import { useState } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { ThemeProvider } from "@/components/theme-provider"
import { MainHeader } from "@/components/main-header"
import { VideoTab } from "@/components/video-tab"
import { SettingsTab } from "@/components/settings-tab"
import { ProfileTab } from "@/components/account-tab"

export default function App() {
  const [activeTab, setActiveTab] = useState("account")
  const [isDarkMode, setIsDarkMode] = useState(false)
  const toggleTheme = () => setIsDarkMode((prev) => !prev)

  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <div className="flex flex-col h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-50">
        <MainHeader toggleTheme={toggleTheme} isDarkMode={isDarkMode} />
        <div className="flex-1 p-4 overflow-auto">
          <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full h-full flex flex-col">
            <TabsList className="grid w-full grid-cols-3 bg-gray-200 dark:bg-gray-800 rounded-lg p-1 mb-4">
              <TabsTrigger
                value="account"
                className="data-[state=active]:bg-white data-[state=active]:text-gray-900 dark:data-[state=active]:bg-gray-700 dark:data-[state=active]:text-gray-50 transition-colors duration-200"
              >
                Quản lý Account
              </TabsTrigger>
              <TabsTrigger
                value="video"
                className="data-[state=active]:bg-white data-[state=active]:text-gray-900 dark:data-[state=active]:bg-gray-700 dark:data-[state=active]:text-gray-50 transition-colors duration-200"
              >
                Quản lý Video
              </TabsTrigger>
              <TabsTrigger
                value="settings"
                className="data-[state=active]:bg-white data-[state=active]:text-gray-900 dark:data-[state=active]:bg-gray-700 dark:data-[state=active]:text-gray-50 transition-colors duration-200"
              >
                Cài đặt
              </TabsTrigger>
            </TabsList>
            <TabsContent
              value="account"
              className="flex-1 p-4 bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-auto"
            >
              <ProfileTab />
            </TabsContent>
            <TabsContent
              value="video"
              className="flex-1 p-4 bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-auto"
            >
              <VideoTab />
            </TabsContent>
            <TabsContent
              value="settings"
              className="flex-1 p-4 bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-auto"
            >
              <SettingsTab />
            </TabsContent>
          </Tabs>
        </div>
        <div className="flex items-center justify-between p-2 bg-gray-200 dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700 text-sm">
          <div className="flex items-center gap-2">
            <span className="font-medium">Đỗ Hữu Ben (Vĩnh viễn)</span>
          </div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1">
              <input
                type="checkbox"
                id="auto-delete"
                className="form-checkbox h-4 w-4 text-blue-600 rounded focus:ring-blue-500"
              />
              <label htmlFor="auto-delete" className="text-gray-700 dark:text-gray-300">
                Auto Delete
              </label>
            </div>
            <div className="flex items-center gap-1">
              <input
                type="checkbox"
                id="auto-shutdown"
                className="form-checkbox h-4 w-4 text-blue-600 rounded focus:ring-blue-500"
              />
              <label htmlFor="auto-shutdown" className="text-gray-700 dark:text-gray-300">
                Auto Shutdown
              </label>
            </div>
          </div>
        </div>
      </div>
    </ThemeProvider>
  )
}
