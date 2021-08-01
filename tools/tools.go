// Package tools provides tools for development.
//
// It follows the pattern set-forth in the wiki:
//  - https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
//  - https://github.com/go-modules-by-example/index/tree/4ea90b07f9/010_tools
//
// Copyright © 2021 Trevor N. Suarez (Rican7)
//
// +build tools
package tools

import (
	// Tools for development
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/cmd/goimports"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
