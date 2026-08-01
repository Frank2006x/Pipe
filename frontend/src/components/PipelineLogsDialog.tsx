"use client";

import React, { useState, useEffect } from "react";
import axios from "axios";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import {
  GitBranch,
  GitCommit,
  Clock,
  CheckCircle2,
  XCircle,
  Loader2,
  RefreshCw,
  Copy,
  Terminal,
  Layers,
  Container,
  User,
  ExternalLink,
  Check,
} from "lucide-react";

export interface PipelineItem {
  id: number;
  repository_id: number;
  github_delivery_id?: string;
  commit_sha: string;
  commit_message?: string;
  branch: string;
  event_type: string;
  status: "pending" | "running" | "success" | "failed" | "cancelled";
  trigger_username?: string;
  started_at?: string;
  finished_at?: string;
  created_at: string;
}

export interface JobItem {
  id: number;
  pipeline_id: number;
  template_id?: number;
  status: "pending" | "running" | "success" | "failed" | "cancelled";
  name: string;
  order_index: number;
  image: string;
  working_directory: string;
  commands: string[];
  logs?: string;
  exit_code?: number | { Int32: number; Valid: boolean };
  started_at?: string;
  finished_at?: string;
}

interface PipelineLogsDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  repoId: number;
  repoName: string;
  repoFullName: string;
}

export default function PipelineLogsDialog({
  open,
  onOpenChange,
  repoId,
  repoName,
  repoFullName,
}: PipelineLogsDialogProps) {
  const [pipelines, setPipelines] = useState<PipelineItem[]>([]);
  const [selectedPipelineId, setSelectedPipelineId] = useState<number | null>(null);
  const [jobs, setJobs] = useState<JobItem[]>([]);
  const [selectedJobId, setSelectedJobId] = useState<number | null>(null);

  const [loadingPipelines, setLoadingPipelines] = useState(false);
  const [loadingJobs, setLoadingJobs] = useState(false);
  const [copiedLogs, setCopiedLogs] = useState(false);

  const backendBaseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  const fetchPipelines = async () => {
    if (!repoId) return;
    setLoadingPipelines(true);
    try {
      const response = await axios.get(
        `${backendBaseUrl}/repositories/${repoId}/pipelines`,
        { withCredentials: true }
      );
      const fetched: PipelineItem[] = response.data?.pipelines || [];
      setPipelines(fetched);
      if (fetched.length > 0 && (!selectedPipelineId || !fetched.some(p => p.id === selectedPipelineId))) {
        setSelectedPipelineId(fetched[0].id);
      }
    } catch (err) {
      console.error("Failed to load pipelines:", err);
    } finally {
      setLoadingPipelines(false);
    }
  };

  const fetchJobs = async (pipelineId: number) => {
    setLoadingJobs(true);
    try {
      const response = await axios.get(
        `${backendBaseUrl}/pipelines/${pipelineId}/jobs`,
        { withCredentials: true }
      );
      const fetchedJobs: JobItem[] = response.data?.jobs || [];
      setJobs(fetchedJobs);
      if (fetchedJobs.length > 0) {
        setSelectedJobId(fetchedJobs[0].id);
      } else {
        setSelectedJobId(null);
      }
    } catch (err) {
      console.error("Failed to load jobs for pipeline:", err);
    } finally {
      setLoadingJobs(false);
    }
  };

  useEffect(() => {
    if (open && repoId) {
      fetchPipelines();
    }
  }, [open, repoId]);

  useEffect(() => {
    if (selectedPipelineId) {
      fetchJobs(selectedPipelineId);
    }
  }, [selectedPipelineId]);

  const selectedPipeline = pipelines.find((p) => p.id === selectedPipelineId);
  const selectedJob = jobs.find((j) => j.id === selectedJobId) || jobs[0];

  const extractExitCode = (code?: number | { Int32: number; Valid: boolean }) => {
    if (code === undefined || code === null) return null;
    if (typeof code === "number") return code;
    if (typeof code === "object" && code.Valid) return code.Int32;
    return null;
  };

  const handleCopyLogs = () => {
    if (!selectedJob?.logs) return;
    navigator.clipboard.writeText(selectedJob.logs);
    setCopiedLogs(true);
    setTimeout(() => setCopiedLogs(false), 2000);
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "success":
        return (
          <span className="inline-flex items-center gap-1 text-[11px] font-semibold px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-500 border border-emerald-500/20">
            <CheckCircle2 className="h-3 w-3" /> Success
          </span>
        );
      case "failed":
        return (
          <span className="inline-flex items-center gap-1 text-[11px] font-semibold px-2 py-0.5 rounded-full bg-destructive/10 text-destructive border border-destructive/20">
            <XCircle className="h-3 w-3" /> Failed
          </span>
        );
      case "running":
        return (
          <span className="inline-flex items-center gap-1 text-[11px] font-semibold px-2 py-0.5 rounded-full bg-primary/10 text-primary border border-primary/20 animate-pulse">
            <Loader2 className="h-3 w-3 animate-spin" /> Running
          </span>
        );
      default:
        return (
          <span className="inline-flex items-center gap-1 text-[11px] font-semibold px-2 py-0.5 rounded-full bg-amber-500/10 text-amber-500 border border-amber-500/20">
            <Clock className="h-3 w-3" /> Pending
          </span>
        );
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-5xl max-h-[90vh] flex flex-col p-0 gap-0 overflow-hidden border border-border bg-card shadow-2xl">
        {/* Header */}
        <DialogHeader className="p-4 border-b border-border bg-muted/20 flex flex-row items-center justify-between">
          <div className="space-y-0.5">
            <div className="flex items-center gap-2">
              <Layers className="h-5 w-5 text-primary" />
              <DialogTitle className="text-base font-bold text-foreground">
                Pipeline Execution History
              </DialogTitle>
            </div>
            <DialogDescription className="text-xs text-muted-foreground">
              Repository: <span className="font-mono font-semibold text-foreground">{repoFullName || repoName}</span>
            </DialogDescription>
          </div>

          <Button
            variant="outline"
            size="sm"
            className="h-8 text-xs gap-1.5"
            onClick={fetchPipelines}
            disabled={loadingPipelines}
          >
            <RefreshCw className={`h-3.5 w-3.5 ${loadingPipelines ? "animate-spin text-primary" : ""}`} />
            Refresh
          </Button>
        </DialogHeader>

        {/* Body Grid */}
        <div className="flex-1 overflow-hidden grid grid-cols-1 md:grid-cols-12 divide-y md:divide-y-0 md:divide-x divide-border">
          {/* Left Column: Pipeline Runs List (4 cols) */}
          <div className="md:col-span-4 p-3 space-y-2 bg-muted/10 overflow-y-auto max-h-[60vh] md:max-h-none">
            <div className="text-xs font-semibold text-muted-foreground uppercase tracking-wider px-1 pb-1">
              Runs ({pipelines.length})
            </div>

            {loadingPipelines && pipelines.length === 0 ? (
              <div className="py-12 flex flex-col items-center justify-center gap-2 text-muted-foreground">
                <Loader2 className="h-6 w-6 animate-spin text-primary" />
                <span className="text-xs">Loading pipeline runs...</span>
              </div>
            ) : pipelines.length > 0 ? (
              <div className="space-y-2">
                {pipelines.map((pipe) => {
                  const isActive = pipe.id === selectedPipelineId;
                  return (
                    <div
                      key={pipe.id}
                      onClick={() => setSelectedPipelineId(pipe.id)}
                      className={`p-3 rounded-lg border text-left cursor-pointer transition-all space-y-2 ${
                        isActive
                          ? "bg-card border-primary ring-1 ring-primary/30 shadow-xs"
                          : "bg-background/50 border-border hover:bg-card hover:border-muted-foreground/30"
                      }`}
                    >
                      <div className="flex items-center justify-between">
                        <span className="font-mono text-xs font-bold text-foreground">
                          Pipeline #{pipe.id}
                        </span>
                        {getStatusBadge(pipe.status)}
                      </div>

                      <div className="space-y-1 text-xs text-muted-foreground">
                        {pipe.commit_message && (
                          <p className="text-xs text-foreground font-medium line-clamp-1">
                            {pipe.commit_message}
                          </p>
                        )}
                        <div className="flex flex-wrap items-center gap-x-2 gap-y-1 text-[11px]">
                          <span className="flex items-center gap-1 font-mono text-primary">
                            <GitBranch className="h-3 w-3" /> {pipe.branch || "main"}
                          </span>
                          {pipe.commit_sha && (
                            <span className="flex items-center gap-1 font-mono text-muted-foreground">
                              <GitCommit className="h-3 w-3" /> {pipe.commit_sha.substring(0, 7)}
                            </span>
                          )}
                        </div>
                      </div>

                      <div className="text-[10px] text-muted-foreground pt-1 border-t border-border/40 flex justify-between">
                        <span>{new Date(pipe.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })}</span>
                        {pipe.trigger_username && <span>by @{pipe.trigger_username}</span>}
                      </div>
                    </div>
                  );
                })}
              </div>
            ) : (
              <div className="p-6 text-center text-xs text-muted-foreground border border-dashed rounded-lg">
                No pipeline runs found for this repository. Push code or trigger a webhook to start a build.
              </div>
            )}
          </div>

          {/* Right Column: Execution Details & Terminal Logs (8 cols) */}
          <div className="md:col-span-8 p-4 flex flex-col gap-3 overflow-y-auto max-h-[60vh] md:max-h-none bg-background">
            {selectedPipeline ? (
              <>
                {/* Pipeline Info Banner */}
                <div className="p-3 rounded-lg border border-border bg-card/40 flex flex-wrap items-center justify-between gap-3 text-xs">
                  <div className="space-y-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="font-bold text-sm">Pipeline #{selectedPipeline.id}</span>
                      {getStatusBadge(selectedPipeline.status)}
                      <span className="text-[10px] uppercase font-mono px-2 py-0.5 bg-muted rounded border border-border">
                        {selectedPipeline.event_type}
                      </span>
                    </div>
                    <p className="text-muted-foreground truncate">
                      {selectedPipeline.commit_message || "No commit message"}
                    </p>
                  </div>

                  <div className="flex items-center gap-3 text-muted-foreground font-mono text-[11px]">
                    <span className="flex items-center gap-1">
                      <GitBranch className="h-3.5 w-3.5 text-primary" /> {selectedPipeline.branch}
                    </span>
                    {selectedPipeline.commit_sha && (
                      <span className="flex items-center gap-1">
                        <GitCommit className="h-3.5 w-3.5 text-primary" /> {selectedPipeline.commit_sha.substring(0, 7)}
                      </span>
                    )}
                  </div>
                </div>

                {/* Jobs Steps Bar */}
                {jobs.length > 0 ? (
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 border-b border-border pb-2 overflow-x-auto">
                      {jobs.map((job, idx) => {
                        const isJobActive = job.id === selectedJobId;
                        const exitCode = extractExitCode(job.exit_code);
                        return (
                          <button
                            key={job.id}
                            onClick={() => setSelectedJobId(job.id)}
                            className={`px-3 py-1.5 rounded-md text-xs font-semibold flex items-center gap-2 border transition-all shrink-0 ${
                              isJobActive
                                ? "bg-primary text-primary-foreground border-primary shadow-xs"
                                : "bg-card hover:bg-muted border-border text-foreground"
                            }`}
                          >
                            <span>Step {idx + 1}: {job.name}</span>
                            {exitCode !== null && (
                              <span className={`text-[10px] px-1.5 rounded font-mono ${
                                exitCode === 0 ? "bg-emerald-500/20 text-emerald-300" : "bg-destructive/20 text-destructive-foreground"
                              }`}>
                                Exit: {exitCode}
                              </span>
                            )}
                          </button>
                        );
                      })}
                    </div>

                    {/* Active Job Meta */}
                    {selectedJob && (
                      <div className="flex items-center justify-between text-xs px-1 text-muted-foreground">
                        <div className="flex items-center gap-3 font-mono text-[11px]">
                          <span className="flex items-center gap-1 text-primary">
                            <Container className="h-3.5 w-3.5" /> {selectedJob.image}
                          </span>
                          <span>Directory: {selectedJob.working_directory}</span>
                        </div>
                        {selectedJob.logs && (
                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-7 text-xs gap-1.5 text-muted-foreground hover:text-foreground"
                            onClick={handleCopyLogs}
                          >
                            {copiedLogs ? (
                              <>
                                <Check className="h-3.5 w-3.5 text-emerald-500" /> Copied!
                              </>
                            ) : (
                              <>
                                <Copy className="h-3.5 w-3.5" /> Copy Logs
                              </>
                            )}
                          </Button>
                        )}
                      </div>
                    )}

                    {/* Dark Terminal Log Output Box */}
                    <div className="flex-1 rounded-lg border border-border/80 bg-black/95 p-4 font-mono text-xs text-emerald-400 overflow-x-auto max-h-[420px] shadow-inner space-y-1">
                      <div className="text-[11px] text-muted-foreground border-b border-border/40 pb-2 mb-3 flex items-center gap-2">
                        <Terminal className="h-3.5 w-3.5 text-emerald-500" />
                        <span>CONTAINER STDOUT & STDERR LOG STREAM</span>
                      </div>

                      {loadingJobs ? (
                        <div className="py-12 flex flex-col items-center justify-center gap-2 text-muted-foreground">
                          <Loader2 className="h-6 w-6 animate-spin text-primary" />
                          <span>Fetching job stdout logs...</span>
                        </div>
                      ) : selectedJob?.logs ? (
                        <pre className="whitespace-pre-wrap leading-relaxed selection:bg-emerald-500/30 selection:text-emerald-200">
                          {selectedJob.logs}
                        </pre>
                      ) : (
                        <div className="py-8 text-center text-muted-foreground italic">
                          No log output recorded yet for this step.
                        </div>
                      )}
                    </div>
                  </div>
                ) : (
                  <div className="p-8 text-center text-xs text-muted-foreground border border-dashed rounded-lg">
                    No execution jobs found for this pipeline run.
                  </div>
                )}
              </>
            ) : (
              <div className="py-20 flex flex-col items-center justify-center gap-2 text-muted-foreground">
                <Layers className="h-10 w-10 opacity-40" />
                <span className="text-xs">Select a pipeline run from the left panel to view logs.</span>
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
