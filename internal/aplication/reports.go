package aplication

import (
	"fmt"
	"project/internal/datamanagment"
	"project/internal/formats"
	"project/internal/types"
	"project/internal/utiles"
	"sort"
	"strings"
)
func (self *Aplication) Mbr_repo(ioservice *datamanagment.IOService)(string,error){
	mbr:=types.CreateMasterBootRecord(ioservice,0)
	ebrs:= make([]string,0,10)
	var extended_partition *types.Partition = nil
	for _,part:= range mbr.Mbr_partitions().Spread(){
		if part.Part_type().Get() != string(utiles.Extendend) {
			continue
		}
		if part.Part_start().Get() == -1{continue}
		extended_partition = &part

	}
	if extended_partition != nil {
		begin := extended_partition.Part_start().Get()
		ebr := types.CreateExtendedBootRecord(extended_partition.Super_service,begin)
		next_ebr_index := ebr.Part_next().Get()
	
		if ebr.Part_start().Get() != -1{
			ebrs = append(ebrs, ebr.Dot_label())
		}
		for next_ebr_index != -1{
			ebr = types.CreateExtendedBootRecord(ebr.Super_service,next_ebr_index)
			next_ebr_index = ebr.Part_next().Get()
			if ebr.Part_start().Get() != -1{
				ebrs = append(ebrs, ebr.Dot_label())
			}
		}

		ebr_nodes := ""
		for i := 0; i < len(ebrs); i++ {
			ebr_nodes = fmt.Sprintf("y%d[color=black,shape=box,label=<%s>];\n",i,ebrs[i]) + ebr_nodes
		}
		return fmt.Sprintf(`
		digraph G {
			subgraph cluster_0{
				label = "Reporte MBR";
				node [color=white,label=<%s>];
				a0;
			}
			subgraph cluster_1{
				label = "Reporte EBR";
				%s
			}
		}`,mbr.Dot_label(),ebr_nodes), nil
	}

	return fmt.Sprintf(`
	digraph G {
		subgraph cluster_0{
			label = "Reporte MBR";
			node [color=white,label=<%s>];
			a0;
		}
	}`,mbr.Dot_label()), nil
}



func (self *Aplication) Disk_repos(ioservice *datamanagment.IOService)(string,error){
	type LabeledSpace struct{
		label string
		space datamanagment.Space
	}
	new:=func(index int32, size int32,title string)LabeledSpace{
		return LabeledSpace{
			label: title,
			space: datamanagment.New_Space(index,size),
		}
	}
	mbr:=types.CreateMasterBootRecord(ioservice,0)
	var extend_space *datamanagment.Space = nil
	spaces:=[]LabeledSpace{new(0,mbr.Size,"MBR")}
	ext_spaces:=[]LabeledSpace{}
	for _,part:= range mbr.Mbr_partitions().Spread(){
		if part.Part_start().Get() == -1{continue}
		if part.Part_type().Get() == string(utiles.Extendend) {
			begin := part.Part_start().Get()
			some:=datamanagment.New_Space(begin,part.Part_s().Get())
			extend_space = &some
			ebr := types.CreateExtendedBootRecord(part.Super_service,begin)
			next_ebr_index := ebr.Part_next().Get()
			if ebr.Part_start().Get() != -1{
				ext_spaces = append(ext_spaces,new(0,ebr.Size,"EBR"))
				ext_spaces = append(ext_spaces,new(ebr.Part_start().Get()-begin,ebr.Part_s().Get(),"Logica"))
			}
			for next_ebr_index != -1{
				ebr = types.CreateExtendedBootRecord(ebr.Super_service,next_ebr_index)
				ext_spaces = append(ext_spaces,new(next_ebr_index-begin,ebr.Size,"EBR"))
				if ebr.Part_start().Get() != -1{
					ext_spaces = append(ext_spaces,new(ebr.Part_start().Get()-begin,ebr.Part_s().Get(),"Logica"))
				}
				next_ebr_index = ebr.Part_next().Get()
			}					
		}else{
			spaces = append(spaces,new(part.Part_start().Get(),part.Part_s().Get(),"Primaria"))
		}
	}
	if extend_space !=nil{
		spacemanger:=datamanagment.SpaceManager_from_free_spaces([]datamanagment.Space{datamanagment.New_Space(0,extend_space.Length)},extend_space.Length)
		for _,s:=range ext_spaces{
			err:=spacemanger.Ocupe_raw_space(s.space.Length,s.space.Index)
			if err!=nil{
				return "",err
			}
		}
		for _,free:=range spacemanger.Free_spaces(){
			ext_spaces = append(ext_spaces, new(free.Index,free.Length,"Free"))
		}
		for i := 0; i < len(ext_spaces); i++ {
			ext_spaces[i].space.Index+=extend_space.Index
		}
		spaces = append(spaces, ext_spaces...)
	}

	disk_total:=mbr.Mbr_tamano().Get()

	// for i := 0; i < len(spaces); i++ {
	// 	fmt.Print(spaces[i].label)
	// 	fmt.Println(spaces[i].space.Show())
	// }
	// fmt.Println(disk_total)
	spacemanger:=datamanagment.SpaceManager_from_free_spaces([]datamanagment.Space{datamanagment.New_Space(0,disk_total)},disk_total)
	for _,s:=range spaces{
		err:=spacemanger.Ocupe_raw_space(s.space.Length,s.space.Index)
		if err!=nil{
			return "",err
		}
	}
	for _,free:=range spacemanger.Free_spaces(){
		spaces = append(spaces, new(free.Index,free.Length,"Free"))
	}
	sort.SliceStable(spaces,func(i, j int) bool {
		return spaces[i].space.Index < spaces[j].space.Index
	})
	partition_cols_ant:=""
	partition_cols_pos:=""
	logic_cols:=""
	span_counter:=0
	for i := 0; i < len(spaces); i++ {
		porcentage:=float32(100*spaces[i].space.Length)/float32(disk_total)
		if extend_space != nil {
			if extend_space.Index <= spaces[i].space.Index && extend_space.Boundary() > spaces[i].space.Index{
				logic_cols +=fmt.Sprintf("<TD BORDER=\"1\" >%s<BR/>%f %%</TD>\n",spaces[i].label,porcentage)
				span_counter++
			}else if extend_space.Index > spaces[i].space.Index{
				partition_cols_ant+=fmt.Sprintf("<TD BORDER=\"1\" COLOR=\" #135F8A\">%s<BR/>%f %%</TD>\n",spaces[i].label,porcentage)
			}else{
				partition_cols_pos+=fmt.Sprintf("<TD BORDER=\"1\" COLOR=\" #135F8A\">%s<BR/>%f %%</TD>\n",spaces[i].label,porcentage)

			}
		}else{
			partition_cols_ant+=fmt.Sprintf("<TD BORDER=\"1\" COLOR=\" #135F8A\">%s<BR/>%f %%</TD>\n",spaces[i].label,porcentage)
		}
	}
	dot_str:=""
	if extend_space == nil{
		dot_str = fmt.Sprintf(`
		digraph G {
		y2[color=black,shape=box,label=<<TABLE BGCOLOR=" #45B3F1" BORDER="1"  COLOR="BLACK">
			<TR>
				%s 
			</TR>	
		</TABLE>>];
		}`,partition_cols_ant)
		
	}else{
		dot_str = fmt.Sprintf(`
		digraph G {
		y2[color=black,shape=box,label=<<TABLE BGCOLOR=" #45B3F1" BORDER="1"  COLOR="BLACK">
			<TR>
				%s 
				<TD BORDER="1" COLOR=" #135F8A"> 
					<TABLE BORDER="0"> 
						<TR>
							<TD COLSPAN="%d" BORDER="1" BGCOLOR="WHITE"> Extendida </TD>
						</TR> 
						<TR>
							%s
						</TR>
					</TABLE>
				</TD>
				%s
			</TR>	
		</TABLE>>];
		}`,partition_cols_ant,span_counter,logic_cols,partition_cols_pos)
	}
	return dot_str,nil
}



func (self *Aplication) Inode_repos(part_id string)(string,error){
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).Id != part_id{continue}
		mounted := &self.mounted_partitions[i]
		io_service:=mounted.io
		var super_block_start int32
		var err error
		var fit utiles.FitCriteria
		
		switch mounted.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
			super_block_start = logic.Part_start().Get()
			fit,err = utiles.Translate_fit(logic.Part_fit().Get())
			if err!=nil{return "",err}
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
			
			fit,err = utiles.Translate_fit(partition.Part_fit().Get())
			if err!=nil{return "",err}
		}
		format := formats.Recover_Format(io_service,super_block_start,fit)
		format.Init_bitmap_mapping()
		spaceman:=format.Inodes_bitmap.Get_SpaceManager()
		spaceman,err=datamanagment.SpaceManager_from_occuped_spaces(spaceman.Free_spaces(),spaceman.Get_length())
		if err !=nil{return "",err}
		builder:=""
		counter:=0
		for _,occuped_space:=range spaceman.Free_spaces(){
			for i := occuped_space.Index; i < occuped_space.Length; i++ {
				inode_abs_index:=format.Inodes_section.Index_for(i)
				inode:=types.CreateIndexNode(io_service,inode_abs_index)
				builder += fmt.Sprintf("subgraph cluster_%d {label=\"Inode %d\\n index=%d\";a%d[color=white,label=<%s>];}\n",counter,i,inode_abs_index,counter,inode.Dot_label())
				counter++
			}
		}
		return fmt.Sprintf("graph{%s}",builder),nil
	}
	return "",fmt.Errorf("there's no partition with this id")
}







type BlockRepoBuilder struct{
	acumulator string
	dot_block_counter int
	format *formats.Format
}
func(self *BlockRepoBuilder) IoService()*datamanagment.IOService{
	return self.format.Super_block.Super_service
}
func(self *BlockRepoBuilder) Build()string{
	return fmt.Sprintf("graph{%s}",self.acumulator)
}
func(self *BlockRepoBuilder) AppendDirectoryBlock(blck_dot string,block_abs_index int32){
	self.appendBlock(blck_dot,block_abs_index,"Directory")
}
func(self *BlockRepoBuilder) AppendFileBlock(blck_dot string,block_abs_index int32){
	self.appendBlock(blck_dot,block_abs_index,"File")
}
func(self *BlockRepoBuilder) AppendPointerBlock(blck_dot string,block_abs_index int32){
	self.appendBlock(blck_dot,block_abs_index,"Pointer")
}
func(self *BlockRepoBuilder) appendBlock(blck_dot string,block_abs_index int32,name string){
	self.acumulator += fmt.Sprintf("subgraph cluster_%d {label=\"%s Block %d\\n index=%d\";a%d[color=white,label=<%s>];}\n",
	self.dot_block_counter,
	name,
	self.format.Block_section.Bit_no_for(block_abs_index),
	block_abs_index,
	self.dot_block_counter,
	blck_dot)
	self.dot_block_counter++
}



func Recursive_block_repo(root_inode *types.IndexNode,inode_type utiles.InodeType,builder *BlockRepoBuilder){
	switch inode_type{
	case (utiles.Directory):
		for n,ptr := range root_inode.I_block().Get(){
			if ptr == -1 {continue}
			if n >= 13 {
				ptr_block := types.CreatePointerBlock(builder.IoService(),ptr)
				builder.AppendPointerBlock(ptr_block.Dot_label(),ptr)
				recursive_block_repo_in_pointer(&ptr_block,int32(n-13),utiles.Directory,builder)
			}else{	
				dir_block := types.CreateDirectoryBlock(builder.IoService(),ptr)
				builder.AppendDirectoryBlock(dir_block.Dot_label(),ptr)
				for _,cont:=range dir_block.B_content().Get(){
					if cont.B_inodo == -1{continue}
					another_inode:=types.CreateIndexNode(builder.IoService(),cont.B_inodo)
					if cont.B_name == formats.OWN_DIR_NAME {continue}
					if cont.B_name == formats.PARENT_DIR_NAME {continue}
					switch another_inode.I_type().Get(){
					case string(utiles.Directory):
						Recursive_block_repo(&another_inode,utiles.Directory,builder)
					case string(utiles.File):
						Recursive_block_repo(&another_inode,utiles.File,builder)
					default: panic("should have matched one of them")
					} 
				}
			}
		}	
	case (utiles.File):
		for n,ptr := range root_inode.I_block().Get(){
			if ptr == -1 {continue}
			if n >= 13 {
				ptr_block := types.CreatePointerBlock(builder.IoService(),ptr)
				builder.AppendPointerBlock(ptr_block.Dot_label(),ptr)
				recursive_block_repo_in_pointer(&ptr_block,int32(n-13),utiles.File,builder)
			}else{	
				file_block := types.CreateFileBlock(builder.IoService(),ptr)
				builder.AppendFileBlock(file_block.Dot_label(),ptr)
			}
		}	
		default: panic("Not expected llllllllllllllllllllllllllllllllll")
	}
}

func recursive_block_repo_in_pointer(pointer_block *types.PointerBlock, level int32, inode_type utiles.InodeType, builder *BlockRepoBuilder){
	if level != 0 { //recursive thing
		for _,pointer_index := range pointer_block.Get().B_pointers{
			if pointer_index == -1 {continue}
			ptr_blck := types.CreatePointerBlock(builder.IoService(),pointer_index)
			builder.AppendPointerBlock(ptr_blck.Dot_label(),pointer_index)
			recursive_block_repo_in_pointer(&ptr_blck,level-1,inode_type,builder)
		}
		return
	}
	switch inode_type{
	case (utiles.Directory):
		for _,pointer_index := range pointer_block.B_pointers().Get(){
			if pointer_index == -1 {continue}
			dir_block := types.CreateDirectoryBlock(builder.IoService(),pointer_index)
			builder.AppendDirectoryBlock(dir_block.Dot_label(),pointer_index)
			for _,cont:=range dir_block.B_content().Get(){
				if cont.B_inodo == -1{continue}
				another_inode:=types.CreateIndexNode(builder.IoService(),cont.B_inodo)
				switch another_inode.I_type().Get(){
				case string(utiles.Directory):
					if cont.B_name == formats.OWN_DIR_NAME {continue}
					if cont.B_name == formats.PARENT_DIR_NAME {continue}
					Recursive_block_repo(&another_inode,utiles.Directory,builder)
				case string(utiles.File):
					Recursive_block_repo(&another_inode,utiles.File,builder)
				default: panic("should have matched one of them")
				} 
			}
		}
	case (utiles.File):
		for _,pointer_index := range pointer_block.Get().B_pointers{
			if pointer_index == -1 {continue}
			file_block := types.CreateFileBlock(builder.IoService(),pointer_index)
			builder.AppendFileBlock(file_block.Dot_label(),pointer_index)
		}
	}
}


func (self *Aplication) Block_repos(part_id string)(string,error){
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).Id != part_id{continue}
		mounted := &self.mounted_partitions[i]
		io_service:=mounted.io
		var super_block_start int32
		var err error
		var fit utiles.FitCriteria
		switch mounted.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
			super_block_start = logic.Part_start().Get()
			fit,err = utiles.Translate_fit(logic.Part_fit().Get())
			if err!=nil{return "",err}
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
			
			fit,err = utiles.Translate_fit(partition.Part_fit().Get())
			if err!=nil{return "",err}
		}
		format := formats.Recover_Format(io_service,super_block_start,fit)
		format.Init_bitmap_mapping()
		root_inode:=format.First_Inode()
		builder:=BlockRepoBuilder{
			acumulator:        "",
			dot_block_counter: 0,
			format:            &format,
		}
		// builder.AppendDirectoryBlock(root_inode.Dot_label(),root_inode.Index)
		Recursive_block_repo(&root_inode,utiles.Directory,&builder)
		return builder.Build(),nil
	}
	return "",fmt.Errorf("there's no partition with this id")
}


func (self *Aplication) Inode_bitmap_repos(part_id string)(string,error){
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).Id != part_id{continue}
		mounted := &self.mounted_partitions[i]
		io_service:=mounted.io
		var super_block_start int32
		var err error
		var fit utiles.FitCriteria
		
		switch mounted.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
			super_block_start = logic.Part_start().Get()
			fit,err = utiles.Translate_fit(logic.Part_fit().Get())
			if err!=nil{return "",err}
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
			
			fit,err = utiles.Translate_fit(partition.Part_fit().Get())
			if err!=nil{return "",err}
		}
		format := formats.Recover_Format(io_service,super_block_start,fit)
		format.Init_bitmap_mapping()
		space_man:=format.Inodes_bitmap.Get_SpaceManager()
		bitmap:=make([]byte,0,1024)
		char_counter:=1

		for i := 0; i < int(space_man.Get_length()); i++ {
			bitmap = append(bitmap,  byte(49))
			if char_counter%20==0{
				bitmap = append(bitmap,  byte(10))
				char_counter = 1
			}else{
				char_counter++
			}
		}
		for _,free_space:=range space_man.Free_spaces(){
			till:=free_space.Boundary()

			for i := free_space.Index; i < till; i++ {
				offset:=i/20
				bitmap[i+offset]=byte(48)
			}
		}
		return string(bitmap),nil
	}
	return "",fmt.Errorf("there's no partition with this id")
}
func (self *Aplication) Block_bitmap_repos(part_id string)(string,error){
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).Id != part_id{continue}
		mounted := &self.mounted_partitions[i]
		io_service:=mounted.io
		var super_block_start int32
		var err error
		var fit utiles.FitCriteria
		
		switch mounted.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
			super_block_start = logic.Part_start().Get()
			fit,err = utiles.Translate_fit(logic.Part_fit().Get())
			if err!=nil{return "",err}
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
			
			fit,err = utiles.Translate_fit(partition.Part_fit().Get())
			if err!=nil{return "",err}
		}
		format := formats.Recover_Format(io_service,super_block_start,fit)
		format.Init_bitmap_mapping()
		space_man:=format.Block_bitmap.Get_SpaceManager()
		bitmap:=make([]byte,0,1024)
		char_counter:=1

		for i := 0; i < int(space_man.Get_length()); i++ {
			bitmap = append(bitmap,  byte(49))
			if char_counter%20==0{
				bitmap = append(bitmap,  byte(10))
				char_counter = 1
			}else{
				char_counter++
			}
		}
		for _,free_space:=range space_man.Free_spaces(){
			till:=free_space.Boundary()

			for i := free_space.Index; i < till; i++ {
				offset:=i/20
				bitmap[i+offset]=byte(48)
			}
		}
		return string(bitmap),nil
	}
	return "",fmt.Errorf("there's no partition with this id")
}







type TreeRepoBuilder struct{
	nodes string
	connections string
	nodes_counter int
	format *formats.Format
}
func(self *TreeRepoBuilder) IoService()*datamanagment.IOService{
	return self.format.Super_block.Super_service
}
func(self *TreeRepoBuilder) Build()string{
	return fmt.Sprintf(`
	digraph G {rankdir=LR;subgraph cluster_0 {
		label="Reporte";node[shape=Mrecord,style=filled];
			%s		
			%s
	}}`,self.nodes,self.connections)
}


type InodeNodeBuilder struct{
	inode_type utiles.InodeType
	index int32
	name string
	id string
	rows string
	connections string
}
func (self *InodeNodeBuilder)Connect(no int,to string){
	self.connections+= fmt.Sprintf("%s:p%d->%s:h;\n",self.id,no,to)
}
func (self *InodeNodeBuilder)Push_row(no int, index int32){
	self.rows +=fmt.Sprintf("|{%d|<p%d>%d}",no,no,index)
}
func(self *TreeRepoBuilder) Collect_Inode_node(inode *InodeNodeBuilder){
	self.nodes+=fmt.Sprintf("%s[label=\"{/%s}|{<h>Inodo|%d}%s\",fillcolor=\"#F89152\"];\n",
	inode.id,
	inode.name,
	self.format.Inodes_section.Bit_no_for(inode.index),
	inode.rows)
	self.connections+=inode.connections
}
func(self *TreeRepoBuilder) New_Inode_node(name string, inode types.IndexNode)(InodeNodeBuilder){
	node_id := fmt.Sprintf(`node%d`,self.nodes_counter)
	self.nodes_counter++
	i_type:=utiles.File
	if inode.I_type().Get() == string(utiles.Directory){
		i_type=utiles.Directory
	}
	return InodeNodeBuilder{
		inode_type:  i_type,
		index:       inode.Index,
		name:        name,
		id:          node_id,
		rows:        "",
		connections: "",
	}
}
type DirBlockNodeBuilder struct{
	index int32
	id string
	rows string
	connections string
}
func (self *DirBlockNodeBuilder)Connect(no int,to string){
	self.connections+= fmt.Sprintf("%s:p%d->%s:h;\n",self.id,no,to)
}
func (self *DirBlockNodeBuilder)Push_row(no int, name string,index int32){
	self.rows +=fmt.Sprintf("|{%s|<p%d>%d}",name,no,index)
}
func(self *TreeRepoBuilder) Collect_DirBlock_node(dir_blck *DirBlockNodeBuilder){
	self.nodes+=fmt.Sprintf("%s[label=\"{<h>Dir Block|%d}%s\",fillcolor=\"#84F852\"];\n",
	dir_blck.id,
	self.format.Block_section.Bit_no_for(dir_blck.index),
	dir_blck.rows)
	self.connections+=dir_blck.connections
}
func(self *TreeRepoBuilder) New_DirBLock_node(dir_block types.DirectoryBlock)(DirBlockNodeBuilder){
	node_id := fmt.Sprintf(`node%d`,self.nodes_counter)
	self.nodes_counter++
	return DirBlockNodeBuilder{
		index:       dir_block.Index,
		id:          node_id,
		rows:        "",
		connections: "",
	}
}
type PntrBlockNodeBuilder struct{
	inode_type utiles.InodeType
	index int32
	id string
	rows string
	connections string
}
func (self *PntrBlockNodeBuilder)Connect(no int,to string){
	self.connections+= fmt.Sprintf("%s:p%d->%s:h;\n",self.id,no,to)
}
func (self *PntrBlockNodeBuilder)Push_row(no int, index int32){
	self.rows +=fmt.Sprintf("|{%d|<p%d>%d}",no,no,index)
}
func(self *TreeRepoBuilder) Collect_PntrBlock_node(pntr_blck *PntrBlockNodeBuilder){
	self.nodes+=fmt.Sprintf("%s[label=\"{<h>Pntr Block|%d}%s\",fillcolor=\"#E1F852\"];\n",
	pntr_blck.id,
	self.format.Block_section.Bit_no_for(pntr_blck.index),
	pntr_blck.rows)
	self.connections+=pntr_blck.connections
}
func(self *TreeRepoBuilder) New_PntrBlock_node(pntr_block types.PointerBlock,inode_type utiles.InodeType)(PntrBlockNodeBuilder){
	node_id := fmt.Sprintf(`node%d`,self.nodes_counter)
	self.nodes_counter++
	return PntrBlockNodeBuilder{
		inode_type: inode_type,
		index:       pntr_block.Index,
		id:          node_id,
		rows:        "",
		connections: "",
	}
}
func(self *TreeRepoBuilder) New_FileBlock_node(file_block types.FileBlock)(string){
	node_id := fmt.Sprintf(`node%d`,self.nodes_counter)
	self.nodes_counter++
	content:=""
	counter:=1
	for _,char:= range file_block.B_content().Get(){
		ascii_val:=char[0]
		if ascii_val >= 32 && ascii_val <= 126{
			switch ascii_val{
			case 34:
				char = `\"`
			case 92:
				char = `\\`
			case 123:
				char = `\{`
			case 125:
				char = `\}`
			case 124:
				char = `\|`
			case 60:
				char = `\<`
			case 62:
				char = `\>`
			}
		}else{
			switch ascii_val{
			case 10:
				char = `\\n`
			case 11:
				char = `\\t`
			default: 
				char = "-"
			}
		}
		content+= char
		if (counter%16)==0{content+="\\n";counter=1}else{counter++}
	}
	self.nodes+=fmt.Sprintf("%s[label=\"{<h>File Block|%d}|{%s}\",fillcolor=\"#52ECF8\"];\n",
	node_id,
	self.format.Block_section.Bit_no_for(file_block.Index),
	content)
	return node_id
}

func (self *Aplication) Tree_repos(part_id string)(string,error){
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).Id != part_id{continue}
		mounted := &self.mounted_partitions[i]
		io_service:=mounted.io
		var super_block_start int32
		var err error
		var fit utiles.FitCriteria
		switch mounted.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
			super_block_start = logic.Part_start().Get()
			fit,err = utiles.Translate_fit(logic.Part_fit().Get())
			if err!=nil{return "",err}
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
			
			fit,err = utiles.Translate_fit(partition.Part_fit().Get())
			if err!=nil{return "",err}
		}
		sp_blck:=types.CreateSuperBlock(io_service,super_block_start)
		if sp_blck.S_firts_ino().Get() == -1 {return "",fmt.Errorf("can not generate tree report for corrupted ext3 filesystem")}

		format := formats.Recover_Format(io_service,super_block_start,fit)
		format.Init_bitmap_mapping()
		root_inode:=format.First_Inode()
		builder:=TreeRepoBuilder{
			nodes:         "",
			connections:   "",
			nodes_counter: 0,
			format:        &format,
		}
		ibuilder:=builder.New_Inode_node("",root_inode)
		Recursive_tree_repo(&root_inode,&ibuilder,&builder)
		builder.Collect_Inode_node(&ibuilder)
		return builder.Build(),nil
	}
	return "",fmt.Errorf("there's no partition with this id")
}

func normalize(name [12]string)string{
	str:=""
	for i := 0; i < 12; i++ {
		char:=name[i]
		ascii_val:=char[0]
		if ascii_val >= 32 && ascii_val <= 126{
			switch ascii_val{
			case 34:
				char = `\"`
			case 92:
				char = `\\`
			case 123:
				char = `\{`
			case 125:
				char = `\}`
			case 124:
				char = `\|`
			case 60:
				char = `\<`
			case 62:
				char = `\>`
			}
		}else{
			switch ascii_val{
			case 10:
				char = `\\n`
			case 11:
				char = `\\t`
			default: 
				char = "-"
			}
		}
		str+=char
	}
	return strings.TrimSpace(str)
}
func Recursive_tree_repo(root_inode *types.IndexNode,inode_bldr *InodeNodeBuilder,builder *TreeRepoBuilder){
	switch inode_bldr.inode_type{
	case (utiles.Directory):
		for n,ptr := range root_inode.I_block().Get(){
			inode_bldr.Push_row(n,ptr)
			if ptr == -1 {continue}
			if n >= 13 {
				ptr_block := types.CreatePointerBlock(builder.IoService(),ptr)
				ptr_block_bldr:=builder.New_PntrBlock_node(ptr_block,utiles.Directory)
				inode_bldr.Connect(n,ptr_block_bldr.id)
				recursive_tree_repo_in_pointer(&ptr_block,int32(n-13),&ptr_block_bldr,builder)
				builder.Collect_PntrBlock_node(&ptr_block_bldr)
			}else{	
				dir_block := types.CreateDirectoryBlock(builder.IoService(),ptr)
				dir_bldr:=builder.New_DirBLock_node(dir_block)
				inode_bldr.Connect(n,dir_bldr.id)

				for i,cont:=range dir_block.B_content().Get(){
					name:=normalize(cont.B_name)
					dir_bldr.Push_row(i,name,cont.B_inodo)
					if cont.B_inodo == -1{continue}
					if cont.B_name == formats.OWN_DIR_NAME {continue}
					if cont.B_name == formats.PARENT_DIR_NAME {continue}


					another_inode:=types.CreateIndexNode(builder.IoService(),cont.B_inodo)
					another_inode_bldr:=builder.New_Inode_node(name,another_inode)
					dir_bldr.Connect(i,another_inode_bldr.id)
					Recursive_tree_repo(&another_inode,&another_inode_bldr,builder) 
					builder.Collect_Inode_node(&another_inode_bldr)
				}
				builder.Collect_DirBlock_node(&dir_bldr)
			}
		}	
	case (utiles.File):
		for n,ptr := range root_inode.I_block().Get(){
			inode_bldr.Push_row(n,ptr)
			if ptr == -1 {continue}
			if n >= 13 {
				ptr_block := types.CreatePointerBlock(builder.IoService(),ptr)
				ptr_block_bldr:=builder.New_PntrBlock_node(ptr_block,utiles.File)
				inode_bldr.Connect(n,ptr_block_bldr.id)
				recursive_tree_repo_in_pointer(&ptr_block,int32(n-13),&ptr_block_bldr,builder)
				builder.Collect_PntrBlock_node(&ptr_block_bldr)
			}else{	
				file_block := types.CreateFileBlock(builder.IoService(),ptr)
				inode_bldr.Connect(n,builder.New_FileBlock_node(file_block))
			}
		}	
	}
}

func recursive_tree_repo_in_pointer(pointer_block *types.PointerBlock, level int32, ptr_blck_bldr *PntrBlockNodeBuilder, builder *TreeRepoBuilder){
	if level != 0 { //recursive thing
		for n,pointer_index := range pointer_block.Get().B_pointers{
			ptr_blck_bldr.Push_row(n,pointer_index)
			if pointer_index == -1 {continue}
			ptr_blck := types.CreatePointerBlock(builder.IoService(),pointer_index)
			new_ptr_blck_bldr:=builder.New_PntrBlock_node(ptr_blck,ptr_blck_bldr.inode_type)
			ptr_blck_bldr.Connect(n,new_ptr_blck_bldr.id)
			recursive_tree_repo_in_pointer(&ptr_blck,level-1,&new_ptr_blck_bldr,builder)
			builder.Collect_PntrBlock_node(&new_ptr_blck_bldr)
		}
		return
	}
	switch ptr_blck_bldr.inode_type{
	case (utiles.Directory):
		for n,pointer_index := range pointer_block.B_pointers().Get(){
			ptr_blck_bldr.Push_row(n,pointer_index)
			if pointer_index == -1 {continue}
			dir_block := types.CreateDirectoryBlock(builder.IoService(),pointer_index)
			dir_block_builder:=builder.New_DirBLock_node(dir_block)
			ptr_blck_bldr.Connect(n,dir_block_builder.id)
			
			for i,cont:=range dir_block.B_content().Get(){
				name:=normalize(cont.B_name)
				dir_block_builder.Push_row(i,name,cont.B_inodo)
				if cont.B_inodo == -1{continue}
				if cont.B_name == formats.OWN_DIR_NAME {continue}
				if cont.B_name == formats.PARENT_DIR_NAME {continue}

				another_inode:=types.CreateIndexNode(builder.IoService(),cont.B_inodo)
				another_inode_bldr:=builder.New_Inode_node(name,another_inode)
				dir_block_builder.Connect(i,another_inode_bldr.id)
				Recursive_tree_repo(&another_inode,&another_inode_bldr,builder)
				builder.Collect_Inode_node(&another_inode_bldr)
				
			}
			builder.Collect_DirBlock_node(&dir_block_builder)
		}
	case (utiles.File):
		for n,pointer_index := range pointer_block.Get().B_pointers{
			ptr_blck_bldr.Push_row(n,pointer_index)
			if pointer_index == -1 {continue}
			file_block := types.CreateFileBlock(builder.IoService(),pointer_index)
			ptr_blck_bldr.Connect(n,builder.New_FileBlock_node(file_block))
		}
	}
}











func (self *Aplication)Journaling(part_id string)(string,error){
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).Id != part_id{continue}
		mounted := &self.mounted_partitions[i]
		io_service:=mounted.io
		var super_block_start int32
		var err error
		var fit utiles.FitCriteria
		switch mounted.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
			super_block_start = logic.Part_start().Get()
			fit,err = utiles.Translate_fit(logic.Part_fit().Get())
			if err!=nil{return "",err}
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
			
			fit,err = utiles.Translate_fit(partition.Part_fit().Get())
			if err!=nil{return "",err}
		}
		format := formats.Recover_Format(io_service,super_block_start,fit)
		if !format.Has_journaling() {return "",fmt.Errorf("cant generate report for non ext3 formated partition")}
		rep := fmt.Sprintf(`digraph G {rankdir=LR;subgraph cluster_0 {
			label="Journaling";node[shape=none];
				node3[label=<<TABLE  CELLSPACING="0" BORDER="0" CELLBORDER = "1" >
				<TR><TD>Commands</TD><TD>Args</TD><TD>Date</TD></TR>
				%s
				</TABLE>>];}}
				`,format.Get_dot_journal_rep())
		
		return rep,nil
	}
	return "",fmt.Errorf("there's no partition with this id")
}









func (self *Aplication) Super_block_repo(part_id string)(string,error){
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).Id != part_id{continue}
		mounted := &self.mounted_partitions[i]
		io_service:=mounted.io
		var super_block_start int32
		switch mounted.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
			super_block_start = logic.Part_start().Get()
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
		}
		sb:=types.CreateSuperBlock(io_service,super_block_start)
		return fmt.Sprintf(`digraph G {subgraph cluster_0 {node1[style=filled,color=white,label=<%s>];}}`,sb.Dot_label()),nil
	}
	return "",fmt.Errorf("there's no partition with this id")
}





func (self *Aplication) Ls_report(part_id string,folders_trgt [][12]string) (string,error) {
	current_time := utiles.Current_Time()
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).Id != part_id{continue}
		mounted := &self.mounted_partitions[i]
		io_service:=mounted.io
		var super_block_start int32
		var err error
		var fit utiles.FitCriteria
		switch mounted.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
			super_block_start = logic.Part_start().Get()
			fit,err = utiles.Translate_fit(logic.Part_fit().Get())
			if err!=nil{return "",err}
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
			
			fit,err = utiles.Translate_fit(partition.Part_fit().Get())
			if err!=nil{return "",err}
		}
		format := formats.Recover_Format(io_service,super_block_start,fit)
		session,err:=parse_into_session_manager(&format)
		if err!=nil{return "",err}
		root:=format.First_Inode()
		err0,trgt_dir := format.Get_nested_dir(root,folders_trgt,false,1,1,current_time,false,false)
		if err0 != nil{return "",err0}
		date_str_conv:= func(date types.TimeHolder)string{
			return fmt.Sprintf("%d/%d/%d %d:%d",date.Day,date.Month,date.Year,date.Hour,date.Minute)
		}
		ugo_str_conv:= func(str [3]string)string{
			ugo:=utiles.UGOPermision_from_str(str)
			return ugo.Canonical_repr()
		}
		to_plain_name := func (str [12]string)string{
			final:=""
			for _, c := range str {final += c}
			return strings.TrimSpace(final)
		}
		var recursive_call func(types.Content)string
		recursive_call = func (content types.Content)string{
			inode:=types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
			name:= to_plain_name(content.B_name().Get())
			rows:=""
			permision_repr:=ugo_str_conv(inode.I_perm().Get())
			usr,err:=session.Get_User_by_id(int(inode.I_uid().Get()))
			if err != nil {panic("no user with that id")}
			owner_repr:=usr.name
			group_repr:=usr.group
			creation_date_repr:=date_str_conv(inode.I_ctime().Get())
			access_date_repr:=date_str_conv(inode.I_atime().Get())
			modified_date_repr:=date_str_conv(inode.I_mtime().Get())
			var type_repr string
			inode_type:=inode.I_type().Get()
			switch inode_type{
			case string(utiles.Directory):
				type_repr = "Directory"
			case string(utiles.File):
				type_repr = "File"
			}
			rows += fmt.Sprintf(`<TR><TD>%s</TD><TD>%s</TD><TD>%s</TD><TD>%d</TD><TD>%s</TD><TD>%s</TD><TD>%s</TD><TD>%s</TD><TD>%s</TD></TR>`,
			permision_repr,owner_repr,group_repr,inode.I_s().Get(),creation_date_repr,modified_date_repr,access_date_repr,type_repr,name)
			if inode_type ==string(utiles.Directory){
				all_childs:=format.Get_strict_shallow_tree_of_childs(inode)
				for _, child := range all_childs {
					row := recursive_call(child)
					rows += row
				}
			}
			return rows
		}
		all_childs:=format.Get_strict_shallow_tree_of_childs(trgt_dir)
		rows:=""
		for _, child := range all_childs {
			rows += recursive_call(child)
		}
		
		return fmt.Sprintf(`digraph G {subgraph cluster_0 {node2 [shape=none,label=<<TABLE ><TR><TD>PERMISION</TD><TD>OWNER</TD>
					  <TD>GROUP</TD><TD>SIZE</TD><TD>CREATION DATE</TD><TD>MOD DATE</TD>
					  <TD>ACCES DATE</TD><TD>TYPE</TD><TD>NAME</TD></TR>%s</TABLE>>];}}`,rows),nil
	}
	return "",fmt.Errorf("there's no partition with this id")
}


// func (self *Aplication) File_repo(part_id string,file_path string)(string,error){
// 	for i := 0; i < len(self.mounted_partitions); i++ {
// 		if (self.mounted_partitions[i]).id != part_id{continue}
// 		mounted := &self.mounted_partitions[i]
// 		io_service:=mounted.io
// 		var super_block_start int32
// 		var err error
// 		var fit utiles.FitCriteria
// 		switch mounted.part_type {
// 		case utiles.Logic:
// 			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
// 			super_block_start = logic.Part_start().Get()
// 			fit,err = utiles.Translate_fit(logic.Part_fit().Get())
// 			if err!=nil{return "",err}
			
// 		case utiles.Primary:
// 			partition := types.CreatePartition(io_service,mounted.index)
// 			super_block_start = partition.Part_start().Get()
			
// 			fit,err = utiles.Translate_fit(partition.Part_fit().Get())
// 			if err!=nil{return "",err}
// 		}
// 		format := formats.Recover_Format(io_service,super_block_start,fit)
		
// 		sb:=types.CreateSuperBlock(io_service,super_block_start)
// 		return fmt.Sprintf(`digraph G {subgraph cluster_0 {node1[style=filled,color=white,label=<%s>];}}`,sb.Dot_label()),nil
// 	}
// 	return "",fmt.Errorf("there's no partition with this id")
// }
