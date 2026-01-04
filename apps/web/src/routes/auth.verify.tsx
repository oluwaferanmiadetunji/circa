import { createFileRoute, redirect } from '@tanstack/react-router'

import Verify from '@/lib/pages/verify'
import { api_url } from '@/lib/constants'

export const Route = createFileRoute('/auth/verify')({
  beforeLoad: async () => {
    // Redirect to /app if user already has a valid session
    const res = await fetch(`${api_url}/me`, {
      method: 'GET',
      credentials: 'include',
    }).catch(() => {
      // Network error - allow access (user might not be logged in)
      return { ok: false } as Response
    })
    
    if (res.ok) {
      throw redirect({ to: '/app' })
    }
  },
  component: Verify,
})
