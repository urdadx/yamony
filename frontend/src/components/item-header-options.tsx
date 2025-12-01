import { EditPencilIcon } from "@/assets/icons/edit-pencil-icon"
import { Button } from "./ui/button"
import { MoreVertical } from "lucide-react"
import { ShareVaultItemDialog } from "./share-vault-item-dialog"
import { ItemOptionsPopover } from "./item-options-popover"
import { PopoverTrigger } from "./ui/popover"

export const ItemHeaderOptions = () => {
  return (
    <>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-lg font-bold text-gray-900">wd1.myworkdayjobs.com</h1>
        <div className="flex items-center gap-2">
          <Button variant="outline" className=" text-primary rounded-lg shadow-none p-4 hover:bg-rose-50/40 hover:text-primary/80">
            <EditPencilIcon color="#f43f5e" className="hover:scale-110 transition-transform" />
            <span className="hover:opacity-80 font-normal transition-opacity">Edit</span>
          </Button>
          <ShareVaultItemDialog itemTitle="wd1.myworkdayjobs.com" />

          <ItemOptionsPopover>
            <PopoverTrigger asChild>
              <Button variant="outline" className="text-primary rounded-lg shadow-none p-3 hover:bg-rose-50/40 hover:text-primary/80">
                <MoreVertical size={12} className="hover:scale-110 text-primary transition-transform" />
              </Button>
            </PopoverTrigger>
          </ItemOptionsPopover>
        </div>
      </div>
    </>
  )

}