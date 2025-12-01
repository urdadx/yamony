import { Copy, Pin, Trash2 } from "lucide-react"
import { Button } from "./ui/button"
import { Popover, PopoverContent } from "./ui/popover"

interface ItemOptionsPopoverProps {
  children: React.ReactNode
}

export const ItemOptionsPopover = ({ children }: ItemOptionsPopoverProps) => {
  return (
    <Popover>
      {children}
      <PopoverContent className="w-48 p-2" align="end">
        <div className="flex flex-col gap-1">
          <Button
            variant="ghost"
            className="justify-start gap-3 hover:bg-gray-100 dark:hover:bg-gray-800"
            onClick={() => console.log('Duplicate')}
          >
            <Copy size={16} />
            <span>Duplicate</span>
          </Button>
          <Button
            variant="ghost"
            className="justify-start gap-3 hover:bg-gray-100 dark:hover:bg-gray-800"
            onClick={() => console.log('Pin item')}
          >
            <Pin size={16} />
            <span>Pin item</span>
          </Button>
          <Button
            variant="ghost"
            className="justify-start gap-3 text-red-600 hover:bg-red-50 hover:text-red-700 dark:hover:bg-red-950"
            onClick={() => console.log('Move to trash')}
          >
            <Trash2 size={16} />
            <span>Move to trash</span>
          </Button>
        </div>
      </PopoverContent>
    </Popover>
  )
}
