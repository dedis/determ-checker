package checker

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"strings"
)

func ReadList(path *string) (listMap map[string]bool) {
	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
		return listMap
	}
	defer file.Close()

	listMap = make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		listMap[scanner.Text()] = true
	}
	return listMap
}

func ReadStrBetweenFileOffsets(path *string, start token.Pos, end token.Pos) string {
	file, err := os.Open(*path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	_ = fset.AddFile(*path, -1, int(fi.Size()))
	crtFile := fset.File(start)
	off1 := crtFile.Offset(start)

	buf := make([]byte, end-start+1)
	_, err = file.ReadAt(buf, int64(off1))
	if err != nil {
        log.Fatal(err)
    }

	return string(buf[:])
}

func AnalyzeSource(spath *string, wlistPkg map[string]bool, blistTypes map[string]bool) bool {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, *spath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return false
	}

	fmt.Println("--- Imports ---")
	for _, imp := range node.Imports {
		lib := strings.Replace(imp.Path.Value, "\"", "", -1)
		if exists, _ := wlistPkg[lib]; !exists {
			fmt.Println("!!!", lib, "is NOT a whitelisted package")
		}
	}

	fmt.Println("--- Types ---")

	// contains map literals
	mapVars := make(map[string]bool)

	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CompositeLit:
		case *ast.BasicLit:
			//val := reflect.ValueOf(n).Elem()
			if exists, _ := blistTypes[n.Kind.String()]; exists {
				fmt.Println("!!!", n.Kind.String(), "is in types blacklist")
			} else {
				//fmt.Println("!!!", n, val, n.Kind, n.Kind.String(), "is in types whitelist")
			}
		// detect transcedental fct that return NaNs
		case *ast.CallExpr:
			// assume this is a math function
			startCallOffset := n.Fun.Pos()
			endCallOffset := n.Fun.End()-1
			startArgOffset := n.Lparen
			endArgOffset := n.Rparen
			callStr := ReadStrBetweenFileOffsets(spath, startCallOffset, endCallOffset)
			argStr := ReadStrBetweenFileOffsets(spath, startArgOffset, endArgOffset)
			if strings.Contains(callStr, "math") {
				fmt.Println("call", callStr, argStr)
			}
		// shorthand map declaration
		case *ast.AssignStmt:
			// if lhs is an identifier
			for idx, ident := range n.Lhs {
				if reflect.ValueOf(ident).Elem().Type().String() == "ast.Ident" {
					rhsIdx := idx
					// rhs index is 0 if multiple vars have a single value
					if idx >= len(n.Rhs) {
						rhsIdx = 0
					}
					// check for maps defined with make()
					if reflect.TypeOf(n.Rhs[rhsIdx]).String() == "*ast.CallExpr" {
						if len(n.Rhs[rhsIdx].(*ast.CallExpr).Args) != 0 && reflect.ValueOf(n.Rhs[rhsIdx].(*ast.CallExpr).Args[0]).Elem().Type().String() == "ast.MapType" {
							mapVars[ident.(*ast.Ident).Name] = true
						}
					}
					// check for maps defined both by type and with make()
					// TODO check if this is really needed
					if reflect.TypeOf(n.Rhs[rhsIdx]).String() == "*ast.CompositeLit" {
						if reflect.ValueOf(n.Rhs[rhsIdx].(*ast.CompositeLit).Type).Elem().Type().String() == "ast.MapType" {
							mapVars[ident.(*ast.Ident).Name] = true
						}
					}
				}
			}
			// map declaration using var
		case *ast.GenDecl:
			if n.Tok == token.CONST || n.Tok == token.VAR {
				for _, s := range n.Specs {
					// iterate the identifiers
					for idx, ident := range s.(*ast.ValueSpec).Names {
						// check that they are truly identifiers
						if reflect.TypeOf(ident).String() == "*ast.Ident" {
							// check if there is a single type, because the variable(s) doesn't receive an initial value
							//if reflect.ValueOf(s.(*ast.ValueSpec).Type).IsValid() {
							if len(s.(*ast.ValueSpec).Values) == 0 {
								if reflect.ValueOf(s.(*ast.ValueSpec).Type).Elem().Type().String() == "ast.MapType" {
									mapVars[ident.Name] = true
								}
							} else {
								// variables are initialized
								// if all variables are initialized to the same value, then rhs index is 0
								rhsIdx := idx
								if idx >= len(s.(*ast.ValueSpec).Values) {
									rhsIdx = 0
								}
								// the rhs can be a call to make()
								if reflect.TypeOf(s.(*ast.ValueSpec).Values[rhsIdx]).String() == "*ast.CallExpr" {
									if len(s.(*ast.ValueSpec).Values[rhsIdx].(*ast.CallExpr).Args) != 0 && reflect.ValueOf(s.(*ast.ValueSpec).Values[rhsIdx].(*ast.CallExpr).Args[0]).Elem().Type().String() == "ast.MapType" {
										mapVars[ident.Name] = true
									}
								}
								// or it can be a composite literal of type definition and a call to make
								if reflect.TypeOf(s.(*ast.ValueSpec).Values[rhsIdx]).String() == "*ast.CompositeLit" {
									if reflect.ValueOf(s.(*ast.ValueSpec).Values[rhsIdx].(*ast.CompositeLit).Type).Elem().Type().String() == "ast.MapType" {
										mapVars[ident.Name] = true
									}
								}
							}
						}
					}
				}
				fmt.Println("DEBUG defined map vars", mapVars)
			}

		// check for ranges along maps
		case *ast.RangeStmt:
			rangeLit := n.X.(*ast.Ident)
			// is the iterated variable a map?
			if mapVars[rangeLit.Name] == true {
				fmt.Println("!!! Potential problem: variable", rangeLit.Name, "is an iterated map")
			}
		default:
			if n != nil {
				val := reflect.ValueOf(n).Elem()
				valType := val.Type().Name()
				if exists, _ := blistTypes[valType]; exists {
					fmt.Println("!!!", valType, "is in types blacklist")
				} else {
					//fmt.Println("!!!", n, val, val.Type(), "is in types whitelist")
				}
			}
		}
		return true
	})
	return true
}
