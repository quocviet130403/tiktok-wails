"use client"

import { useState } from "react"
import { ThemeProvider } from "@/components/theme-provider"
import { ModeToggle } from "@/components/ui/mode-toggle"
import { Toaster } from "@/components/ui/sonner"
import { VideoTab } from "@/components/video-tab"
import { SettingsTab } from "@/components/settings-tab"
import { ProfileTab } from "@/components/account-tab"
import { ProfileDouyinTab } from "@/components/profile-douyin-tab"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import { Badge } from "@/components/ui/badge"
import { Chrome, Video, Settings, CircleDot } from "lucide-react"

type TabKey = "profiles" | "douyin" | "videos" | "settings"

const navItems: { key: TabKey; icon: React.ElementType; label: string; color: string }[] = [
  { key: "profiles", icon: Chrome, label: "Chrome Profiles", color: "text-blue-500" },
  { key: "douyin", icon: CircleDot, label: "Douyin Profiles", color: "text-pink-500" },
  { key: "videos", icon: Video, label: "Videos", color: "text-violet-500" },
  { key: "settings", icon: Settings, label: "Cài đặt", color: "text-gray-400" },
]

export default function App() {
  const [activeTab, setActiveTab] = useState<TabKey>("profiles")

  const renderContent = () => {
    switch (activeTab) {
      case "profiles": return <ProfileTab />
      case "douyin": return <ProfileDouyinTab />
      case "videos": return <VideoTab />
      case "settings": return <SettingsTab />
    }
  }

  return (
    <ThemeProvider attribute="class" defaultTheme="dark" enableSystem>
      <TooltipProvider delayDuration={100}>
        <div className="flex h-screen bg-background text-foreground">

          {/* Sidebar */}
          <aside className="w-[60px] flex flex-col items-center border-r border-border bg-card py-4 gap-1">
            {/* Logo */}
            <div className="mb-4 w-9 h-9 rounded-lg bg-gradient-to-br from-violet-500 to-blue-500 flex items-center justify-center">
              <span className="text-white font-bold text-sm">TR</span>
            </div>

            {/* Nav Items */}
            <nav className="flex flex-col gap-1 flex-1">
              {navItems.map((item) => {
                const isActive = activeTab === item.key
                const Icon = item.icon
                return (
                  <Tooltip key={item.key}>
                    <TooltipTrigger asChild>
                      <button
                        onClick={() => setActiveTab(item.key)}
                        className={`
                          w-10 h-10 rounded-lg flex items-center justify-center transition-all duration-200
                          ${isActive
                            ? "bg-primary/10 shadow-sm"
                            : "hover:bg-accent"
                          }
                        `}
                      >
                        <Icon className={`h-5 w-5 ${isActive ? item.color : "text-muted-foreground"}`} />
                      </button>
                    </TooltipTrigger>
                    <TooltipContent side="right" sideOffset={8}>
                      {item.label}
                    </TooltipContent>
                  </Tooltip>
                )
              })}
            </nav>

            {/* Bottom: Theme Toggle */}
            <div className="mt-auto">
              <ModeToggle />
            </div>
          </aside>

          {/* Main Area */}
          <div className="flex-1 flex flex-col min-w-0">

            {/* Top Bar */}
            <header className="h-12 border-b border-border bg-card/50 backdrop-blur-sm flex items-center justify-between px-4 shrink-0">
              <div className="flex items-center gap-3">
                <h1 className="text-sm font-semibold tracking-tight">
                  {navItems.find(n => n.key === activeTab)?.label}
                </h1>
              </div>
              <div className="flex items-center gap-3">
                <Badge variant="secondary" className="text-xs gap-1.5 font-normal">
                  <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                  Pipeline Idle
                </Badge>
              </div>
            </header>

            {/* Content */}
            <main className="flex-1 overflow-auto p-4">
              {renderContent()}
            </main>

          </div>
        </div>
        <Toaster />
      </TooltipProvider>
    </ThemeProvider>
  )
}
