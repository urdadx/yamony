import { AppSidebar } from '@/components/app-sidebar';
import { createFileRoute, redirect } from '@tanstack/react-router'
import { z } from 'zod'

const adminSearchSchema = z.object({
  vaultId: z.number().optional(),
})

export const Route = createFileRoute('/admin')({
  component: RouteComponent,
  validateSearch: adminSearchSchema,
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
