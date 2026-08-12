"use client";

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  ArrowLeft,
  Key as LockKeyhole,
  SpinnerGap as LoaderCircle,
  UploadSimple as UploadCloud,
  WarningCircle as CircleAlert,
} from "@phosphor-icons/react/ssr";

import { JobStatusPanel } from "@/components/jobs/job-status-panel";
import { ResultPreview } from "@/components/jobs/result-preview";
import { ToolOptions } from "@/components/tools/tool-options";
import { FileDropzone } from "@/components/upload/file-dropzone";
import { FilePreview } from "@/components/upload/file-preview";
import { getToolBySlug, type MediaTool } from "@/config/tools";
import {
  defaultOptionsFor,
  type JobOptions,
  type OptionErrors,
  validateOptions,
} from "@/features/media-jobs/options";
import { useSelectedFile } from "@/hooks/use-selected-file";
import { isActiveStatus, type Job } from "@/lib/api/contracts";
import {
  APIError,
  cancelJob,
  createJob,
  deleteJob,
  getJob,
  getJobResult,
} from "@/lib/api/jobs";
import { useWorkspaceStore } from "@/stores/use-workspace-store";

interface MediaJobWorkspaceProps {
  toolSlug: string;
}

export function MediaJobWorkspace({ toolSlug }: MediaJobWorkspaceProps) {
  const tool = getToolBySlug(toolSlug);
  if (!tool) {
    return null;
  }

  return <Workspace toolSlug={toolSlug} />;
}

function Workspace({ toolSlug }: { toolSlug: string }) {
  const tool = getToolBySlug(toolSlug)!;

  const Icon = tool.icon;
  const queryClient = useQueryClient();
  const uploadControllerRef = useRef<AbortController | null>(null);

  const {
    selectedFile,
    error: fileError,
    selectFile,
    clearFile,
  } = useSelectedFile(tool.category);
  const [options, setOptions] = useState<JobOptions>(() =>
    defaultOptionsFor(tool.operation),
  );
  const [optionErrors, setOptionErrors] = useState<OptionErrors>({});
  const [uploadProgress, setUploadProgress] = useState(0);

  const imagePreviewStyle = previewStyle(tool.operation, options);
  const activeJob = useWorkspaceStore((state) => state.activeJob);
  const hydrated = useWorkspaceStore((state) => state.hydrated);
  const setActiveJob = useWorkspaceStore((state) => state.setActiveJob);
  const clearActiveJob = useWorkspaceStore((state) => state.clearActiveJob);

  const credentials = activeJob?.toolSlug === tool.slug ? activeJob : null;
  const otherActiveJob =
    activeJob && activeJob.toolSlug !== tool.slug ? activeJob : null;

  const jobQuery = useQuery({
    queryKey: ["jobs", credentials?.id],
    enabled: hydrated && Boolean(credentials),
    queryFn: ({ signal }) =>
      getJob(credentials!.id, credentials!.ownerToken, signal),
    refetchInterval: (query) => {
      const job = query.state.data;
      return job && isActiveStatus(job.status) ? 1500 : false;
    },
    retry: (failureCount, error) => {
      if (error instanceof APIError && [401, 404].includes(error.status)) {
        return false;
      }
      return failureCount < 2;
    },
  });

  const createMutation = useMutation({
    mutationFn: async ({
      file,
      jobOptions,
    }: {
      file: File;
      jobOptions: JobOptions;
    }) => {
      const controller = new AbortController();
      uploadControllerRef.current = controller;
      setUploadProgress(0);

      return createJob({
        operation: tool.operation,
        options: jobOptions,
        file,
        signal: controller.signal,
        onProgress: setUploadProgress,
      });
    },
    onSuccess: (created) => {
      setActiveJob({
        id: created.id,
        ownerToken: created.owner_token,
        toolSlug: tool.slug,
        operation: tool.operation,
        createdAt: new Date().toISOString(),
      });
      queryClient.setQueryData<Job>(["jobs", created.id], created);
    },
    onSettled: () => {
      uploadControllerRef.current = null;
    },
  });

  const cancelMutation = useMutation({
    mutationFn: () => cancelJob(credentials!.id, credentials!.ownerToken),
    onSuccess: (job) => {
      queryClient.setQueryData(["jobs", credentials!.id], job);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (keepFile: boolean) => {
      if (credentials) {
        await deleteJob(credentials.id, credentials.ownerToken);
      }
      return keepFile;
    },
    onSuccess: (keepFile) => {
      if (credentials) {
        queryClient.removeQueries({ queryKey: ["jobs", credentials.id] });
        queryClient.removeQueries({ queryKey: ["job-result", credentials.id] });
      }
      clearActiveJob();
      if (!keepFile) {
        clearFile();
      }
      setOptions(defaultOptionsFor(tool.operation));
      setOptionErrors({});
      setUploadProgress(0);
    },
  });

  const resultQuery = useQuery({
    queryKey: ["job-result", credentials?.id],
    enabled:
      Boolean(credentials) &&
      jobQuery.data?.status === "completed" &&
      jobQuery.data.result_ready,
    queryFn: ({ signal }) =>
      getJobResult(credentials!.id, credentials!.ownerToken, signal),
    staleTime: Number.POSITIVE_INFINITY,
    retry: 1,
  });

  useEffect(() => {
    return () => uploadControllerRef.current?.abort();
  }, []);

  function submitJob() {
    if (!selectedFile || activeJob) {
      return;
    }

    const errors = validateOptions(tool.operation, options);
    setOptionErrors(errors);
    if (Object.keys(errors).length > 0) {
      return;
    }

    createMutation.mutate({
      file: selectedFile.file,
      jobOptions: options,
    });
  }

  function clearUnrecoverableSession() {
    if (credentials) {
      queryClient.removeQueries({ queryKey: ["jobs", credentials.id] });
      queryClient.removeQueries({ queryKey: ["job-result", credentials.id] });
    }
    clearActiveJob();
    clearFile();
  }

  const apiError =
    createMutation.error ??
    jobQuery.error ??
    cancelMutation.error ??
    deleteMutation.error ??
    resultQuery.error;

  return (
    <main
      id="main-content"
      className="mx-auto w-full max-w-6xl px-4 py-5 sm:px-6 lg:px-8 lg:py-7"
    >
      <Link
        href="/"
        className="inline-flex min-h-9 items-center gap-2 text-xs font-semibold text-ink-muted hover:text-ink"
      >
        <ArrowLeft aria-hidden="true" className="size-4" />
        Back to all tools
      </Link>

      <header className="mt-4 grid overflow-hidden rounded-lg border border-line bg-surface shadow-hard sm:grid-cols-[minmax(0,1fr)_8rem]">
        <div className="p-5 sm:p-6">
          <p className="text-xs font-semibold capitalize text-accent-strong">
            {tool.category} tool
          </p>
          <h1 className="mt-1.5 max-w-4xl text-[clamp(1.75rem,4vw,2.4rem)] leading-tight font-semibold tracking-[-0.045em]">
            {tool.title}
          </h1>
          <p className="mt-2 max-w-2xl text-sm leading-6 text-ink-muted">
            {tool.description}
          </p>
        </div>
        <div className="hidden place-items-center border-l border-line bg-canvas text-accent-strong sm:grid">
          <span className="grid size-12 place-items-center rounded-lg bg-accent"><Icon aria-hidden="true" weight="bold" className="size-5" /></span>
        </div>
      </header>

      {otherActiveJob ? (
        <OtherJobNotice
          activeToolSlug={otherActiveJob.toolSlug}
          requestedToolTitle={tool.title}
        />
      ) : (
        <div className="mt-4 grid gap-4 lg:grid-cols-[minmax(0,1fr)_15rem]">
          <section className="min-w-0">
            {!hydrated && <WorkspaceLoading />}

            {hydrated && !credentials && !createMutation.isPending && (
              <div className="rounded-lg border border-line bg-surface p-3 shadow-hard sm:p-4">
                {selectedFile ? (
                  <FilePreview
                    key={selectedFile.previewUrl}
                    imageStyle={imagePreviewStyle}
                    selectedFile={selectedFile}
                    onRemove={clearFile}
                  />
                ) : (
                  <FileDropzone
                    category={tool.category}
                    error={fileError}
                    onFileSelected={selectFile}
                  />
                )}

                <div className="mt-3 rounded-md bg-canvas p-4">
                  <h2 className="text-sm font-semibold">Output settings</h2>
                  <div className="mt-4">
                    <ToolOptions
                      operation={tool.operation}
                      options={options}
                      errors={optionErrors}
                      onChange={(nextOptions) => {
                        setOptions(nextOptions);
                        setOptionErrors({});
                      }}
                    />
                  </div>
                </div>

                <button
                  type="button"
                  disabled={!selectedFile}
                  className="mt-3 inline-flex min-h-11 w-full items-center justify-center gap-2 rounded-md bg-ink px-5 py-2.5 text-sm font-semibold text-canvas shadow-hard hover:-translate-y-0.5 hover:shadow-soft disabled:cursor-not-allowed disabled:opacity-45"
                  onClick={submitJob}
                >
                  <UploadCloud aria-hidden="true" className="size-4" />
                  Process image
                </button>
              </div>
            )}

            {createMutation.isPending && (
              <UploadProgress
                progress={uploadProgress}
                onCancel={() => uploadControllerRef.current?.abort()}
              />
            )}

            {credentials && jobQuery.isPending && <WorkspaceLoading />}

            {credentials && jobQuery.data?.status === "completed" && (
              <>
                {resultQuery.isPending && (
                  <WorkspaceLoading label="Loading result…" />
                )}
                {resultQuery.data && (
                  <ResultPreview
                    key={credentials.id}
                    tool={tool}
                    result={resultQuery.data}
                    deleting={deleteMutation.isPending}
                    onStartAnother={() => deleteMutation.mutate(false)}
                  />
                )}
              </>
            )}

            {credentials &&
              jobQuery.data &&
              jobQuery.data.status !== "completed" && (
                <JobStatusPanel
                  job={jobQuery.data}
                  cancelling={cancelMutation.isPending}
                  deleting={deleteMutation.isPending}
                  onCancel={() => cancelMutation.mutate()}
                  onReset={(keepFile) => {
                    if (jobQuery.data.status === "expired") {
                      clearUnrecoverableSession();
                      return;
                    }
                    deleteMutation.mutate(keepFile);
                  }}
                />
              )}

            {credentials && jobQuery.isError && (
              <RequestError
                error={jobQuery.error}
                onClear={clearUnrecoverableSession}
              />
            )}

            {apiError &&
              !jobQuery.isError &&
              !createMutation.isPending &&
              !resultQuery.isPending && <InlineError error={apiError} />}
          </section>

          <WorkspaceAside tool={tool} job={jobQuery.data} />
        </div>
      )}
    </main>
  );
}

function UploadProgress({
  progress,
  onCancel,
}: {
  progress: number;
  onCancel: () => void;
}) {
  return (
    <section
      aria-live="polite"
      className="grid min-h-88 place-items-center rounded-lg border border-line bg-surface p-6 text-center"
    >
      <div className="w-full max-w-xl">
        <p className="text-sm font-medium text-ink-muted">
          Uploading / {progress}%
        </p>
        <h2 className="mt-3 text-2xl font-semibold tracking-[-0.035em]">
          Sending your file
        </h2>
        <div
          role="progressbar"
          aria-valuemin={0}
          aria-valuemax={100}
          aria-valuenow={progress}
          className="mt-8 h-2 overflow-hidden rounded bg-line"
        >
          <div
            className="h-full bg-accent-strong transition-[width] duration-200"
            style={{ width: `${progress}%` }}
          />
        </div>
        <button
          type="button"
          className="mt-8 min-h-11 rounded-md border border-line px-4 py-2 text-sm font-semibold hover:bg-canvas"
          onClick={onCancel}
        >
          Cancel upload
        </button>
      </div>
    </section>
  );
}

function WorkspaceLoading({ label = "Restoring job…" }: { label?: string }) {
  return (
        <div className="grid min-h-88 place-items-center rounded-lg border border-line bg-surface">
      <div className="text-center">
        <LoaderCircle
          aria-hidden="true"
          className="mx-auto size-8 animate-spin text-accent-strong"
        />
        <p className="mt-4 text-sm text-ink-muted">{label}</p>
      </div>
    </div>
  );
}

function OtherJobNotice({
  activeToolSlug,
  requestedToolTitle,
}: {
  activeToolSlug: string;
  requestedToolTitle: string;
}) {
  return (
    <section className="mt-6 grid min-h-[25rem] place-items-center border-2 border-ink bg-inverse p-8 text-center text-on-inverse">
      <div className="max-w-xl">
        <LockKeyhole
          aria-hidden="true"
          className="mx-auto size-10 text-accent"
        />
        <p className="mt-6 font-mono text-[0.7rem] uppercase tracking-[0.14em] text-accent">
          One job at a time
        </p>
        <h2 className="mt-4 text-3xl font-black uppercase tracking-[-0.05em]">
          Finish the active job first
        </h2>
        <p className="mt-4 leading-7 text-on-inverse/65">
          {requestedToolTitle} is ready, but another workspace still owns the
          current anonymous job.
        </p>
        <Link
          href={`/tools/${activeToolSlug}`}
          className="mt-8 inline-flex min-h-12 items-center border-2 border-on-inverse bg-accent px-5 py-3 text-sm font-bold text-[#171a17]"
        >
          Return to active job
        </Link>
      </div>
    </section>
  );
}

function RequestError({
  error,
  onClear,
}: {
  error: Error;
  onClear: () => void;
}) {
  const unrecoverable =
    error instanceof APIError && [401, 404].includes(error.status);

  return (
    <section className="border-2 border-danger bg-surface p-8 text-center">
      <CircleAlert aria-hidden="true" className="mx-auto size-9 text-danger" />
      <h2 className="mt-5 text-2xl font-bold">Job unavailable</h2>
      <p className="mt-3 text-sm leading-6 text-ink-muted">{error.message}</p>
      {unrecoverable && (
        <button
          type="button"
          className="mt-6 min-h-11 border-2 border-ink bg-accent px-4 py-2 text-sm font-bold text-[#171a17]"
          onClick={onClear}
        >
          Clear local session
        </button>
      )}
    </section>
  );
}

function InlineError({ error }: { error: Error }) {
  return (
    <div
      role="alert"
      className="mt-4 flex gap-3 border-2 border-danger bg-surface p-4 text-sm"
    >
      <CircleAlert aria-hidden="true" className="size-5 shrink-0 text-danger" />
      <span>{error.message}</span>
    </div>
  );
}

function WorkspaceAside({ tool, job }: { tool: MediaTool; job?: Job }) {
  return (
    <aside className="space-y-4">
      <dl className="overflow-hidden rounded-md border border-line bg-surface shadow-hard lg:sticky lg:top-6">
        <MetaRow term="Input" value={tool.accepts} />
        <MetaRow term="Output" value={tool.output} />
        {job && <MetaRow term="Status" value={job.status} last />}
      </dl>
    </aside>
  );
}

function previewStyle(operation: MediaTool["operation"], options: JobOptions) {
  if (operation === "image_flip") {
    return {
      transform:
        options.flip_direction === "vertical" ? "scaleY(-1)" : "scaleX(-1)",
    };
  }
  if (operation === "image_rotate") {
    return { transform: `rotate(${options.angle ?? 90}deg)` };
  }
  if (operation === "image_blur") {
    return { filter: `blur(${Math.min(options.strength ?? 2, 10)}px)` };
  }
  if (operation === "image_adjust") {
    return {
      filter: `brightness(${100 + (options.brightness ?? 0)}%) contrast(${100 + (options.contrast ?? 0)}%) saturate(${options.saturation ?? 100}%)`,
    };
  }
  return undefined;
}

function MetaRow({
  term,
  value,
  last = false,
}: {
  term: string;
  value: string;
  last?: boolean;
}) {
  return (
    <div className={`p-3.5 ${last ? "" : "border-b border-line"}`}>
      <dt className="text-[0.65rem] font-semibold text-ink-muted">
        {term}
      </dt>
      <dd className="mt-1 text-xs font-semibold capitalize">{value}</dd>
    </div>
  );
}
