package grader

import (
	"autograder/internal/config"
	"autograder/internal/models"
	"encoding/json"
	"fmt"
	"os"
)

// MarkerInput is the JSON payload sent to the marker container.
// Uses generated FunctionT from the FlatBuffers schema.
type MarkerInput struct {
	Filename  string            `json:"filename"`
	Stdout    bool              `json:"stdout"`
	Functions []config.Function `json:"functions"`
}

// MarkerOutput is the JSON result from the marker container.
type MarkerOutput struct {
	Stdout    string `json:"stdout"`
	Functions []struct {
		Result       interface{} `json:"result"`
		Status       int         `json:"status"`
		Buffer       string      `json:"buffer"`
		TestcaseName string      `json:"testcase_name"`
	} `json:"functions"`
}

// Grader evaluates marker output against lab test cases.
// Returns generated FlatBuffers model types.
type Grader interface {
	BuildInput(lab *config.Lab) *MarkerInput
	ReadOutput(path string) (*MarkerOutput, error)
	Evaluate(lab *config.Lab, output *MarkerOutput) (*models.GradeResultT, error)
}

// DefaultGrader is the standard implementation.
type DefaultGrader struct{}

// BuildInput converts a lab's test cases into a MarkerInput.
func (g *DefaultGrader) BuildInput(lab *config.Lab) *MarkerInput {
	input := &MarkerInput{}
	for _, tc := range lab.Testcase {
		if tc.Type == "stdout" {
			input.Stdout = true
		}
		if tc.Type == "function" {
			for i := range tc.Functions {
				tc.Functions[i].TestcaseName = tc.Name
			}
			input.Functions = append(input.Functions, tc.Functions...)
		}
	}
	return input
}

// ReadOutput loads marker output from a JSON file.
func (g *DefaultGrader) ReadOutput(path string) (*MarkerOutput, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading output %s: %w", path, err)
	}
	o := &MarkerOutput{}
	if err := json.Unmarshal(data, o); err != nil {
		return nil, fmt.Errorf("parsing output JSON: %w", err)
	}
	return o, nil
}

// Evaluate grades marker output against a lab's expected results.
// Returns a generated GradeResultT.
func (g *DefaultGrader) Evaluate(lab *config.Lab, output *MarkerOutput) (*models.GradeResultT, error) {
	result := &models.GradeResultT{}

	for _, tc := range lab.Testcase {
		var eval *models.EvaluationT

		switch tc.Type {
		case "stdout":
			eval = evaluateStdout(&tc, output)
		case "function":
			eval = evaluateFunction(&tc, output)
		default:
			continue
		}

		result.Evaluations = append(result.Evaluations, eval)
		result.TotalPoints += eval.Points
	}

	return result, nil
}

func evaluateStdout(tc *config.Testcase, output *MarkerOutput) *models.EvaluationT {
	eval := &models.EvaluationT{
		Actual: output.Stdout,
		Type:   tc.Type,
		Name:   tc.Name,
	}
	for _, exp := range tc.Expected {
		for _, v := range exp.Values {
			if v == output.Stdout {
				eval.Points = exp.Points
				eval.Status = exp.Feedback
				return eval
			}
		}
	}
	eval.Status = "Incorrect output"
	return eval
}

func evaluateFunction(tc *config.Testcase, output *MarkerOutput) *models.EvaluationT {
	eval := &models.EvaluationT{Type: tc.Type, Name: tc.Name}

	for _, fn := range output.Functions {
		if tc.Name == fn.TestcaseName {
			eval.Actual = fmt.Sprintf("%v", fn.Result)
			for _, exp := range tc.Expected {
				for _, v := range exp.Values {
					if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", fn.Result) {
						eval.Points = exp.Points
						eval.Status = exp.Feedback
						return eval
					}
				}
			}
		}
	}

	if eval.Actual == "" {
		eval.Status = "Could not match function to test case"
	} else if eval.Status == "" {
		eval.Status = "Incorrect return value"
	}
	return eval
}

// WriteInput writes a MarkerInput to a JSON file.
func WriteInput(input *MarkerInput, path string) error {
	data, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshaling input: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}
