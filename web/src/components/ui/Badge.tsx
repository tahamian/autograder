import { type JSX, splitProps } from "solid-js";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "../../lib/cn.ts";

const badgeVariants = cva(
  "inline-flex items-center rounded-md px-2.5 py-0.5 text-xs font-semibold transition-colors",
  {
    variants: {
      variant: {
        default: "bg-slate-800 text-slate-200",
        success: "bg-emerald-500/15 text-emerald-400 border border-emerald-500/25",
        destructive: "bg-red-500/15 text-red-400 border border-red-500/25",
        outline: "border border-slate-700 text-slate-300",
        secondary: "bg-slate-800 text-slate-400",
      },
    },
    defaultVariants: { variant: "default" },
  },
);

type BadgeProps = JSX.HTMLAttributes<HTMLSpanElement> & VariantProps<typeof badgeVariants>;

export default function Badge(props: BadgeProps) {
  const [local, rest] = splitProps(props, ["variant", "class", "children"]);
  return (
    <span class={cn(badgeVariants({ variant: local.variant }), local.class)} {...rest}>
      {local.children}
    </span>
  );
}
