import { useState } from "react";
import { Link } from "@tanstack/react-router";
import { CircleDashed, ArrowRight } from "lucide-react";
import toast from "react-hot-toast";
import { api_url } from "@/lib/constants";

const Signin = () => {
  const [loading, setLoading] = useState(false);

  const handleSignin = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);

    const formData = new FormData(e.currentTarget);
    const emailValue = formData.get("email") as string;

    try {
      const res = await fetch(`${api_url}/auth/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ email: emailValue }),
      });

      if (!res.ok) {
        const error = await res
          .json()
          .catch(() => ({ message: "Failed to send login link" }));
        throw new Error(error.message || "Failed to send login link");
      }

      const data = await res.json();
      toast.success(data.message || "Login link sent! Check your email.");
    } catch (error) {
      console.error("Login error:", error);
      toast.error(
        error instanceof Error
          ? error.message
          : "Failed to send login link. Please try again.",
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen w-full relative z-20 bg-white dark:bg-[#1a1625] fade-in flex flex-col antialiased overflow-x-hidden">
      <div className="flex flex-col items-center justify-center min-h-screen bg-white dark:bg-[#1a1625] px-4 fade-in">
        <div className="w-full max-w-sm bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] shadow-sm rounded-xl p-8">
          <div className="mb-6 flex justify-center">
            <div className="w-12 h-12 bg-gradient-primary rounded-lg flex items-center justify-center text-white">
              <CircleDashed className="w-6 h-6" strokeWidth={1.5} />
            </div>
          </div>
          <h1 className="text-xl font-semibold tracking-tight text-center text-[#333] dark:text-[#f5f3ff] mb-1">
            Sign in to your account
          </h1>
          <p className="text-sm text-[#666] dark:text-[#c4b5fd] text-center mb-8">
            Enter your email to continue.
          </p>

          <form onSubmit={handleSignin} className="space-y-4">
            <div>
              <label className="block text-xs font-medium text-[#666] dark:text-[#c4b5fd] mb-1.5 uppercase tracking-wide">
                Email Address
              </label>
              <input
                type="email"
                name="email"
                id="signin-email"
                required
                placeholder="alice@example.com"
                className="w-full px-3 py-2.5 bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#667eea]/20 dark:focus:ring-[#c4b5fd]/20 focus:border-[#667eea] dark:focus:border-[#c4b5fd] transition-all placeholder:text-[#999] dark:placeholder:text-[#a78bfa] text-[#333] dark:text-[#f5f3ff]"
              />
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full cursor-pointer mt-2 bg-gradient-primary hover:opacity-90 text-white font-medium text-sm py-2.5 rounded-lg transition-all flex items-center justify-center gap-2 group shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span>{loading ? "Signing in..." : "Continue"}</span>
              <ArrowRight className="w-4 h-4 text-white/80 group-hover:text-white transition-colors" />
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-xs text-[#999] dark:text-[#a78bfa]">
              Don't have an account?{" "}
              <Link
                to="/auth/signup"
                className="text-[#333] dark:text-[#f5f3ff] font-medium hover:underline"
              >
                Sign up
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Signin;
