package nxboundary

import (
	"encoding/json"
	"go/ast"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

type NxProjectFile struct {
	Tags []string `json:"tags"`
}

var (
	nxPackageTags = make(map[string][]string)
	config        = &Config{
		DepConstraints: make(map[string][]string),
	}
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "nxboundary",
		Doc:  "enforce package boundaries for nx monorepos",
		Run:  run,

		Flags: flags(config),

		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	return runWithConfig(*config, pass)
}

func runWithConfig(config Config, pass *analysis.Pass) (interface{}, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	gomod, err := os.ReadFile(cwd + "/go.mod")
	if err != nil {
		return nil, err
	}

	modPath := modfile.ModulePath(gomod) + "/"

	for _, file := range pass.Files {
		packageTags := fetchTags(modPath, pass.Pkg.Path())

		ast.Inspect(file, func(node ast.Node) bool {
			importSpec, ok := node.(*ast.ImportSpec)
			if !ok {
				return true
			}

			importPath := strings.ReplaceAll(importSpec.Path.Value, "\"", "")

			importTags := fetchTags(modPath, importPath)

			// If the current package has tags, the import has tags, the import is in the scope of the current module, and the import is not part of the current package, check if there's any overlap between them.
			if (len(packageTags) > 0 || len(importTags) > 0) && strings.Contains(importPath, modPath) && !strings.Contains(importPath, pass.Pkg.Path()) {
				overlap := false

				for _, packageTag := range packageTags {
					for _, importTag := range importTags {
						if config.IsTagAllowed(packageTag, importTag) {
							overlap = true
							break
						}
					}
					if overlap {
						break
					}
				}

				if !overlap {
					pass.Reportf(importSpec.Pos(), "package %s is not allowed to import package %s", pass.Pkg.Path(), importPath)
				}
			}

			return true
		})
	}

	return nil, nil
}

func fetchTags(modulePath string, filePath string) []string {
	// try loading project.json from package path
	// if not found, try loading from parent directory

	if !strings.Contains(filePath, modulePath) {
		return []string{}
	}

	tags := nxPackageTags[filePath]
	if len(tags) > 0 {
		return tags
	}

	cwd, err := os.Getwd()
	if err != nil {
		return []string{}
	}

	path := filePath
	if modulePath != "" {
		path = filepath.Join(cwd, strings.Replace(filePath, modulePath, "", 1))
	}

	nxProjectFile, err := os.ReadFile(path + "/project.json")
	// If we can't find a project.json file, recurse up the directory tree until we find a project.json file
	if err != nil {
		// If we're at the root of the project, return an empty array
		if path == cwd {
			return []string{}
		}

		parentDir := strings.Split(filePath, "/")
		parentDir = parentDir[:len(parentDir)-1]

		return fetchTags(modulePath, strings.Join(parentDir, "/"))
	}

	var projectFile NxProjectFile
	err = json.Unmarshal(nxProjectFile, &projectFile)
	if err != nil {
		return []string{}
	}

	nxPackageTags[filePath] = projectFile.Tags

	return projectFile.Tags
}
