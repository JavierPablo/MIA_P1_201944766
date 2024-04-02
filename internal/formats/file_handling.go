package formats

import (
	"fmt"
	"project/internal/types"
)

func fill_fileblock(file_inode *types.FileBlock,at int32, data *[]string){
	char_spaces := file_inode.B_content().Spread()
	data_writen := 0
	for i := int(at); i < len(char_spaces) && data_writen < len(*data); i++ {
		char_spaces[i].Set((*data)[data_writen])
		data_writen++
	}
	*data = (*data)[data_writen:]
}

func (self *Format) erase_pointer_block_content(pointer_block *types.PointerBlock, level int32,till int32){
	pointers := pointer_block.Get().B_pointers
	if level != 0 { //recursive thing
		for i := till; i < 16; i++ {
			if pointers[i] == -1 {continue}
			pnt_blck := types.CreatePointerBlock(self.super_service,pointers[i])
			self.erase_pointer_block_content(&pnt_blck,level-1,0)
			bit := self.Block_section.Bit_no_for(pointers[i])
			self.Block_bitmap.Erase(1,bit)
			pointer_block.B_pointers().No(i).Set(-1)
		}
		return
	}
	for i := till; i < 16; i++ {
		if pointers[i] == -1 {continue}
		bit := self.Block_section.Bit_no_for(pointers[i])
		self.Block_bitmap.Erase(1,bit)
		pointer_block.B_pointers().No(i).Set(-1)
	}
}
func (self *Format) write_file_content_in_pointer(pointer_block *types.PointerBlock, level int32,at *int32, data *[]string){
	if level != 0 { //recursive thing
		for n,pointer_index := range pointer_block.Get().B_pointers{
			if len(*data) == 0{
				self.erase_pointer_block_content(pointer_block,level,int32(n))
				return
			}
			var ptr_blck types.PointerBlock
			if pointer_index == -1 {
				index, new_pointer_block := self.Create_PointerBlock()
				pointer_block.B_pointers().No(int32(n)).Set(index)
				ptr_blck = new_pointer_block
			}else{
				ptr_blck = types.CreatePointerBlock(self.super_service,pointer_index)
			}
			self.write_file_content_in_pointer(&ptr_blck,level-1,at,data)
		}
		return
	}
	
	for n,pointer_index := range pointer_block.Get().B_pointers{
		if *at > 64 { //skip
			*at =- 64
			continue
		}
		if len(*data) == 0{
			self.erase_pointer_block_content(pointer_block,level,int32(n))
			return
		}
		var file_block types.FileBlock
		if pointer_index == -1 {
			index, new_file_block := self.Create_FileBlock()
			pointer_block.B_pointers().No(int32(n)).Set(index)
			file_block = new_file_block
		}else{
			file_block = types.CreateFileBlock(self.super_service,pointer_index)
		}
		fill_fileblock(&file_block,*at,data)
		*at = 0;
	}
}
func (self *Format) Update_file(file_inode *types.IndexNode,at int32, data []string){
	current_size := file_inode.I_s().Get() - at
	expected_final_size := int32(len(data))
	if current_size < 0 {
		panic("Attempting to write data in file at unexpected far position")
	}
	
	for n,ptr := range file_inode.I_block().Get(){
		if len(data) == 0{
			ptrs := file_inode.I_block().Get()
			for i := n; i < 16; i++ {
				if ptrs[i] == -1 {continue}
				if i >= 13 {
					ptr_block := types.CreatePointerBlock(self.super_service,ptrs[i])
					self.erase_pointer_block_content(&ptr_block,int32(i - 13),0)
				}
				
				bit := self.Block_section.Bit_no_for(ptrs[i])
				// fmt.Printf("DEBUG: bit value is = %d\n",bit)
				self.Block_bitmap.Erase(1,bit)
				file_inode.I_block().No(int32(i)).Set(-1)
			}
			// file_inode.I_s().Set(expected_final_size)
			break
		}
		if n >= 13 {
			var ptr_block types.PointerBlock
			if ptr == -1 {
				ptr_index,new_ptr_block := self.Create_PointerBlock()
				file_inode.I_block().No(int32(n)).Set(ptr_index)
				ptr_block = new_ptr_block
			}else{
				ptr_block = types.CreatePointerBlock(self.super_service,ptr)
			}
			self.write_file_content_in_pointer(&ptr_block,int32(n-13),&at,&data)
		}else{
			if at > 64 { //skip
				at =- 64
				continue
			}
			var file_block types.FileBlock
			if ptr == -1 {
				ptr_index,new_ptr_block := self.Create_FileBlock()
				file_inode.I_block().No(int32(n)).Set(ptr_index)
				file_block = new_ptr_block
			}else{
				file_block = types.CreateFileBlock(self.super_service,ptr)
			}
			fill_fileblock(&file_block,at,&data)
			at = 0
		}
	}	
	if len(data) == 0{
		file_inode.I_s().Set(expected_final_size)
	}else{
		panic(fmt.Sprintf("Inode exausted and more space is needed. Bytes remaining = %d" ,len(data)))
	}
}






func fill_buff_from_fileblock(file_inode *types.FileBlock, amout int32,data *[]string){
	char_spaces := file_inode.B_content().Get()
	char_slice:= char_spaces[:amout]
	*data = append(*data, char_slice...)
}
func (self *Format) read_file_content_in_pointer(pointer_block *types.PointerBlock, level int32, amout *int32, data *[]string){
	if level != 0 { //recursive thing
		for _,pointer_index := range pointer_block.Get().B_pointers{
			if pointer_index == -1 {
				panic("Pointer should exist since bytes are left")
			}
			ptr_blck := types.CreatePointerBlock(self.super_service,pointer_index)
			self.read_file_content_in_pointer(&ptr_blck,level-1,amout,data)
			if *amout == 0{
				return
			}
		}
		return
	}
	
	for _,pointer_index := range pointer_block.Get().B_pointers{
		if pointer_index == -1 {
			panic("Pointer should exist since bytes are left")
		}
		file_block := types.CreateFileBlock(self.super_service,pointer_index)
		if *amout > 64 {
			fill_buff_from_fileblock(&file_block,64,data)
			*amout -= 64
		}else{
			fill_buff_from_fileblock(&file_block,*amout,data)
			*amout = 0 
		}
		if *amout == 0{
			return
		}
		
	}
}
func (self *Format) Read_file(file_inode *types.IndexNode)[]string{
	bytes_to_read := file_inode.I_s().Get()
	data := make([]string,0,100)
	for n,ptr := range file_inode.I_block().Get(){
		if bytes_to_read == 0{
			break
		}
		if n >= 13 {
			if ptr == -1 {
				panic("Pointer should exist since bytes are left")
			}
			ptr_block := types.CreatePointerBlock(self.super_service,ptr)
			self.read_file_content_in_pointer(&ptr_block,int32(n-13),&bytes_to_read,&data)
		}else{	
			if ptr == -1 {panic("Pointer should exist since bytes are left")}
			file_block := types.CreateFileBlock(self.super_service,ptr)
			if bytes_to_read > 64 {
				fill_buff_from_fileblock(&file_block,64,&data)
				bytes_to_read -= 64
			}else{
				fill_buff_from_fileblock(&file_block,bytes_to_read,&data)
				bytes_to_read = 0 
			}	
		}
	}	
	if bytes_to_read == 0{
		return data
	}else{
		panic("There are bytes left to read, and somehow that couldnt be done succesfully")
	}
}