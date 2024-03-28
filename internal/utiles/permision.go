package utiles

import (
	"fmt"
	"strconv"
)
var ALL_PERMITION Permision = New_Permision(1,1,1)
var NO_PERMITION Permision = New_Permision(0,0,0)
var UGO_PERMITION_664 UGOPermision = UGOPermision{
	User:  New_Permision(1,1,0),
	Group: New_Permision(1,1,0),
	Other: New_Permision(0,1,0),
}
var NO_UGO_PERMITION UGOPermision = UGOPermision{
	User:  NO_PERMITION,
	Group: NO_PERMITION,
	Other: NO_PERMITION,
}
type UGOPermision struct{
	User Permision
	Group Permision
	Other Permision
}
func (self *UGOPermision) To_arr_string()[3]string{
	return [3]string{self.User.To_string(),self.Group.To_string(),self.Other.To_string()}
}
func UGOPermision_from_str(str [3]string)UGOPermision{
	return UGOPermision{
		User:  Permision_from_str(str[0]),
		Group: Permision_from_str(str[1]),
		Other: Permision_from_str(str[2]),
	}
}
func New_Permision(read int, write int, execute int)Permision{
	return Permision{
		read:    read,
		write:   write,
		execute: execute,
	}
}

func Permision_from_str(str string)Permision{
	num,err := strconv.Atoi(str)
	if err!=nil{
		panic(fmt.Sprintf("Problems in converting permition: %s is not a number",str))
	}
	new_perm := Permision{}
	new_perm.execute = num%2
	num = num/2
	new_perm.read = num%2
	num = num/2
	new_perm.write = num%2
	return new_perm
}

type Permision struct{
	read int
	write int
	execute int
}
func (self *Permision) Can_read()bool{
	return self.read == 1
}
func (self *Permision) Can_write()bool{
	return self.write == 1
}
func (self *Permision) Can_exec()bool{
	return self.execute == 1
}
func (self *Permision) To_string()string{
	return strconv.Itoa(self.execute*1 + self.read*2 + self.write*4)
	
}




func Parse_to_ugo(perm string)(UGOPermision,error){
	if len(perm) != 3 {return UGOPermision{},fmt.Errorf("wrong permition format")}
	return UGOPermision{
		User:  Permision_from_str(string(perm[0])),
		Group: Permision_from_str(string(perm[1])),
		Other: Permision_from_str(string(perm[2])),
	},nil
	
}