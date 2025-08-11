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
                Home <ChevronDown className="ml-1 h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem>
                <Home className="mr-2 h-4 w-4" /> Trang chủ
              </DropdownMenuItem>
              <DropdownMenuItem>
                <MonitorDot className="mr-2 h-4 w-4" /> Dashboard
              </DropdownMenuItem>
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
              <DropdownMenuItem>
                <Edit className="mr-2 h-4 w-4" /> Chỉnh sửa
              </DropdownMenuItem>
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
            <span className="text-xs">Open</span>
          </Button>
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <Save className="h-5 w-5" />
            <span className="text-xs">Save As...</span>
          </Button>
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <Link className="h-5 w-5" />
            <span className="text-xs">Export Link</span>
          </Button>
        </div>
        <Separator orientation="vertical" className="h-10 bg-gray-300 dark:bg-gray-700" />
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <UserPlus className="h-5 w-5" />
            <span className="text-xs">Add Account</span>
          </Button>
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <Plus className="h-5 w-5" />
            <span className="text-xs">Add Campaign</span>
          </Button>
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <Settings className="h-5 w-5" />
            <span className="text-xs">API Setting</span>
          </Button>
          <Button
            variant="ghost"
            className="flex flex-col items-center gap-1 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-700 transition-colors duration-200"
          >
            <Settings className="h-5 w-5" />
            <span className="text-xs">Render Setting</span>
          </Button>
        </div>
        <Separator orientation="vertical" className="h-10 bg-gray-300 dark:bg-gray-700" />
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2">
            <Clock className="h-5 w-5 text-gray-700 dark:text-gray-300" />
            <span className="text-sm text-gray-700 dark:text-gray-300">Delay</span>
            <Input
              type="text"
              defaultValue="00:00:00"
              className="w-24 h-8 text-center bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
            />
          </div>
          <div className="flex items-center gap-2">
            <List className="h-5 w-5 text-gray-700 dark:text-gray-300" />
            <span className="text-sm text-gray-700 dark:text-gray-300">Thread</span>
            <Input
              type="number"
              defaultValue="1"
              className="w-16 h-8 text-center bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-gray-50"
            />
          </div>
        </div>
      </div>
    </header>
  )
}
