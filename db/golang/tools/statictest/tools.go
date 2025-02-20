//go:build tools
// +build tools

package tools

import (
	_ "github.com/jingyugao/rowserrcheck/passes/rowserr"
	_ "github.com/timakin/bodyclose/passes/bodyclose"
	_ "golang.org/x/tools/go/analysis"
	_ "honnef.co/go/tools/simple"
	_ "honnef.co/go/tools/staticcheck"
	_ "honnef.co/go/tools/stylecheck"
)
