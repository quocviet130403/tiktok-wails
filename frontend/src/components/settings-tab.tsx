"use client"

import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Separator } from "@/components/ui/separator"
import { Save, RotateCcw, Chrome, Clock, Info } from "lucide-react"
import { toast } from "sonner"
import { GetAllSettings, SetSetting } from "../../wailsjs/go/backend/App"

export function SettingsTab() {
  const [settings, setSettings] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  const fetchSettings = async () => {
    try {
      const result = await GetAllSettings()
      if (result) setSettings(result)
    } catch (err) {
      console.error("Lỗi khi tải settings:", err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchSettings() }, [])

  const handleChange = (key: string, value: string) => {
    setSettings((prev) => ({ ...prev, [key]: value }))
  }

  const handleSaveAll = async () => {
    setSaving(true)
    try {
      for (const [key, value] of Object.entries(settings)) {
        await SetSetting(key, value)
      }
      toast.success("Đã lưu tất cả cài đặt!")
    } catch (err) {
      toast.error("Lỗi khi lưu cài đặt")
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full text-muted-foreground">
        Đang tải cài đặt...
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6 max-w-3xl">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Chrome */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-medium flex items-center gap-2">
              <Chrome className="h-4 w-4 text-blue-500" />
              Chrome Browser
            </CardTitle>
            <CardDescription className="text-xs">Đường dẫn tới Chrome executable</CardDescription>
          </CardHeader>
          <CardContent>
            <Input
              value={settings.path_chrome || ""}
              onChange={(e) => handleChange("path_chrome", e.target.value)}
              placeholder="C:/Program Files/Google/Chrome/..."
              className="h-9 text-sm font-mono"
            />
          </CardContent>
        </Card>

        {/* Schedule */}
        <Card>
          <CardHeader className="pb-3">
            <CardTitle className="text-sm font-medium flex items-center gap-2">
              <Clock className="h-4 w-4 text-violet-500" />
              Lịch tự động
            </CardTitle>
            <CardDescription className="text-xs">Pipeline: Scrape → Upload → Xóa → Auth</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="space-y-1.5">
              <Label className="text-xs text-muted-foreground">Kiểu lịch</Label>
              <Select value={settings.schedule_time || "daily"} onValueChange={(v) => handleChange("schedule_time", v)}>
                <SelectTrigger className="h-9">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="daily">Hàng ngày</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-1.5">
              <Label className="text-xs text-muted-foreground">Giờ bắt đầu (đặt 24 để tắt)</Label>
              <Input
                type="number"
                min="0"
                max="24"
                value={settings.run_at_time || "24"}
                onChange={(e) => handleChange("run_at_time", e.target.value)}
                className="h-9 font-mono"
              />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Info */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-sm font-medium flex items-center gap-2">
            <Info className="h-4 w-4 text-muted-foreground" />
            Tổng quan cài đặt
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
            {Object.entries(settings).map(([key, value]) => (
              <div key={key} className="flex items-center justify-between px-3 py-2 rounded-md bg-muted/50">
                <span className="font-mono text-xs text-muted-foreground">{key}</span>
                <span className="text-xs truncate max-w-[180px] font-medium">{value}</span>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Actions */}
      <Separator />
      <div className="flex items-center gap-2">
        <Button onClick={handleSaveAll} disabled={saving} size="sm" className="gap-1.5">
          <Save className="h-3.5 w-3.5" />
          {saving ? "Đang lưu..." : "Lưu tất cả"}
        </Button>
        <Button variant="outline" onClick={fetchSettings} size="sm" className="gap-1.5">
          <RotateCcw className="h-3.5 w-3.5" /> Tải lại
        </Button>
      </div>
    </div>
  )
}
