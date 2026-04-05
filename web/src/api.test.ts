import { describe, it, beforeEach, afterEach } from "@std/testing/bdd";
import { expect } from "chai";
import { fetchLabs, fetchLab, submitFile, submitSolution } from "./api.ts";

describe("api", () => {
  let originalFetch: typeof globalThis.fetch;

  beforeEach(() => {
    originalFetch = globalThis.fetch;
  });

  afterEach(() => {
    globalThis.fetch = originalFetch;
  });

  // --- fetchLabs ---

  describe("fetchLabs", () => {
    it("returns labs from a successful response", async () => {
      const mockLabs = [{ id: "lab_1", name: "Lab 1", problem_statement: "Do stuff" }];
      globalThis.fetch = (async () =>
        new Response(JSON.stringify(mockLabs), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        })) as typeof fetch;

      const labs = await fetchLabs();
      expect(labs).to.deep.equal(mockLabs);
    });

    it("throws on non-ok response", async () => {
      globalThis.fetch = (async () => new Response("not found", { status: 404 })) as typeof fetch;

      try {
        await fetchLabs();
        expect.fail("should have thrown");
      } catch (err: unknown) {
        expect((err as Error).message).to.equal("Failed to fetch labs");
      }
    });
  });

  // --- fetchLab ---

  describe("fetchLab", () => {
    it("returns a single lab on success", async () => {
      const mockLab = { id: "lab_1", name: "Lab 1", problem_statement: "Do stuff" };
      globalThis.fetch = (async () =>
        new Response(JSON.stringify(mockLab), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        })) as typeof fetch;

      const lab = await fetchLab("lab_1");
      expect(lab.id).to.equal("lab_1");
      expect(lab.name).to.equal("Lab 1");
    });

    it("throws with server error message on 404", async () => {
      globalThis.fetch = (async () =>
        new Response(JSON.stringify({ error: 'Lab "lab_99" not found' }), {
          status: 404,
          headers: { "Content-Type": "application/json" },
        })) as typeof fetch;

      try {
        await fetchLab("lab_99");
        expect.fail("should have thrown");
      } catch (err: unknown) {
        expect((err as Error).message).to.equal('Lab "lab_99" not found');
      }
    });

    it("throws generic message when server returns no JSON", async () => {
      globalThis.fetch = (async () =>
        new Response("internal error", { status: 500 })) as typeof fetch;

      try {
        await fetchLab("lab_1");
        expect.fail("should have thrown");
      } catch (err: unknown) {
        expect((err as Error).message).to.include("not found");
      }
    });

    it("calls the correct URL", async () => {
      let capturedUrl = "";
      globalThis.fetch = (async (input: RequestInfo | URL) => {
        capturedUrl = String(input);
        return new Response(JSON.stringify({ id: "lab_2", name: "Lab 2", problem_statement: "" }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }) as typeof fetch;

      await fetchLab("lab_2");
      expect(capturedUrl).to.equal("/api/labs/lab_2");
    });
  });

  // --- submitSolution ---

  describe("submitSolution", () => {
    it("sends code as form field", async () => {
      let capturedBody: FormData | undefined;
      globalThis.fetch = (async (_input: RequestInfo | URL, init?: RequestInit) => {
        capturedBody = init?.body as FormData;
        return new Response(JSON.stringify({ evaluations: [], total_points: 0 }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }) as typeof fetch;

      await submitSolution({ labId: "lab_1", code: 'print("hi")' });
      expect(capturedBody!.get("lab_id")).to.equal("lab_1");
      expect(capturedBody!.get("code")).to.equal('print("hi")');
      expect(capturedBody!.has("file")).to.be.false;
    });

    it("sends code with custom filename", async () => {
      let capturedBody: FormData | undefined;
      globalThis.fetch = (async (_input: RequestInfo | URL, init?: RequestInit) => {
        capturedBody = init?.body as FormData;
        return new Response(JSON.stringify({ evaluations: [], total_points: 0 }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }) as typeof fetch;

      await submitSolution({ labId: "lab_1", code: "x=1", filename: "my.py" });
      expect(capturedBody!.get("filename")).to.equal("my.py");
    });

    it("sends file when file param provided", async () => {
      let capturedBody: FormData | undefined;
      globalThis.fetch = (async (_input: RequestInfo | URL, init?: RequestInit) => {
        capturedBody = init?.body as FormData;
        return new Response(JSON.stringify({ evaluations: [], total_points: 0 }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }) as typeof fetch;

      const file = new File(["code"], "test.py");
      await submitSolution({ labId: "lab_1", file });
      expect(capturedBody!.has("file")).to.be.true;
      expect(capturedBody!.has("code")).to.be.false;
    });
  });

  // --- submitFile (legacy) ---

  describe("submitFile", () => {
    it("returns grade result on success", async () => {
      const mockResult = {
        evaluations: [{ type: "stdout", name: "hw", actual: "Hello", status: "OK", points: 1 }],
        total_points: 1,
      };
      globalThis.fetch = (async () =>
        new Response(JSON.stringify(mockResult), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        })) as typeof fetch;

      const file = new File(["print('hello')"], "hello.py");
      const result = await submitFile("lab_1", file);
      expect(result.total_points).to.equal(1);
      expect(result.evaluations).to.have.lengthOf(1);
    });

    it("throws with error message on failure", async () => {
      globalThis.fetch = (async () =>
        new Response(JSON.stringify({ error: "Invalid lab ID" }), {
          status: 400,
          headers: { "Content-Type": "application/json" },
        })) as typeof fetch;

      const file = new File(["x"], "test.py");
      try {
        await submitFile("bad_lab", file);
        expect.fail("should have thrown");
      } catch (err: unknown) {
        expect((err as Error).message).to.equal("Invalid lab ID");
      }
    });

    it("throws generic message when server returns no error field", async () => {
      globalThis.fetch = (async () =>
        new Response(JSON.stringify({}), {
          status: 500,
          headers: { "Content-Type": "application/json" },
        })) as typeof fetch;

      const file = new File(["x"], "test.py");
      try {
        await submitFile("lab_1", file);
        expect.fail("should have thrown");
      } catch (err: unknown) {
        expect((err as Error).message).to.equal("Submission failed");
      }
    });

    it("sends POST to /api/submit with form data", async () => {
      let capturedUrl = "";
      let capturedMethod = "";
      let capturedBody: FormData | undefined;

      globalThis.fetch = (async (input: RequestInfo | URL, init?: RequestInit) => {
        capturedUrl = String(input);
        capturedMethod = init?.method ?? "GET";
        capturedBody = init?.body as FormData;
        return new Response(JSON.stringify({ evaluations: [], total_points: 0 }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }) as typeof fetch;

      const file = new File(["code"], "script.py");
      await submitFile("lab_2", file);

      expect(capturedUrl).to.equal("/api/submit");
      expect(capturedMethod).to.equal("POST");
      expect(capturedBody).to.not.be.undefined;
      expect(capturedBody!.get("lab_id")).to.equal("lab_2");
      expect(capturedBody!.has("file")).to.be.true;
    });
  });
});
