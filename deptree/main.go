package main

import (
	"flag"
	"fmt"
	"go/token"
	"log"

	"github.com/dedis/deter-checker/checker"
	pkgs "golang.org/x/tools/go/packages"
)

var stdPkgs = make(map[string]struct{})
var exists = struct{}{}

const mode pkgs.LoadMode = pkgs.NeedName |
	pkgs.NeedFiles |
	pkgs.NeedImports

func init() {
	packages, err := pkgs.Load(nil, "std")
	if err != nil {
		panic(err)
	}

	for _, p := range packages {
		stdPkgs[p.PkgPath] = exists
	}
}

func checkFiles(deps map[string]*pkgs.Package, wl, bl map[string]bool) {
	for name, pkg := range deps {
		fmt.Println("Checking source files of:", name)
		files := pkg.GoFiles
		for _, f := range files {
			fmt.Println("====", f, "====")
			checker.AnalyzeSource(&f, wl, bl)
		}
	}
}

func findDependencies(pkgList []*pkgs.Package) map[string]*pkgs.Package {
	pkgSet := make(map[string]*pkgs.Package)
	for _, pkg := range pkgList {
		imports := pkg.Imports
		for _, imp := range imports {
			impStr := imp.String()
			if _, ok := stdPkgs[impStr]; !ok {
				if _, ok := pkgSet[impStr]; !ok {
					pkgSet[impStr] = imp
				}
			}
		}
	}
	return pkgSet
}

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintln(out, "usage: analyzer [options] <module dir>")
		fmt.Fprintln(out, "Options:")
		flag.PrintDefaults()
	}

	pattern := flag.String("pattern", "./...", "Go package")
	flag.Parse()

	wfile := flag.Args()[1]
	bfile := flag.Args()[2]

	fset := token.NewFileSet()
	cfg := &pkgs.Config{Mode: mode, Dir: flag.Args()[0], Fset: fset}
	pkgList, err := pkgs.Load(cfg, *pattern)
	if err != nil {
		log.Fatal(err)
	}

	wl := checker.ReadList(&wfile)
	bl := checker.ReadList(&bfile)

	deps := findDependencies(pkgList)
	checkFiles(deps, wl, bl)
}
