import * as React from "react"
import { useIsMobile } from "@/hooks/use-mobile"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Drawer,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer"
import { Button } from "@/components/ui/button"
import { UserIcon } from "@/assets/icons/user-icon"
import { CreditCardIcon } from "@/assets/icons/credit-card-icon"
import { SolarNotebookBoldDuotone } from "@/assets/icons/note-icon"
import { IcognitoIcon } from "@/assets/icons/icognito-icon"
import { AddIcon } from "@/assets/icons/add-icon"
import { NewLoginItem } from "./new-login-item"
import { NewCardItem } from "./new-card-item"
import { NewNoteItem } from "./new-note-item"
import { NewAliasItem } from "./new-alias-item"

type ItemType = "login" | "card" | "note" | "alias"

interface CreateItemDialogProps {
  open?: boolean
  onOpenChange?: (open: boolean) => void
  onSelectType?: (type: ItemType) => void
}

export function CreateItemDialog({
  open,
  onOpenChange,
  onSelectType
}: CreateItemDialogProps) {
  const [isOpen, setIsOpen] = React.useState(false)
  const [loginItemOpen, setLoginItemOpen] = React.useState(false)
  const [cardItemOpen, setCardItemOpen] = React.useState(false)
  const [noteItemOpen, setNoteItemOpen] = React.useState(false)
  const [aliasItemOpen, setAliasItemOpen] = React.useState(false)
  const isMobile = useIsMobile()

  const actualOpen = open !== undefined ? open : isOpen
  const actualOnOpenChange = onOpenChange || setIsOpen

  const handleSelectType = (type: ItemType) => {
    if (type === "login") {
      setLoginItemOpen(true)
      actualOnOpenChange(false)
    } else if (type === "card") {
      setCardItemOpen(true)
      actualOnOpenChange(false)
    } else if (type === "note") {
      setNoteItemOpen(true)
      actualOnOpenChange(false)
    } else if (type === "alias") {
      setAliasItemOpen(true)
      actualOnOpenChange(false)
    } else {
      onSelectType?.(type)
      actualOnOpenChange(false)
    }
  }

  const itemTypes = [
    {
      type: "login" as ItemType,
      label: "Login",
      icon: UserIcon,
      iconColor: "#3b82f6",
      bgColor: "bg-blue-50 dark:bg-blue-950/30",
      borderColor: "border-blue-200 dark:border-blue-800",
    },
    {
      type: "card" as ItemType,
      label: "Card",
      icon: CreditCardIcon,
      iconColor: "#8b5cf6",
      bgColor: "bg-purple-50 dark:bg-purple-950/30",
      borderColor: "border-purple-200 dark:border-purple-800",
    },
    {
      type: "note" as ItemType,
      label: "Note",
      icon: SolarNotebookBoldDuotone,
      iconColor: "#10b981",
      bgColor: "bg-emerald-50 dark:bg-emerald-950/30",
      borderColor: "border-emerald-200 dark:border-emerald-800",
    },
    {
      type: "alias" as ItemType,
      label: "Alias",
      icon: IcognitoIcon,
      iconColor: "#f59e0b",
      bgColor: "bg-amber-50 dark:bg-amber-950/30",
      borderColor: "border-amber-200 dark:border-amber-800",
    },
  ]

  const ItemTypeGrid = () => (
    <div className="grid grid-cols-2 gap-4 py-4">
      {itemTypes.map(({ type, label, icon: Icon, iconColor, bgColor, borderColor }) => (
        <button
          key={type}
          onClick={() => handleSelectType(type)}
          className={`flex flex-col items-center justify-center gap-2 p-2 rounded-lg border-2 transition-all hover:scale-101 ${bgColor} ${borderColor} hover:shadow-sm`}
        >
          <div className="text-4xl">
            <Icon color={iconColor} />
          </div>
          <span className="text-sm font-medium text-foreground">{label}</span>
        </button>
      ))}
    </div>
  )

  if (!isMobile) {
    return (
      <>
        <Dialog open={actualOpen} onOpenChange={actualOnOpenChange}>
          <DialogTrigger >
            <Button>
              <AddIcon color="#FFFFFF" className="w-5 h-5 " />
              Create item
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px] rounded-xl border-none bg-clip-padding shadow-2xl ring-4 ring-neutral-200/80 outline-none md:max-w-lg dark:bg-neutral-800 dark:ring-neutral-900">
            <DialogHeader>
              <DialogTitle>Create New Item</DialogTitle>
            </DialogHeader>
            <ItemTypeGrid />
          </DialogContent>
        </Dialog>
        <NewLoginItem open={loginItemOpen} onOpenChange={setLoginItemOpen} />
        <NewCardItem open={cardItemOpen} onOpenChange={setCardItemOpen} />
        <NewNoteItem open={noteItemOpen} onOpenChange={setNoteItemOpen} />
        <NewAliasItem open={aliasItemOpen} onOpenChange={setAliasItemOpen} />
      </>
    )
  }

  return (
    <>
      <Drawer open={actualOpen} onOpenChange={actualOnOpenChange}>
        <DrawerTrigger>
          <Button>
            <AddIcon color="#FFFFFF" className="w-5 h-5 " />
            Create item
          </Button>      </DrawerTrigger>
        <DrawerContent>
          <DrawerHeader className="text-left">
            <DrawerTitle>Create New Item</DrawerTitle>
            <DrawerDescription>
              Choose the type of item you want to create
            </DrawerDescription>
          </DrawerHeader>
          <div className="px-4 pb-8">
            <ItemTypeGrid />
          </div>
        </DrawerContent>
      </Drawer>
      <NewLoginItem open={loginItemOpen} onOpenChange={setLoginItemOpen} />
      <NewCardItem open={cardItemOpen} onOpenChange={setCardItemOpen} />
      <NewNoteItem open={noteItemOpen} onOpenChange={setNoteItemOpen} />
      <NewAliasItem open={aliasItemOpen} onOpenChange={setAliasItemOpen} />
    </>
  )
}