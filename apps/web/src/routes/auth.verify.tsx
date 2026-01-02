import { createFileRoute } from '@tanstack/react-router'

import Verify from '@/lib/pages/verify'

export const Route = createFileRoute('/auth/verify')({
  component: Verify,
})
