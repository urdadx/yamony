import { AppSidebar } from '@/components/app-sidebar';
import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/admin')({
  component: RouteComponent,
  beforeLoad: ({ context, location }) => {
    if (context.authState.isLoading) {
      return;
    }

    if (!context.authState.isAuthenticated) {
      throw redirect({ to: "/login", search: location.search });
    }
  }
})

function RouteComponent() {

  return (
    <>
      <AppSidebar />

    </>
  )
}
