package aplication

import (
	// "fmt"
	"project/internal/datamanagment"
	"project/internal/formats/ext2"
	"project/internal/types"
	"project/internal/utiles"
	"strings"
)

func (self *Aplication) Recover_EXT2_Format(super_service *datamanagment.IOService) (int32,utiles.FitCriteria) {
	var super_block_index int32 = -1
	var fit utiles.FitCriteria
	switch self.active_partition.part_type {
	case utiles.Logic:
		ebr := types.CreateExtendedBootRecord(super_service, self.active_partition.index)
		super_block_index = ebr.Part_start().Get()
		fit = utiles.Translate_fit(ebr.Part_fit().Get())
	case utiles.Primary:
		partition := types.CreatePartition(super_service, self.active_partition.index)
		super_block_index = partition.Part_start().Get()
		fit = utiles.Translate_fit(partition.Part_fit().Get())
	}
	return super_block_index,fit
}



func (self *Aplication) Make_file(super_service *datamanagment.IOService, folders [][12]string, data []string, file_name [12]string, create_recursive bool) bool {
	current_time := utiles.Current_Time()
	super_block_index,fit := self.Recover_EXT2_Format(super_service)
	format := ext2.Recover_FormatEXT2(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	var dir types.IndexNode
	_, dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	result,dir := format.Get_nested_dir(dir,folders,create_recursive,user_corr,user_g_corr,current_time,false,true)
	if !result{return false}
	file_result,file:=format.Put_in_dir(dir,types.IndexNodeHolder{
		I_uid:   user_corr,
		I_gid:   user_g_corr,
		I_s:     int32(len(data)),
		I_atime: utiles.NO_TIME,
		I_ctime: current_time,
		I_mtime: current_time,
		I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  string(utiles.File),
		I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
	},file_name)
	if file_result == -1 {
		return false
	}
	format.Update_file(&file,0,data)
	dir.I_mtime().Set(current_time)
	return true
}

func (self *Aplication) Show_file(super_service *datamanagment.IOService, folders [][12]string, file_name [12]string) (bool,string) {
	super_block_index,fit := self.Recover_EXT2_Format(super_service)
	format := ext2.Recover_FormatEXT2(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	current_time := utiles.Current_Time()
	_, dir := format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	result,dir := format.Get_nested_dir(dir,folders,false,user_corr,user_g_corr,current_time,true,false)
	if !result{return false,""}
	file_r,content:=format.Search_for_inode(dir,file_name)
	if file_r == -1 {return false,""}
	file:=types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
	r:=format.Read_file(&file)
	return true,strings.Join(r,"")
}


func (self *Aplication) Remove(super_service *datamanagment.IOService, folders [][12]string, with_name [12]string) (bool) {
	super_block_index,fit := self.Recover_EXT2_Format(super_service)
	format := ext2.Recover_FormatEXT2(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	current_time := utiles.Current_Time()
	_, dir := format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	result,dir := format.Get_nested_dir(dir,folders,false,user_corr,user_g_corr,current_time,true,false)
	if !result {return false}
	remove_result := format.Remove_dir_if_possilbe(dir,with_name,user_corr,user_g_corr,current_time)
	return remove_result
}


func (self *Aplication) Edit_file(super_service *datamanagment.IOService, folders [][12]string, data []string, file_name [12]string) bool {
	current_time := utiles.Current_Time()
	super_block_index,fit := self.Recover_EXT2_Format(super_service)
	format := ext2.Recover_FormatEXT2(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	_, dir := format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	result,dir := format.Get_nested_dir(dir,folders,false,user_corr,user_g_corr,current_time,true,true)
	if !result{return false}
	file_result,content:=format.Search_for_inode(dir,file_name)
	if file_result == -1 {
		return false
	}
	file:=types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
	format.Update_file(&file,0,data)
	dir.I_mtime().Set(current_time)
	return true
}

func (self *Aplication)Rename_inode(super_service *datamanagment.IOService, folders [][12]string, trgt_name [12]string,new_name [12]string) bool {
	current_time := utiles.Current_Time()
	super_block_index,fit := self.Recover_EXT2_Format(super_service)
	format := ext2.Recover_FormatEXT2(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	_, dir := format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	result,dir := format.Get_nested_dir(dir,folders,false,user_corr,user_g_corr,current_time,false,true)
	if !result{return false}
	file_result,content:=format.Search_for_inode(dir,trgt_name)
	if file_result == -1 {
		return false
	}
	file_result2,_:=format.Search_for_inode(dir,new_name)
	if file_result2 != -1 {
		return false
	}
	content.B_name().Set(new_name)
	dir.I_mtime().Set(current_time)
	return true
}



func (self *Aplication) Make_dir(super_service *datamanagment.IOService, folders [][12]string, dir_name [12]string, create_recursive bool) bool {
	current_time := utiles.Current_Time()
	super_block_index,fit := self.Recover_EXT2_Format(super_service)
	format := ext2.Recover_FormatEXT2(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	var dir types.IndexNode
	_, dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	result,dir := format.Get_nested_dir(dir,folders,create_recursive,user_corr,user_g_corr,current_time,false,true)
	if !result{return false}

	dir_result,_:=format.Put_in_dir(dir,types.IndexNodeHolder{
		I_uid:   user_corr,
		I_gid:   user_g_corr,
		I_s:     0,
		I_atime: utiles.NO_TIME,
		I_ctime: current_time,
		I_mtime: current_time,
		I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  string(utiles.Directory),
		I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
	},dir_name)
	if dir_result == -1 {
		return false
	}
	dir.I_mtime().Set(current_time)
	return true
}















func (self *Aplication) Copy(super_service *datamanagment.IOService, folders_trgt [][12]string, for_name [12]string, folders_dest [][12]string) bool {
	current_time := utiles.Current_Time()
	super_block_index,fit := self.Recover_EXT2_Format(super_service)
	format := ext2.Recover_FormatEXT2(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	var dir types.IndexNode
	_, dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	result,srcs_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,true,false)
	if !result {return false}
	search_r,trgt_dir := format.Search_for_inode(srcs_dir,for_name)
	if search_r == -1 {return false}
	result,dest_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,false,true)
	if !result{return false}
	return format.Directory_deep_copy(dest_dir,types.IndexNode(trgt_dir),for_name,user_corr,user_g_corr,current_time)
}


func (self *Aplication) Move(super_service *datamanagment.IOService, folders_trgt [][12]string, for_name [12]string, folders_dest [][12]string) bool {
	current_time := utiles.Current_Time()
	super_block_index,fit := self.Recover_EXT2_Format(super_service)
	format := ext2.Recover_FormatEXT2(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	var dir types.IndexNode
	_, dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	result,srcs_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,true,false)
	if !result {return false}
	search_r,_ := format.Extract_inode(srcs_dir,for_name)
	if search_r == -1 {return false}
	result,dest_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,false,true)
	if !result{return false}
	return format.Put_indx_in_dir(dest_dir,search_r,for_name) != -1
}




// func (self *Aplication) Find(super_service *datamanagment.IOService, folders_trgt [][12]string, name_criteria utiles.NameCriteria) bool {
// 	current_time := utiles.Current_Time()
// 	super_block_index,fit := self.Recover_EXT2_Format(super_service)
// 	format := ext2.Recover_FormatEXT2(super_service, super_block_index, fit)
// 	format.Init_bitmap_mapping()
// 	var dir types.IndexNode
// 	_, dir = format.First_Inode()
// 	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
// 	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
// 	result,srcs_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,true,false)
// 	if !result {return false}
// 	all_content := format.Find_in_dir(srcs_dir,name_criteria,user_corr,user_g_corr,current_time)
// 	de
// 	for _,content := range all_content{
// 		inode := types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
// 		switch inode.I_type().Get() {
// 		case string(utiles.File):
// 			fmt.Printf("\n")
	
// 		case string(utiles.Directory):
		
// 			if result == -1 {return false}
// 			sub_folders := self.Get_shallow_tree(trgt_dir)
// 			for _,folder:=range sub_folders{
// 				inode:=types.CreateIndexNode(folder.Super_service,folder.B_inodo().Get())
// 				self.Directory_deep_copy(new_inode,inode,folder.B_name().Get(),usr_id,usr_grp_id,time)
// 			}
// 		default:panic("Didnt match neither file nor directory")
// 		}
// 	}
	
// }



