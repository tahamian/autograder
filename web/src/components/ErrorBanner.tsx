import Alert from "./ui/Alert.tsx";

type Props = {
  message: string;
  onDismiss: () => void;
};

export default function ErrorBanner(props: Props) {
  return (
    <Alert variant="destructive">
      <div class="flex items-start justify-between gap-3">
        <div class="flex items-start gap-3">
          <svg
            class="mt-0.5 h-4 w-4 shrink-0"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <circle cx="12" cy="12" r="10" />
            <line x1="15" y1="9" x2="9" y2="15" />
            <line x1="9" y1="9" x2="15" y2="15" />
          </svg>
          <p class="text-sm">{props.message}</p>
        </div>
        <button
          onClick={props.onDismiss}
          class="shrink-0 rounded-lg p-1 text-red-400 transition-colors hover:bg-red-500/20 hover:text-red-300"
        >
          <svg
            class="h-4 w-4"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>
      </div>
    </Alert>
  );
}
