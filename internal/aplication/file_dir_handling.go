package aplication

import (
	// "fmt"
	"fmt"
	"project/internal/datamanagment"
	"project/internal/formats"
	"project/internal/types"
	"project/internal/utiles"
	"strings"
)

func (self *Aplication) Recover_EXT2_Format(super_service *datamanagment.IOService) (int32,utiles.FitCriteria,error) {
	var super_block_index int32 = -1
	var fit utiles.FitCriteria
	var err error
	switch self.active_partition.part_type {
	case utiles.Logic:
		ebr := types.CreateExtendedBootRecord(super_service, self.active_partition.index)
		super_block_index = ebr.Part_start().Get()
		fit,err = utiles.Translate_fit(ebr.Part_fit().Get())
		if err!=nil{return -1,utiles.Best,err}

	case utiles.Primary:
		partition := types.CreatePartition(super_service, self.active_partition.index)
		super_block_index = partition.Part_start().Get()
		fit,err = utiles.Translate_fit(partition.Part_fit().Get())
		if err!=nil{return -1,utiles.Best,err}
	}
	return super_block_index,fit,nil
}



func (self *Aplication) Make_file(folders [][12]string, data []string, file_name [12]string, create_recursive bool) (*formats.JournalingManager,error) {
	if self.active_partition == nil {return nil,fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	current_time := utiles.Current_Time()
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err!=nil{return nil,err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()

	var dir types.IndexNode
	dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)

	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,dir := format.Get_nested_dir(dir,folders,create_recursive,user_corr,user_g_corr,current_time,false,true)
	if err0 != nil{return nil,err0}


	file:=format.Put_in_dir(dir,format.Wrap_holder_in_template(types.IndexNodeHolder{
		I_uid:   user_corr,
		I_gid:   user_g_corr,
		I_s:     int32(len(data)),
		I_atime: utiles.NO_TIME,
		I_ctime: current_time,
		I_mtime: current_time,
		I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  string(utiles.File),
		I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
	}),file_name)
	if file.Index == -1 {
	}
	// old_size:=file.I_s().Get()
	format.Update_file(&file,0,data)
	format.Update_dir_and_ancestors_size(dir,int32(len(data)))

	dir.I_mtime().Set(current_time)
	return format.Get_journaling(),nil
}

func (self *Aplication) Show_file(folders [][12]string, file_name [12]string) (string,error) {
	if self.active_partition == nil {return "",fmt.Errorf("there's no active partition")}
	if !self.active_partition.has_session {return "",fmt.Errorf("there's no session yet")}
	if self.active_partition.session.active_user == nil {return "",fmt.Errorf("there's no user logged")}
	super_service := self.active_partition.io
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err != nil {return "",err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	current_time := utiles.Current_Time()
	dir := format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,dir := format.Get_nested_dir(dir,folders,false,user_corr,user_g_corr,current_time,true,false)
	if err0 != nil{return "",err0}
	file_r,content:=format.Search_for_inode(dir,file_name)
	if file_r == -1 {return "",fmt.Errorf("File not found")}
	file:=types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
	perm:=format.User_allowed_actions(user_corr,user_g_corr,&file)
	if !perm.Can_read(){return "",fmt.Errorf("user is not allowed to read this file")}
	r:=format.Read_file(&file)
	return strings.Join(r,""),nil
}


func (self *Aplication) Remove(folders [][12]string, with_name [12]string) (*formats.JournalingManager,error) {
	if self.active_partition == nil {return nil,fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err != nil{return nil,err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	current_time := utiles.Current_Time()
	dir := format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,dir := format.Get_nested_dir(dir,folders,false,user_corr,user_g_corr,current_time,true,false)
	if err0 != nil{return nil,err0}

	remove_result := format.Remove_inode_if_possilbe(dir,with_name,user_corr,user_g_corr,current_time)
	if !remove_result{return nil,fmt.Errorf("Dir was not removed")}
	return format.Get_journaling(),nil
}


func (self *Aplication) Edit_file(folders [][12]string, data []string, file_name [12]string) (*formats.JournalingManager,error) {
	if self.active_partition == nil {return nil,fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	current_time := utiles.Current_Time()
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err!=nil{return nil,err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	dir := format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,dir := format.Get_nested_dir(dir,folders,false,user_corr,user_g_corr,current_time,true,true)
	if err0 != nil{return nil,err0}
	file_result,content:=format.Search_for_inode(dir,file_name)
	if file_result == -1 {
		return nil,fmt.Errorf("File not found")
	}
	file:=types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
	old_size:=file.I_s().Get()
	format.Update_file(&file,0,data)
	dir.I_mtime().Set(current_time)
	format.Update_dir_and_ancestors_size(dir,int32(len(data))-old_size)

	return format.Get_journaling(),nil
}

func (self *Aplication)Rename_inode(folders [][12]string, trgt_name [12]string,new_name [12]string) (*formats.JournalingManager,error) {
	if self.active_partition == nil {return nil,fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	current_time := utiles.Current_Time()
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err!=nil{return nil,err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	dir := format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,dir := format.Get_nested_dir(dir,folders,false,user_corr,user_g_corr,current_time,false,true)
	if err0 != nil{return nil,err0}
	file_result,content:=format.Search_for_inode(dir,trgt_name)
	if file_result == -1 {
		return nil,fmt.Errorf("File with that name doesnt exist")
	}
	file_result2,_:=format.Search_for_inode(dir,new_name)
	if file_result2 != -1 {
		return nil,fmt.Errorf("File with that name already exists")
	}
	content.B_name().Set(new_name)
	dir.I_mtime().Set(current_time)
	return format.Get_journaling(),nil
}



func (self *Aplication) Make_dir(folders [][12]string, dir_name [12]string, create_recursive bool) (*formats.JournalingManager,error) {
	if self.active_partition == nil {return nil,fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	current_time := utiles.Current_Time()
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err!=nil{return nil,err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	var dir types.IndexNode
	dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,dir := format.Get_nested_dir(dir,folders,create_recursive,user_corr,user_g_corr,current_time,false,true)
	if err0 != nil{return nil,err0}

	dir_appended:=format.Put_in_dir(dir,format.Wrap_holder_in_template(types.IndexNodeHolder{
		I_uid:   user_corr,
		I_gid:   user_g_corr,
		I_s:     0,
		I_atime: utiles.NO_TIME,
		I_ctime: current_time,
		I_mtime: current_time,
		I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
		I_type:  string(utiles.Directory),
		I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
	}),dir_name)
	if dir_appended.Index == -1 {
		return nil,fmt.Errorf("Dir coulndt be created")
	}
	format.Set_parent_child_relation(dir,dir_appended)
	dir.I_mtime().Set(current_time)
	return format.Get_journaling(),nil
}















func (self *Aplication) Copy(folders_trgt [][12]string, for_name [12]string, folders_dest [][12]string) (*formats.JournalingManager,error) {
	if self.active_partition == nil {return nil,fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	current_time := utiles.Current_Time()
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err!=nil{return nil,err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	var dir types.IndexNode
	dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	
	err0,srcs_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,true,false)
	if err0 != nil{return nil,err0}
	search_r,trgt_dir_content := format.Search_for_inode(srcs_dir,for_name)
	if search_r == -1 {return nil,fmt.Errorf("")}
	trgt_dir:=types.CreateIndexNode(trgt_dir_content.Super_service,trgt_dir_content.B_inodo().Get())

	err1,dest_dir := format.Get_nested_dir(dir,folders_dest,false,user_corr,user_g_corr,current_time,false,true)
	if err1 != nil{return nil,err1}
	res:=format.Directory_deep_copy(dest_dir,trgt_dir,for_name,user_corr,user_g_corr,current_time)
	format.Update_dir_and_ancestors_size(dest_dir,trgt_dir.I_s().Get())
	if res{return format.Get_journaling(),nil}
	return nil,fmt.Errorf("No se pudo realizar la copia")
}


func (self *Aplication) Move(folders_trgt [][12]string, for_name [12]string, folders_dest [][12]string) (*formats.JournalingManager,error) {
	if self.active_partition == nil {return nil,fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	current_time := utiles.Current_Time()
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err!=nil{return nil,err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	var dir types.IndexNode
	dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,srcs_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,true,false)
	if err0 != nil{return nil,err0}
	search_r,_ := format.Extract_inode(srcs_dir,for_name)
	if search_r == -1 {return nil,fmt.Errorf("")}
	err1,dest_dir := format.Get_nested_dir(dir,folders_dest,false,user_corr,user_g_corr,current_time,false,true)
	if err1 != nil{return nil,err1}
	appended:= format.Put_in_dir(dest_dir,format.Wrap_indx_in_template(search_r),for_name) 
	if appended.Index !=-1{
		if appended.I_type().Get() == string(utiles.Directory){
			format.Set_parent_child_relation(dest_dir,appended)
		}
		indoe_size:=appended.I_s().Get()
		if indoe_size != 0{
			format.Update_dir_and_ancestors_size(srcs_dir,-indoe_size)
			format.Update_dir_and_ancestors_size(dest_dir,indoe_size)
		}

		return format.Get_journaling(),nil
	}
	return nil,fmt.Errorf("no se puedo mover el fichero")
}


type StrFile struct{
	name string
	childs []*StrFile
	logger string
}
func (self *StrFile) print_as_root()string{
	self.logger += self.name + "\n"
	if len(self.childs) == 0{
		return ""
	}
	for i := 0; i < len(self.childs)-1; i++ {
		self.childs[i].print_as_branch("|  ")
	}
	self.childs[len(self.childs)-1].print_as_branch("   ")
	return self.logger
};
func (self *StrFile) print_as_branch(history string){
	self.logger += "|__"
	self.logger += self.name + "\n"
	if len(self.childs) == 0{
		return
	}
	for i := 0; i < len(self.childs)-1; i++ {
		self.logger += history
		self.childs[i].print_as_branch(history+"|  ")
	}
	self.logger += history
	self.childs[len(self.childs)-1].print_as_branch(history+"   ")
};
func (self *Aplication) Find(folders_trgt [][12]string, name_criteria utiles.NameCriteria) (string,error) {
	if self.active_partition == nil {return "",fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	current_time := utiles.Current_Time()
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err!=nil{return "",err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	var dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,srcs_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,true,false)
	if err0 != nil{return "",err0}

	to_plain_name := func (str [12]string)string{
		final:=""
		for _, c := range str {final += c}
		return strings.TrimSpace(final)
	}
	
	var print_element func(types.Content)*StrFile
	print_element = func (content types.Content)*StrFile{
		inode:=types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
		raw_name :=content.B_name().Get()
		switch inode.I_type().Get() {
		case string(utiles.File):
			if name_criteria.Match(raw_name){
				return &StrFile{name: to_plain_name(raw_name),childs: []*StrFile{}}
			}
		case string(utiles.Directory):
			all_childs:=format.Get_strict_shallow_tree_of_childs(inode)
			childs := make([]*StrFile,0,5)
			for _, child := range all_childs {
				result := print_element(child)
				if result == nil{continue}
				childs = append(childs, result)
			}
			if len(childs) == 0{
				if name_criteria.Match(raw_name){
					return &StrFile{name: to_plain_name(raw_name),childs: []*StrFile{}}
				}
				return nil
			}
			return &StrFile{name: to_plain_name(raw_name),childs: childs}	
		}
		return nil
	}
	all_childs:=format.Get_strict_shallow_tree_of_childs(srcs_dir)
	childs := make([]*StrFile,0,5)
	for _, child := range all_childs {
		result := print_element(child)
		if result == nil{continue}
		childs = append(childs, result)
	}
	if len(childs) == 0{return "",nil}
	tree := StrFile{name: ".",childs: childs}
	
	return tree.print_as_root(),nil
}



func (self *Aplication) Change_own(folders_trgt [][12]string, for_name [12]string, user string,recursive bool) (*formats.JournalingManager,error) {
	if self.active_partition == nil {return nil,fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	current_time := utiles.Current_Time()
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err!=nil{return nil,err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	// format.Init_bitmap_mapping()
	var dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,srcs_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,true,false)
	if err0 != nil{return nil,err0}
	search_r,trgt_inode_content := format.Search_for_inode(srcs_dir,for_name)
	if search_r == -1 {return nil,fmt.Errorf("")}
	trgt_inode:=types.CreateIndexNode(trgt_inode_content.Super_service,trgt_inode_content.B_inodo().Get())

	new_user,err:=self.active_partition.session.Get_User(user)
	if err!= nil {return nil,err}
	new_user_corr:= int32(new_user.correlative_number)
	new_user_g_corr :=  int32(self.active_partition.session.Get_user_group(new_user).correlative_number)
	changed:=format.Change_owner(trgt_inode,user_corr,user_g_corr,current_time,new_user_corr,new_user_g_corr,recursive)
	if changed != 0 {
		return format.Get_journaling(),nil
	}
	return nil,fmt.Errorf("No se han realizado cambios")
}


func (self *Aplication) Chagne_UGO(folders_trgt [][12]string, for_name [12]string, perm string,recursive bool)(*formats.JournalingManager,error){
	if self.active_partition == nil {return nil,fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	current_time := utiles.Current_Time()
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err!=nil{return nil,err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	// format.Init_bitmap_mapping()
	var dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,srcs_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,true,false)
	if err0 != nil{return nil,err0}
	search_r,trgt_inode_content := format.Search_for_inode(srcs_dir,for_name)
	if search_r == -1 {return nil,fmt.Errorf("")}
	trgt_inode:=types.CreateIndexNode(trgt_inode_content.Super_service,trgt_inode_content.B_inodo().Get())
	ugo,err:=utiles.Parse_to_ugo(perm)
	if err!= nil {return nil,err}
	new_perm := ugo.To_arr_string()
	changed:=format.Change_ugo_permition(trgt_inode,current_time,recursive,new_perm)
	if changed != 0 {
		return format.Get_journaling(),nil
	}
	return nil,fmt.Errorf("No se han realizado cambios")

}

