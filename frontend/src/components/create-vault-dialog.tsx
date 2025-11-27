import * as React from "react";
import { cn } from "@/lib/utils";
import { Button } from "./ui/button";
import { Plus } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "./ui/drawer";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { useIsMobile } from "@/hooks/use-mobile";
import { HomeIcon } from "@/assets/icons/home-icon";
import { ShieldStrongIcon } from "@/assets/icons/shield-strong-icon";
import { SolarStarIcon } from "@/assets/icons/star-icon";
import { BrushIcon } from "@/assets/icons/brush-icon";
import { CameraIcon } from "@/assets/icons/camera-icon";
import { BoltIcon } from "@/assets/icons/bolt-icon";
import { PaletteIcon } from "@/assets/icons/palette-icon";
import { LinkDuoIcon } from "@/assets/icons/link-icon";
import { PhoneIcon } from "@/assets/icons/phone-icon";
import { LetterIcon } from "@/assets/icons/letter";

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

const VAULT_ICONS = [
  { name: "Home", component: HomeIcon },
  { name: "Shield", component: ShieldStrongIcon },
  { name: "Star", component: SolarStarIcon },
  { name: "Brush", component: BrushIcon },
  { name: "Camera", component: CameraIcon },
  { name: "Bolt", component: BoltIcon },
  { name: "Palette", component: PaletteIcon },
  { name: "Link", component: LinkDuoIcon },
  { name: "Phone", component: PhoneIcon },
  { name: "Letter", component: LetterIcon },
];

export const CreateVaultDialog = () => {
  const [open, setOpen] = React.useState(false);
  const isMobile = useIsMobile();

  if (!isMobile) {
    return (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <Button variant="ghost" className="size-8 p-0 text-muted-foreground hover:text-foreground">
            <Plus className="size-5" />
          </Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[550px] rounded-xl border-none bg-clip-padding shadow-2xl ring-4 ring-neutral-200/80 outline-none md:max-w-2xl dark:bg-neutral-800 dark:ring-neutral-900">
          <DialogHeader>
            <DialogTitle>Create Vault</DialogTitle>

          </DialogHeader>
          <VaultForm />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Drawer open={open} onOpenChange={setOpen}>
      <DrawerTrigger asChild>
        <Button variant="ghost" className="size-8 p-0 text-muted-foreground hover:text-foreground">
          <Plus className="size-5" />
        </Button>
      </DrawerTrigger>
      <DrawerContent>
        <DrawerHeader className="text-left">
          <DrawerTitle>Create Vault</DrawerTitle>
          <DrawerDescription>
            Create a new vault to organize your items. Choose a title, color, and icon.
          </DrawerDescription>
        </DrawerHeader>
        <VaultForm className="px-4" />
        <DrawerFooter className="pt-2">
          <DrawerClose asChild>
            <Button variant="outline">Cancel</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
};

function VaultForm({ className }: React.ComponentProps<"form">) {
  const [title, setTitle] = React.useState("");
  const [selectedColor, setSelectedColor] = React.useState(VAULT_COLORS[0].value);
  const [selectedIcon, setSelectedIcon] = React.useState("Home");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log({ title, selectedColor, selectedIcon });
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
              className={cn(
                "size-10 rounded-full transition-all hover:scale-110",
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
          {VAULT_ICONS.map((icon) => {
            const IconComponent = icon.component;
            return (
              <button
                key={icon.name}
                type="button"
                onClick={() => setSelectedIcon(icon.name)}
                className={cn(
                  "flex items-center justify-center p-3 rounded-lg border border-rose-50 transition-all hover:bg-accent",
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

      <Button type="submit">Create Vault</Button>
    </form>
  );
}