package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"reflect"
	"strings"
)

func getBlacklist(bpath *string) (bl map[string]bool) {
	file, err := os.Open(*bpath)
	if err != nil {
		log.Fatal(err)
		return bl
	}
	defer file.Close()

	bl = make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bl[scanner.Text()] = true
	}
	return bl
}

func analyzeSource(spath *string, blistPkg map[string]bool, blistTypes map[string]bool) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, *spath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("--- Imports ---")
	for _, imp := range node.Imports {
		lib := strings.Replace(imp.Path.Value, "\"", "", -1)
		if val, _ := blistPkg[lib]; val {
			fmt.Println(lib, "is in package blacklist")
		} else {
			fmt.Println(lib, "is a safe package ")
		}
	}

	ast.Inspect(node, func(n ast.Node) bool {

		switch n := n.(type) {

		case *ast.CompositeLit:
		case *ast.BasicLit:
			fmt.Println(n)
			if exists, _ := blistTypes[n.Kind.String()]; exists {
				fmt.Println(n.Kind.String(), "is in types blacklist")
			} else {
				fmt.Println(n.Kind.String(), "is a safe type")
			}
		default:
			if n != nil {
				val := reflect.ValueOf(n).Elem()
				valType := val.Type().Name()
				if exists, _ := blistTypes[valType]; exists {
					fmt.Println(valType, "is in types blacklist")
				} else {
					fmt.Println(valType, "is a safe type")
				}

				//fmt.Println("----" + val.Type().Name())
			}


		}
		return true
	})

	// info.Uses allows to lookup import paths for identifiers.
	//info := &types.Info{
	//	Uses: make(map[*ast.Ident]types.Object),
	//}

	//ast.Walk(astTypes{info}, node)


	//fmt.Println("--- Function declarations ---")
	//for _, decl := range node.Decls {
	//fnDecl, ok := decl.(*ast.FuncDecl)
	//if ok {
	//fmt.Println(fnDecl.Name.Name)
	//}
	//}
	//ast.Print(fset, node)

	//fmt.Println("--- Inspect ---")
	//ast.Inspect(node, func(n ast.Node) bool {
	//switch x := n.(type) {
	//case *ast.Ident:
	//fmt.Println("Identifier:", x.Name)
	//case *ast.BasicLit:
	//fmt.Println("Literal:", x.Kind.String(), x.ValuePos, x.Value)
	//}
	//return true
	//})
}

type astTypes struct {
	info *types.Info
}


func (v astTypes) Visit(node ast.Node) (w ast.Visitor) {
	switch node := node.(type) {

	case *ast.CallExpr:
		// Get some kind of *ast.Ident for the CallExpr that represents the
		// package. Then we can look it up in v.info. Where exactly it sits in
		// the ast depends on the form of the function call.

		switch node := node.Fun.(type) {
		case *ast.SelectorExpr: // foo.ReadFile
			pkgID := node.X.(*ast.Ident)
			fmt.Println(v.info.Uses[pkgID].(*types.PkgName).Imported().Path())

		case *ast.Ident: // ReadFile
			pkgID := node
			fmt.Println(v.info.Uses[pkgID].Pkg().Path())

		}
	case *ast.CompositeLit:
		fmt.Println(node.Type)

	}

	return v
}


func main() {
	spath := flag.String("s", "", "source file")
	bppath := flag.String("bp", "", "blacklist packages file")
	btpath := flag.String("bt", "", "blacklist types file")
	flag.Parse()
	blistPkg := getBlacklist(bppath)
	blistTypes := getBlacklist(btpath)
	if blistPkg == nil {
		os.Exit(1)
	}
	analyzeSource(spath, blistPkg, blistTypes)
}
