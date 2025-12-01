import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { Switch } from "@/components/ui/switch"
import { Input } from "@/components/ui/input"
import { useState } from "react"
import { Info } from "lucide-react"

export function AccountSettings() {
  const [twoFactorEnabled, setTwoFactorEnabled] = useState(false)
  const [displayName, setDisplayName] = useState("")

  return (
    <div className="space-y-8 py-3">
      <div className="space-y-6 border rounded-md p-6">
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <Label className="text-sm font-normal">Username</Label>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm">nerdyshinobi12@gmail.com</span>
              <div className="h-5 w-5 rounded-full bg-green-600 flex items-center justify-center">
                <svg
                  className="h-3 w-3 text-white"
                  fill="none"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <polyline points="20 6 9 17 4 12" />
                </svg>
              </div>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <Label className="text-sm font-normal">Display name</Label>
            </div>
            <div className="flex items-center gap-2">
              <Input
                type="text"
                placeholder="Enter display name"
                value={displayName}
                onChange={(e) => setDisplayName(e.target.value)}
                className="h-9 w-[200px]"
              />
            </div>
          </div>
        </div>

        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <Label className="text-sm font-normal">Password</Label>
            </div>
            <Button variant="outline" size="sm">
              Change password
            </Button>
          </div>
        </div>

        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Label className="text-sm font-normal">Two-password mode</Label>
              <Info className="h-4 w-4 text-muted-foreground" />
            </div>
            <Switch
              checked={twoFactorEnabled}
              onCheckedChange={setTwoFactorEnabled}
            />
          </div>
        </div>
      </div>

      {/* Delete Account Section */}
      <div className="space-y-4 border rounded-md p-6 border-destructive/50">
        <h3 className="font-semibold text-md text-destructive">Delete account</h3>

        <div className="space-y-4">
          <p className="text-sm text-muted-foreground">
            Once you delete your account, there is no going back. All your data will be permanently deleted. Please be certain.
          </p>

          <Button variant="destructive" className="w-full sm:w-auto">
            Delete my account
          </Button>
        </div>
      </div>
    </div>
  )
}
