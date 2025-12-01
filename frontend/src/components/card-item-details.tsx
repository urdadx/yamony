import { useState } from "react"
import { CreditCard, Calendar, Lock, KeyRound, User, StickyNote } from "lucide-react"
import { Button } from "@/components/ui/button"
import { ItemHeaderOptions } from "./item-header-options"

export function CardItemDetails() {
  const [showNumber, setShowNumber] = useState(false)
  const [showPin, setShowPin] = useState(false)

  const maskedNumber = "•••• •••• •••• ••••"
  const cardNumber = "4242 4242 4242 4242"
  const maskedPin = "••••"
  const pinValue = "1234"

  return (
    <main className="min-h-screen w-full p-4">
      <div className="w-full">
        <ItemHeaderOptions />

        <div className="w-full overflow-hidden bg-white border rounded-lg">
          <div className="flex items-center gap-4 px-4 py-4 border-b hover:bg-gray-50 transition-colors cursor-pointer">
            <User className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">Name on card</label>
              <div className="text-sm text-gray-900">John Doe</div>
            </div>
          </div>

          <div className="flex items-center gap-4 px-4 py-4 border-b hover:bg-gray-50 transition-colors cursor-pointer">
            <CreditCard className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">Card number</label>
              <div className="flex items-center justify-between">
                <div className="text-sm text-gray-900 flex-1">{showNumber ? cardNumber : maskedNumber}</div>
                <Button variant="ghost" size="sm" className="h-8">{showNumber ? "Hide" : "Show"}</Button>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-4 px-4 py-4 border-b hover:bg-gray-50 transition-colors cursor-pointer">
            <Calendar className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">Expiration date</label>
              <div className="text-sm text-gray-900">08/28</div>
            </div>
          </div>

          <div className="flex items-center gap-4 px-4 py-4 border-b hover:bg-gray-50 transition-colors cursor-pointer">
            <Lock className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">Security code</label>
              <div className="text-sm text-gray-900">•••</div>
            </div>
          </div>

          <div className="flex items-center gap-4 px-4 py-4 hover:bg-gray-50 transition-colors cursor-pointer">
            <KeyRound className="h-4 w-4 text-gray-500 shrink-0" />
            <div className="flex-1">
              <label className="block text-xs font-medium text-gray-500 mb-1">PIN</label>
              <div className="flex items-center justify-between">
                <div className="text-sm text-gray-900 flex-1">{showPin ? pinValue : maskedPin}</div>
                <Button variant="ghost" size="sm" className="h-8" onClick={() => setShowPin(!showPin)}>{showPin ? "Hide" : "Show"}</Button>
              </div>
            </div>
          </div>
        </div>

        <div className="flex mt-3 bg-white rounded-lg border items-center gap-4 px-4 py-4 hover:bg-gray-50 transition-colors cursor-pointer">
          <StickyNote className="h-4 w-4 text-gray-500 shrink-0" />
          <div className="flex-1">
            <label className="block text-xs font-medium text-gray-500 mb-1">Notes</label>
            <div className="text-sm text-gray-900">Primary business card</div>
          </div>
        </div>
      </div>
    </main>
  )
}
