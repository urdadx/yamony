import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { createFileRoute } from '@tanstack/react-router'
import { SecuritySettings } from '@/components/security-settings'
import { ExportSettings } from '@/components/export-settings'
import { AccountSettings } from '@/components/account-settings'
import { ImportSettings } from '@/components/import-settings'

export const Route = createFileRoute('/admin/settings')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
      <Tabs className="items-start p-6" defaultValue="account">
        <TabsList className="h-auto gap-2 rounded-none border-b bg-transparent px-0 text-foreground">
          <TabsTrigger
            className="after:-mb-1 relative after:absolute after:inset-x-0 after:bottom-0 after:h-0.5 hover:bg-accent hover:text-foreground data-[state=active]:bg-transparent data-[state=active]:shadow-none data-[state=active]:hover:bg-accent data-[state=active]:after:bg-primary"
            value="account"
          >
            Account
          </TabsTrigger>
          <TabsTrigger
            className="after:-mb-1 relative after:absolute after:inset-x-0 after:bottom-0 after:h-0.5 hover:bg-accent hover:text-foreground data-[state=active]:bg-transparent data-[state=active]:shadow-none data-[state=active]:hover:bg-accent data-[state=active]:after:bg-primary"
            value="security"
          >
            Security
          </TabsTrigger>
          <TabsTrigger
            className="after:-mb-1 relative after:absolute after:inset-x-0 after:bottom-0 after:h-0.5 hover:bg-accent hover:text-foreground data-[state=active]:bg-transparent data-[state=active]:shadow-none data-[state=active]:hover:bg-accent data-[state=active]:after:bg-primary"
            value="import"
          >
            Import
          </TabsTrigger>
          <TabsTrigger
            className="after:-mb-1 relative after:absolute after:inset-x-0 after:bottom-0 after:h-0.5 hover:bg-accent hover:text-foreground data-[state=active]:bg-transparent data-[state=active]:shadow-none data-[state=active]:hover:bg-accent data-[state=active]:after:bg-primary"
            value="export"
          >
            Export
          </TabsTrigger>

        </TabsList>
        <TabsContent className='p-2 m-0' value="account">
          <div className=" max-w-2xl">
            <AccountSettings />
          </div>
        </TabsContent>
        <TabsContent className='p-2 m-0' value="security">
          <div className=" max-w-2xl">
            <SecuritySettings />
          </div>
        </TabsContent>
        <TabsContent className='p-2 m-0' value="import">
          <div className=" max-w-2xl">
            <ImportSettings />
          </div>
        </TabsContent>
        <TabsContent className='p-2 m-0' value="export">
          <div className=" max-w-2xl">
            <ExportSettings />
          </div>
        </TabsContent>
      </Tabs>
    </>
  )
}
