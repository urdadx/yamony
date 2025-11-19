import { AddIcon } from "@/assets/icons/add-icon";
import { GlobalSearch } from "./global-search";
import { Button } from "./ui/button";
import { SidebarTrigger } from "./ui/sidebar";

export const Navbar = () => {

  return (
    <>

      <header className="px-4 sticky top-0 flex justify-between h-14 shrink-0 items-center  bg-background/50 backdrop-blur-lg border-b transition-[width,height] ease-linear z-10 group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">

        <div className="flex items-center px-0  capitalize">
          <div className="flex sm:hidden">
            <SidebarTrigger className=" w-8 h-8 text-muted-foreground" />
          </div>
        </div>
        <div className="flex items-center gap-3 justify-between flex-1">
          <GlobalSearch />
          <div className="flex gap-2 items-center">
            <Button variant="outline">
              Share
            </Button>
            <Button>
              <AddIcon color="#FFFFFF" className="w-5 h-5 " />
              Create item
            </Button>
          </div>
        </div>
      </header>
    </>
  );
};