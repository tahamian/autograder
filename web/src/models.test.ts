import { describe, it } from "@std/testing/bdd";
import { expect } from "chai";
import type { Lab, Evaluation, GradeResult, ApiError } from "./models.ts";

describe("models", () => {
  describe("Lab", () => {
    it("has the correct shape", () => {
      const lab: Lab = { id: "lab_1", name: "Test", problem_statement: "Do stuff" };
      expect(lab.id).to.equal("lab_1");
      expect(lab.name).to.equal("Test");
      expect(lab.problem_statement).to.equal("Do stuff");
    });
  });

  describe("Evaluation", () => {
    it("has the correct shape for a pass", () => {
      const ev: Evaluation = {
        type: "stdout",
        name: "hw",
        actual: "Hello",
        status: "OK",
        points: 1,
      };
      expect(ev.points).to.equal(1);
      expect(ev.actual).to.equal("Hello");
    });

    it("supports null actual", () => {
      const ev: Evaluation = {
        type: "function",
        name: "fn",
        actual: null,
        status: "No match",
        points: 0,
      };
      expect(ev.actual).to.be.null;
    });
  });

  describe("GradeResult", () => {
    it("has evaluations and total_points", () => {
      const result: GradeResult = {
        evaluations: [
          { type: "stdout", name: "a", actual: "hi", status: "OK", points: 0.5 },
          { type: "function", name: "b", actual: "42", status: "OK", points: 0.5 },
        ],
        total_points: 1.0,
      };
      expect(result.evaluations).to.have.lengthOf(2);
      expect(result.total_points).to.equal(1.0);
    });

    it("handles empty evaluations", () => {
      const result: GradeResult = { evaluations: [], total_points: 0 };
      expect(result.evaluations).to.be.empty;
    });
  });

  describe("ApiError", () => {
    it("has an error field", () => {
      const err: ApiError = { error: "something went wrong" };
      expect(err.error).to.equal("something went wrong");
    });
  });
});
