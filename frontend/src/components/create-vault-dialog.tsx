import { DialogTrigger } from "@radix-ui/react-dialog";
import { Dialog, DialogContent } from "./ui/dialog";
import { Button } from "./ui/button";
import { Plus } from "lucide-react";

export const CreateVaultDialog = () => {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="ghost" className="size-8 p-0 text-muted-foreground hover:text-foreground">
          <Plus className="size-5" />
        </Button>
      </DialogTrigger>
      <DialogContent className="w-[320px] h-[200px] rounded-xl border-none bg-clip-padding shadow-2xl ring-4 ring-neutral-200/80 outline-none md:max-w-2xl dark:bg-neutral-800 dark:ring-neutral-900">

      </DialogContent>
    </Dialog>
  );
}