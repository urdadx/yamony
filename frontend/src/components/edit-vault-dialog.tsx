import * as React from "react";
import { cn } from "@/lib/utils";
import { Button } from "./ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "./ui/dialog";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
} from "./ui/drawer";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { useIsMobile } from "@/hooks/use-mobile";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { VAULT_ICONS } from "@/lib/vault-icons";
import { toast } from "sonner";

const VAULT_COLORS = [
  { name: "Pink", value: "#ec4899" },
  { name: "Orange", value: "#f97316" },
  { name: "Yellow", value: "#eab308" },
  { name: "Green", value: "#22c55e" },
  { name: "Blue", value: "#3b82f6" },
  { name: "Indigo", value: "#6366f1" },
  { name: "Purple", value: "#a855f7" },
  { name: "Gray", value: "#6b7280" },
  { name: "Red", value: "#ef4444" },
  { name: "Cyan", value: "#06b6d4" },
  { name: "Brown", value: "#a0522d" },
];

const VAULT_ICON_LIST = Object.entries(VAULT_ICONS).map(([name, component]) => ({
  name,
  component,
}));

interface EditVaultDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  vault?: {
    id: string;
    name: string;
    icon?: string;
    theme?: string;
  };
}

export const EditVaultDialog = ({ open, onOpenChange, vault }: EditVaultDialogProps) => {
  const isMobile = useIsMobile();

  if (!vault) {
    return null;
  }

  if (!isMobile) {
    return (
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="rounded-xl border-none bg-clip-padding shadow-2xl ring-4 ring-neutral-200/80 outline-none md:max-w-2xl dark:bg-neutral-800 dark:ring-neutral-900">
          <DialogHeader>
            <DialogTitle>Edit Vault</DialogTitle>
          </DialogHeader>
          <VaultForm key={vault.id} vault={vault} onSuccess={() => onOpenChange(false)} />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Drawer open={open} onOpenChange={onOpenChange}>
      <DrawerContent>
        <DrawerHeader className="text-left">
          <DrawerTitle>Edit Vault</DrawerTitle>
          <DrawerDescription>
            Update your vault details. Change the title, color, and icon.
          </DrawerDescription>
        </DrawerHeader>
        <VaultForm key={vault.id} className="px-4" vault={vault} onSuccess={() => onOpenChange(false)} />
        <DrawerFooter className="pt-2">
          <DrawerClose asChild>
            <Button variant="outline">Cancel</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
};

function VaultForm({
  className,
  vault,
  onSuccess
}: React.ComponentProps<"form"> & {
  vault: {
    id: string;
    name: string;
    icon?: string;
    theme?: string;
  };
  onSuccess?: () => void;
}) {
  const [title, setTitle] = React.useState(vault.name);
  const [selectedColor, setSelectedColor] = React.useState(vault.theme || VAULT_COLORS[0].value);
  const [selectedIcon, setSelectedIcon] = React.useState(vault.icon || "Home");

  const queryClient = useQueryClient();

  const updateVault = useMutation({
    mutationFn: async (data: {
      name: string;
      icon?: string;
      theme?: string;
    }) => {
      const response = await api.put(`/vaults/${vault.id}`, data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaults"] });
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!title.trim()) {
      toast.error("Please enter a vault title");
      return;
    }

    try {
      await updateVault.mutateAsync({
        name: title,
        icon: selectedIcon,
        theme: selectedColor,
      });

      toast.success("Vault updated successfully!");
      onSuccess?.();
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : "Failed to update vault";
      toast.error(errorMessage);
    }
  };

  return (
    <form className={cn("grid items-start gap-6", className)} onSubmit={handleSubmit}>
      <div className="grid gap-3">
        <Label htmlFor="title">Title</Label>
        <Input
          id="title"
          placeholder="Enter vault title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          disabled={updateVault.isPending}
        />
      </div>

      <div className="grid gap-3">
        <Label>Color</Label>
        <div className="flex flex-wrap gap-4">
          {VAULT_COLORS.map((color) => (
            <button
              key={color.value}
              type="button"
              onClick={() => setSelectedColor(color.value)}
              disabled={updateVault.isPending}
              className={cn(
                "size-10 rounded-full transition-all hover:scale-110 disabled:opacity-50 disabled:cursor-not-allowed",
                selectedColor === color.value && "ring-2 ring-offset-2 ring-foreground"
              )}
              style={{ backgroundColor: color.value }}
              aria-label={color.name}
            />
          ))}
        </div>
      </div>

      <div className="grid gap-3">
        <Label>Icon</Label>
        <div className="grid grid-cols-5 gap-2">
          {VAULT_ICON_LIST.map((icon) => {
            const IconComponent = icon.component;
            return (
              <button
                key={icon.name}
                type="button"
                onClick={() => setSelectedIcon(icon.name)}
                disabled={updateVault.isPending}
                className={cn(
                  "flex items-center justify-center p-3 rounded-lg border border-rose-50 transition-all hover:bg-accent disabled:opacity-50 disabled:cursor-not-allowed",
                  selectedIcon === icon.name && "bg-accent border-rose-100"
                )}
                aria-label={icon.name}
              >
                <IconComponent className="size-6" color={selectedColor} />
              </button>
            );
          })}
        </div>
      </div>

      <Button type="submit" disabled={updateVault.isPending}>
        {updateVault.isPending ? "Updating..." : "Update Vault"}
      </Button>
    </form>
  );
}
