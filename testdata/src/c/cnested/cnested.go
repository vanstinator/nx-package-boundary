package cnested

import "github.com/vanstinator/nxboundary/testdata/src/a" // want `package github.com/vanstinator/nxboundary/testdata/src/c/cnested is not allowed to import package github.com/vanstinator/nxboundary/testdata/src/a`

func CNestedFunc() {
	a.PackageAFunc()
}
