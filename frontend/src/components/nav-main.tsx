import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Link, useLocation } from "@tanstack/react-router";

export function NavMain({
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
                        ? "bg-green-50 dark:bg-green-900/30 hover:bg-green-100 dark:hover:bg-green-900/40 active:bg-green-100 dark:active:bg-green-900/40"
                        : ""
                    }
                  >
                    {item.icon && (
                      <item.icon
                        color="currentColor"
                        className={`size-4! text-muted-foreground ${isMainItemActive
                          ? "text-green-600  dark:text-green-400"
                          : ""
                          }`}
                      />
                    )}
                    <span
                      className={
                        isMainItemActive
                          ? "text-green-600  dark:text-green-400 font-medium"
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
                            ? "bg-green-100 dark:bg-green-900/30"
                            : ""
                        }
                      >
                        <Link to={subItem.url} className="w-full">
                          <span
                            className={
                              isSubItemActive
                                ? "text-green-600 dark:text-green-400 font-medium"
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