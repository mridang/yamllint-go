package main_test

import (
	"fmt"
	"testing"

	"github.com/mridang/yamllint-go/internal/parser"
	"github.com/mridang/yamllint-go/internal/rules/colons"
	"github.com/mridang/yamllint-go/internal/types"
)

func TestFoo(t *testing.T) {
	// 1. Define input YAML with a violation
	// Violation: "key : value" has 1 space before the colon.
	yamlContent := "key : value\n"

	// 2. Parse the content
	tokens, _, _, err := parser.Parse(yamlContent)
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}

	// --- DEBUG PRINT ---
	fmt.Println("--- TOKEN DUMP ---")
	for i, tok := range tokens {
		// Calculate the gap from previous token if possible
		gap := ""
		if i > 0 {
			prev := tokens[i-1]
			if prev.EndMark.Line == tok.StartMark.Line {
				gapVal := tok.StartMark.Column - prev.EndMark.Column
				gap = fmt.Sprintf("(Gap from prev: %d)", gapVal)
			}
		}

		fmt.Printf("[%d] Type: %d | Val: %q | Line: %d | StartCol: %d | EndCol: %d %s\n",
			i, tok.Type, tok.Value, tok.StartMark.Line, tok.StartMark.Column, tok.EndMark.Column, gap)
	}
	fmt.Println("------------------")
	// -------------------

	// 3. Initialize the Rule
	rule := colons.Rule{}
	config := rule.DefaultConfig() // Default: max-spaces-before: 0

	// 4. Create the shared context
	context := make(map[string]any)

	// 5. Run the Check loop
	fmt.Println("--- Starting Check ---")

	var problems []types.Problem

	for i, tok := range tokens {
		var prev, next, nextNext *types.Token

		if i > 0 {
			prev = tokens[i-1]
		}
		if i < len(tokens)-1 {
			next = tokens[i+1]
		}
		if i < len(tokens)-2 {
			nextNext = tokens[i+2]
		}

		// Call the rule logic
		for problem := range rule.Check(config, tok, prev, next, nextNext, context) {
			problems = append(problems, problem)
			fmt.Printf("FOUND PROBLEM: [Line %d, Col %d] %s\n",
				problem.Line, problem.Column, problem.Desc)
		}
	}

	// 6. Assert results
	if len(problems) == 0 {
		t.Error("Expected to find a problem, but found none!")
	} else {
		p := problems[0]
		if p.Desc != "too many spaces before colon" {
			t.Errorf("Unexpected problem description: %s", p.Desc)
		}
	}

	fmt.Println("--- End Check ---")
}
