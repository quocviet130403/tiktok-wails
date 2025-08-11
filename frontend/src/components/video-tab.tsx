"use client"

import { useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Edit, Trash2, Play } from "lucide-react"

interface Video {
  id: string
  title: string
  description: string
  duration: string
  thumbnail: string
  schedule: string
  account: string
  status: string
  download: string
  render: string
  upload: string
}

export function VideoTab() {
  const [videos, setVideos] = useState<Video[]>([
    {
      id: "1",
      title: "Jonas Blue and Sabrina Carpenter Live ...",
      description: "Alien with Sabrina Carpe...",
      duration: "00:01:38",
      thumbnail: "https://i.ytimg.com/vi/...",
      schedule: "None",
      account: "",
      status: "Not Start",
      download: "0%",
      render: "0%",
      upload: "0%",
    },
    {
      id: "2",
      title: "ONLY THE BRAVE by DIESEL",
      description: "Join us to meet the brav...",
      duration: "00:00:50",
      thumbnail: "https://i.ytimg.com/vi/...",
      schedule: "None",
      account: "",
      status: "Not Start",
      download: "0%",
      render: "0%",
      upload: "0%",
    },
    {
      id: "3",
      title: "Jonas Blue Explores Hong Kong with A...",
      description: "When American Airlines ...",
      duration: "00:05:23",
      thumbnail: "https://i.ytimg.com/vi/...",
      schedule: "None",
      account: "",
      status: "Not Start",
      download: "0%",
      render: "0%",
      upload: "0%",
    },
    {
      id: "4",
      title: "Jonas Blue's #MarimbaMama",
      description: "Time to get creative guy...",
      duration: "00:03:08",
      thumbnail: "https://i.ytimg.com/vi/...",
      schedule: "None",
      account: "",
      status: "Not Start",
      download: "0%",
      render: "0%",
      upload: "0%",
    },
    {
      id: "5",
      title: "Jonas Blue Live at Heaven London",
      description: "Such a special experienc...",
      duration: "00:02:48",
      thumbnail: "https://i.ytimg.com/vi/...",
      schedule: "None",
      account: "",
      status: "Not Start",
      download: "0%",
      render: "0%",
      upload: "0%",
    },
    {
      id: "6",
      title: "Jonas Blue in France & Germany June 2...",
      description: "Fun times last weekend i...",
      duration: "00:01:42",
      thumbnail: "https://i.ytimg.com/vi/...",
      schedule: "None",
      account: "",
      status: "Not Start",
      download: "0%",
      render: "0%",
      upload: "0%",
    },
    {
      id: "7",
      title: "Jonas Blue at Ministry Of Sound 2017",
      description: "So sick to have my first ...",
      duration: "00:01:22",
      thumbnail: "https://i.ytimg.com/vi/...",
      schedule: "None",
      account: "",
      status: "Not Start",
      download: "0%",
      render: "0%",
      upload: "0%",
    },
    {
      id: "8",
      title: "Jonas Blue at Lollapalooza Berlin 2016",
      description: "Had such an awesome b...",
      duration: "00:00:59",
      thumbnail: "https://i.ytimg.com/vi/...",
      schedule: "None",
      account: "",
      status: "Not Start",
      download: "0%",
      render: "0%",
      upload: "0%",
    },
    {
      id: "9",
      title: "Jonas Blue at V Festival 2016",
      description: "An awesome recap of a ...",
      duration: "00:01:01",
      thumbnail: "https://i.ytimg.com/vi/...",
      schedule: "None",
      account: "",
      status: "Not Start",
      download: "0%",
      render: "0%",
      upload: "0%",
    },
    {
      id: "10",
      title: "Jonas Blue - David Guetta Listen Tour S...",
      description: "A huge moment for me ...",
      duration: "00:01:10",
      thumbnail: "https://i.ytimg.com/vi/...",
      schedule: "None",
      account: "",
      status: "Not Start",
      download: "0%",
      render: "0%",
      upload: "0%",
    },
  ])

  const handleDelete = (id: string) => {
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
                <TableCell className="text-gray-800 dark:text-gray-200">{video.description}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.duration}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.thumbnail}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.schedule}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.account}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.status}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.download}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.render}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{video.upload}</TableCell>
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
