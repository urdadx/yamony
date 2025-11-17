import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/chatbot')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/chatbot"!</div>
}
