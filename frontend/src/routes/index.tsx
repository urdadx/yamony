import { Button } from '@/components/ui/button'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  return (
    <div className="text-center mt-20 text-2xl font-bold">
      Hi there! Welcome to Yamony!
      <Button>
        Click me !
      </Button>
    </div>
  )
}
