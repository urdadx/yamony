import { Label } from "@/components/ui/label"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Checkbox } from "@/components/ui/checkbox"
import { useState } from "react"

export function SecuritySettings() {
  const [unlockMethod, setUnlockMethod] = useState("password")
  const [autoLockTime, setAutoLockTime] = useState("1hour")
  const [extraPassword, setExtraPassword] = useState(false)

  return (
    <div className="space-y-8 py-3">
      <div className="space-y-4 border rounded-md p-6">
        <h3 className="font-semibold text-base">Unlock with</h3>

        <RadioGroup value={unlockMethod} onValueChange={setUnlockMethod}>
          <div className="space-y-4">
            <div className="flex items-start gap-3">
              <RadioGroupItem value="none" id="none" className="mt-0.5" />
              <div className="flex flex-col gap-1">
                <Label htmlFor="none" className="font-medium cursor-pointer">
                  None
                </Label>
                <p className="text-sm text-muted-foreground">
                  Nucleo Pass will always be accessible
                </p>
              </div>
            </div>

            {/* PIN Code Option */}
            <div className="flex items-start gap-3">
              <RadioGroupItem value="pin" id="pin" className="mt-0.5" />
              <div className="flex flex-col gap-1">
                <Label htmlFor="pin" className="font-medium cursor-pointer">
                  PIN code
                </Label>
                <p className="text-sm text-muted-foreground">
                  Online access to Nucleo Pass will require a PIN code. You'll be logged out after 3 failed attempts.
                </p>
              </div>
            </div>

            {/* Password Option */}
            <div className="flex items-start gap-3">
              <RadioGroupItem value="password" id="password" className="mt-0.5" />
              <div className="flex flex-col gap-1">
                <Label htmlFor="password" className="font-medium cursor-pointer">
                  Password
                </Label>
                <p className="text-sm text-muted-foreground">
                  Access to Nucleo Pass will always require your Nucleo password.
                </p>
              </div>
            </div>
          </div>
        </RadioGroup>
        {/* Auto-lock after Section */}
        <div className="space-y-4">
          <h3 className="font-semibold text-base">Auto-lock after</h3>

          <Select value={autoLockTime} onValueChange={setAutoLockTime}>
            <SelectTrigger className="w-full">
              <SelectValue placeholder="Select auto-lock time" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1min">1 minute</SelectItem>
              <SelectItem value="5min">5 minutes</SelectItem>
              <SelectItem value="15min">15 minutes</SelectItem>
              <SelectItem value="30min">30 minutes</SelectItem>
              <SelectItem value="1hour">1 hour</SelectItem>
              <SelectItem value="4hours">4 hours</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Extra password Section */}
      <div className="space-y-4 border rounded-md p-6">
        <h3 className="font-semibold text-base">Extra password</h3>

        <div className="flex items-start gap-3">
          <Checkbox
            id="extra-password"
            checked={extraPassword}
            onCheckedChange={(checked) => setExtraPassword(checked as boolean)}
          />
          <div className="flex flex-col gap-1">
            <Label htmlFor="extra-password" className="font-medium cursor-pointer">
              Protect Nucleo Pass with an extra password
            </Label>
            <p className="text-sm text-muted-foreground">
              The extra password will be required to use Nucleo Pass. It acts as an additional password on top of your Nucleo password.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
