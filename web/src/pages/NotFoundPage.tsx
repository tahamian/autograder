import { A } from "@solidjs/router";
import Button from "../components/ui/Button.tsx";

export default function NotFoundPage() {
  return (
    <div class="flex flex-col items-center justify-center py-20 text-center animate-fade-in">
      <div class="text-6xl font-bold text-slate-700">404</div>
      <p class="mt-3 text-lg text-slate-400">Page not found</p>
      <p class="mt-1 text-sm text-slate-500">The page you're looking for doesn't exist.</p>
      <A href="/" class="mt-6">
        <Button variant="outline">
          <svg
            class="h-4 w-4"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="m15 18-6-6 6-6" />
          </svg>
          Back to Labs
        </Button>
      </A>
    </div>
  );
}
