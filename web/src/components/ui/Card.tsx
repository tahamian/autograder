import { type JSX, splitProps } from "solid-js";
import { cn } from "../../lib/cn.ts";

type CardProps = JSX.HTMLAttributes<HTMLDivElement>;

export function Card(props: CardProps) {
  const [local, rest] = splitProps(props, ["class", "children"]);
  return (
    <div
      class={cn(
        "rounded-xl border border-slate-800 bg-slate-900/50 backdrop-blur-sm shadow-xl",
        local.class,
      )}
      {...rest}
    >
      {local.children}
    </div>
  );
}

export function CardHeader(props: CardProps) {
  const [local, rest] = splitProps(props, ["class", "children"]);
  return (
    <div class={cn("flex flex-col space-y-1.5 p-6", local.class)} {...rest}>
      {local.children}
    </div>
  );
}

export function CardTitle(props: CardProps) {
  const [local, rest] = splitProps(props, ["class", "children"]);
  return (
    <h3
      class={cn("text-lg font-semibold leading-none tracking-tight text-white", local.class)}
      {...rest}
    >
      {local.children}
    </h3>
  );
}

export function CardContent(props: CardProps) {
  const [local, rest] = splitProps(props, ["class", "children"]);
  return (
    <div class={cn("p-6 pt-0", local.class)} {...rest}>
      {local.children}
    </div>
  );
}
