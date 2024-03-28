package parser

import (
	"fmt"

	"strings"

	// "strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type INI struct {
	Tasks []*Task `@@* `
}

type Task struct {
	Command string  `@Word`
	Flags   []*Flag `@@*`
}

type Flag struct {
	Key   string `Minus @(Word|Ident)`
	Value string `(Equal @(String|Named|Word|Ident|Any))?`
}
func Get_parser()*participle.Parser[INI]{
	graphQLLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Named", `([a-zA-Z]|[0-9])+\.([a-zA-Z]|[0-9])*`},
		{"Ident", `[a-zA-Z]+[0-9]+`},
		{"Word", `[a-zA-Z]+`},
		{"Comment", `(#)[^\n]*(\n|^)?`},
		{"Whitespace", `[ \t\n\r]+`},
		{"String", `"(\\"|[^"])*"`},
		{"EOL", `[\n\r]+`},
		{"Minus", `-`},
		{"Equal", `=`},
		{"Any", `[^ \t\n\r]+`},
	} )
	parser := participle.MustBuild[INI](
		participle.Lexer(graphQLLexer),
		participle.Elide("Whitespace","EOL","Comment"),
		participle.UseLookahead(2),
	)
	
	return parser
}
func Some_test() {
	
	line := `
	execute -path=./asdf.txt
	rep -id=191b -Path=/home/serchiboi/archivos/reports/reporte3.jpg -name=mbr #-> id2
	rep -id=192a -path="/home/serchiboi/archivos/reportes/reporte11_tree.jpg" -name=tree # asdf - sdf=:"sdf" asdfsdafsda;kljfds -we -sdf 
	mkfile -path=/2-1/FFFF.txt -size=280320
	withString -path="/2 -1/ FF FF.txt" 
	WithoutVal -path="sdfghadfgt" -some
	cualeuri -algo=A299
	execute -path=./asdf.txt
	Pause
	rename -path=/Pelis/canon.txt -name=solar.txt
    `
	// line := `rename -path=/Pelis/canon.txt -name=solar.txt`
	// graphQLLexer := lexer.MustSimple([]lexer.SimpleRule{
	// 	{"Ident", `[a-zA-Z]+[0-9]+`},
	// 	{"Word", `[a-zA-Z]+`},
	// 	{"Comment", `(#)[^\n]*(\n|^)?`},
	// 	{"Whitespace", `[ \t\n\r]+`},
	// 	{"String", `"(\\"|[^"])*"`},
	// 	{"EOL", `[\n\r]+`},
	// 	{"Minus", `-`},
	// 	{"Equal", `=`},
	// 	{"Any", `[^ \t\n\r]+`},
	// 	// {"Any", `(\d|[a-zA-Z0-9]|[^ \t\n\r])+`},
	// } )
	// parser := participle.MustBuild[INI](
	// 	participle.Lexer(graphQLLexer),
	// 	participle.Elide("Whitespace","EOL","Comment"),
	// 	participle.UseLookahead(2),
	// )
	parser:=Get_parser()
	ini, err := parser.ParseString("", line)
	if err != nil {
		fmt.Println(line)
		for i := 0; i < 42; i++ {
			fmt.Print(" ")
		}
		fmt.Println("|")
		panic(fmt.Sprintf("Error al parsear la línea '%s': %v\n", line, err))
	}
	for _, task := range ini.Tasks {
		fmt.Printf("Command: %s\n", task.Command)
		for _, flag := range task.Flags {
			fmt.Printf("  Flag Key: %s, Flag Value: <%s>\n", flag.Key, flag.Value)
		}
	}
}

func (self *Task)get_value(param string) (string,error){
	for _, f := range self.Flags {
		if strings.ToLower(f.Key) == param{
			if strings.HasPrefix(f.Value,"\""){
				return f.Value[1 : len(f.Value)-1],nil
			}
			return f.Value,nil
		}
	}
	return "",fmt.Errorf("The parameter %s was not defined",param)
}


// rmdisk
// fdisk
// mount
// unmount
// mkfs
// ------------------------------------------------------------
// login
// logout
// mkgrp
// rmgrp
// mkusr
// rmusr
// ------------------------------------------------------------
// mkfile
// cat
// remove
// edit
// rename
// mkdir
// copy
// move
// find
// chown
// chgrp
// chmod
// pause