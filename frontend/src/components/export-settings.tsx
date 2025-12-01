import { Label } from "@/components/ui/label"
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { useState } from "react"
import { Eye, EyeOff } from "lucide-react"

export function ExportSettings() {
  const [fileFormat, setFileFormat] = useState("pgp")
  const [passphrase, setPassphrase] = useState("")
  const [showPassphrase, setShowPassphrase] = useState(false)

  const handleExport = () => {
    // Handle export logic here
    console.log("Exporting with format:", fileFormat)
  }

  return (
    <div className="space-y-8 py-3">
      <div className="space-y-4 border rounded-md p-6">
        <h3 className="font-semibold text-base">File format</h3>

        <RadioGroup value={fileFormat} onValueChange={setFileFormat}>
          <div className="space-y-4 pb-4">
            <div className="flex items-start gap-3">
              <RadioGroupItem value="pgp" id="pgp" className="mt-0.5" />
              <div className="flex flex-col gap-1">
                <Label htmlFor="pgp" className="font-medium cursor-pointer">
                  PGP-encrypted (recommended)
                </Label>
                <p className="text-sm text-muted-foreground">
                  Export your data in an encrypted format that requires a passphrase to access
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <RadioGroupItem value="zip" id="zip" className="mt-0.5" />
              <div className="flex flex-col gap-1">
                <Label htmlFor="zip" className="font-medium cursor-pointer">
                  ZIP
                </Label>
                <p className="text-sm text-muted-foreground">
                  Export your data as a compressed ZIP archive
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <RadioGroupItem value="csv" id="csv" className="mt-0.5" />
              <div className="flex flex-col gap-1">
                <Label htmlFor="csv" className="font-medium cursor-pointer">
                  CSV
                </Label>
                <p className="text-sm text-muted-foreground">
                  Export your data as a comma-separated values file
                </p>
              </div>
            </div>
          </div>
        </RadioGroup>
        {fileFormat === "pgp" && (
          <div className="space-y-4 border-t pt-4">
            <h3 className="font-semibold text-base">Passphrase</h3>
            <p className="text-sm text-muted-foreground -mt-2">
              The exported file will be encrypted using PGP and requires a strong passphrase.
            </p>

            <div className="relative">
              <Input
                type={showPassphrase ? "text" : "password"}
                value={passphrase}
                onChange={(e) => setPassphrase(e.target.value)}
                placeholder="Enter a strong passphrase"
                className="pr-10"
              />
              <button
                type="button"
                onClick={() => setShowPassphrase(!showPassphrase)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              >
                {showPassphrase ? (
                  <EyeOff className="h-4 w-4" />
                ) : (
                  <Eye className="h-4 w-4" />
                )}
              </button>
            </div>
          </div>
        )}
        <div className="flex justify-end">
          <Button
            onClick={handleExport}
            disabled={fileFormat === "pgp" && !passphrase}
            className="min-w-32 w-full"
          >
            Export
          </Button>
        </div>
      </div>




    </div>
  )
}
