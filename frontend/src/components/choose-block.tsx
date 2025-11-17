import * as React from "react"

export function ChooseBlock({
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      className="sticky top-0 h-screen w-[400px] border-l bg-white"
      {...props}
    >
      <div className="p-4">
        {/* Content goes here */}
      </div>
    </div>
  )
}
