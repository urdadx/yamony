import { User } from "lucide-react"
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar"
import { getApexDomain, GOOGLE_FAVICON_URL } from "@/lib/utils"

export const VaultItem = () => {

  const apexDomain = getApexDomain("https://yahoo.com");
  return (
    <div className="w-full cursor-pointer flex items-center gap-2 hover:bg-gray-50 rounded-md p-2">
      <Avatar className="flex aspect-square size-10 items-center justify-center rounded-xl text-sidebar-primary-foreground">
        <div className="bg-rose-50 w-full h-full flex items-center justify-center">
          <AvatarImage
            width={24}
            height={24}
            className="w-5 h-5"
            src={`${GOOGLE_FAVICON_URL}${apexDomain}`}
            alt="vault avatar"
          />
        </div>
        <AvatarFallback className="bg-rose-50 ">
          <User className="size-5 text-primary" />
        </AvatarFallback>
      </Avatar>
      <div className="grid flex-1 text-left text-sm leading-tight">
        <span className="truncate ">
          yahoo.com
        </span>
        <span className="truncate text-muted-foreground text-xs">
          abassabdulwahab3@gmail.com
        </span>
      </div>
    </div>
  )
}