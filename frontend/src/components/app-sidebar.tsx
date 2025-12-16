
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
import { Logo } from "./logo-image";
import { Outlet } from "@tanstack/react-router";
import { Navbar } from "./navbar";
import { UserDropdown } from "./user-dropdown";
import { CreateVaultDialog } from "./create-vault-dialog";
import { BarChart } from "@/assets/icons/bar-chart-icon";
import { SettingsIcon } from "@/assets/icons/settings-icon";
import { SidebarBottomItems } from "./bottom-sidebar";
import { PhoneIcon } from "@/assets/icons/phone-icon";
import { VaultList } from "./vault-list";



const bottomSidebarItems = [
  {
    title: "Insights",
    url: "/admin/insights",
    icon: BarChart,
  },
  {
    title: "Settings",
    url: "/admin/settings",
    icon: SettingsIcon,
  },
  {
    title: "Get on other devices",
    url: "/admin/get-on-other-devices",
    icon: PhoneIcon,
  }
];

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
          <VaultList />
        </SidebarContent>
        <SidebarFooter>
          <SidebarBottomItems items={bottomSidebarItems} />
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