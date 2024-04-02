package aplication

import (
	"fmt"
	"project/internal/formats"
	"project/internal/types"
	"project/internal/utiles"
)


func (self *Aplication) Log_in_user(part_id string, user string, password string)(*formats.JournalingManager,error){
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
			if err!=nil{return nil,err}
			
		case utiles.Primary:
			partition := types.CreatePartition(io_service,mounted.index)
			super_block_start = partition.Part_start().Get()
			
			fit,err = utiles.Translate_fit(partition.Part_fit().Get())
			if err!=nil{return nil,err}
		}
		if !mounted.has_session{
			format := formats.Recover_Format(io_service,super_block_start,fit)
			
			// format.Init_bitmap_mapping()
			session,err := parse_into_session_manager(&format)
			if err!=nil{return nil,err}
			mounted.session = session
			mounted.has_session = true
		}
		self.active_partition = &self.mounted_partitions[i]
		result := mounted.session.Try_log_user(user,password)
		return formats.Get_only_journaling(io_service,super_block_start),result
	}
	return nil,fmt.Errorf("there's no partition with this id")
}
func (self *Aplication) Log_out()(*formats.JournalingManager,error){
	if self.active_partition == nil {return nil,nil}
	if !self.active_partition.has_session {return nil,nil}
	if self.active_partition.session.active_user == nil {return nil,nil}
	var super_block_start int32
		switch self.active_partition.part_type {
		case utiles.Logic:
			logic := types.CreateExtendedBootRecord(self.active_partition.io,self.active_partition.index)
			super_block_start = logic.Part_start().Get()
		case utiles.Primary:
			partition := types.CreatePartition(self.active_partition.io,self.active_partition.index)
			super_block_start = partition.Part_start().Get()
		}
	self.active_partition.session.active_user = nil
	return formats.Get_only_journaling(self.active_partition.io,super_block_start),nil
}


func (self *Aplication) Make_group(group_name string)(*formats.JournalingManager,error){
	if self.active_partition == nil {return nil,fmt.Errorf("There's no active partition")}
	io_service := self.active_partition.io
	if !self.active_partition.has_session {return nil,fmt.Errorf("There's no active session. Login first")}
	if self.active_partition.session.active_user == nil {return nil,fmt.Errorf("There's no active user, log in")}
	if self.active_partition.session.active_user.name != "root" {return nil,fmt.Errorf("User loged in is not root user")}

	result := self.active_partition.session.New_Group(group_name)
	if result !=nil {return nil,result}
	return self.active_partition.Write_session_to_root_file(io_service)
}

func (self *Aplication) Remove_group(group_name string)(*formats.JournalingManager,error){
	if self.active_partition == nil {return nil,fmt.Errorf("There's no active partition")}
	io_service := self.active_partition.io
	if !self.active_partition.has_session {return nil,fmt.Errorf("There's no active session. Login first")}
	if self.active_partition.session.active_user == nil {return nil,fmt.Errorf("There's no active user, log in")}
	if self.active_partition.session.active_user.name != "root" {return nil,fmt.Errorf("User loged in is not root user")}

	result := self.active_partition.session.Remove_Group(group_name)

	if result !=nil {return nil,result}
	return self.active_partition.Write_session_to_root_file(io_service)
}



func (self *Aplication) Make_user(name string, pass string, group string)(*formats.JournalingManager,error){
	if self.active_partition == nil {return nil,fmt.Errorf("There's no active partition")}
	io_service := self.active_partition.io

	if !self.active_partition.has_session {return nil,fmt.Errorf("There's no active session. Login first")}
	if self.active_partition.session.active_user == nil {return nil,fmt.Errorf("There's no active user, log in")}
	if self.active_partition.session.active_user.name != "root" {return nil,fmt.Errorf("User loged in is not root user")}

	result := self.active_partition.session.New_User(name,pass,group)
	if result !=nil {return nil,result}
	return self.active_partition.Write_session_to_root_file(io_service)
}

func (self *Aplication) Remove_user(name string)(*formats.JournalingManager,error){
	if self.active_partition == nil {return nil,fmt.Errorf("There's no active partition")}
	io_service := self.active_partition.io

	if !self.active_partition.has_session {return nil,fmt.Errorf("There's no active session. Login first")}
	if self.active_partition.session.active_user == nil {return nil,fmt.Errorf("There's no active user, log in")}
	if self.active_partition.session.active_user.name != "root" {return nil,fmt.Errorf("User loged in is not root user")}

	result := self.active_partition.session.Remove_User(name)
	if result !=nil {return nil,result}
	return self.active_partition.Write_session_to_root_file(io_service)
}





func (self *Aplication) Chagne_User_Group(user string, group_name string)(*formats.JournalingManager,error){
	if self.active_partition == nil {return nil,fmt.Errorf("There's no active partition")}
	io_service := self.active_partition.io
	if !self.active_partition.has_session {return nil,fmt.Errorf("There's no active session. Login first")}
	if self.active_partition.session.active_user == nil {return nil,fmt.Errorf("There's no active user, log in")}
	if self.active_partition.session.active_user.name != "root" {return nil,fmt.Errorf("User loged in is not root user")}

	result := self.active_partition.session.Change_user_grp(user,group_name)
	if result !=nil {return nil,result}
	return self.active_partition.Write_session_to_root_file(io_service)
}

