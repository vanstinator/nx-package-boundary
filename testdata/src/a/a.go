package a

import (
	"github.com/vanstinator/nxboundary/testdata/src/b" // want `package github.com/vanstinator/nxboundary/testdata/src/a is not allowed to import package github.com/vanstinator/nxboundary/testdata/src/b`
	"github.com/vanstinator/nxboundary/testdata/src/c"
)

func PackageAFunc() {
	b.PackageBFunc()
	c.PackageCFunc()
}
