import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuAction,
} from "@/components/ui/sidebar";
import { Link, useLocation } from "@tanstack/react-router";
import { MoreVertical } from "lucide-react";
import { Button } from "./ui/button";
import { cn } from "@/lib/utils";

export function NavMain({
  items,
}: {
  items: {
    title: string;
    url: string;
    icon?: any;
    subtitle?: string;
    items?: {
      title: string;
      url: string;
    }[];
  }[];
}) {
  const pathname = useLocation({
    select: (location) => location.pathname,
  });

  return (
    <SidebarGroup className="pt-0 mt-0">
      <SidebarMenu>
        {items.map((item) => {
          const isMainItemActive = pathname === item.url;

          return (
            <div key={item.title}>
              <SidebarMenuItem>
                <Link to={item.url} className="w-full">
                  <SidebarMenuButton
                    tooltip={item.title}
                    className={
                      isMainItemActive
                        ? "bg-rose-50/80! border border-rose-50 dark:bg-rose-900/30!  active:bg-rose-100! dark:active:bg-rose-900/40!"
                        : "flex items-center justify-between hover:bg-rose-50"
                    }
                  >
                    <div className="flex items-center gap-3 flex-1">
                      <div className={cn(
                        "p-2 rounded-full",
                        isMainItemActive && "bg-white dark:bg-gray-800 border border-rose-50 dark:border-rose-900"
                      )}>
                        {item.icon && (
                          <item.icon
                            color="currentColor"
                            className={`size-5! font-light text-muted-foreground ${isMainItemActive
                              ? "text-rose-600!  dark:text-rose-400!"
                              : ""
                              }`}
                          />
                        )}
                      </div>
                      <div className="flex flex-col gap-0.5">
                        <span

                        >
                          {item.title}
                        </span>
                        {item.subtitle && (
                          <span className="text-sm text-muted-foreground">
                            {item.subtitle}
                          </span>
                        )}
                      </div>
                    </div>
                    <Button
                      size="sm"
                      variant="ghost"
                      className="rounded-full  "
                    >
                      <MoreVertical className="size-4" />
                      <span className="sr-only">More options</span>
                    </Button>
                  </SidebarMenuButton>

                </Link>

                <div className="py-1">
                  {item.items?.map((subItem) => {
                    const isSubItemActive = pathname === subItem.url;

                    return (
                      <SidebarMenuButton
                        key={subItem.title}
                        className={
                          isSubItemActive
                            ? "bg-rose-100! dark:bg-rose-900/30"
                            : ""
                        }
                      >
                        <Link to={subItem.url} className="w-full">
                          <span
                            className={
                              isSubItemActive
                                ? "text-rose-600 dark:text-rose-400 font-medium"
                                : ""
                            }
                          >
                            {subItem.title}
                          </span>
                        </Link>
                      </SidebarMenuButton>
                    );
                  })}
                </div>
              </SidebarMenuItem>
            </div>
          );
        })}
      </SidebarMenu>
    </SidebarGroup>
  );
}