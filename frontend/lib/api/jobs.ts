import type { JobOperation } from "@/config/tools";
import type { JobOptions } from "@/features/media-jobs/options";
import {
  createdJobSchema,
  errorEnvelopeSchema,
  jobSchema,
  type CreatedJob,
  type Job,
} from "@/lib/api/contracts";

export class APIError extends Error {
  constructor(
    message: string,
    readonly status: number,
    readonly code: string,
    readonly requestID?: string,
  ) {
    super(message);
    this.name = "APIError";
  }
}

export interface CreateJobInput {
  operation: JobOperation;
  options: JobOptions;
  file: File;
  signal?: AbortSignal;
  onProgress?: (percentage: number) => void;
}

export interface ResultFile {
  blob: Blob;
  filename: string;
}

export function createJob(input: CreateJobInput): Promise<CreatedJob> {
  return new Promise((resolve, reject) => {
    const request = new XMLHttpRequest();
    const body = new FormData();
    body.append("operation", input.operation);
    body.append("options", JSON.stringify(input.options));
    body.append("file", input.file, input.file.name);

    const abort = () => request.abort();
    input.signal?.addEventListener("abort", abort, { once: true });

    request.open("POST", "/api/v1/jobs");
    request.setRequestHeader("Accept", "application/json");

    request.upload.addEventListener("progress", (event) => {
      if (!event.lengthComputable) {
        return;
      }
      input.onProgress?.(
        Math.min(100, Math.round((event.loaded / event.total) * 100)),
      );
    });

    request.addEventListener("load", () => {
      input.signal?.removeEventListener("abort", abort);

      if (request.status < 200 || request.status >= 300) {
        reject(parseAPIError(request.responseText, request.status));
        return;
      }

      try {
        resolve(createdJobSchema.parse(JSON.parse(request.responseText)));
      } catch {
        reject(
          new APIError(
            "The server returned an invalid job response.",
            request.status,
            "invalid_response",
          ),
        );
      }
    });

    request.addEventListener("error", () => {
      input.signal?.removeEventListener("abort", abort);
      reject(new APIError("The API could not be reached.", 0, "network_error"));
    });

    request.addEventListener("abort", () => {
      input.signal?.removeEventListener("abort", abort);
      reject(new DOMException("Upload cancelled.", "AbortError"));
    });

    request.send(body);
  });
}

export async function getJob(
  id: string,
  ownerToken: string,
  signal?: AbortSignal,
): Promise<Job> {
  const data = await requestJSON(`/api/v1/jobs/${encodeURIComponent(id)}`, {
    ownerToken,
    signal,
  });
  return jobSchema.parse(data);
}

export async function cancelJob(
  id: string,
  ownerToken: string,
): Promise<Job> {
  const data = await requestJSON(
    `/api/v1/jobs/${encodeURIComponent(id)}/cancel`,
    { method: "POST", ownerToken },
  );
  return jobSchema.parse(data);
}

export async function deleteJob(
  id: string,
  ownerToken: string,
): Promise<void> {
  await requestJSON(`/api/v1/jobs/${encodeURIComponent(id)}`, {
    method: "DELETE",
    ownerToken,
    expectJSON: false,
  });
}

export async function getJobResult(
  id: string,
  ownerToken: string,
  signal?: AbortSignal,
): Promise<ResultFile> {
  const response = await fetch(
    `/api/v1/jobs/${encodeURIComponent(id)}/result`,
    {
      headers: authorizationHeaders(ownerToken),
      signal,
      cache: "no-store",
    },
  );

  if (!response.ok) {
    throw await responseAPIError(response);
  }

  return {
    blob: await response.blob(),
    filename:
      filenameFromDisposition(response.headers.get("content-disposition")) ??
      `recode-${id}`,
  };
}

interface RequestOptions {
  method?: "GET" | "POST" | "DELETE";
  ownerToken: string;
  signal?: AbortSignal;
  expectJSON?: boolean;
}

async function requestJSON(
  url: string,
  {
    method = "GET",
    ownerToken,
    signal,
    expectJSON = true,
  }: RequestOptions,
): Promise<unknown> {
  const response = await fetch(url, {
    method,
    headers: authorizationHeaders(ownerToken),
    signal,
    cache: "no-store",
  });

  if (!response.ok) {
    throw await responseAPIError(response);
  }

  if (!expectJSON || response.status === 204) {
    return undefined;
  }
  return response.json();
}

function authorizationHeaders(ownerToken: string): HeadersInit {
  return {
    Accept: "application/json",
    Authorization: `Bearer ${ownerToken}`,
  };
}

async function responseAPIError(response: Response): Promise<APIError> {
  return parseAPIError(await response.text(), response.status);
}

function parseAPIError(body: string, status: number): APIError {
  try {
    const envelope = errorEnvelopeSchema.parse(JSON.parse(body));
    return new APIError(
      envelope.error.message,
      status,
      envelope.error.code,
      envelope.error.request_id,
    );
  } catch {
    return new APIError(
      "The request could not be completed.",
      status,
      "unexpected_error",
    );
  }
}

function filenameFromDisposition(value: string | null): string | null {
  if (!value) {
    return null;
  }

  const match = /filename="([^"]+)"/i.exec(value);
  return match?.[1] ?? null;
}

