import type { ReactNode } from 'react'
import { Toaster } from 'react-hot-toast'

import { ThemeProvider } from '@/lib/components/theme-provider'
import { ThemeToggle } from '@/lib/components/theme-toggle'

export { AppLayout } from './app-layout'

type LayoutProps = {
  children: ReactNode
}

export const Layout = ({ children }: LayoutProps) => {
  return (
    <ThemeProvider>
      <div
        className="flex min-h-screen flex-col bg-white dark:bg-[#1a1625] text-[#333] dark:text-[#f5f3ff] overflow-x-hidden"
        style={{ overscrollBehavior: 'none' }}
      >
        <main
          className="bg-white dark:bg-[#1a1625] flex-1"
          style={{ overscrollBehavior: 'none' }}
        >
          {children}
        </main>

        <div className="fixed top-6 right-6 z-50">
          <div className="bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] rounded-full h-12 w-12 flex items-center justify-center p-2 shadow-lg hover:shadow-xl transition-shadow">
            <ThemeToggle />
          </div>
        </div>

        <Toaster
          position="top-right"
          toastOptions={{
            duration: 4000,
            className:
              '!bg-white dark:!bg-[#241f2e] !text-[#333] dark:!text-[#f5f3ff] !border !border-[#eeeeee] dark:!border-[#3d3551]',
            style: {
              borderRadius: '0.5rem',
              boxShadow:
                '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
            },
            success: {
              iconTheme: {
                primary: '#10b981',
                secondary: '#fff',
              },
            },
            error: {
              iconTheme: {
                primary: '#ef4444',
                secondary: '#fff',
              },
            },
          }}
        />
      </div>
    </ThemeProvider>
  )
}
