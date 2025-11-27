import * as React from "react"
import { BoltIcon } from "@/assets/icons/bolt-icon"
import { ShieldStrongIcon } from "@/assets/icons/shield-strong-icon"
import { ShieldWeakIcon } from "@/assets/icons/shield-weak-icon"
import { ShieldVulnerableIcon } from "@/assets/icons/shield-vulnerable.tsx"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Slider } from "@/components/ui/slider"
import { Checkbox } from "@/components/ui/checkbox"
import { useIsMobile } from "@/hooks/use-mobile"
import { cn } from "@/lib/utils"
import { Copy, Check } from "lucide-react"
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetTrigger,
} from "@/components/ui/sheet"
import { toast } from "sonner"

export function PasswordGenerator() {
  const isMobile = useIsMobile()
  const [triggerGenerate, setTriggerGenerate] = React.useState(0)

  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button variant="outline">
          <BoltIcon color="#888888" />
          Generate password
        </Button>
      </SheetTrigger>
      <SheetContent side={isMobile ? "bottom" : "right"} className="overflow-y-auto w-full sm:max-w-[550px]">
        <div className="flex px-4 pt-3 p-2 items-center justify-between">
          <SheetClose />
          <Button onClick={() => setTriggerGenerate(prev => prev + 1)} >
            Generate password
          </Button>
        </div>
        <PasswordGeneratorForm className="px-5 py-0" triggerGenerate={triggerGenerate} />
      </SheetContent>
    </Sheet>
  )
}

interface PasswordGeneratorFormProps extends React.ComponentProps<"div"> {
  onGenerate?: () => void
  triggerGenerate?: number
}

function PasswordGeneratorForm({ className, onGenerate, triggerGenerate }: PasswordGeneratorFormProps) {
  const [password, setPassword] = React.useState("")
  const [length, setLength] = React.useState([16])
  const [includeUppercase, setIncludeUppercase] = React.useState(true)
  const [includeLowercase, setIncludeLowercase] = React.useState(true)
  const [includeNumbers, setIncludeNumbers] = React.useState(true)
  const [includeSymbols, setIncludeSymbols] = React.useState(true)
  const [copied, setCopied] = React.useState(false)

  const generatePassword = React.useCallback(() => {
    const uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    const lowercase = "abcdefghijklmnopqrstuvwxyz"
    const numbers = "0123456789"
    const symbols = "!@#$%^&*()_+-=[]{}|;:,.<>?"

    let charset = ""
    if (includeUppercase) charset += uppercase
    if (includeLowercase) charset += lowercase
    if (includeNumbers) charset += numbers
    if (includeSymbols) charset += symbols

    if (charset === "") {
      setPassword("")
      return
    }

    let newPassword = ""
    for (let i = 0; i < length[0]; i++) {
      newPassword += charset.charAt(Math.floor(Math.random() * charset.length))
    }
    setPassword(newPassword)
    setCopied(false)
    onGenerate?.()
  }, [length, includeUppercase, includeLowercase, includeNumbers, includeSymbols, onGenerate])

  React.useEffect(() => {
    generatePassword()
  }, [generatePassword])

  React.useEffect(() => {
    if (triggerGenerate && triggerGenerate > 0) {
      generatePassword()
    }
  }, [triggerGenerate, generatePassword])

  const getPasswordStrength = () => {
    if (!password) return {
      label: "None",
      icon: null,
      color: "text-gray-400"
    }

    let strength = 0
    if (password.length >= 8) strength++
    if (password.length >= 12) strength++
    if (password.length >= 16) strength++
    if (includeUppercase && /[A-Z]/.test(password)) strength++
    if (includeLowercase && /[a-z]/.test(password)) strength++
    if (includeNumbers && /[0-9]/.test(password)) strength++
    if (includeSymbols && /[^A-Za-z0-9]/.test(password)) strength++

    if (strength <= 2) return {
      label: "Vulnerable",
      icon: ShieldVulnerableIcon,
      color: "text-red-500"
    }
    if (strength <= 4) return {
      label: "Weak",
      icon: ShieldWeakIcon,
      color: "text-orange-500"
    }
    return {
      label: "Strong",
      icon: ShieldStrongIcon,
      color: "text-green-500"
    }
  }

  const strength = getPasswordStrength()

  const copyToClipboard = async () => {
    if (!password) return
    await navigator.clipboard.writeText(password)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
    toast.success("Password copied to clipboard!")
  }

  return (
    <div className={cn("grid items-start gap-5", className)}>
      <div className="grid gap-3">
        <Label htmlFor="generated-password">Generated Password</Label>
        <div className="relative">
          <Input
            id="generated-password"
            value={password}
            readOnly
            className="pr-10 font-mono text-lg"
          />
          <Button
            type="button"
            variant="ghost"
            size="icon"
            className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7"
            onClick={copyToClipboard}
          >
            {copied ? (
              <Check className="h-4 w-4 text-green-500" />
            ) : (
              <Copy className="h-4 w-4" />
            )}
          </Button>
        </div>
      </div>

      <div className="grid gap-2">
        <div className="flex items-center justify-center">
          <div className="flex items-center gap-2">
            {strength.icon && (
              <strength.icon className={cn("h-5 w-5", strength.color)} color="currentColor" />
            )}
            <span className={cn("text-sm font-medium", strength.color)}>
              {strength.label}
            </span>
          </div>
        </div>
      </div>

      <div className="grid gap-3">
        <div className="flex items-center justify-between">
          <Label htmlFor="password-length">Character Length</Label>
          <span className="text-sm font-medium">{length[0]}</span>
        </div>
        <Slider
          id="password-length"
          min={4}
          max={64}
          step={1}
          value={length}
          onValueChange={setLength}
        />
      </div>

      <div className="grid gap-6 pt-2">
        <div className="flex items-center space-x-3">
          <Checkbox
            id="uppercase"
            checked={includeUppercase}
            onCheckedChange={(checked) => setIncludeUppercase(checked as boolean)}
          />
          <Label htmlFor="uppercase" className="cursor-pointer font-normal">
            Include Uppercase Letters (A-Z)
          </Label>
        </div>

        <div className="flex items-center space-x-3">
          <Checkbox
            id="lowercase"
            checked={includeLowercase}
            onCheckedChange={(checked) => setIncludeLowercase(checked as boolean)}
          />
          <Label htmlFor="lowercase" className="cursor-pointer font-normal">
            Include Lowercase Letters (a-z)
          </Label>
        </div>

        <div className="flex items-center space-x-3">
          <Checkbox
            id="numbers"
            checked={includeNumbers}
            onCheckedChange={(checked) => setIncludeNumbers(checked as boolean)}
          />
          <Label htmlFor="numbers" className="cursor-pointer font-normal">
            Include Numbers (0-9)
          </Label>
        </div>

        <div className="flex items-center space-x-3">
          <Checkbox
            id="symbols"
            checked={includeSymbols}
            onCheckedChange={(checked) => setIncludeSymbols(checked as boolean)}
          />
          <Label htmlFor="symbols" className="cursor-pointer font-normal">
            Include Symbols (!@#$%^&*)
          </Label>
        </div>
      </div>
    </div>
  )
}
