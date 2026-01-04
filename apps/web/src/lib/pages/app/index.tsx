const App = () => {
  return (
    <div className="max-w-4xl mx-auto space-y-8">
      <div>
        <h2 className="text-lg font-medium text-zinc-900 dark:text-zinc-50 tracking-tight">
          Overview
        </h2>
        <p className="text-sm text-zinc-500 dark:text-zinc-400">
          Welcome back. You have 0 pending contributions.
        </p>
      </div>

      <div className="bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 border-dashed rounded-xl p-8 text-center">
        <div className="w-12 h-12 bg-white dark:bg-zinc-950 border border-zinc-200 dark:border-zinc-800 rounded-full flex items-center justify-center text-zinc-400 dark:text-zinc-500 mx-auto mb-3 shadow-sm">
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M5 13l4 4L19 7"
            />
          </svg>
        </div>
        <h3 className="text-sm font-medium text-zinc-900 dark:text-zinc-50">
          All caught up
        </h3>
        <p className="text-xs text-zinc-500 dark:text-zinc-400 mt-1">
          No pending payments for any active groups.
        </p>
      </div>
    </div>
  )
}

export default App
