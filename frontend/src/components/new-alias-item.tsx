import { useIsMobile } from "@/hooks/use-mobile"
import { Button } from "./ui/button"
import { Sheet, SheetClose, SheetContent } from "./ui/sheet"
import { PaperclipIcon, StickyNote, AtSign, Forward } from "lucide-react"
import { Input } from "./ui/input"
import { Textarea } from "./ui/textarea"

interface NewAliasItemProps {
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

export const NewAliasItem = ({ open, onOpenChange }: NewAliasItemProps) => {
  const isMobile = useIsMobile()

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent>
        <SheetContent side={isMobile ? "bottom" : "right"} className="overflow-y-auto w-full sm:max-w-[550px]">
          <div className="flex px-4 pt-3 p-2 items-center justify-between">
            <SheetClose />
            <Button>
              Create alias
            </Button>
          </div>
          <div className="min-h-screen w-full px-4 flex flex-col gap-4">
            <div className="flex bg-white rounded-lg border items-center gap-4 px-4 p-2">
              <div className="flex-1 flex flex-col">
                <label className="block text-sm text-gray-500">
                  Title
                </label>
                <Input
                  placeholder="Untitled"
                  className="text-xl! rounded-none font-semibold border-gray-300 bg-transparent border-0 focus-visible:ring-offset-0 focus-visible:ring-0 p-0 h-auto"
                />
              </div>
            </div>

            <div className="w-full overflow-hidden bg-white border rounded-lg">
              <div className="flex items-center gap-4 px-4 py-4 border-b">
                <AtSign className="h-4 w-4 text-gray-500 shrink-0" />
                <div className="flex-1 flex flex-col">
                  <label className="block text-sm font-medium text-gray-500">
                    Alias
                  </label>
                  <Input
                    type="email"
                    placeholder="alias@example.com"
                    className="text-md! rounded-none border-gray-300 bg-transparent border-0 focus-visible:ring-offset-0 focus-visible:ring-0 p-0 h-auto"
                  />
                </div>
              </div>

              <div className="flex items-center gap-4 px-4 py-4">
                <Forward className="h-4 w-4 text-gray-500 shrink-0" />
                <div className="flex-1 flex flex-col">
                  <label className="block text-sm font-medium text-gray-500">
                    Forwards to
                  </label>
                  <Input
                    type="email"
                    placeholder="your-email@example.com"
                    className="text-md! rounded-none border-gray-300 bg-transparent border-0 focus-visible:ring-offset-0 focus-visible:ring-0 p-0 h-auto"
                  />
                </div>
              </div>
            </div>

            <div className="flex items-start gap-4 px-4 py-3 rounded-lg border">
              <StickyNote className="h-4 w-4 text-gray-500 shrink-0 mt-1" />
              <div className="flex-1 flex flex-col">
                <label className="block text-sm font-medium text-gray-500">
                  Note
                </label>
                <Textarea
                  placeholder="Add your note here..."
                  className="text-md! shadow-none rounded-none outline-0 min-h-[30px] border-gray-300 bg-transparent border-0 focus-visible:ring-offset-0 focus-visible:ring-0 p-0 resize-none"
                />
              </div>
            </div>

            <div className="flex flex-col gap-4 px-4 py-3 rounded-lg border">
              <div className="flex items-center gap-4">
                <PaperclipIcon className="h-4 w-4 text-gray-500 shrink-0" />
                <div className="flex-1 flex flex-col">
                  <label className="block text-sm font-medium text-gray-500">
                    Attachments
                  </label>
                  <div className="flex items-center gap-2">
                    <div className="text-sm text-gray-900">
                      Upload files from your device
                    </div>
                  </div>
                </div>
              </div>
              <Button>
                Choose a file or drag it here
              </Button>
            </div>
          </div>
        </SheetContent>
      </SheetContent>
    </Sheet>
  )
}
