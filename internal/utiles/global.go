package utiles

import (
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
func Translate_fit(fit string)FitCriteria{
	switch fit {
	case string(First):
		return First
	case string(Worst):
		return Worst
	case string(Best):
		return Best
	}
	panic("Fit criteria not valid")
}




type NameCriteria struct{
	chars [12]Char
}
func (self *NameCriteria) Match(name [12]string)bool{
	for i := 0; i < 12; i++ {
		if !self.chars[i].Matches(name[i]){return false}
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
