import { CircleDashed, ArrowRight, ShieldCheck, Zap, Coins } from 'lucide-react'

const Home = () => {
  return (
    <div
      className="min-h-screen w-full relative z-20 bg-zinc-50 dark:bg-zinc-950 mesh-bg fade-in flex flex-col antialiased overflow-x-hidden"
      style={{
        overscrollBehavior: 'none',
        overscrollBehaviorY: 'none',
        overscrollBehaviorX: 'none',
      }}
    >
      {/* Navbar */}
      <nav className="sticky top-0 z-50 w-full backdrop-blur-md bg-zinc-50/80 dark:bg-zinc-950/80 border-b border-zinc-200/50 dark:border-zinc-800/50">
        <div className="mx-auto px-4 sm:px-6 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-zinc-900 dark:bg-zinc-100 rounded-lg flex items-center justify-center text-white dark:text-zinc-900 shadow-sm">
              <CircleDashed className="w-[18px] h-[18px]" />
            </div>
            <span className="font-semibold text-lg tracking-tight text-zinc-900 dark:text-zinc-100">
              Circa
            </span>
          </div>
          <div className="flex items-center gap-6">
            <div className="hidden md:flex items-center gap-6 text-sm font-medium text-zinc-500 dark:text-zinc-400">
              <a
                href="#features"
                className="hover:text-zinc-900 dark:hover:text-zinc-100 transition-colors"
              >
                Features
              </a>
            </div>
            <button className="bg-zinc-900 dark:bg-zinc-100 hover:bg-zinc-800 dark:hover:bg-zinc-200 text-white dark:text-zinc-900 text-sm font-medium px-4 py-2 rounded-full transition-all shadow-sm flex items-center gap-2 group">
              Go to App
              <ArrowRight className="w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" />
            </button>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <main className="flex-1 flex flex-col items-center pt-20 pb-20 px-4 sm:px-6">
        <div className=" text-center mb-16 space-y-6">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm mb-4">
            <span className="flex h-2 w-2 rounded-full bg-emerald-500"></span>
            <span className="text-xs font-medium text-zinc-600 dark:text-zinc-400">
              v2.0 is now live on mainnet
            </span>
          </div>

          <h1 className="text-5xl md:text-7xl font-semibold tracking-tighter text-zinc-900 dark:text-zinc-100 leading-[1.1]">
            Rotating savings,
            <br />
            <span className="text-zinc-400 dark:text-zinc-500">
              reinvented on-chain.
            </span>
          </h1>

          <p className="text-lg text-zinc-500 dark:text-zinc-400 max-w-xl mx-auto font-light leading-relaxed">
            Circa creates trustless credit circles for you and your friends.
            Pool funds, automate payouts, and save towards your goals without
            the headache.
          </p>

          <div className="flex items-center justify-center gap-4 pt-4">
            <button className="bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 px-8 py-3.5 rounded-full font-medium text-sm hover:scale-105 transition-transform duration-200 shadow-xl shadow-zinc-900/10">
              Start a Circle
            </button>
            <button className="bg-white dark:bg-zinc-900 text-zinc-900 dark:text-zinc-100 border border-zinc-200 dark:border-zinc-800 px-8 py-3.5 rounded-full font-medium text-sm hover:bg-zinc-50 dark:hover:bg-zinc-800 transition-colors">
              Read the Docs
            </button>
          </div>
        </div>

        {/* UI Visualization */}
        <div className="perspective-container w-full max-w-5xl px-4 mb-24">
          <div className="hero-card bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 p-2 md:p-4 mx-auto max-w-4xl relative overflow-hidden">
            <div className="absolute inset-0 bg-linear-to-tr from-zinc-100/50 dark:from-zinc-800/50 to-transparent pointer-events-none"></div>
            {/* Mock App Header */}
            <div className="flex items-center justify-between border-b border-zinc-100 dark:border-zinc-800 pb-4 mb-6 px-4 pt-2">
              <div className="flex gap-2">
                <div className="w-3 h-3 rounded-full bg-zinc-200 dark:bg-zinc-700"></div>
                <div className="w-3 h-3 rounded-full bg-zinc-200 dark:bg-zinc-700"></div>
              </div>
              <div className="h-2 w-20 bg-zinc-100 dark:bg-zinc-800 rounded-full"></div>
            </div>
            {/* Mock Content Grid */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 px-4 pb-4">
              <div className="col-span-2 space-y-4">
                <div className="h-32 bg-zinc-50 dark:bg-zinc-950 border border-zinc-100 dark:border-zinc-800 rounded-xl w-full flex flex-col justify-center p-6 space-y-3">
                  <div className="h-2 w-24 bg-zinc-200 dark:bg-zinc-700 rounded"></div>
                  <div className="h-8 w-48 bg-zinc-900 dark:bg-zinc-100 rounded-md opacity-90"></div>
                </div>
                <div className="space-y-2">
                  <div className="h-14 w-full border border-zinc-100 dark:border-zinc-800 rounded-lg flex items-center px-4 justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-zinc-100 dark:bg-zinc-800"></div>
                      <div className="h-2 w-24 bg-zinc-100 dark:bg-zinc-800 rounded"></div>
                    </div>
                    <div className="h-2 w-12 bg-emerald-100 dark:bg-emerald-900/30 rounded"></div>
                  </div>
                  <div className="h-14 w-full border border-zinc-100 dark:border-zinc-800 rounded-lg flex items-center px-4 justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-zinc-100 dark:bg-zinc-800"></div>
                      <div className="h-2 w-24 bg-zinc-100 dark:bg-zinc-800 rounded"></div>
                    </div>
                    <div className="h-2 w-12 bg-zinc-100 dark:bg-zinc-800 rounded"></div>
                  </div>
                </div>
              </div>
              <div className="space-y-4">
                <div className="h-full bg-zinc-900 dark:bg-zinc-100 rounded-xl p-6 text-white/20 dark:text-zinc-900/20 flex flex-col justify-between">
                  <div className="space-y-2">
                    <div className="h-2 w-full bg-white/10 dark:bg-zinc-900/10 rounded"></div>
                    <div className="h-2 w-2/3 bg-white/10 dark:bg-zinc-900/10 rounded"></div>
                  </div>
                  <div className="h-8 w-full bg-white/10 dark:bg-zinc-900/10 rounded mt-8"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Features Grid */}
        <div
          className="max-w-6xl w-full grid grid-cols-1 md:grid-cols-3 gap-8 px-4 border-t border-zinc-200 dark:border-zinc-800 pt-16"
          id="features"
        >
          <div className="space-y-3">
            <div className="w-10 h-10 rounded-lg bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center text-zinc-900 dark:text-zinc-100 border border-zinc-200 dark:border-zinc-700">
              <ShieldCheck className="w-5 h-5" />
            </div>
            <h3 className="font-medium text-zinc-900 dark:text-zinc-100">
              Trustless Design
            </h3>
            <p className="text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed">
              Smart contracts handle the pooling and distribution. No central
              authority holds your funds.
            </p>
          </div>
          <div className="space-y-3">
            <div className="w-10 h-10 rounded-lg bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center text-zinc-900 dark:text-zinc-100 border border-zinc-200 dark:border-zinc-700">
              <Zap className="w-5 h-5" />
            </div>
            <h3 className="font-medium text-zinc-900 dark:text-zinc-100">
              Automated Rounds
            </h3>
            <p className="text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed">
              Payouts are calculated and executed automatically. Notifications
              keep everyone in sync.
            </p>
          </div>
          <div className="space-y-3">
            <div className="w-10 h-10 rounded-lg bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center text-zinc-900 dark:text-zinc-100 border border-zinc-200 dark:border-zinc-700">
              <Coins className="w-5 h-5" />
            </div>
            <h3 className="font-medium text-zinc-900 dark:text-zinc-100">
              Stablecoin Native
            </h3>
            <p className="text-sm text-zinc-500 dark:text-zinc-400 leading-relaxed">
              Save in USDC, DAI, or ETH. Avoid volatility while you reach your
              savings goals.
            </p>
          </div>
        </div>

        <footer className="mt-20 py-8 border-t border-zinc-200 dark:border-zinc-800 w-full max-w-6xl flex justify-between items-center text-xs text-zinc-400 dark:text-zinc-500">
          <div>&copy; 2026 Circa Finance.</div>
          <div className="flex gap-4">
            <a
              href="#"
              className="hover:text-zinc-900 dark:hover:text-zinc-100"
            >
              Twitter
            </a>
            <a
              href="#"
              className="hover:text-zinc-900 dark:hover:text-zinc-100"
            >
              Github
            </a>
            <a
              href="#"
              className="hover:text-zinc-900 dark:hover:text-zinc-100"
            >
              Discord
            </a>
          </div>
        </footer>
      </main>
    </div>
  )
}

export default Home
