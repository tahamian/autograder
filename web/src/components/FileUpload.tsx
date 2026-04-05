import { createSignal, Show } from "solid-js";

type Props = {
  onFileChange: (file: File | null) => void;
};

export default function FileUpload(props: Props) {
  const [fileName, setFileName] = createSignal<string | null>(null);
  const [dragging, setDragging] = createSignal(false);

  let inputRef!: HTMLInputElement;

  const handleFile = (files: FileList | null) => {
    const file = files?.[0] ?? null;
    setFileName(file?.name ?? null);
    props.onFileChange(file);
  };

  return (
    <div class="space-y-2">
      <label class="text-sm font-medium text-slate-400">Python File</label>
      <div
        class={`group relative flex cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed px-6 py-8 transition-all duration-200 ${
          dragging()
            ? "border-blue-500 bg-blue-500/10"
            : "border-slate-700 bg-slate-900/50 hover:border-slate-500 hover:bg-slate-800/50"
        }`}
        onClick={() => inputRef.click()}
        onDragOver={(e) => {
          e.preventDefault();
          setDragging(true);
        }}
        onDragLeave={() => setDragging(false)}
        onDrop={(e) => {
          e.preventDefault();
          setDragging(false);
          handleFile(e.dataTransfer?.files ?? null);
        }}
      >
        <input
          ref={inputRef}
          type="file"
          accept=".py"
          class="hidden"
          onChange={(e) => handleFile(e.currentTarget.files)}
        />

        <Show
          when={fileName()}
          fallback={
            <>
              <div class="mb-3 flex h-12 w-12 items-center justify-center rounded-xl bg-slate-800 transition-colors group-hover:bg-slate-700">
                <svg
                  class="h-6 w-6 text-slate-400"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.5"
                >
                  <path d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                </svg>
              </div>
              <p class="text-sm text-slate-400">
                <span class="font-medium text-blue-400">Click to upload</span> or drag and drop
              </p>
              <p class="mt-1 text-xs text-slate-500">Python files only (.py) — max 20 KB</p>
            </>
          }
        >
          {(name) => (
            <div class="flex items-center gap-3">
              <div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-500/15">
                <svg
                  class="h-5 w-5 text-blue-400"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path d="M14.5 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V7.5L14.5 2z" />
                  <polyline points="14 2 14 8 20 8" />
                </svg>
              </div>
              <div>
                <p class="text-sm font-medium text-white">{name()}</p>
                <p class="text-xs text-slate-500">Click or drop to replace</p>
              </div>
            </div>
          )}
        </Show>
      </div>
    </div>
  );
}
