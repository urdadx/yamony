import { SolarLetterBoldDuotone } from "@/assets/icons/letter";
import { Button } from "./ui/button";
import { SidebarTrigger } from "./ui/sidebar";

export const Navbar = () => {

  return (
    <>

      <header className="px-4 sticky top-0 flex justify-between h-16 shrink-0 items-center  bg-background/50 backdrop-blur-lg border-b transition-[width,height] ease-linear z-10 group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
        <SidebarTrigger className="-ml-3 hidden sm:flex text-muted-foreground" />

        <div className="flex items-center px-0 sm:px-3 capitalize">
          <div className="flex sm:hidden">
            <SidebarTrigger className=" w-8 h-8 text-muted-foreground" />
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline">
            Share
          </Button>
          <Button
          >
            <SolarLetterBoldDuotone color="#FFFFFF" className="w-5 h-5 " />
            Invite team
          </Button>
        </div>
      </header>
    </>
  );
};