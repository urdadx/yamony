import { useId } from "react"
import { RiApps2Line, RiBankCardLine, RiLockPasswordLine, RiStickyNoteLine } from "@remixicon/react"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

export function FilterByCategories() {
  const id = useId()
  return (
    <div className="*:not-first:mt-2">
      <Select defaultValue="2">
        <SelectTrigger
          id={id}
          className="[&>span]:flex [&>span]:items-center [&>span]:gap-2 [&>span_svg]:shrink-0 [&>span_svg]:text-muted-foreground/80 w-[115px]"
        >
          <SelectValue placeholder="Filter" />
        </SelectTrigger>
        <SelectContent className="[&_*[role=option]]:ps-2 [&_*[role=option]]:pe-8 [&_*[role=option]>span]:start-auto [&_*[role=option]>span]:end-2 [&_*[role=option]>span]:flex [&_*[role=option]>span]:items-center [&_*[role=option]>span]:gap-2 [&_*[role=option]>span>svg]:shrink-0 [&_*[role=option]>span>svg]:text-muted-foreground/80 w-[150px]">
          <SelectItem value="1">
            <RiApps2Line size={16} aria-hidden="true" />
            <span className="truncate">All</span>
          </SelectItem>
          <SelectItem value="2">
            <RiBankCardLine size={16} aria-hidden="true" />
            <span className="truncate">Cards</span>
          </SelectItem>
          <SelectItem value="3">
            <RiLockPasswordLine size={16} aria-hidden="true" />
            <span className="truncate">Logins</span>
          </SelectItem>
          <SelectItem value="4">
            <RiStickyNoteLine size={16} aria-hidden="true" />
            <span className="truncate">Notes</span>
          </SelectItem>
        </SelectContent>
      </Select>
    </div>
  )
}
