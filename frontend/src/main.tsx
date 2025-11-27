import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'
import './styles.css'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { TooltipProvider } from './components/ui/tooltip'
import { Toaster } from './components/ui/sonner'
import { AuthProvider, useAuth } from './context/auth-context'
import { GlobalLoader } from './components/global-loader'

const router = createRouter({
  routeTree,
  context: {
    authState: undefined!,
  },
  defaultPreload: 'intent',
  scrollRestoration: true,
  defaultStructuralSharing: true,
  defaultPreloadStaleTime: 0,
  defaultPendingComponent: () => (
    <GlobalLoader />
  ),
})

const queryClient = new QueryClient();

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

function InnerApp() {
  const { authState } = useAuth();

  if (authState.isLoading) {
    return (
      <GlobalLoader />
    );
  }

  return <RouterProvider router={router} context={{ authState }} />;
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <TooltipProvider>
          <InnerApp />
          <Toaster position='bottom-center' richColors theme='light' />
        </TooltipProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}

const rootElement = document.getElementById("app")!;
if (!rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement);
  root.render(
    <StrictMode>
      <App />
    </StrictMode>,
  );
}

