import type { Lab, GradeResult, Evaluation } from "./api.ts";

/**
 * Build <option> elements for the lab <select>.
 * Returns the options to append (including the placeholder).
 */
export function buildLabOptions(labs: Lab[]): HTMLOptionElement[] {
  const opts: HTMLOptionElement[] = [];

  const placeholder = document.createElement("option");
  placeholder.value = "";
  placeholder.disabled = true;
  placeholder.selected = true;
  placeholder.textContent = "Select a lab…";
  opts.push(placeholder);

  for (const lab of labs) {
    const opt = document.createElement("option");
    opt.value = lab.id;
    opt.textContent = lab.name;
    opt.dataset.problem = lab.problem_statement;
    opts.push(opt);
  }

  return opts;
}

/**
 * Render the problem statement text, replacing literal \n with <br>.
 */
export function formatProblemStatement(raw: string): string {
  return raw.replace(/\\n/g, "<br>");
}

/**
 * Build a single evaluation card element.
 */
export function buildEvalCard(ev: Evaluation): HTMLDivElement {
  const pass = ev.points > 0;
  const card = document.createElement("div");
  card.className = `eval-card rounded-lg border bg-slate-900 px-4 py-3 ${
    pass
      ? "border-l-[3px] border-l-green-500 border-slate-700"
      : "border-l-[3px] border-l-red-500 border-slate-700"
  }`;
  card.setAttribute("data-testid", `eval-${ev.name}`);
  card.innerHTML = `
    <div class="flex items-center gap-3 text-sm">
      <span class="eval-type rounded bg-slate-700 px-2 py-0.5 text-xs font-semibold uppercase">${ev.type}</span>
      <span class="eval-name flex-1 font-medium">${ev.name}</span>
      <span class="eval-points font-bold ${pass ? "text-green-400" : "text-red-400"}">${ev.points} pts</span>
    </div>
    <div class="mt-1.5 flex gap-4 text-sm text-slate-400">
      <span class="eval-status">${ev.status}</span>
      ${ev.actual !== null ? `<span class="eval-actual">Output: <code class="rounded bg-slate-950 px-1.5 py-0.5 text-xs">${ev.actual}</code></span>` : ""}
    </div>
  `;
  return card;
}

/**
 * Build the full results container content from a GradeResult.
 * Returns a DocumentFragment ready to append.
 */
export function buildResultsFragment(result: GradeResult): DocumentFragment {
  const frag = document.createDocumentFragment();

  for (const ev of result.evaluations) {
    frag.appendChild(buildEvalCard(ev));
  }

  const total = document.createElement("div");
  total.className = "total mt-4 text-right text-lg font-bold text-white";
  total.textContent = `Total: ${result.total_points} points`;
  frag.appendChild(total);

  return frag;
}

/**
 * Validate a file before submission.
 * Returns null if valid, or an error message string.
 */
export function validateFile(
  file: File | undefined,
  allowedExtensions: string[] = ["py"],
  maxSizeBytes: number = 20 * 1024,
): string | null {
  if (!file) {
    return "Please select a file.";
  }
  const ext = file.name.split(".").pop() ?? "";
  if (!allowedExtensions.includes(ext)) {
    return `Only .${allowedExtensions.join(", .")} files are allowed.`;
  }
  if (file.size > maxSizeBytes) {
    return `File too large (max ${Math.round(maxSizeBytes / 1024)} KB).`;
  }
  return null;
}
