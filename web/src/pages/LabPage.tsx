import { createSignal, createResource, Show } from "solid-js";
import { useParams, A } from "@solidjs/router";
import { fetchLab, submitSolution } from "../api.ts";
import type { GradeResult } from "../models.ts";
import ProblemStatement from "../components/ProblemStatement.tsx";
import FileUpload from "../components/FileUpload.tsx";
import CodeEditor from "../components/CodeEditor.tsx";
import Results from "../components/Results.tsx";
import ErrorBanner from "../components/ErrorBanner.tsx";
import LoadingSpinner from "../components/LoadingSpinner.tsx";
import Button from "../components/ui/Button.tsx";

type Mode = "editor" | "upload";

export default function LabPage() {
  const params = useParams<{ id: string }>();
  const [lab] = createResource(() => params.id, fetchLab);
  const [mode, setMode] = createSignal<Mode>("editor");
  const [code, setCode] = createSignal("");
  const [file, setFile] = createSignal<File | null>(null);
  const [result, setResult] = createSignal<GradeResult | null>(null);
  const [error, setError] = createSignal<string | null>(null);
  const [submitting, setSubmitting] = createSignal(false);

  const handleSubmit = async () => {
    const currentLab = lab();
    if (!currentLab) return;
    setError(null);
    setResult(null);

    if (mode() === "editor") {
      const c = code().trim();
      if (!c) {
        setError("Please write some code before submitting.");
        return;
      }
    } else {
      const f = file();
      if (!f) {
        setError("Please select a file.");
        return;
      }
      if (!f.name.endsWith(".py")) {
        setError("Only .py files are allowed.");
        return;
      }
      if (f.size > 20 * 1024) {
        setError("File too large (max 20 KB).");
        return;
      }
    }

    setSubmitting(true);
    try {
      const res =
        mode() === "editor"
          ? await submitSolution({ labId: currentLab.id, code: code() })
          : await submitSolution({ labId: currentLab.id, file: file()! });
      setResult(res);
    } catch (err) {
      setError(String(err));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div class="space-y-6">
      {/* Back link */}
      <A
        href="/"
        class="inline-flex items-center gap-1.5 text-sm text-slate-400 transition-colors hover:text-white"
      >
        <svg class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="m15 18-6-6 6-6" />
        </svg>
        All Labs
      </A>

      <Show when={lab.loading}>
        <div class="flex items-center gap-3 text-slate-400 animate-pulse">
          <LoadingSpinner />
          <span>Loading lab…</span>
        </div>
      </Show>

      <Show when={lab.error}>
        <ErrorBanner message={`${lab.error}`} onDismiss={() => {}} />
      </Show>

      <Show when={lab()}>
        {(currentLab) => (
          <div class="space-y-6 animate-fade-in">
            {/* Title */}
            <div>
              <h2 class="text-xl font-semibold text-white">{currentLab().name}</h2>
              <p class="mt-1 text-sm text-slate-500">{currentLab().id}</p>
            </div>

            {/* Problem */}
            <div class="animate-slide-down">
              <ProblemStatement text={currentLab().problem_statement} />
            </div>

            {/* Mode tabs */}
            <div class="flex gap-1 rounded-lg bg-slate-900 p-1">
              <button
                class={`flex-1 rounded-md px-4 py-2 text-sm font-medium transition-all duration-200 ${
                  mode() === "editor"
                    ? "bg-slate-800 text-white shadow-sm"
                    : "text-slate-400 hover:text-slate-200"
                }`}
                onClick={() => setMode("editor")}
              >
                <span class="inline-flex items-center gap-2">
                  <svg
                    class="h-4 w-4"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <polyline points="16 18 22 12 16 6" />
                    <polyline points="8 6 2 12 8 18" />
                  </svg>
                  Write Code
                </span>
              </button>
              <button
                class={`flex-1 rounded-md px-4 py-2 text-sm font-medium transition-all duration-200 ${
                  mode() === "upload"
                    ? "bg-slate-800 text-white shadow-sm"
                    : "text-slate-400 hover:text-slate-200"
                }`}
                onClick={() => setMode("upload")}
              >
                <span class="inline-flex items-center gap-2">
                  <svg
                    class="h-4 w-4"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                  >
                    <path d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                  </svg>
                  Upload File
                </span>
              </button>
            </div>

            {/* Editor or Upload */}
            <div class="animate-fade-in">
              <Show when={mode() === "editor"}>
                <CodeEditor onCodeChange={setCode} />
              </Show>
              <Show when={mode() === "upload"}>
                <FileUpload onFileChange={setFile} />
              </Show>
            </div>

            {/* Submit */}
            <Button onClick={handleSubmit} disabled={submitting()} size="lg" class="w-full">
              <Show
                when={submitting()}
                fallback={
                  <>
                    <svg
                      class="h-4 w-4"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                      stroke-width="2"
                    >
                      <path d="m5 12 5 5L20 7" />
                    </svg>
                    Submit Solution
                  </>
                }
              >
                <LoadingSpinner />
                Grading…
              </Show>
            </Button>

            {/* Error */}
            <Show when={error()}>
              {(msg) => (
                <div class="animate-scale-in">
                  <ErrorBanner message={msg()} onDismiss={() => setError(null)} />
                </div>
              )}
            </Show>

            {/* Results */}
            <Show when={result()}>
              {(res) => (
                <div class="animate-fade-in">
                  <Results result={res()} />
                </div>
              )}
            </Show>
          </div>
        )}
      </Show>
    </div>
  );
}
