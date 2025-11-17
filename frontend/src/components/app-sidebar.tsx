
import {
  Bot,
  ChartNoAxesColumnIncreasing,
  Link,
  Settings2,
  Store,
} from "lucide-react";
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
import { UpgradeBanner } from "./upgrade-banner";
import { NavMain } from "./nav-main";
import { Logo } from "./logo-image";
import { Outlet } from "@tanstack/react-router";
import { Navbar } from "./navbar";
import { ChooseBlock } from "./choose-block";

const data = {
  navMain: [
    {
      title: "Home",
      url: "/admin/home",
      icon: Link,
    },
    {
      title: "Customize",
      url: "/admin/customize",
      icon: Store,
    },

    {
      title: "Insights",
      url: "/admin/insights",
      icon: ChartNoAxesColumnIncreasing,
    },
    {
      title: "My Chatbot",
      url: "/admin/chatbot",
      icon: Bot,
    },
    {
      title: "Settings",
      url: "/admin/settings",
      icon: Settings2,
    },

  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {

  return (
    <SidebarProvider>
      <Sidebar variant="floating" collapsible="offcanvas" {...props}>
        <SidebarHeader>
          <div className="flex gap-2 items-center pb-4">
            <Logo />
            <div className="hidden lg:flex">
              <h1 className="text-xl font-bold instrument-serif-regular-italic">
                Yamony
              </h1>
            </div>
          </div>
        </SidebarHeader>
        <SidebarContent >
          <NavMain items={data.navMain} />
        </SidebarContent>
        <SidebarFooter>
          <UpgradeBanner />
          {/* <PageSwitcher /> */}
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