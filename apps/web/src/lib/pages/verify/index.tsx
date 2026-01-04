import { useEffect, useRef } from "react";
import { useNavigate } from "@tanstack/react-router";
import { CircleDashed } from "lucide-react";
import toast from "react-hot-toast";
import { api_url } from "@/lib/constants";

const Verify = () => {
  const navigate = useNavigate();
  const hasVerified = useRef(false);

  useEffect(() => {
    if (hasVerified.current) {
      return;
    }

    const params = new URLSearchParams(window.location.search);
    const token = params.get("token");

    if (!token) {
      navigate({ to: "/" });
      return;
    }

    const verify = async () => {
      hasVerified.current = true;

      try {
        const res = await fetch(`${api_url}/auth/verify`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({ token }),
        });

        if (!res.ok) {
          const errorData = await res.json().catch(() => ({
            message: "Invalid or expired verification link",
          }));
          toast.error(
            errorData.message || "Invalid or expired verification link",
          );
          navigate({ to: "/" });
          return;
        }

        const data = await res.json();

        if (data.needsWallet) {
          toast.success(
            "Email verified! Please connect your wallet to continue.",
          );
          navigate({ to: "/auth/connect_wallet" });
        } else {
          toast.success("Login successful!");
          navigate({ to: "/app" });
        }
      } catch (error) {
        console.error("Verification error:", error);
        toast.error("Failed to verify email. Please try again.");
        navigate({ to: "/" });
      }
    };

    verify();
  }, [navigate]);

  return (
    <div className="min-h-screen w-full relative z-20 bg-white dark:bg-[#1a1625] fade-in flex flex-col antialiased overflow-x-hidden">
      <div className="flex flex-col items-center justify-center min-h-screen px-4">
        <div className="w-12 h-12 bg-gradient-primary rounded-lg flex items-center justify-center text-white">
          <CircleDashed
            className="w-6 h-6 animate-spin [animation-duration:5s]"
            strokeWidth={1.5}
          />
        </div>
        <p className="mt-4 text-sm text-[#666] dark:text-[#c4b5fd]">
          Verifying your email…
        </p>
      </div>
    </div>
  );
};

export default Verify;
