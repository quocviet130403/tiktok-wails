"use client"

import { useEffect, useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Edit, Trash2, Play } from "lucide-react"
import { GetAllVideos } from "../../wailsjs/go/backend/App"

interface Video {
  id:           number
	title:        string
	videoURL:     string
	thumbnailURL: string
	duration:     number
	likeCount:    number
	accountID:    number
	status:       string
}

export function VideoTab() {
  const [videos, setVideos] = useState<any[]>([])

  const [pagination, _] = useState({
    page: 1,
    pageSize: 20,
  })

  const fetchVideos = async () => {
    const result = await GetAllVideos(pagination.page, pagination.pageSize)
    if (result) {
      setVideos(result)
    }
  }

  useEffect(() => {
    fetchVideos()
  }, [])

  const handleDelete = (id: number) => {
    setVideos(videos.filter((video) => video.id !== id))
  }

  // In a real app, you'd open a dialog for editing
  const handleEdit = (video: Video) => {
    alert(`Chỉnh sửa video: ${video.title}`)
    // Implement a dialog similar to AccountTab for actual editing
  }

  return (
    <div className="flex flex-col h-full">
      <div className="border rounded-md overflow-auto flex-1 bg-white dark:bg-gray-800">
        <Table className="w-full">
          <TableHeader className="bg-gray-100 dark:bg-gray-700 sticky top-0">
            <TableRow className="hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200">
              <TableHead className="w-[40px] text-center text-gray-700 dark:text-gray-300"></TableHead>
              <TableHead className="text-gray-700 dark:text-gray-300">Title</TableHead>
              <TableHead className="text-gray-700 dark:text-gray-300">Description</TableHead>
              <TableHead className="w-[100px] text-gray-700 dark:text-gray-300">Duration</TableHead>
              <TableHead className="text-gray-700 dark:text-gray-300">Thumbnail</TableHead>
              <TableHead className="w-[100px] text-gray-700 dark:text-gray-300">Schedule</TableHead>
              <TableHead className="w-[100px] text-gray-700 dark:text-gray-300">Account</TableHead>
              <TableHead className="w-[100px] text-gray-700 dark:text-gray-300">Status</TableHead>
              <TableHead className="w-[80px] text-gray-700 dark:text-gray-300">Download</TableHead>
              <TableHead className="w-[80px] text-gray-700 dark:text-gray-300">Render</TableHead>
              <TableHead className="w-[80px] text-gray-700 dark:text-gray-300">Upload</TableHead>
              <TableHead className="w-[60px] text-center text-gray-700 dark:text-gray-300">Edit</TableHead>
              <TableHead className="w-[60px] text-center text-gray-700 dark:text-gray-300">Delete</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {videos.map((video, index) => (
              <TableRow
                key={video.id}
                className="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors duration-200"
              >
                <TableCell className="text-center text-gray-800 dark:text-gray-200">
                  <Play className="h-4 w-4 inline-block text-blue-500" />
                  <span className="ml-1">{index + 1}</span>
                </TableCell>
                <TableCell className="font-medium text-gray-800 dark:text-gray-200">{video.title}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.duration}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.thumbnailURL}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.accountID}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.likeCount}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.status}</TableCell>
                <TableCell className="text-center">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleEdit(video)}
                    className="hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                  >
                    <Edit className="h-4 w-4 text-blue-500" />
                    <span className="sr-only">Edit</span>
                  </Button>
                </TableCell>
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
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
