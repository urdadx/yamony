import { ChevronsUpDown, LogOut, Plus } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { RiStackFill } from "@remixicon/react";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { useIsMobile } from "@/hooks/use-mobile";
import { useAuth } from "@/context/auth-context";
import { useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";
import { useMe } from "@/hooks/use-me";

export function UserDropdown() {

  const isMobile = useIsMobile();
  const { logout } = useAuth();

  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    toast.success("Logged out successfully");
    navigate({ to: "/login" });

  }

  const { data: user } = useMe();

  return (
    <SidebarMenu>

      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              {user?.image ? (
                <Avatar className="flex aspect-square border-2 border-primary size-9 items-center justify-center rounded-lg text-sidebar-primary-foreground">
                  <AvatarImage src={user.image} alt="Profile image" />
                  <AvatarFallback>
                    {user?.username ? user.username.charAt(0).toUpperCase() : "U"}
                  </AvatarFallback>
                </Avatar>
              ) : (
                <div className="flex aspect-square size-9 items-center justify-center rounded-lg bg-rose-500 text-white font-semibold">
                  {user?.username ? user.username.charAt(0).toUpperCase() : "U"}
                </div>
              )}
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate ">
                  {user?.username || "User"}
                </span>
                <span className="truncate text-xs">
                  {user?.email || "No email"}
                </span>
              </div>
              <ChevronsUpDown className="ml-auto" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width] min-w-72 rounded-lg"
            align="start"
            side={isMobile ? "bottom" : "top"}
            sideOffset={4}
          >
            <DropdownMenuItem
              className="gap-2 p-2"
            >
              <div className="flex size-6 items-center justify-center rounded-md border bg-background">
                <Plus className="size-4" />
              </div>
              <div className="font-medium text-muted-foreground">
                Create a new profile
              </div>
            </DropdownMenuItem>

            <DropdownMenuItem

              className="gap-2 p-2"
            >
              <div className="flex size-6 items-center justify-center rounded-sm border">
                <Avatar className="h-6 w-6 size-4 shrink-0">
                  <AvatarImage
                    src="https://api.dicebear.com/9.x/glass/svg?seed=Christopher"
                    alt="Profile image"
                  />
                  <AvatarFallback></AvatarFallback>
                </Avatar>
              </div>
              Default
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="gap-2 p-2"
            >
              <div className="flex size-6 items-center justify-center rounded-sm border">
                <RiStackFill className="size-4 shrink-0" />
              </div>
              Switch profiles
            </DropdownMenuItem>

            <DropdownMenuSeparator />

            <DropdownMenuItem
              className="gap-2 p-2 text-red-500 hover:text-red-600"
              onClick={handleLogout}
            >

              <div className="flex size-6 items-center justify-center">
                <LogOut className="size-4 shrink-0" />
              </div>
              Logout
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}