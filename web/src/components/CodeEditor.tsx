import { createSignal } from "solid-js";

type Props = {
  onCodeChange: (code: string) => void;
  placeholder?: string;
};

export default function CodeEditor(props: Props) {
  const [lineCount, setLineCount] = createSignal(1);

  const handleInput = (e: InputEvent & { currentTarget: HTMLTextAreaElement }) => {
    const val = e.currentTarget.value;
    props.onCodeChange(val);
    setLineCount(val.split("\n").length);
  };

  const handleKeyDown = (e: KeyboardEvent & { currentTarget: HTMLTextAreaElement }) => {
    // Tab inserts 4 spaces instead of moving focus
    if (e.key === "Tab") {
      e.preventDefault();
      const ta = e.currentTarget;
      const start = ta.selectionStart;
      const end = ta.selectionEnd;
      ta.value = ta.value.substring(0, start) + "    " + ta.value.substring(end);
      ta.selectionStart = ta.selectionEnd = start + 4;
      props.onCodeChange(ta.value);
    }
  };

  return (
    <div class="space-y-2">
      <label class="text-sm font-medium text-slate-400">Python Code</label>
      <div class="group relative overflow-hidden rounded-xl border border-slate-700 bg-slate-950 transition-all duration-200 focus-within:border-blue-500 focus-within:ring-2 focus-within:ring-blue-500/25">
        {/* Line numbers gutter */}
        <div class="flex">
          <div class="select-none border-r border-slate-800 bg-slate-900/50 px-3 py-4 text-right font-mono text-xs leading-6 text-slate-600">
            {Array.from({ length: Math.max(lineCount(), 12) }, (_, i) => (
              <div>{i + 1}</div>
            ))}
          </div>

          {/* Editor */}
          <textarea
            class="min-h-[300px] w-full resize-none bg-transparent px-4 py-4 font-mono text-sm leading-6 text-slate-200 placeholder:text-slate-600 focus:outline-none"
            placeholder={
              props.placeholder ?? "def solution():\n    # Write your code here\n    pass"
            }
            spellcheck={false}
            onInput={handleInput}
            onKeyDown={handleKeyDown}
          />
        </div>

        {/* Bottom bar */}
        <div class="flex items-center justify-between border-t border-slate-800 bg-slate-900/30 px-4 py-1.5">
          <span class="text-[10px] font-medium uppercase tracking-wider text-slate-600">
            Python
          </span>
          <span class="text-[10px] text-slate-600">{lineCount()} lines</span>
        </div>
      </div>
    </div>
  );
}
