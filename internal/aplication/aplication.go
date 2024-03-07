package aplication

import (
	"fmt"
	"project/internal/formats/ext2"
	// "project/internal/types"
	"project/internal/utiles"
	"strconv"
	"strings"
)


type Aplication struct{
	active_partition *MountedPartition
	mounted_partitions []MountedPartition
}


type MountedPartition struct{
	id string
	format ext2.FormatEXT2
	session SessionManager
}


type SessionManager struct{
	active_user *User
	users []User
	groups []Group
	insertion_order []Entity
	
}
func parse_into_session_manager(format *ext2.FormatEXT2)SessionManager{
	_,root_dir := format.First_Inode()
	indx,file := format.Search_for_inode(root_dir,utiles.Into_ArrayChar12("users.txt"))
	if indx == -1{panic("\"users.txt\" file not created wich is abnormal")}
	if file.I_type().Get() != string(utiles.File){panic("\"users.txt\" is not file type")}
	
	lines := strings.Split(strings.Join(format.Read_file(&file), ""),"\n")
	const USER_INDICATOR string = "U"
	const GROUP_INDICATOR string = "H"
	user_list := make([]User,0,10)
	group_list := make([]Group,0,10)
	order := make([]Entity,0,10)
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
			order = append(order, &user)
		}else if entity_comp[1] == GROUP_INDICATOR{
			group := Group{
				correlative_number: correlative,
				name:               entity_comp[2],
			}
			group_list = append(group_list, group)
			order = append(order, &group)
		}
	}
	return SessionManager{
		active_user:     nil,
		users:           user_list,
		groups:          group_list,
		insertion_order: order,
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
type Group struct{
	correlative_number int
	name string
}
func (self *Group) to_string_line()string{
	return fmt.Sprintf("%d,G,%s\n",self.correlative_number,self.name)
}
type Entity interface{
	to_string_line()string;
}
