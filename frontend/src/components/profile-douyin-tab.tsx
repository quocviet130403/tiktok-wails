"use client"

import { useEffect, useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Edit, Trash2, Plus } from "lucide-react"
// import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { GetAllDouyinProfiles, AddDouyinProfile, UpdateDouyinProfile, DeleteDouyinProfile } from "../../wailsjs/go/backend/App"

interface ProfileDouyin {
  id: number
  nickname: string
  url: string
}

export function ProfileDouyinTab() {

  const [profiles, setProfiles] = useState<any[]>([])

  const fetchProfiles = async () => {
    const result = await GetAllDouyinProfiles()
    if (result) {
      setProfiles(result)
    }
  }

  useEffect(() => {
    fetchProfiles()
  }, [])

  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [currentProfile, setCurrentProfile] = useState<ProfileDouyin | null>(null)
  const [editPost, setEditPost] = useState({
    nickname: "",
    url: "",
  })
  const [createPost, setCreatePost] = useState({
    nickname: "",
    url: "",
  })

  const handleEdit = (profile: ProfileDouyin) => {
    setCurrentProfile(profile)
    setEditPost({
      nickname: profile.nickname,
      url: profile.url,
    })
    setIsEditDialogOpen(true)
  }

  const handleDelete = (id: number) => {
    DeleteDouyinProfile(id)
      .then(() => {
        fetchProfiles()
      })
      .catch((error: any) => {
        console.error("Error deleting profile:", error)
      })
  }

  const handleSaveEdit = () => {
    if (currentProfile) {

      UpdateDouyinProfile(currentProfile.id, editPost.nickname, editPost.url)
        .then(() => {
          fetchProfiles()
          setIsEditDialogOpen(false)
          setCurrentProfile(null)
          setEditPost({
            nickname: "",
            url: "",
          })
        })
        .catch((error: any) => {
          console.error("Error updating profile:", error)
        })
    }
  }

  const handleSaveCreate = () => {
    console.log("Creating new profile with data:", createPost)

    AddDouyinProfile(createPost.nickname, createPost.url)
    .then(() => {
      console.log("Profile created successfully")
      fetchProfiles()
      
      setCreatePost({
        nickname: "",
        url: "",
      })
      
      setIsCreateDialogOpen(false)
    })
    .catch((error: any) => {
      console.error("Error creating profile:", error)
    })
  }

  // const handleMove = (id: string, direction: "up" | "down") => {
  //   const index = profiles.findIndex((acc) => acc.id === id)
  //   if (index === -1) return

  //   const newProfiles = [...profiles]
  //   if (direction === "up" && index > 0) {
  //     ;[newProfiles[index - 1], newProfiles[index]] = [newProfiles[index], newProfiles[index - 1]]
  //   } else if (direction === "down" && index < newProfiles.length - 1) {
  //     ;[newProfiles[index + 1], newProfiles[index]] = [newProfiles[index], newProfiles[index + 1]]
  //   }
  //   setProfiles(newProfiles)
  // }

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center gap-2 mb-4 p-2 border rounded-md bg-gray-50 dark:bg-gray-700">
        {/* <Select defaultValue="Ben Hữu Đỗ (dohuubenbmt@...)">
          <SelectTrigger className="w-[200px] bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50">
            <SelectValue placeholder="Select an profile" />
          </SelectTrigger>
          <SelectContent className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
            <SelectItem value="Ben Hữu Đỗ (dohuubenbmt@...)">Ben Hữu Đỗ (dohuubenbmt@...)</SelectItem>
            <SelectItem value="Profile 2">Profile 2</SelectItem>
          </SelectContent>
        </Select> */}
        {/* <Button
          variant="ghost"
          size="icon"
          className="hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
        >
          <X className="h-4 w-4 text-red-500" />
          <span className="sr-only">Delete Selected</span>
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
        >
          <Check className="h-4 w-4 text-green-500" />
          <span className="sr-only">Confirm</span>
        </Button> */}
        {/* <div className="flex items-center gap-1 ml-auto">
          <input
            type="checkbox"
            id="page-checkbox"
            className="form-checkbox h-4 w-4 text-blue-600 rounded focus:ring-blue-500"
          />
          <label htmlFor="page-checkbox" className="text-gray-700 dark:text-gray-300 text-sm">
            Page
          </label>
        </div>
        <div className="flex items-center gap-1">
          <input
            type="checkbox"
            id="group-checkbox"
            className="form-checkbox h-4 w-4 text-blue-600 rounded focus:ring-blue-500"
          />
          <label htmlFor="group-checkbox" className="text-gray-700 dark:text-gray-300 text-sm">
            Group
          </label>
        </div> */}
        <Button
          variant="ghost"
          className="flex items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
          onClick={() => setIsCreateDialogOpen(true)}
        >
          <Plus className="h-4 w-4" /> Add
        </Button>
      </div>

      <div className="border rounded-md overflow-auto flex-1 bg-white dark:bg-gray-800">
        <Table className="w-full">
          <TableHeader className="bg-gray-100 dark:bg-gray-700 sticky top-0">
            <TableRow className="hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200">
              <TableHead className="w-[100px] text-gray-700 dark:text-gray-300">Type</TableHead>
              <TableHead className="text-gray-700 dark:text-gray-300">Name</TableHead>
              <TableHead className="w-[80px] text-center text-gray-700 dark:text-gray-300">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {profiles.map((profile) => (
              <TableRow
                key={profile.id}
                className="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors duration-200"
              >
                <TableCell className="font-medium text-gray-800 dark:text-gray-200">{profile.id}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{profile.name}</TableCell>
                <TableCell className="text-center flex items-center justify-center gap-1">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleEdit(profile)}
                    className="hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                  >
                    <Edit className="h-4 w-4 text-blue-500" />
                    <span className="sr-only">Edit</span>
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDelete(profile.id)}
                    className="hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                  >
                    <Trash2 className="h-4 w-4 text-red-500" />
                    <span className="sr-only">Delete</span>
                  </Button>
                  {/* <div className="flex flex-col gap-0.5">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleMove(profile.id, "up")}
                      className="h-6 w-6 p-0 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                    >
                      <ChevronUp className="h-3 w-3" />
                      <span className="sr-only">Move Up</span>
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleMove(profile.id, "down")}
                      className="h-6 w-6 p-0 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                    >
                      <ChevronDown className="h-3 w-3" />
                      <span className="sr-only">Move Down</span>
                    </Button>
                  </div> */}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-[425px] bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
          <DialogHeader>
            <DialogTitle className="text-gray-900 dark:text-gray-50">Chỉnh sửa Profile</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="name" className="text-right text-gray-700 dark:text-gray-300">
                Tên
              </Label>
              <Input
                id="nickname"
                value={editPost?.nickname}
                onChange={(e) => setEditPost({ ...editPost, nickname: e.target.value })}
                className="col-span-3 bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="url" className="text-right text-gray-700 dark:text-gray-300">
                URL
              </Label>
              <Input
                id="url"
                value={editPost?.url}
                onChange={(e) => setEditPost({ ...editPost, url: e.target.value })}
                className="col-span-3 bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
              />
            </div>
          </div>
          <DialogFooter>
            <Button
              type="submit"
              onClick={handleSaveEdit}
              className="bg-blue-500 hover:bg-blue-600 text-white transition-colors duration-200"
            >
              Lưu thay đổi
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-[425px] bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
          <DialogHeader>
            <DialogTitle className="text-gray-900 dark:text-gray-50">Tạo mới Profile</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="name" className="text-right text-gray-700 dark:text-gray-300">
                Tên
              </Label>
              <Input
                id="nickname"
                value={createPost.nickname}
                onChange={(e) => setCreatePost({ ...createPost, nickname: e.target.value })}
                className="col-span-3 bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="url" className="text-right text-gray-700 dark:text-gray-300">
                URL
              </Label>
              <Input
                id="url"
                value={createPost.url}
                onChange={(e) => setCreatePost({ ...createPost, url: e.target.value })}
                className="col-span-3 bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
              />
            </div>
          </div>
          <DialogFooter>
            <Button
              type="submit"
              onClick={handleSaveCreate}
              className="bg-blue-500 hover:bg-blue-600 text-white transition-colors duration-200"
            >
              Lưu thay đổi
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
