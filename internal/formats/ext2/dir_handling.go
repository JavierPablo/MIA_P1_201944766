package ext2

import (
	// "fmt"
	"project/internal/types"
	"project/internal/utiles"
)

func (self *FormatEXT2)search_in_ptr_block(pointer_block *types.PointerBlock, level int32, name [12]string) (int32,types.Content) {
	if level != 0 {
		for _,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_index,trgt_block := self.search_in_ptr_block(&ptr_block,level-1,name)
			if trgt_index != -1 {
				return trgt_index, trgt_block
			} 
		}
		return -1,types.CreateContent(self.super_service,0)
	}
	for _,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
		for b, content := range dir_block.B_content().Get(){
			if content.B_inodo == -1 {continue}
			if name == content.B_name {
				return content.B_inodo,dir_block.B_content().No(int32(b))
			}
		}
	}
	return -1,types.CreateContent(self.super_service,0)

}
func (self *FormatEXT2)Search_for_inode(dir types.IndexNode, name [12]string) (int32,types.Content) {
	for n,ptr := range dir.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_index,trgt_block := self.search_in_ptr_block(&ptr_block,int32(n-13),name)
			if trgt_index != -1 {
				return trgt_index, trgt_block
			} 
		}else{
			dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
			for d, content := range dir_block.B_content().Get(){
				if content.B_inodo == -1 {continue}
				if name == content.B_name {
					return content.B_inodo,dir_block.B_content().No(int32(d))
				}
			}
		}
	}
	return -1,types.CreateContent(self.super_service,0)
}










func (self *FormatEXT2)extract_inode_in_ptr_block(pointer_block *types.PointerBlock, level int32, name [12]string) (int32,types.IndexNode) {
	if level != 0 {
		for i,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_index,trgt_block := self.extract_inode_in_ptr_block(&ptr_block,level-1,name)
			if trgt_index != -1 {
				must_erase := true
				for _,sub_ptr := range ptr_block.B_pointers().Get(){
					if sub_ptr != -1{must_erase = false;break}
				}
				if must_erase{
					self.Block_bitmap.Erase(1,self.Block_section.Bit_no_for(ptr))
					pointer_block.B_pointers().No(int32(i)).Set(-1)
				}

				return trgt_index, trgt_block
			} 
		}
		return -1,types.CreateIndexNode(self.super_service,0)
	}
	for n,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
		for i, content := range dir_block.B_content().Get(){
			if content.B_inodo == -1 {continue}
			if name == content.B_name {
				inode_result := types.CreateIndexNode(self.super_service,content.B_inodo)
				
				dir_block.B_content().No(int32(i)).B_inodo().Set(-1)
				must_erase := true
				for _,cntnt := range dir_block.B_content().Get(){
					if cntnt.B_inodo != -1 {must_erase = false;break}
				}
				if must_erase{
					self.Block_bitmap.Erase(1,self.Block_section.Bit_no_for(ptr))
					pointer_block.B_pointers().No(int32(n)).Set(-1)
				}
			
				return content.B_inodo,inode_result
			}
		}
	}
	return -1,types.CreateIndexNode(self.super_service,0)

}
func (self *FormatEXT2)Extract_inode(dir types.IndexNode, name [12]string) (int32,types.IndexNode) {
	for n,ptr := range dir.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_index,trgt_block := self.extract_inode_in_ptr_block(&ptr_block,int32(n-13),name)
			if trgt_index != -1 {
				must_erase := true
				for _,sub_ptr := range ptr_block.B_pointers().Get(){
					if sub_ptr != -1{must_erase = false;break}
				}
				if must_erase{
					self.Block_bitmap.Erase(1,self.Block_section.Bit_no_for(ptr))
					dir.I_block().No(int32(n)).Set(-1)
				}	

				return trgt_index, trgt_block
			} 
		}else{
			dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
			for i, content := range dir_block.B_content().Get(){
				if content.B_inodo == -1 {continue}
				if name == content.B_name {
					inode_result := types.CreateIndexNode(self.super_service,content.B_inodo)
				
					dir_block.B_content().No(int32(i)).B_inodo().Set(-1)
					must_erase := true
					for _,cntnt := range dir_block.B_content().Get(){
						if cntnt.B_inodo != -1 {must_erase = false;break}
					}
					if must_erase{
						self.Block_bitmap.Erase(1,self.Block_section.Bit_no_for(ptr))
						dir.I_block().No(int32(n)).Set(-1)
					}
				
					return content.B_inodo,inode_result
				}
			}
		}
	}
	return -1,types.CreateIndexNode(self.super_service,0)
}












func (self *FormatEXT2)Put_in_dir(dir types.IndexNode, trgt_inode types.IndexNodeHolder,name [12]string) (int32,types.IndexNode){
	var index int32 
	var inode types.IndexNode
	index,inode = self.try_put_in_existing_dir_block(dir,trgt_inode,name)
	if index == -1 {
		// fmt.Printf("Debug: is comming in second if\n")
		index,inode = self.try_put_in_new_dir_block_in_existing_ptr(dir,trgt_inode,name)
	}
	if index == -1 {
		// fmt.Printf("Debug: is comming in last if\n")
		index,inode = self.try_put_in_new_dir_block_in_new_ptr(dir,trgt_inode,name)
	}
	return index,inode
}
func (self *FormatEXT2)case3_ptr_block(pointer_block *types.PointerBlock, level int32, trgt_inode types.IndexNodeHolder,name [12]string) (int32,types.IndexNode){
	if level != 0 {
		for n,ptr := range pointer_block.B_pointers().Get(){
			if ptr != -1 { continue }
			ptr_block_indx,ptr_block := self.Create_PointerBlock()
			pointer_block.B_pointers().No(int32(n)).Set(ptr_block_indx)
			trgt_index,trgt_block := self.case3_ptr_block(&ptr_block,level-1,trgt_inode,name)
			if trgt_index != -1 {
				return trgt_index, trgt_block
			} 
		}
		return -1,types.CreateIndexNode(self.super_service,0)
	}
	for n,ptr := range pointer_block.B_pointers().Get(){
		if ptr != -1 { continue }
		dir_blck_index,dir_block := self.Create_DirectoryBlock()
		
		pointer_block.B_pointers().No(int32(n)).Set(dir_blck_index)
		inode_indx,inode := self.Create_Inode(trgt_inode)
		dir_block.B_content().No(0).Set(types.ContentHolder{
			B_name:  name,
			B_inodo: inode_indx,
		})
		return inode_indx,inode
	}
	return -1,types.CreateIndexNode(self.super_service,0)

}
func (self *FormatEXT2)try_put_in_new_dir_block_in_new_ptr(dir types.IndexNode, trgt_inode types.IndexNodeHolder,name [12]string) (int32,types.IndexNode){
	for n,ptr := range dir.I_block().Get(){
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			indx,inode := self.case3_ptr_block(&ptr_block,int32(n-13),trgt_inode,name)
			if indx != -1 {
				return indx, inode
			} 
		}else{
			if ptr != -1 { continue }
			dir_blck_index,dir_block := self.Create_DirectoryBlock()
			dir.I_block().No(int32(n)).Set(dir_blck_index)
			inode_indx,inode := self.Create_Inode(trgt_inode)
			dir_block.B_content().No(0).Set(types.ContentHolder{
				B_name:  name,
				B_inodo: inode_indx,
			})
			return inode_indx,inode
		}
	}
	return -1,types.CreateIndexNode(self.super_service,0)
}

func (self *FormatEXT2)case2_ptr_block(pointer_block *types.PointerBlock, level int32, trgt_inode types.IndexNodeHolder,name [12]string) (int32,types.IndexNode){
	if level != 0 {
		for _,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_index,trgt_block := self.case2_ptr_block(&ptr_block,level-1,trgt_inode,name)
			if trgt_index != -1 {
				return trgt_index, trgt_block
			} 
		}
		return -1,types.CreateIndexNode(self.super_service,0)
	}
	for n,ptr := range pointer_block.B_pointers().Get(){
		if ptr != -1 { continue }
		dir_blck_index,dir_block := self.Create_DirectoryBlock()
		
		pointer_block.B_pointers().No(int32(n)).Set(dir_blck_index)
		inode_indx,inode := self.Create_Inode(trgt_inode)
		dir_block.B_content().No(0).Set(types.ContentHolder{
			B_name:  name,
			B_inodo: inode_indx,
		})
		return inode_indx,inode
	}
	return -1,types.CreateIndexNode(self.super_service,0)

}
func (self *FormatEXT2)try_put_in_new_dir_block_in_existing_ptr(dir types.IndexNode, trgt_inode types.IndexNodeHolder,name [12]string) (int32,types.IndexNode){
	for n,ptr := range dir.I_block().Get(){
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			indx,inode := self.case2_ptr_block(&ptr_block,int32(n-13),trgt_inode,name)
			if indx != -1 {
				return indx, inode
			} 
		}else{
			if ptr != -1 { continue }
			dir_blck_index,dir_block := self.Create_DirectoryBlock()
			dir.I_block().No(int32(n)).Set(dir_blck_index)
			inode_indx,inode := self.Create_Inode(trgt_inode)
			dir_block.B_content().No(0).Set(types.ContentHolder{
				B_name:  name,
				B_inodo: inode_indx,
			})
			return inode_indx,inode
		}
	}
	return -1,types.CreateIndexNode(self.super_service,0)
}
func (self *FormatEXT2)case1_ptr_blck(pointer_block *types.PointerBlock, level int32, trgt_inode types.IndexNodeHolder,name [12]string) (int32,types.IndexNode){
	if level != 0 {
		for _,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_index,trgt_block := self.case1_ptr_blck(&ptr_block,level-1,trgt_inode,name)
			if trgt_index != -1 {
				return trgt_index, trgt_block
			} 
		}
		return -1,types.CreateIndexNode(self.super_service,0)

	}
	for _,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
		for _, content := range dir_block.B_content().Spread(){
			if content.B_inodo().Get() != -1 {continue}
			indx,inode := self.Create_Inode(trgt_inode)
			content.B_inodo().Set(indx)
			content.B_name().Set(name)
			return indx,inode
		}
	}
	return -1,types.CreateIndexNode(self.super_service,0)
}
func (self *FormatEXT2)try_put_in_existing_dir_block(dir types.IndexNode, trgt_inode types.IndexNodeHolder,name [12]string) (int32,types.IndexNode){
	for n,ptr := range dir.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			indx,inode := self.case1_ptr_blck(&ptr_block,int32(n-13),trgt_inode,name)
			if indx != -1 {
				return indx, inode
			} 
		}else{
			dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
			for _, content := range dir_block.B_content().Spread(){
				if content.B_inodo().Get() != -1 {continue}
				indx,inode := self.Create_Inode(trgt_inode)
				content.B_inodo().Set(indx)
				content.B_name().Set(name)
				return indx,inode
			}
		}
	}
	return -1,types.CreateIndexNode(self.super_service,0)
}













func (self *FormatEXT2)shallow_tree_from_ptr_blck(pointer_block *types.PointerBlock,level int32,all_content *[]types.Content) {
	if level != 0 {
		for _,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			self.shallow_tree_from_ptr_blck(&ptr_block,level-1,all_content)
		}
		return
	}
	for _,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
		for _, content := range dir_block.B_content().Spread(){
			if content.B_inodo().Get() == -1 {continue}
			*all_content = append(*all_content, content)
		}
	}
}
func (self *FormatEXT2)Get_shallow_tree(dir types.IndexNode) []types.Content {
	all_content := make([]types.Content,0,10)
	for n,ptr := range dir.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			self.shallow_tree_from_ptr_blck(&ptr_block,int32(n-13),&all_content)
			
		}else{
			dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
			for _, content := range dir_block.B_content().Spread(){
				if content.B_inodo().Get() == -1 {continue}
				all_content = append(all_content, content)
			}
		}
	}
	return all_content
}



func (self *FormatEXT2) User_allowed_actions(user int32, grp int32,inode *types.IndexNode)utiles.Permision{
	if user == 1 && grp == 1{return utiles.ALL_PERMITION}
	user_corr := inode.I_uid().Get()
	user_grp_corr := inode.I_gid().Get()
	if user == user_corr {
		return utiles.Permision_from_str(inode.I_perm().No(0).Get())
	}
	if grp == user_grp_corr{
		return utiles.Permision_from_str(inode.I_perm().No(1).Get())
	}
	return utiles.Permision_from_str(inode.I_perm().No(3).Get())
}


func (self *FormatEXT2) Get_nested_dir(init_dir types.IndexNode, folders [][12]string, create_recursive bool,usr_id int32,usr_grp_id int32,time types.TimeHolder,for_read bool,for_write bool) (bool,types.IndexNode) {
	init_dir.I_atime().Set(time)
	dir := init_dir
	for _, dir_name := range folders {
		var result int32
		result, content := self.Search_for_inode(dir, dir_name)

		if result == -1 {
			if !create_recursive {
				return false,types.IndexNode{}
			}
			result, new_dir := self.Put_in_dir(dir,types.IndexNodeHolder{
				I_uid:   usr_id,
				I_gid:   usr_grp_id,
				I_s:     0,
				I_atime: utiles.NO_TIME,
				I_ctime: time,
				I_mtime: utiles.NO_TIME,
				I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
				I_type:  string(utiles.Directory),
				I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
			},dir_name)
			if result == -1 {
				return false,types.IndexNode{}
			}
			dir.I_mtime().Set(time)
			dir = new_dir
		} else {
			dir = types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
			if dir.I_type().Get() == string(utiles.File) {return false,types.IndexNode{}}
			permisions := self.User_allowed_actions(usr_id,usr_grp_id,&dir)
			if for_read && !permisions.Can_read(){
				return false,types.IndexNode{}
			}
			if for_write && !permisions.Can_write(){
				return false,types.IndexNode{}
			}
			dir.I_atime().Set(time)
		}
	}
	return true,dir
}









func (self *FormatEXT2)erase_dir_from_ptr_blck(pointer_block *types.PointerBlock,level int32,usr_id int32,usr_grp_id int32,time types.TimeHolder)bool {
	can_erase_self:=true
	if level != 0 {
		for i,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			erase := self.erase_dir_from_ptr_blck(&ptr_block,level-1,usr_id,usr_grp_id,time)
			if !erase{
				can_erase_self = false
				continue
			}
			bit := self.Block_section.Bit_no_for(ptr_block.Index)
			self.Block_bitmap.Erase(1,bit)
			pointer_block.B_pointers().No(int32(i)).Set(-1)
		}
		return can_erase_self
	}
	for _,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
		for i, content := range dir_block.B_content().Get(){
			if content.B_inodo == -1 {continue}
			inode := types.CreateIndexNode(pointer_block.Super_service,content.B_inodo)
			permision := self.User_allowed_actions(usr_id,usr_grp_id,&inode)
			if !permision.Can_write(){
				can_erase_self=false
				continue
			}
			if inode.I_type().Get()== string(utiles.Directory){
				erase := self.Erase_dir(inode,usr_id,usr_grp_id,time)
				if !erase{
					can_erase_self = false
					inode.I_mtime().Set(time)
					continue
				} 
			}
			bit := self.Inodes_section.Bit_no_for(content.B_inodo)
			self.Inodes_bitmap.Erase(1,bit)
			dir_block.B_content().No(int32(i)).B_inodo().Set(-1)
		}
	}
	return can_erase_self
}
func (self *FormatEXT2)Erase_dir(dir types.IndexNode,usr_id int32,usr_grp_id int32,time types.TimeHolder) bool {
	can_erase_self := false
	for n,ptr := range dir.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			erase := self.erase_dir_from_ptr_blck(&ptr_block,int32(n)-13,usr_id,usr_grp_id,time)
			if !erase{
				can_erase_self = false
				continue
			}
			bit := self.Block_section.Bit_no_for(ptr_block.Index)
			self.Block_bitmap.Erase(1,bit)
			dir.I_block().No(int32(n)).Set(-1)
			
		}else{
			dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
			for i, content := range dir_block.B_content().Get(){
				if content.B_inodo == -1 {continue}
				inode := types.CreateIndexNode(dir.Super_service,content.B_inodo)
				permision := self.User_allowed_actions(usr_id,usr_grp_id,&inode)
				if !permision.Can_write(){
					can_erase_self=false
					continue
				}
				if inode.I_type().Get()== string(utiles.Directory){
					erase := self.Erase_dir(inode,usr_id,usr_grp_id,time)
					if !erase{
						can_erase_self = false
						inode.I_mtime().Set(time)
						continue
					} 
				}
				bit := self.Inodes_section.Bit_no_for(content.B_inodo)
				self.Inodes_bitmap.Erase(1,bit)
				dir_block.B_content().No(int32(i)).B_inodo().Set(-1)
			}
		}
	}
	return can_erase_self
}



func (self *FormatEXT2)Remove_dir_if_possilbe(in_dir types.IndexNode,with_name [12]string,usr_id int32,usr_grp_id int32,time types.TimeHolder) bool {
	for n,ptr := range in_dir.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			found := self.renove_dir_if_possilbe_in_ptr_block(&ptr_block,int32(n-13),with_name,usr_id,usr_grp_id,time)
			if found {
				for _,sub_ptr := range ptr_block.B_pointers().Get(){
					if sub_ptr != -1{return true}
				}
				self.Block_bitmap.Erase(1,self.Block_section.Bit_no_for(ptr))
				in_dir.I_block().No(int32(n)).Set(-1)
				return true
			}
		}else{
			dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
			for i, content := range dir_block.B_content().Get(){
				if content.B_inodo == -1 {continue}
				if with_name == content.B_name {
					inode_result := types.CreateIndexNode(self.super_service,content.B_inodo)
					permission := self.User_allowed_actions(usr_id,usr_grp_id,&inode_result)
					if !permission.Can_write() {return true}
					must_erase := self.Erase_dir(inode_result,usr_id,usr_grp_id,time)
					if must_erase{
						dir_block.B_content().No(int32(i)).B_inodo().Set(-1)
						self.Inodes_bitmap.Erase(1,self.Inodes_section.Bit_no_for(content.B_inodo))
						for _,cntnt := range dir_block.B_content().Get(){
							if cntnt.B_inodo != -1 {return true}
						}
						self.Block_bitmap.Erase(1,self.Block_section.Bit_no_for(ptr))
						in_dir.I_block().No(int32(n)).Set(-1)
						return true
					}
				
					return true
				}
			}
		}
	}
	return false
}








func (self *FormatEXT2)renove_dir_if_possilbe_in_ptr_block(pointer_block *types.PointerBlock, level int32, name [12]string,usr_id int32,usr_grp_id int32,time types.TimeHolder) bool {
	if level != 0 {
		for i,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			found := self.renove_dir_if_possilbe_in_ptr_block(&ptr_block,level-1,name,usr_id,usr_grp_id,time)
			if found {
				for _,sub_ptr := range ptr_block.B_pointers().Get(){
					if sub_ptr != -1{return true}
				}
				self.Block_bitmap.Erase(1,self.Block_section.Bit_no_for(ptr))
				pointer_block.B_pointers().No(int32(i)).Set(-1)
				return true
			}
		}
		return false
	}
	for n,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
		for i, content := range dir_block.B_content().Get(){
			if content.B_inodo == -1 {continue}
			if name == content.B_name {
				inode_result := types.CreateIndexNode(self.super_service,content.B_inodo)
				permission := self.User_allowed_actions(usr_id,usr_grp_id,&inode_result)
				if !permission.Can_write() {return true}
				must_erase := self.Erase_dir(inode_result,usr_id,usr_grp_id,time)
				if must_erase{
					dir_block.B_content().No(int32(i)).B_inodo().Set(-1)
					self.Inodes_bitmap.Erase(1,self.Inodes_section.Bit_no_for(content.B_inodo))
					for _,cntnt := range dir_block.B_content().Get(){
						if cntnt.B_inodo != -1 {return true}
					}
					self.Block_bitmap.Erase(1,self.Block_section.Bit_no_for(ptr))
					pointer_block.B_pointers().No(int32(n)).Set(-1)
					return true
				}
			
				return true
			}
		}
	}
	return false

}


func (self *FormatEXT2)Directory_deep_copy(dest_dir types.IndexNode,trgt_dir types.IndexNode,with_name [12]string,usr_id int32,usr_grp_id int32,time types.TimeHolder) bool {
	perm_to_dest := self.User_allowed_actions(usr_id,usr_grp_id,&dest_dir)
	if !perm_to_dest.Can_write(){return false}
	perm_to_trgt := self.User_allowed_actions(usr_id,usr_grp_id,&trgt_dir)
	if !perm_to_trgt.Can_read(){return false}
	switch trgt_dir.I_type().Get() {
	case string(utiles.File):
		result,new_inode := self.Put_in_dir(dest_dir,types.IndexNodeHolder{
			I_uid:   usr_id,
			I_gid:   usr_grp_id,
			I_s:     trgt_dir.I_s().Get(),
			I_atime: utiles.NO_TIME,
			I_ctime: time,
			I_mtime: time,
			I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
			I_type:  string(utiles.File),
			I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
		},with_name)
		if result == -1 {return false}
		data:=self.Read_file(&trgt_dir)
		self.Update_file(&new_inode,0,data)

	case string(utiles.Directory):
		result,new_inode := self.Put_in_dir(dest_dir,types.IndexNodeHolder{
			I_uid:   usr_id,
			I_gid:   usr_grp_id,
			I_s:     trgt_dir.I_s().Get(),
			I_atime: utiles.NO_TIME,
			I_ctime: time,
			I_mtime: time,
			I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
			I_type:  string(utiles.Directory),
			I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
		},with_name)
		if result == -1 {return false}
		sub_folders := self.Get_shallow_tree(trgt_dir)
		for _,folder:=range sub_folders{
			inode:=types.CreateIndexNode(folder.Super_service,folder.B_inodo().Get())
			self.Directory_deep_copy(new_inode,inode,folder.B_name().Get(),usr_id,usr_grp_id,time)
		}
	default:panic("Didnt match neither file nor directory")
	}

	return true
}














func (self *FormatEXT2)Put_indx_in_dir(dir types.IndexNode, trgt_inode int32,name [12]string) (int32){
	var index int32 
	index = self.try_put_indx_in_existing_dir_block(dir,trgt_inode,name)
	if index == -1 {
		// fmt.Printf("Debug: is comming in second if\n")
		index = self.try_put_indx_in_new_dir_block_in_existing_ptr(dir,trgt_inode,name)
	}
	if index == -1 {
		// fmt.Printf("Debug: is comming in last if\n")
		index = self.try_put_indx_in_new_dir_block_in_new_ptr(dir,trgt_inode,name)
	}
	return index
}
func (self *FormatEXT2)case3_ptr_block_indx(pointer_block *types.PointerBlock, level int32, trgt_inode int32,name [12]string) (int32){
	if level != 0 {
		for n,ptr := range pointer_block.B_pointers().Get(){
			if ptr != -1 { continue }
			ptr_block_indx,ptr_block := self.Create_PointerBlock()
			pointer_block.B_pointers().No(int32(n)).Set(ptr_block_indx)
			trgt_index := self.case3_ptr_block_indx(&ptr_block,level-1,trgt_inode,name)
			if trgt_index != -1 {
				return trgt_index
			} 
		}
		return -1
	}
	for n,ptr := range pointer_block.B_pointers().Get(){
		if ptr != -1 { continue }
		dir_blck_index,dir_block := self.Create_DirectoryBlock()
		
		pointer_block.B_pointers().No(int32(n)).Set(dir_blck_index)
		dir_block.B_content().No(0).Set(types.ContentHolder{
			B_name:  name,
			B_inodo: trgt_inode,
		})
		return trgt_inode
	}
	return -1

}
func (self *FormatEXT2)try_put_indx_in_new_dir_block_in_new_ptr(dir types.IndexNode, trgt_inode int32,name [12]string) (int32){
	for n,ptr := range dir.I_block().Get(){
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			indx := self.case3_ptr_block_indx(&ptr_block,int32(n-13),trgt_inode,name)
			if indx != -1 {
				return indx
			} 
		}else{
			if ptr != -1 { continue }
			dir_blck_index,dir_block := self.Create_DirectoryBlock()
			dir.I_block().No(int32(n)).Set(dir_blck_index)
			dir_block.B_content().No(0).Set(types.ContentHolder{
				B_name:  name,
				B_inodo: trgt_inode,
			})
			return trgt_inode
		}
	}
	return -1
}

func (self *FormatEXT2)case2_ptr_block_indx(pointer_block *types.PointerBlock, level int32, trgt_inode int32,name [12]string) (int32){
	if level != 0 {
		for _,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_index := self.case2_ptr_block_indx(&ptr_block,level-1,trgt_inode,name)
			if trgt_index != -1 {
				return trgt_index
			} 
		}
		return -1
	}
	for n,ptr := range pointer_block.B_pointers().Get(){
		if ptr != -1 { continue }
		dir_blck_index,dir_block := self.Create_DirectoryBlock()
		
		pointer_block.B_pointers().No(int32(n)).Set(dir_blck_index)
		dir_block.B_content().No(0).Set(types.ContentHolder{
			B_name:  name,
			B_inodo: trgt_inode,
		})
		return trgt_inode
	}
	return -1

}
func (self *FormatEXT2)try_put_indx_in_new_dir_block_in_existing_ptr(dir types.IndexNode, trgt_inode int32,name [12]string) (int32){
	for n,ptr := range dir.I_block().Get(){
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			indx := self.case2_ptr_block_indx(&ptr_block,int32(n-13),trgt_inode,name)
			if indx != -1 {
				return indx
			} 
		}else{
			if ptr != -1 { continue }
			dir_blck_index,dir_block := self.Create_DirectoryBlock()
			dir.I_block().No(int32(n)).Set(dir_blck_index)
			dir_block.B_content().No(0).Set(types.ContentHolder{
				B_name:  name,
				B_inodo: trgt_inode,
			})
			return trgt_inode
		}
	}
	return -1
}
func (self *FormatEXT2)case1_ptr_blck_indx(pointer_block *types.PointerBlock, level int32, trgt_inode int32,name [12]string) (int32){
	if level != 0 {
		for _,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_index := self.case1_ptr_blck_indx(&ptr_block,level-1,trgt_inode,name)
			if trgt_index != -1 {
				return trgt_index
			} 
		}
		return -1

	}
	for _,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
		for _, content := range dir_block.B_content().Spread(){
			if content.B_inodo().Get() != -1 {continue}
			content.B_inodo().Set(trgt_inode)
			content.B_name().Set(name)
			return trgt_inode
		}
	}
	return -1
}
func (self *FormatEXT2)try_put_indx_in_existing_dir_block(dir types.IndexNode, trgt_inode int32,name [12]string) (int32){
	for n,ptr := range dir.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			indx := self.case1_ptr_blck_indx(&ptr_block,int32(n-13),trgt_inode,name)
			if indx != -1 {
				return indx
			} 
		}else{
			dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
			for _, content := range dir_block.B_content().Spread(){
				if content.B_inodo().Get() != -1 {continue}
				content.B_inodo().Set(trgt_inode)
				content.B_name().Set(name)
				return trgt_inode
			}
		}
	}
	return -1
}




















func (self *FormatEXT2)Find_in_dir(in_dir types.IndexNode, name_criteria utiles.NameCriteria,usr_id int32,usr_grp_id int32,time types.TimeHolder) []types.Content {
	matched_files := make([]types.Content,0,10)
	perm_to_dest := self.User_allowed_actions(usr_id,usr_grp_id,&in_dir)
	if !perm_to_dest.Can_read(){return matched_files}
	in_dir.I_atime().Set(time)
	for _,content:= range self.Get_shallow_tree(in_dir){
		if !name_criteria.Match(content.B_name().Get()){continue}
		matched_files = append(matched_files, content)
	}
	return matched_files
}