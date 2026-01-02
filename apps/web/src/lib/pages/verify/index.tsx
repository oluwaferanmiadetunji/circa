import { useState } from 'react'
import { useNavigate } from '@tanstack/react-router'
import { CircleDashed, Wallet } from 'lucide-react'

const Verify = () => {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleWalletConnect = async () => {
    setLoading(true)
    // Simulate wallet connection
    await new Promise((resolve) => setTimeout(resolve, 800))
    setLoading(false)
    // Navigate to dashboard or home after successful connection
    navigate({ to: '/' })
  }

  return (
    <div className="min-h-screen w-full relative z-20 bg-zinc-50 dark:bg-zinc-950 fade-in flex flex-col antialiased overflow-x-hidden">
      <LoginWalletView onConnect={handleWalletConnect} loading={loading} />
    </div>
  )
}

type LoginWalletViewProps = {
  onConnect: () => void
  loading: boolean
}

const LoginWalletView = ({ onConnect, loading }: LoginWalletViewProps) => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-zinc-50 dark:bg-zinc-950 px-4 fade-in">
      <div className="w-full max-w-sm bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 shadow-sm rounded-xl p-8 text-center">
        <div className="mb-6 flex justify-center">
          <div className="w-12 h-12 bg-zinc-900 dark:bg-zinc-100 rounded-lg flex items-center justify-center text-white dark:text-zinc-900">
            <CircleDashed className="w-6 h-6" strokeWidth={1.5} />
          </div>
        </div>
        <h1 className="text-xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-100 mb-2">
          Connect Wallet
        </h1>
        <p className="text-sm text-zinc-500 dark:text-zinc-400 mb-8">
          Connect your wallet to verify your identity.
        </p>

        <button
          onClick={onConnect}
          disabled={loading}
          className="w-full bg-zinc-900 dark:bg-zinc-100 hover:bg-zinc-800 dark:hover:bg-zinc-200 text-white dark:text-zinc-900 font-medium text-sm py-2.5 rounded-lg transition-all flex items-center justify-center gap-2 group shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span>{loading ? 'Connecting...' : 'Connect Wallet'}</span>
          <Wallet className="w-4 h-4 text-zinc-400 dark:text-zinc-600 group-hover:text-white dark:group-hover:text-zinc-900 transition-colors" />
        </button>
        <p className="mt-4 text-xs text-zinc-400 dark:text-zinc-500">
          By connecting, you agree to our Terms of Service.
        </p>
      </div>
    </div>
  )
}

export default Verify

