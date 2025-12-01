import * as React from "react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer"
import { Input } from "@/components/ui/input"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Switch } from "@/components/ui/switch"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useIsMobile } from "@/hooks/use-mobile"
import { Globe, Copy, Check } from "lucide-react"
import { UserPlusIcon } from "@/assets/icons/user-plus-icon"

interface ShareVaultItemDialogProps {
  itemTitle?: string
}

export function ShareVaultItemDialog({ itemTitle = "Untitled document 2025-11-26 16.23.12" }: ShareVaultItemDialogProps) {
  const isMobile = useIsMobile()

  if (!isMobile) {
    return (
      <Dialog>
        <DialogTrigger asChild>
          <Button variant="outline" className="text-primary rounded-lg shadow-none p-3 hover:bg-rose-50/40 hover:text-primary/80">
            <UserPlusIcon color="#f43f5e" className="hover:scale-110 transition-transform" />
          </Button>
        </DialogTrigger>
        <DialogContent className="rounded-xl border-none bg-clip-padding shadow-2xl ring-4 ring-neutral-200/80 outline-none md:max-w-lg dark:bg-neutral-800 dark:ring-neutral-900">
          <DialogHeader className="space-y-0">
            <div className="flex items-center justify-between">
              <DialogTitle className="text-xl font-semibold">
                Share {itemTitle}
              </DialogTitle>

            </div>
          </DialogHeader>
          <ShareForm />
        </DialogContent>
      </Dialog>
    )
  }

  return (
    <Drawer >
      <DrawerTrigger asChild>
        <Button variant="outline" className="text-primary rounded-lg shadow-none p-3 hover:bg-rose-50/40 hover:text-primary/80">
          <UserPlusIcon color="#f43f5e" className="hover:scale-110 transition-transform" />
        </Button>
      </DrawerTrigger>
      <DrawerContent>
        <DrawerHeader className="text-left">
          <DrawerTitle className="text-xl font-semibold">
            Share {itemTitle}
          </DrawerTitle>
        </DrawerHeader>
        <ShareForm className="px-4" />
      </DrawerContent>
    </Drawer>
  )
}

function ShareForm({ className }: React.ComponentProps<"div">) {
  const [publicLinkEnabled, setPublicLinkEnabled] = React.useState(false)
  const [copied, setCopied] = React.useState(false)
  const shareableLink = "https://yamony.app/share/abc123xyz" // Replace with actual link generation logic

  const handleCopyLink = async () => {
    try {
      await navigator.clipboard.writeText(shareableLink)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  return (
    <div className={cn("space-y-6", className)}>
      <div>
        <div className="flex items-center gap-2">
          <Input
            placeholder="Add people or groups to share"
            className="flex-1"
          />

        </div>
      </div>

      {/* People with access */}
      <div>
        <h3 className="text-sm font-semibold mb-3">People with access</h3>
        <div className="flex items-center justify-between py-2">
          <div className="flex items-center gap-3">
            <Avatar className="h-10 w-10">
              <AvatarFallback className="bg-orange-200 text-orange-700 text-sm font-medium">
                N
              </AvatarFallback>
            </Avatar>
            <div>
              <div className="text-sm font-medium">(you)</div>
              <div className="text-sm text-gray-500">nerdyshinobi12@gmail.com</div>
            </div>
          </div>
          <div className="text-sm text-gray-600 font-medium">Owner</div>
        </div>
      </div>

      {/* Divider */}
      <div className="border-t border-gray-200" />

      {/* Create public link */}
      <div>
        <div className="flex items-center justify-between">
          <h3 className="text-sm font-semibold">Create public link</h3>
          <Switch
            checked={publicLinkEnabled}
            onCheckedChange={setPublicLinkEnabled}
          />
        </div>

        {publicLinkEnabled && (
          <div className="mt-4 space-y-3">
            <div className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg dark:bg-neutral-900">
              <Globe className="h-5 w-5 text-gray-400 mt-0.5 shrink-0" />
              <div className="flex-1">
                <div className="text-sm font-medium text-gray-600 dark:text-gray-300">Anyone with the link</div>
                <div className="text-sm text-gray-500 dark:text-gray-400">Anyone on the Internet with the link can view</div>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <Input
                value={shareableLink}
                readOnly
                className="flex-1 font-mono text-sm"
              />
              <Button
                type="button"
                variant="outline"
                size="icon"
                onClick={handleCopyLink}
                className="shrink-0"
              >
                {copied ? (
                  <Check className="h-4 w-4 text-green-600" />
                ) : (
                  <Copy className="h-4 w-4" />
                )}
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
