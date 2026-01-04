import { createFileRoute } from '@tanstack/react-router'

import App from '@/lib/pages/app'

export const Route = createFileRoute('/app/')({
  component: App,
})
