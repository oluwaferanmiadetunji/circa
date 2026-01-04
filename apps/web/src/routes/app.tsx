import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'

import { AppLayout } from '@/lib/layout/app-layout'
import { api_url } from '@/lib/constants'

export const Route = createFileRoute('/app')({
  beforeLoad: async () => {
    // Redirect to /auth/signup if user doesn't have a valid session
    try {
      const res = await fetch(`${api_url}/me`, {
        method: 'GET',
        credentials: 'include',
      })
      if (!res.ok) {
        throw redirect({ to: '/auth/signup' })
      }
    } catch (error) {
      if (error instanceof Error && error.name === 'Redirect') {
        throw error
      }
      // If check fails, redirect to signup
      throw redirect({ to: '/auth/signup' })
    }
  },
  component: () => (
    <AppLayout>
      <Outlet />
    </AppLayout>
  ),
})
