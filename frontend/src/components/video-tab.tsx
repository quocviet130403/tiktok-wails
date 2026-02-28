"use client"

import { useEffect, useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Trash2, ChevronLeft, ChevronRight, RefreshCw } from "lucide-react"
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
  const [loading, setLoading] = useState(false)
  const pageSize = 20

  const fetchVideos = async () => {
    setLoading(true)
    try {
      const result = await GetAllVideos(page, pageSize)
      if (result) {
        setVideos(result)
      } else {
        setVideos([])
      }
    } catch (err) {
      console.error("Lỗi khi tải videos:", err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchVideos()
  }, [page])

  const handleDelete = async (id: number) => {
    if (!confirm("Bạn có chắc muốn xóa video này?")) return
    // TODO: Call backend DeleteVideo when exposed via app.go
    setVideos(videos.filter((video) => video.id !== id))
  }

  const statusBadge = (status: string) => {
    const color = status === "done"
      ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300"
      : "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300"
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium ${color}`}>
        {status}
      </span>
    )
  }

  return (
    <div className="flex flex-col h-full gap-3">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">
          Trang {page} • {videos.length} video
        </h3>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={fetchVideos}
            disabled={loading}
            className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600"
          >
            <RefreshCw className={`h-4 w-4 mr-1 ${loading ? "animate-spin" : ""}`} />
            Tải lại
          </Button>
        </div>
      </div>

      {/* Table */}
      <div className="border rounded-md overflow-auto flex-1 bg-white dark:bg-gray-800">
        <Table className="w-full">
          <TableHeader className="bg-gray-100 dark:bg-gray-700 sticky top-0">
            <TableRow className="hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200">
              <TableHead className="w-[40px] text-center text-gray-700 dark:text-gray-300">#</TableHead>
              <TableHead className="text-gray-700 dark:text-gray-300">Tiêu đề</TableHead>
              <TableHead className="w-[80px] text-gray-700 dark:text-gray-300">Thời lượng</TableHead>
              <TableHead className="w-[80px] text-gray-700 dark:text-gray-300">Lượt thích</TableHead>
              <TableHead className="w-[80px] text-gray-700 dark:text-gray-300">Trạng thái</TableHead>
              <TableHead className="w-[60px] text-center text-gray-700 dark:text-gray-300">Xóa</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {videos.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center py-8 text-gray-500">
                  {loading ? "Đang tải..." : "Không có video nào"}
                </TableCell>
              </TableRow>
            ) : (
              videos.map((video, index) => (
                <TableRow
                  key={video.id}
                  className="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors duration-200"
                >
                  <TableCell className="text-center text-gray-800 dark:text-gray-200">
                    {(page - 1) * pageSize + index + 1}
                  </TableCell>
                  <TableCell className="font-medium text-gray-800 dark:text-gray-200 max-w-[300px] truncate">
                    {video.title}
                  </TableCell>
                  <TableCell className="text-gray-800 dark:text-gray-200">
                    {formatDuration(video.duration)}
                  </TableCell>
                  <TableCell className="text-gray-800 dark:text-gray-200">
                    {video.like_count.toLocaleString()}
                  </TableCell>
                  <TableCell>{statusBadge(video.status)}</TableCell>
                  <TableCell className="text-center">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleDelete(video.id)}
                      className="hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                    >
                      <Trash2 className="h-4 w-4 text-red-500" />
                      <span className="sr-only">Delete</span>
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-between pt-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => setPage((p) => Math.max(1, p - 1))}
          disabled={page <= 1}
          className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600"
        >
          <ChevronLeft className="h-4 w-4 mr-1" />
          Trước
        </Button>
        <span className="text-sm text-gray-500 dark:text-gray-400">Trang {page}</span>
        <Button
          variant="outline"
          size="sm"
          onClick={() => setPage((p) => p + 1)}
          disabled={videos.length < pageSize}
          className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600"
        >
          Sau
          <ChevronRight className="h-4 w-4 ml-1" />
        </Button>
      </div>
    </div>
  )
}
