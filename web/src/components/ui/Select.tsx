import { Select as KSelect } from "@kobalte/core/select";

import { cn } from "../../lib/cn.ts";

type SelectOption = { value: string; label: string; description?: string };

type SelectProps = {
  options: SelectOption[];
  placeholder?: string;
  value?: string;
  onChange?: (value: string) => void;
  class?: string;
};

export default function Select(props: SelectProps) {
  return (
    <KSelect<SelectOption>
      options={props.options}
      optionValue="value"
      optionTextValue="label"
      value={props.options.find((o) => o.value === props.value)}
      onChange={(opt) => {
        if (opt) props.onChange?.(opt.value);
      }}
      placeholder={props.placeholder ?? "Select…"}
      itemComponent={(itemProps) => (
        <KSelect.Item
          item={itemProps.item}
          class="relative flex cursor-pointer select-none items-center rounded-lg px-3 py-2.5 text-sm text-slate-300 outline-none transition-colors data-[highlighted]:bg-slate-700/50 data-[highlighted]:text-white data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
        >
          <KSelect.ItemLabel class="flex-1">{itemProps.item.rawValue.label}</KSelect.ItemLabel>
          <KSelect.ItemIndicator class="ml-2 text-blue-400">
            <svg
              class="h-4 w-4"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <polyline points="20 6 9 17 4 12" />
            </svg>
          </KSelect.ItemIndicator>
        </KSelect.Item>
      )}
    >
      <KSelect.Trigger
        class={cn(
          "inline-flex w-full items-center justify-between rounded-lg border border-slate-700 bg-slate-900 px-3 py-2.5 text-sm text-slate-200 shadow-sm transition-all duration-200 hover:border-slate-600 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/25 data-[placeholder-shown]:text-slate-500",
          props.class,
        )}
      >
        <KSelect.Value<SelectOption>>{(state) => state.selectedOption()?.label}</KSelect.Value>
        <KSelect.Icon class="ml-2 text-slate-500">
          <svg
            class="h-4 w-4"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="m6 9 6 6 6-6" />
          </svg>
        </KSelect.Icon>
      </KSelect.Trigger>
      <KSelect.Portal>
        <KSelect.Content class="animate-scale-in z-50 overflow-hidden rounded-xl border border-slate-700 bg-slate-800 shadow-2xl shadow-black/50">
          <KSelect.Listbox class="max-h-64 overflow-y-auto p-1.5" />
        </KSelect.Content>
      </KSelect.Portal>
    </KSelect>
  );
}
