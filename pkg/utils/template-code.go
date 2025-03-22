package utils

import (
	"fmt"
	"regexp"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/models"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
)

func extractFunctionName(codeSnippet string) string {
	patterns := []string{
		`def\s+(\w+)\s*\(`,       // Python function: def func_name(a, b):
		`const\s+(\w+)\s*=\s*\(`, // JavaScript function: const func_name = (a, b) => {}
		`function\s+(\w+)\s*\(`,  // JavaScript function: function func_name(a, b) {}
	}

	for _, pattern := range patterns {
		rg := regexp.MustCompile(pattern)
		matches := rg.FindStringSubmatch(codeSnippet)

		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

func GenerateCodeTemplate(testCases []models.TestCase, language, codeSnippet, userCode string) string {
	// get function name from the code snippet
	functionName := extractFunctionName(codeSnippet)

	// create string format for testcases
	testcases := "["
	for _, tc := range testCases {
		testCase := ""
		if language == "python" {
			testCase = "{\"input\": \"" + tc.Input + "\", \"output\": \"" + tc.Output + "\"},"
		}
		if language == "javascript" {
			testCase = "{input: " + tc.Input + ", output: " + tc.Output + "},"
		}

		testcases += testCase
	}
	testcases += "]"

	codeTemplate := fmt.Sprintf(constants.GOLANG_CODE_TEMPLATE[language], userCode, testcases, functionName)

	return codeTemplate
}