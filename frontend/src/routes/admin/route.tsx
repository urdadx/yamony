import { useSession } from '@/hooks/use-session'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin')({
  component: RouteComponent,
  beforeLoad: () => {

  }
})

function RouteComponent() {

  const { currentSession } = useSession()
  console.log("current session: ", currentSession)

  return <div>Hello admin</div>
}
