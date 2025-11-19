import * as React from "react"

export function Playground({
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      className="sticky top-0 h-screen "
      {...props}
    >
      <div className="p-4">
      </div>
    </div>
  )
}
