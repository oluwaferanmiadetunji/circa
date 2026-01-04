import { useState, useRef } from "react";
import { useNavigate } from "@tanstack/react-router";
import { CheckCircle2, CircleDashed, Wallet } from "lucide-react";
import toast from "react-hot-toast";
import { api_url } from "@/lib/constants";

interface EthereumProvider {
  request: (args: { method: string; params?: unknown[] }) => Promise<unknown>;
  isMetaMask?: boolean;
}

declare global {
  interface Window {
    ethereum?: EthereumProvider;
  }
}

const ConnectWallet = () => {
  const [loading, setLoading] = useState(false);
  const [signing, setSigning] = useState(false);
  const [connectedAddress, setConnectedAddress] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const navigate = useNavigate();
  const isSigningRef = useRef(false);

  const handleWalletConnect = async () => {
    if (!window.ethereum) {
      toast.error("Please install MetaMask or another Ethereum wallet");
      return;
    }

    setLoading(true);
    try {
      const accounts = (await window.ethereum.request({
        method: "eth_requestAccounts",
      })) as string[];

      if (accounts.length === 0) {
        throw new Error("No accounts found");
      }

      const address = accounts[0];
      setConnectedAddress(address);

      const nonceRes = await fetch(`${api_url}/auth/nonce`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ address }),
      });

      if (!nonceRes.ok) {
        const errorData = await nonceRes.json().catch(() => ({
          message: "Failed to get nonce",
        }));
        throw new Error(errorData.message || "Failed to get nonce");
      }

      const nonceData = await nonceRes.json();

      if (nonceData.messageTemplate) {
        setMessage(nonceData.messageTemplate);
      } else if (nonceData.nonce) {
        const fallbackMessage = `circa wants you to sign in with your Ethereum account:\n${address}\n\nSign in to Circa\n\nURI: ${
          window.location.origin
        }\nVersion: 1\nNonce: ${
          nonceData.nonce
        }\nIssued At: ${new Date().toISOString()}`;
        setMessage(fallbackMessage);
      } else {
        throw new Error("Invalid nonce response");
      }
    } catch (error) {
      console.error("Wallet connection error:", error);
      if (error instanceof Error && error.message.includes("rejected")) {
        return;
      }
      const errorMessage =
        error instanceof Error
          ? error.message
          : "Failed to connect wallet. Please try again.";
      toast.error(errorMessage);
      // Reset state on error
      setConnectedAddress(null);
      setMessage(null);
    } finally {
      setLoading(false);
    }
  };

  const handleSignMessage = async () => {
    if (isSigningRef.current || signing) {
      return;
    }

    if (!window.ethereum || !connectedAddress || !message) {
      return;
    }

    isSigningRef.current = true;
    setSigning(true);
    try {
      const signature = (await window.ethereum.request({
        method: "personal_sign",
        params: [message, connectedAddress],
      })) as string;

      const res = await fetch(`${api_url}/auth/signup/complete`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          address: connectedAddress,
          signature,
          message,
        }),
      });

      if (!res.ok) {
        const error = await res.json().catch(() => ({
          message: "Failed to complete signup",
        }));

        if (res.status === 401) {
          const errorMsg = error.message || "Session expired";
          if (
            errorMsg.includes("session") ||
            errorMsg.includes("expired") ||
            errorMsg.includes("nonce")
          ) {
            toast.error(
              "Your session has expired. Please start the signup process again.",
            );
            navigate({ to: "/auth/signup" });
            return;
          }
        }
        throw new Error(error.message || "Failed to complete signup");
      }

      localStorage.removeItem("circa_signup_data");
      toast.success("Account created successfully!");
      navigate({ to: "/app" });
    } catch (error) {
      console.error("Message signing error:", error);

      if ((error as { code?: number })?.code === 4001) {
        return;
      }
      const errorMessage =
        error instanceof Error
          ? error.message
          : "Failed to sign message. Please try again.";
      toast.error(errorMessage);
    } finally {
      setSigning(false);
      isSigningRef.current = false;
    }
  };

  return (
    <div className="min-h-screen w-full relative z-20 bg-white dark:bg-[#1a1625] fade-in flex flex-col antialiased overflow-x-hidden">
      <LoginWalletView
        onConnect={handleWalletConnect}
        onSign={handleSignMessage}
        loading={loading}
        signing={signing}
        connectedAddress={connectedAddress}
        hasMessage={!!message}
      />
    </div>
  );
};

type LoginWalletViewProps = {
  onConnect: () => void;
  onSign: () => void;
  loading: boolean;
  signing: boolean;
  connectedAddress: string | null;
  hasMessage: boolean;
};

const LoginWalletView = ({
  onConnect,
  onSign,
  loading,
  signing,
  connectedAddress,
  hasMessage,
}: LoginWalletViewProps) => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-white dark:bg-[#1a1625] px-4 fade-in">
      <div className="w-full max-w-sm bg-white dark:bg-[#241f2e] border border-[#eeeeee] dark:border-[#3d3551] shadow-sm rounded-xl p-8 text-center">
        <div className="mb-6 flex justify-center">
          <div className="w-12 h-12 bg-gradient-primary rounded-lg flex items-center justify-center text-white">
            {connectedAddress ? (
              <CheckCircle2 className="w-6 h-6" strokeWidth={1.5} />
            ) : (
              <CircleDashed className="w-6 h-6" strokeWidth={1.5} />
            )}
          </div>
        </div>

        {connectedAddress ? (
          <>
            <div className="mb-4 flex items-center justify-center gap-2 text-green-600 dark:text-green-400">
              <CheckCircle2 className="w-5 h-5" strokeWidth={2} />
              <span className="text-sm font-medium">Email verified</span>
            </div>
            <h1 className="text-xl font-semibold tracking-tight text-[#333] dark:text-[#f5f3ff] mb-2">
              Connect your wallet to finish signup
            </h1>
            <p className="text-sm text-[#666] dark:text-[#c4b5fd] mb-6">
              Wallet connected: {connectedAddress.slice(0, 6)}...
              {connectedAddress.slice(-4)}
            </p>

            {hasMessage ? (
              <button
                onClick={(e) => {
                  if (signing) {
                    e.preventDefault();
                    e.stopPropagation();
                    return;
                  }
                  onSign();
                }}
                disabled={signing}
                className="w-full cursor-pointer bg-gradient-primary hover:opacity-90 text-white font-medium text-sm py-2.5 rounded-lg transition-all flex items-center justify-center gap-2 group shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
                style={signing ? { pointerEvents: "none" } : {}}
              >
                <span>{signing ? "Signing..." : "Sign message"}</span>
              </button>
            ) : (
              <p className="text-sm text-[#666] dark:text-[#c4b5fd]">
                Preparing message...
              </p>
            )}
          </>
        ) : (
          <>
            <div className="mb-4 flex items-center justify-center gap-2 text-green-600 dark:text-green-400">
              <CheckCircle2 className="w-5 h-5" strokeWidth={2} />
              <span className="text-sm font-medium">Email verified</span>
            </div>
            <h1 className="text-xl font-semibold tracking-tight text-[#333] dark:text-[#f5f3ff] mb-2">
              Connect your wallet to finish signup
            </h1>
            <p className="text-sm text-[#666] dark:text-[#c4b5fd] mb-8">
              Connect your wallet to verify your identity and complete your
              account setup.
            </p>

            <button
              onClick={onConnect}
              disabled={loading}
              className="w-full bg-gradient-primary cursor-pointer hover:opacity-90 text-white font-medium text-sm py-2.5 rounded-lg transition-all flex items-center justify-center gap-2 group shadow-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span>{loading ? "Connecting..." : "Connect wallet"}</span>
              <Wallet className="w-4 h-4 text-white/80 group-hover:text-white transition-colors" />
            </button>
          </>
        )}

        <p className="mt-4 text-xs text-[#999] dark:text-[#a78bfa]">
          By connecting, you agree to our Terms of Service.
        </p>
      </div>
    </div>
  );
};

export default ConnectWallet;
