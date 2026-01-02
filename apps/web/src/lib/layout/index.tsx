import type { ReactNode } from 'react';

import { ThemeProvider } from '@/lib/components/theme-provider';

import { Header } from './components/header';

type LayoutProps = {
  children: ReactNode;
};

export const Layout = ({ children }: LayoutProps) => {
  return (
    <ThemeProvider>
      <div 
        className="flex min-h-screen flex-col bg-zinc-50 dark:bg-zinc-950 dark:text-white overflow-x-hidden"
        style={{ overscrollBehavior: 'none' }}
      >
        <Header />
        <main className="bg-zinc-50 dark:bg-zinc-950 flex-1" style={{ overscrollBehavior: 'none' }}>{children}</main>
      </div>
    </ThemeProvider>
  );
};
