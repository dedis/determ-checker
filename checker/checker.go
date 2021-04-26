package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
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

func analyzeSource(spath *string, blist map[string]bool) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, *spath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("--- Imports ---")
	for _, imp := range node.Imports {
		lib := strings.Replace(imp.Path.Value, "\"", "", -1)
		if val, _ := blist[lib]; val {
			fmt.Println(lib, "is in blacklist")
		} else {
			fmt.Println(lib, "is safe")
		}
	}

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

func main() {
	spath := flag.String("s", "", "source file")
	bpath := flag.String("b", "", "blacklist file")
	flag.Parse()
	blist := getBlacklist(bpath)
	if blist == nil {
		os.Exit(1)
	}
	analyzeSource(spath, blist)
}
