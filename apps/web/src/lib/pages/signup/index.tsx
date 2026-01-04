import { useState, useEffect } from 'react'
import { Link } from '@tanstack/react-router'
import { CircleDashed, ArrowRight, Mail } from 'lucide-react'
import toast from 'react-hot-toast'
import { api_url } from '@/lib/constants'

type SignupView = 'signup-form' | 'verify-email'

const Signup = () => {
  const [view, setView] = useState<SignupView>('signup-form')
  const [signupEmail, setSignupEmail] = useState('')
  const [loading, setLoading] = useState(false)

  // Get saved form data for pre-filling
  const [savedFormData, setSavedFormData] = useState<{
    fullName?: string
    email?: string
    displayName?: string
  } | null>(null)

  useEffect(() => {
    const savedData = localStorage.getItem('circa_signup_data')
    if (savedData) {
      try {
        const data = JSON.parse(savedData)
        // Check if data is recent (within 24 hours)
        if (
          data.timestamp &&
          Date.now() - data.timestamp < 24 * 60 * 60 * 1000
        ) {
          setSavedFormData({
            fullName: data.fullName,
            email: data.email,
            displayName: data.displayName,
          })
          if (data.email) {
            setSignupEmail(data.email)
          }
        } else {
          // Clear old data
          localStorage.removeItem('circa_signup_data')
        }
      } catch {
        localStorage.removeItem('circa_signup_data')
      }
    }
  }, [])

  const handleSignup = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setLoading(true)

    const formData = new FormData(e.currentTarget)
    const fullName = formData.get('name') as string
    const email = formData.get('email') as string
    const displayName = formData.get('displayName') as string
    setSignupEmail(email)

    // Save form data to localStorage
    localStorage.setItem(
      'circa_signup_data',
      JSON.stringify({
        fullName,
        email,
        displayName,
        timestamp: Date.now(),
      }),
    )

    try {
      const res = await fetch(`${api_url}/auth/signup`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          full_name: fullName,
          email,
          display_name: displayName,
        }),
      })

      if (!res.ok) {
        const error = await res
          .json()
          .catch(() => ({ message: 'Failed to create account' }))
        throw new Error(error.message || 'Failed to create account')
      }

      const data = await res.json()
      toast.success(
        data.message ||
          'Signup successful! Check your email for the magic link.',
      )
      // Keep data in localStorage until signup is fully complete
      setView('verify-email')
    } catch (error) {
      console.error('Signup error:', error)
      toast.error(
        error instanceof Error
          ? error.message
          : 'Failed to create account. Please try again.',
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen w-full relative z-20 bg-white dark:bg-[#1a1625] fade-in flex flex-col antialiased overflow-x-hidden">
      {view === 'signup-form' && (
        <SignupFormView
          onSubmit={handleSignup}
          loading={loading}
          initialData={savedFormData}
        />
      )}
      {view === 'verify-email' && (
        <VerifyEmailView
          email={signupEmail}
          onBack={() => setView('signup-form')}
        />
      )}
    </div>
  )
}

type SignupFormViewProps = {
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void
  loading: boolean
  initialData?: {
    fullName?: string
    email?: string
    displayName?: string
  } | null
}

const SignupFormView = ({
  onSubmit,
  loading,
  initialData,
}: SignupFormViewProps) => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-white dark:bg-[#1a1625] px-4 fade-in">
      <div className="w-full max-w-sm bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] shadow-sm rounded-xl p-8">
        <div className="mb-6 flex justify-center">
          <div className="w-12 h-12 bg-gradient-primary rounded-lg flex items-center justify-center text-white">
            <CircleDashed className="w-6 h-6" strokeWidth={1.5} />
          </div>
        </div>
        <h1 className="text-xl font-semibold tracking-tight text-center text-[#333] dark:text-[#f5f3ff] mb-1">
          Create your account
        </h1>
        <p className="text-sm text-[#666] dark:text-[#c4b5fd] text-center mb-8">
          Enter your details to get started.
        </p>

        <form onSubmit={onSubmit} className="space-y-4">
          <div>
            <label className="block text-xs font-medium text-[#666] dark:text-[#c4b5fd] mb-1.5 uppercase tracking-wide">
              Full Name
            </label>
            <input
              type="text"
              name="name"
              required
              placeholder="Alice Smith"
              defaultValue={initialData?.fullName || ''}
              className="w-full px-3 py-2.5 bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#667eea]/20 dark:focus:ring-[#c4b5fd]/20 focus:border-[#667eea] dark:focus:border-[#c4b5fd] transition-all placeholder:text-[#999] dark:placeholder:text-[#a78bfa] text-[#333] dark:text-[#f5f3ff]"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-[#666] dark:text-[#c4b5fd] mb-1.5 uppercase tracking-wide">
              Email Address
            </label>
            <input
              type="email"
              name="email"
              id="signup-email"
              required
              placeholder="alice@example.com"
              defaultValue={initialData?.email || ''}
              className="w-full px-3 py-2.5 bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#667eea]/20 dark:focus:ring-[#c4b5fd]/20 focus:border-[#667eea] dark:focus:border-[#c4b5fd] transition-all placeholder:text-[#999] dark:placeholder:text-[#a78bfa] text-[#333] dark:text-[#f5f3ff]"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-[#666] dark:text-[#c4b5fd] mb-1.5 uppercase tracking-wide">
              Display Name
            </label>
            <div className="relative">
              <span className="absolute left-3 top-2.5 text-[#999] dark:text-[#a78bfa] text-sm">
                @
              </span>
              <input
                type="text"
                name="displayName"
                required
                placeholder="alice"
                defaultValue={initialData?.displayName || ''}
                className="w-full pl-7 pr-3 py-2.5 bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#667eea]/20 dark:focus:ring-[#c4b5fd]/20 focus:border-[#667eea] dark:focus:border-[#c4b5fd] transition-all placeholder:text-[#999] dark:placeholder:text-[#a78bfa] text-[#333] dark:text-[#f5f3ff]"
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            role="button"
            className="w-full mt-2 cursor-pointer bg-gradient-primary hover:opacity-90 text-white font-medium text-sm py-2.5 rounded-lg transition-all flex items-center justify-center gap-2 group shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span>{loading ? 'Creating account...' : 'Continue'}</span>
            <ArrowRight className="w-4 h-4 text-white/80 group-hover:text-white transition-colors" />
          </button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-xs text-[#999] dark:text-[#a78bfa]">
            Already have an account?{' '}
            <Link
              to="/auth/signin"
              className="text-[#333] dark:text-[#f5f3ff] font-medium hover:underline"
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
  onBack: () => void
}

const VerifyEmailView = ({ email, onBack }: VerifyEmailViewProps) => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-white dark:bg-[#1a1625] px-4 fade-in">
      <div className="w-full max-w-sm bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] shadow-sm rounded-xl p-8 text-center">
        <div className="w-12 h-12 bg-[#f5f5f5] dark:bg-[#3d3551] border border-[#eeeeee] dark:border-[#3d3551] rounded-full flex items-center justify-center text-[#333] dark:text-[#f5f3ff] mx-auto mb-6">
          <Mail className="w-6 h-6" />
        </div>
        <h1 className="text-xl font-semibold tracking-tight text-[#333] dark:text-[#f5f3ff] mb-2">
          Check your email
        </h1>
        <p className="text-sm text-[#666] dark:text-[#c4b5fd] mb-6">
          We sent a magic link to{' '}
          <span className="font-medium text-[#333] dark:text-[#f5f3ff]">
            {email}
          </span>
          . Click the link to complete your signup.
        </p>

        <button
          onClick={onBack}
          className="text-xs text-[#999] dark:text-[#a78bfa] hover:text-[#666] dark:hover:text-[#c4b5fd] flex items-center justify-center gap-1 mx-auto"
        >
          <ArrowRight className="w-3 h-3 rotate-180" />
          Back to Signup
        </button>
      </div>
    </div>
  )
}

export default Signup
