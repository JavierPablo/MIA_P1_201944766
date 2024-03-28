package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"project/internal/aplication"
	"project/internal/datamanagment"
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
		execute(input,&app,&ioservice_pool,DISK_CONTAINER_PATH,parser)
	}
	
	

}
func execute(input string,app *aplication.Aplication, ioservice_pool *datamanagment.IOServicePool,disk_path string,parser *participle.Parser[parser.INI]) {
	
	parsing, err := parser.ParseString("", input)
	if err != nil {
		panic(fmt.Sprintf("Error al parsear: %v\n", err))
	}
	for _, task := range parsing.Tasks {
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
			done := app.Remove_disk(params.Driverletter,disk_path)
			if !done{
				Err.Printf("Command \"%s\" failed in execution:\n",task.Command)
				continue
			}
		case "fdisk":
			params,err:=task.Get_FdiskParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			io,err := ioservice_pool.Get_service_with(params.Driverletter)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			if params.Delete{
				
			}else if params.Add != 0{
				err = app.Modify_partition_size_in_disk(params.Size,io,params.Name,params.Unit)

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
			io,err := ioservice_pool.Get_service_with(params.Driverletter)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}			
			app.Mount_partition(io,params.Name,params.Driverletter)
		case "unmount":
			params,err:=task.Get_UnmountParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			app.Unmount_partition(params.Id)
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
			err = app.Log_in_user(params.Id,params.User,params.Pass)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "logout":
			app.Log_out()
		case "mkgrp":
			params,err:=task.Get_MkgrpParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			err = app.Make_group(params.Name)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "rmgrp":
			params,err:=task.Get_RmgrpParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			err = app.Remove_group(params.Name)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "mkusr":
			params,err:=task.Get_MkusrParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			err = app.Make_user(params.User,params.Pass,params.Grp)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "rmusr":
			params,err:=task.Get_RmusrParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			err = app.Remove_user(params.User)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
// -----------------------------------------------------------------------
		case "mkfile":
			params,err:=task.Get_MkfileParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			var content []string

			if params.Cont != ""{
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
				Err.Printf("Command failed in execution:\nThere is no content for file")
				continue
			}
			dirs:=strings.Split(params.Path, "/")[1:]
			file_name:=utiles.Into_ArrayChar12(dirs[len(dirs) - 1])
			dirs = dirs[:len(dirs) - 1]
			folders:= make([][12]string,0,len(dirs))
			for i := 0; i < len(dirs); i++ {
				folders = append(folders, utiles.Into_ArrayChar12(dirs[i]))
			}
			err = app.Make_file(folders,content,file_name,params.R)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
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
			err = app.Remove(folders,file_name)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "edit":
			params,err:=task.Get_EditParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			var content []string
			if params.Cont != ""{
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
			err = app.Edit_file(folders,content,file_name)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
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
			err = app.Rename_inode(folders,file_name,utiles.Into_ArrayChar12(params.Name))
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
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
			err = app.Make_dir(folders,file_name,params.R)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
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
			err = app.Copy(folders,file_name,folders2)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
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
			err = app.Move(folders,file_name,folders2)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "find": // REQUIERE IMPLEMENTATION
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
			// err = app.
			// if err!=nil{
			// 	Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
			// 	continue
			// }
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
			err = app.Change_own(folders,file_name,params.User,params.R)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "chgrp":
		params,err:=task.Get_ChgrpParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			err=app.Chagne_User_Group(params.User,params.Grp)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
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
			err = app.Chagne_UGO(folders,file_name,params.Ugo,params.R)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
		case "pause":
			var nothing string
			fmt.Scanln(&nothing)
		case "recovery":// REQUIERE IMPLEMENTATION
			var nothing string
			fmt.Scanln(&nothing)
		case "loss":// REQUIERE IMPLEMENTATION
			var nothing string
			fmt.Scanln(&nothing)
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
			execute(instructions,app,ioservice_pool,disk_path,parser)
		case "mountid":
			app.Print_mounted()
		case "rep":// REQUIERE IMPLEMENTATION
			params,err:=task.Get_RepParam()
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			file, err :=os.Create(params.Path)
			if err!=nil{
				Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
				continue
			}
			var dot_str string
			switch params.Name {
			case "mbr":
				ioservice,err:=ioservice_pool.Get_service_with(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					file.Close()
					continue
				}
				dot_str, err = app.Mbr_repo(ioservice)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					file.Close()
					continue
				}

			case "disk":
				ioservice,err:=ioservice_pool.Get_service_with(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					file.Close()
					continue
				}
				dot_str,err=app.Disk_repos(ioservice)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					file.Close()
					continue
				}
			case "inode":
				dot_str,err=app.Inode_repos(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					file.Close()
					continue
				}
			case "block":
				dot_str,err=app.Block_repos(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					file.Close()
					continue
				}
			case "bm_inode":
				dot_str,err=app.Inode_bitmap_repos(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					file.Close()
					continue
				}
			case "bm_block":
				dot_str,err=app.Block_bitmap_repos(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					file.Close()
					continue
				}
			case "tree":
				dot_str,err=app.Tree_repos(params.Id)
				if err!=nil{
					Err.Printf("Command \"%s\" failed in execution:\n%s\n",task.Command,err)
					file.Close()
					continue
				}
			case "sb":
			case "file":
			case "ls":
			case "journaling":
			}
			file.Write([]byte(dot_str))
			file.Close()
		default: 
			Err.Println("Command not recognized")
			continue
		}
		Ok.Printf("Comando \"%s\" ejecutado con exito\n",task.Command)
		// ioservice_pool.Flush_changes()
	}
}


func bestfit(){
	sm := datamanagment.SpaceManager_from_free_spaces([]datamanagment.Space{
		datamanagment.New_Space(0,3),datamanagment.New_Space(5,2),datamanagment.New_Space(10,10) },20)
	sm.Log_chunks_state()
	fmt.Printf("Best fit for 1 in index = %d\n",sm.Best_fit(1))
	sm.Log_chunks_state()
	fmt.Printf("Worst fit for 1 in index = %d\n",sm.Worst_fit(1))
	sm.Log_chunks_state()
	inx := sm.First_fit(6)
	fmt.Printf("First fit for 6 in index = %d\n",inx)
	sm.Log_chunks_state()
	fmt.Println("Occuping that space")
	sm.Ocupe_space_unchecked(int(inx),6)
	sm.Log_chunks_state()
	fmt.Println("Erasing [12,13]")
	sm.Free_space(1,12)
	sm.Log_chunks_state()
	sm.Ocupe_space_unchecked(2,1)
	sm.Log_chunks_state()
	fmt.Println("-----------------------------------------------")
	fmt.Println("Simulating simetric case with occuped spaces")
	fmt.Println("-----------------------------------------------")
	sm,_ = datamanagment.SpaceManager_from_occuped_spaces([]datamanagment.Space{
		datamanagment.New_Space(3,2),datamanagment.New_Space(7,3),
	},20)
	sm.Log_chunks_state()
	fmt.Printf("Best fit for 1 in index = %d\n",sm.Best_fit(1))
	sm.Log_chunks_state()
	fmt.Printf("Worst fit for 1 in index = %d\n",sm.Worst_fit(1))
	sm.Log_chunks_state()
	inx = sm.First_fit(6)
	fmt.Printf("First fit for 6 in index = %d\n",inx)
	sm.Log_chunks_state()
	fmt.Println("Occuping that space")
	sm.Ocupe_space_unchecked(int(inx),6)
	sm.Log_chunks_state()
	fmt.Println("Erasing [12,13]")
	sm.Free_space(1,12)
	sm.Log_chunks_state()
	sm.Ocupe_space_unchecked(2,1)
	sm.Log_chunks_state()
	fmt.Println("Free space in [7,15]")
	sm.Free_space(9,7)
	sm.Log_chunks_state()
	fmt.Println("Free space in [3,5]")
	sm.Free_space(2,3)
	sm.Log_chunks_state()
	fmt.Println("Ocupe space in [9,13]")
	sm.Ocupe_raw_space(5,9)
	sm.Log_chunks_state()
	fmt.Println("Ocupe space in [5,5]")
	sm.Ocupe_raw_space(1,5)
	sm.Log_chunks_state()
}