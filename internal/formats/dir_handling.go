package formats

import (
	"fmt"
	"project/internal/types"
	"project/internal/utiles"
)

var OWN_DIR_NAME [12]string = [12]string{"."," "," "," "," "," "," "," "," "," "," "," "}
var PARENT_DIR_NAME [12]string = [12]string{".","."," "," "," "," "," "," "," "," "," "," "}
func (self *Format)Set_parent_child_relation(parent types.IndexNode, child types.IndexNode) {
	first_dir_block_indx:=child.I_block().No(0).Get()
	var dir_blck types.DirectoryBlock
	if first_dir_block_indx == -1{
		_,dir_blck=self.Create_DirectoryBlock()
		child.I_block().No(0).Set(dir_blck.Index)
	}else{
		dir_blck=types.CreateDirectoryBlock(self.super_service,first_dir_block_indx)
	}
	trgt_space:=dir_blck.B_content().No(0)
	trgt_space.B_name().Set(OWN_DIR_NAME)
	trgt_space.B_inodo().Set(child.Index)
	trgt_space=dir_blck.B_content().No(1)
	trgt_space.B_name().Set(PARENT_DIR_NAME)
	trgt_space.B_inodo().Set(parent.Index)
}


func (self *Format)search_in_ptr_block(pointer_block *types.PointerBlock, level int32, name [12]string) (int32,types.Content) {
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
func (self *Format)Search_for_inode(dir types.IndexNode, name [12]string) (int32,types.Content) {
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










func (self *Format)extract_inode_in_ptr_block(pointer_block *types.PointerBlock, level int32, name [12]string) (int32,types.IndexNode) {
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
func (self *Format)Extract_inode(dir types.IndexNode, name [12]string) (int32,types.IndexNode) {
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








type InodeTemplate struct{
	index int32
	sketch types.IndexNodeHolder
}
func (self *InodeTemplate) is_index()bool{return self.index != -1}
func (self *InodeTemplate) is_sketch()bool{return self.index == -1}



func (self *Format)Put_in_dir(dir types.IndexNode, trgt_inode InodeTemplate,name [12]string) (types.IndexNode){ 
	var inode types.IndexNode
	inode = self.try_put_in_existing_dir_block(dir,trgt_inode,name)
	if inode.Index == -1 {
		inode = self.try_put_in_new_dir_block_in_existing_ptr(dir,trgt_inode,name)
	}
	return inode
}
func (self *Format)create_ptr_block_recursive_and_append(pointer_block *types.PointerBlock, level int32, trgt_inode InodeTemplate,name [12]string) (types.IndexNode){
	if level != 0 {
		for n,ptr := range pointer_block.B_pointers().Get(){
			if ptr != -1 { continue }
			_,ptr_block := self.Create_PointerBlock()
			pointer_block.B_pointers().No(int32(n)).Set(ptr_block.Index)
			trgt_block := self.create_ptr_block_recursive_and_append(&ptr_block,level-1,trgt_inode,name)
			if trgt_block.Index != -1 {
				return trgt_block
			} 
		}
		return types.CreateIndexNode(nil,-1)
	}
	for n,ptr := range pointer_block.B_pointers().Get(){
		if ptr != -1 { continue }
		dir_blck_index,dir_block := self.Create_DirectoryBlock()
		
		pointer_block.B_pointers().No(int32(n)).Set(dir_blck_index)
		var result_inode types.IndexNode
		if trgt_inode.is_sketch(){
			_,result_inode = self.Create_Inode(trgt_inode.sketch)
		}else{
			result_inode = types.CreateIndexNode(self.super_service,trgt_inode.index)
		}
		dir_block.B_content().No(0).Set(types.ContentHolder{
			B_name:  name,
			B_inodo: result_inode.Index,
		})
		return result_inode
	}
	return types.CreateIndexNode(nil,-1)

}


func (self *Format)case2_ptr_block(pointer_block *types.PointerBlock, level int32, trgt_inode InodeTemplate,name [12]string) (types.IndexNode,*PntrBlockAvailable){
	var first_empty *PntrBlockAvailable = nil
	var last_available *PntrBlockAvailable = nil
	
	if level != 0 {
		for _,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 {
				if first_empty == nil{
					available:=PntrBlockAvailable{ptr_block:*pointer_block,level:level}
					first_empty = &available
				} 
				continue
			}
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_block,sub_available := self.case2_ptr_block(&ptr_block,level-1,trgt_inode,name)
			if trgt_block.Index != -1 {
				return trgt_block,nil
			}
			if sub_available != nil{
				if last_available != nil{
					if sub_available.level < last_available.level{
						last_available = sub_available
					} 
				}else{
					last_available = sub_available
				}
			}
		}
		if last_available != nil {return types.CreateIndexNode(nil,-1),last_available}
		return types.CreateIndexNode(nil,-1),first_empty
	}
	for n,ptr := range pointer_block.B_pointers().Get(){
		if ptr != -1 { continue }
		dir_blck_index,dir_block := self.Create_DirectoryBlock()
		
		pointer_block.B_pointers().No(int32(n)).Set(dir_blck_index)
		var result_inode types.IndexNode
		if trgt_inode.is_sketch(){
			_,result_inode = self.Create_Inode(trgt_inode.sketch)
		}else{
			result_inode = types.CreateIndexNode(self.super_service,trgt_inode.index)
		}
		dir_block.B_content().No(0).Set(types.ContentHolder{
			B_name:  name,
			B_inodo: result_inode.Index,
		})
		return result_inode,nil
	}
	if last_available != nil {return types.CreateIndexNode(nil,-1),last_available}
	return types.CreateIndexNode(nil,-1),first_empty

}
type PntrBlockAvailable struct{
	ptr_block types.PointerBlock
	level int32
}
func (self *Format)try_put_in_new_dir_block_in_existing_ptr(dir types.IndexNode, trgt_inode InodeTemplate,name [12]string) (types.IndexNode){
	var pntr_block_available *PntrBlockAvailable = nil
	for n,ptr := range dir.I_block().Get(){
		if n >= 13{
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			result_inode,sub_available := self.case2_ptr_block(&ptr_block,int32(n-13),trgt_inode,name)
			if result_inode.Index != -1 {
				return result_inode
			} 
			pntr_block_available=sub_available
		}else{
			if ptr != -1 { continue }
			dir_blck_index,dir_block := self.Create_DirectoryBlock()
			dir.I_block().No(int32(n)).Set(dir_blck_index)
			var result_inode types.IndexNode
			if trgt_inode.is_sketch(){
				_,result_inode = self.Create_Inode(trgt_inode.sketch)
			}else{
				result_inode = types.CreateIndexNode(self.super_service,trgt_inode.index)
			}
			dir_block.B_content().No(0).Set(types.ContentHolder{
				B_name:  name,
				B_inodo: result_inode.Index,
			})
			return result_inode
		}
	}
	if pntr_block_available != nil{
		return self.create_new_ptr_and_blocks(&pntr_block_available.ptr_block,pntr_block_available.level,trgt_inode,name)
	}else{
		if dir.I_block().No(13).Get() == -1 {
			_,new_ptr_blck:=self.Create_PointerBlock()
			dir.I_block().No(int32(13)).Set(new_ptr_blck.Index)
			return self.create_new_ptr_and_blocks(&new_ptr_blck,0,trgt_inode,name)
		}else if dir.I_block().No(14).Get() == -1 {
			_,new_ptr_blck:=self.Create_PointerBlock()
			dir.I_block().No(int32(14)).Set(new_ptr_blck.Index)
			return self.create_new_ptr_and_blocks(&new_ptr_blck,1,trgt_inode,name)
		}else if dir.I_block().No(15).Get() == -1 {
			_,new_ptr_blck:=self.Create_PointerBlock()
			dir.I_block().No(int32(15)).Set(new_ptr_blck.Index)
			return self.create_new_ptr_and_blocks(&new_ptr_blck,2,trgt_inode,name)
		}
	}
	return types.CreateIndexNode(nil,-1)
}
func (self *Format)create_new_ptr_and_blocks(pointer_block *types.PointerBlock, level int32, trgt_inode InodeTemplate,name [12]string) types.IndexNode{
	if level != 0 {
		for i,ptr := range pointer_block.B_pointers().Get(){
			if ptr != -1 { continue }
			_,new_ptr_blck:=self.Create_PointerBlock()
			pointer_block.B_pointers().No(int32(i)).Set(new_ptr_blck.Index)
			trgt_block := self.create_new_ptr_and_blocks(&new_ptr_blck,level-1,trgt_inode,name)
			return trgt_block
		}
		panic("Wrong determination for free block pointer")

	}
	for n,ptr := range pointer_block.B_pointers().Get(){
		if ptr != -1 { continue }
		_,new_dir_block := self.Create_DirectoryBlock()
		pointer_block.B_pointers().No(int32(n)).Set(new_dir_block.Index)
		var new_inode types.IndexNode
		if trgt_inode.is_index(){
			new_inode = types.CreateIndexNode(self.super_service,trgt_inode.index)
		}else{
			_,new_inode=self.Create_Inode(trgt_inode.sketch)
		}
		new_dir_block.B_content().No(0).Set(types.ContentHolder{
			B_name:  name,
			B_inodo: new_inode.Index,
		})
		return new_inode
	}
	panic("Wrong determination for free block pointer")
}
func (self *Format)case1_ptr_blck(pointer_block *types.PointerBlock, level int32, trgt_inode InodeTemplate,name [12]string) types.IndexNode{
	if level != 0 {
		for _,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_block := self.case1_ptr_blck(&ptr_block,level-1,trgt_inode,name)
			if trgt_block.Index != -1 {
				return trgt_block
			} 
		}
		return types.CreateIndexNode(nil,-1)

	}
	for _,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
		for _, content := range dir_block.B_content().Spread(){
			if content.B_inodo().Get() != -1 {continue}
			var result_inode types.IndexNode
			if trgt_inode.is_sketch(){
				_,result_inode = self.Create_Inode(trgt_inode.sketch)
				content.B_inodo().Set(result_inode.Index)
				content.B_name().Set(name)
			}else{
				content.B_inodo().Set(trgt_inode.index)
				content.B_name().Set(name)
				result_inode = types.CreateIndexNode(self.super_service,trgt_inode.index)
			}
			return result_inode
		}
	}
	return types.CreateIndexNode(nil,-1)
}
func (self *Format)try_put_in_existing_dir_block(dir types.IndexNode, trgt_inode InodeTemplate,name [12]string) types.IndexNode{
	for n,ptr := range dir.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			result_inode := self.case1_ptr_blck(&ptr_block,int32(n-13),trgt_inode,name)
			if result_inode.Index != -1 {
				return result_inode
			} 
		}else{
			dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
			for _, content := range dir_block.B_content().Spread(){
				if content.B_inodo().Get() != -1 {continue}
				var result_inode types.IndexNode
				if trgt_inode.is_sketch(){
					_,result_inode = self.Create_Inode(trgt_inode.sketch)
					content.B_inodo().Set(result_inode.Index)
					content.B_name().Set(name)
				}else{
					content.B_inodo().Set(trgt_inode.index)
					content.B_name().Set(name)
					result_inode = types.CreateIndexNode(self.super_service,trgt_inode.index)
				}
				return result_inode
			}
		}
	}
	return types.CreateIndexNode(nil,-1)
}













func (self *Format)shallow_tree_from_ptr_blck(pointer_block *types.PointerBlock,level int32,all_content *[]types.Content) {
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
func (self *Format)Get_shallow_tree_of_childs(dir types.IndexNode) []types.Content {
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
func (self *Format)strict_shallow_tree_from_ptr_blck(pointer_block *types.PointerBlock,level int32,all_content *[]types.Content) {
	if level != 0 {
		for _,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			self.strict_shallow_tree_from_ptr_blck(&ptr_block,level-1,all_content)
		}
		return
	}
	for _,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
		for _, content := range dir_block.B_content().Spread(){
			if content.B_inodo().Get() == -1 {continue}
			if content.B_name().Get() == OWN_DIR_NAME {continue}
			if content.B_name().Get() == PARENT_DIR_NAME {continue}
			*all_content = append(*all_content, content)
		}
	}
}
func (self *Format)Get_strict_shallow_tree_of_childs(dir types.IndexNode) []types.Content {
	all_content := make([]types.Content,0,10)
	for n,ptr := range dir.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			self.strict_shallow_tree_from_ptr_blck(&ptr_block,int32(n-13),&all_content)
			
		}else{
			dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
			for _, content := range dir_block.B_content().Spread(){
				if content.B_inodo().Get() == -1 {continue}
				if content.B_name().Get() == OWN_DIR_NAME {continue}
				if content.B_name().Get() == PARENT_DIR_NAME {continue}
				all_content = append(all_content, content)
			}
		}
	}
	return all_content
}



type PermisionCase struct{
	permision utiles.Permision
	status int
}
func (self *PermisionCase) Is_owner()bool{return self.status == 1}
func (self *PermisionCase) Is_grp_part()bool{return self.status == 2}
func (self *PermisionCase) Is_other()bool{return self.status == 3}

func (self *Format) User_allowed_actions_with_ownership(user int32, grp int32,inode *types.IndexNode)PermisionCase{
	user_corr := inode.I_uid().Get()
	user_grp_corr := inode.I_gid().Get()
	case_stat := 3
	if user == user_corr && grp == user_grp_corr{
		case_stat = 1
	}else if grp == user_grp_corr{
		case_stat = 2
	}
	if user == 1 && grp == 1{return PermisionCase{
		permision: utiles.ALL_PERMITION,
		status:    case_stat,
	}}

	if user == user_corr {return PermisionCase{
		permision: utiles.Permision_from_str(inode.I_perm().No(0).Get()),
		status:    case_stat,
	}}
	if grp == user_grp_corr{return PermisionCase{
		permision: utiles.Permision_from_str(inode.I_perm().No(1).Get()),
		status:    case_stat,
	}}
	return PermisionCase{
		permision: utiles.Permision_from_str(inode.I_perm().No(3).Get()),
		status:    case_stat,
	}
}
func (self *Format) User_allowed_actions(user int32, grp int32,inode *types.IndexNode)utiles.Permision{
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

func (self *Format) Wrap_indx_in_template(index int32)InodeTemplate{
	return InodeTemplate{
		index:  index,
		sketch: types.IndexNodeHolder{},
	}
}
func (self *Format) Wrap_holder_in_template(template types.IndexNodeHolder)InodeTemplate{
	return InodeTemplate{
		index:  -1,
		sketch: template,
	}
}
func (self *Format) Get_nested_dir(init_dir types.IndexNode, folders [][12]string, create_recursive bool,usr_id int32,usr_grp_id int32,time types.TimeHolder,for_read bool,for_write bool) (error,types.IndexNode) {
	init_dir.I_atime().Set(time)
	dir := init_dir
	for _, dir_name := range folders {
		var result int32
		result, content := self.Search_for_inode(dir, dir_name)

		if result == -1 {
			if !create_recursive {
				return fmt.Errorf("path doesnt exist and recursive creation is disabled"),types.IndexNode{}
			}
			new_dir := self.Put_in_dir(dir,self.Wrap_holder_in_template(types.IndexNodeHolder{
					I_uid:   usr_id,
					I_gid:   usr_grp_id,
					I_s:     0,
					I_atime: utiles.NO_TIME,
					I_ctime: time,
					I_mtime: utiles.NO_TIME,
					I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
					I_type:  string(utiles.Directory),
					I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
				}),dir_name)
			if new_dir.Index == -1 {
				return fmt.Errorf("problems in creating new dirs"),types.IndexNode{}
			}
			self.Set_parent_child_relation(dir,new_dir)
			dir.I_mtime().Set(time)
			dir = new_dir
		} else {
			dir = types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
			if dir.I_type().Get() == string(utiles.File) {return fmt.Errorf("dir inode is actually a file, cant be used for nesting other dirs"),types.IndexNode{}}
			permisions := self.User_allowed_actions(usr_id,usr_grp_id,&dir)
			if for_read && !permisions.Can_read(){
				return fmt.Errorf("read denied due lack of read permitions over this dir"),types.IndexNode{}
			}
			if for_write && !permisions.Can_write(){
				return fmt.Errorf("write denied due lack of write permitions over this dir"),types.IndexNode{}
			}
			dir.I_atime().Set(time)
		}
	}
	return nil,dir
}









func (self *Format)erase_inode_content_from_ptr_blck(pointer_block *types.PointerBlock,main_inode *types.IndexNode,level int32,inode_type utiles.InodeType,usr_id int32,usr_grp_id int32,time types.TimeHolder)bool {
	can_erase_self:=true
	if level != 0 {
		for i,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			can_erase := self.erase_inode_content_from_ptr_blck(&ptr_block,main_inode,level-1,inode_type,usr_id,usr_grp_id,time)
			if can_erase{
				bit := self.Block_section.Bit_no_for(ptr_block.Index)
				self.Block_bitmap.Erase(1,bit)
				pointer_block.B_pointers().No(int32(i)).Set(-1)
			}else{		
				can_erase_self = false
			}
		}
		return can_erase_self
	}
	for i,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		switch inode_type{
		case utiles.Directory:
			dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
			can_erase_block:=true
			for i, content := range dir_block.B_content().Get(){
				if content.B_inodo == -1 {continue}
				if content.B_name == OWN_DIR_NAME {continue}
				if content.B_name == PARENT_DIR_NAME {continue}
				sub_inode := types.CreateIndexNode(pointer_block.Super_service,content.B_inodo)
				permision := self.User_allowed_actions(usr_id,usr_grp_id,&sub_inode)
				can_erase_inode:=false
				if sub_inode.I_type().Get()== string(utiles.Directory){
					can_erase_inode = self.Erase_inode_content(sub_inode,utiles.Directory,usr_id,usr_grp_id,time)
				}else if permision.Can_write(){
					can_erase_inode = self.Erase_inode_content(sub_inode,utiles.File,usr_id,usr_grp_id,time)
					self.Update_dir_and_ancestors_size(*main_inode,-sub_inode.I_s().Get())
				}
				if permision.Can_write(){
					if can_erase_inode{
						bit := self.Inodes_section.Bit_no_for(content.B_inodo)
						self.Inodes_bitmap.Erase(1,bit)
						dir_block.B_content().No(int32(i)).B_inodo().Set(-1)
					} else{can_erase_block=false}
				}else{
					can_erase_block=false
				}
			}
			if can_erase_block{
				bit := self.Block_section.Bit_no_for(dir_block.Index)
				self.Block_bitmap.Erase(1,bit)
				pointer_block.B_pointers().No(int32(i)).Set(-1)
			}else{can_erase_self = false}
		case utiles.File:
			file_block := types.CreateFileBlock(self.super_service,ptr)
			bit := self.Block_section.Bit_no_for(file_block.Index)
			self.Block_bitmap.Erase(1,bit)
			pointer_block.B_pointers().No(int32(i)).Set(-1)
		}
	}
	return can_erase_self
}
func (self *Format)Erase_inode_content(inode types.IndexNode, inode_type utiles.InodeType,usr_id int32,usr_grp_id int32,time types.TimeHolder) bool {
	can_erase_self := true
	for n,ptr := range inode.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			can_erase := self.erase_inode_content_from_ptr_blck(&ptr_block,&inode,int32(n)-13,inode_type,usr_id,usr_grp_id,time)
			if can_erase{
				bit := self.Block_section.Bit_no_for(ptr_block.Index)
				self.Block_bitmap.Erase(1,bit)
				inode.I_block().No(int32(n)).Set(-1)
			}else{		
				can_erase_self = false
			}
			
		}else{
			switch inode_type{
			case utiles.Directory:
				dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
				can_erase_block:=true
				for i, content := range dir_block.B_content().Get(){
					if content.B_inodo == -1 {continue}
					if content.B_name == OWN_DIR_NAME {continue}
					if content.B_name == PARENT_DIR_NAME {continue}
					sub_inode := types.CreateIndexNode(inode.Super_service,content.B_inodo)
					permision := self.User_allowed_actions(usr_id,usr_grp_id,&sub_inode)
					can_erase_inode:=false
					if sub_inode.I_type().Get()== string(utiles.Directory){
						// fmt.Println(content.B_name)
						can_erase_inode = self.Erase_inode_content(sub_inode,utiles.Directory,usr_id,usr_grp_id,time)
					}else if permision.Can_write(){
						can_erase_inode = self.Erase_inode_content(sub_inode,utiles.File,usr_id,usr_grp_id,time)
						self.Update_dir_and_ancestors_size(inode,-sub_inode.I_s().Get())
					}
					if permision.Can_write(){
						if can_erase_inode{
							inode.I_mtime().Set(time)
							bit := self.Inodes_section.Bit_no_for(content.B_inodo)
							self.Inodes_bitmap.Erase(1,bit)
							dir_block.B_content().No(int32(i)).B_inodo().Set(-1)
						} else{can_erase_block=false}
					}else{
						can_erase_block=false
					}
				}
				if can_erase_block{
					bit := self.Block_section.Bit_no_for(dir_block.Index)
					self.Block_bitmap.Erase(1,bit)
					inode.I_block().No(int32(n)).Set(-1)
				}else{can_erase_self = false}
			case utiles.File:
				file_block := types.CreateFileBlock(self.super_service,ptr)
				bit := self.Block_section.Bit_no_for(file_block.Index)
				self.Block_bitmap.Erase(1,bit)
				inode.I_block().No(int32(n)).Set(-1)
			}
		}
	}
	return can_erase_self
}



func (self *Format)Remove_inode_if_possilbe(in_dir types.IndexNode,with_name [12]string,usr_id int32,usr_grp_id int32,time types.TimeHolder) bool {
	for n,ptr := range in_dir.I_block().Get(){
		if ptr == -1 { continue }
		if n >= 13{
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			found := self.renove_inode_if_possilbe_in_ptr_block(&ptr_block,&in_dir,int32(n-13),with_name,usr_id,usr_grp_id,time)
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
					sub_inode := types.CreateIndexNode(self.super_service,content.B_inodo)
					permission := self.User_allowed_actions(usr_id,usr_grp_id,&sub_inode)
					can_erase_inode:=false
					if sub_inode.I_type().Get()== string(utiles.Directory){
						// fmt.Println(content.B_name)
						can_erase_inode = self.Erase_inode_content(sub_inode,utiles.Directory,usr_id,usr_grp_id,time)
					}else if permission.Can_write(){
						can_erase_inode = self.Erase_inode_content(sub_inode,utiles.File,usr_id,usr_grp_id,time)
						self.Update_dir_and_ancestors_size(in_dir,-sub_inode.I_s().Get())
					}
					if permission.Can_write(){
						if can_erase_inode{
							in_dir.I_mtime().Set(time)
							bit := self.Inodes_section.Bit_no_for(content.B_inodo)
							self.Inodes_bitmap.Erase(1,bit)
							dir_block.B_content().No(int32(i)).B_inodo().Set(-1)
						} 
					}
					for _,cntnt := range dir_block.B_content().Get(){
						if cntnt.B_inodo != -1 {return true}
					}						
					self.Block_bitmap.Erase(1,self.Block_section.Bit_no_for(ptr))
					in_dir.I_block().No(int32(n)).Set(-1)
					return true
				}
			}
		}
	}
	return false
}








func (self *Format)renove_inode_if_possilbe_in_ptr_block(pointer_block *types.PointerBlock, main_inode *types.IndexNode,level int32, name [12]string,usr_id int32,usr_grp_id int32,time types.TimeHolder) bool {
	if level != 0 {
		for i,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			found := self.renove_inode_if_possilbe_in_ptr_block(&ptr_block,main_inode,level-1,name,usr_id,usr_grp_id,time)
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
					sub_inode := types.CreateIndexNode(self.super_service,content.B_inodo)
					permission := self.User_allowed_actions(usr_id,usr_grp_id,&sub_inode)
					can_erase_inode:=false
					if sub_inode.I_type().Get()== string(utiles.Directory){
						can_erase_inode = self.Erase_inode_content(sub_inode,utiles.Directory,usr_id,usr_grp_id,time)
					}else if permission.Can_write(){
						can_erase_inode = self.Erase_inode_content(sub_inode,utiles.File,usr_id,usr_grp_id,time)
						self.Update_dir_and_ancestors_size(*main_inode,-sub_inode.I_s().Get())
					}
					if permission.Can_write(){
						if can_erase_inode{
							bit := self.Inodes_section.Bit_no_for(content.B_inodo)
							self.Inodes_bitmap.Erase(1,bit)
							dir_block.B_content().No(int32(i)).B_inodo().Set(-1)
						} 
					}
					for _,cntnt := range dir_block.B_content().Get(){
						if cntnt.B_inodo != -1 {return true}
					}						
					self.Block_bitmap.Erase(1,self.Block_section.Bit_no_for(ptr))
					pointer_block.B_pointers().No(int32(n)).Set(-1)
					return true
				}
			}
	}
	return false

}


func (self *Format)Directory_deep_copy(dest_dir types.IndexNode,trgt_dir types.IndexNode,with_name [12]string,usr_id int32,usr_grp_id int32,time types.TimeHolder) bool {
	perm_to_dest := self.User_allowed_actions(usr_id,usr_grp_id,&dest_dir)
	if !perm_to_dest.Can_write(){return false}
	perm_to_trgt := self.User_allowed_actions(usr_id,usr_grp_id,&trgt_dir)
	if !perm_to_trgt.Can_read(){return false}
	// fmt.Println("=====================")
	// fmt.Println(trgt_dir.I_type().Get())
	switch trgt_dir.I_type().Get() {
	case string(utiles.File):
		new_inode := self.Put_in_dir(dest_dir,self.Wrap_holder_in_template(types.IndexNodeHolder{
			I_uid:   usr_id,
			I_gid:   usr_grp_id,
			I_s:     trgt_dir.I_s().Get(),
			I_atime: utiles.NO_TIME,
			I_ctime: time,
			I_mtime: time,
			I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
			I_type:  string(utiles.File),
			I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
		}),with_name)
		if new_inode.Index == -1 {return false}
		data:=self.Read_file(&trgt_dir)
		self.Update_file(&new_inode,0,data)

	case string(utiles.Directory):
		new_inode := self.Put_in_dir(dest_dir,self.Wrap_holder_in_template(types.IndexNodeHolder{
			I_uid:   usr_id,
			I_gid:   usr_grp_id,
			I_s:     trgt_dir.I_s().Get(),
			I_atime: utiles.NO_TIME,
			I_ctime: time,
			I_mtime: time,
			I_block: [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
			I_type:  string(utiles.Directory),
			I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
		}),with_name)
		self.Set_parent_child_relation(dest_dir,new_inode)
		if new_inode.Index == -1 {return false}
		sub_folders := self.Get_strict_shallow_tree_of_childs(trgt_dir)
		for _,folder:=range sub_folders{
			inode:=types.CreateIndexNode(folder.Super_service,folder.B_inodo().Get())
			self.Directory_deep_copy(new_inode,inode,folder.B_name().Get(),usr_id,usr_grp_id,time)
		}
	default:panic("Didnt match neither file nor directory")
	}

	return true
}



func (self *Format)Find_in_dir(in_dir types.IndexNode, name_criteria utiles.NameCriteria,usr_id int32,usr_grp_id int32,time types.TimeHolder) []types.Content {
	matched_files := make([]types.Content,0,10)
	perm_to_dest := self.User_allowed_actions(usr_id,usr_grp_id,&in_dir)
	if !perm_to_dest.Can_read(){return matched_files}
	in_dir.I_atime().Set(time)
	for _,content:= range self.Get_shallow_tree_of_childs(in_dir){
		if !name_criteria.Match(content.B_name().Get()){continue}
		matched_files = append(matched_files, content)
	}
	return matched_files
}
















func (self *Format)Change_owner(dest_dir types.IndexNode,usr_id int32,usr_grp_id int32,time types.TimeHolder,new_usr_id int32,new_usr_grp_id int32,recursive bool) int {
	changed:=0
	perm_to_dest := self.User_allowed_actions_with_ownership(usr_id,usr_grp_id,&dest_dir)
	if perm_to_dest.Is_owner(){
		dest_dir.I_uid().Set(new_usr_id)
		dest_dir.I_gid().Set(new_usr_grp_id)
		dest_dir.I_mtime().Set(time)
		changed++
	}
	if perm_to_dest.permision.Can_read(){
		if dest_dir.I_type().Get() == string(utiles.Directory) && recursive{
			dest_dir.I_atime().Set(time)
			sub_folders := self.Get_strict_shallow_tree_of_childs(dest_dir)
			for _,folder:=range sub_folders{
				inode:=types.CreateIndexNode(folder.Super_service,folder.B_inodo().Get())
				changed += self.Change_owner(inode,usr_id,usr_grp_id,time,new_usr_id,new_usr_grp_id,recursive)
			}
		}
	}
	return changed
}
func (self *Format)Change_ugo_permition(dest_dir types.IndexNode,time types.TimeHolder,recursive bool,new_perm [3]string) int {
	changed:=0
	dest_dir.I_mtime().Set(time)
	dest_dir.I_perm().Set(new_perm)
	changed++
	if dest_dir.I_type().Get() == string(utiles.Directory) && recursive{
		dest_dir.I_atime().Set(time)
		sub_folders := self.Get_strict_shallow_tree_of_childs(dest_dir)
		for _,folder:=range sub_folders{
			inode:=types.CreateIndexNode(folder.Super_service,folder.B_inodo().Get())
			changed += self.Change_ugo_permition(inode,time,recursive,new_perm)
		}
	}
	return changed
}



func (self *Format)Update_dir_and_ancestors_size(dir types.IndexNode,add_size int32) {
	// fmt.Printf("%d ============\n",dir.I_s().Get())
	dir.I_s().Set(dir.I_s().Get()+add_size)
	// fmt.Printf("%d ++++++++++++\n",dir.I_s().Get())
	dir_blck:=types.CreateDirectoryBlock(dir.Super_service,dir.I_block().No(0).Get())
	// fmt.Println(dir_blck.B_content().No(1).B_name().Get())
	parent:=types.CreateIndexNode(dir.Super_service,dir_blck.B_content().No(1).B_inodo().Get())
	if parent.Index == dir.Index {return}
	self.Update_dir_and_ancestors_size(parent,add_size)
}










