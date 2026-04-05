import { createResource, Show, For } from "solid-js";
import { A } from "@solidjs/router";
import { fetchLabs } from "../api.ts";
import { Card, CardContent } from "../components/ui/Card.tsx";
import Badge from "../components/ui/Badge.tsx";
import LoadingSpinner from "../components/LoadingSpinner.tsx";
import ErrorBanner from "../components/ErrorBanner.tsx";

export default function LabListPage() {
  const [labs] = createResource(fetchLabs);

  return (
    <div class="space-y-6">
      <div>
        <h2 class="text-lg font-semibold text-white">Available Labs</h2>
        <p class="mt-1 text-sm text-slate-400">
          Select a lab to view the problem and submit your solution.
        </p>
      </div>

      <Show when={labs.loading}>
        <div class="flex items-center gap-3 text-slate-400 animate-pulse">
          <LoadingSpinner />
          <span>Loading labs…</span>
        </div>
      </Show>

      <Show when={labs.error}>
        <ErrorBanner message={`Failed to load labs: ${labs.error}`} onDismiss={() => {}} />
      </Show>

      <Show when={labs()}>
        {(labList) => (
          <div class="stagger-children space-y-3">
            <For each={labList()}>
              {(lab) => (
                <A href={`/labs/${lab.id}`} class="block group">
                  <Card class="transition-all duration-200 group-hover:border-blue-500/40 group-hover:bg-slate-800/50 group-hover:-translate-y-0.5 group-hover:shadow-lg group-hover:shadow-blue-500/10">
                    <CardContent class="p-5">
                      <div class="flex items-center justify-between">
                        <div class="flex-1">
                          <div class="flex items-center gap-2.5">
                            <h3 class="font-medium text-white group-hover:text-blue-400 transition-colors">
                              {lab.name}
                            </h3>
                            <Badge variant="outline">{lab.id}</Badge>
                          </div>
                          <p class="mt-1.5 text-sm text-slate-400 line-clamp-2">
                            {lab.problem_statement.replace(/\\n/g, " ")}
                          </p>
                        </div>
                        <svg
                          class="ml-4 h-5 w-5 shrink-0 text-slate-600 transition-all group-hover:text-blue-400 group-hover:translate-x-1"
                          viewBox="0 0 24 24"
                          fill="none"
                          stroke="currentColor"
                          stroke-width="2"
                        >
                          <path d="m9 18 6-6-6-6" />
                        </svg>
                      </div>
                    </CardContent>
                  </Card>
                </A>
              )}
            </For>
          </div>
        )}
      </Show>
    </div>
  );
}
