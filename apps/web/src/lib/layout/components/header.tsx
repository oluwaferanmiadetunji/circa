import { ThemeToggle } from '@/lib/components/theme-toggle';

export const Header = () => {
  return (
    <header className="bg-zinc-50/80 dark:bg-zinc-950/80 sticky top-0 z-10 w-full backdrop-blur-md border-b border-zinc-200/50 dark:border-zinc-800/50">
      <section className="wrapper mx-auto flex items-center justify-between py-2">
        <div className="ml-auto">
          <ThemeToggle />
        </div>
      </section>
    </header>
  );
};
