"use client";

import React, { useState } from "react";
import axios from "axios";
import { GitBranch, ArrowRight, Shield, AlertCircle } from "lucide-react";

const GithubIcon = (props: React.SVGProps<SVGSVGElement>) => (
  <svg viewBox="0 0 24 24" fill="currentColor" {...props}>
    <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/>
  </svg>
);

export default function LoginPage() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleGithubLogin = async () => {
    setLoading(true);
    setError(null);
    try {
      const backendBaseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      const backendUrl = `${backendBaseUrl}/auth/github`;
      
      // Axios request to initialize authentication endpoint
      const response = await axios.get(backendUrl, {
        withCredentials: true,
      });

      // If the backend returns a URL in data (e.g. if modified to return JSON)
      if (response.data && response.data.url) {
        window.location.href = response.data.url;
      } else {
        // Fallback: If it redirects or completes, navigate standard window
        window.location.href = backendUrl;
      }
    } catch (err: any) {
      const backendBaseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
      console.warn("Axios request redirected or threw CORS (expected for OAuth redirects). Falling back to direct redirect.", err);
      // Fallback navigation to complete OAuth login in browser window
      window.location.href = `${backendBaseUrl}/auth/github`;
    }
  };

  return (
    <div className="min-h-screen bg-background text-foreground font-sans flex items-center justify-center p-4 relative overflow-hidden">
      {/* Background soft glow meshes */}
      <div className="absolute top-1/4 left-1/2 -translate-x-1/2 w-[500px] h-[300px] bg-[radial-gradient(circle,_var(--primary)_0%,_transparent_75%)] opacity-[0.06] blur-3xl pointer-events-none -z-10" />

      <div className="w-full max-w-md bg-card border border-border rounded-xl shadow-lg p-8 relative z-10">
        {/* Branding header */}
        <div className="flex flex-col items-center text-center mb-8">
          <div className="h-12 w-12 rounded-lg bg-card border border-border flex items-center justify-center text-primary mb-4 shadow-sm">
            <GitBranch className="h-6 w-6" />
          </div>
          <h2 className="text-3xl font-extrabold tracking-tight text-foreground">
            Welcome to <span className="bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">PIPE</span>
          </h2>
          <p className="text-sm text-muted-foreground mt-2">
            Automate software delivery with AI-powered CI/CD.
          </p>
        </div>

        {error && (
          <div className="mb-6 flex items-start gap-2.5 p-3.5 rounded-lg border border-destructive/20 bg-destructive/5 text-destructive text-sm leading-relaxed">
            <AlertCircle className="h-4.5 w-4.5 shrink-0 mt-0.5" />
            <span>{error}</span>
          </div>
        )}

        <div className="space-y-4">
          <button
            onClick={handleGithubLogin}
            disabled={loading}
            className="w-full flex items-center justify-center gap-3 rounded-lg border border-border bg-background py-3.5 text-sm font-semibold text-foreground transition-all hover:bg-muted hover:border-muted-foreground/30 active:scale-[0.99] disabled:opacity-75 disabled:cursor-not-allowed"
          >
            {loading ? (
              <div className="h-5 w-5 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
            ) : (
              <GithubIcon className="h-5 w-5" />
            )}
            <span>{loading ? "Redirecting to GitHub..." : "Continue with GitHub"}</span>
            <ArrowRight className="h-4 w-4 text-muted-foreground" />
          </button>
        </div>

        {/* Security / Privacy notice */}
        <div className="mt-8 pt-6 border-t border-border flex items-start gap-2.5 text-xs text-muted-foreground leading-relaxed">
          <Shield className="h-4 w-4 text-primary shrink-0 mt-0.5" />
          <p>
            Secure by design. We encrypt all access tokens at rest and enforce webhook verification.
          </p>
        </div>
      </div>
    </div>
  );
}
