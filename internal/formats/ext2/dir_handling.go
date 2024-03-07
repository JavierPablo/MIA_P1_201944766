package ext2

import (
	"fmt"
	"project/internal/types"
)

func (self *FormatEXT2)search_in_ptr_block(pointer_block *types.PointerBlock, level int32, name [12]string) (int32,types.IndexNode) {
	if level != 0 {
		for _,ptr := range pointer_block.B_pointers().Get(){
			if ptr == -1 { continue }
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			trgt_index,trgt_block := self.search_in_ptr_block(&ptr_block,level-1,name)
			if trgt_index != -1 {
				return trgt_index, trgt_block
			} 
		}
		return -1,types.CreateIndexNode(self.super_service,0)
	}
	for _,ptr := range pointer_block.B_pointers().Get(){
		if ptr == -1 { continue }
		dir_block := types.CreateDirectoryBlock(self.super_service,ptr)
		for _, content := range dir_block.B_content().Get(){
			if content.B_inodo == -1 {continue}
			if name == content.B_name {
				inode_result := types.CreateIndexNode(self.super_service,content.B_inodo)
				return content.B_inodo,inode_result
			}
		}
	}
	return -1,types.CreateIndexNode(self.super_service,0)

}
func (self *FormatEXT2)Search_for_inode(dir types.IndexNode, name [12]string) (int32,types.IndexNode) {
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
			for _, content := range dir_block.B_content().Get(){
				if content.B_inodo == -1 {continue}
				if name == content.B_name {
					inode_result := types.CreateIndexNode(self.super_service,content.B_inodo)
					return content.B_inodo,inode_result
				}
			}
		}
	}
	return -1,types.CreateIndexNode(self.super_service,0)
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
		fmt.Printf("Debug: is comming in second if\n")
		index,inode = self.try_put_in_new_dir_block_in_existing_ptr(dir,trgt_inode,name)
	}
	if index == -1 {
		fmt.Printf("Debug: is comming in last if\n")
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