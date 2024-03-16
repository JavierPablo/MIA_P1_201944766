package main

import (
	// "project/internal/parser"
	"fmt"
	"project/internal/aplication"
	"project/internal/datamanagment"
	"project/internal/formats/ext2"
	"project/internal/types"
	"project/internal/utiles"
	"strings"
)

func main() {
	// a.Make_disk(3,utiles.Mb,utiles.First)
	
	// format_and_write_file_EXT2()
	// format_and_test_dir_and_file_EXT2()
	// recover_and_test_dir_and_file_EXT2()
	// recover_and_erase_file_EXT2()
	
	// a:=aplication.Aplication{}
	// a.Make_disk(10,utiles.Kb,utiles.First)
	// partition_creation_and_mod_test()
	// partition_recovery_test()
	
	
	// parser.Some_test()

	// a.Make_disk(1,utiles.Kb,utiles.First)
	io := datamanagment.IOService_from(aplication.DISK_CONTAINER_PATH+"/B.dsk")
	mbr := types.CreateMasterBootRecord(&io,0)
	fmt.Println(mbr.Dot_label())
}

func parser_testing(){

}



func partition_recovery_test(){
	a:=aplication.Aplication{}
	// a.Make_disk(1,utiles.Kb,utiles.First)
	io := datamanagment.IOService_from(aplication.DISK_CONTAINER_PATH+"/B.dsk")
	mbr := types.CreateMasterBootRecord(&io,0)

	fmt.Println("Try reading 3 logic space ")
	result,ext_part:=a.Get_extended_partition(mbr)
	if !result{panic("There was no extended partition")}
	result,ebr := a.Find_logical_partition_by_name(ext_part,utiles.Into_ArrayChar16("3 logic"))
	if !result{panic("There was no logical partition found")}
	fmt.Println(ebr.Get())

	fmt.Println("Try reading extended space ")
	fmt.Println(ext_part.Get())
	sm:=a.Get_disk_partitions_space_manager(mbr)
	sm.Log_chunks_state()
	sm=a.Get_extended_part_space_manager(ext_part)
	sm.Log_chunks_state()
}
func partition_creation_and_mod_test(){
	a:=aplication.Aplication{}
	// a.Make_disk(1,utiles.Kb,utiles.First)
	io := datamanagment.IOService_from(aplication.DISK_CONTAINER_PATH+"/B.dsk")
	mbr := types.CreateMasterBootRecord(&io,0)
	abs_index := a.Partition_disk(1000,&io,"first",utiles.B,utiles.Primary,utiles.Best)
	fmt.Printf("Primary partition \"first\" at %d\n",abs_index)
	abs_index2 := a.Partition_disk(1000,&io,"second",utiles.B,utiles.Primary,utiles.Best)
	fmt.Printf("Primary partition \"second\" at %d\n",abs_index2)
	abs_index3 := a.Partition_disk(6000,&io,"third extd",utiles.B,utiles.Extendend,utiles.Best)
	fmt.Printf("Extended partition \"third extd\" at %d\n",abs_index3)
	for _,partition := range mbr.Mbr_partitions().Spread(){
		if partition.Part_start().Get() == -1{continue}
		name_frags := partition.Part_name().Get()
		name := strings.Join(name_frags[:],"")
		fmt.Printf("Partition found at %d with name \"%s\"\n",partition.Part_start().Get(),name)
	}
	abs_index4 := a.Partition_disk(1000,&io,"first logic",utiles.B,utiles.Logic,utiles.Best)
	fmt.Printf("Logic partition \"first logic\" at %d\n",abs_index4)
	ebr := types.CreateExtendedBootRecord(&io,abs_index4)
	fmt.Println(ebr.Get())

	l_name := "sec logic"
	abs_index4 = a.Partition_disk(1000,&io,l_name,utiles.B,utiles.Logic,utiles.Best)
	fmt.Printf("Logic partition \"%s\" at %d\n",l_name,abs_index4)
	ebr = types.CreateExtendedBootRecord(&io,abs_index4)
	fmt.Println(ebr.Get())

	l_name = "3 logic"
	abs_index4 = a.Partition_disk(1000,&io,l_name,utiles.B,utiles.Logic,utiles.Best)
	fmt.Printf("Logic partition \"%s\" at %d\n",l_name,abs_index4)
	ebr = types.CreateExtendedBootRecord(&io,abs_index4)
	fmt.Println(ebr.Get())

	fmt.Println("Try to modify 3 logic space ")
	algo:=a.Modify_partition_size_in_disk(-200,&io,l_name,utiles.B)
	fmt.Println(algo)
	result,ext_part:=a.Get_extended_partition(mbr)
	if !result{panic("There was no extended partition")}
	result,ebr = a.Find_logical_partition_by_name(ext_part,utiles.Into_ArrayChar16(l_name))
	if !result{panic("There was no logical partition found")}
	fmt.Println(ebr.Get())

	fmt.Println("Try to modify extended space ")
	algo = a.Modify_partition_size_in_disk(-200,&io,"third extd",utiles.B)
	fmt.Println(algo)
	fmt.Println(ext_part.Get())
	sm:=a.Get_disk_partitions_space_manager(mbr)
	sm.Log_chunks_state()
	sm=a.Get_extended_part_space_manager(ext_part)
	sm.Log_chunks_state()
	io.Flush()
	// mbr := types.CreateMasterBootRecord(&io,0)
}
func recover_and_erase_file_EXT2(){
	io := datamanagment.IOService_from(aplication.DISK_CONTAINER_PATH+"/A.dsk")

	f:=ext2.Recover_FormatEXT2(&io,0,utiles.First)
	fmt.Println("Initizal state")
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	inode_index,inode := f.First_Inode()
	fmt.Print("Inode index = ")
	fmt.Println(inode_index)
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Searching for home directory")
	searched_indx,home_dir_content := f.Search_for_inode(inode,utiles.Into_ArrayChar12("home"))
	home_dir := types.CreateIndexNode(&io,home_dir_content.B_inodo().Get())
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Printf("Is searched dir found succsesfully: %t\n",(searched_indx != -1))
	if searched_indx == -1{panic("Aborting test")}
	fmt.Println("Searching for nested dir user")
	searched_nested_indx,user_dir_content := f.Search_for_inode(home_dir,utiles.Into_ArrayChar12("user"))
	user_dir := types.CreateIndexNode(&io,user_dir_content.B_inodo().Get())
	fmt.Printf("Is searched \"user\" dir found succsesfully : %t\n",searched_nested_indx!=-1)
	if searched_nested_indx == -1{panic("Aborting test")}
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Reading file usuarios.txt in user dir")
	file_inode_index,file_inode := f.Extract_inode(user_dir,utiles.Into_ArrayChar12("usuarios.txt"))
	

	fmt.Printf("Is file \"usuarios.txt\" extracted succsesfully : %t\n",file_inode_index!=-1)
	if file_inode_index == -1{panic("Aborting test")}
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Reading file")
	content:=f.Read_file(&file_inode)
	fmt.Printf("(%s) ... (%s)\n",content[:2],content[len(content)-2:])
	
	fmt.Println("Erasing file content, including blocks")
	f.Update_file(&file_inode,0,[]string{})
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
}
func recover_and_test_dir_and_file_EXT2(){
	io := datamanagment.IOService_from(aplication.DISK_CONTAINER_PATH+"/A.dsk")

	f:=ext2.Recover_FormatEXT2(&io,0,utiles.First)
	fmt.Println("Initizal state")
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	inode_index,inode := f.First_Inode()
	fmt.Print("Inode index = ")
	fmt.Println(inode_index)
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Searching for home directory")
	searched_indx,home_dir_content := f.Search_for_inode(inode,utiles.Into_ArrayChar12("home"))
	home_dir := types.CreateIndexNode(&io,home_dir_content.B_inodo().Get())
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Printf("Is searched dir found succsesfully: %t\n",(searched_indx != -1))
	if searched_indx == -1{panic("Aborting test")}
	fmt.Println("Searching for nested dir user")
	searched_nested_indx,user_dir_content := f.Search_for_inode(home_dir,utiles.Into_ArrayChar12("user"))
	fmt.Printf("Is searched \"user\" dir found succsesfully : %t\n",searched_nested_indx!=-1)
	user_dir := types.CreateIndexNode(&io,user_dir_content.B_inodo().Get())
	if searched_nested_indx == -1{panic("Aborting test")}
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Reading file usuarios.txt in user dir")
	file_inode_index,file_inode_content := f.Search_for_inode(user_dir,utiles.Into_ArrayChar12("usuarios.txt"))
	fmt.Printf("Is searched \"usuarios.txt\" file found succsesfully : %t\n",file_inode_index!=-1)
	file_inode := types.CreateIndexNode(&io,file_inode_content.B_inodo().Get())

	if file_inode_index == -1{panic("Aborting test")}
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Reading file")
	content:=f.Read_file(&file_inode)
	fmt.Printf("(%s) ... (%s)\n",content[:2],content[len(content)-2:])
}
func format_and_test_dir_and_file_EXT2(){
	io := datamanagment.IOService_from(aplication.DISK_CONTAINER_PATH+"/A.dsk")

	f:=ext2.Format_new_FormatEXT2_and_extract(&io,utiles.First,0,1*int32(utiles.Mb))
	fmt.Println("Initizal state")
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	inode_index,inode := f.Create_Inode(types.IndexNodeHolder{
		I_uid:   0,
		I_gid:   0,
		I_s:     0,
		I_atime: utiles.Current_Time(), // change latter
		I_ctime: utiles.Current_Time(),
		I_mtime: utiles.Current_Time(), // change latter
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.Directory),
		I_perm:  [3]string{"a","a","a"}, // change latter
	})
	fmt.Print("Inode index = ")
	fmt.Println(inode_index)
	fmt.Println("Inode Created")
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Creating two directorys in root")
	f.Put_in_dir(inode,types.IndexNodeHolder{
		I_uid:   0,
		I_gid:   0,
		I_s:     0,
		I_atime: utiles.Current_Time(),
		I_ctime: utiles.Current_Time(),
		I_mtime: utiles.Current_Time(),
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.Directory),
		I_perm:  [3]string{"a","a","a"}, // change latter
	},utiles.Into_ArrayChar12("dev"))
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	created_indx,_ := f.Put_in_dir(inode,types.IndexNodeHolder{
		I_uid:   0,
		I_gid:   0,
		I_s:     0,
		I_atime: utiles.Current_Time(),
		I_ctime: utiles.Current_Time(),
		I_mtime: utiles.Current_Time(),
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.Directory),
		I_perm:  [3]string{"a","a","a"}, // change latter
	},utiles.Into_ArrayChar12("home"))
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Searching for home directory")
	searched_indx,home_dir_content := f.Search_for_inode(inode,utiles.Into_ArrayChar12("home"))
	home_dir := types.CreateIndexNode(&io,home_dir_content.B_inodo().Get())
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Printf("Comparing both: searched and created: %t\n",created_indx==searched_indx)
	fmt.Println("Appending nested directory")
	nested_indx,_ := f.Put_in_dir(home_dir,types.IndexNodeHolder{
		I_uid:   0,
		I_gid:   0,
		I_s:     0,
		I_atime: utiles.Current_Time(),
		I_ctime: utiles.Current_Time(),
		I_mtime: utiles.Current_Time(),
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.Directory),
		I_perm:  [3]string{"a","a","a"}, // change latter
		},utiles.Into_ArrayChar12("user"))
	fmt.Println("Searching and comparing created and searched nested dir")
	searched_nested_indx,user_dir_content := f.Search_for_inode(home_dir,utiles.Into_ArrayChar12("user"))
	user_dir := types.CreateIndexNode(&io,user_dir_content.B_inodo().Get())
	fmt.Printf("Result is : %t\n",searched_nested_indx==nested_indx)
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Writing file usuarios.txt in user dir")
	
	_,file := f.Put_in_dir(user_dir,types.IndexNodeHolder{
		I_uid:   0,
		I_gid:   0,
		I_s:     0,
		I_atime: utiles.Current_Time(),
		I_ctime: utiles.Current_Time(),
		I_mtime: utiles.Current_Time(),
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.File),
		I_perm:  [3]string{"a","a","a"}, // change latter
	},utiles.Into_ArrayChar12("usuarios.txt"))
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Updating file")
	const size = 64*13 + 64*16 + 64*16*16 + 64*16*16*16
	data := make([]string,size)
	data[0] = "+"
	for i := 1; i < size-1; i++ {
		data[i] = "."
	}
	data[size-1] = "+"
	f.Update_file(&file,0,data)
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	
	// f.Log_block_bitmap()
	// f.Log_inode_bitmap()
	// fmt.Println("Reading file just created")
	// content:=f.Read_file(&inode)

	// fmt.Printf("(%s) ... (%s)\n",content[:2],content[len(content)-2:])
	// // for _,ch := range f.Read_file(&inode){
	// // 	fmt.Print(ch)
	// // }
	// // fmt.Println()
	io.Flush()
}
func format_and_write_file_EXT2(){
	io := datamanagment.IOService_from(aplication.DISK_CONTAINER_PATH+"/A.dsk")
	f:=ext2.Format_new_FormatEXT2_and_extract(&io,utiles.First,0,1*int32(utiles.Mb))
	fmt.Println("Initizal state")
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	file_inode := types.IndexNodeHolder{
		I_uid:   0,
		I_gid:   0,
		I_s:     0,
		I_atime: utiles.Current_Time(), // change latter
		I_ctime: utiles.Current_Time(),
		I_mtime: utiles.Current_Time(), // change latter
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.File),
		I_perm:  [3]string{"a","a","a"}, // change latter
	}
	inode_index,inode := f.Create_Inode(file_inode)
	fmt.Print("Inode index = ")
	fmt.Println(inode_index)
	fmt.Println("Inode Created")
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Updating file")
	const size = 64*13 + 64*16 + 64*16*16 + 64*16*16*16
	data := make([]string,size)
	data[0] = "+"
	for i := 1; i < size-1; i++ {
		data[i] = "."
	}
	data[size-1] = "+"
	f.Update_file(&inode,0,data)
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Reading file just created")
	content:=f.Read_file(&inode)

	fmt.Printf("(%s) ... (%s)\n",content[:2],content[len(content)-2:])
	// for _,ch := range f.Read_file(&inode){
	// 	fmt.Print(ch)
	// }
	// fmt.Println()
	io.Flush()
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
	sm = datamanagment.SpaceManager_from_occuped_spaces([]datamanagment.Space{
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