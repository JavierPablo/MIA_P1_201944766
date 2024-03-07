package ext2

import (
	"fmt"
	"project/internal/datamanagment"
	"project/internal/types"
	"project/internal/utiles"
)

// import "math"



type FormatEXT2 struct{
	Fit utiles.FitCriteria
	super_service *datamanagment.IOService
	Super_block types.SuperBlock
	Block_bitmap Bitmap
	Inodes_bitmap Bitmap
	Block_section Section
	Inodes_section Section
}

func Format_new_FormatEXT2(super_service *datamanagment.IOService,fit utiles.FitCriteria,index int32,partition_size int32)FormatEXT2{
	super_block:=    types.CreateSuperBlock(super_service,index)

	var inodes_size int32 = types.CreateIndexNode(super_service,0).Size
	var blocks_size int32 =types.CreateFileBlock(super_service,0).Size

	n:=(partition_size-super_block.Size)/(int32(4)+inodes_size+blocks_size*int32(2))
	
	bm_inode_start:=index+super_block.Size
	bm_block_start:=bm_inode_start+n
	sec_inode_start:=bm_block_start+3*n
	sec_block_start:=sec_inode_start+inodes_size*n
	super_block.Set(types.SuperBlockHolder{
		S_filesystem_type:   2,
		S_inodes_count:      n,
		S_blocks_count:      n*3,
		S_free_blocks_count: n*3,
		S_free_inodes_count: n,
		S_mtime:             types.TimeHolder{},
		S_umtime:            types.TimeHolder{},
		S_mnt_count:         0,
		S_magic:             0,
		S_inode_s:           inodes_size,
		S_block_s:           blocks_size,
		S_firts_ino:         sec_inode_start,
		S_first_blo:         sec_block_start,
		S_bm_inode_start:    bm_inode_start,
		S_bm_block_start:    bm_block_start,
		S_inode_start:       sec_inode_start,
		S_block_start:       sec_block_start,
	})
	
	format := FormatEXT2{
		Fit: fit,
		super_service: super_service,
		Super_block:    super_block,
		Block_bitmap:  New_Bitmap(super_service,bm_block_start,3*n,super_block.S_free_blocks_count()),
		Inodes_bitmap:   New_Bitmap(super_service,bm_inode_start,n,super_block.S_free_inodes_count()),
		Block_section:  New_section(super_service,sec_block_start,3*n*blocks_size,blocks_size),
		Inodes_section: New_section(super_service,sec_inode_start,n*inodes_size,inodes_size),
	}
	format.Block_bitmap.Clear()
	format.Inodes_bitmap.Clear()
	format.Block_section.Clear()
	format.Inodes_section.Clear()

	format.Block_bitmap.Init_mapping()
	format.Inodes_bitmap.Init_mapping()
	
	return format
}
func Recover_FormatEXT2(super_service *datamanagment.IOService,index int32, fit utiles.FitCriteria)FormatEXT2{
	super_block:=    types.CreateSuperBlock(super_service,index)

	inodes_size:= super_block.S_inode_s().Get()
	blocks_size :=super_block.S_block_s().Get()
	
	bm_inode_start:=super_block.S_bm_inode_start().Get()
	bm_block_start:=super_block.S_bm_block_start().Get()
	sec_inode_start:=super_block.S_inode_start().Get()
	sec_block_start:=super_block.S_block_start().Get()

	inodes_count:=super_block.S_inodes_count().Get()
	blocks_count:=super_block.S_blocks_count().Get()
	
	inodes_free:=super_block.S_free_inodes_count()
	blocks_free:=super_block.S_free_blocks_count()

	super_block.S_mtime().Set(utiles.Current_Time())//maybe not here
	
	format := FormatEXT2{
		Fit: fit,
		super_service: super_service,
		Super_block:    super_block,
		Block_bitmap:   New_Bitmap(super_service,bm_block_start,blocks_count,blocks_free),
		Inodes_bitmap:  New_Bitmap(super_service,bm_inode_start,inodes_count,inodes_free),
		
		Block_section:  New_section(super_service,sec_block_start,blocks_count*blocks_size,blocks_size),
		Inodes_section: New_section(super_service,sec_inode_start,inodes_count*inodes_size,inodes_size),
	}

	format.Block_bitmap.Init_mapping()
	format.Inodes_bitmap.Init_mapping()

	return format
}











func (self *FormatEXT2) Create_DirectoryBlock() (int32,types.DirectoryBlock){
	index_result := int32(0)
	if self.Fit == utiles.Best{
		index_result = self.Block_bitmap.Best_fit(1)
	}else if self.Fit == utiles.First{
		index_result = self.Block_bitmap.First_fit(1)
	}else if self.Fit == utiles.Worst{
		index_result = self.Block_bitmap.Worst_fit(1)
	}
	if index_result == -1{panic("Empty space in block bitmap")}
	abs_index := self.Block_section.Index_for(index_result)
	dir_block := types.CreateDirectoryBlock(self.super_service,abs_index)
	for _, v := range dir_block.B_content().Spread() {
		v.B_inodo().Set(-1)	
	}
	return abs_index,dir_block
}
func (self *FormatEXT2) Create_FileBlock() (int32,types.FileBlock){
	index_result := int32(0)
	if self.Fit == utiles.Best{
		index_result = self.Block_bitmap.Best_fit(1)
	}else if self.Fit == utiles.First{
		index_result = self.Block_bitmap.First_fit(1)
	}else if self.Fit == utiles.Worst{
		index_result = self.Block_bitmap.Worst_fit(1)
	}
	if index_result == -1{panic("Empty space in block bitmap")}
	abs_index := self.Block_section.Index_for(index_result)
	file_block := types.CreateFileBlock(self.super_service,abs_index)
	return abs_index, file_block
}
func (self *FormatEXT2) Create_PointerBlock() (int32,types.PointerBlock){
	index_result := int32(0)
	if self.Fit == utiles.Best{
		index_result = self.Block_bitmap.Best_fit(1)
	}else if self.Fit == utiles.First{
		index_result = self.Block_bitmap.First_fit(1)
	}else if self.Fit == utiles.Worst{
		index_result = self.Block_bitmap.Worst_fit(1)
	}
	if index_result == -1{panic("Empty space in block bitmap")}
	abs_index := self.Block_section.Index_for(index_result)
	// {-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1}
	pointer_block := types.CreatePointerBlock(self.super_service,abs_index)
	pointer_block.Set(types.PointerBlockHolder{
		B_pointers: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
	})
	return abs_index, pointer_block
}






func (self *FormatEXT2) First_Inode() (int32,types.IndexNode){
	index := self.Super_block.S_firts_ino().Get()
	return index, types.CreateIndexNode(self.super_service,index)
}
func (self *FormatEXT2) Create_Inode(inode_trgt types.IndexNodeHolder) (int32,types.IndexNode){
	index_result := int32(0)
	if self.Fit == utiles.Best{
		index_result = self.Inodes_bitmap.Best_fit(1)
	}else if self.Fit == utiles.First{
		index_result = self.Inodes_bitmap.First_fit(1)
	}else if self.Fit == utiles.Worst{
		index_result = self.Inodes_bitmap.Worst_fit(1)
	}
	if index_result == -1{panic("Empty space in inode bitmap")}
	abs_index := self.Inodes_section.Index_for(index_result)
	
	inode := types.CreateIndexNode(self.super_service,abs_index)
	inode.Set(inode_trgt)
	return abs_index, inode
}














func (self *FormatEXT2) Log_block_bitmap(){
	free := self.Block_bitmap.free
	used := self.Block_bitmap.length-free
	fmt.Printf("\tblock: (occuped %d, free %d)\n",used,free)
	// fmt.Print("\t")
	// self.Block_bitmap.Log_bitmap_state()
}
func (self *FormatEXT2) Log_inode_bitmap(){
	free := self.Inodes_bitmap.free
	used := self.Inodes_bitmap.length-free
	fmt.Printf("\tinodes: (occuped %d, free %d)\n",used,free)
	// fmt.Print("\t")
	// self.Inodes_bitmap.Log_bitmap_state()
}
func (self *FormatEXT2) Prueba(){
	self.Block_bitmap.Best_fit(20)
	self.Super_block.S_free_blocks_count().Set(self.Block_bitmap.free)
}