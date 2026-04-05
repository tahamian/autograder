export default function Header() {
  return (
    <header class="space-y-2">
      <div class="flex items-center gap-3">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-blue-700 shadow-lg shadow-blue-500/25">
          <svg
            class="h-5 w-5 text-white"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="M14.5 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V7.5L14.5 2z" />
            <polyline points="14 2 14 8 20 8" />
            <path d="m9 15 2 2 4-4" />
          </svg>
        </div>
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-white">Autograder</h1>
          <p class="text-sm text-slate-400">Upload your Python script and get instant feedback</p>
        </div>
      </div>
    </header>
  );
}
