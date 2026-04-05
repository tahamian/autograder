// Re-export model types from the single source of truth.
export type { Lab, Evaluation, GradeResult, ApiError } from "./models.ts";

import type { Lab, GradeResult, ApiError } from "./models.ts";

const BASE = "/api";

export async function fetchLabs(): Promise<Lab[]> {
  const res = await fetch(`${BASE}/labs`);
  if (!res.ok) throw new Error("Failed to fetch labs");
  return res.json();
}

export async function fetchLab(id: string): Promise<Lab> {
  const res = await fetch(`${BASE}/labs/${id}`);
  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error((data as ApiError).error || `Lab "${id}" not found`);
  }
  return res.json();
}

export type SubmitParams =
  | { labId: string; file: File }
  | { labId: string; code: string; filename?: string };

export async function submitSolution(params: SubmitParams): Promise<GradeResult> {
  const form = new FormData();
  form.append("lab_id", params.labId);

  if ("file" in params) {
    form.append("file", params.file);
  } else {
    form.append("code", params.code);
    if (params.filename) {
      form.append("filename", params.filename);
    }
  }

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

// Legacy alias
export const submitFile = (labId: string, file: File) => submitSolution({ labId, file });
