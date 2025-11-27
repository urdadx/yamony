import { useId } from "react"
import { SearchIcon } from "lucide-react"

import { Input } from "@/components/ui/input"

export function GlobalSearch() {
  const id = useId()
  return (
    <div className="*:not-first:mt-2 w-full">
      <div className="relative">
        <Input
          id={id}
          className="peer ps-9 pe-9"
          placeholder="Search in Personal"
          type="search"
        />
        <div className="pointer-events-none absolute inset-y-0 start-0 flex items-center justify-center ps-3 text-muted-foreground/80 peer-disabled:opacity-50">
          <SearchIcon size={16} />
        </div>

      </div>
    </div>
  )
}
