import { type ParentProps } from "solid-js";
import Header from "./Header.tsx";

export default function Layout(props: ParentProps) {
  return (
    <div class="mx-auto max-w-2xl px-4 py-10 animate-fade-in">
      <Header />
      <main class="mt-8">{props.children}</main>
      <footer class="mt-16 text-center text-xs text-slate-600">
        Autograder — code runs sandboxed in Docker with no network access
      </footer>
    </div>
  );
}
