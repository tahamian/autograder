import { For } from "solid-js";
import type { GradeResult, Evaluation } from "../models.ts";
import { Card, CardHeader, CardTitle, CardContent } from "./ui/Card.tsx";
import Badge from "./ui/Badge.tsx";

type Props = { result: GradeResult };

function EvalCard(props: { eval: Evaluation; index: number }) {
  const pass = () => props.eval.points > 0;

  return (
    <div
      class="group rounded-xl border bg-slate-900/50 p-4 transition-all duration-200 hover:bg-slate-800/50"
      classList={{
        "border-emerald-500/30 hover:border-emerald-500/50": pass(),
        "border-red-500/30 hover:border-red-500/50": !pass(),
      }}
      style={{ "animation-delay": `${props.index * 80}ms` }}
    >
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div
            class="flex h-8 w-8 items-center justify-center rounded-lg text-sm"
            classList={{
              "bg-emerald-500/15 text-emerald-400": pass(),
              "bg-red-500/15 text-red-400": !pass(),
            }}
          >
            {pass() ? (
              <svg
                class="h-4 w-4"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.5"
              >
                <polyline points="20 6 9 17 4 12" />
              </svg>
            ) : (
              <svg
                class="h-4 w-4"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.5"
              >
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            )}
          </div>
          <div>
            <p class="text-sm font-medium text-white">{props.eval.name}</p>
            <p class="text-xs text-slate-500">{props.eval.status}</p>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <Badge variant="secondary">{props.eval.type}</Badge>
          <Badge variant={pass() ? "success" : "destructive"}>{props.eval.points} pts</Badge>
        </div>
      </div>
      {props.eval.actual !== null && (
        <div class="mt-3 rounded-lg bg-slate-950/50 px-3 py-2">
          <span class="text-xs text-slate-500">Output: </span>
          <code class="text-xs text-slate-300">{String(props.eval.actual)}</code>
        </div>
      )}
    </div>
  );
}

export default function Results(props: Props) {
  const total = () => props.result.total_points;
  const allPassed = () => props.result.evaluations.every((e) => e.points > 0);

  return (
    <Card class={allPassed() ? "border-emerald-500/20 animate-pulse-glow" : ""}>
      <CardHeader>
        <div class="flex items-center justify-between">
          <CardTitle>
            <span class="flex items-center gap-2">
              {allPassed() ? <span class="text-emerald-400">🎉</span> : <span>📋</span>}
              Results
            </span>
          </CardTitle>
          <div class="flex items-center gap-2">
            <span
              class="text-2xl font-bold tabular-nums"
              classList={{
                "text-emerald-400": allPassed(),
                "text-white": !allPassed(),
              }}
            >
              {total()}
            </span>
            <span class="text-sm text-slate-500">points</span>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div class="stagger-children space-y-3">
          <For each={props.result.evaluations}>{(ev, i) => <EvalCard eval={ev} index={i()} />}</For>
        </div>
      </CardContent>
    </Card>
  );
}
