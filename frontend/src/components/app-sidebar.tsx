
import type * as React from "react";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarInset,
  SidebarProvider,
  SidebarRail,
} from "@/components/ui/sidebar";
import { NavMain } from "./nav-main";
import { Logo } from "./logo-image";
import { Outlet } from "@tanstack/react-router";
import { Navbar } from "./navbar";
import { HomeIcon } from "@/assets/icons/home-icon";
import { PaletteIcon } from "@/assets/icons/palette-icon";
import { UserDropdown } from "./user-dropdown";
import { Plus } from "lucide-react";
import { Button } from "./ui/button";
import { CreateVaultDialog } from "./create-vault-dialog";

const data = {
  navMain: [
    {
      title: "Personal",
      url: "/admin/home",
      icon: HomeIcon,
      subtitle: "338 items",
    },
    {
      title: "My Vault",
      url: "/admin/customize",
      icon: PaletteIcon,
      subtitle: "338 items",
    },


  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {

  return (
    <SidebarProvider>
      <Sidebar collapsible="offcanvas" {...props}>
        <SidebarHeader className="px-4">
          <div className="flex gap-1 items-center py-2">
            <Logo />
            <div className="hidden lg:flex">
              <h1 className="text-lg font-bold recoleta-bold ">
                Nucleopass
              </h1>
            </div>
          </div>
        </SidebarHeader>
        <SidebarContent >
          <div className="flex items-center justify-between px-4 ">
            <h2 className="text-sm font-medium text-muted-foreground">
              Vaults
            </h2>
            <CreateVaultDialog />
          </div>
          <NavMain items={data.navMain} />
        </SidebarContent>
        <SidebarFooter>
          <UserDropdown />
        </SidebarFooter>
        <SidebarRail />
      </Sidebar>
      <SidebarInset >
        <Navbar />
        <Outlet />
      </SidebarInset>
    </SidebarProvider>
  );
}