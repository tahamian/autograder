import { type JSX, splitProps } from "solid-js";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "../../lib/cn.ts";

const alertVariants = cva("relative w-full rounded-xl border p-4 text-sm", {
  variants: {
    variant: {
      default: "border-slate-800 bg-slate-900/50 text-slate-200",
      destructive: "border-red-500/30 bg-red-500/10 text-red-400",
      success: "border-emerald-500/30 bg-emerald-500/10 text-emerald-400",
      info: "border-blue-500/30 bg-blue-500/10 text-blue-400",
    },
  },
  defaultVariants: { variant: "default" },
});

type AlertProps = JSX.HTMLAttributes<HTMLDivElement> & VariantProps<typeof alertVariants>;

export default function Alert(props: AlertProps) {
  const [local, rest] = splitProps(props, ["variant", "class", "children"]);
  return (
    <div role="alert" class={cn(alertVariants({ variant: local.variant }), local.class)} {...rest}>
      {local.children}
    </div>
  );
}
