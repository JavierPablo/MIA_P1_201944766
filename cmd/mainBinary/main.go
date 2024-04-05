package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"project/internal/aplication"
	"project/internal/datamanagment"
	"project/internal/formats"
	"project/internal/parser"
	"project/internal/utiles"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/fatih/color"
)
var Ok = color.New(color.FgGreen)
var Result = color.New(color.FgCyan)
var Err = color.New(color.FgRed)
var Line = color.New(color.FgBlack)
var Print = color.New(color.FgYellow)
func main(){
	
	
	main_program()
	// fmt.Println(string(utiles.File))
	// parser.Some_test()
}

func main_program(){
	Ok.Println("Bienvenido al systema")
	app:=aplication.Aplication{}
	parser:=parser.Get_parser()
	const DISK_CONTAINER_PATH string = "./MIA/P1"
	ioservice_pool := datamanagment.New_IOServicePool(DISK_CONTAINER_PATH)
	reader := bufio.NewReader(os.Stdin)
	for{
		input,err := reader.ReadString('\n')
		if err != nil {
			Err.Println("Error reading input:", err)
			return
		}
		execute(input,nil,&app,&ioservice_pool,DISK_CONTAINER_PATH,parser,true)
	}
	
	

}
func execute(input string,or_tasks *[]*parser.Task,app *aplication.Aplication, ioservice_pool *datamanagment.IOServicePool,disk_path string,parser_engine *participle.Parser[parser.INI],record_activity bool) {
	var tasks []*parser.Task 
	if or_tasks != nil{
		tasks = *or_tasks
	}else{
		parsing, err := parser_engine.ParseString("", input)
		if err != nil {
			panic(fmt.Sprintf("Error al parsear: %v\n", err))
		}
		tasks=parsing.Tasks
	}
	for _, task := range tasks {
		if task.Command != "print"{
			Line.Println(task.To_raw_string())
		}
		switch strings.ToLower(task.Command) {
		case "mkdisk":
			params,err:=task.Get_MkdiskParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			_,_,err = app.Make_disk(params.Size,params.Unit,params.Fit,disk_path)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "rmdisk":
			params,err:=task.Get_RmdiskParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			err = app.Remove_disk(params.Driveletter,disk_path)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "fdisk":
			params,err:=task.Get_FdiskParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			io,err := ioservice_pool.Get_service_with(params.Driveletter)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if params.Delete{
				err=app.Remove_partition_disk(io,params.Name)
			}else if params.Add != 0{
				err = app.Modify_partition_size_in_disk(params.Add,io,params.Name,params.Unit)

			}else{	
				_,edrr := app.Partition_disk(params.Size,io,params.Name,params.Unit,params.Type,params.Fit)
				err = edrr
			}
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "mount":
			params,err:=task.Get_MountParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			io,err := ioservice_pool.Get_service_with(params.Driveletter)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}			
			err= app.Mount_partition(io,params.Name,params.Driveletter)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}			

		case "unmount":
			params,err:=task.Get_UnmountParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			err = app.Unmount_partition(params.Id)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "mkfs":
			params,err:=task.Get_MkfsParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			err = app.Format_mounted_partition(params.Id,params.Type,params.Fs)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
// -----------------------------------------------------------------------
		case "login":
			params,err:=task.Get_LoginParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			journal,err := app.Log_in_user(params.Id,params.User,params.Pass)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Login,[]string{params.User,params.Pass}))
			}
		case "logout":
			journal,err :=app.Log_out()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Unlog,[]string{}))
			}
		case "mkgrp":
			params,err:=task.Get_MkgrpParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			journal,err := app.Make_group(params.Name)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Make_group,[]string{params.Name}))
			}
		case "rmgrp":
			params,err:=task.Get_RmgrpParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			journal,err := app.Remove_group(params.Name)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Remove_group,[]string{params.Name}))
			}
		case "mkusr":
			params,err:=task.Get_MkusrParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			journal,err := app.Make_user(params.User,params.Pass,params.Grp)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Make_user,[]string{params.User,params.Pass,params.Grp}))
			}
		case "rmusr":
			params,err:=task.Get_RmusrParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			journal,err := app.Remove_user(params.User)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Remove_user,[]string{params.User}))
			}
// -----------------------------------------------------------------------
		case "mkfile":
			params,err:=task.Get_MkfileParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			var content []string

			if params.Fixedcont != ""{
				content = strings.Split(params.Fixedcont, "")
			}else if params.Cont != ""{
				b, err := os.ReadFile(params.Cont)
				if err != nil {
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					continue
				}
				content = make([]string, 0,len(b))
				for i := 0; i < len(b); i++ {
					content = append(content, string(b[i]))
				}

			}else if params.Size != 0 {
				content = make([]string, 0,params.Size)
				for i := 0; i < int(params.Size); i++ {
					content = append(content, strconv.Itoa(rand.Intn(10)))
				}
			}else{
				Err.Printf("Command failed in execution:\nThere is no content for file\n")
				continue
			}
			dirs:=strings.Split(params.Path, "/")[1:]
			file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
			dirs = dirs[:len(dirs) - 1]
			folders:= make([][12]string,0,len(dirs))
			for i := 0; i < len(dirs); i++ {
				folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
			}
			journal,err := app.Make_file(folders,content,file_name,params.R)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				if params.R{
					journal.Push_instruction(formats.New_inst(formats.Make_file,[]string{"Y",params.Path,strings.Join(content, "")}))
				}else{
					journal.Push_instruction(formats.New_inst(formats.Make_file,[]string{"N",params.Path,strings.Join(content, "")}))
				}
			}
		case "cat":
			params,err:=task.Get_CatParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			for i := 0; i < len(params.Paths); i++ {
				dirs:=strings.Split(params.Paths[i], "/")[1:]
				file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
				dirs = dirs[:len(dirs) - 1]
				folders:= make([][12]string,0,len(dirs))
				for i := 0; i < len(dirs); i++ {
					folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
				}
				content,err := app.Show_file(folders,file_name)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					continue
				}
				Result.Println(content)
				Result.Println()

			}
		case "remove":
			params,err:=task.Get_RemoveParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			dirs:=strings.Split(params.Path, "/")[1:]
			file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
			dirs = dirs[:len(dirs) - 1]
			folders:= make([][12]string,0,len(dirs))
			for i := 0; i < len(dirs); i++ {
				folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
			}
			journal,err := app.Remove(folders,file_name)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Remove,[]string{params.Path}))
			}
		case "edit":
			params,err:=task.Get_EditParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			var content []string
			if params.Fixedcont != ""{
				content = strings.Split(params.Fixedcont, "")
			}else if params.Cont != ""{
				b, err := os.ReadFile(params.Cont)
				if err != nil {
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					continue
				}
				content = make([]string, 0,len(b))
				for i := 0; i < len(b); i++ {
					content = append(content, string(b[i]))
				}

			}else{
				Err.Printf("Command \"%s\" failed in execution:\nThere is no content for file",task.Command)
				continue
			}
			dirs:=strings.Split(params.Path, "/")[1:]
			file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
			dirs = dirs[:len(dirs) - 1]
			folders:= make([][12]string,0,len(dirs))
			for i := 0; i < len(dirs); i++ {
				folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
			}
			journal,err := app.Edit_file(folders,content,file_name)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Edit_file,[]string{params.Path, strings.Join(content, "")}))
			}
		case "rename":
			params,err:=task.Get_RenameParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			
			dirs:=strings.Split(params.Path, "/")[1:]
			file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
			dirs = dirs[:len(dirs) - 1]
			folders:= make([][12]string,0,len(dirs))
			for i := 0; i < len(dirs); i++ {
				folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
			}
			journal,err := app.Rename_inode(folders,file_name,utiles.Into_ArrayChar12(params.Name))
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Rename_inode,[]string{params.Path, params.Name}))
			}
		case "mkdir":
			params,err:=task.Get_MkdirParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			
			dirs:=strings.Split(params.Path, "/")[1:]
			file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
			dirs = dirs[:len(dirs) - 1]
			folders:= make([][12]string,0,len(dirs))
			for i := 0; i < len(dirs); i++ {
				folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
			}
			journal,err := app.Make_dir(folders,file_name,params.R)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				if params.R{
					journal.Push_instruction(formats.New_inst(formats.Make_dir,[]string{"Y",params.Path}))
				}else{
					journal.Push_instruction(formats.New_inst(formats.Make_dir,[]string{"N",params.Path}))
				}
			}
		case "copy":
			params,err:=task.Get_CopyParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			
			dirs:=strings.Split(params.Path, "/")[1:]
			file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
			dirs = dirs[:len(dirs) - 1]
			folders:= make([][12]string,0,len(dirs))
			for i := 0; i < len(dirs); i++ {
				folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
			}
			dirs2:=strings.Split(params.Destino, "/")[1:]
			folders2:= make([][12]string,0,len(dirs2))
			for i := 0; i < len(dirs2); i++ {
				folders2 = append(folders2, utiles.Into_ArrayChar12(dirs2[i]))
			}
			journal,err := app.Copy(folders,file_name,folders2)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Copy,[]string{params.Path, params.Destino}))
			}
		case "move":
			params,err:=task.Get_MoveParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			
			dirs:=strings.Split(params.Path, "/")[1:]
			file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
			dirs = dirs[:len(dirs) - 1]
			folders:= make([][12]string,0,len(dirs))
			for i := 0; i < len(dirs); i++ {
				folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
			}
			dirs2:=strings.Split(params.Destino, "/")[1:]
			folders2:= make([][12]string,0,len(dirs2))
			for i := 0; i < len(dirs2); i++ {
				folders2 = append(folders2, utiles.Into_ArrayChar12(dirs2[i]))
			}
			journal,err := app.Move(folders,file_name,folders2)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Move,[]string{params.Path, params.Destino}))
			}
		case "find": 
			params,err:=task.Get_FindParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			var folders [][12]string
			if params.Path == "/"{
				folders= make([][12]string,0,0)
			}else{
				dirs := strings.Split(params.Path, "/")[1:]
				folders= make([][12]string,0,len(dirs))
				for i := 0; i < len(dirs); i++ {
					folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
				}
			}
			SPACE_CHAR:=utiles.New_Char(" ")
			criteria_chars:=[12]utiles.Char{SPACE_CHAR,SPACE_CHAR,SPACE_CHAR,SPACE_CHAR,SPACE_CHAR,SPACE_CHAR,SPACE_CHAR,SPACE_CHAR,SPACE_CHAR,SPACE_CHAR,SPACE_CHAR,SPACE_CHAR}
			chars:=strings.Split(params.Name,"")
			for i := 0; i < len(chars) && i < 12; i++ {
				if chars[i] == "?"{
					criteria_chars[i] = utiles.ANY_CHAR
				}else if chars[i] == "*"{
					for j := i; j < 12; j++ {
						criteria_chars[j] = utiles.ANY_CHAR
					}
					break
				}else{
					criteria_chars[i] = utiles.New_Char(chars[i])
				}
			}
			criteria:=utiles.NameCriteria{
				Chars: criteria_chars,
			}
			app.Find(folders,criteria)

		case "chown":
			params,err:=task.Get_ChownParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			dirs:=strings.Split(params.Path, "/")[1:]
			file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
			dirs = dirs[:len(dirs) - 1]
			folders:= make([][12]string,0,len(dirs))
			for i := 0; i < len(dirs); i++ {
				folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
			}
			journal,err := app.Change_own(folders,file_name,params.User,params.R)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				if params.R{
					journal.Push_instruction(formats.New_inst(formats.Chown,[]string{"Y",params.Path,params.User}))
				}else{
					journal.Push_instruction(formats.New_inst(formats.Chown,[]string{"N",params.Path,params.User}))
				}
			}
		case "chgrp":
		params,err:=task.Get_ChgrpParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			journal,err:=app.Chagne_User_Group(params.User,params.Grp)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				journal.Push_instruction(formats.New_inst(formats.Chown,[]string{params.User,params.Grp}))
			}
		case "chmod":
		params,err:=task.Get_ChmodParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			dirs:=strings.Split(params.Path, "/")[1:]
			file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
			dirs = dirs[:len(dirs) - 1]
			folders:= make([][12]string,0,len(dirs))
			for i := 0; i < len(dirs); i++ {
				folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
			}
			journal,err := app.Chagne_UGO(folders,file_name,params.Ugo,params.R)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if journal!=nil && record_activity{
				if params.R{
					journal.Push_instruction(formats.New_inst(formats.Chown,[]string{"Y",params.Path,params.Ugo}))
				}else{
					journal.Push_instruction(formats.New_inst(formats.Chown,[]string{"N",params.Path,params.Ugo}))
				}
			}
		case "pause":
			var nothing string
			fmt.Scanln(&nothing)
		case "recovery":
			params,err:=task.Get_RecoveryParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			new_tasks,err := app.Recovery_ext3(params.Id)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			execute("",&new_tasks,app,ioservice_pool,disk_path,parser_engine,false)
		case "loss":
			params,err:=task.Get_LossParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			err = app.Loss_ext3(params.Id)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			
		case "execute":
			params,err:=task.Get_ExecuteParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			
			b, err := os.ReadFile(params.Path)
			if err != nil {
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			instructions := string(b)
			execute(instructions,nil,app,ioservice_pool,disk_path,parser_engine,true)
		case "mountid":
			app.Print_mounted()
		case "rep":
			params,err:=task.Get_RepParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			
			var report_content string
			switch params.Name {
			case "mbr":
				ioservice,err:=ioservice_pool.Get_service_with(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					continue
				}
				report_content, err = app.Mbr_repo(ioservice)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					continue
				}
			case "disk":
				ioservice,err:=ioservice_pool.Get_service_with(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
				report_content,err=app.Disk_repos(ioservice)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
			case "inode":
				report_content,err=app.Inode_repos(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
			case "block":
				report_content,err=app.Block_repos(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
			case "bm_inode":
				report_content,err=app.Inode_bitmap_repos(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
			case "bm_block":
				report_content,err=app.Block_bitmap_repos(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
			case "tree":
				report_content,err=app.Tree_repos(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
			case "sb":
				report_content,err=app.Super_block_repo(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
			case "file":
				if params.Ruta == ""{
					Err.Printf("Command \"%s\" failed in execution:\n Path vas not specified",task.Command)
					

					continue
				}
				dirs:=strings.Split(params.Ruta, "/")[1:]
				file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
				dirs = dirs[:len(dirs) - 1]
				folders:= make([][12]string,0,len(dirs))
				for i := 0; i < len(dirs); i++ {
					folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
				}
				content,err := app.Show_file(folders,file_name)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
				report_content = content
			case "ls":
				if params.Ruta == ""{
					Err.Printf("Command \"%s\" failed in execution:\n Path vas not specified",task.Command)
					
					continue
				}
				var folders [][12]string
				if params.Ruta == "/"{
					folders= make([][12]string,0,0)
				}else{
					dirs := strings.Split(params.Ruta, "/")[1:]
					folders= make([][12]string,0,len(dirs))
					for i := 0; i < len(dirs); i++ {
						folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
					}
				}
				report_content,err = app.Ls_report(params.Id,folders)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
			case "journaling":
				report_content,err=app.Journaling(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					
					continue
				}
			}
			err=render_reports(params.Path,report_content)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "print":
			params,err:=task.Get_PrintParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			Print.Println(params.Val)
			continue

		default: 
		Err.Println("Command not recognized for")
		Err.Println(task.Command)
		continue
		}
		Ok.Printf("Comando \"%s\" ejecutado con exito\n",task.Command)
		ioservice_pool.Flush_changes()
	}
}



func render_reports(dest_path string,content string)error{
	dir := filepath.Dir(dest_path)
    if err := os.MkdirAll(dir, 0777); err != nil {return err}
	splited_path:=strings.Split(dest_path, "/")
	last_name:=splited_path[len(splited_path)-1]
	DOT_TEMP_PATH:="./temp.dot"
	switch strings.ToLower(strings.Split(last_name, ".")[1]) {
	case "txt":
		file, err :=os.Create(dest_path)
		if err!=nil{return err}
		file.Write([]byte(content))
		file.Close()
	case "pdf":
		file, err :=os.Create(DOT_TEMP_PATH)
		if err!=nil{return err}
		file.Write([]byte(content))
		file.Close()
		cmd := exec.Command("dot", "-Tpdf",DOT_TEMP_PATH,"-o",dest_path)
		_, err = cmd.Output()
    	if err != nil {return err}
	case "jpg","jpeg":
		file, err :=os.Create(DOT_TEMP_PATH)
		if err!=nil{return err}
		file.Write([]byte(content))
		file.Close()
		cmd := exec.Command("dot", "-Tjpg",DOT_TEMP_PATH,"-o",dest_path)
		b, err := cmd.Output()
		
    	if err != nil {
			Err.Println(string(b))
			return err
		}
	case "png":
		file, err :=os.Create(DOT_TEMP_PATH)
		if err!=nil{return err}
		file.Write([]byte(content))
		file.Close()
		cmd := exec.Command("dot", "-Tpng",DOT_TEMP_PATH,"-o",dest_path)
		_, err = cmd.Output()
    	if err != nil {return err}
		
	}    
    return nil
}
