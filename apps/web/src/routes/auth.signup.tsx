import { createFileRoute } from '@tanstack/react-router'

import Signup from '@/lib/pages/signup'

export const Route = createFileRoute('/auth/signup')({
  component: Signup,
})
