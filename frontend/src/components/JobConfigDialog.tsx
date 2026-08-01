"use client";

import React, { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Plus,
  Trash2,
  MoveUp,
  MoveDown,
  Terminal,
  Container,
  Folder,
  Layers,
  Sparkles,
  CheckCircle2,
  Sliders,
  Code2,
} from "lucide-react";

export interface CustomJobConfig {
  id: string;
  name: string;
  image: string;
  workingDirectory: string;
  commands: string[];
  env: { key: string; value: string }[];
}

interface JobConfigDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  repoName: string;
  repoFullName: string;
}

const PRESET_IMAGES = [
  { label: "Go", value: "golang:1.24-alpine" },
  { label: "Node.js", value: "node:20-alpine" },
  { label: "Python", value: "python:3.11-alpine" },
  { label: "Rust", value: "rust:1.75-alpine" },
  { label: "Java / Maven", value: "maven:3.9-eclipse-temurin" },
  { label: "Ubuntu", value: "ubuntu:22.04" },
];

const DEFAULT_MOCK_JOBS: CustomJobConfig[] = [
  {
    id: "job-1",
    name: "Build & Compile",
    image: "golang:1.24-alpine",
    workingDirectory: "/workspace",
    commands: ["pwd", "ls -la", "go mod download", "go build -o app ./cmd/api"],
    env: [{ key: "CGO_ENABLED", value: "0" }],
  },
  {
    id: "job-2",
    name: "Unit & Integration Tests",
    image: "golang:1.24-alpine",
    workingDirectory: "/workspace",
    commands: ["go test -v ./..."],
    env: [{ key: "ENV", value: "test" }],
  },
];

export default function JobConfigDialog({
  open,
  onOpenChange,
  repoName,
  repoFullName,
}: JobConfigDialogProps) {
  const [jobs, setJobs] = useState<CustomJobConfig[]>(DEFAULT_MOCK_JOBS);
  const [activeJobId, setActiveJobId] = useState<string>(DEFAULT_MOCK_JOBS[0].id);
  const [saveSuccess, setSaveSuccess] = useState(false);

  const activeJob = jobs.find((j) => j.id === activeJobId) || jobs[0];

  const handleAddJob = () => {
    const newJob: CustomJobConfig = {
      id: `job-${Date.now()}`,
      name: `Custom Step ${jobs.length + 1}`,
      image: "golang:1.24-alpine",
      workingDirectory: "/workspace",
      commands: ["echo 'Running custom pipeline step...'"],
      env: [],
    };
    setJobs([...jobs, newJob]);
    setActiveJobId(newJob.id);
  };

  const handleDeleteJob = (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (jobs.length <= 1) {
      alert("A pipeline must have at least one job.");
      return;
    }
    const filtered = jobs.filter((j) => j.id !== id);
    setJobs(filtered);
    if (activeJobId === id) {
      setActiveJobId(filtered[0].id);
    }
  };

  const handleMoveJob = (index: number, direction: "up" | "down", e: React.MouseEvent) => {
    e.stopPropagation();
    const newJobs = [...jobs];
    const targetIndex = direction === "up" ? index - 1 : index + 1;
    if (targetIndex < 0 || targetIndex >= newJobs.length) return;

    const temp = newJobs[index];
    newJobs[index] = newJobs[targetIndex];
    newJobs[targetIndex] = temp;
    setJobs(newJobs);
  };

  const updateActiveJob = (updatedFields: Partial<CustomJobConfig>) => {
    setJobs((prev) =>
      prev.map((j) => (j.id === activeJobId ? { ...j, ...updatedFields } : j))
    );
  };

  const handleAddCommand = () => {
    if (!activeJob) return;
    updateActiveJob({ commands: [...activeJob.commands, ""] });
  };

  const handleCommandChange = (index: number, val: string) => {
    if (!activeJob) return;
    const newCmds = [...activeJob.commands];
    newCmds[index] = val;
    updateActiveJob({ commands: newCmds });
  };

  const handleDeleteCommand = (index: number) => {
    if (!activeJob) return;
    updateActiveJob({
      commands: activeJob.commands.filter((_, i) => i !== index),
    });
  };

  const handleAddEnv = () => {
    if (!activeJob) return;
    updateActiveJob({ env: [...activeJob.env, { key: "", value: "" }] });
  };

  const handleEnvChange = (index: number, field: "key" | "value", val: string) => {
    if (!activeJob) return;
    const newEnv = [...activeJob.env];
    newEnv[index][field] = val;
    updateActiveJob({ env: newEnv });
  };

  const handleDeleteEnv = (index: number) => {
    if (!activeJob) return;
    updateActiveJob({
      env: activeJob.env.filter((_, i) => i !== index),
    });
  };

  const applyTemplate = (type: "go" | "node" | "python") => {
    if (type === "go") {
      setJobs([
        {
          id: "job-1",
          name: "Compile Binary",
          image: "golang:1.24-alpine",
          workingDirectory: "/workspace",
          commands: ["pwd", "ls -la", "go mod download", "go build ./..."],
          env: [{ key: "CGO_ENABLED", value: "0" }],
        },
        {
          id: "job-2",
          name: "Run Test Suite",
          image: "golang:1.24-alpine",
          workingDirectory: "/workspace",
          commands: ["go test -v ./..."],
          env: [],
        },
      ]);
    } else if (type === "node") {
      setJobs([
        {
          id: "job-1",
          name: "Install Dependencies & Build",
          image: "node:20-alpine",
          workingDirectory: "/workspace",
          commands: ["npm ci", "npm run build"],
          env: [{ key: "NODE_ENV", value: "production" }],
        },
        {
          id: "job-2",
          name: "Run ESLint & Tests",
          image: "node:20-alpine",
          workingDirectory: "/workspace",
          commands: ["npm test"],
          env: [],
        },
      ]);
    } else if (type === "python") {
      setJobs([
        {
          id: "job-1",
          name: "Install Requirements & Test",
          image: "python:3.11-alpine",
          workingDirectory: "/workspace",
          commands: ["pip install -r requirements.txt", "pytest"],
          env: [{ key: "PYTHONUNBUFFERED", value: "1" }],
        },
      ]);
    }
  };

  const handleSave = () => {
    setSaveSuccess(true);
    setTimeout(() => {
      setSaveSuccess(false);
      onOpenChange(false);
    }, 1200);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-4xl max-h-[90vh] flex flex-col p-0 gap-0 overflow-hidden border border-border bg-card shadow-2xl">
        {/* Header */}
        <DialogHeader className="p-5 border-b border-border/80 bg-muted/20">
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <Sliders className="h-5 w-5 text-primary" />
                <DialogTitle className="text-lg font-bold text-foreground">
                  Configure Pipeline Jobs
                </DialogTitle>
              </div>
              <DialogDescription className="text-xs text-muted-foreground">
                Repository: <span className="font-mono font-semibold text-foreground">{repoFullName || repoName}</span>
              </DialogDescription>
            </div>

            {/* Quick Templates */}
            <div className="flex items-center gap-1.5 bg-background p-1 rounded-lg border border-border">
              <span className="text-[10px] font-semibold uppercase px-2 text-muted-foreground flex items-center gap-1">
                <Sparkles className="h-3 w-3 text-primary" /> Presets:
              </span>
              <Button
                variant="ghost"
                size="sm"
                className="h-7 text-xs px-2"
                onClick={() => applyTemplate("go")}
              >
                Go App
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="h-7 text-xs px-2"
                onClick={() => applyTemplate("node")}
              >
                Node.js
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="h-7 text-xs px-2"
                onClick={() => applyTemplate("python")}
              >
                Python
              </Button>
            </div>
          </div>
        </DialogHeader>

        {/* Body Split View */}
        <div className="flex-1 overflow-hidden grid grid-cols-1 md:grid-cols-12 divide-y md:divide-y-0 md:divide-x divide-border">
          {/* Left Column: Jobs List Manager (4 cols) */}
          <div className="md:col-span-4 p-4 space-y-3 bg-muted/10 overflow-y-auto max-h-[60vh] md:max-h-none">
            <div className="flex items-center justify-between pb-1">
              <span className="text-xs font-semibold text-muted-foreground uppercase tracking-wider flex items-center gap-1.5">
                <Layers className="h-3.5 w-3.5 text-primary" /> Jobs ({jobs.length})
              </span>
              <Button
                variant="outline"
                size="sm"
                className="h-7 text-xs gap-1"
                onClick={handleAddJob}
              >
                <Plus className="h-3.5 w-3.5" /> Add Job
              </Button>
            </div>

            <div className="space-y-2">
              {jobs.map((job, idx) => {
                const isActive = job.id === activeJobId;
                return (
                  <div
                    key={job.id}
                    onClick={() => setActiveJobId(job.id)}
                    className={`p-3 rounded-lg border text-left cursor-pointer transition-all flex items-start justify-between gap-2 group ${
                      isActive
                        ? "bg-card border-primary ring-1 ring-primary/30 shadow-xs"
                        : "bg-background/50 border-border hover:bg-card hover:border-muted-foreground/30"
                    }`}
                  >
                    <div className="space-y-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <span className="h-5 w-5 rounded-full bg-primary/10 text-primary text-[11px] font-bold flex items-center justify-center shrink-0">
                          {idx + 1}
                        </span>
                        <span className="text-sm font-bold text-foreground truncate block">
                          {job.name}
                        </span>
                      </div>
                      <div className="flex flex-wrap items-center gap-x-2 gap-y-1 text-[11px] text-muted-foreground pl-7">
                        <span className="font-mono text-[10px] text-primary truncate max-w-[140px]">
                          {job.image}
                        </span>
                        <span>•</span>
                        <span>{job.commands.length} cmd{job.commands.length !== 1 ? "s" : ""}</span>
                      </div>
                    </div>

                    <div className="flex items-center gap-0.5 opacity-80 group-hover:opacity-100 shrink-0">
                      <button
                        className="p-1 hover:text-foreground text-muted-foreground disabled:opacity-30"
                        disabled={idx === 0}
                        onClick={(e) => handleMoveJob(idx, "up", e)}
                        title="Move Up"
                      >
                        <MoveUp className="h-3.5 w-3.5" />
                      </button>
                      <button
                        className="p-1 hover:text-foreground text-muted-foreground disabled:opacity-30"
                        disabled={idx === jobs.length - 1}
                        onClick={(e) => handleMoveJob(idx, "down", e)}
                        title="Move Down"
                      >
                        <MoveDown className="h-3.5 w-3.5" />
                      </button>
                      <button
                        className="p-1 hover:text-destructive text-muted-foreground ml-1"
                        onClick={(e) => handleDeleteJob(job.id, e)}
                        title="Delete Job"
                      >
                        <Trash2 className="h-3.5 w-3.5" />
                      </button>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          {/* Right Column: Job Configuration Form (8 cols) */}
          <div className="md:col-span-8 p-5 space-y-5 overflow-y-auto max-h-[60vh] md:max-h-none">
            {activeJob ? (
              <>
                {/* Job General Info */}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-1.5">
                    <label className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                      <Code2 className="h-3.5 w-3.5 text-primary" /> Job Name
                    </label>
                    <Input
                      value={activeJob.name}
                      onChange={(e) => updateActiveJob({ name: e.target.value })}
                      placeholder="e.g. Build & Test"
                      className="h-9 text-sm"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <label className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                      <Folder className="h-3.5 w-3.5 text-primary" /> Working Directory
                    </label>
                    <Input
                      value={activeJob.workingDirectory}
                      onChange={(e) =>
                        updateActiveJob({ workingDirectory: e.target.value })
                      }
                      placeholder="/workspace"
                      className="h-9 text-sm font-mono"
                    />
                  </div>
                </div>

                {/* Docker Image Selection */}
                <div className="space-y-1.5">
                  <label className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                    <Container className="h-3.5 w-3.5 text-primary" /> Docker Container Image
                  </label>
                  <div className="flex gap-2">
                    <select
                      value={
                        PRESET_IMAGES.some((p) => p.value === activeJob.image)
                          ? activeJob.image
                          : "custom"
                      }
                      onChange={(e) => {
                        if (e.target.value !== "custom") {
                          updateActiveJob({ image: e.target.value });
                        }
                      }}
                      className="h-9 rounded-md border border-input bg-background px-3 py-1 text-xs shadow-xs focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                    >
                      {PRESET_IMAGES.map((preset) => (
                        <option key={preset.value} value={preset.value}>
                          {preset.label}
                        </option>
                      ))}
                      <option value="custom">Custom...</option>
                    </select>

                    <Input
                      value={activeJob.image}
                      onChange={(e) => updateActiveJob({ image: e.target.value })}
                      placeholder="Enter docker image tag (e.g. node:20-alpine)"
                      className="h-9 text-xs font-mono flex-1"
                    />
                  </div>
                </div>

                {/* Commands Editor */}
                <div className="space-y-2 pt-2 border-t border-border/60">
                  <div className="flex items-center justify-between">
                    <label className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                      <Terminal className="h-3.5 w-3.5 text-primary" /> Shell Commands (Executed in order)
                    </label>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-7 text-xs gap-1 text-primary hover:text-primary"
                      onClick={handleAddCommand}
                    >
                      <Plus className="h-3 w-3" /> Add Command
                    </Button>
                  </div>

                  <div className="space-y-2">
                    {activeJob.commands.map((cmd, cIdx) => (
                      <div key={cIdx} className="flex items-center gap-2">
                        <span className="text-xs font-mono text-muted-foreground w-4 text-right">
                          {cIdx + 1}.
                        </span>
                        <Input
                          value={cmd}
                          onChange={(e) => handleCommandChange(cIdx, e.target.value)}
                          placeholder="e.g. go build ./..."
                          className="h-8 text-xs font-mono flex-1 bg-background"
                        />
                        <button
                          onClick={() => handleDeleteCommand(cIdx)}
                          className="text-muted-foreground hover:text-destructive p-1"
                          title="Remove Command"
                        >
                          <Trash2 className="h-3.5 w-3.5" />
                        </button>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Environment Variables Editor */}
                <div className="space-y-2 pt-2 border-t border-border/60">
                  <div className="flex items-center justify-between">
                    <label className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                      <Sliders className="h-3.5 w-3.5 text-primary" /> Environment Variables (Optional)
                    </label>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-7 text-xs gap-1 text-primary hover:text-primary"
                      onClick={handleAddEnv}
                    >
                      <Plus className="h-3 w-3" /> Add Env Var
                    </Button>
                  </div>

                  {activeJob.env.length > 0 ? (
                    <div className="space-y-2">
                      {activeJob.env.map((eItem, eIdx) => (
                        <div key={eIdx} className="flex items-center gap-2">
                          <Input
                            value={eItem.key}
                            onChange={(e) =>
                              handleEnvChange(eIdx, "key", e.target.value)
                            }
                            placeholder="KEY (e.g. CGO_ENABLED)"
                            className="h-8 text-xs font-mono flex-1"
                          />
                          <span className="text-xs text-muted-foreground font-bold">=</span>
                          <Input
                            value={eItem.value}
                            onChange={(e) =>
                              handleEnvChange(eIdx, "value", e.target.value)
                            }
                            placeholder="VALUE (e.g. 0)"
                            className="h-8 text-xs font-mono flex-1"
                          />
                          <button
                            onClick={() => handleDeleteEnv(eIdx)}
                            className="text-muted-foreground hover:text-destructive p-1"
                            title="Remove Env Var"
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                          </button>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-xs text-muted-foreground italic">
                      No custom environment variables configured.
                    </p>
                  )}
                </div>
              </>
            ) : null}
          </div>
        </div>

        {/* Footer Actions */}
        <DialogFooter className="p-4 border-t border-border bg-muted/20 flex flex-row items-center justify-between">
          <div className="text-xs text-muted-foreground">
            {saveSuccess ? (
              <span className="text-emerald-500 font-semibold flex items-center gap-1.5">
                <CheckCircle2 className="h-4 w-4" /> Pipeline jobs configuration saved!
              </span>
            ) : (
              <span>Configured steps will run sequentially in Docker containers.</span>
            )}
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button size="sm" onClick={handleSave}>
              Save Pipeline Config
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
