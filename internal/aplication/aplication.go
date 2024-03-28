package aplication

import (
	"fmt"
	"project/internal/datamanagment"
	"project/internal/formats/ext2"
	"project/internal/types"
	"project/internal/utiles"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// type LoadedService struct{
// 	letter string
// }
var Ok = color.New(color.FgGreen)
var Result = color.New(color.FgCyan)
var Err = color.New(color.FgRed)

type Aplication struct{
	// loaded_services []LoadedService
	partition_correlative int32
	logic_correlative int32
	active_partition *MountedPartition
	mounted_partitions []MountedPartition
}
func (self *Aplication) Print_mounted(){
	for i, mounted := range self.mounted_partitions {
		Result.Printf("%d | %s -> %s = %s\n",i,mounted.id,mounted.name, string(mounted.part_type))
	}
}
// func (self *Aplication) Put_service(letter string,io datamanagment.IOService){
// 	for i := 0; i < len(self.loaded_services); i++ {
// 		if self.loaded_services[i].letter == letter{
// 			return
// 		}
// 	}
// 	self.loaded_services = append(self.loaded_services, LoadedService{
// 		letter: letter,
// 		io:     io,
// 	})
// }
// func (self *Aplication) Get_service(letter string)(*datamanagment.IOService,error){
// 	for i := 0; i < len(self.loaded_services); i++ {
// 		if self.loaded_services[i].letter == letter{
// 			return &self.loaded_services[i].io,nil
// 		}
// 	}
// 	return nil,fmt.Errorf("There is no loaded service for %s",letter)
// }

type MountedPartition struct{
	name string
	io *datamanagment.IOService
	id string
	part_type utiles.PartitionType
	index int32
	has_session bool
	session SessionManager
}
func (self *MountedPartition) Write_session_to_root_file(io_service *datamanagment.IOService)(error){
	var err error
	var abs_index int32
	var fit utiles.FitCriteria
	switch self.part_type {
	case utiles.Logic:
		ebr := types.CreateExtendedBootRecord(io_service,self.index)
		abs_index = ebr.Part_start().Get()
		fit, err = utiles.Translate_fit(ebr.Part_fit().Get())
		if err != nil {return err}
	case utiles.Primary:
		partition := types.CreatePartition(io_service,self.index)
		abs_index = partition.Part_start().Get()
		fit,err = utiles.Translate_fit(partition.Part_fit().Get())
		if err != nil {return err}
	}
	format := ext2.Recover_FormatEXT2(io_service,abs_index,fit)
	format.Init_bitmap_mapping()
	root_folder := format.First_Inode()
	if root_folder.Index == -1 {return fmt.Errorf("There is no root directory for this partition")}
	indx,content :=format.Search_for_inode(root_folder,utiles.Into_ArrayChar12("users.txt"))
	if indx == -1 {return fmt.Errorf("There is no users.txt file for this partition")}
	file:=types.CreateIndexNode(content.Super_service,content.B_inodo().Get())
	session:=self.session.To_file()
	format.Update_file(&file,0,session)
	return nil
}

type SessionManager struct{
	active_user *User
	users []User
	groups []Group
	correlative_groups int
	correlative_users int
	insertion_order []Entity
	
}
func (self *SessionManager) Try_log_user(user string, pass string)error{
	for n, u := range self.users {
		if u.name == user {
			if u.password != pass {
				return fmt.Errorf("wrong password for user")
			}
		}else{continue}
		
		self.active_user = &self.users[n]
		return nil
	}
	return fmt.Errorf("There is no user with that name")
}
func (self *SessionManager) Remove_Group(group string)error{
	for i,g := range self.groups {
		if g.name != group {continue}
		(self.groups[i]).correlative_number = 0
		return nil
	}
	return fmt.Errorf("There's no group with that name")
}
func (self *SessionManager) New_Group(group string)error{
	for _,g := range self.groups {
		if g.name == group {
			return fmt.Errorf("There's already a group with that name")
		}
	}
	self.correlative_groups++
	new_group :=Group{
		correlative_number: self.correlative_groups,
		name:               group,
	}
	self.groups = append(self.groups, new_group)
	self.insertion_order = append(self.insertion_order, &self.groups[len(self.groups)-1])
	return nil
}
func (self *SessionManager) Get_User(name string)(*User,error){
	for i,u := range self.users {
		if u.name != name {continue}
		return &(self.users[i]),nil
	}
	return nil,fmt.Errorf("There's no user with that name")

}
func (self *SessionManager) Remove_User(name string)error{
	for i,u := range self.users {
		if u.name != name {continue}
		(self.users[i]).correlative_number = 0
		return nil
	}
	return fmt.Errorf("There's no user with that name")

}
func (self *SessionManager) Change_user_grp(name string, group string)error{
	var grp *Group = nil
	for i := 0; i < len(self.groups); i++ {
		if self.groups[i].name != group{continue}
		grp = &self.groups[i]
		break
	}
	if grp == nil {return fmt.Errorf("no such group with that name")}
	if !grp.available() {return fmt.Errorf("the group is not available")}

	for i := 0; i < len(self.users); i++ {
		if self.users[i].name != name{continue}
		self.users[i].group = group
		return nil
	}
	
	return fmt.Errorf("user doesnt exist")
}
func (self *SessionManager) New_User(name string, pass string, group string)error{
	for _,u := range self.users {
		if u.name == name {
			return fmt.Errorf("There's already a user with that name")
		}
	}
	grp_available := false
	for _,g := range self.groups {
		if g.name == group {
			if !g.available() {
				return fmt.Errorf("Group assigned has already been removed before")
			}
			grp_available = true
			break
		}
	}
	if !grp_available {
		return fmt.Errorf("There's no group with that name")
	}
	self.correlative_users++
	new_user :=User{
		correlative_number: self.correlative_users,
		name:               name,
		password:           pass,
		group:              group,
	}

	self.users = append(self.users, new_user)
	self.insertion_order = append(self.insertion_order, &self.users[len(self.users)-1])
	return nil
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
			if !g.available() {return nil}
			return &g
		}
	}
	return nil
}

func parse_into_session_manager(format *ext2.FormatEXT2)(SessionManager,error){
	root_dir := format.First_Inode()
	indx,content := format.Search_for_inode(root_dir,utiles.Into_ArrayChar12("users.txt"))
	if indx == -1{
		return SessionManager{},fmt.Errorf("\"users.txt\" file not created wich is abnormal")
	}
	file:=types.CreateIndexNode(content.Super_service,content.B_inodo().Get())

	if file.I_type().Get() != string(utiles.File){return SessionManager{}, fmt.Errorf("\"users.txt\" is not file type")}
	
	lines := strings.Split(strings.Join(format.Read_file(&file), ""),"\n")
	const USER_INDICATOR string = "U"
	const GROUP_INDICATOR string = "G"
	user_list := make([]User,0,10)
	group_list := make([]Group,0,10)
	order := make([]Entity,0,10)
	user_corr := 1
	group_corr := 1
	for _,line := range lines{
		if len(line) ==0 {continue}
		entity_comp := strings.Split(line, ",")
		correlative,err:= strconv.Atoi(entity_comp[0])
		
		if err != nil{return SessionManager{}, fmt.Errorf("Correlative number was not a valid numeric string")}
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
	},nil

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



