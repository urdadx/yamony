import { createFileRoute } from '@tanstack/react-router'
import { VaultList } from '@/components/vault-list'
import { VaultItemDetails } from '@/components/vault-item-details'


export const Route = createFileRoute('/admin/home')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div className="flex h-screen w-full">
      <div className="shrink-0 border-r p-2">
        <VaultList />
      </div>
      <div className="flex-1 w-full bg-gray-50/50 smooth-div ">
        <VaultItemDetails />
      </div>
    </div>
  )
}
