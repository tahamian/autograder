package grader

import (
	"autograder/internal/config"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

var g = &DefaultGrader{}

func TestBuildInput_StdoutOnly(t *testing.T) {
	lab := &config.Lab{Testcase: []config.Testcase{{Type: "stdout"}}}
	input := g.BuildInput(lab)
	if !input.Stdout {
		t.Error("expected Stdout true")
	}
	if len(input.Functions) != 0 {
		t.Error("expected no functions")
	}
}

func TestBuildInput_FunctionOnly(t *testing.T) {
	lab := &config.Lab{
		Testcase: []config.Testcase{{
			Type: "function", Name: "tc1",
			Functions: []config.Function{{FunctionName: "add"}},
		}},
	}
	input := g.BuildInput(lab)
	if input.Stdout {
		t.Error("expected Stdout false")
	}
	if len(input.Functions) != 1 || input.Functions[0].TestcaseName != "tc1" {
		t.Errorf("unexpected functions: %+v", input.Functions)
	}
}

func TestBuildInput_Mixed(t *testing.T) {
	lab := &config.Lab{
		Testcase: []config.Testcase{
			{Type: "stdout"},
			{Type: "function", Name: "f", Functions: []config.Function{{FunctionName: "a"}, {FunctionName: "b"}}},
		},
	}
	input := g.BuildInput(lab)
	if !input.Stdout || len(input.Functions) != 2 {
		t.Errorf("unexpected: stdout=%v functions=%d", input.Stdout, len(input.Functions))
	}
}

func TestEvaluate_StdoutCorrect(t *testing.T) {
	lab := &config.Lab{
		Testcase: []config.Testcase{{
			Type: "stdout", Name: "hw",
			Expected: []config.Expected{{Feedback: "OK", Points: 1, Values: []string{"Hello"}}},
		}},
	}
	result, _ := g.Evaluate(lab, &MarkerOutput{Stdout: "Hello"})
	if result.TotalPoints != 1 || result.Evaluations[0].Status != "OK" {
		t.Errorf("unexpected: %+v", result)
	}
}

func TestEvaluate_StdoutIncorrect(t *testing.T) {
	lab := &config.Lab{
		Testcase: []config.Testcase{{
			Type: "stdout", Name: "hw",
			Expected: []config.Expected{{Feedback: "OK", Points: 1, Values: []string{"Hello"}}},
		}},
	}
	result, _ := g.Evaluate(lab, &MarkerOutput{Stdout: "Wrong"})
	if result.TotalPoints != 0 || result.Evaluations[0].Status != "Incorrect output" {
		t.Errorf("unexpected: %+v", result)
	}
}

func TestEvaluate_FunctionCorrect(t *testing.T) {
	lab := &config.Lab{
		Testcase: []config.Testcase{{
			Type: "function", Name: "tc",
			Expected: []config.Expected{{Feedback: "Nice", Points: 1, Values: []string{"5.0"}}},
		}},
	}
	output := &MarkerOutput{
		Functions: []FunctionOutput{
			{Result: "5.0", TestcaseName: "tc"},
		},
	}
	result, _ := g.Evaluate(lab, output)
	if result.TotalPoints != 1 || result.Evaluations[0].Status != "Nice" {
		t.Errorf("unexpected: %+v", result)
	}
}

func TestEvaluate_FunctionNoMatch(t *testing.T) {
	lab := &config.Lab{
		Testcase: []config.Testcase{{
			Type: "function", Name: "tc",
			Expected: []config.Expected{{Feedback: "OK", Points: 1, Values: []string{"5"}}},
		}},
	}
	output := &MarkerOutput{
		Functions: []FunctionOutput{
			{Result: "5", TestcaseName: "other"},
		},
	}
	result, _ := g.Evaluate(lab, output)
	if result.Evaluations[0].Status != "Could not match function to test case" {
		t.Errorf("unexpected status: %q", result.Evaluations[0].Status)
	}
}

func TestEvaluate_Multiple(t *testing.T) {
	lab := &config.Lab{
		Testcase: []config.Testcase{
			{Type: "stdout", Name: "s", Expected: []config.Expected{{Feedback: "OK", Points: 0.5, Values: []string{"hi"}}}},
			{Type: "function", Name: "f", Expected: []config.Expected{{Feedback: "OK", Points: 0.5, Values: []string{"42"}}}},
		},
	}
	output := &MarkerOutput{
		Stdout: "hi",
		Functions: []FunctionOutput{
			{Result: "42", TestcaseName: "f"},
		},
	}
	result, _ := g.Evaluate(lab, output)
	if result.TotalPoints != 1.0 || len(result.Evaluations) != 2 {
		t.Errorf("unexpected: total=%f evals=%d", result.TotalPoints, len(result.Evaluations))
	}
}

// #3 — Numeric comparison: float64(5) from JSON should match expected "5.0"
func TestEvaluate_NumericFloatMatchesString(t *testing.T) {
	lab := &config.Lab{
		Testcase: []config.Testcase{{
			Type: "function", Name: "tc",
			Expected: []config.Expected{{Feedback: "OK", Points: 1, Values: []string{"5.0"}}},
		}},
	}
	// JSON unmarshals numbers as float64 — "5.0" in Python becomes float64(5) in Go
	output := &MarkerOutput{
		Functions: []FunctionOutput{
			{Result: float64(5), TestcaseName: "tc"},
		},
	}
	result, _ := g.Evaluate(lab, output)
	if result.TotalPoints != 1 {
		t.Errorf("expected 1 point for float64(5) matching '5.0', got %f", result.TotalPoints)
	}
}

func TestEvaluate_NumericIntMatchesString(t *testing.T) {
	lab := &config.Lab{
		Testcase: []config.Testcase{{
			Type: "function", Name: "tc",
			Expected: []config.Expected{{Feedback: "OK", Points: 1, Values: []string{"42"}}},
		}},
	}
	output := &MarkerOutput{
		Functions: []FunctionOutput{
			{Result: float64(42), TestcaseName: "tc"},
		},
	}
	result, _ := g.Evaluate(lab, output)
	if result.TotalPoints != 1 {
		t.Errorf("expected 1 point for float64(42) matching '42', got %f", result.TotalPoints)
	}
}

func TestFunctionOutput_Serialization(t *testing.T) {
	fo := FunctionOutput{
		Result: "5.0", Status: 0, Buffer: "hello", TestcaseName: "tc1",
	}
	data, err := json.Marshal(fo)
	if err != nil {
		t.Fatal(err)
	}
	var loaded FunctionOutput
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatal(err)
	}
	if loaded.Result != "5.0" || loaded.Buffer != "hello" || loaded.TestcaseName != "tc1" {
		t.Errorf("unexpected: %+v", loaded)
	}
}

func TestReadOutput(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "output.json")
	os.WriteFile(p, []byte(`{"stdout":"ok","functions":[]}`), 0644)
	out, err := g.ReadOutput(p)
	if err != nil {
		t.Fatal(err)
	}
	if out.Stdout != "ok" {
		t.Errorf("expected 'ok', got %q", out.Stdout)
	}
}

func TestReadOutput_MissingFile(t *testing.T) {
	if _, err := g.ReadOutput("/nope"); err == nil {
		t.Error("expected error")
	}
}

func TestReadOutput_BadJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	os.WriteFile(p, []byte("nope"), 0644)
	if _, err := g.ReadOutput(p); err == nil {
		t.Error("expected error")
	}
}

func TestWriteInput(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "input.json")
	input := &MarkerInput{Filename: "test.py", Stdout: true}
	if err := WriteInput(input, p); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(p)
	var loaded MarkerInput
	json.Unmarshal(data, &loaded)
	if loaded.Filename != "test.py" || !loaded.Stdout {
		t.Errorf("unexpected: %+v", loaded)
	}
}
