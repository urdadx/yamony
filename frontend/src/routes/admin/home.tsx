import { createFileRoute } from '@tanstack/react-router'
import { ChooseBlock } from '@/components/choose-block'

export const Route = createFileRoute('/admin/home')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div className="flex h-screen w-full">
      {/* Preview Container */}
      <div className="flex-1 p-4 sm:p-6 lg:p-8 overflow-auto">
        <h1 className="text-2xl font-bold mb-4">Preview will be here</h1>
      </div>

      {/* Choose Block Container */}
      <div className="shrink-0">
        <ChooseBlock />
      </div>
    </div>
  )
}
