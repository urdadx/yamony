import { StickyNote, PencilIcon, ZapIcon } from "lucide-react"
import { ItemHeaderOptions } from "./item-header-options"

export function NoteItemDetails() {
  return (
    <main className="min-h-screen w-full p-4">
      <div className="w-full">
        <ItemHeaderOptions />

        <div className="flex bg-white rounded-lg border items-start gap-4 px-4 py-4 hover:bg-gray-50 transition-colors cursor-pointer">
          <StickyNote className="h-4 w-4 text-gray-500 shrink-0 mt-1" />
          <div className="flex-1">
            <label className="block text-xs font-medium text-gray-500 mb-1">Note</label>
            <div className="text-sm text-gray-900 whitespace-pre-wrap">
              This is my secure note content. It can contain multiple lines and important information that I want to keep safe.
            </div>
          </div>
        </div>

        <div className="mt-3 w-full overflow-hidden border rounded-lg">
          <div className="flex items-center gap-4 px-4 py-3 hover:bg-gray-50 transition-colors cursor-pointer">
            <PencilIcon className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">Last edited</label>
              <div className="text-sm text-gray-900">Nov 23, 2025</div>
            </div>
          </div>

          <div className="flex items-center gap-4 px-4 py-3 hover:bg-gray-50 transition-colors cursor-pointer">
            <ZapIcon className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">Created on</label>
              <div className="text-sm text-gray-900">Oct 15, 2025</div>
            </div>
          </div>
        </div>
      </div>
    </main>
  )
}
