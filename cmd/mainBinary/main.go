package main

import (
	"fmt"
	"project/internal/datamanagment"
	"project/internal/formats/ext2"
	"project/internal/types"
	"project/internal/utiles"
)

func main() {
	// fmt.Println(int32(15)/64)
	// simulate_reading()
	// simulate_writing()
	// bitmap_firstfit()
	// bitmap_worstfit()
	// bitmap_bestfit()
	// recover_and_read_format_EXT2()

	// format_and_write_file_EXT2()
	// first := [12]string{"1","1","1","1","1","1","1","1","1","1","1","",}
	// second := [12]string{"1","1","1","1","1","1","1","1","1","1","1","",}
	// fmt.Println(first == second)
	// format_and_test_dir_and_file_EXT2()
	// recover_and_test_dir_and_file_EXT2()
	// recover_and_erase_file_EXT2()
	bestfit()
	// a:=utiles.New_Space(5,5)
	// b:=utiles.New_Space(0,5)
	// fmt.Printf("%s %s with %s\n",a.Show(),a.Contains(b),b.Show())
	// b=utiles.New_Space(10,5)
	// fmt.Printf("%s %s with %s\n",a.Show(),a.Contains(b),b.Show())
	// b=utiles.New_Space(5,5)
	// fmt.Printf("%s %s with %s\n",a.Show(),a.Contains(b),b.Show())
	// b=utiles.New_Space(5,4)
	// fmt.Printf("%s %s with %s\n",a.Show(),a.Contains(b),b.Show())
	// b=utiles.New_Space(9,1)
	// fmt.Printf("%s %s with %s\n",a.Show(),a.Contains(b),b.Show())
	// b=utiles.New_Space(6,1)
	// fmt.Printf("%s %s with %s\n",a.Show(),a.Contains(b),b.Show())
	// b=utiles.New_Space(6,10)
	// fmt.Printf("%s %s with %s\n",a.Show(),a.Contains(b),b.Show())
	// b=utiles.New_Space(4,10)
	// fmt.Printf("%s %s with %s\n",a.Show(),a.Contains(b),b.Show())
	// b=utiles.New_Space(4,5)
	// fmt.Printf("%s %s with %s\n",a.Show(),a.Contains(b),b.Show())
	
	// c:=utiles.New_Space(10,9)
	// sd:=utiles.New_Space(13,3)
	// fmt.Printf("%s %s with %s\n",c.Show(),c.Contains(sd),sd.Show())
	
	// fmt.Printf("split %s with %s\n",c.Show(),sd.Show())
	// a1,a2:=c.Split(sd)
	// fmt.Printf("Result in two peaces %s AND %s\n",a1.Show(),a2.Show())



}
func check(){

}
func recover_and_read_format_EXT2(){
	io := datamanagment.IOService_from("./disco.txt")
	f:=ext2.Recover_FormatEXT2(&io,0,utiles.First)
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	inode := types.CreateIndexNode(&io,160)
	fmt.Println("Reading file just created")
	for _,ch := range f.Read_file(&inode){
		fmt.Print(ch)
	}
	fmt.Println()
}
func recover_and_erase_file_EXT2(){
	io := datamanagment.IOService_from("./disco.txt")
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
	searched_indx,home_dir := f.Search_for_inode(inode,utiles.Into_ArrayChar12("home"))
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Printf("Is searched dir found succsesfully: %t\n",(searched_indx != -1))
	if searched_indx == -1{panic("Aborting test")}
	fmt.Println("Searching for nested dir user")
	searched_nested_indx,user_dir := f.Search_for_inode(home_dir,utiles.Into_ArrayChar12("user"))
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
	io := datamanagment.IOService_from("./disco.txt")
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
	searched_indx,home_dir := f.Search_for_inode(inode,utiles.Into_ArrayChar12("home"))
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Printf("Is searched dir found succsesfully: %t\n",(searched_indx != -1))
	if searched_indx == -1{panic("Aborting test")}
	fmt.Println("Searching for nested dir user")
	searched_nested_indx,user_dir := f.Search_for_inode(home_dir,utiles.Into_ArrayChar12("user"))
	fmt.Printf("Is searched \"user\" dir found succsesfully : %t\n",searched_nested_indx!=-1)
	if searched_nested_indx == -1{panic("Aborting test")}
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Reading file usuarios.txt in user dir")
	file_inode_index,file_inode := f.Search_for_inode(user_dir,utiles.Into_ArrayChar12("usuarios.txt"))
	fmt.Printf("Is searched \"usuarios.txt\" file found succsesfully : %t\n",file_inode_index!=-1)
	if file_inode_index == -1{panic("Aborting test")}
	f.Log_block_bitmap()
	f.Log_inode_bitmap()
	fmt.Println("Reading file")
	content:=f.Read_file(&file_inode)
	fmt.Printf("(%s) ... (%s)\n",content[:2],content[len(content)-2:])
}
func format_and_test_dir_and_file_EXT2(){
	io := datamanagment.IOService_from("./disco.txt")
	f:=ext2.Format_new_FormatEXT2(&io,utiles.First,0,1_000_000)
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
	searched_indx,home_dir := f.Search_for_inode(inode,utiles.Into_ArrayChar12("home"))
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
	searched_nested_indx,user_dir := f.Search_for_inode(home_dir,utiles.Into_ArrayChar12("user"))
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
	io := datamanagment.IOService_from("./disco.txt")
	f:=ext2.Format_new_FormatEXT2(&io,utiles.First,0,1_000_000)
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
	const size = 64*13 + 64*16 + 64*16*16 + 64*16*16*16 + 1
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
	sm := utiles.Spacemanager_from_free_spaces([]utiles.Space{
		utiles.New_Space(0,3),utiles.New_Space(5,2),utiles.New_Space(10,10) },20)
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
	sm = utiles.Spacemanager_from_occuped_spaces([]utiles.Space{
		utiles.New_Space(3,2),utiles.New_Space(7,3),
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
	fmt.Println("free last space")
	// sm.Free_space(1,19)
	sm.Log_chunks_state()
	// bitmap.Init_mapping()
	// bitmap.Log_chunks_state()
	// bitmap.Log_bitmap_state()
	// fmt.Print("fittin for 2 in index = ")
	// fmt.Println(bitmap.Best_fit(2))
	// bitmap.Log_chunks_state()
	// bitmap.Log_bitmap_state()
	// fmt.Println("Erase 1 at 0")
	// bitmap.Erase(1,0)
	// bitmap.Log_chunks_state()
	// bitmap.Log_bitmap_state()
	// fmt.Print("fittin for 1 in index = ")
	// fmt.Println(bitmap.Best_fit(1))
	// bitmap.Log_chunks_state()
	// bitmap.Log_bitmap_state()
}
// func bitmap_worstfit(){
// 	io := datamanagment.IOService_from("./disco.txt")
// 	bitmap := ext2.New_Bitmap(&io,0,15,9)
// 	bitmap.Init_mapping()
// 	bitmap.Log_chunks_state()
// 	bitmap.Log_bitmap_state()
// 	fmt.Print("fittin for 2 in index = ")
// 	fmt.Println(bitmap.Worst_fit(2))
// 	bitmap.Log_chunks_state()
// 	bitmap.Log_bitmap_state()
// // 	fmt.Println("Erase 1 at 2")
// // 	bitmap.Erase(1,2)
// // 	bitmap.Log_chunks_state()
// // 	bitmap.Log_bitmap_state()
// // 	fmt.Print("fittin for 1 in index = ")
// // 	fmt.Println(bitmap.Worst_fit(1))
// // 	bitmap.Log_chunks_state()
// // 	bitmap.Log_bitmap_state()
// }
// func bitmap_firstfit(){
// 	io := datamanagment.IOService_from("./disco.txt")
// 	bitmap := ext2.New_Bitmap(&io,0,10,9)
// 	bitmap.Init_mapping()
// 	bitmap.Log_chunks_state()
// 	bitmap.Log_bitmap_state()
// 	fmt.Println(bitmap.First_fit(4))
// 	bitmap.Log_chunks_state()
// 	bitmap.Log_bitmap_state()
// 	fmt.Println("Erase 1 at 2")
// 	bitmap.Erase(1,2)
// 	bitmap.Log_chunks_state()
// 	bitmap.Log_bitmap_state()
// 	fmt.Println("Erase 1 at 5")
// 	bitmap.Erase(1,5)
// 	bitmap.Log_chunks_state()
// 	bitmap.Log_bitmap_state()
// 	fmt.Println("Fit of 3")
// 	fmt.Println(bitmap.First_fit(3))
// 	bitmap.Log_chunks_state()
// 	bitmap.Log_bitmap_state()
// 	fmt.Println("Erase 1 at 6")
// 	bitmap.Erase(1,6)
// 	bitmap.Log_chunks_state()
// 	bitmap.Log_bitmap_state()
// 	fmt.Println("Fitting for 2")
// 	fmt.Println(bitmap.First_fit(2))
// 	bitmap.Log_chunks_state()
// 	bitmap.Log_bitmap_state()
// }



func simulate_writing() {
	io := datamanagment.IOService_from("./disco.txt")
	mbr := types.CreateMasterBootRecord(&io, 0)
	mbr.Mbr_tamano().Set(1034)
	partition1 := mbr.Mbr_partitions().No(0)
	partition1.Part_status().Set("Y")
	partition1.Part_name().Set(utiles.Into_ArrayChar16("1234567890******************"))
	io.Flush()
}
func simulate_reading() {
	io := datamanagment.IOService_from("./disco.txt")
	mbr := types.CreateMasterBootRecord(&io, 0)
	fmt.Println("This is the partition name")
	fmt.Print("\"")
	for _, c := range mbr.Mbr_partitions().No(0).Part_name().Spread() {
		fmt.Print(c.Get())
	}
	fmt.Print("\"")
	fmt.Println()
}
