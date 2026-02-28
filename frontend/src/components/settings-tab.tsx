"use client"

import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Separator } from "@/components/ui/separator"
import { Save, RotateCcw } from "lucide-react"
import { GetAllSettings, SetSetting } from "../../wailsjs/go/backend/App"

export function SettingsTab() {
  const [settings, setSettings] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  const fetchSettings = async () => {
    try {
      const result = await GetAllSettings()
      if (result) {
        setSettings(result)
      }
    } catch (err) {
      console.error("Lỗi khi tải settings:", err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSettings()
  }, [])

  const handleSave = async (key: string, value: string) => {
    setSaving(true)
    try {
      await SetSetting(key, value)
      setSettings((prev) => ({ ...prev, [key]: value }))
    } catch (err) {
      console.error(`Lỗi khi lưu ${key}:`, err)
    } finally {
      setSaving(false)
    }
  }

  const handleChange = (key: string, value: string) => {
    setSettings((prev) => ({ ...prev, [key]: value }))
  }

  const handleSaveAll = async () => {
    setSaving(true)
    try {
      for (const [key, value] of Object.entries(settings)) {
        await SetSetting(key, value)
      }
      alert("Đã lưu tất cả settings!")
    } catch (err) {
      console.error("Lỗi khi lưu settings:", err)
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <div className="flex items-center justify-center h-full text-gray-500">Đang tải settings...</div>
  }

  return (
    <div className="flex flex-col gap-6 p-4">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {/* Chrome Path */}
        <div className="grid gap-4 p-4 border rounded-lg bg-gray-50 dark:bg-gray-700 shadow-sm">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-50">Cài đặt Chrome</h3>
          <div className="grid gap-2">
            <Label htmlFor="path_chrome" className="text-gray-700 dark:text-gray-300">
              Đường dẫn Chrome
            </Label>
            <Input
              id="path_chrome"
              value={settings.path_chrome || ""}
              onChange={(e) => handleChange("path_chrome", e.target.value)}
              placeholder="C:/Program Files/Google/Chrome/Application/chrome.exe"
              className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
            />
          </div>
        </div>

        {/* Schedule Settings */}
        <div className="grid gap-4 p-4 border rounded-lg bg-gray-50 dark:bg-gray-700 shadow-sm">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-50">Lịch chạy tự động</h3>
          <div className="grid gap-2">
            <Label htmlFor="schedule_time" className="text-gray-700 dark:text-gray-300">
              Loại lịch
            </Label>
            <Select
              value={settings.schedule_time || "daily"}
              onValueChange={(value) => handleChange("schedule_time", value)}
            >
              <SelectTrigger
                id="schedule_time"
                className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
              >
                <SelectValue placeholder="Chọn loại lịch" />
              </SelectTrigger>
              <SelectContent className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
                <SelectItem value="daily">Hàng ngày</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="grid gap-2">
            <Label htmlFor="run_at_time" className="text-gray-700 dark:text-gray-300">
              Giờ chạy (0-23, đặt 24 để tắt)
            </Label>
            <Input
              id="run_at_time"
              type="number"
              min="0"
              max="24"
              value={settings.run_at_time || "24"}
              onChange={(e) => handleChange("run_at_time", e.target.value)}
              className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
            />
            <p className="text-xs text-gray-500 dark:text-gray-400">
              Pipeline sẽ chạy tuần tự: Scrape → Upload → Xóa → Kiểm tra Auth
            </p>
          </div>
        </div>

        {/* Info */}
        <div className="grid gap-4 p-4 border rounded-lg bg-gray-50 dark:bg-gray-700 shadow-sm">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-50">Thông tin</h3>
          <div className="text-sm text-gray-600 dark:text-gray-400 space-y-2">
            <p>📋 Tổng settings: <strong>{Object.keys(settings).length}</strong></p>
            {Object.entries(settings).map(([key, value]) => (
              <div key={key} className="flex justify-between items-center py-1 border-b border-gray-200 dark:border-gray-600">
                <span className="font-mono text-xs">{key}</span>
                <span className="text-xs truncate max-w-[150px]">{value}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Save Button */}
      <Separator className="bg-gray-300 dark:bg-gray-600" />
      <div className="flex items-center gap-2">
        <Button
          onClick={handleSaveAll}
          disabled={saving}
          className="bg-green-500 hover:bg-green-600 text-white transition-colors duration-200"
        >
          <Save className="h-4 w-4 mr-2" />
          {saving ? "Đang lưu..." : "Lưu tất cả"}
        </Button>
        <Button
          variant="outline"
          onClick={fetchSettings}
          className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
        >
          <RotateCcw className="h-4 w-4 mr-2" />
          Tải lại
        </Button>
      </div>
    </div>
  )
}
