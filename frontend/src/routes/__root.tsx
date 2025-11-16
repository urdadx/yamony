import { Outlet, createRootRoute } from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'
import { ReactQueryDevtoolsPanel } from '@tanstack/react-query-devtools'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { TooltipProvider } from '@/components/ui/tooltip'
import { Toaster } from '@/components/ui/sonner'


const queryClient = new QueryClient()

export const Route = createRootRoute({
  component: () => (
    <>
      <QueryClientProvider client={queryClient}>
        <TooltipProvider>

          <Outlet />
          <Toaster richColors theme="light" />

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
        </TooltipProvider>
      </QueryClientProvider>

    </>
  ),
})
