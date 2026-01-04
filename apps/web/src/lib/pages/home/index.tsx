import { useEffect, useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import {
  CircleDashed,
  ArrowRight,
  ShieldCheck,
  Zap,
  Coins,
} from "lucide-react";
import { api_url } from "@/lib/constants";

const Home = () => {
  const [hasSession, setHasSession] = useState<boolean | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const checkSession = async () => {
      try {
        const res = await fetch(`${api_url}/me`, {
          method: "GET",
          credentials: "include",
        });
        setHasSession(res.ok);
      } catch {
        setHasSession(false);
      }
    };
    checkSession();
  }, []);
  return (
    <div
      className="min-h-screen w-full relative z-20 bg-white dark:bg-[#1a1625] mesh-bg fade-in flex flex-col antialiased overflow-x-hidden"
      style={{
        overscrollBehavior: "none",
        overscrollBehaviorY: "none",
        overscrollBehaviorX: "none",
      }}
    >
      {/* Navbar */}
      <nav className="sticky top-0 z-50 w-full backdrop-blur-md bg-white/80 dark:bg-[#1a1625]/80 border-b border-[#eeeeee] dark:border-[#3d3551] py-2">
        <div className="mx-auto px-4 sm:px-6 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-gradient-primary rounded-lg flex items-center justify-center text-white shadow-sm">
              <CircleDashed className="w-[18px] h-[18px]" />
            </div>
            <span className="font-semibold text-lg tracking-tight text-[#333] dark:text-[#f5f3ff]">
              Circa
            </span>
          </div>
          <div className="flex items-center gap-6">
            {hasSession ? (
              <button
                onClick={() => navigate({ to: "/app" })}
                role="link"
                className="cursor-pointer bg-gradient-primary hover:opacity-90 text-white text-sm font-medium px-4 py-2 rounded-full transition-all shadow-sm flex items-center gap-2 group"
              >
                Go to App
                <ArrowRight className="w-3.5 h-3.5 group-hover:translate-x-0.5 transition-transform" />
              </button>
            ) : (
              <>
                <button
                  onClick={() => navigate({ to: "/auth/signin" })}
                  role="link"
                  className="bg-[#f5f5f5] cursor-pointer dark:bg-[#241f2e] hover:bg-[#eeeeee] dark:hover:bg-[#3d3551] text-[#333] dark:text-[#f5f3ff] text-sm font-medium px-4 py-2 rounded-full transition-all shadow-sm flex items-center gap-2 group"
                >
                  Sign in
                </button>
                <button
                  onClick={() => navigate({ to: "/auth/signup" })}
                  role="link"
                  className="bg-gradient-primary hover:opacity-90 text-white text-sm font-medium px-4 py-2 rounded-full transition-all shadow-sm flex items-center gap-2 group"
                >
                  Create Account
                </button>
              </>
            )}
          </div>
        </div>
      </nav>

      <main className="flex-1 flex flex-col items-center pt-20 pb-20 px-4 sm:px-6">
        <div className=" text-center mb-16 space-y-6">
          <h1 className="text-5xl md:text-7xl font-semibold tracking-tighter text-[#333] dark:text-[#f5f3ff] leading-[1.1]">
            Rotating savings,
            <br />
            <span className="text-[#666] dark:text-[#c4b5fd]">
              reinvented on-chain.
            </span>
          </h1>

          <p className="text-lg text-[#666] dark:text-[#c4b5fd] max-w-xl mx-auto font-light leading-relaxed">
            Circa creates trustless credit circles for you and your friends.
            Pool funds, automate payouts, and save towards your goals without
            the headache.
          </p>

          <div className="flex items-center justify-center gap-4 pt-4"></div>
        </div>

        <div className="perspective-container w-full max-w-5xl px-4 mb-24">
          <div className="hero-card bg-white dark:bg-[#241f2e] rounded-2xl border border-[#eeeeee] dark:border-[#3d3551] p-2 md:p-4 mx-auto max-w-4xl relative overflow-hidden">
            <div className="absolute inset-0 bg-linear-to-tr from-[#f5f5f5]/50 dark:from-[#3d3551]/50 to-transparent pointer-events-none"></div>

            <div className="flex items-center justify-between border-b border-[#eeeeee] dark:border-[#3d3551] pb-4 mb-6 px-4 pt-2">
              <div className="flex gap-2">
                <div className="w-3 h-3 rounded-full bg-[#999] dark:bg-[#a78bfa]"></div>
                <div className="w-3 h-3 rounded-full bg-[#999] dark:bg-[#a78bfa]"></div>
              </div>
              <div className="h-2 w-20 bg-[#f5f5f5] dark:bg-[#3d3551] rounded-full"></div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 px-4 pb-4">
              <div className="col-span-2 space-y-4">
                <div className="h-32 bg-[#f5f5f5] dark:bg-[#1a1625] border border-[#eeeeee] dark:border-[#3d3551] rounded-xl w-full flex flex-col justify-center p-6 space-y-3">
                  <div className="h-2 w-24 bg-[#999] dark:bg-[#a78bfa] rounded"></div>
                  <div className="h-8 w-48 bg-gradient-primary rounded-md opacity-90"></div>
                </div>
                <div className="space-y-2">
                  <div className="h-14 w-full border border-[#eeeeee] dark:border-[#3d3551] rounded-lg flex items-center px-4 justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-[#f5f5f5] dark:bg-[#3d3551]"></div>
                      <div className="h-2 w-24 bg-[#f5f5f5] dark:bg-[#3d3551] rounded"></div>
                    </div>
                    <div className="h-2 w-12 bg-emerald-100 dark:bg-emerald-900/30 rounded"></div>
                  </div>
                  <div className="h-14 w-full border border-[#eeeeee] dark:border-[#3d3551] rounded-lg flex items-center px-4 justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-[#f5f5f5] dark:bg-[#3d3551]"></div>
                      <div className="h-2 w-24 bg-[#f5f5f5] dark:bg-[#3d3551] rounded"></div>
                    </div>
                    <div className="h-2 w-12 bg-[#f5f5f5] dark:bg-[#3d3551] rounded"></div>
                  </div>
                </div>
              </div>
              <div className="space-y-4">
                <div className="h-full bg-gradient-primary rounded-xl p-6 text-white/20 flex flex-col justify-between">
                  <div className="space-y-2">
                    <div className="h-2 w-full bg-white/10 rounded"></div>
                    <div className="h-2 w-2/3 bg-white/10 rounded"></div>
                  </div>
                  <div className="h-8 w-full bg-white/10 rounded mt-8"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div
          className="max-w-6xl w-full grid grid-cols-1 md:grid-cols-3 gap-8 px-4 border-t border-[#eeeeee] dark:border-[#3d3551] pt-16"
          id="features"
        >
          <div className="space-y-3">
            <div className="w-10 h-10 rounded-lg bg-[#f5f5f5] dark:bg-[#3d3551] flex items-center justify-center text-[#333] dark:text-[#f5f3ff] border border-[#eeeeee] dark:border-[#3d3551]">
              <ShieldCheck className="w-5 h-5" />
            </div>
            <h3 className="font-medium text-[#333] dark:text-[#f5f3ff]">
              Trustless Design
            </h3>
            <p className="text-sm text-[#666] dark:text-[#c4b5fd] leading-relaxed">
              Smart contracts handle the pooling and distribution. No central
              authority holds your funds.
            </p>
          </div>
          <div className="space-y-3">
            <div className="w-10 h-10 rounded-lg bg-[#f5f5f5] dark:bg-[#3d3551] flex items-center justify-center text-[#333] dark:text-[#f5f3ff] border border-[#eeeeee] dark:border-[#3d3551]">
              <Zap className="w-5 h-5" />
            </div>
            <h3 className="font-medium text-[#333] dark:text-[#f5f3ff]">
              Automated Rounds
            </h3>
            <p className="text-sm text-[#666] dark:text-[#c4b5fd] leading-relaxed">
              Payouts are calculated and executed automatically. Notifications
              keep everyone in sync.
            </p>
          </div>
          <div className="space-y-3">
            <div className="w-10 h-10 rounded-lg bg-[#f5f5f5] dark:bg-[#3d3551] flex items-center justify-center text-[#333] dark:text-[#f5f3ff] border border-[#eeeeee] dark:border-[#3d3551]">
              <Coins className="w-5 h-5" />
            </div>
            <h3 className="font-medium text-[#333] dark:text-[#f5f3ff]">
              Stablecoin Native
            </h3>
            <p className="text-sm text-[#666] dark:text-[#c4b5fd] leading-relaxed">
              Save in USDC, DAI, or ETH. Avoid volatility while you reach your
              savings goals.
            </p>
          </div>
        </div>

        <footer className="mt-20 py-8 border-t border-[#eeeeee] dark:border-[#3d3551] w-full max-w-6xl flex justify-between items-center text-xs text-[#999] dark:text-[#a78bfa]">
          <div>&copy; {new Date().getFullYear()} Circa Finance.</div>
        </footer>
      </main>
    </div>
  );
};

export default Home;
