import type { ReactNode } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link, useNavigate, useRouterState } from "@tanstack/react-router";
import { CircleDashed, LayoutGrid, Menu, PlusCircle } from "lucide-react";

import { api_url } from "@/lib/constants";

type User = {
  id: string;
  address: string;
  displayName: string;
  createdAt: string;
  updatedAt?: string;
};

type AppLayoutProps = {
  children: ReactNode;
};

export const AppLayout = ({ children }: AppLayoutProps) => {
  const navigate = useNavigate();
  const router = useRouterState();
  const currentPath = router.location.pathname;

  const { data: user, isLoading } = useQuery<User>({
    queryKey: ["/me"],
    queryFn: async () => {
      const res = await fetch(`${api_url}/me`, {
        method: "GET",
        credentials: "include",
      });
      if (!res.ok) {
        throw new Error("Failed to fetch user");
      }
      return res.json();
    },
  });

  const mockGroups: Array<{
    id: string;
    name: string;
    status: "pending" | "paid";
  }> = [];

  const isActiveRoute = (path: string) => {
    return currentPath === path || currentPath.startsWith(path + "/");
  };

  return (
    <div className="flex h-screen">
      <aside className="hidden md:flex w-64 bg-zinc-50 dark:bg-zinc-950 border-r border-zinc-200 dark:border-zinc-800 flex-col justify-between pt-6 pb-4">
        <div>
          <div className="px-6 mb-8 flex items-center gap-2">
            <div className="w-6 h-6 bg-zinc-900 dark:bg-zinc-50 rounded flex items-center justify-center text-white dark:text-zinc-900">
              <CircleDashed className="w-3.5 h-3.5" />
            </div>
            <span className="font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">
              Circa
            </span>
          </div>

          <nav className="px-3 space-y-0.5">
            <Link
              to="/app"
              className={`w-full flex items-center gap-3 px-3 py-2 text-sm rounded-md transition-colors ${
                isActiveRoute("/app") && !currentPath.includes("/app/")
                  ? "bg-zinc-200/50 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 font-medium"
                  : "text-zinc-500 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-900 hover:text-zinc-900 dark:hover:text-zinc-50"
              }`}
            >
              <LayoutGrid className="w-[18px] h-[18px]" strokeWidth={1.5} />
              Dashboard
            </Link>

            <Link
              to="/app/create_circle"
              className={`w-full flex items-center gap-3 px-3 py-2 text-sm rounded-md transition-colors ${
                isActiveRoute("/app/create_group")
                  ? "bg-zinc-200/50 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 font-medium"
                  : "text-zinc-500 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-900 hover:text-zinc-900 dark:hover:text-zinc-50"
              }`}
            >
              <PlusCircle className="w-[18px] h-[18px]" strokeWidth={1.5} />
              New Circle
            </Link>

            {mockGroups.length > 0 && (
              <>
                <div className="px-3 pt-6 pb-2 text-xs font-medium text-zinc-400 dark:text-zinc-500 uppercase tracking-wider">
                  Your Groups
                </div>
                {mockGroups.map((group) => (
                  <button
                    key={group.id}
                    onClick={() => navigate({ to: `/app/groups/${group.id}` })}
                    className={`w-full flex items-center gap-3 px-3 py-2 text-sm rounded-md transition-colors ${
                      isActiveRoute(`/app/groups/${group.id}`)
                        ? "bg-zinc-100 dark:bg-zinc-900 text-zinc-900 dark:text-zinc-50 font-medium"
                        : "text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100/50 dark:hover:bg-zinc-900 hover:text-zinc-900 dark:hover:text-zinc-50"
                    }`}
                  >
                    <div
                      className={`w-1.5 h-1.5 rounded-full ${
                        group.status === "pending"
                          ? "bg-amber-500"
                          : "bg-emerald-500"
                      }`}
                    />
                    <span className="truncate">{group.name}</span>
                  </button>
                ))}
              </>
            )}
          </nav>
        </div>

        {/* User Profile Card */}
        <div className="px-4">
          <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg p-3 shadow-sm">
            {isLoading ? (
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-full bg-zinc-200 dark:bg-zinc-800 animate-pulse" />
                <div className="flex-1 space-y-2">
                  <div className="h-3 w-24 bg-zinc-200 dark:bg-zinc-800 rounded animate-pulse" />
                  <div className="h-2 w-32 bg-zinc-200 dark:bg-zinc-800 rounded animate-pulse" />
                </div>
              </div>
            ) : user ? (
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-full bg-linear-to-tr from-zinc-200 to-zinc-100 dark:from-zinc-700 dark:to-zinc-800 border border-zinc-100 dark:border-zinc-800 flex items-center justify-center text-xs font-semibold text-zinc-600 dark:text-zinc-400">
                  {user.displayName?.charAt(0).toUpperCase() || "U"}
                </div>
                <div className="flex-1 overflow-hidden">
                  <p className="text-xs font-medium text-zinc-900 dark:text-zinc-50 truncate">
                    {user.displayName || "User"}
                  </p>
                  <p className="text-xs text-zinc-500 dark:text-zinc-400 truncate font-mono">
                    {user.address
                      ? `${user.address.slice(0, 6)}...${user.address.slice(
                          -4,
                        )}`
                      : "No address"}
                  </p>
                </div>
              </div>
            ) : null}
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <main className="flex-1 flex flex-col bg-white dark:bg-zinc-950 overflow-hidden relative">
        {/* Mobile Header */}
        <header className="md:hidden h-14 bg-white dark:bg-zinc-950 border-b border-zinc-200 dark:border-zinc-800 flex items-center justify-between px-4 sticky top-0 z-20">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 bg-zinc-900 dark:bg-zinc-50 rounded flex items-center justify-center text-white dark:text-zinc-900">
              <CircleDashed className="w-3.5 h-3.5" />
            </div>
            <span className="font-semibold tracking-tight text-zinc-900 dark:text-zinc-50">
              Circa
            </span>
          </div>
          <button className="text-zinc-500 dark:text-zinc-400">
            <Menu className="w-6 h-6" />
          </button>
        </header>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-4 md:p-8 fade-in">
          {children}
        </div>
      </main>
    </div>
  );
};
