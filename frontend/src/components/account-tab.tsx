"use client"

import { useState } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Edit, Trash2, ChevronUp, ChevronDown, Plus, X, Check } from "lucide-react"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"

interface Account {
  id: string
  type: string
  name: string
}

export function AccountTab() {
  const [accounts, setAccounts] = useState<Account[]>([
    { id: "1", type: "Profile", name: "Ben Hữu Đỗ" },
    { id: "2", type: "Page", name: "Foot Center" },
    { id: "3", type: "Page", name: "Ẩm Thực Thế Giới" },
    { id: "4", type: "Page", name: "DHB Tools" },
    { id: "5", type: "Page", name: "Adbreaks Marketing" },
    { id: "6", type: "Page", name: "Cộng đồng Pokémon Go" },
  ])
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [currentAccount, setCurrentAccount] = useState<Account | null>(null)
  const [editName, setEditName] = useState("")
  const [editType, setEditType] = useState("")

  const handleEdit = (account: Account) => {
    setCurrentAccount(account)
    setEditName(account.name)
    setEditType(account.type)
    setIsEditDialogOpen(true)
  }

  const handleDelete = (id: string) => {
    setAccounts(accounts.filter((account) => account.id !== id))
  }

  const handleSaveEdit = () => {
    if (currentAccount) {
      setAccounts(
        accounts.map((account) =>
          account.id === currentAccount.id ? { ...account, name: editName, type: editType } : account,
        ),
      )
      setIsEditDialogOpen(false)
      setCurrentAccount(null)
    }
  }

  const handleMove = (id: string, direction: "up" | "down") => {
    const index = accounts.findIndex((acc) => acc.id === id)
    if (index === -1) return

    const newAccounts = [...accounts]
    if (direction === "up" && index > 0) {
      ;[newAccounts[index - 1], newAccounts[index]] = [newAccounts[index], newAccounts[index - 1]]
    } else if (direction === "down" && index < newAccounts.length - 1) {
      ;[newAccounts[index + 1], newAccounts[index]] = [newAccounts[index], newAccounts[index + 1]]
    }
    setAccounts(newAccounts)
  }

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center gap-2 mb-4 p-2 border rounded-md bg-gray-50 dark:bg-gray-700">
        <Select defaultValue="Ben Hữu Đỗ (dohuubenbmt@...)">
          <SelectTrigger className="w-[200px] bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50">
            <SelectValue placeholder="Select an account" />
          </SelectTrigger>
          <SelectContent className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
            <SelectItem value="Ben Hữu Đỗ (dohuubenbmt@...)">Ben Hữu Đỗ (dohuubenbmt@...)</SelectItem>
            <SelectItem value="Account 2">Account 2</SelectItem>
          </SelectContent>
        </Select>
        <Button
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
        </Button>
        <div className="flex items-center gap-1 ml-auto">
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
        </div>
        <Button
          variant="ghost"
          className="flex items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
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
              <TableHead className="w-[80px] text-center text-gray-700 dark:text-gray-300">Add</TableHead>
              <TableHead className="w-[80px] text-center text-gray-700 dark:text-gray-300">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {accounts.map((account) => (
              <TableRow
                key={account.id}
                className="hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors duration-200"
              >
                <TableCell className="font-medium text-gray-800 dark:text-gray-200">{account.type}</TableCell>
                <TableCell className="text-gray-800 dark:text-gray-200">{account.name}</TableCell>
                <TableCell className="text-center">
                  <Button
                    variant="ghost"
                    size="icon"
                    className="hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                  >
                    <Plus className="h-4 w-4 text-blue-500" />
                    <span className="sr-only">Add</span>
                  </Button>
                </TableCell>
                <TableCell className="text-center flex items-center justify-center gap-1">
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleEdit(account)}
                    className="hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                  >
                    <Edit className="h-4 w-4 text-blue-500" />
                    <span className="sr-only">Edit</span>
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDelete(account.id)}
                    className="hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                  >
                    <Trash2 className="h-4 w-4 text-red-500" />
                    <span className="sr-only">Delete</span>
                  </Button>
                  <div className="flex flex-col gap-0.5">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleMove(account.id, "up")}
                      className="h-6 w-6 p-0 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                    >
                      <ChevronUp className="h-3 w-3" />
                      <span className="sr-only">Move Up</span>
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleMove(account.id, "down")}
                      className="h-6 w-6 p-0 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
                    >
                      <ChevronDown className="h-3 w-3" />
                      <span className="sr-only">Move Down</span>
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-[425px] bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
          <DialogHeader>
            <DialogTitle className="text-gray-900 dark:text-gray-50">Chỉnh sửa Account</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="name" className="text-right text-gray-700 dark:text-gray-300">
                Tên
              </Label>
              <Input
                id="name"
                value={editName}
                onChange={(e) => setEditName(e.target.value)}
                className="col-span-3 bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="type" className="text-right text-gray-700 dark:text-gray-300">
                Loại
              </Label>
              <Select value={editType} onValueChange={setEditType}>
                <SelectTrigger className="col-span-3 bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50">
                  <SelectValue placeholder="Chọn loại" />
                </SelectTrigger>
                <SelectContent className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
                  <SelectItem value="Profile">Profile</SelectItem>
                  <SelectItem value="Page">Page</SelectItem>
                  <SelectItem value="Group">Group</SelectItem>
                </SelectContent>
              </Select>
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
    </div>
  )
}
