// Re-export model types from the single source of truth.
export type { Lab, Evaluation, GradeResult, ApiError } from "./models.ts";

import type { Lab, GradeResult, ApiError } from "./models.ts";

const BASE = "/api";

export async function fetchLabs(): Promise<Lab[]> {
  const res = await fetch(`${BASE}/labs`);
  if (!res.ok) throw new Error("Failed to fetch labs");
  return res.json();
}

export async function submitFile(labId: string, file: File): Promise<GradeResult> {
  const form = new FormData();
  form.append("lab_id", labId);
  form.append("file", file);

  const res = await fetch(`${BASE}/submit`, {
    method: "POST",
    body: form,
  });

  const data = await res.json();
  if (!res.ok) {
    throw new Error((data as ApiError).error || "Submission failed");
  }
  return data as GradeResult;
}
