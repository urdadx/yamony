import { Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'
import { ReactQueryDevtoolsPanel } from '@tanstack/react-query-devtools'
import { Toaster } from '@/components/ui/sonner'
import type { AuthState } from '@/context/auth-context'

interface MyRouterContext {
  authState: AuthState;
}
export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: () => (
    <>

      <Outlet />
      <Toaster richColors theme="light" />

      {
        import.meta.env.VITE_ENV === 'development' && (
          <TanStackDevtools
            config={{
              position: 'bottom-right',
            }}
            plugins={[
              {
                name: 'TanStack Query',
                render: <ReactQueryDevtoolsPanel />,
                defaultOpen: true
              },
              {
                name: 'TanStack Router',
                render: <TanStackRouterDevtoolsPanel />,
                defaultOpen: false
              },
            ]}
          />
        )
      }

    </>
  ),
})
