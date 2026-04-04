import { fetchLabs, submitFile } from "./api.ts";
import {
  buildLabOptions,
  buildResultsFragment,
  formatProblemStatement,
  validateFile,
} from "./ui.ts";

// --- DOM refs ---
const labSelect = document.getElementById("lab-select") as HTMLSelectElement;
const problemBox = document.getElementById("problem-statement") as HTMLDivElement;
const problemText = document.getElementById("problem-text") as HTMLParagraphElement;
const form = document.getElementById("submit-form") as HTMLFormElement;
const fileInput = document.getElementById("file-input") as HTMLInputElement;
const submitBtn = document.getElementById("submit-btn") as HTMLButtonElement;
const resultsDiv = document.getElementById("results") as HTMLDivElement;
const resultsBody = document.getElementById("results-body") as HTMLDivElement;
const errorBox = document.getElementById("error-box") as HTMLDivElement;

// --- Init ---
async function init(): Promise<void> {
  try {
    const labs = await fetchLabs();
    labSelect.innerHTML = "";
    for (const opt of buildLabOptions(labs)) {
      labSelect.appendChild(opt);
    }
  } catch (err) {
    showError(`Failed to load labs: ${err}`);
  }
}

// --- Events ---
labSelect.addEventListener("change", () => {
  const selected = labSelect.selectedOptions[0];
  const problem = selected?.dataset.problem;
  if (problem) {
    problemText.innerHTML = formatProblemStatement(problem);
    problemBox.classList.remove("hidden");
  } else {
    problemBox.classList.add("hidden");
  }
  hideResults();
  hideError();
});

form.addEventListener("submit", async (e: Event) => {
  e.preventDefault();
  hideResults();
  hideError();

  const labId = labSelect.value;
  if (!labId) {
    showError("Please select a lab.");
    return;
  }

  const file = fileInput.files?.[0];
  const validationError = validateFile(file);
  if (validationError) {
    showError(validationError);
    return;
  }

  submitBtn.disabled = true;
  submitBtn.textContent = "Grading…";

  try {
    const result = await submitFile(labId, file!);
    resultsBody.innerHTML = "";
    resultsBody.appendChild(buildResultsFragment(result));
    resultsDiv.classList.remove("hidden");
  } catch (err) {
    showError(`${err}`);
  } finally {
    submitBtn.disabled = false;
    submitBtn.textContent = "Submit";
  }
});

// --- Helpers ---
function showError(msg: string): void {
  errorBox.textContent = msg;
  errorBox.classList.remove("hidden");
}

function hideError(): void {
  errorBox.classList.add("hidden");
}

function hideResults(): void {
  resultsDiv.classList.add("hidden");
}

// --- Boot ---
init();
