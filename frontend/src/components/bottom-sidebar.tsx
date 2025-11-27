import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Link, useLocation } from "@tanstack/react-router";

export function SidebarBottomItems({
  items,
}: {
  items: {
    title: string;
    url: string;
    icon?: any;
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
    <SidebarGroup>
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
                        ? "bg-rose-100! dark:bg-rose-900/30!"
                        : ""
                    }
                  >
                    {item.icon && (
                      <item.icon
                        color="currentColor"
                        className={
                          isMainItemActive
                            ? "text-rose-600 dark:text-rose-400"
                            : ""
                        }
                      />
                    )}
                    <span
                      className={
                        isMainItemActive
                          ? "text-rose-600 dark:text-rose-400 font-medium"
                          : ""
                      }
                    >
                      {item.title}
                    </span>
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
                            ? "bg-rose-100 dark:bg-rose-900/30"
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