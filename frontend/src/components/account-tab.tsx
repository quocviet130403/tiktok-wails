"use client"

import { useEffect, useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Label } from "@/components/ui/label"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Edit, Trash2, Plus, Link, Search, ShieldCheck, ShieldX } from "lucide-react"
import {
  GetAllProfiles,
  AddProfile,
  UpdateProfile,
  DeleteProfile,
  GetAllDouyinProfiles,
  GetAllDouyinProfilesFromProfile,
  ConnectWithProfileDouyin,
} from "../../wailsjs/go/backend/App"

interface Profile {
  id: number
  name: string
  hashtag: string
  first_comment: string
  is_authenticated: boolean
  proxy_ip: string
  proxy_port: string
}

export function ProfileTab() {
  const [profiles, setProfiles] = useState<Profile[]>([])
  const [searchQuery, setSearchQuery] = useState("")
  const [profileDouyins, setProfileDouyins] = useState<any[]>([])

  const fetchProfiles = async () => {
    const result = await GetAllProfiles()
    if (result) setProfiles(result)
  }

  const fetchProfileDouyins = async () => {
    const result = await GetAllDouyinProfiles()
    if (result) setProfileDouyins(result)
  }

  useEffect(() => {
    fetchProfiles()
    fetchProfileDouyins()
  }, [])

  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false)
  const [isLinkDialogOpen, setIsLinkDialogOpen] = useState(false)
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)
  const [deleteProfileId, setDeleteProfileId] = useState<number | null>(null)
  const [currentProfile, setCurrentProfile] = useState<Profile | null>(null)
  const [editPost, setEditPost] = useState({ name: "", hashtag: "", firstComment: "", proxyIp: "", proxyPort: "" })
  const [createPost, setCreatePost] = useState({ name: "", hashtag: "", firstComment: "", proxyIp: "", proxyPort: "" })
  const [selectedProfileDouyin, setSelectedProfileDouyin] = useState<any[]>([])

  const filteredProfiles = profiles.filter((p) =>
    p.name.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const handleEdit = (profile: Profile) => {
    setCurrentProfile(profile)
    setEditPost({
      name: profile.name,
      hashtag: profile.hashtag,
      firstComment: profile.first_comment,
      proxyIp: profile.proxy_ip,
      proxyPort: profile.proxy_port,
    })
    setIsEditDialogOpen(true)
  }

  const handleDelete = (id: number) => {
    setDeleteProfileId(id)
    setIsDeleteDialogOpen(true)
  }

  const handleConfirmDelete = () => {
    if (deleteProfileId === null) return
    DeleteProfile(deleteProfileId)
      .then(() => fetchProfiles())
      .catch(console.error)
      .finally(() => {
        setIsDeleteDialogOpen(false)
        setDeleteProfileId(null)
      })
  }

  const handleLink = async (profile: Profile) => {
    const linked = await GetAllDouyinProfilesFromProfile(profile.id)
    setSelectedProfileDouyin(linked || [])
    setCurrentProfile(profile)
    setIsLinkDialogOpen(true)
  }

  const handleSaveEdit = () => {
    if (!currentProfile) return
    UpdateProfile(currentProfile.id, editPost.name, editPost.hashtag, editPost.firstComment, editPost.proxyIp, editPost.proxyPort)
      .then(() => { fetchProfiles(); setIsEditDialogOpen(false) })
      .catch(console.error)
  }

  const handleSaveCreate = () => {
    AddProfile(createPost.name, createPost.hashtag, createPost.firstComment, createPost.proxyIp, createPost.proxyPort)
      .then(() => {
        fetchProfiles()
        setCreatePost({ name: "", hashtag: "", firstComment: "", proxyIp: "", proxyPort: "" })
        setIsCreateDialogOpen(false)
      })
      .catch(console.error)
  }

  const handleLinkProfile = () => {
    if (!currentProfile) return
    ConnectWithProfileDouyin(currentProfile.id, selectedProfileDouyin.map((p) => p.id))
      .then(() => { setIsLinkDialogOpen(false); setSelectedProfileDouyin([]) })
      .catch(console.error)
  }

  return (
    <div className="flex flex-col h-full gap-4">
      {/* Toolbar */}
      <div className="flex items-center justify-between gap-3">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Tìm profile..."
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
                  <TableHead>Tên</TableHead>
                  <TableHead className="w-[120px]">Xác thực</TableHead>
                  <TableHead className="w-[160px]">Proxy</TableHead>
                  <TableHead className="w-[120px] text-center">Thao tác</TableHead>
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
                      <TableCell className="font-medium">{profile.name}</TableCell>
                      <TableCell>
                        {profile.is_authenticated ? (
                          <Badge variant="secondary" className="gap-1 bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-0">
                            <ShieldCheck className="h-3 w-3" /> Đã xác thực
                          </Badge>
                        ) : (
                          <Badge variant="secondary" className="gap-1 bg-red-500/10 text-red-600 dark:text-red-400 border-0">
                            <ShieldX className="h-3 w-3" /> Chưa
                          </Badge>
                        )}
                      </TableCell>
                      <TableCell className="font-mono text-xs">
                        {profile.proxy_ip ? `${profile.proxy_ip}:${profile.proxy_port}` : (
                          <span className="text-muted-foreground">—</span>
                        )}
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
                              <Button variant="ghost" size="icon" className="h-7 w-7" onClick={() => handleLink(profile)}>
                                <Link className="h-3.5 w-3.5" />
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>Liên kết Douyin</TooltipContent>
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
            <DialogTitle>Chỉnh sửa Profile</DialogTitle>
          </DialogHeader>
          <div className="grid gap-3 py-2">
            {[
              { id: "name", label: "Tên", value: editPost.name, key: "name" as const },
              { id: "hashtag", label: "Hashtag", value: editPost.hashtag, key: "hashtag" as const },
              { id: "first_comment", label: "Comment đầu", value: editPost.firstComment, key: "firstComment" as const },
              { id: "proxy_ip", label: "Proxy IP", value: editPost.proxyIp, key: "proxyIp" as const },
              { id: "proxy_port", label: "Proxy Port", value: editPost.proxyPort, key: "proxyPort" as const },
            ].map((field) => (
              <div key={field.id} className="grid grid-cols-4 items-center gap-3">
                <Label htmlFor={field.id} className="text-right text-sm">{field.label}</Label>
                <Input
                  id={field.id}
                  value={field.value}
                  onChange={(e) => setEditPost({ ...editPost, [field.key]: e.target.value })}
                  className="col-span-3 h-9"
                />
              </div>
            ))}
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
            <DialogTitle>Tạo mới Profile</DialogTitle>
          </DialogHeader>
          <div className="grid gap-3 py-2">
            {[
              { id: "c-name", label: "Tên", value: createPost.name, key: "name" as const },
              { id: "c-hashtag", label: "Hashtag", value: createPost.hashtag, key: "hashtag" as const },
              { id: "c-fc", label: "Comment đầu", value: createPost.firstComment, key: "firstComment" as const },
              { id: "c-pip", label: "Proxy IP", value: createPost.proxyIp, key: "proxyIp" as const },
              { id: "c-pport", label: "Proxy Port", value: createPost.proxyPort, key: "proxyPort" as const },
            ].map((field) => (
              <div key={field.id} className="grid grid-cols-4 items-center gap-3">
                <Label htmlFor={field.id} className="text-right text-sm">{field.label}</Label>
                <Input
                  id={field.id}
                  value={field.value}
                  onChange={(e) => setCreatePost({ ...createPost, [field.key]: e.target.value })}
                  className="col-span-3 h-9"
                />
              </div>
            ))}
          </div>
          <DialogFooter>
            <Button onClick={handleSaveCreate} size="sm">Tạo Profile</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Link Dialog */}
      <Dialog open={isLinkDialogOpen} onOpenChange={setIsLinkDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Liên kết Douyin Profile</DialogTitle>
          </DialogHeader>
          <div className="grid gap-3 py-2">
            <div className="grid grid-cols-4 items-center gap-3">
              <Label className="text-right text-sm">Douyin</Label>
              <Select onValueChange={(value) => {
                if (value === "default") return
                const p = profileDouyins.find((pd) => String(pd.id) === value)
                if (p && !selectedProfileDouyin.find((s) => s.id === p.id)) {
                  setSelectedProfileDouyin([...selectedProfileDouyin, p])
                }
              }}>
                <SelectTrigger className="col-span-3 h-9">
                  <SelectValue placeholder="Chọn profile..." />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="default">-- Chọn Profile --</SelectItem>
                  {profileDouyins.map((p) => (
                    <SelectItem key={p.id} value={String(p.id)}>{p.nickname}</SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="space-y-2 max-h-40 overflow-auto">
              {selectedProfileDouyin.length > 0 ? (
                selectedProfileDouyin.map((p) => (
                  <div key={p.id} className="flex items-center justify-between px-3 py-2 rounded-md bg-muted">
                    <span className="text-sm">{p.nickname}</span>
                    <Button variant="ghost" size="icon" className="h-6 w-6"
                      onClick={() => setSelectedProfileDouyin(selectedProfileDouyin.filter((s) => s.id !== p.id))}
                    >
                      <Trash2 className="h-3 w-3 text-destructive" />
                    </Button>
                  </div>
                ))
              ) : (
                <p className="text-sm text-muted-foreground text-center py-4">Chưa có liên kết nào</p>
              )}
            </div>
          </div>
          <DialogFooter>
            <Button onClick={handleLinkProfile} size="sm">Lưu liên kết</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Xác nhận xóa Profile</AlertDialogTitle>
            <AlertDialogDescription>
              Bạn có chắc chắn muốn xóa profile này không? Hành động này không thể hoàn tác.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={() => setDeleteProfileId(null)}>Hủy</AlertDialogCancel>
            <AlertDialogAction onClick={handleConfirmDelete} className="bg-destructive text-destructive-foreground hover:bg-destructive/90">
              Xóa
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
