package utiles

import (
	"fmt"
	"project/internal/types"
	"strings"
	"time"
)

type PartitionType string
const (
	Primary PartitionType = "P"
	Extendend PartitionType = "E"
	Logic PartitionType = "L"
)

type SizeUnit int32
const (
	Mb SizeUnit = 1024*1024
	Kb SizeUnit = 1024
	B SizeUnit = 1
)

type InodeType string
const (
	Directory InodeType = "0"
	File InodeType = "1"
)

type Format int32
const (
	Ext2 Format = 2
	Ext3 Format = 3
)

type FitCriteria string
const (
	First FitCriteria = "F"
	Best FitCriteria = "B"
	Worst FitCriteria = "W"
)

func Into_ArrayChar12(str string) [12]string{
	var stringArray [12]string
	chars:=strings.Split(str,"")

	copy(stringArray[:], chars)
	str_len :=len(chars)
	needed := 12 - str_len
	for i := 0;i<needed;i++{
		stringArray[i+str_len] = " "
	}
	return stringArray
}
func Into_ArrayChar4(str string) [4]string{
	var stringArray [4]string
	chars:=strings.Split(str,"")

	copy(stringArray[:], chars)
	str_len :=len(chars)
	needed := 4 - str_len
	for i := 0;i<needed;i++{
		stringArray[i+str_len] = " "
	}
	return stringArray
}
func Into_ArrayChar16(str string) [16]string{
	var stringArray [16]string
	chars:=strings.Split(str,"")

	copy(stringArray[:], chars)
	str_len :=len(chars)
	needed := 16 - str_len
	for i := 0;i<needed;i++{
		stringArray[i+str_len] = " "
	}
	return stringArray
}

func Current_Time()types.TimeHolder{
	time := time.Now()
	return types.TimeHolder{
		Hour:   int32(time.Hour()),
		Minute: int32(time.Minute()),
		Second: int32(time.Second()),
		Day:    int32(time.Day()),
		Month:  int32(time.Month()),
		Year:   int32(time.Year()),
	}
}
var NO_TIME types.TimeHolder = types.TimeHolder{
	Hour:   0,
	Minute: 0,
	Second: 0,
	Day:    0,
	Month:  0,
	Year:   0,
}
func Translate_size_unit(unit string)(SizeUnit,error){
	switch unit {
	case "K":
		return Kb,nil
	case "M":
		return Mb,nil
	case "B":
		return B,nil
	}
	return B,fmt.Errorf("size unit criteria not valid for %s",unit)
}
func Translate_fit(fit string)(FitCriteria,error){
	switch fit {
	case string(First):
		return First,nil
	case string(Worst):
		return Worst,nil
	case string(Best):
		return Best,nil
	}
	return "",fmt.Errorf("Fit criteria not valid for %s",fit)
}
func Translate_partition_type(type_ string)(PartitionType,error){
	switch type_ {
	case string(Primary):
		return Primary,nil
	case string(Extendend):
		return Extendend,nil
	case string(Logic):
		return Logic,nil
	}
	return "",fmt.Errorf("Partition type not valid for %s",type_)
}
func Translate_format_type(format string)(Format,error){
	switch format {
	case "2fs":
		return Ext2,nil
	case "3fs":
		return Ext3,nil
	}
	
	return 1,fmt.Errorf("format type not valid for %s",format)
}




type NameCriteria struct{
	Chars [12]Char
}
func (self *NameCriteria) Match(name [12]string)bool{
	for i := 0; i < 12; i++ {
		if !self.Chars[i].Matches(name[i]){return false}
	}
	return true
}
type Char struct{
	char string
	any_case bool
}
func New_Char(ch string)Char{
	return Char{
		char:     ch,
		any_case: false,
	}
}
var ANY_CHAR Char=Char{
	char:     "",
	any_case: true,
}
func (self *Char) Matches(trgt string) bool {
	if self.any_case {return true}
	return  self.char == trgt
}
