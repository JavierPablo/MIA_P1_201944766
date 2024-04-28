package aplication

import (
	"fmt"
	"project/internal/formats"
	"project/internal/types"
	"project/internal/utiles"
)
func (self *Aplication) All_shallow_childs_for(folders_trgt [][12]string) ([]*types.Content,error) {
	if self.active_partition == nil {return []*types.Content{},fmt.Errorf("there's no active partition")}
	super_service := self.active_partition.io
	current_time := utiles.Current_Time()
	super_block_index,fit,err := self.Recover_EXT2_Format(super_service)
	if err!=nil{return []*types.Content{},err}
	format := formats.Recover_Format(super_service, super_block_index, fit)
	format.Init_bitmap_mapping()
	var dir = format.First_Inode()
	user_corr:= int32(self.active_partition.session.active_user.correlative_number)
	user_g_corr :=  int32(self.active_partition.session.Get_user_group(self.active_partition.session.active_user).correlative_number)
	err0,srcs_dir := format.Get_nested_dir(dir,folders_trgt,false,user_corr,user_g_corr,current_time,true,false)
	if err0 != nil{return []*types.Content{},err0}
	

	all_childs:=format.Get_strict_shallow_tree_of_childs(srcs_dir)
	filtered_childs := make([]*types.Content,0,5)
	for _, child := range all_childs {
		child_inode:=types.CreateIndexNode(child.Super_service,child.B_inodo().Get())
		permisions := format.User_allowed_actions(user_corr,user_g_corr,&child_inode)
		if !permisions.Can_read(){continue}
		filtered_childs = append(filtered_childs, &child)
	}
	return filtered_childs,nil
}