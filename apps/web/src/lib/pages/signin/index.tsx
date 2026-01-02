import { useState } from 'react'
import { useNavigate, Link } from '@tanstack/react-router'
import { CircleDashed, ArrowRight } from 'lucide-react'

const Signin = () => {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleSignin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setLoading(true)

    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 800))

    setLoading(false)
    // Navigate to wallet verification after email signin
    navigate({ to: '/auth/verify' })
  }

  return (
    <div className="min-h-screen w-full relative z-20 bg-zinc-50 dark:bg-zinc-950 fade-in flex flex-col antialiased overflow-x-hidden">
      <div className="flex flex-col items-center justify-center min-h-screen bg-zinc-50 dark:bg-zinc-950 px-4 fade-in">
        <div className="w-full max-w-sm bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm rounded-xl p-8">
          <div className="mb-6 flex justify-center">
            <div className="w-12 h-12 bg-zinc-900 dark:bg-zinc-100 rounded-lg flex items-center justify-center text-white dark:text-zinc-900">
              <CircleDashed className="w-6 h-6" strokeWidth={1.5} />
            </div>
          </div>
          <h1 className="text-xl font-semibold tracking-tight text-center text-zinc-900 dark:text-zinc-100 mb-1">
            Sign in to your account
          </h1>
          <p className="text-sm text-zinc-500 dark:text-zinc-400 text-center mb-8">
            Enter your email to continue.
          </p>

          <form onSubmit={handleSignin} className="space-y-4">
            <div>
              <label className="block text-xs font-medium text-zinc-700 dark:text-zinc-300 mb-1.5 uppercase tracking-wide">
                Email Address
              </label>
              <input
                type="email"
                name="email"
                id="signin-email"
                required
                placeholder="alice@example.com"
                className="w-full px-3 py-2.5 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-zinc-900/10 dark:focus:ring-zinc-100/10 focus:border-zinc-900 dark:focus:border-zinc-100 transition-all placeholder:text-zinc-400 dark:placeholder:text-zinc-600 text-zinc-900 dark:text-zinc-100"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full mt-2 bg-zinc-900 dark:bg-zinc-100 hover:bg-zinc-800 dark:hover:bg-zinc-200 text-white dark:text-zinc-900 font-medium text-sm py-2.5 rounded-lg transition-all flex items-center justify-center gap-2 group shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span>{loading ? 'Signing in...' : 'Continue'}</span>
              <ArrowRight className="w-4 h-4 text-zinc-400 dark:text-zinc-600 group-hover:text-white dark:group-hover:text-zinc-900 transition-colors" />
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-xs text-zinc-400 dark:text-zinc-500">
              Don't have an account?{' '}
              <Link
                to="/auth/signup"
                className="text-zinc-900 dark:text-zinc-100 font-medium hover:underline"
              >
                Sign up
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default Signin
