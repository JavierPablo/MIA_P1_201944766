package aplication

import (
	"fmt"
	"os"
	"project/internal/datamanagment"
	"project/internal/formats"
	"project/internal/types"
	"project/internal/utiles"
	"strconv"
	"strings"
)
func (self *Aplication) Make_disk(size int32, unite utiles.SizeUnit, fit utiles.FitCriteria,DISK_CONTAINER_PATH string) (string,types.MasterBootRecord,error){
	var file *os.File
	ascii := 65
	path := DISK_CONTAINER_PATH+"/"+string(ascii)+".dsk"
	file, err :=os.Open(path)
	for err == nil{
		file.Close()
		ascii += 1
		if ascii == 91{return "",types.MasterBootRecord{}, fmt.Errorf("There are no abc names available for disks")}
		path = DISK_CONTAINER_PATH+"/"+string(ascii)+".dsk"
		file, err = os.Open(path)
	}
	file, err = os.Create(path)
	if err != nil {return "",types.MasterBootRecord{}, fmt.Errorf("Problems creating disk even thouhg the name is available")}
	defer file.Close()
	total_bytes := int(size*int32(unite))
	buffer := make([]byte,total_bytes)
	
	_, err = file.Write(buffer)
	if err != nil {return "",types.MasterBootRecord{}, fmt.Errorf("Problems writing 0s in disk")}
	// for i := 0; i < 10 + 1; i++ {
	// 	_, err = file.Write(buffer)
	// 	if err != nil {return "",types.MasterBootRecord{}, fmt.Errorf("Problems writing 0s in disk")}
	// }
	io_service,err:=datamanagment.IOService_from(path)	
	if err!=nil{return "",types.MasterBootRecord{},err}
	mbr := types.CreateMasterBootRecord(&io_service,0)
	mbr.Mbr_tamano().Set(int32(total_bytes))
	mbr.Mbr_fecha_creacion().Set(utiles.Current_Time())
	mbr.Dsk_fit().Set(string(fit))
	for _,partition := range mbr.Mbr_partitions().Spread(){
		partition.Part_start().Set(-1)
	}
	io_service.Flush()
	return path,mbr,nil
}

func (self *Aplication) Remove_disk(letter string,DISK_CONTAINER_PATH string) error{
	path := DISK_CONTAINER_PATH+"/"+letter+".dsk"
	err := os.Remove(path)
	return err

}











func (self *Aplication) potential_last_available_ebr_in_ext_part(extended_partition types.Partition)(types.ExtendedBootRecord){
	begin := extended_partition.Part_start().Get()
	
	ebr := types.CreateExtendedBootRecord(extended_partition.Super_service,begin)
	ebr_part_start := ebr.Part_start().Get()
	for ebr_part_start != -1{
		next_ebr_index := ebr.Part_next().Get()
		if next_ebr_index == -1 {break}
		ebr = types.CreateExtendedBootRecord(ebr.Super_service,next_ebr_index)
		ebr_part_start = ebr.Part_start().Get()
	}
	return ebr
}
func (self *Aplication) Get_extended_part_space_manager(extended_partition types.Partition)(datamanagment.SpaceManager,error){
	begin := extended_partition.Part_start().Get()
	part_spaces := make([]datamanagment.Space, 0,5)
	
	ebr := types.CreateExtendedBootRecord(extended_partition.Super_service,begin)
	part_spaces = append(part_spaces, datamanagment.New_Space(0,ebr.Size))
	
	next_ebr_index := ebr.Part_next().Get()
	ebr_part_start := ebr.Part_start().Get()
	if ebr_part_start != -1{
		part_spaces = append(part_spaces, datamanagment.New_Space(ebr_part_start-begin,ebr.Part_s().Get()))
	}
	for next_ebr_index != -1{
		ebr = types.CreateExtendedBootRecord(ebr.Super_service,next_ebr_index)
		part_spaces = append(part_spaces, datamanagment.New_Space(next_ebr_index-begin,ebr.Size))
		next_ebr_index = ebr.Part_next().Get()
		ebr_part_start = ebr.Part_start().Get()
		if ebr_part_start != -1{
			part_spaces = append(part_spaces, datamanagment.New_Space(ebr_part_start-begin,ebr.Part_s().Get()))
		}
	}
	space_manager,err := datamanagment.SpaceManager_from_occuped_spaces(part_spaces,extended_partition.Part_s().Get())
	if err!=nil{return datamanagment.SpaceManager{},err}
	
	return space_manager,nil
}
func (self *Aplication) Get_extended_partition(mbr types.MasterBootRecord)(error,types.Partition){
	var extended_partition types.Partition
	found := false
	for _,partition := range mbr.Mbr_partitions().Spread(){
		if partition.Part_start().Get() == -1{continue}
		if partition.Part_type().Get() == string(utiles.Extendend){
			extended_partition = partition
			found= true
			break
		}
	}
	if !found{
		return fmt.Errorf("Extendes partition has not been created yet"),types.Partition{}
	}
	return nil,extended_partition
}
func (self *Aplication) Get_disk_partitions_space_manager(mbr types.MasterBootRecord)(datamanagment.SpaceManager,error){
	partitins_spaces := make([]datamanagment.Space, 0,4)
	for _,partition := range mbr.Mbr_partitions().Spread(){
		begin := partition.Part_start().Get()
		if begin == -1{continue}
		begin -= mbr.Size
		size := partition.Part_s().Get()
		space := datamanagment.New_Space(begin,size)
		partitins_spaces = append(partitins_spaces,space)
	}
	space_manager,err := datamanagment.SpaceManager_from_occuped_spaces(partitins_spaces,mbr.Mbr_tamano().Get()-mbr.Size)
	if err !=nil {return datamanagment.SpaceManager{},err}
	return space_manager,nil
}
func (self *Aplication) Find_p_or_e_partition_by_name(mbr types.MasterBootRecord,p_name [16]string)(bool,types.Partition){
	for _,partition := range mbr.Mbr_partitions().Spread(){
		begin := partition.Part_start().Get()
		if begin == -1{continue}
		if partition.Part_name().Get() == p_name{
			return true,partition
		}
	}
	return false,types.Partition{}
}
func (self *Aplication) Find_logical_partition_by_name(extended_partition types.Partition,p_name [16]string)(error,types.ExtendedBootRecord){
	begin := extended_partition.Part_start().Get()
	ebr := types.CreateExtendedBootRecord(extended_partition.Super_service,begin)
	next_ebr_index := ebr.Part_next().Get()
	if ebr.Part_start().Get() != -1 && ebr.Part_name().Get() == p_name{
		return nil,ebr
	}
	
	for next_ebr_index != -1{
		ebr = types.CreateExtendedBootRecord(ebr.Super_service,next_ebr_index)
		next_ebr_index = ebr.Part_next().Get()
		if ebr.Part_start().Get() != -1 && ebr.Part_name().Get() == p_name{
			return nil,ebr
		}
	}
	return fmt.Errorf("No logical partition found"),types.ExtendedBootRecord{}
}
func (self *Aplication) modify_logical_partition(mbr types.MasterBootRecord,add int32,p_name [16]string)error{
	
	found_err,extended_partition:= self.Get_extended_partition(mbr)
	begin := extended_partition.Part_start().Get()

	if found_err != nil{
		return found_err
	}
	space_manager,err := self.Get_extended_part_space_manager(extended_partition)
	if err !=nil {return err}

	
	err,ebr:=self.Find_logical_partition_by_name(extended_partition,p_name)
	if err!=nil{
		return err
	}
	if add < 0 {
		if -add >ebr.Part_s().Get(){return fmt.Errorf("requested size for substraction is greather than actual size")}

		result := space_manager.Free_space(-add,ebr.Part_start().Get()+ebr.Part_s().Get()-begin+add)
		if result != nil{
			return result
		}
	}else{
		result := space_manager.Ocupe_raw_space(add,ebr.Part_start().Get()+ebr.Part_s().Get()-begin)
		if result != nil{
			new_start:=ebr.Part_start().Get() - add
			result := space_manager.Ocupe_raw_space(add,new_start - mbr.Size)
			if result != nil {return result}
			ebr.Part_start().Set(new_start)
			return result
		}
	}
	ebr.Part_s().Set(ebr.Part_s().Get()+add)
	return nil
}
func (self *Aplication) Modify_partition_size_in_disk(size int32, io_service *datamanagment.IOService, 
	p_name string, unite utiles.SizeUnit) error{
	partition_name := utiles.Into_ArrayChar16(p_name)
	add := size*int32(unite)
	mbr := types.CreateMasterBootRecord(io_service,0)
	partition_found,partition_trgt := self.Find_p_or_e_partition_by_name(mbr,partition_name)
	if !partition_found {
		// fmt.Println("------------------------------modifyng logical partition-------------------------")
		return self.modify_logical_partition(mbr,add,partition_name)
	}
	space_manager,err := self.Get_disk_partitions_space_manager(mbr)
	if err !=nil {return err}
	// fmt.Println("-----------------------modifyng primary partition--------------------------------")

	// space_manager.Log_chunks_state()
	if add < 0 {
		// fmt.Println("--------------------removing sapce-----------------------------------")

		if -add >partition_trgt.Part_s().Get(){return fmt.Errorf("requested size for substraction is greather than actual size")}
		result := space_manager.Free_space(-add,partition_trgt.Part_start().Get()+partition_trgt.Part_s().Get()-mbr.Size+add)
		if result != nil{
			return result
		}
	}else{
		// fmt.Println("--------------------adding space-----------------------------------")
		// fmt.Println(add)
		result := space_manager.Ocupe_raw_space(add,partition_trgt.Part_start().Get()+partition_trgt.Part_s().Get()-mbr.Size)
		if result != nil{
			// try in begin then
			new_start:=partition_trgt.Part_start().Get() - add
			result := space_manager.Ocupe_raw_space(add,new_start - mbr.Size)
			if result != nil {return result}
			partition_trgt.Part_start().Set(new_start)
		}
	}
	partition_trgt.Part_s().Set(partition_trgt.Part_s().Get()+add)
	return nil
}

func (self *Aplication) Remove_partition_disk(io_service *datamanagment.IOService, p_name string) error{
	update_index_for_remain_ebr:=func (old_index int32,new_index int32){
		for i := 0; i < len(self.mounted_partitions); i++ {
			mntd:=&self.mounted_partitions[i]
			if mntd.index == old_index{
				mntd.index = new_index
			}
		}
		panic("Unexpected behavior")

	}
	unmount_partition:=func (index int32){
		for i := 0; i < len(self.mounted_partitions); i++ {
			mntd:=&self.mounted_partitions[i]
			if mntd.index == index{
				if self.active_partition == mntd{
					self.active_partition = nil
				}
				self.mounted_partitions = append(self.mounted_partitions[:i],self.mounted_partitions[i+1:]... )
				return
			}
		}
		panic("Unexpected behavior")
	}
	mbr := types.CreateMasterBootRecord(io_service,0)
	partition_name := utiles.Into_ArrayChar16(p_name)
	partition_found,partition_trgt := self.Find_p_or_e_partition_by_name(mbr,partition_name)
	if !partition_found {
		// fmt.Println("------------------------------modifyng logical partition-------------------------")
		found_err,extended_partition:= self.Get_extended_partition(mbr)
		if found_err != nil{return found_err}
		begin := extended_partition.Part_start().Get()
		parent_ebr := types.CreateExtendedBootRecord(extended_partition.Super_service,begin)
		next_ebr_index := parent_ebr.Part_next().Get()

		if parent_ebr.Part_start().Get() != -1 && parent_ebr.Part_name().Get() == partition_name{
			if parent_ebr.Part_mount().Get() == "Y"{
				parent_ebr.Part_mount().Set("N")
				unmount_partition(parent_ebr.Index)
			}
			if next_ebr_index != -1{
				child_ebr := types.CreateExtendedBootRecord(parent_ebr.Super_service,next_ebr_index)
				parent_ebr.Set(child_ebr.Get())
				child_ebr.Part_s().Set(-1)
				child_ebr.Part_start().Set(-1)
				if child_ebr.Part_mount().Get() == "Y"{
					child_ebr.Part_mount().Set("N")
					update_index_for_remain_ebr(child_ebr.Index,parent_ebr.Index)
				}
			}else{
				parent_ebr.Part_s().Set(-1)
				parent_ebr.Part_start().Set(-1)
			}
			return nil
		}
		
		for next_ebr_index != -1{
			child_ebr := types.CreateExtendedBootRecord(parent_ebr.Super_service,next_ebr_index)
			next_ebr_index = child_ebr.Part_next().Get()
			if child_ebr.Part_start().Get() != -1 && child_ebr.Part_name().Get() == partition_name{
				child_ebr.Part_start().Set(-1)
				child_ebr.Part_s().Set(-1)
				child_ebr.Part_next().Set(-1)

				if child_ebr.Part_mount().Get() == "Y"{
					child_ebr.Part_mount().Set("N")
					unmount_partition(child_ebr.Index)
				}
				
				parent_ebr.Part_next().Set(next_ebr_index)
				return nil
			}
			parent_ebr = child_ebr
		}
		return fmt.Errorf("no primary/extended partition nor logical partition found for removal")
	}
	
	partition_trgt.Part_s().Set(-1)
	partition_trgt.Part_start().Set(-1)
	if partition_trgt.Part_status().Get() == "Y"{
		partition_trgt.Part_status().Set("N")
		unmount_partition(partition_trgt.Index)
	}
	return nil
}
func (self *Aplication) Partition_disk(size int32, io_service *datamanagment.IOService, 
	p_name string, unite utiles.SizeUnit, partition_type utiles.PartitionType,
	fit utiles.FitCriteria) (int32,error){

	mbr := types.CreateMasterBootRecord(io_service,0)
	partition_size := size*int32(unite)
	if partition_type == utiles.Primary{
		return self.create_EP_partitions(mbr,utiles.Primary,partition_size,fit,p_name)
	}else if partition_type == utiles.Extendend{
		for _,partition := range mbr.Mbr_partitions().Spread(){
			if partition.Part_start().Get() == -1{continue}
			if partition.Part_type().Get() == string(utiles.Extendend){
				return -1,fmt.Errorf("Extendend partition already exists")
			}
		}
		return self.create_EP_partitions(mbr,utiles.Extendend,partition_size,fit,p_name)
	}else if partition_type == utiles.Logic{
		found_err,extended_partition:= self.Get_extended_partition(mbr)
		if found_err != nil{
			return -1,found_err
		}
		// if extended_partition.Part_s().Get() < partition_size {return fmt.Errorf("can not make partition ")}
		begin :=extended_partition.Part_start().Get()
		space_manager,err := self.Get_extended_part_space_manager(extended_partition)
		if err !=nil {return -1,err}

		potential_ebr := self.potential_last_available_ebr_in_ext_part(extended_partition)
	
		ebr_abs_index := potential_ebr.Index
		if potential_ebr.Part_start().Get() != -1{
			ebr_space_index,err := Try_fit(&space_manager,potential_ebr.Size,extended_partition.Part_fit().Get()+"F")
			if err!=nil{return -1,err}
			if ebr_space_index == -1{return -1,fmt.Errorf("There is no enough space for the ebr")}
	
			ebr_abs_index = space_manager.Ocupe_space_unchecked(int(ebr_space_index),potential_ebr.Size) + begin
			potential_ebr.Part_next().Set(ebr_abs_index)
		}

		part_space_index,err := Try_fit(&space_manager,partition_size,extended_partition.Part_fit().Get() + "F")
		if err!=nil{return -1,err}
		if part_space_index == -1{return -1,fmt.Errorf("There is no enough space for the logic partition itself")}

		part_abs_index := space_manager.Ocupe_space_unchecked(int(part_space_index),partition_size) + begin
		new_ebr := types.CreateExtendedBootRecord(potential_ebr.Super_service,ebr_abs_index)
		new_ebr.Set(types.ExtendedBootRecordHolder{
			Part_mount: "N",
			Part_fit:   string(fit),
			Part_start: part_abs_index,
			Part_s:     partition_size,
			Part_next:  -1,
			Part_name:  utiles.Into_ArrayChar16(p_name),
		})

		return ebr_abs_index,nil
	}
	return -1,fmt.Errorf("Unknown partition type")
}
func (self *Aplication) create_EP_partitions(mbr types.MasterBootRecord,p_type utiles.PartitionType,partition_size int32,fit utiles.FitCriteria, p_name string)(int32,error){
	partitins_spaces := make([]datamanagment.Space, 0,4)
	var partition_available types.Partition
	partition_found := false

	for _,partition := range mbr.Mbr_partitions().Spread(){
		begin := partition.Part_start().Get()
		if begin == -1{
			if !partition_found{
				partition_found = true
				partition_available = partition
			}
			continue
		}
		begin = begin - mbr.Size
		size := partition.Part_s().Get()
		space := datamanagment.New_Space(begin,size)
		partitins_spaces = append(partitins_spaces,space)
	
	}
	

	if !partition_found {
		return -1,fmt.Errorf("There are no partition spaces available: the 4 partition spaces are already in use")
	}
	
	space_manager,err := datamanagment.SpaceManager_from_occuped_spaces(partitins_spaces,mbr.Mbr_tamano().Get()-mbr.Size)
	if err!=nil{return -1,err}
	
	space_index,err := Try_fit(&space_manager,partition_size,string(fit))
	if err != nil{return -1,err} 
	// fmt.Printf("partition_size = %d at index = %d for fit %s \n ",partition_size,space_index,fit)
	if space_index == -1{return -1,fmt.Errorf("There is no enough space for that operation")}

	abs_indx := space_manager.Ocupe_space_unchecked(int(space_index),partition_size) + mbr.Size
	
	partition_available.Set(types.PartitionHolder{
		Part_status:      "N",
		Part_type:        string(p_type),
		Part_fit:         string(fit),
		Part_start:       abs_indx,
		Part_s:           partition_size,
		Part_name:        utiles.Into_ArrayChar16(p_name),
		Part_correlative: int32(-1),
		Part_id:          utiles.Into_ArrayChar4("    "),
	})
	if p_type == utiles.Extendend{
		first_ebr := types.CreateExtendedBootRecord(partition_available.Super_service,abs_indx)
		first_ebr.Part_start().Set(-1)
		first_ebr.Part_next().Set(-1)
	}
	return abs_indx,nil
}
func Try_fit(self *datamanagment.SpaceManager,for_length int32,of_type string) (int32,error){
	switch strings.ToUpper(of_type){
		case string(utiles.Best):
			return self.Best_fit(for_length),nil
		case string(utiles.Worst):
			return self.Worst_fit(for_length),nil
		case string(utiles.First):
			return self.First_fit(for_length),nil
	}
	return -1,fmt.Errorf("Not matching any branch")
}








func (self *Aplication) Mount_partition(io_service *datamanagment.IOService, p_name string, letter string) error{
	mbr := types.CreateMasterBootRecord(io_service,0)
	result,partition := self.Find_p_or_e_partition_by_name(mbr,utiles.Into_ArrayChar16(p_name))
	if !result{
		err,extp := self.Get_extended_partition(mbr)
		if err != nil{ return fmt.Errorf("partition with that name doesnt exist") }
		err,ebr := self.Find_logical_partition_by_name(extp,utiles.Into_ArrayChar16(p_name))
		if err != nil{ return err }
		self.logic_correlative++
		ebr.Part_mount().Set("Y")
		id := letter + strconv.Itoa(int(self.logic_correlative)) + "66" + "L"
		extp.Part_status().Set("Y")

		temp:=ebr.Part_name().Get()
		name:= strings.Join(temp[:],"")
		
		self.mounted_partitions = append(self.mounted_partitions, MountedPartition{
			Name:	name,
			Id:        id,
			io:          io_service,
			part_type: utiles.Logic,
			index:     ebr.Index,
			has_session: false,
			session:   SessionManager{},
		})
		super_block := types.CreateSuperBlock(ebr.Super_service,ebr.Part_start().Get())
		format := super_block.S_filesystem_type().Get()
		if format != int32(utiles.Ext2) && format != int32(utiles.Ext3){return nil}
		super_block.S_mtime().Set(utiles.Current_Time())
		super_block.S_mnt_count().Set(super_block.S_mnt_count().Get()+1)
		return nil
	}
	if partition.Part_type().Get() == string(utiles.Extendend){return fmt.Errorf("can not mout extended partition")}
	self.partition_correlative++
	id := letter + strconv.Itoa(int(self.partition_correlative)) + "66"
	partition.Part_status().Set("Y")
	partition.Part_correlative().Set(self.partition_correlative)
	
	partition.Part_id().Set(utiles.Into_ArrayChar4(id))
	temp:=partition.Part_name().Get()
	name:= strings.Join(temp[:],"")
	
	self.mounted_partitions = append(self.mounted_partitions, MountedPartition{
		Name:	name,
		io:          io_service,
		Id:          id,
		part_type:   utiles.Primary,
		index:       partition.Index,
		has_session: false,
		session:     SessionManager{},
	})
	super_block := types.CreateSuperBlock(partition.Super_service,partition.Part_start().Get())
	format := super_block.S_filesystem_type().Get()
	if format == int32(utiles.Ext2) || format == int32(utiles.Ext3){
		super_block.S_umtime().Set(utiles.Current_Time())
		super_block.S_mtime().Set(utiles.Current_Time())
		super_block.S_mnt_count().Set(super_block.S_mnt_count().Get()+1)
	}
	// self.Put_service(letter,*io_service)
	return nil
}

func (self *Aplication) Unmount_partition(id string) error{
	// io_service := self.Get_service(id)
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).Id != id{continue}
		mounted := self.mounted_partitions[i]
		io_service := mounted.io

		var part_start int32
		switch mounted.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
			logic.Part_mount().Set("N")
			part_start = logic.Part_start().Get()
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			partition.Part_status().Set("N")
			partition.Part_correlative().Set(-1)
			partition.Part_id().Set(utiles.Into_ArrayChar4("     "))
			part_start = partition.Part_s().Get()
			
		}
		self.mounted_partitions = append(self.mounted_partitions[:i],self.mounted_partitions[i+1:]... )
		super_block := types.CreateSuperBlock(io_service,part_start)
		format := super_block.S_filesystem_type().Get()
		if format == int32(utiles.Ext2) || format == int32(utiles.Ext3){
			super_block.S_umtime().Set(utiles.Current_Time())
		}
		return nil
	

	}
	return fmt.Errorf("mounted partition with that id not found")
}








func (self *Aplication) Format_mounted_partition(id string, p_type bool, format utiles.Format) error{
	for i := 0; i < len(self.mounted_partitions); i++ {
		// fmt.Printf("%s == %s",(self.mounted_partitions[i]).id,id)
		if (self.mounted_partitions[i]).Id == id{
			mounted := &self.mounted_partitions[i]
			io_service := mounted.io
			switch mounted.part_type {
			case utiles.Logic:
				logic := types.CreateExtendedBootRecord(io_service,mounted.index)
				super_block_start := logic.Part_start().Get()
				fit,err:=utiles.Translate_fit(logic.Part_fit().Get())
				if err!=nil{return err}
				switch format {
				case utiles.Ext2:
					formats.Format_new_fresh_FormatEXT2(io_service,fit,super_block_start,logic.Part_s().Get())	
				case utiles.Ext3:
					formats.Format_new_fresh_FormatEXT3(io_service,fit,super_block_start,logic.Part_s().Get())	
				}

			case utiles.Primary:
				partition := types.CreatePartition(io_service,mounted.index)
				super_block_start := partition.Part_start().Get()
				
				fit,err:=utiles.Translate_fit(partition.Part_fit().Get())
				if err!=nil{return err}
				switch format {
				case utiles.Ext2:
					formats.Format_new_fresh_FormatEXT2(io_service,fit,super_block_start,partition.Part_s().Get())				
				case utiles.Ext3:
					formats.Format_new_fresh_FormatEXT3(io_service,fit,super_block_start,partition.Part_s().Get())				
				}
			}
			return nil
		}

	}
	return fmt.Errorf("There is no mounted partition with that name")
}











