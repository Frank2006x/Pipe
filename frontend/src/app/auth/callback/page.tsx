"use client";

import React, { useEffect, useState, Suspense } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import axios from "axios";
import { Loader2, AlertCircle } from "lucide-react";

function CallbackHandler() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const code = searchParams.get("code");
    if (!code) {
      setError("No authorization code provided from GitHub.");
      return;
    }

    const swapCodeForToken = async () => {
      try {
        const backendBaseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
        // Invoke backend callback via Axios in the background
        await axios.get(`${backendBaseUrl}/auth/github/callback?code=${code}`, {
          withCredentials: true,
        });
        
        // Authenticated successfully. Redirect back to homepage
        router.push("/");
      } catch (err: any) {
        console.error("Callback verification failed:", err);
        setError("Failed to verify GitHub authorization. Please try again.");
      }
    };

    swapCodeForToken();
  }, [searchParams, router]);

  if (error) {
    return (
      <div className="w-full max-w-md bg-card border border-border rounded-xl p-8 text-center space-y-4 shadow-lg">
        <div className="flex justify-center text-destructive">
          <AlertCircle className="h-12 w-12" />
        </div>
        <h2 className="text-2xl font-bold text-foreground">Authentication Error</h2>
        <p className="text-sm text-muted-foreground">{error}</p>
        <a 
          href="/login" 
          className="inline-block w-full rounded-lg bg-primary py-2.5 text-sm font-semibold text-primary-foreground transition-all hover:opacity-90"
        >
          Back to Login
        </a>
      </div>
    );
  }

  return (
    <div className="flex flex-col items-center justify-center space-y-4 text-center">
      <Loader2 className="h-12 w-12 text-primary animate-spin" />
      <h2 className="text-2xl font-bold text-foreground">Verifying GitHub Account...</h2>
      <p className="text-sm text-muted-foreground">Completing your secure handshake. Please hold on.</p>
    </div>
  );
}

export default function CallbackPage() {
  return (
    <div className="min-h-screen bg-background text-foreground flex items-center justify-center p-4 relative overflow-hidden">
      {/* Background soft glow mesh */}
      <div className="absolute top-1/4 left-1/2 -translate-x-1/2 w-[500px] h-[300px] bg-[radial-gradient(circle,_var(--primary)_0%,_transparent_75%)] opacity-[0.06] blur-3xl pointer-events-none -z-10" />

      <Suspense fallback={
        <div className="flex flex-col items-center justify-center space-y-4 text-center">
          <Loader2 className="h-12 w-12 text-primary animate-spin" />
          <h2 className="text-2xl font-bold text-foreground">Loading...</h2>
        </div>
      }>
        <CallbackHandler />
      </Suspense>
    </div>
  );
}
