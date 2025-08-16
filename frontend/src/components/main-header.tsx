import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { Separator } from "@/components/ui/separator"
import {
  FolderOpen,
  Save,
  Link,
  UserPlus,
  Plus,
  Settings,
  Clock,
  List,
  ChevronDown,
  FileText,
  HelpCircle,
  Home,
  Edit,
  Menu,
  LogOut,
  MonitorDot,
} from "lucide-react"
import { ModeToggle } from "@/components/ui/mode-toggle"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "./ui/select"

interface MainHeaderProps {
  toggleTheme: () => void
  isDarkMode: boolean
}

export function MainHeader(_props: MainHeaderProps) {
  return (
    <header className="bg-gray-200 dark:bg-gray-800 border-b border-gray-300 dark:border-gray-700 p-2 flex flex-col gap-2">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
              >
                Trang chủ <ChevronDown className="ml-1 h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem>
                <Home className="mr-2 h-4 w-4" /> Trang chủ
              </DropdownMenuItem>
              {/* <DropdownMenuItem>
                <MonitorDot className="mr-2 h-4 w-4" /> Dashboard
              </DropdownMenuItem> */}
              <DropdownMenuItem>
                <LogOut className="mr-2 h-4 w-4" /> Đăng xuất
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
              >
                Editor <ChevronDown className="ml-1 h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              {/* <DropdownMenuItem>
                <Edit className="mr-2 h-4 w-4" /> Chỉnh sửa
              </DropdownMenuItem> */}
              <DropdownMenuItem>
                <FileText className="mr-2 h-4 w-4" /> Xem log
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                className="text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
              >
                Help <ChevronDown className="ml-1 h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem>
                <HelpCircle className="mr-2 h-4 w-4" /> Trợ giúp
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
        <div className="flex items-center gap-2">
          <ModeToggle />
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
              >
                <Menu className="h-5 w-5" />
                <span className="sr-only">Menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem>Log Error</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <div className="flex items-center gap-4">
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <FolderOpen className="h-5 w-5" />
            <span className="text-xs">Import</span>
          </Button>
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <Save className="h-5 w-5" />
            <span className="text-xs">Export</span>
          </Button>
          {/* <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <Link className="h-5 w-5" />
            <span className="text-xs">Export Link</span>
          </Button> */}
        </div>
        <Separator orientation="vertical" className="h-10 bg-gray-300 dark:bg-gray-700" />
        <div className="flex items-center gap-2">
          {/* <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <UserPlus className="h-5 w-5" />
            <span className="text-xs">Add Account</span>
          </Button> */}
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <Plus className="h-5 w-5" />
            <span className="text-xs">Thêm chiến dịch</span>
          </Button>
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <Settings className="h-5 w-5" />
            <span className="text-xs">Cài đặt API</span>
          </Button>
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <Settings className="h-5 w-5" />
            <span className="text-xs">Cài đặt Render</span>
          </Button>
        </div>
        <Separator orientation="vertical" className="h-10 bg-gray-300 dark:bg-gray-700" />
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <Clock className="h-5 w-5 text-gray-700 dark:text-gray-300" />
            <span className="text-sm text-gray-700 dark:text-gray-300">Thời gian chạy</span>
            {/* <Input
              type="text"
              defaultValue="Mỗi ngày"
              className="w-24 h-8 text-center bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
            /> */}
            <Select>
              <SelectTrigger className="w-30 h-8 text-center bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50">
                <SelectValue placeholder="Select time" />
              </SelectTrigger>
              <SelectContent className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
                <SelectItem value="1">Every Day</SelectItem>
                <SelectItem value="7">Every Week</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="flex items-center gap-2">
            <Clock className="h-5 w-5 text-gray-700 dark:text-gray-300" />
            <span className="text-sm text-gray-700 dark:text-gray-300">Chạy tại thời gian</span>
            {/* <Input
              type="text"
              defaultValue="Mỗi ngày"
              className="w-24 h-8 text-center bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
            /> */}
            <Select>
              <SelectTrigger className="w-30 h-8 text-center bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50">
                <SelectValue placeholder="Select time" />
              </SelectTrigger>
              <SelectContent className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-50">
                <SelectItem value="00:00:00">00:00:00</SelectItem>
                <SelectItem value="01:00:00">01:00:00</SelectItem>
                <SelectItem value="02:00:00">02:00:00</SelectItem>
                <SelectItem value="03:00:00">03:00:00</SelectItem>
                <SelectItem value="04:00:00">04:00:00</SelectItem>
                <SelectItem value="05:00:00">05:00:00</SelectItem>
                <SelectItem value="06:00:00">06:00:00</SelectItem>
                <SelectItem value="07:00:00">07:00:00</SelectItem>
                <SelectItem value="08:00:00">08:00:00</SelectItem>
                <SelectItem value="09:00:00">09:00:00</SelectItem>
                <SelectItem value="10:00:00">10:00:00</SelectItem>
                <SelectItem value="11:00:00">11:00:00</SelectItem>
                <SelectItem value="12:00:00">12:00:00</SelectItem>
                <SelectItem value="13:00:00">13:00:00</SelectItem>
                <SelectItem value="14:00:00">14:00:00</SelectItem>
                <SelectItem value="15:00:00">15:00:00</SelectItem>
                <SelectItem value="16:00:00">16:00:00</SelectItem>
                <SelectItem value="17:00:00">17:00:00</SelectItem>
                <SelectItem value="18:00:00">18:00:00</SelectItem>
                <SelectItem value="19:00:00">19:00:00</SelectItem>
                <SelectItem value="20:00:00">20:00:00</SelectItem>
                <SelectItem value="21:00:00">21:00:00</SelectItem>
                <SelectItem value="22:00:00">22:00:00</SelectItem>
                <SelectItem value="23:00:00">23:00:00</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </div>
    </header>
  )
}
