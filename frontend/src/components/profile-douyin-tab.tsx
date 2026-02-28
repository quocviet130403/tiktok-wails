"use client"

import { useEffect, useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip"
import { Switch } from "@/components/ui/switch"
import { Edit, Trash2, Plus, Search } from "lucide-react"
import {
  GetAllDouyinProfiles,
  AddDouyinProfile,
  UpdateDouyinProfile,
  DeleteDouyinProfile,
  ToggleHasTranslate,
} from "../../wailsjs/go/backend/App"

interface ProfileDouyin {
  id: number
  nickname: string
  url: string
  last_video_reup: any
  has_translate: boolean
}

export function ProfileDouyinTab() {
  const [profiles, setProfiles] = useState<ProfileDouyin[]>([])
  const [searchQuery, setSearchQuery] = useState("")

  const fetchProfiles = async () => {
    const result = await GetAllDouyinProfiles()
    if (result) setProfiles(result)
  }

  useEffect(() => { fetchProfiles() }, [])

  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [currentProfile, setCurrentProfile] = useState<ProfileDouyin | null>(null)
  const [editPost, setEditPost] = useState({ nickname: "", url: "" })
  const [createPost, setCreatePost] = useState({ nickname: "", url: "" })

  const filteredProfiles = profiles.filter((p) =>
    p.nickname.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const handleEdit = (profile: ProfileDouyin) => {
    setCurrentProfile(profile)
    setEditPost({ nickname: profile.nickname, url: profile.url })
    setIsEditDialogOpen(true)
  }

  const handleDelete = (id: number) => {
    DeleteDouyinProfile(id).then(() => fetchProfiles()).catch(console.error)
  }

  const handleSaveEdit = () => {
    if (!currentProfile) return
    UpdateDouyinProfile(currentProfile.id, editPost.nickname, editPost.url)
      .then(() => { fetchProfiles(); setIsEditDialogOpen(false) })
      .catch(console.error)
  }

  const handleSaveCreate = () => {
    AddDouyinProfile(createPost.nickname, createPost.url)
      .then(() => {
        fetchProfiles()
        setCreatePost({ nickname: "", url: "" })
        setIsCreateDialogOpen(false)
      })
      .catch(console.error)
  }

  const handleToggleTranslate = (id: number) => {
    ToggleHasTranslate(id).then(() => fetchProfiles()).catch(console.error)
  }

  return (
    <div className="flex flex-col h-full gap-4">
      {/* Toolbar */}
      <div className="flex items-center justify-between gap-3">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Tìm Douyin profile..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9 h-9"
          />
        </div>
        <Button onClick={() => setIsCreateDialogOpen(true)} size="sm" className="gap-1.5">
          <Plus className="h-4 w-4" /> Thêm Profile
        </Button>
      </div>

      {/* Table */}
      <Card className="flex-1 overflow-hidden">
        <CardContent className="p-0 h-full">
          <div className="overflow-auto h-full">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="w-[60px]">ID</TableHead>
                  <TableHead>Nickname</TableHead>
                  <TableHead className="w-[200px]">Video cuối</TableHead>
                  <TableHead className="w-[100px] text-center">Dịch</TableHead>
                  <TableHead className="w-[100px] text-center">Thao tác</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredProfiles.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                      Không có profile nào
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredProfiles.map((profile) => (
                    <TableRow key={profile.id} className="group">
                      <TableCell className="font-mono text-xs text-muted-foreground">{profile.id}</TableCell>
                      <TableCell className="font-medium">{profile.nickname}</TableCell>
                      <TableCell className="text-sm text-muted-foreground">
                        {profile.last_video_reup ?? (
                          <Badge variant="secondary" className="text-xs font-normal">Chưa reup</Badge>
                        )}
                      </TableCell>
                      <TableCell className="text-center">
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <div className="flex items-center justify-center">
                              <Switch
                                checked={profile.has_translate}
                                onCheckedChange={() => handleToggleTranslate(profile.id)}
                              />
                            </div>
                          </TooltipTrigger>
                          <TooltipContent>
                            {profile.has_translate ? "Tắt dịch phụ đề" : "Bật dịch phụ đề"}
                          </TooltipContent>
                        </Tooltip>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center justify-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => handleEdit(profile)}>
                                <Edit className="h-3.5 w-3.5" />
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>Chỉnh sửa</TooltipContent>
                          </Tooltip>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button variant="ghost" size="icon" className="h-7 w-7 text-destructive" onClick={() => handleDelete(profile.id)}>
                                <Trash2 className="h-3.5 w-3.5" />
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>Xóa</TooltipContent>
                          </Tooltip>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Chỉnh sửa Douyin Profile</DialogTitle>
          </DialogHeader>
          <div className="grid gap-3 py-2">
            <div className="grid grid-cols-4 items-center gap-3">
              <Label className="text-right text-sm">Nickname</Label>
              <Input value={editPost.nickname} onChange={(e) => setEditPost({ ...editPost, nickname: e.target.value })} className="col-span-3 h-9" />
            </div>
            <div className="grid grid-cols-4 items-center gap-3">
              <Label className="text-right text-sm">URL</Label>
              <Input value={editPost.url} onChange={(e) => setEditPost({ ...editPost, url: e.target.value })} className="col-span-3 h-9" />
            </div>
          </div>
          <DialogFooter>
            <Button onClick={handleSaveEdit} size="sm">Lưu thay đổi</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Thêm Douyin Profile</DialogTitle>
          </DialogHeader>
          <div className="grid gap-3 py-2">
            <div className="grid grid-cols-4 items-center gap-3">
              <Label className="text-right text-sm">Nickname</Label>
              <Input value={createPost.nickname} onChange={(e) => setCreatePost({ ...createPost, nickname: e.target.value })} className="col-span-3 h-9" />
            </div>
            <div className="grid grid-cols-4 items-center gap-3">
              <Label className="text-right text-sm">URL</Label>
              <Input value={createPost.url} onChange={(e) => setCreatePost({ ...createPost, url: e.target.value })} className="col-span-3 h-9" />
            </div>
          </div>
          <DialogFooter>
            <Button onClick={handleSaveCreate} size="sm">Tạo Profile</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
