import { createFileRoute } from '@tanstack/react-router'

import Signin from '@/lib/pages/signin'

export const Route = createFileRoute('/auth/signin')({
  component: Signin,
})
