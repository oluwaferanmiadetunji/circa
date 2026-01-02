import { useState } from 'react'
import { useNavigate, Link } from '@tanstack/react-router'
import { CircleDashed, ArrowRight, Mail } from 'lucide-react'

type SignupView = 'signup-form' | 'verify-email'

const Signup = () => {
  const [view, setView] = useState<SignupView>('signup-form')
  const [signupEmail, setSignupEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleSignup = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setLoading(true)

    const formData = new FormData(e.currentTarget)
    const email = formData.get('email') as string
    setSignupEmail(email)

    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 800))

    setLoading(false)
    setView('verify-email')
  }

  const simulateLinkClick = () => {
    // Navigate to auth verify page after email verification
    navigate({ to: '/auth/verify' })
  }

  return (
    <div className="min-h-screen w-full relative z-20 bg-zinc-50 dark:bg-zinc-950 fade-in flex flex-col antialiased overflow-x-hidden">
      {view === 'signup-form' && (
        <SignupFormView onSubmit={handleSignup} loading={loading} />
      )}
      {view === 'verify-email' && (
        <VerifyEmailView
          email={signupEmail}
          onSimulateLink={simulateLinkClick}
          onBack={() => setView('signup-form')}
        />
      )}
    </div>
  )
}

type SignupFormViewProps = {
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void
  loading: boolean
}

const SignupFormView = ({ onSubmit, loading }: SignupFormViewProps) => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-zinc-50 dark:bg-zinc-950 px-4 fade-in">
      <div className="w-full max-w-sm bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm rounded-xl p-8">
        <div className="mb-6 flex justify-center">
          <div className="w-12 h-12 bg-zinc-900 dark:bg-zinc-100 rounded-lg flex items-center justify-center text-white dark:text-zinc-900">
            <CircleDashed className="w-6 h-6" strokeWidth={1.5} />
          </div>
        </div>
        <h1 className="text-xl font-semibold tracking-tight text-center text-zinc-900 dark:text-zinc-100 mb-1">
          Create your account
        </h1>
        <p className="text-sm text-zinc-500 dark:text-zinc-400 text-center mb-8">
          Enter your details to get started.
        </p>

        <form onSubmit={onSubmit} className="space-y-4">
          <div>
            <label className="block text-xs font-medium text-zinc-700 dark:text-zinc-300 mb-1.5 uppercase tracking-wide">
              Full Name
            </label>
            <input
              type="text"
              name="name"
              required
              placeholder="Alice Smith"
              className="w-full px-3 py-2.5 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-zinc-900/10 dark:focus:ring-zinc-100/10 focus:border-zinc-900 dark:focus:border-zinc-100 transition-all placeholder:text-zinc-400 dark:placeholder:text-zinc-600 text-zinc-900 dark:text-zinc-100"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-zinc-700 dark:text-zinc-300 mb-1.5 uppercase tracking-wide">
              Email Address
            </label>
            <input
              type="email"
              name="email"
              id="signup-email"
              required
              placeholder="alice@example.com"
              className="w-full px-3 py-2.5 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-zinc-900/10 dark:focus:ring-zinc-100/10 focus:border-zinc-900 dark:focus:border-zinc-100 transition-all placeholder:text-zinc-400 dark:placeholder:text-zinc-600 text-zinc-900 dark:text-zinc-100"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-zinc-700 dark:text-zinc-300 mb-1.5 uppercase tracking-wide">
              Display Name
            </label>
            <div className="relative">
              <span className="absolute left-3 top-2.5 text-zinc-400 dark:text-zinc-600 text-sm">
                @
              </span>
              <input
                type="text"
                name="displayName"
                required
                placeholder="alice"
                className="w-full pl-7 pr-3 py-2.5 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-zinc-900/10 dark:focus:ring-zinc-100/10 focus:border-zinc-900 dark:focus:border-zinc-100 transition-all placeholder:text-zinc-400 dark:placeholder:text-zinc-600 text-zinc-900 dark:text-zinc-100"
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full mt-2 bg-zinc-900 dark:bg-zinc-100 hover:bg-zinc-800 dark:hover:bg-zinc-200 text-white dark:text-zinc-900 font-medium text-sm py-2.5 rounded-lg transition-all flex items-center justify-center gap-2 group shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span>{loading ? 'Creating account...' : 'Continue'}</span>
            <ArrowRight className="w-4 h-4 text-zinc-400 dark:text-zinc-600 group-hover:text-white dark:group-hover:text-zinc-900 transition-colors" />
          </button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-xs text-zinc-400 dark:text-zinc-500">
            Already have an account?{' '}
            <Link
              to="/auth/signin"
              className="text-zinc-900 dark:text-zinc-100 font-medium hover:underline"
            >
              Log in
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}

type VerifyEmailViewProps = {
  email: string
  onSimulateLink: () => void
  onBack: () => void
}

const VerifyEmailView = ({
  email,
  onSimulateLink,
  onBack,
}: VerifyEmailViewProps) => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-zinc-50 dark:bg-zinc-950 px-4 fade-in">
      <div className="w-full max-w-sm bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm rounded-xl p-8 text-center">
        <div className="w-12 h-12 bg-zinc-50 dark:bg-zinc-800 border border-zinc-100 dark:border-zinc-700 rounded-full flex items-center justify-center text-zinc-900 dark:text-zinc-100 mx-auto mb-6">
          <Mail className="w-6 h-6" />
        </div>
        <h1 className="text-xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-100 mb-2">
          Check your email
        </h1>
        <p className="text-sm text-zinc-500 dark:text-zinc-400 mb-6">
          We sent a magic link to{' '}
          <span className="font-medium text-zinc-900 dark:text-zinc-100">
            {email}
          </span>
          . Click the link to complete your signup.
        </p>

        <div className="bg-zinc-50 dark:bg-zinc-800 p-4 rounded-lg border border-zinc-100 dark:border-zinc-700 mb-6">
          <p className="text-xs text-zinc-400 dark:text-zinc-500 italic mb-2">
            Simulation: Click below to emulate clicking the email link.
          </p>
          <button
            onClick={onSimulateLink}
            className="w-full bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-700 text-zinc-900 dark:text-zinc-100 font-medium text-xs py-2 rounded shadow-sm hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors"
          >
            Simulate "Magic Link" Click
          </button>
        </div>

        <button
          onClick={onBack}
          className="text-xs text-zinc-400 dark:text-zinc-500 hover:text-zinc-600 dark:hover:text-zinc-400 flex items-center justify-center gap-1 mx-auto"
        >
          <ArrowRight className="w-3 h-3 rotate-180" />
          Back to Signup
        </button>
      </div>
    </div>
  )
}

export default Signup
