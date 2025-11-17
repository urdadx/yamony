import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/insights')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/insights"!</div>
}
