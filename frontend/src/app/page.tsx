"use client";

import React, { useEffect, useState } from "react";
import Link from "next/link";
import axios from "axios";
import Navbar from "@/components/Navbar";
import { 
  GitPullRequest, 
  GitBranch,
  Terminal, 
  Settings, 
  ShieldCheck, 
  Activity, 
  Bell, 
  Layers, 
  Zap, 
  ArrowRight,
  Lock,
  Cpu
} from "lucide-react";

export default function Home() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const backendBaseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
        await axios.get(`${backendBaseUrl}/auth/me`, {
          withCredentials: true,
        });
        setIsLoggedIn(true);
      } catch (err) {
        setIsLoggedIn(false);
      }
    };
    checkAuth();
  }, []);

  return (
    <div className="min-h-screen bg-background text-foreground font-sans antialiased relative">
      {/* Subtle Glowing Backdrop */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full max-w-7xl h-[600px] bg-[radial-gradient(circle_at_top,_var(--primary)_0%,_transparent_65%)] opacity-[0.08] blur-3xl pointer-events-none -z-10" />
      
      <Navbar />

      {/* Hero Section */}
      <section className="py-20 md:py-32 border-b border-border bg-card/20 relative overflow-hidden">
        {/* Subtle Accent Glow */}
        <div className="absolute top-1/2 right-10 -translate-y-1/2 w-[300px] h-[300px] bg-[radial-gradient(circle,_var(--accent)_0%,_transparent_70%)] opacity-[0.05] blur-2xl pointer-events-none" />

        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 text-center relative z-10">
          {/* Creative Accent Badge */}
          <div className="inline-flex items-center gap-1.5 rounded-full border border-accent/20 bg-accent/5 px-3.5 py-1.5 text-xs font-semibold text-accent-foreground backdrop-blur-sm mb-8">
            <span className="flex h-2 w-2 rounded-full bg-accent animate-pulse" />
            <span>PIPE AI-Powered CI/CD is now in Public Beta</span>
          </div>

          <h1 className="text-4xl font-extrabold tracking-tight sm:text-6xl max-w-3xl mx-auto leading-tight text-foreground">
            The AI-Powered Pipeline that{" "}
            <span className="bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
              Automates & Audits
            </span>{" "}
            Releases.
          </h1>

          <p className="mt-6 text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto leading-relaxed">
            PIPE automates software delivery from code to deployment. Connect your GitHub repository, build custom pipelines, monitor stage runs, and let AI review pull requests before they merge.
          </p>

          <div className="mt-10 flex flex-col sm:flex-row justify-center gap-4">
            <Link 
              href={isLoggedIn ? "/home" : "/login"}
              className="rounded-lg bg-primary px-8 py-3.5 text-base font-semibold text-primary-foreground shadow-lg shadow-primary/20 transition-all hover:opacity-90 hover:shadow-primary/30 hover:scale-[1.01] flex items-center justify-center"
            >
              {isLoggedIn ? "Go to Dashboard" : "Connect Repository"}
            </Link>
            <a 
              href="#features"
              className="flex items-center justify-center gap-2 rounded-lg border border-border bg-card px-8 py-3.5 text-base font-semibold text-foreground transition-all hover:bg-muted hover:border-muted-foreground/30"
            >
              <span>Explore Features</span>
              <ArrowRight className="h-4 w-4" />
            </a>
          </div>
        </div>
      </section>

      {/* Features Grid Section */}
      <section id="features" className="py-20 border-b border-border bg-background">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-3xl mx-auto mb-16">
            <h2 className="text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
              Key <span className="text-primary font-extrabold">Capabilities</span>
            </h2>
            <p className="mt-4 text-muted-foreground">Minimal configuration. Standard, secure-by-default execution containers.</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {/* Feature 1 */}
            <div className="bg-card border border-border p-6 rounded-lg transition-all duration-300 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5 hover:translate-y-[-2px]">
              <Zap className="h-6 w-6 text-primary mb-4" />
              <h3 className="text-lg font-bold text-foreground mb-2">Automated CI/CD</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">Trigger pipelines automatically on every push or pull request event.</p>
            </div>

            {/* Feature 2 */}
            <div className="bg-card border border-border p-6 rounded-lg transition-all duration-300 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5 hover:translate-y-[-2px]">
              <GitPullRequest className="h-6 w-6 text-primary mb-4" />
              <h3 className="text-lg font-bold text-foreground mb-2">AI PR Reviews</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">Receive intelligent code reviews, security bug detection, and suggestions before merging.</p>
            </div>

            {/* Feature 3 */}
            <div className="bg-card border border-border p-6 rounded-lg transition-all duration-300 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5 hover:translate-y-[-2px]">
              <Settings className="h-6 w-6 text-primary mb-4" />
              <h3 className="text-lg font-bold text-foreground mb-2">One-Click Integration</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">Connect your GitHub repositories and configure build tasks in minutes.</p>
            </div>

            {/* Feature 4 */}
            <div className="bg-card border border-border p-6 rounded-lg transition-all duration-300 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5 hover:translate-y-[-2px]">
              <Activity className="h-6 w-6 text-primary mb-4" />
              <h3 className="text-lg font-bold text-foreground mb-2">Real-Time Monitoring</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">Track build progress, deployment status, and live container logs.</p>
            </div>

            {/* Feature 5 */}
            <div className="bg-card border border-border p-6 rounded-lg transition-all duration-300 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5 hover:translate-y-[-2px]">
              <Bell className="h-6 w-6 text-primary mb-4" />
              <h3 className="text-lg font-bold text-foreground mb-2">Instant Notifications</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">Stay updated on pipeline failures, successful deployments, and event updates.</p>
            </div>

            {/* Feature 6 */}
            <div className="bg-card border border-border p-6 rounded-lg transition-all duration-300 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5 hover:translate-y-[-2px]">
              <Layers className="h-6 w-6 text-primary mb-4" />
              <h3 className="text-lg font-bold text-foreground mb-2">Container-Based Builds</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">Run isolated, reproducible build environments using Docker containers.</p>
            </div>

            {/* Feature 7 */}
            <div className="bg-card border border-border p-6 rounded-lg transition-all duration-300 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5 hover:translate-y-[-2px]">
              <ShieldCheck className="h-6 w-6 text-primary mb-4" />
              <h3 className="text-lg font-bold text-foreground mb-2">Secure by Design</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">Protected with GitHub OAuth, webhook verification, and encrypted secret tokens.</p>
            </div>

            {/* Feature 8 */}
            <div className="bg-card border border-border p-6 rounded-lg transition-all duration-300 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/5 hover:translate-y-[-2px]">
              <Cpu className="h-6 w-6 text-primary mb-4" />
              <h3 className="text-lg font-bold text-foreground mb-2">Faster Development</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">Reduce manual operational tasks and accelerate deployment cycle speed.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Security Callout Section */}
      <section id="security" className="py-16 bg-card/10 border-b border-border relative">
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_bottom_left,_var(--primary)_0%,_transparent_40%)] opacity-[0.03] pointer-events-none" />
        
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 relative z-10">
          <div className="max-w-3xl mx-auto border border-border bg-card p-8 rounded-lg shadow-sm hover:border-accent/30 transition-colors duration-300">
            <div className="flex items-center gap-3 mb-4">
              <Lock className="h-6 w-6 text-accent shrink-0" />
              <h3 className="text-2xl font-bold text-foreground">Secure Code Execution & Safe Tokens</h3>
            </div>
            <p className="text-muted-foreground text-sm leading-relaxed">
              PIPE builds run inside isolated VMs using secure container wrappers. Your API keys, deployment target credentials, and GitHub OAuth scopes are secured with AES-256 GCM encryption at rest.
            </p>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-card/40 py-12 text-center text-muted-foreground text-xs">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 space-y-4">
          <div className="flex items-center justify-center gap-2">
            <GitBranch className="h-4 w-4 text-primary" />
            <span className="font-semibold text-foreground tracking-tight">PIPE CI/CD Platform</span>
          </div>
          <p>© {new Date().getFullYear()} PIPE Inc. All rights reserved.</p>
          <div className="flex justify-center gap-6">
            <a href="#" className="hover:text-foreground">Privacy Policy</a>
            <a href="#" className="hover:text-foreground">Terms of Service</a>
            <a href="#" className="hover:text-foreground">GitHub Organization</a>
          </div>
        </div>
      </footer>
    </div>
  );
}
