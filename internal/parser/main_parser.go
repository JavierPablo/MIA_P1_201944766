package parser

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

type INI struct {
	Properties []*Property `@@*`
	Sections   []*Section  `@@*`
  }
//   ---------------------
  type Task struct {
	Command string      `@Ident`
	Flags []*Flag `@@*`
  }
  type Flag struct {
	Key   string `@Ident "="`
	Value string `@Ident | @String`
  }
//   type ValueS interface{ value() }


//   ---------------------
  type Section struct {
	Identifier string      `"[" @Ident "]"`
	Properties []*Property `@@*`
  }
  
  
  type Property struct {
	Key   string `@Ident "="`
	Value Value `@@`
  }
  
  type Value interface{ value() }
  
  type String struct {
	 String string `@String`
  }
  
  func (String) value() {}
  
  type Number struct {
	 Number float64 `@Float | @Int`
  }
  
  func (Number) value() {}

func Some_test(){
	parser, _ := participle.Build[INI](
		participle.Unquote("String"),
		participle.Union[Value](String{}, Number{}),
	  )
	//   _s := `COMMAND -<ident> -<ident> = <val>(int,ident,String) `
	// Task(String,[]Flag)
	// Flag(String,String)
	  ini, _ := parser.ParseString("", `
	  age = 21
	  name = "Bob Smith"
	  
	  [address]
	  city = "Beverly Hills"
	  postal_code = 90210
	  `)
	  fmt.Println(ini.Properties[1])
}
