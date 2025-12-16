import { User } from "lucide-react"
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar"
import { getApexDomain, GOOGLE_FAVICON_URL } from "@/lib/utils"

export const VaultItem = () => {

  const apexDomain = getApexDomain("https://google.com");

  return (
    <div className="w-full cursor-pointer flex items-center gap-2 hover:bg-gray-50 rounded-md p-2">
      <Avatar className="size-10 rounded-xl">
        <AvatarImage
          className="rounded-xl object-cover"
          src={`${GOOGLE_FAVICON_URL}${apexDomain}`}
          alt="vault avatar"
        />
        <AvatarFallback className="bg-rose-50 rounded-xl">
          <User className="size-5 text-primary" />
        </AvatarFallback>
      </Avatar>
      <div className="grid flex-1 text-left text-sm leading-tight">
        <span className="truncate ">
          google.com
        </span>
        <span className="truncate text-muted-foreground text-xs">
          abassabdulwahab3@gmail.com
        </span>
      </div>
    </div>
  )
}