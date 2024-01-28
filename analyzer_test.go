package nxboundary

import (
	"fmt"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func makeAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "nxboundary",
		Doc:  "enforce package boundaries for nx monorepos",
		// Flags: flags(&cnf),
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return run(pass)
		},
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	fmt.Println("testdata", testdata)

	testCases := []struct {
		desc string
		pkg  string
	}{
		{
			desc: "Invalid imports",
			pkg:  "a",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			a := makeAnalyzer()

			analysistest.Run(t, testdata, a, "./src/"+test.pkg+"/...")
		})
	}
}
