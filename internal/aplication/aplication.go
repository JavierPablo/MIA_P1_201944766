package aplication

import (
	"fmt"
	"project/internal/datamanagment"
	"project/internal/formats/ext2"
	"project/internal/types"
	"project/internal/utiles"
	"strconv"
	"strings"
)


type Aplication struct{
	partition_correlative int32
	logic_correlative int32
	active_partition *MountedPartition
	mounted_partitions []MountedPartition
}

type MountedPartition struct{
	id string
	part_type utiles.PartitionType
	index int32
	has_session bool
	session SessionManager
}
func (self *MountedPartition) Write_session_to_root_file(io_service *datamanagment.IOService)bool{
	var abs_index int32
	var fit utiles.FitCriteria
	switch self.part_type {
	case utiles.Logic:
		ebr := types.CreateExtendedBootRecord(io_service,self.index)
		abs_index = ebr.Part_start().Get()
		fit = utiles.Translate_fit(ebr.Part_fit().Get())
	case utiles.Primary:
		partition := types.CreatePartition(io_service,self.index)
		abs_index = partition.Part_start().Get()
		fit = utiles.Translate_fit(partition.Part_fit().Get())
	}
	format := ext2.Recover_FormatEXT2(io_service,abs_index,fit)
	format.Init_bitmap_mapping()
	r,root_folder := format.First_Inode()
	if r == -1 {return false}
	r,content :=format.Search_for_inode(root_folder,utiles.Into_ArrayChar12("user.txt"))
	if r == -1 {return false}
	file:=types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
	session:=self.session.To_file()
	format.Update_file(&file,0,session)
	return true
}

type SessionManager struct{
	active_user *User
	users []User
	groups []Group
	correlative_groups int
	correlative_users int
	insertion_order []Entity
	
}
func (self *SessionManager) Try_log_user(user string, pass string)bool{
	for n, u := range self.users {
		if u.name != user || u.password != pass {continue}
		self.active_user = &self.users[n]
		return true
	}
	return false
}
func (self *SessionManager) Remove_Group(group string)bool{
	for _,g := range self.groups {
		if g.name != group {continue}
		g.correlative_number = 0
		return true
	}
	return false
}
func (self *SessionManager) New_Group(group string)bool{
	for _,g := range self.groups {
		if g.name == group {return false}
	}
	self.correlative_groups++
	new_group :=Group{
		correlative_number: self.correlative_groups,
		name:               group,
	}
	self.groups = append(self.groups, new_group)
	self.insertion_order = append(self.insertion_order, &new_group)
	return true
}
func (self *SessionManager) Remove_User(name string)bool{
	for _,u := range self.users {
		if u.name != name {continue}
		u.correlative_number = 0
		return true
	}
	return false
}
func (self *SessionManager) New_User(name string, pass string, group string)bool{
	for _,u := range self.users {
		if u.name == name {return false}
	}
	grp_available := false
	for _,g := range self.groups {
		if g.name == group {
			if !g.available() {return false}
			grp_available = true
			break
		}
	}
	if !grp_available {return false}
	self.correlative_users++
	new_user :=User{
		correlative_number: self.correlative_users,
		name:               name,
		password:           pass,
		group:              group,
	}

	self.users = append(self.users, new_user)
	self.insertion_order = append(self.insertion_order, &new_user)
	return true
}
func (self *SessionManager) To_file()[]string{
	data := make([]string,0,30)
	for _,e := range self.insertion_order {
		data = append(data, strings.Split(e.to_string_line(),"")...)	
	}
	return data
}
func (self *SessionManager) Get_user_group(user *User)*Group{
	for _,g := range self.groups {
		if g.name == user.group{
			// if !g.available() {return nil}
			return &g
		}
	}
	return nil
}

func parse_into_session_manager(format *ext2.FormatEXT2)SessionManager{
	_,root_dir := format.First_Inode()
	indx,content := format.Search_for_inode(root_dir,utiles.Into_ArrayChar12("users.txt"))
	if indx == -1{panic("\"users.txt\" file not created wich is abnormal")}
	file:=types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
	if file.I_type().Get() != string(utiles.File){panic("\"users.txt\" is not file type")}
	
	lines := strings.Split(strings.Join(format.Read_file(&file), ""),"\n")
	const USER_INDICATOR string = "U"
	const GROUP_INDICATOR string = "H"
	user_list := make([]User,0,10)
	group_list := make([]Group,0,10)
	order := make([]Entity,0,10)
	user_corr := 1
	group_corr := 1
	for _,line := range lines{
		entity_comp := strings.Split(line, ",")
		correlative,err:= strconv.Atoi(entity_comp[0])
		if err != nil{panic("Correlative number was not a valid numeric string")}
		if entity_comp[1] == USER_INDICATOR{
			user := User{
				correlative_number: correlative,
				group:              entity_comp[2],
				name:               entity_comp[3],
				password:           entity_comp[4],
			}
			user_list = append(user_list, user)
			user_corr = correlative
			order = append(order, &user)
		}else if entity_comp[1] == GROUP_INDICATOR{
			group := Group{
				correlative_number: correlative,
				name:               entity_comp[2],
			}
			group_corr = correlative
			group_list = append(group_list, group)
			order = append(order, &group)
		}
	}
	file.I_atime().Set(utiles.Current_Time())
	return SessionManager{
		active_user:     nil,
		users:           user_list,
		groups:          group_list,
		insertion_order: order,
		correlative_groups: group_corr,
		correlative_users: user_corr,
	}

}

type User struct{
	correlative_number int
	name string
	password string
	group string
}
func (self *User) to_string_line()string{
	return fmt.Sprintf("%d,U,%s,%s,%s\n",self.correlative_number,self.name,self.group,self.password)
}
func (self *User) available()bool{
	return self.correlative_number != 0
}




type Group struct{
	correlative_number int
	name string
}
func (self *Group) available()bool{
	return self.correlative_number != 0
}
func (self *Group) to_string_line()string{
	return fmt.Sprintf("%d,G,%s\n",self.correlative_number,self.name)
}
type Entity interface{
	to_string_line()string;
	available()bool;
}



