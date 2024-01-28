package nxboundary

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func makeAnalyzer() *analysis.Analyzer {
	cnf := Config{
		DepConstraints: make(stringMap),
	}

	return &analysis.Analyzer{
		Name:  "nxboundary",
		Doc:   "enforce package boundaries for nx monorepos",
		Flags: flags(&cnf),
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return runWithConfig(cnf, pass)
		},
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// func TestIncorrectFlags(t *testing.T) {
// 	assertWrongAllowedTagsError := func(msg string, err error) {
// 		if err == nil || err.Error() != errWrongAllowedTags.Error() {
// 			t.Errorf("Wrong error for invalid usage[%q]: %v", msg, err)
// 		}
// 	}
// 	a := makeAnalyzer()
// 	flg := a.Flags.Lookup(FlagAllowedTags)
// 	assertWrongAllowedTagsError("empty flag", flg.Value.Set(""))
// 	assertWrongAllowedTagsError("white space only", flg.Value.Set("   "))
// 	assertWrongAllowedTagsError("no colons", flg.Value.Set("no colons"))
// 	assertWrongAllowedTagsError("no tags", flg.Value.Set("scope:test|"))
// }

// func TestCorrectFlags(t *testing.T) {
// 	a := makeAnalyzer()
// 	flg := a.Flags.Lookup(FlagAllowedTags)
// 	if err := flg.Value.Set("scope:test|tag1,tag2"); err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}
// 	if err := flg.Value.Set("scope:test|scope:test2"); err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}
// 	if err := flg.Value.Set("scope:test|scope:test2,scope:test3"); err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}
// }

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	testCases := []struct {
		desc        string
		pkg         string
		allowedTags stringMap
	}{
		{
			desc: "Package A",
			pkg:  "a",
			allowedTags: stringMap{
				"scope:a": {"scope:c"},
			},
		},
		{
			desc: "Package B",
			pkg:  "b",
			allowedTags: stringMap{
				"scope:b": {"scope:c"},
			},
		},
		{
			desc: "Package C",
			pkg:  "c",
			allowedTags: stringMap{
				"scope:c": {"scope:d"},
			},
		},
		{
			desc: "Package D",
			pkg:  "d",
			allowedTags: stringMap{
				"scope:d": {"scope:b"},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			// t.Parallel()

			a := makeAnalyzer()
			flg := a.Flags.Lookup(FlagAllowedTags)
			for k, v := range test.allowedTags {
				err := flg.Value.Set(fmt.Sprintf("%s|%s", k, strings.Join(v, ",")))
				if err != nil {
					t.Fatal(err)
				}
			}

			analysistest.RunWithSuggestedFixes(t, testdata, a, "./src/"+test.pkg+"/...")
		})
	}
}
