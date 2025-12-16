import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { MoreVertical } from "lucide-react";
import { Button } from "./ui/button";
import { api } from "@/lib/api";
import { useQuery } from "@tanstack/react-query";
import { VAULT_ICONS, type VaultIconName } from "@/lib/vault-icons";
import { Skeleton } from "./ui/skeleton";
import { useNavigate, useSearch } from "@tanstack/react-router";
import { cn } from "@/lib/utils";
import { useEffect } from "react";
import { VaultOptionsPopover } from "./vault-options-popover";

interface Vault {
  id: number;
  user_id: number;
  name: string;
  description?: string;
  icon?: string;
  theme?: string;
  is_favorite: boolean;
  item_count: number;
  created_at: string;
  updated_at: string;
}

export function VaultList() {

  const { data: vaults, isLoading } = useQuery<Vault[]>({
    queryKey: ["vaults"],
    queryFn: async () => {
      const response = await api.get("/vaults");
      return response.data as Vault[];
    },
    staleTime: 5 * 60 * 1000,
  });

  const navigate = useNavigate();
  const { vaultId: selectedVaultId } = useSearch({ from: '/admin' });

  useEffect(() => {
    if (vaults && vaults.length > 0 && !selectedVaultId) {
      navigate({
        to: '.',
        search: { vaultId: vaults[0].id },
        replace: true,
      });
    }
  }, [vaults, selectedVaultId, navigate]);

  if (isLoading) {
    return (
      <SidebarMenu>
        {[1, 2, 3].map((i) => (
          <SidebarMenuItem key={i}>
            <div className="flex items-center gap-3 p-2">
              <Skeleton className="size-10 rounded-full" />
              <div className="flex-1 space-y-2">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-3 w-16" />
              </div>
            </div>
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    );
  }

  if (!vaults || vaults.length === 0) {
    return (
      <div className="px-4 py-6 text-center text-sm text-muted-foreground">
        No vaults yet. Create your first vault!
      </div>
    );
  }

  return (
    <SidebarMenu className="px-2">
      {vaults.map((vault) => {
        const IconComponent = vault.icon ? VAULT_ICONS[vault.icon as VaultIconName] : VAULT_ICONS.Home;
        const iconColor = vault.theme || "#ec4899";
        const isActive = selectedVaultId === vault.id;

        return (
          <SidebarMenuItem key={vault.id}>
            <SidebarMenuButton
              tooltip={vault.name}
              className={cn(
                "flex items-center justify-between px-2 py-7",
                isActive
                  ? "bg-rose-50/80 border border-rose-50 dark:bg-rose-900/30 active:bg-rose-100 dark:active:bg-rose-900/40"
                  : "flex items-center justify-between hover:bg-rose-50 px-2 py-7"
              )}
              onClick={() => {
                navigate({
                  to: '.',
                  search: { vaultId: vault.id },
                });
              }}
            >
              <div className="flex items-center gap-3 flex-1">
                <div className={cn(
                  "p-2 rounded-full",
                  isActive && "bg-white dark:bg-gray-800 border border-rose-50 dark:border-rose-900"
                )}>
                  {IconComponent && (
                    <IconComponent
                      color={iconColor}
                      className={cn(
                        "size-5 font-light",
                        isActive && "text-rose-600 dark:text-rose-400"
                      )}
                    />
                  )}
                </div>
                <div className="flex flex-col gap-0.5">
                  <span>{vault.name}</span>
                  <span className="text-sm text-muted-foreground">
                    {vault.item_count} {vault.item_count === 1 ? 'item' : 'items'}
                  </span>
                </div>
              </div>
              <VaultOptionsPopover
                vault={{
                  ...vault,
                  id: String(vault.id),
                }}

              >
                <Button
                  size="sm"
                  variant="ghost"
                  className="rounded-full"
                  onClick={(e) => {
                    e.stopPropagation();
                  }}
                >
                  <MoreVertical className="size-4" />
                  <span className="sr-only">More options</span>
                </Button>
              </VaultOptionsPopover>
            </SidebarMenuButton>
          </SidebarMenuItem>
        );
      })}
    </SidebarMenu>
  );
}
