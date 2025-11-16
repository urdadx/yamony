import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/links')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/links"!</div>
}
