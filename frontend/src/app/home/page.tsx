"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import axios from "axios";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  CardFooter
} from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { LogOut, Mail, Shield, User, GitBranch, Sun, Moon, Plus, Loader2, Globe, Lock, Check, AlertCircle, ExternalLink } from "lucide-react";
import {
  CommandDialog,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
} from "@/components/ui/command";

const GithubIcon = (props: React.SVGProps<SVGSVGElement>) => (
  <svg viewBox="0 0 24 24" fill="currentColor" {...props}>
    <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12" />
  </svg>
);

interface UserProfile {
  id: number;
  github_id: number;
  username: string;
  email: string;
  avatar_url: string | { String: string; Valid: boolean };
}

interface ImportedRepository {
  id: number;
  user_id: number;
  github_repo_id: number;
  name: string;
  full_name: string;
  html_url: string;
  default_branch: string;
  private: boolean;
  owner: string;
  description: string | { String: string; Valid: boolean } | null;
  clone_url: string;
}

interface GithubRepoOwner {
  id: number;
  login: string;
}

interface GithubRepository {
  id: number;
  name: string;
  full_name: string;
  owner: GithubRepoOwner;
  description: string;
  default_branch: string;
  private: boolean;
  html_url: string;
  clone_url: string;
  language: string;
  visibility: string;
}

export default function HomePage() {
  const router = useRouter();
  const [user, setUser] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [theme, setTheme] = useState<"light" | "dark">("light");

  const [importedRepos, setImportedRepos] = useState<ImportedRepository[]>([]);
  const [loadingImported, setLoadingImported] = useState(true);

  // Dialog & GitHub repository selection state
  const [isImportOpen, setIsImportOpen] = useState(false);
  const [githubRepos, setGithubRepos] = useState<GithubRepository[]>([]);
  const [loadingGithubRepos, setLoadingGithubRepos] = useState(false);
  const [githubReposError, setGithubReposError] = useState<string | null>(null);
  const [importingRepoId, setImportingRepoId] = useState<number | null>(null);

  const backendBaseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  // Sync theme with DOM and localStorage
  useEffect(() => {
    const savedTheme = localStorage.getItem("theme") as "light" | "dark" | null;
    const systemPrefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    const initialTheme = savedTheme || (systemPrefersDark ? "dark" : "light");

    setTheme(initialTheme);
    if (initialTheme === "dark") {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
  }, []);

  const toggleTheme = () => {
    const nextTheme = theme === "light" ? "dark" : "light";
    setTheme(nextTheme);
    localStorage.setItem("theme", nextTheme);
    if (nextTheme === "dark") {
      document.documentElement.classList.add("dark");
    } else {
      document.documentElement.classList.remove("dark");
    }
  };

  const fetchImportedRepos = async () => {
    setLoadingImported(true);
    try {
      const response = await axios.get(`${backendBaseUrl}/repositories`, {
        withCredentials: true,
      });
      setImportedRepos(response.data.repositories || []);
    } catch (err) {
      console.error("Failed to fetch imported repositories:", err);
    } finally {
      setLoadingImported(false);
    }
  };

  const fetchGithubRepos = async () => {
    setLoadingGithubRepos(true);
    setGithubReposError(null);
    try {
      const response = await axios.get(`${backendBaseUrl}/github/repositories`, {
        withCredentials: true,
      });
      setGithubRepos(response.data.repositories || []);
    } catch (err: any) {
      console.error("Failed to fetch GitHub repositories:", err);
      setGithubReposError("Failed to fetch repositories from GitHub. Check your integration credentials.");
    } finally {
      setLoadingGithubRepos(false);
    }
  };

  const handleImportRepo = async (repo: GithubRepository) => {
    setImportingRepoId(repo.id);
    try {
      await axios.post(
        `${backendBaseUrl}/repositories/import`,
        {
          owner: repo.owner.login,
          repo: repo.name,
        },
        {
          withCredentials: true,
        }
      );
      await fetchImportedRepos();
      setIsImportOpen(false);
    } catch (err: any) {
      console.error("Import failed:", err);
      alert(err.response?.data?.message || "Failed to import repository. It may already be imported.");
    } finally {
      setImportingRepoId(null);
    }
  };

  const handleDeleteRepo = async (id: number) => {
    if (!confirm("Are you sure you want to remove this repository from Pipe?")) return;
    try {
      await axios.delete(`${backendBaseUrl}/repositories/${id}`, {
        withCredentials: true,
      });
      await fetchImportedRepos();
    } catch (err) {
      console.error("Failed to delete repository:", err);
      alert("Failed to remove repository.");
    }
  };

  useEffect(() => {
    const fetchUserProfile = async () => {
      try {
        const response = await axios.get(`${backendBaseUrl}/auth/me`, {
          withCredentials: true,
        });
        setUser(response.data);
        fetchImportedRepos();
      } catch (err: any) {
        console.error("Failed to fetch user profile, redirecting to login:", err);
        router.push("/login");
      } finally {
        setLoading(false);
      }
    };

    fetchUserProfile();
  }, [router, backendBaseUrl]);

  useEffect(() => {
    if (isImportOpen) {
      fetchGithubRepos();
    }
  }, [isImportOpen]);

  const handleLogout = async () => {
    try {
      await axios.get(`${backendBaseUrl}/auth/logout`, {
        withCredentials: true,
      });
      router.push("/login");
    } catch (err) {
      console.error("Logout failed:", err);
      // Fallback: Force redirect anyway
      router.push("/login");
    }
  };

  // Helper to extract image URL safely from potential pgtype.Text structure
  const getAvatarUrl = (avatarUrlField: any): string => {
    if (!avatarUrlField) return "";
    if (typeof avatarUrlField === "object" && "String" in avatarUrlField) {
      return avatarUrlField.String;
    }
    return String(avatarUrlField);
  };

  return (
    <div className="min-h-screen bg-background text-foreground font-sans relative flex flex-col">
      {/* Background soft glow mesh */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full max-w-7xl h-[400px] bg-[radial-gradient(circle_at_top,_var(--primary)_0%,_transparent_65%)] opacity-[0.06] blur-3xl pointer-events-none -z-10" />

      {/* Simple Professional Header */}
      <header className="border-b border-border bg-card/30 backdrop-blur-sm sticky top-0 z-40">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <span className="text-xl font-bold tracking-tight text-foreground">
              PIPE<span className="text-primary">.</span> Dashboard
            </span>
          </div>
          <div className="flex items-center gap-4">
            {/* Theme Toggle Button */}
            <button
              onClick={toggleTheme}
              aria-label="Toggle Theme"
              className="rounded-lg p-2 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus:outline-none"
            >
              {theme === "dark" ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
            </button>
            {user && (
              <Button variant="ghost" size="sm" onClick={handleLogout} className="text-muted-foreground hover:text-foreground">
                <LogOut className="h-4 w-4 mr-2" />
                Sign Out
              </Button>
            )}
          </div>
        </div>
      </header>

      {/* Main Layout Grid (80% / 20%) */}
      <main className="flex-1 max-w-7xl w-full mx-auto px-4 py-8 md:py-12">
        {loading ? (
          <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
            <div className="lg:col-span-4 space-y-4">
              <Skeleton className="h-40 w-full" />
            </div>
            <div className="lg:col-span-1">
              <Card className="border border-border bg-card">
                <CardHeader className="flex flex-col items-center">
                  <Skeleton className="h-16 w-16 rounded-full" />
                  <Skeleton className="h-4 w-2/3 mt-3" />
                </CardHeader>
              </Card>
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-5 gap-6 items-start">
            {/* Left Workspace Panel (80% / 4 Columns) */}
            <div className="lg:col-span-4 space-y-6">
              <div className="flex items-center justify-between">
                <h1 className="text-xl font-bold tracking-tight text-foreground">Repositories</h1>
                <Button className="h-10 w-40 rounded-md font-medium text-sm flex items-center justify-center gap-1.5" onClick={() => setIsImportOpen(true)}>
                  <Plus className="h-4 w-4" />
                  Import Repository
                </Button>
              </div>

              {loadingImported ? (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <Skeleton className="h-40 w-full rounded-xl" />
                  <Skeleton className="h-40 w-full rounded-xl" />
                </div>
              ) : importedRepos.length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {importedRepos.map((repo) => (
                    <Card key={repo.id} className="border border-border bg-card/30 hover:bg-card/60 transition-all duration-300 shadow-xs hover:border-primary/20 flex flex-col justify-between group">
                      <CardHeader className="pb-3">
                        <div className="flex items-start justify-between gap-4">
                          <div className="space-y-1 min-w-0">
                            <CardTitle className="text-base font-bold flex items-center gap-2 truncate text-foreground group-hover:text-primary transition-colors">
                              <GitBranch className="h-4 w-4 text-primary shrink-0" />
                              <span className="truncate">{repo.name}</span>
                            </CardTitle>
                            <span className="text-xs text-muted-foreground font-mono truncate block">{repo.full_name}</span>
                          </div>
                          <span className={`text-[10px] px-2 py-0.5 rounded-full font-semibold uppercase border shrink-0 ${
                            repo.private
                              ? "bg-amber-500/10 text-amber-500 border-amber-500/20"
                              : "bg-emerald-500/10 text-emerald-500 border-emerald-500/20"
                          }`}>
                            {repo.private ? "Private" : "Public"}
                          </span>
                        </div>
                        {repo.description && (
                          <p className="text-xs text-muted-foreground line-clamp-2 mt-2 leading-relaxed">
                            {typeof repo.description === "object" && repo.description !== null && "String" in repo.description
                              ? repo.description.String
                              : String(repo.description)}
                          </p>
                        )}
                      </CardHeader>
                      <CardContent className="pb-3 text-xs text-muted-foreground flex flex-wrap gap-x-4 gap-y-2 border-t border-border/40 pt-3">
                        <div className="flex items-center gap-1.5">
                          <span className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                          <span>Default branch: <span className="font-mono text-foreground font-medium">{repo.default_branch}</span></span>
                        </div>
                      </CardContent>
                      <CardFooter className="pt-0 pb-4 flex justify-between gap-2">
                        <a 
                          href={repo.html_url} 
                          target="_blank" 
                          rel="noopener noreferrer"
                          className="text-xs text-muted-foreground hover:text-primary transition-colors flex items-center gap-1"
                        >
                          View GitHub <ExternalLink className="h-3.5 w-3.5" />
                        </a>
                        <Button 
                          variant="ghost" 
                          size="sm" 
                          className="h-8 text-xs hover:bg-destructive/10 hover:text-destructive text-muted-foreground transition-all"
                          onClick={() => handleDeleteRepo(repo.id)}
                        >
                          Remove
                        </Button>
                      </CardFooter>
                    </Card>
                  ))}
                </div>
              ) : (
                <div className="p-8 rounded-lg border border-dashed border-border bg-card/10 flex flex-col items-center justify-center text-center min-h-[350px] transition-all">
                  <GitBranch className="h-12 w-12 text-muted-foreground mb-4 opacity-80" />
                  <h2 className="text-2xl font-bold text-foreground">No active repositories</h2>
                  <p className="text-muted-foreground text-sm max-w-sm mt-2">
                    Connect your GitHub repositories to set up pipelines, enable AI pull request audits, and monitor your automated container builds.
                  </p>
                  <Button className="mt-6" onClick={() => setIsImportOpen(true)}>
                    Import Repository
                  </Button>
                </div>
              )}
            </div>

            {/* Right Profile Panel (20% / 1 Column) */}
            <div className="lg:col-span-1">
              <Card className="border border-border bg-card shadow-sm hover:border-primary/20 transition-all duration-300">
                <CardHeader className="flex flex-col items-center pb-4 text-center border-b border-border/50">
                  <Avatar className="h-16 w-16 border border-border shadow-sm mb-3">
                    <AvatarImage src={getAvatarUrl(user?.avatar_url)} alt={user?.username} />
                    <AvatarFallback className="bg-primary/10 text-primary font-bold text-lg">
                      {user?.username?.substring(0, 2).toUpperCase() || <User />}
                    </AvatarFallback>
                  </Avatar>
                  <CardTitle className="text-base font-bold truncate max-w-full">{user?.username}</CardTitle>
                  <CardDescription className="text-xs">GitHub Member</CardDescription>
                </CardHeader>
                <CardContent className="py-4 space-y-4 text-xs">
                  {user?.email && (
                    <div className="flex items-center gap-2 text-muted-foreground truncate">
                      <Mail className="h-3.5 w-3.5 shrink-0" />
                      <span className="truncate">{user.email}</span>
                    </div>
                  )}
                  <div className="flex items-center gap-2 text-muted-foreground">
                    <GithubIcon className="h-3.5 w-3.5 shrink-0" />
                    <span>ID: {user?.github_id}</span>
                  </div>
                  <div className="flex items-start gap-2 text-muted-foreground">
                    <Shield className="h-3.5 w-3.5 shrink-0 mt-0.5" />
                    <span>Session: HttpOnly Secured</span>
                  </div>
                </CardContent>
                <CardFooter className="pt-0 pb-4 px-6 flex justify-stretch">
                  <Button variant="outline" size="sm" onClick={handleLogout} className="w-full text-xs">
                    <LogOut className="h-3.5 w-3.5 mr-1.5" />
                    Sign Out
                  </Button>
                </CardFooter>
              </Card>
            </div>
          </div>
        )}
      </main>

      {/* Repository Selector Dialog */}
      <CommandDialog open={isImportOpen} onOpenChange={setIsImportOpen} title="Import Repository" description="Select a repository from your GitHub account to connect.">
        <CommandInput placeholder="Search your GitHub repositories..." />
        <CommandList className="max-h-[400px]">
          {loadingGithubRepos && (
            <div className="py-6 flex flex-col items-center justify-center gap-2 text-muted-foreground">
              <Loader2 className="h-6 w-6 animate-spin text-primary" />
              <span className="text-xs">Fetching repositories from GitHub...</span>
            </div>
          )}
          
          {githubReposError && !loadingGithubRepos && (
            <div className="py-6 px-4 flex flex-col items-center justify-center text-center gap-2 text-destructive">
              <AlertCircle className="h-8 w-8" />
              <p className="text-sm font-semibold">{githubReposError}</p>
              <Button size="sm" variant="outline" className="mt-2" onClick={fetchGithubRepos}>
                Retry Fetching
              </Button>
            </div>
          )}

          {!loadingGithubRepos && !githubReposError && (
            <>
              <CommandEmpty>No repositories found on GitHub.</CommandEmpty>
              
              <CommandGroup heading="Your GitHub Repositories">
                {githubRepos.map((repo) => {
                  const isAlreadyImported = importedRepos.some(
                    (imported) => imported.github_repo_id === repo.id
                  );
                  const isCurrentlyImporting = importingRepoId === repo.id;

                  return (
                    <CommandItem
                      key={repo.id}
                      value={`${repo.full_name} ${repo.name}`}
                      disabled={isAlreadyImported || isCurrentlyImporting}
                      onSelect={() => handleImportRepo(repo)}
                      className="flex items-center justify-between gap-4 p-3 cursor-pointer hover:bg-muted/80"
                    >
                      <div className="flex items-start gap-3 min-w-0">
                        <div className="p-2 rounded-md bg-muted border border-border shrink-0 mt-0.5">
                          <GitBranch className="h-4 w-4 text-foreground" />
                        </div>
                        <div className="min-w-0 space-y-0.5">
                          <span className="text-sm font-semibold text-foreground truncate block">
                            {repo.name}
                          </span>
                          <span className="text-xs text-muted-foreground truncate block">
                            {repo.full_name}
                          </span>
                          {repo.description && (
                            <p className="text-xs text-muted-foreground line-clamp-1">
                              {repo.description}
                            </p>
                          )}
                        </div>
                      </div>

                      <div className="flex items-center gap-2 shrink-0">
                        {repo.language && (
                          <span className="text-[10px] px-1.5 py-0.5 rounded-md bg-primary/10 text-primary border border-primary/20 font-medium">
                            {repo.language}
                          </span>
                        )}
                        <span className={`text-[9px] font-semibold uppercase px-1.5 py-0.5 rounded-md border ${
                          repo.private
                            ? "bg-amber-500/10 text-amber-500 border-amber-500/20"
                            : "bg-emerald-500/10 text-emerald-500 border-emerald-500/20"
                        }`}>
                          {repo.private ? "Private" : "Public"}
                        </span>
                        
                        {isAlreadyImported ? (
                          <span className="text-xs font-medium text-emerald-500 flex items-center gap-1">
                            <Check className="h-3.5 w-3.5" /> Imported
                          </span>
                        ) : isCurrentlyImporting ? (
                          <Loader2 className="h-4 w-4 animate-spin text-primary" />
                        ) : (
                          <Plus className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
                        )}
                      </div>
                    </CommandItem>
                  );
                })}
              </CommandGroup>
            </>
          )}
        </CommandList>
      </CommandDialog>
    </div>
  );
}
