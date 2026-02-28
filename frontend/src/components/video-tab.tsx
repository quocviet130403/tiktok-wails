"use client"

import { useEffect, useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip"
import { Skeleton } from "@/components/ui/skeleton"
import { Trash2, ChevronLeft, ChevronRight, RefreshCw, Heart, Clock } from "lucide-react"
import { GetAllVideos } from "../../wailsjs/go/backend/App"

interface Video {
  id: number
  title: string
  video_url: string
  thumbnail_url: string
  duration: number
  like_count: number
  profile_douyin_id: number
  status: string
}

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${m}:${s.toString().padStart(2, "0")}`
}

export function VideoTab() {
  const [videos, setVideos] = useState<Video[]>([])
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const pageSize = 20

  const fetchVideos = async () => {
    setLoading(true)
    try {
      const result = await GetAllVideos(page, pageSize)
      setVideos(result || [])
    } catch (err) {
      console.error("Lỗi khi tải videos:", err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchVideos() }, [page])

  const handleDelete = async (id: number) => {
    if (!confirm("Bạn có chắc muốn xóa video này?")) return
    setVideos(videos.filter((v) => v.id !== id))
  }

  const statusBadge = (status: string) => {
    if (status === "done") {
      return <Badge variant="secondary" className="bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-0 text-xs">Hoàn thành</Badge>
    }
    return <Badge variant="secondary" className="bg-amber-500/10 text-amber-600 dark:text-amber-400 border-0 text-xs">Chờ xử lý</Badge>
  }

  return (
    <div className="flex flex-col h-full gap-4">
      {/* Toolbar */}
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Trang {page} • {videos.length} video
        </p>
        <Button variant="outline" size="sm" onClick={fetchVideos} disabled={loading} className="gap-1.5">
          <RefreshCw className={`h-3.5 w-3.5 ${loading ? "animate-spin" : ""}`} />
          Tải lại
        </Button>
      </div>

      {/* Table */}
      <Card className="flex-1 overflow-hidden">
        <CardContent className="p-0 h-full">
          <div className="overflow-auto h-full">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="w-[50px]">#</TableHead>
                  <TableHead>Tiêu đề</TableHead>
                  <TableHead className="w-[90px]">
                    <div className="flex items-center gap-1"><Clock className="h-3 w-3" /> Thời lượng</div>
                  </TableHead>
                  <TableHead className="w-[80px]">
                    <div className="flex items-center gap-1"><Heart className="h-3 w-3" /> Likes</div>
                  </TableHead>
                  <TableHead className="w-[100px]">Trạng thái</TableHead>
                  <TableHead className="w-[50px] text-center">Xóa</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell><Skeleton className="h-4 w-6" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-48" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-10" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-12" /></TableCell>
                      <TableCell><Skeleton className="h-5 w-16 rounded-full" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-4 mx-auto" /></TableCell>
                    </TableRow>
                  ))
                ) : videos.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center py-8 text-muted-foreground">
                      Không có video nào
                    </TableCell>
                  </TableRow>
                ) : (
                  videos.map((video, index) => (
                    <TableRow key={video.id} className="group">
                      <TableCell className="font-mono text-xs text-muted-foreground">
                        {(page - 1) * pageSize + index + 1}
                      </TableCell>
                      <TableCell className="font-medium max-w-[300px] truncate">{video.title}</TableCell>
                      <TableCell className="font-mono text-xs">{formatDuration(video.duration)}</TableCell>
                      <TableCell className="font-mono text-xs">{video.like_count.toLocaleString()}</TableCell>
                      <TableCell>{statusBadge(video.status)}</TableCell>
                      <TableCell className="text-center">
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-7 w-7 opacity-0 group-hover:opacity-100 transition-opacity text-destructive"
                              onClick={() => handleDelete(video.id)}
                            >
                              <Trash2 className="h-3.5 w-3.5" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>Xóa video</TooltipContent>
                        </Tooltip>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      {/* Pagination */}
      <div className="flex items-center justify-between">
        <Button variant="outline" size="sm" onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page <= 1} className="gap-1">
          <ChevronLeft className="h-3.5 w-3.5" /> Trước
        </Button>
        <span className="text-sm text-muted-foreground">Trang {page}</span>
        <Button variant="outline" size="sm" onClick={() => setPage((p) => p + 1)} disabled={videos.length < pageSize} className="gap-1">
          Sau <ChevronRight className="h-3.5 w-3.5" />
        </Button>
      </div>
    </div>
  )
}
