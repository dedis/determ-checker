package main

import (
	"flag"
	"fmt"
	"go/token"
	"log"

	"golang.org/x/tools/go/packages"
)

var wlFile string = "../inputs/whitelist-pkg.txt"
var blTypes string = "../inputs/blacklist-types.txt"

const mode packages.LoadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedImports

//func findDependencies(pkg *packages.Package) {
//fmt.Println("Printing dependencies of:", pkg.Name)
//if len(pkg.Imports) == 0 {
//return
//}
//for _, p := range pkg.Imports {
//fmt.Println(">>>", p.Name)
//findDependencies(p)
//}
//return
//}

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintln(out, "usage: analyzer [options] <module dir>")
		fmt.Fprintln(out, "Options:")
		flag.PrintDefaults()
	}

	pattern := flag.String("pattern", "./...", "Go package")
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("Expecting a single argument: directory of module")
	}

	fmt.Println(*pattern)

	var fset = token.NewFileSet()
	cfg := &packages.Config{Mode: mode, Dir: flag.Args()[0], Fset: fset}
	pkgs, err := packages.Load(cfg, *pattern)
	if err != nil {
		log.Fatal(err)
	}

	//blPkg := checker.ReadList(&wlFile)
	//blTypes := checker.ReadList(&blTypes)

	for _, pkg := range pkgs {
		fmt.Println(pkg.ID, pkg.PkgPath)
		fmt.Println(pkg.Imports)
		//fmt.Println(pkg.Imports["go.dedis.ch/kyber/v3"].GoFiles)
		imports := pkg.Imports["go.dedis.ch/kyber/v3/proof"]
		files := imports.GoFiles
		fmt.Println("Other files:", imports.OtherFiles)
		for _, f := range files {
			fmt.Println("Source file:", f)
			//checker.AnalyzeSource(&f, blPkg, blTypes)
		}
		for _, p := range imports.Imports {
			fmt.Println("Import dependency:", p.Name)
		}
	}
}
