import { z } from "zod";

export const jobOperationSchema = z.enum([
  "video_grayscale",
  "video_extract_audio",
  "video_remove_audio",
  "video_convert",
  "video_clip",
  "image_grayscale",
  "image_convert",
  "image_compress",
  "image_resize",
]);

export const jobStatusSchema = z.enum([
  "queued",
  "processing",
  "cancelling",
  "completed",
  "failed",
  "cancelled",
  "expired",
]);

export type JobStatus = z.infer<typeof jobStatusSchema>;

export const jobSchema = z.object({
  id: z.string().min(1),
  operation: jobOperationSchema,
  status: jobStatusSchema,
  progress: z.number().int().min(0).max(100),
  attempt: z.number().int().min(0),
  expires_at: z.string().nullish(),
  result_ready: z.boolean(),
  failure_code: z.string().optional(),
});

export const createdJobSchema = jobSchema.extend({
  owner_token: z.string().min(1),
});

export const errorEnvelopeSchema = z.object({
  error: z.object({
    code: z.string(),
    message: z.string(),
    request_id: z.string().optional(),
  }),
});

export type Job = z.infer<typeof jobSchema>;
export type CreatedJob = z.infer<typeof createdJobSchema>;

export function isActiveStatus(status: JobStatus): boolean {
  return (
    status === "queued" ||
    status === "processing" ||
    status === "cancelling"
  );
}

export function isTerminalStatus(status: JobStatus): boolean {
  return !isActiveStatus(status);
}
