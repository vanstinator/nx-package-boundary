package main

import (
	"github.com/vanstinator/nxboundary"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(nxboundary.NewAnalyzer())
}
