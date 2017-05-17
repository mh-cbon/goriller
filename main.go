// Package goriller generate gorilla routers.
package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/mh-cbon/astutil"
	"github.com/mh-cbon/httper/utils"
)

var name = "goriller"
var version = "0.0.0"

func main() {

	var help bool
	var h bool
	var ver bool
	var v bool
	var outPkg string
	var mode string
	flag.BoolVar(&help, "help", false, "Show help.")
	flag.BoolVar(&h, "h", false, "Show help.")
	flag.BoolVar(&ver, "version", false, "Show version.")
	flag.BoolVar(&v, "v", false, "Show version.")
	flag.StringVar(&outPkg, "p", "", "Package name of the new code.")
	flag.StringVar(&mode, "mode", "std", "The generation mode std|rpc.")

	flag.Parse()

	if ver || v {
		showVer()
		return
	}
	if help || h {
		showHelp()
		return
	}

	if flag.NArg() < 1 {
		panic("wrong usage")
	}
	args := flag.Args()

	out := ""
	if args[0] == "-" {
		args = args[1:]
		out = "-"
	}

	todos, err := utils.NewTransformsArgs(utils.GetPkgToLoad()).Parse(args)
	if err != nil {
		panic(err)
	}

	filesOut := utils.NewFilesOut("github.com/mh-cbon/" + name)

	for _, todo := range todos.Args {
		if todo.FromPkgPath == "" {
			log.Println("Skipped ", todo.FromTypeName)
			continue
		}

		fileOut := filesOut.Get(todo.ToPath)

		fileOut.PkgName = outPkg
		if fileOut.PkgName == "" {
			fileOut.PkgName = findOutPkg(todo)
		}

		if err := processType(mode, todo, fileOut); err != nil {
			log.Println(err)
		}
	}
	filesOut.Write(out)
}

func showVer() {
	fmt.Printf("%v %v\n", name, version)
}

func showHelp() {
	showVer()
	fmt.Println()
	fmt.Println("Usage")
	fmt.Printf("	%v [-p name] [-mode name] [...types]\n\n", name)
	fmt.Printf("  types:  A list of types such as src:dst.\n")
	fmt.Printf("          A type is defined by its package path and its type name,\n")
	fmt.Printf("          [pkgpath/]name\n")
	fmt.Printf("          If the Package path is empty, it is set to the package name being generated.\n")
	// fmt.Printf("          If the Package path is a directory relative to the cwd, and the Package name is not provided\n")
	// fmt.Printf("          the package path is set to this relative directory,\n")
	// fmt.Printf("          the package name is set to the name of this directory.\n")
	fmt.Printf("          Name can be a valid type identifier such as TypeName, *TypeName, []TypeName \n")
	fmt.Printf("  -p:     The name of the package output.\n")
	fmt.Printf("  -mode:  The generation mode std|rpc.\n")
	fmt.Println()
}

func findOutPkg(todo utils.TransformArg) string {
	if todo.ToPkgPath != "" {
		prog := astutil.GetProgramFast(todo.ToPkgPath)
		if prog != nil {
			pkg := prog.Package(todo.ToPkgPath)
			return pkg.Pkg.Name()
		}
	}
	if todo.ToPkgPath == "" {
		prog := astutil.GetProgramFast(utils.GetPkgToLoad())
		if len(prog.Imported) < 1 {
			panic("impossible, add [-p name] option")
		}
		for _, p := range prog.Imported {
			return p.Pkg.Name()
		}
	}
	if strings.Index(todo.ToPkgPath, "/") > -1 {
		return filepath.Base(todo.ToPkgPath)
	}
	return todo.ToPkgPath
}

func processType(mode string, todo utils.TransformArg, fileOut *utils.FileOut) error {
	dest := &fileOut.Body
	srcName := todo.FromTypeName
	destName := todo.ToTypeName

	prog := astutil.GetProgramFast(todo.FromPkgPath)
	pkg := prog.Package(todo.FromPkgPath)
	foundMethods := astutil.FindMethods(pkg)

	srcConcrete := astutil.GetUnpointedType(srcName)
	// the json input must provide a key/value for each params.
	structType := astutil.FindStruct(pkg, srcConcrete)
	structComment := astutil.GetComment(prog, structType.Pos())
	// todo: might do better to send only annotations or do other improvemenets.
	structComment = makeCommentLines(structComment)
	structAnnotations := astutil.GetAnnotations(structComment, "@")

	// fileOut.AddImport("io", "")
	// fileOut.AddImport("strings", "")
	fileOut.AddImport("github.com/gorilla/mux", "")
	// fileOut.AddImport("github.com/mh-cbon/httper/lib", "httper")

	// cheat.
	// fmt.Fprintf(dest, `var xxStringsSplit = strings.Split
	// `)

	// Declare the new type
	fmt.Fprintf(dest, `
// %v is a goriller of %v.
%v
type %v struct{
	embed %v
}
		`, destName, srcName, structComment, destName, srcName)

	// Make the constructor
	fmt.Fprintf(dest, `// New%v constructs a goriller of %v
func New%v(embed %v) *%v {
	ret := &%v{
		embed: embed,
	}
  return ret
}
`, destName, srcName, destName, srcName, destName, destName)

	fmt.Fprintf(dest, `// Bind the given router.
func (t %v) Bind(router *mux.Router) {
`, destName)

	for _, m := range foundMethods[srcConcrete] {
		methodName := astutil.MethodName(m)

		expr := ""
		if mode == "std" {

			comment := astutil.GetComment(prog, m.Pos())
			annotations := astutil.GetAnnotations(comment, "@")
			annotations = mergeAnnotations(structAnnotations, annotations)

			if route, ok := annotations["route"]; ok {
				route = strings.TrimSpace(route)
				expr += fmt.Sprintf(`.HandleFunc(%q, t.embed.%v)`, route, methodName)
			}

			if name, ok := annotations["name"]; ok {
				name = strings.TrimSpace(name)
				expr += fmt.Sprintf(`.Name(%q)`, name)
			}
			if methods, ok := annotations["methods"]; ok {
				methods = strings.TrimSpace(methods)
				methods = stringifyList(methods)
				if methods != "" {
					expr += fmt.Sprintf(`.Methods(%v)`, methods)
				}
			}
			if schemes, ok := annotations["schemes"]; ok {
				schemes = strings.TrimSpace(schemes)
				schemes = stringifyList(schemes)
				if schemes != "" {
					expr += fmt.Sprintf(`.Schemes(%v)`, schemes)
				}
			}
			// if headers, ok := annotations["headers"]; ok {
			// 	headers = strings.TrimSpace(headers)
			// 	headers = stringifyList(headers)
			// 	if headers != "" {
			// 		expr += fmt.Sprintf(`.Headers(%v)`, headers)
			// 	}
			// }
			if host, ok := annotations["host"]; ok {
				host = strings.TrimSpace(host)
				if host != "" {
					expr += fmt.Sprintf(`.Host(%v)`, host)
				}
			}
		} else {
			expr += fmt.Sprintf(`.HandleFunc(%q, t.embed.%v)`, methodName, methodName)

		}

		if expr != "" {
			fmt.Fprintf(dest, "router%v", expr)
		}
		fmt.Fprintln(dest)
	}

	fmt.Fprintf(dest, `
}`)

	return nil
}

func mergeAnnotations(structAnnot, methodAnnot map[string]string) map[string]string {
	ret := map[string]string{}
	for k, v := range methodAnnot {
		ret[k] = v
	}
	for k, v := range structAnnot {
		if _, ok := ret[k]; !ok {
			ret[k] = v
		}
	}
	return ret
}

func makeCommentLines(s string) string {
	s = strings.TrimSpace(s)
	comment := ""
	for _, k := range strings.Split(s, "\n") {
		comment += "// " + k + "\n"
	}
	comment = strings.TrimSpace(comment)
	if comment == "" {
		comment = "//"
	}
	return comment
}

func stringifyList(s string) string {
	ret := ""
	for _, l := range strings.Split(s, ",") {
		l = strings.TrimSpace(l)
		ret += fmt.Sprintf("%q, ", l)
	}
	if ret != "" {
		ret = ret[:len(ret)-2]
	}
	return ret
}
