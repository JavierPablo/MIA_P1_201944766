package aplication

import (
	"project/internal/datamanagment"
	"project/internal/formats/ext2"
	"project/internal/types"
	"project/internal/utiles"
)


func (self *Aplication) Log_in_user(io_service *datamanagment.IOService,part_id string, user string, password string)bool{
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).id != part_id{continue}
		mounted := self.mounted_partitions[i]
		var super_block_start int32
		var fit utiles.FitCriteria

		switch mounted.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(io_service,mounted.index)
			super_block_start = logic.Part_start().Get()
			fit = utiles.Translate_fit(logic.Part_fit().Get())
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
			fit = utiles.Translate_fit(partition.Part_fit().Get())
		}
		if !mounted.has_session{
			format := ext2.Recover_FormatEXT2(io_service,super_block_start,fit)
			session := parse_into_session_manager(&format)
			mounted.session = session
		}
		result := mounted.session.Try_log_user(user,password)
		return result
	}
	return false
}
func (self *Aplication) Log_out(io_service *datamanagment.IOService)bool{
	if self.active_partition == nil {return true}
	if !self.active_partition.has_session {return true}
	if self.active_partition.session.active_user == nil {return true}
	self.active_partition.session.active_user = nil
	return true
}


func (self *Aplication) Make_group(io_service *datamanagment.IOService,group_name string)bool{
	if self.active_partition == nil {return false}
	if !self.active_partition.has_session {return false}
	if self.active_partition.session.active_user == nil {return false}
	if self.active_partition.session.active_user.name != "root" {return false}

	result := self.active_partition.session.New_Group(group_name)
	if !result {return false}
	return self.active_partition.Write_session_to_root_file(io_service)
}

func (self *Aplication) Remove_group(io_service *datamanagment.IOService,group_name string)bool{
	if self.active_partition == nil {return false}
	if !self.active_partition.has_session {return false}
	if self.active_partition.session.active_user == nil {return false}
	if self.active_partition.session.active_user.name != "root" {return false}

	result := self.active_partition.session.Remove_Group(group_name)
	if !result {return false}
	return self.active_partition.Write_session_to_root_file(io_service)
}



func (self *Aplication) Make_user(io_service *datamanagment.IOService,name string, pass string, group string)bool{
	if self.active_partition == nil {return false}
	if !self.active_partition.has_session {return false}
	if self.active_partition.session.active_user == nil {return false}
	if self.active_partition.session.active_user.name != "root" {return false}

	result := self.active_partition.session.New_User(name,pass,group)
	if !result {return false}
	return self.active_partition.Write_session_to_root_file(io_service)
}

func (self *Aplication) Remove_user(io_service *datamanagment.IOService,name string)bool{
	if self.active_partition == nil {return false}
	if !self.active_partition.has_session {return false}
	if self.active_partition.session.active_user == nil {return false}
	if self.active_partition.session.active_user.name != "root" {return false}

	result := self.active_partition.session.Remove_User(name)
	if !result {return false}
	return self.active_partition.Write_session_to_root_file(io_service)
}