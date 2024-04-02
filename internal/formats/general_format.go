package formats

import (
	"fmt"
	"project/internal/datamanagment"
	"project/internal/types"
	"project/internal/utiles"
	"strings"
)




type Format struct{
	Fit utiles.FitCriteria
	super_service *datamanagment.IOService
	Super_block types.SuperBlock
	Block_bitmap Bitmap
	Inodes_bitmap Bitmap
	Block_section Section
	Inodes_section Section
	journaling *JournalingManager
}
func (self *Format) Get_journaling()*JournalingManager{return self.journaling}
func (self *Format) Has_journaling()bool{return self.journaling != nil}
func (self *Format) Save_journal(ins Instruction){self.journaling.Push_instruction(ins)}
func (self *Format) Get_dot_journal_rep()string{
	return self.journaling.Generate_dot_rep()
}

func Format_new_fresh_FormatEXT2(super_service *datamanagment.IOService,fit utiles.FitCriteria,index int32,partition_size int32){
	super_block:=    types.CreateSuperBlock(super_service,index)

	var inodes_size int32 = types.CreateIndexNode(super_service,0).Size
	var blocks_size int32 =types.CreateFileBlock(super_service,0).Size

	n:=(partition_size-super_block.Size)/(int32(4)+inodes_size+blocks_size*int32(3))
	bm_inode_start:=index+super_block.Size
	bm_block_start:=bm_inode_start+n
	sec_inode_start:=bm_block_start+3*n
	sec_block_start:=sec_inode_start+inodes_size*n
	super_block.Set(types.SuperBlockHolder{
		S_filesystem_type:   int32(utiles.Ext2),
		S_inodes_count:      n,
		S_blocks_count:      n*3,
		S_free_blocks_count: n*3,
		S_free_inodes_count: n,
		S_mtime:             utiles.NO_TIME,
		S_umtime:            utiles.NO_TIME,
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
	
	format := Format{
		Fit: fit,
		super_service: super_service,
		Super_block:    super_block,
		Block_bitmap:  New_Bitmap(super_service,bm_block_start,3*n,super_block.S_free_blocks_count()),
		Inodes_bitmap:   New_Bitmap(super_service,bm_inode_start,n,super_block.S_free_inodes_count()),
		Block_section:  New_section(super_service,sec_block_start,3*n*blocks_size,blocks_size),
		Inodes_section: New_section(super_service,sec_inode_start,n*inodes_size,inodes_size),
		journaling: nil,
	}
	format.Block_bitmap.Clear()
	format.Inodes_bitmap.Clear()
	format.Block_section.Clear()
	format.Inodes_section.Clear()

	format.Init_bitmap_mapping()
	current_time:=utiles.Current_Time()
	_,root:=format.Create_Inode(types.IndexNodeHolder{
		I_uid:   1,
		I_gid:   1,
		I_s:     0,
		I_atime: utiles.NO_TIME,
		I_ctime: current_time,
		I_mtime: utiles.NO_TIME,
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.Directory),
		I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
	})
	format.Set_parent_child_relation(root,root)

	content:=strings.Split("1,G,root\n1,U,root,root,123\n","")
	user_file:=format.Put_in_dir(root,format.Wrap_holder_in_template(types.IndexNodeHolder{
		I_uid:   1,
		I_gid:   1,
		I_s:     int32(len(content)),
		I_atime: utiles.NO_TIME,
		I_ctime: current_time,
		I_mtime: utiles.NO_TIME,
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.File),
		I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
	}),utiles.Into_ArrayChar12("users.txt"))
	format.Update_file(&user_file,0,content)
}

const JOURNALING_SIZE int32 = 100
func Get_FormatEXT3_for_heal(super_service *datamanagment.IOService,index int32,fit utiles.FitCriteria)*Format{
	super_block:=    types.CreateSuperBlock(super_service,index)
	inodes_size:= super_block.S_inode_s().Get()
	blocks_size :=super_block.S_block_s().Get()
	
	bm_inode_start:=super_block.S_bm_inode_start().Get()
	if index+super_block.Size == bm_inode_start{return nil}
	bm_block_start:=super_block.S_bm_block_start().Get()
	sec_inode_start:=super_block.S_inode_start().Get()
	sec_block_start:=super_block.S_block_start().Get()

	inodes_count:=super_block.S_inodes_count().Get()
	blocks_count:=super_block.S_blocks_count().Get()
	
	inodes_free:=super_block.S_free_inodes_count()
	blocks_free:=super_block.S_free_blocks_count()

	super_block.S_mtime().Set(utiles.Current_Time())//maybe not here
	journal_manager:=New_Journaling(super_service,index+super_block.Size,bm_inode_start-(index+super_block.Size))
	journal_manager.Init_mapping(false)
	format := Format{
		Fit: fit,
		super_service: super_service,
		Super_block:    super_block,
		Block_bitmap:   New_Bitmap(super_service,bm_block_start,blocks_count,blocks_free),
		Inodes_bitmap:  New_Bitmap(super_service,bm_inode_start,inodes_count,inodes_free),
		
		Block_section:  New_section(super_service,sec_block_start,blocks_count*blocks_size,blocks_size),
		Inodes_section: New_section(super_service,sec_inode_start,inodes_count*inodes_size,inodes_size),

		journaling: &journal_manager,
	}
	format.Init_bitmap_mapping()
	current_time:=utiles.Current_Time()
	_,root:=format.Create_Inode(types.IndexNodeHolder{
		I_uid:   1,
		I_gid:   1,
		I_s:     0,
		I_atime: utiles.NO_TIME,
		I_ctime: current_time,
		I_mtime: utiles.NO_TIME,
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.Directory),
		I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
	})
	format.Set_parent_child_relation(root,root)

	content:=strings.Split("1,G,root\n1,U,root,root,123\n","")
	user_file:=format.Put_in_dir(root,format.Wrap_holder_in_template(types.IndexNodeHolder{
		I_uid:   1,
		I_gid:   1,
		I_s:     int32(len(content)),
		I_atime: utiles.NO_TIME,
		I_ctime: current_time,
		I_mtime: utiles.NO_TIME,
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.File),
		I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
	}),utiles.Into_ArrayChar12("users.txt"))
	format.Update_file(&user_file,0,content)

	return &format
	


}
func Format_new_fresh_FormatEXT3(super_service *datamanagment.IOService,fit utiles.FitCriteria,index int32,partition_size int32){
	super_block:=    types.CreateSuperBlock(super_service,index)

	var inodes_size int32 = types.CreateIndexNode(super_service,0).Size
	var blocks_size int32 =types.CreateFileBlock(super_service,0).Size

	n:=(partition_size-super_block.Size)/(int32(4)+inodes_size+blocks_size*int32(3) + JOURNALING_SIZE)
	bm_inode_start:=index+super_block.Size+n*JOURNALING_SIZE
	bm_block_start:=bm_inode_start+n
	sec_inode_start:=bm_block_start+3*n
	sec_block_start:=sec_inode_start+inodes_size*n
	super_block.Set(types.SuperBlockHolder{
		S_filesystem_type:   int32(utiles.Ext2),
		S_inodes_count:      n,
		S_blocks_count:      n*3,
		S_free_blocks_count: n*3,
		S_free_inodes_count: n,
		S_mtime:             utiles.NO_TIME,
		S_umtime:            utiles.NO_TIME,
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
	journal_manager:=New_Journaling(super_service,index+super_block.Size,bm_inode_start-(index+super_block.Size))
	journal_manager.Restart_count()

	format := Format{
		Fit: fit,
		super_service: super_service,
		Super_block:    super_block,
		Block_bitmap:  New_Bitmap(super_service,bm_block_start,3*n,super_block.S_free_blocks_count()),
		Inodes_bitmap:   New_Bitmap(super_service,bm_inode_start,n,super_block.S_free_inodes_count()),
		Block_section:  New_section(super_service,sec_block_start,3*n*blocks_size,blocks_size),
		Inodes_section: New_section(super_service,sec_inode_start,n*inodes_size,inodes_size),
		journaling: &journal_manager,
	}
	format.Block_bitmap.Clear()
	format.Inodes_bitmap.Clear()
	format.Block_section.Clear()
	format.Inodes_section.Clear()

	format.Init_bitmap_mapping()
	current_time:=utiles.Current_Time()
	_,root:=format.Create_Inode(types.IndexNodeHolder{
		I_uid:   1,
		I_gid:   1,
		I_s:     0,
		I_atime: utiles.NO_TIME,
		I_ctime: current_time,
		I_mtime: utiles.NO_TIME,
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.Directory),
		I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
	})
	format.Set_parent_child_relation(root,root)

	content:=strings.Split("1,G,root\n1,U,root,root,123\n","")
	user_file:=format.Put_in_dir(root,format.Wrap_holder_in_template(types.IndexNodeHolder{
		I_uid:   1,
		I_gid:   1,
		I_s:     int32(len(content)),
		I_atime: utiles.NO_TIME,
		I_ctime: current_time,
		I_mtime: utiles.NO_TIME,
		I_block: [16]int32{-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1,-1},
		I_type:  string(utiles.File),
		I_perm:  utiles.UGO_PERMITION_664.To_arr_string(),
	}),utiles.Into_ArrayChar12("users.txt"))
	format.Update_file(&user_file,0,content)
}
func Get_only_journaling(super_service *datamanagment.IOService,index int32)*JournalingManager{
	super_block:=    types.CreateSuperBlock(super_service,index)
	bm_inode_start:=super_block.S_bm_inode_start().Get()
	if index+super_block.Size != bm_inode_start{
		journal_manager:=New_Journaling(super_service,index+super_block.Size,bm_inode_start-(index+super_block.Size))
		journal_manager.Init_mapping(true)
		return &journal_manager
	}
	return nil
}
func Recover_Format(super_service *datamanagment.IOService,index int32, fit utiles.FitCriteria)Format{
	super_block:=    types.CreateSuperBlock(super_service,index)
	bm_inode_start:=super_block.S_bm_inode_start().Get()
	if index+super_block.Size == bm_inode_start{
		return Recover_FormatEXT2(super_service,index,fit)
	}else{
		return Recover_FormatEXT3(super_service,index,fit)
	}
}
func Recover_FormatEXT2(super_service *datamanagment.IOService,index int32, fit utiles.FitCriteria)Format{
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
	
	format := Format{
		Fit: fit,
		super_service: super_service,
		Super_block:    super_block,
		Block_bitmap:   New_Bitmap(super_service,bm_block_start,blocks_count,blocks_free),
		Inodes_bitmap:  New_Bitmap(super_service,bm_inode_start,inodes_count,inodes_free),
		
		Block_section:  New_section(super_service,sec_block_start,blocks_count*blocks_size,blocks_size),
		Inodes_section: New_section(super_service,sec_inode_start,inodes_count*inodes_size,inodes_size),

		journaling: nil,
	}
	return format
}
func Recover_FormatEXT3(super_service *datamanagment.IOService,index int32, fit utiles.FitCriteria)(Format){
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
	journal_manager:=New_Journaling(super_service,index+super_block.Size,bm_inode_start-(index+super_block.Size))
	journal_manager.Init_mapping(true)
	format := Format{
		Fit: fit,
		super_service: super_service,
		Super_block:    super_block,
		Block_bitmap:   New_Bitmap(super_service,bm_block_start,blocks_count,blocks_free),
		Inodes_bitmap:  New_Bitmap(super_service,bm_inode_start,inodes_count,inodes_free),
		
		Block_section:  New_section(super_service,sec_block_start,blocks_count*blocks_size,blocks_size),
		Inodes_section: New_section(super_service,sec_inode_start,inodes_count*inodes_size,inodes_size),

		journaling: &journal_manager,
	}

	return format
}






func (self *Format) Init_bitmap_mapping() {
	self.Block_bitmap.Init_mapping()
	self.Inodes_bitmap.Init_mapping()
}










func (self *Format) Create_DirectoryBlock() (int32,types.DirectoryBlock){
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
func (self *Format) Create_FileBlock() (int32,types.FileBlock){
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
func (self *Format) Create_PointerBlock() (int32,types.PointerBlock){
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






func (self *Format) First_Inode() (types.IndexNode){
	index := self.Super_block.S_firts_ino().Get()
	return types.CreateIndexNode(self.super_service,index)
}
func (self *Format) Create_Inode(inode_trgt types.IndexNodeHolder) (int32,types.IndexNode){
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






















func (self *Format) Log_block_bitmap(){
	free := self.Block_bitmap.free
	used := self.Block_bitmap.length-free
	fmt.Printf("\tblock: (occuped %d, free %d)\n",used,free)
	// fmt.Print("\t")
	// self.Block_bitmap.Log_bitmap_state()
}
func (self *Format) Log_inode_bitmap(){
	free := self.Inodes_bitmap.free
	used := self.Inodes_bitmap.length-free
	fmt.Printf("\tinodes: (occuped %d, free %d)\n",used,free)
	// fmt.Print("\t")
	// self.Inodes_bitmap.Log_bitmap_state()
}
func (self *Format) Prueba(){
	self.Block_bitmap.Best_fit(20)
	self.Super_block.S_free_blocks_count().Set(self.Block_bitmap.free)
}