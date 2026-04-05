import Alert from "./ui/Alert.tsx";

type Props = { text: string };

export default function ProblemStatement(props: Props) {
  const formatted = () => props.text.replace(/\\n/g, "<br>");

  return (
    <Alert variant="info">
      <div class="flex items-start gap-3">
        <svg
          class="mt-0.5 h-4 w-4 shrink-0 text-blue-400"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <circle cx="12" cy="12" r="10" />
          <path d="M12 16v-4" />
          <path d="M12 8h.01" />
        </svg>
        <div>
          <p class="mb-1 text-xs font-semibold uppercase tracking-wider text-blue-400">
            Problem Statement
          </p>
          <p class="text-sm leading-relaxed text-slate-300" innerHTML={formatted()} />
        </div>
      </div>
    </Alert>
  );
}
