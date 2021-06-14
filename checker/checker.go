package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
		if exists, _ := blistPkg[lib]; !exists {
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
			val := reflect.ValueOf(n).Elem()
			if exists, _ := blistTypes[n.Kind.String()]; exists {
				fmt.Println("basic lit", n, n.Kind)
				fmt.Println("!!!" + n.Kind.String(), "is in types blacklist")
			} else {
				fmt.Println("!!!", n, val, n.Kind, n.Kind.String(), "is in types whitelist")
			}
		// shorthand map declaration
		case *ast.AssignStmt:
			// if lhs is an identifier
			fmt.Println("len lhs=", len(n.Lhs), "len rhs=", len(n.Rhs), n.Lhs, n.Rhs) 
			for idx, ident := range n.Lhs {
				if reflect.ValueOf(ident).Elem().Type().String() == "ast.Ident" {
					rhsIdx := idx
                    // multiple vars have a single value
                    if idx >= len(n.Rhs) {
                    	rhsIdx = 0
                    }
					if reflect.TypeOf(n.Rhs[rhsIdx]).String() == "*ast.CallExpr" {
						if len(n.Rhs[rhsIdx].(*ast.CallExpr).Args) != 0 && reflect.ValueOf(n.Rhs[rhsIdx].(*ast.CallExpr).Args[0]).Elem().Type().String() == "ast.MapType" {
							mapVars[ident.(*ast.Ident).Name] = true
						}
					}
					if reflect.TypeOf(n.Rhs[rhsIdx]).String() == "*ast.CompositeLit" {
                    	if reflect.ValueOf(n.Rhs[rhsIdx].(*ast.CompositeLit).Type).Elem().Type().String() == "ast.MapType" {
                        	mapVars[ident.(*ast.Ident).Name] = true
                        }
                   	}
				}
			}
			fmt.Println("map vars",mapVars)
			if len(n.Lhs) != 0 && reflect.ValueOf(n.Lhs[0]).Elem().Type().String() == "ast.Ident" {
				// if rhs is a call expression defining a map
				if reflect.TypeOf(n.Rhs[0]).String() == "*ast.CallExpr" && len(n.Rhs[0].(*ast.CallExpr).Args) != 0 && reflect.ValueOf(n.Rhs[0].(*ast.CallExpr).Args[0]).Elem().Type().String() == "ast.MapType" { 
					mapVars[n.Lhs[0].(*ast.Ident).Name] = true
                	//fmt.Println("---->",n.Rhs[0].(*ast.CallExpr).Args[0])
                	fmt.Println("map vars",mapVars)
				}
			}
			/*
			if reflect.TypeOf(n.Rhs[0]).String() == "*ast.CallExpr" && len(n.Rhs[0].(*ast.CallExpr).Args) != 0 && len(n.Lhs) != 0 {
				//mapVars[n.Rhs[0].(*ast.CallExpr).Args[0].Name] = true
				if reflect.ValueOf(n.Lhs[0]).Elem().Type().String() == "ast.Ident" && n.Rhs[0].(*ast.CallExpr).Args[0] != nil {
					//fmt.Println(n.Lhs[0].(*ast.Ident).Obj.Type)
					fmt.Println(n.Rhs[0].(*ast.CallExpr).Args[0],  reflect.ValueOf(n.Rhs[0].(*ast.CallExpr).Args[0]).Elem().Type().String() )
					mapVars[n.Lhs[0].(*ast.Ident).Name] = true
				fmt.Println("---->",n.Rhs[0].(*ast.CallExpr).Args[0])
            	fmt.Println("map vars",mapVars)
				}
			}
			*/
		// map declaration using var
        case *ast.GenDecl:
			fmt.Println("++++",n.Tok)
			if n.Tok == token.CONST || n.Tok == token.VAR {
				for _, s := range n.Specs {
					fmt.Println("++++",s, reflect.ValueOf(s.(*ast.ValueSpec).Type))
					for idx, ident := range s.(*ast.ValueSpec).Names {
						if reflect.TypeOf(ident).String() == "*ast.Ident" {
						fmt.Println("on name", ident.Name)
						// there is just one type, because the variable doesn't receive an initial value
						//if reflect.ValueOf(s.(*ast.ValueSpec).Type).IsValid() {
						if len(s.(*ast.ValueSpec).Values) == 0 {
							fmt.Println("s.(*ast.ValueSpec).Type().String()=", reflect.ValueOf(s.(*ast.ValueSpec).Type).Elem().Type().String())
							if reflect.ValueOf(s.(*ast.ValueSpec).Type).Elem().Type().String() == "ast.MapType" {
								mapVars[ident.Name] = true
							}
						} else {
						// variables are initialized
							rhsIdx := idx
							if  idx >= len(s.(*ast.ValueSpec).Values) {
								rhsIdx = 0
							}
							fmt.Println("we are here",reflect.ValueOf(s.(*ast.ValueSpec).Values[rhsIdx]).Elem().Type().String())
							if reflect.TypeOf(s.(*ast.ValueSpec).Values[rhsIdx]).String() == "*ast.CallExpr" {
                        		if len(s.(*ast.ValueSpec).Values[rhsIdx].(*ast.CallExpr).Args) != 0 && reflect.ValueOf(s.(*ast.ValueSpec).Values[rhsIdx].(*ast.CallExpr).Args[0]).Elem().Type().String() == "ast.MapType" {
                            		mapVars[ident.Name] = true
                        		}
                    		}
							if reflect.TypeOf(s.(*ast.ValueSpec).Values[rhsIdx]).String() == "*ast.CompositeLit" {
								if reflect.ValueOf(s.(*ast.ValueSpec).Values[rhsIdx].(*ast.CompositeLit).Type).Elem().Type().String() == "ast.MapType" {
									mapVars[ident.Name] = true
								}
							}
						}
						}
					}
				}
			fmt.Println("map vars",mapVars)
			}

/*
					if reflect.ValueOf(s.(*ast.ValueSpec).Type).IsValid() {
					fmt.Println("s.(*ast.ValueSpec).Type().String()=", reflect.ValueOf(s.(*ast.ValueSpec).Type).Elem().Type().String())
					if reflect.ValueOf(s.(*ast.ValueSpec).Type).Elem().Type().String() == "ast.MapType" {
						for _, ident := range s.(*ast.ValueSpec).Names {
							mapVars[ident.Name] = true
						}
						fmt.Println("map vars",mapVars)
					}
					}
				}
			}
*/
		case *ast.ExprStmt:
			fmt.Println("x=",n.X)
		// value declaration
		case *ast.Ident:
			//val := reflect.ValueOf(n).Elem()
			if n.Obj != nil {
				fmt.Println("***********",n.Name, n.Obj.Type)
			}
			if n.Obj != nil && n.Obj.Type == "MapType" {
				fmt.Println("!!! map name" + n.Name, n.Obj.Name, "---------")
				mapVars[n.Name] = true
			}
			fmt.Println("map vars",mapVars)
		// check for ranges along maps
		case *ast.RangeStmt:
			rangeLit := n.X.(*ast.Ident)
            if mapVars[rangeLit.Name] == true {
                fmt.Println("!!! !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!", rangeLit.Name, "is an iterated map")
           	}
		default:
			if n != nil {
				val := reflect.ValueOf(n).Elem()
				valType := val.Type().Name()


				if exists, _ := blistTypes[valType]; exists {
					fmt.Println("default",n, val.Type())
					fmt.Println("!!!", valType, "is in types blacklist")
				}  else {
					fmt.Println("!!!", n, val, val.Type(), "is in types whitelist")
				}
			}
		}
		return true
	})
}

func main() {
	spath := flag.String("s", "", "source file")
	bppath := flag.String("bp", "", "blacklist packages file")
	btpath := flag.String("bt", "", "blacklist types file")
	flag.Parse()
	blistPkg := getBlacklist(bppath)
	blistTypes := getBlacklist(btpath)
	if blistPkg == nil || blistTypes == nil {
		os.Exit(1)
	}
	analyzeSource(spath, blistPkg, blistTypes)
}
