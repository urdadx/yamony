import chromeIcon from "@/assets/image-icons/chrome.png"
import firefoxIcon from "@/assets/image-icons/firefox.png"
import edgeIcon from "@/assets/image-icons/edge.png"
import braveIcon from "@/assets/image-icons/icons8-brave-web-browser-48.png"
import bitwardenIcon from "@/assets/image-icons/bitwarden.png"

const passwordManagers = [
  {
    id: "chrome",
    name: "Chrome",
    icon: chromeIcon,
    formats: "csv",
  },
  {
    id: "firefox",
    name: "Firefox",
    icon: firefoxIcon,
    formats: "csv",
  },
  {
    id: "edge",
    name: "Edge",
    icon: edgeIcon,
    formats: "csv",
  },
  {
    id: "brave",
    name: "Brave",
    icon: braveIcon,
    formats: "csv",
  },
  {
    id: "bitwarden",
    name: "Bitwarden",
    icon: bitwardenIcon,
    formats: "json, zip",
  },
]

export function ImportSettings() {
  return (
    <div className="space-y-6 px-6 py-4 border rounded-md">
      <div className="space-y-4">
        <h3 className="font-semibold text-base">Import</h3>
        <p className="text-sm text-muted-foreground ">
          To migrate data from another password manager, go to the password manager, export your data, then upload it to Nucleo Pass.
        </p>
      </div>
      <div className="space-y-4">
        <h3 className="font-semibold text-base">Select your password manager</h3>

        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {passwordManagers.map((manager) => (
            <button
              key={manager.id}
              className="flex flex-col items-center gap-3 p-3 border rounded-lg hover:bg-accent hover:border-primary/50 transition-colors"
            >
              <img
                src={manager.icon}
                alt={manager.name}
                loading="lazy"
                className="h-9 w-9 object-contain"
              />
              <div className="text-center">
                <p className="font-medium text-sm">{manager.name}</p>
                <p className="text-xs text-muted-foreground">{manager.formats}</p>
              </div>
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}
