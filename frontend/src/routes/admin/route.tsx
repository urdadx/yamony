import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/admin')({
  component: RouteComponent,
  beforeLoad: ({ context, location }) => {
    const isAuthenticated = context.authState.isAuthenticated
    if (!isAuthenticated) {
      throw redirect({ to: "/login", search: location.search });

    }
  }
})

function RouteComponent() {

  return <div>Hello admin</div>
}
