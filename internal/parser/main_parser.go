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
func (self *Task) To_raw_string()string{
	str:=""
	str += self.Command + " "
	for i := 0; i < len(self.Flags); i++ {
		str += self.Flags[i].Key
		if self.Flags[i].Value != ""{
			str += "="+self.Flags[i].Value
		}
		str+= " "
	}
	return str
}
type Flag struct {
	Key   string `Minus @(Word|Ident)`
	Value string `(Equal @(String|Named|Word|Ident|Ident2|Any|Num|An))?`
}
func Get_parser()*participle.Parser[INI]{
	graphQLLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Named", `([a-zA-Z]|[0-9])+\.([a-zA-Z]|[0-9])*`},
		{"Ident", `[a-zA-Z]+[0-9]+`},
		{"Ident2", `([a-zA-Z]|_)+[0-9]+`},
		{"An", `[0-9]+[a-zA-Z]+`},
		{"Num", `-?[0-9]+`},
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
	#ENTRADA PRIMER PARTE MIA
#Seccion A - Sergie Arizandieta
#1S 2024
#CAMBIAR /home/serchiboi -> POR SU USUARIO EJ ->  /home/SU_USER 
#CAMBIAR LOS IDS

#DISCO X
#mkdisk -param=x -size=30 -path=/home/serchiboi/archivos/Disco.dsk # ERR PARAMETROS

#CREACION DE DISCOS---------------------------------------------------
Mkdisk -size=50 -unit=M -fit=FF                   # 50M A
Mkdisk -unit=k -size=51200 -fit=BF                # 50M B
mkDisk -size=14                                   # 13M C
mkdisk -size=51200 -unit=K                        # 50M D
mkDisk -size=20 -unit=M -fit=WF                   # 20M E
Mkdisk -size=50 -unit=M -fit=FF                   # 50M F X
Mkdisk -size=50 -unit=M -fit=FF                   # 50M G X
mkdisk -size=51200 -unit=K                        # 50M H X
mkdisk -size=51200 -unit=K                        # 50M I X

#ELIMINACION DE DISCOS---------------------------------------------------
rmdisk -driveletter=Z #ERR RUTA NO ENCONTRADA
rmdisk -driveletter=F
rmdisk -driveletter=G
rmdisk -driveletter=H
rmdisk -driveletter=I


#CREACION DE PARTICIONES---------------------------------------------------
#DISCO 1
fdisk -type=P -unit=b -name=Part0 -size=10485760 -driveletter=Z -fit=BF # ERR RUTA NO ENCONTRADA
fdisk -type=P -unit=b -name=Part1 -size=10485760 -driveletter=A -fit=BF # 10M
fdisk -type=P -unit=k -name=Part2 -size=10240 -driveletter=A -fit=BF    # 10M
fdisk -type=P -unit=M -name=Part3 -size=10 -driveletter=A -fit=BF       # 10M
fdisk -type=P -unit=b -name=Part4 -size=10485760 -driveletter=A -fit=BF # 10M
fdisk -type=P -unit=b -name=Part5 -size=10485760 -driveletter=A -fit=BF #ERR PARTICION 5
# LIBRE DISCO 1: 50-4*10 = 10 -> 20%

#DISCO 2
fdisk -type=L -unit=k -name=Part6 -size=10240 -driveletter=B -fit=BF #ERRROR no hay una extendida
fdisk -type=L -unit=k -name=Part7 -size=10240 -driveletter=B -fit=BF #ERRROR no hay una extendida
fDisk -type=P -unit=K -name=Part8 -size=10240 -driveletter=B -fit=BF    # 10M
fDisk -type=P -unit=m -name=Part9 -size=10 -driveletter=B -fit=FF       # 10M
fDisk -type=P -unit=K -name=Part10 -size=5120 -driveletter=B -fit=WF    # 5M
fdisk -type=E -unit=m -name=Part11 -size=20 -driveletter=B            # 20M
fdisk -type=L -unit=k -name=Part12 -size=1536 -driveletter=B          # 1.5M
fdisk -type=L -unit=k -name=Part13 -size=1536 -driveletter=B -fit=BF
fdisk -type=L -unit=k -name=Part14 -size=1536 -driveletter=B -fit=FF
fdisk -type=L -unit=k -name=Part15 -size=1536 -driveletter=B -fit=BF
fdisk -type=L -unit=k -name=Part16 -size=1536 -driveletter=B -fit=WF
fdisk -type=L -unit=k -name=Part17 -size=1536 -driveletter=B -fit=BF
fdisk -type=L -unit=k -name=Part18 -size=1536 -driveletter=B -fit=FF
fdisk -type=L -unit=k -name=Part19 -size=1536 -driveletter=B -fit=BF
fdisk -type=L -unit=k -name=Part20 -size=1536 -driveletter=B -fit=FF
fdisk -type=L -unit=k -name=Part21 -size=1536 -driveletter=B -fit=BF
fdisk -type=L -unit=k -name=Part22 -size=1536 -driveletter=B -fit=wF
fdisk -type=L -unit=k -name=Part23 -size=1536 -driveletter=B -fit=BF
fdisk -type=L -unit=k -name=Part24 -size=1536 -driveletter=B -fit=FF
# LIBRE DISCO 2: 50-45 = 5 -> 10%
# LIBRE EXTENDIDA 2: 20-13*1.5 = 0.5 -> 2.5% (por los EBR deberia ser menos)

#DISCO 3
fdisk -type=P -unit=m -name=Part25 -size=20 -driveletter=C    # 20M #ERR FALTA ESPACIO
fdisk -type=P -unit=m -name=Part26 -size=4 -driveletter=C     #4M
fdisk -type=P -unit=m -name=Part27 -size=4 -driveletter=C     #4M
fdisk -type=P -unit=m -name=Part28 -size=1 -driveletter=C     #1M
#LIBRE DISCO 3: 14-9= 5 -> 35.71%

#ELIMINAR Y AGREGAR ESPACIO DISCO 3
fdisk -add=-1000 -unit=m -driveletter=C -name=Part26 # ERR SIZE NEGATIVO
fdisk -add=1000 -unit=m -driveletter=C -name=Part26 # ERR PARTICION NO TIENE ESPACIO
fdisk -add=-2 -unit=m -driveletter=C -name=Part26 # 4-2= 2M
fdisk -delete=full -name=Part27 -driveletter=C # 0
fdisk -add=4 -unit=m -driveletter=C -name=Part28 # 4+1= 5M
#LIBRE DISCO 3: 14-7 = 3 -> 50%

#DISCO 5
fdisk -type=E -unit=k -name=Part29 -size=5120 -driveletter=E -fit=BF # 5MB
fdisk -type=L -unit=k -name=Part30 -size=1024 -driveletter=E -fit=BF # 1MB
fdisk -type=P -unit=k -name=Part31 -size=5120 -driveletter=E -fit=BF # 5MB
fdisk -type=L -unit=k -name=Part32 -size=1024 -driveletter=E -fit=BF # 1MB
fdisk -type=L -unit=k -name=Part33 -size=1024 -driveletter=E -fit=BF # 1MB
fdisk -type=L -unit=k -name=Part34 -size=1024 -driveletter=E -fit=BF # 1MB
# LIBRE DISCO 5: 20-10 = 5 -> 50%
# LIBRE EXTENDIDA 2: 5-4 = 1 -> 20% (por los EBR deberia ser menos)

#MONTAR PARTICIONES---------------------------------------------------
#DISCO X
mount -driveletter=A -name=Part5 #ERR PARTICION NO EXISTE
#DISCO 1
mount -driveletter=A -name=Part1 #191a -> id1 -> cambiar el 191a por el ID que a ustedes les genera
mount -driveletter=A -name=Part2 #191b -> id2 -> cambiar el 191b por el ID que a ustedes les genera
mount -driveletter=A -name=Part1 #ERR PARTICION YA MONTADA
#DISCO 2
mount -driveletter=B -name=Part11 #ERR MONTAR EXTENDIDA
mount -driveletter=B -name=Part8 #192a -> id3 -> cambiar el 192a por el ID que a ustedes les genera
mount -driveletter=B -name=Part9 #192b -> id4 -> cambiar el 192b por el ID que a ustedes les genera
#DISCO 3
mount -driveletter=C -name=Part26 #193a -> id5 -> cambiar el 193a por el ID que a ustedes les genera
#DISCO 5
mount -driveletter=E -name=Part31 #194a -> id6 -> cambiar el 194a por el ID que a ustedes les genera


#DESMONTAR PARTICIONES---------------------------------------------------
unmount -id=IDx #ERR NO EXISTE ID
#DISCO 1
unmount -id=191a #-> id1
unmount -id=191a #ERR PARTICION YA DESMONTADA -> id1
#DISCO 2
unmount -id=192b #-> id4


#REPORTES---------------------------------------------------
#DISCO 1
rep -id=191a -Path=/home/serchiboi/archivos/reports/reporte1.jpg -name=mbr #ERR ID NO ENCONTRADO -> id1
rep -id=191b -Path=/home/serchiboi/archivos/reports/reporte2.jpg -name=disk #-> id2
rep -id=191b -Path=/home/serchiboi/archivos/reports/reporte3.jpg -name=mbr #-> id2

#DISCO 2
rep -id=192b -Path=/home/serchiboi/archivos/reports/reporte4.jpg -name=mbr #ERR ID NO ENCONTRADO -> id4
rep -id=192a -Path=/home/serchiboi/archivos/reports/reporte5.jpg -name=disk #-> id3
rep -id=192a -Path=/home/serchiboi/archivos/reports/reporte6.jpg -name=mbr #-> id3

#DISCO 3
rep -id=IDx -Path=/home/serchiboi/archivos/reports/reporte7.jpg -name=mbr #ERR ID NO ENCONTRADO
rep -id=193a -Path=/home/serchiboi/archivos/reports/reporte8.jpg -name=disk #-> id5
rep -id=193a -Path=/home/serchiboi/archivos/reports/reporte9.jpg -name=mbr #-> id5

#DISCO 5
rep -id=IDx -Path=/home/serchiboi/archivos/reports/reporte10.jpg -name=mbr #ERR ID NO ENCONTRADO
rep -id=194a -Path=/home/serchiboi/archivos/reports/reporte11.jpg -name=disk #-> id6
rep -id=194a -Path=/home/serchiboi/archivos/reports/reporte12.jpg -name=mbr #-> id6

#exec -path=../basico.mia
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
