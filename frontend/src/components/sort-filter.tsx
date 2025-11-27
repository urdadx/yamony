import { useId } from "react"
import { RiTimeLine, RiSortAlphabetAsc } from "@remixicon/react"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

export function SortFilter() {
  const id = useId()
  return (
    <div className="">
      <Select defaultValue="recent">
        <SelectTrigger
          id={id}
          className="[&>span]:flex [&>span]:items-center [&>span]:gap-2 [&>span_svg]:shrink-0 [&>span_svg]:text-muted-foreground/80 w-[115px]"
        >
          <SelectValue placeholder="Sort by" />
        </SelectTrigger>
        <SelectContent className="[&_*[role=option]]:ps-2 [&_*[role=option]]:pe-8 [&_*[role=option]>span]:start-auto [&_*[role=option]>span]:end-2 [&_*[role=option]>span]:flex [&_*[role=option]>span]:items-center [&_*[role=option]>span]:gap-2 [&_*[role=option]>span>svg]:shrink-0 [&_*[role=option]>span>svg]:text-muted-foreground/80 w-[150px]">
          <SelectItem value="recent">
            <RiTimeLine size={16} aria-hidden="true" />
            <span className="truncate">Recent</span>
          </SelectItem>
          <SelectItem value="alphabetical">
            <RiSortAlphabetAsc size={16} aria-hidden="true" />
            <span className="truncate">Alphabetical</span>
          </SelectItem>
        </SelectContent>
      </Select>
    </div>
  )
}
