import { describe, it, beforeEach, afterEach } from "@std/testing/bdd";
import { expect } from "chai";
import { setupDOM, teardownDOM } from "./test-setup.ts";
import {
  buildLabOptions,
  formatProblemStatement,
  buildEvalCard,
  buildResultsFragment,
  validateFile,
} from "./ui.ts";
import type { Lab, Evaluation, GradeResult } from "./api.ts";

describe("ui", () => {
  beforeEach(() => {
    setupDOM();
  });

  afterEach(() => {
    teardownDOM();
  });

  // --- buildLabOptions ---

  describe("buildLabOptions", () => {
    it("returns a placeholder option when given an empty array", () => {
      const opts = buildLabOptions([]);
      expect(opts).to.have.lengthOf(1);
      expect(opts[0].disabled).to.be.true;
      expect(opts[0].textContent).to.equal("Select a lab…");
    });

    it("creates an option for each lab plus the placeholder", () => {
      const labs: Lab[] = [
        { id: "lab_1", name: "Lab 1", problem_statement: "Do stuff" },
        { id: "lab_2", name: "Lab 2", problem_statement: "Do more stuff" },
      ];
      const opts = buildLabOptions(labs);
      expect(opts).to.have.lengthOf(3);

      expect(opts[1].value).to.equal("lab_1");
      expect(opts[1].textContent).to.equal("Lab 1");
      expect(opts[1].dataset.problem).to.equal("Do stuff");

      expect(opts[2].value).to.equal("lab_2");
      expect(opts[2].textContent).to.equal("Lab 2");
    });

    it("stores the problem_statement in data-problem", () => {
      const labs: Lab[] = [{ id: "x", name: "X", problem_statement: "Write code\\nplease" }];
      const opts = buildLabOptions(labs);
      expect(opts[1].dataset.problem).to.equal("Write code\\nplease");
    });
  });

  // --- formatProblemStatement ---

  describe("formatProblemStatement", () => {
    it("replaces literal \\n with <br>", () => {
      expect(formatProblemStatement("line1\\nline2")).to.equal("line1<br>line2");
    });

    it("handles multiple \\n", () => {
      expect(formatProblemStatement("a\\nb\\nc")).to.equal("a<br>b<br>c");
    });

    it("returns the string unchanged when no \\n", () => {
      expect(formatProblemStatement("no newlines")).to.equal("no newlines");
    });

    it("handles empty string", () => {
      expect(formatProblemStatement("")).to.equal("");
    });
  });

  // --- buildEvalCard ---

  describe("buildEvalCard", () => {
    it("creates a card with pass styling when points > 0", () => {
      const ev: Evaluation = {
        type: "stdout",
        name: "hw",
        actual: "Hello",
        status: "Correct!",
        points: 1.0,
      };
      const card = buildEvalCard(ev);

      expect(card.tagName).to.equal("DIV");
      expect(card.className).to.include("border-l-green-500");
      expect(card.getAttribute("data-testid")).to.equal("eval-hw");
      expect(card.querySelector(".eval-type")!.textContent).to.equal("stdout");
      expect(card.querySelector(".eval-name")!.textContent).to.equal("hw");
      expect(card.querySelector(".eval-points")!.textContent).to.equal("1 pts");
      expect(card.querySelector(".eval-status")!.textContent).to.equal("Correct!");
      expect(card.querySelector(".eval-actual code")!.textContent).to.equal("Hello");
    });

    it("creates a card with fail styling when points = 0", () => {
      const ev: Evaluation = {
        type: "function",
        name: "fn",
        actual: "3",
        status: "Wrong",
        points: 0,
      };
      const card = buildEvalCard(ev);

      expect(card.className).to.include("border-l-red-500");
      expect(card.className).not.to.include("border-l-green-500");
      expect(card.querySelector(".eval-points")!.textContent).to.equal("0 pts");
    });

    it("omits output section when actual is null", () => {
      const ev: Evaluation = {
        type: "function",
        name: "fn",
        actual: null,
        status: "No match",
        points: 0,
      };
      const card = buildEvalCard(ev);

      expect(card.querySelector(".eval-actual")).to.be.null;
    });
  });

  // --- buildResultsFragment ---

  describe("buildResultsFragment", () => {
    it("returns a fragment with eval cards and a total", () => {
      const result: GradeResult = {
        evaluations: [
          { type: "stdout", name: "a", actual: "hi", status: "OK", points: 0.5 },
          { type: "function", name: "b", actual: "42", status: "OK", points: 0.5 },
        ],
        total_points: 1.0,
      };
      const frag = buildResultsFragment(result);
      const container = document.createElement("div");
      container.appendChild(frag);

      const cards = container.querySelectorAll(".eval-card");
      expect(cards).to.have.lengthOf(2);

      const total = container.querySelector(".total");
      expect(total).to.not.be.null;
      expect(total!.textContent).to.equal("Total: 1 points");
    });

    it("handles empty evaluations", () => {
      const result: GradeResult = { evaluations: [], total_points: 0 };
      const frag = buildResultsFragment(result);
      const container = document.createElement("div");
      container.appendChild(frag);

      expect(container.querySelectorAll(".eval-card")).to.have.lengthOf(0);
      expect(container.querySelector(".total")!.textContent).to.equal("Total: 0 points");
    });
  });

  // --- validateFile ---

  describe("validateFile", () => {
    it("returns error when file is undefined", () => {
      expect(validateFile(undefined)).to.equal("Please select a file.");
    });

    it("returns error for disallowed extension", () => {
      const file = new File(["x"], "script.exe");
      expect(validateFile(file)).to.include("Only");
    });

    it("returns error for too-large file", () => {
      const big = new File([new Uint8Array(21 * 1024)], "big.py");
      expect(validateFile(big)).to.include("too large");
    });

    it("returns null for a valid .py file", () => {
      const file = new File(["print('hi')"], "hello.py");
      expect(validateFile(file)).to.be.null;
    });

    it("respects custom extensions list", () => {
      const file = new File(["x"], "code.js");
      expect(validateFile(file, ["js"])).to.be.null;
      expect(validateFile(file, ["py"])).to.not.be.null;
    });

    it("respects custom max size", () => {
      const file = new File(["abc"], "ok.py");
      expect(validateFile(file, ["py"], 2)).to.include("too large");
      expect(validateFile(file, ["py"], 1024)).to.be.null;
    });
  });
});
