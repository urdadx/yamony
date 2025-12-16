import * as React from "react";
import { Edit, Trash2, FolderInput } from "lucide-react";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { EditVaultDialog } from "@/components/edit-vault-dialog";
import { useState } from "react";
import { DeleteVaultDialog } from "@/components/delete-vault-dialog";

interface VaultOptionsPopoverProps {
  children: React.ReactNode;
  vault: {
    id: string;
    name: string;
    icon?: string;
    theme?: string;
  };
  onMove?: () => void;
  onDelete?: () => void;
}

export function VaultOptionsPopover({
  children,
  vault,
  onMove,
  onDelete,
}: VaultOptionsPopoverProps) {
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [popoverOpen, setPopoverOpen] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);

  return (
    <>
      <Popover open={popoverOpen} onOpenChange={setPopoverOpen}>
        <PopoverTrigger asChild>{children}</PopoverTrigger>
        <PopoverContent className="w-48 p-1" align="end">
          <div className="flex flex-col gap-0.5">
            <Button
              variant="ghost"
              className="justify-start gap-2 px-2 py-1.5 h-auto font-normal"
              onClick={() => {
                setPopoverOpen(false);
                setEditDialogOpen(true);
              }}
            >
              <Edit className="size-4" />
              <span>Edit</span>
            </Button>
            <Button
              variant="ghost"
              className="justify-start gap-2 px-2 py-1.5 h-auto font-normal"
              onClick={onMove}
            >
              <FolderInput className="size-4" />
              <span>Move</span>
            </Button>
            <Button
              variant="ghost"
              className="justify-start gap-2 px-2 py-1.5 h-auto font-normal text-destructive hover:text-destructive"
              onClick={() => {
                setPopoverOpen(false);
                setConfirmOpen(true);
              }}
            >
              <Trash2 className="size-4" />
              <span>Delete</span>
            </Button>
          </div>
        </PopoverContent>
      </Popover>

      <EditVaultDialog
        open={editDialogOpen}
        onOpenChange={setEditDialogOpen}
        vault={vault}
      />

      <DeleteVaultDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        vaultId={vault.id}
        vaultName={vault.name}
        onDeleted={onDelete}
      />
    </>
  );
}
