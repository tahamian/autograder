import { JSDOM } from "jsdom";

/**
 * Creates a minimal JSDOM environment and injects its globals
 * so that document.createElement etc. work in tests.
 */
export function setupDOM(html = "<!DOCTYPE html><html><body></body></html>"): JSDOM {
  const dom = new JSDOM(html, { url: "http://localhost" });
  globalThis.document = dom.window.document as unknown as Document;
  globalThis.window = dom.window as unknown as Window & typeof globalThis;
  globalThis.HTMLElement = dom.window.HTMLElement;
  globalThis.HTMLSelectElement = dom.window.HTMLSelectElement;
  globalThis.HTMLOptionElement = dom.window.HTMLOptionElement;
  globalThis.HTMLDivElement = dom.window.HTMLDivElement;
  globalThis.HTMLInputElement = dom.window.HTMLInputElement;
  globalThis.HTMLButtonElement = dom.window.HTMLButtonElement;
  globalThis.DocumentFragment = dom.window.DocumentFragment;
  globalThis.File = dom.window.File as unknown as typeof File;
  return dom;
}

/**
 * Tears down the DOM globals.
 */
export function teardownDOM(): void {
  // @ts-expect-error cleanup global
  delete globalThis.document;
  // @ts-expect-error cleanup global
  delete globalThis.window;
}
