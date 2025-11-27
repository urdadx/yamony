import { SidebarTrigger } from "./ui/sidebar";
import { PasswordGenerator } from "./password-generator";
import { CreateItemDialog } from "./create-item-dialog";

export const Navbar = () => {

  return (
    <>
      <header className="px-4 sticky top-0 flex justify-between h-14 shrink-0 items-center  bg-background/50 backdrop-blur-lg border-b transition-[width,height] ease-linear z-10 group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
        <SidebarTrigger className="hidden sm:flex w-8 h-8 text-muted-foreground" />
        <div className="flex items-center px-0  capitalize">
          <div className="flex sm:hidden">
            <SidebarTrigger className=" w-8 h-8 text-muted-foreground" />
          </div>
        </div>
        <div className="flex items-center gap-3 justify-end flex-1">
          {/* <GlobalSearch /> */}
          <div className="flex gap-2 items-center">
            <PasswordGenerator />
            <CreateItemDialog />
          </div>
        </div>
      </header>
    </>
  );
};