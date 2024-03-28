package aplication

import (
	"fmt"
	"project/internal/formats/ext2"
	"project/internal/types"
	"project/internal/utiles"
)


func (self *Aplication) Log_in_user(part_id string, user string, password string)error{
	for i := 0; i < len(self.mounted_partitions); i++ {
		if (self.mounted_partitions[i]).id != part_id{continue}
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
			if err!=nil{return err}
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
			
			fit,err = utiles.Translate_fit(partition.Part_fit().Get())
			if err!=nil{return err}
		}
		if !mounted.has_session{
			format := ext2.Recover_FormatEXT2(io_service,super_block_start,fit)
			
			// format.Init_bitmap_mapping()
			session,err := parse_into_session_manager(&format)
			if err!=nil{return err}
			mounted.session = session
			mounted.has_session = true
		}
		self.active_partition = &self.mounted_partitions[i]
		result := mounted.session.Try_log_user(user,password)
		return result
	}
	return fmt.Errorf("There's no partition with this id")
}
func (self *Aplication) Log_out()error{
	if self.active_partition == nil {return nil}
	if !self.active_partition.has_session {return nil}
	if self.active_partition.session.active_user == nil {return nil}
	self.active_partition.session.active_user = nil
	return nil
}


func (self *Aplication) Make_group(group_name string)error{
	if self.active_partition == nil {return fmt.Errorf("There's no active partition")}
	io_service := self.active_partition.io
	if !self.active_partition.has_session {return fmt.Errorf("There's no active session. Login first")}
	if self.active_partition.session.active_user == nil {return fmt.Errorf("There's no active user, log in")}
	if self.active_partition.session.active_user.name != "root" {return fmt.Errorf("User loged in is not root user")}

	result := self.active_partition.session.New_Group(group_name)
	if result !=nil {return result}
	return self.active_partition.Write_session_to_root_file(io_service)
}

func (self *Aplication) Remove_group(group_name string)error{
	if self.active_partition == nil {return fmt.Errorf("There's no active partition")}
	io_service := self.active_partition.io
	if !self.active_partition.has_session {return fmt.Errorf("There's no active session. Login first")}
	if self.active_partition.session.active_user == nil {return fmt.Errorf("There's no active user, log in")}
	if self.active_partition.session.active_user.name != "root" {return fmt.Errorf("User loged in is not root user")}

	result := self.active_partition.session.Remove_Group(group_name)

	if result !=nil {return result}
	return self.active_partition.Write_session_to_root_file(io_service)
}



func (self *Aplication) Make_user(name string, pass string, group string)error{
	if self.active_partition == nil {return fmt.Errorf("There's no active partition")}
	io_service := self.active_partition.io

	if !self.active_partition.has_session {return fmt.Errorf("There's no active session. Login first")}
	if self.active_partition.session.active_user == nil {return fmt.Errorf("There's no active user, log in")}
	if self.active_partition.session.active_user.name != "root" {return fmt.Errorf("User loged in is not root user")}

	result := self.active_partition.session.New_User(name,pass,group)
	if result !=nil {return result}
	return self.active_partition.Write_session_to_root_file(io_service)
}

func (self *Aplication) Remove_user(name string)error{
	if self.active_partition == nil {return fmt.Errorf("There's no active partition")}
	io_service := self.active_partition.io

	if !self.active_partition.has_session {return fmt.Errorf("There's no active session. Login first")}
	if self.active_partition.session.active_user == nil {return fmt.Errorf("There's no active user, log in")}
	if self.active_partition.session.active_user.name != "root" {return fmt.Errorf("User loged in is not root user")}

	result := self.active_partition.session.Remove_User(name)
	if result !=nil {return result}
	return self.active_partition.Write_session_to_root_file(io_service)
}





func (self *Aplication) Chagne_User_Group(user string, group_name string)error{
	if self.active_partition == nil {return fmt.Errorf("There's no active partition")}
	io_service := self.active_partition.io
	if !self.active_partition.has_session {return fmt.Errorf("There's no active session. Login first")}
	if self.active_partition.session.active_user == nil {return fmt.Errorf("There's no active user, log in")}
	if self.active_partition.session.active_user.name != "root" {return fmt.Errorf("User loged in is not root user")}

	result := self.active_partition.session.Change_user_grp(user,group_name)
	if result !=nil {return result}
	return self.active_partition.Write_session_to_root_file(io_service)
}

