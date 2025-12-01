import { AtSign, Forward, StickyNote, PencilIcon, ZapIcon } from "lucide-react"
import { ItemHeaderOptions } from "./item-header-options"

export function AliasItemDetails() {
  return (
    <main className="min-h-screen w-full p-4">
      <div className="w-full">
        <ItemHeaderOptions />

        <div className="w-full overflow-hidden bg-white border rounded-lg">
          <div className="flex items-center gap-4 px-4 py-4 border-b hover:bg-gray-50 transition-colors cursor-pointer">
            <AtSign className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">Alias</label>
              <div className="text-sm text-gray-900">privacy.alias@example.com</div>
            </div>
          </div>

          <div className="flex items-center gap-4 px-4 py-4 hover:bg-gray-50 transition-colors cursor-pointer">
            <Forward className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">Forwards to</label>
              <div className="text-sm text-gray-900">myemail@example.com</div>
            </div>
          </div>
        </div>

        <div className="flex mt-3 bg-white rounded-lg border items-center gap-4 px-4 py-4 hover:bg-gray-50 transition-colors cursor-pointer">
          <StickyNote className="h-4 w-4 text-gray-500 shrink-0" />
          <div className="flex-1">
            <label className="block text-xs font-medium text-gray-500 mb-1">Note</label>
            <div className="text-sm text-gray-900">Used for newsletter subscriptions</div>
          </div>
        </div>

        <div className="mt-3 w-full overflow-hidden border rounded-lg">
          <div className="flex items-center gap-4 px-4 py-3 hover:bg-gray-50 transition-colors cursor-pointer">
            <PencilIcon className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">Last edited</label>
              <div className="text-sm text-gray-900">Nov 20, 2025</div>
            </div>
          </div>

          <div className="flex items-center gap-4 px-4 py-3 hover:bg-gray-50 transition-colors cursor-pointer">
            <ZapIcon className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">Created on</label>
              <div className="text-sm text-gray-900">Sep 10, 2025</div>
            </div>
          </div>
        </div>
      </div>
    </main>
  )
}
