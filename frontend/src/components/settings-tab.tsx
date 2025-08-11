import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Separator } from "@/components/ui/separator"
import { Plus, FolderOpen, Settings, Play, Square } from "lucide-react"

export function SettingsTab() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 p-4">
      <div className="grid gap-4 p-4 border rounded-lg bg-gray-50 dark:bg-gray-700 shadow-sm">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-50">Cài đặt chung</h3>
        <div className="grid gap-2">
          <Label htmlFor="type" className="text-gray-700 dark:text-gray-300">
            Loại
          </Label>
          <Select defaultValue="Youtube">
            <SelectTrigger
              id="type"
              className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
            >
              <SelectValue placeholder="Chọn loại" />
            </SelectTrigger>
            <SelectContent className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
              <SelectItem value="Youtube">Youtube</SelectItem>
              <SelectItem value="Facebook">Facebook</SelectItem>
              <SelectItem value="TikTok">TikTok</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="grid gap-2">
          <Label htmlFor="total-video" className="text-gray-700 dark:text-gray-300">
            Tổng số video
          </Label>
          <Input
            id="total-video"
            type="number"
            defaultValue="15"
            className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
          />
        </div>
      </div>

      <div className="grid gap-4 p-4 border rounded-lg bg-gray-50 dark:bg-gray-700 shadow-sm">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-50">Đầu vào & Đầu ra</h3>
        <div className="grid gap-2">
          <Label htmlFor="input-path" className="text-gray-700 dark:text-gray-300">
            Đầu vào
          </Label>
          <div className="flex items-center gap-2">
            <Input
              id="input-path"
              defaultValue="C:\Users\dohuu\Desktop\Project\DHB Reup Facebook v2\DHB"
              className="flex-1 bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
            />
            <Button
              variant="outline"
              size="icon"
              className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
            >
              <FolderOpen className="h-4 w-4 text-gray-700 dark:text-gray-300" />
              <span className="sr-only">Browse</span>
            </Button>
            <Button
              variant="outline"
              size="icon"
              className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
            >
              <Plus className="h-4 w-4 text-gray-700 dark:text-gray-300" />
              <span className="sr-only">Add</span>
            </Button>
          </div>
        </div>
        <div className="grid gap-2">
          <Label htmlFor="save-to" className="text-gray-700 dark:text-gray-300">
            Lưu vào
          </Label>
          <div className="flex items-center gap-2">
            <Input
              id="save-to"
              defaultValue="C:\Users\dohuu\Desktop\Project\DHB Reup Facebook v2\DHB"
              className="flex-1 bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
            />
            <Button
              variant="outline"
              size="icon"
              className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
            >
              <FolderOpen className="h-4 w-4 text-gray-700 dark:text-gray-300" />
              <span className="sr-only">Browse</span>
            </Button>
          </div>
        </div>
      </div>

      <div className="grid gap-4 p-4 border rounded-lg bg-gray-50 dark:bg-gray-700 shadow-sm">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-50">Xử lý & Render</h3>
        <div className="grid gap-2">
          <Label htmlFor="process" className="text-gray-700 dark:text-gray-300">
            Quá trình
          </Label>
          <Select defaultValue="Download -> Render -> Upload">
            <SelectTrigger
              id="process"
              className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
            >
              <SelectValue placeholder="Chọn quá trình" />
            </SelectTrigger>
            <SelectContent className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
              <SelectItem value="Download -> Render -> Upload">
                Download -{">"} Render -{">"} Upload
              </SelectItem>
              <SelectItem value="Render -> Upload">Render -{">"} Upload</SelectItem>
              <SelectItem value="Download Only">Download Only</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="grid gap-2">
          <Label htmlFor="render" className="text-gray-700 dark:text-gray-300">
            Render
          </Label>
          <div className="flex items-center gap-2">
            <Select defaultValue="Cut and Change MD5">
              <SelectTrigger
                id="render"
                className="flex-1 bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
              >
                <SelectValue placeholder="Chọn tùy chọn render" />
              </SelectTrigger>
              <SelectContent className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
                <SelectItem value="Cut and Change MD5">Cut and Change MD5</SelectItem>
                <SelectItem value="No Change">No Change</SelectItem>
              </SelectContent>
            </Select>
            <Select defaultValue="Default">
              <SelectTrigger className="w-[100px] bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50">
                <SelectValue placeholder="Mặc định" />
              </SelectTrigger>
              <SelectContent className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
                <SelectItem value="Default">Default</SelectItem>
                <SelectItem value="Custom">Custom</SelectItem>
              </SelectContent>
            </Select>
            <Button
              variant="outline"
              size="icon"
              className="bg-white dark:bg-gray-800 border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-200"
            >
              <Settings className="h-4 w-4 text-gray-700 dark:text-gray-300" />
              <span className="sr-only">Render Settings</span>
            </Button>
          </div>
        </div>
        <Separator className="my-2 bg-gray-300 dark:bg-gray-600" />
        <div className="flex items-center gap-2">
          <Button className="flex-1 bg-green-500 hover:bg-green-600 text-white transition-colors duration-200">
            <Play className="h-4 w-4 mr-2" /> Bắt đầu
          </Button>
          <Button
            variant="outline"
            className="flex-1 bg-red-500 hover:bg-red-600 text-white transition-colors duration-200"
          >
            <Square className="h-4 w-4 mr-2" /> Dừng
          </Button>
        </div>
      </div>
    </div>
  )
}
