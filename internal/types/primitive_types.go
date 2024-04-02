package types

import (
	"project/internal/datamanagment"
	"strconv"

	byteslib "bytes"
	"encoding/binary"
)
type Integer struct{
	Super_service *datamanagment.IOService
	Index int32
	Size int32
}
func (self Integer)Set(obj int32) {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(obj))
	self.Super_service.Write(bytes,self.Index)
}
func (self Integer)Get() int32{
	var value int32
	chunk := self.Super_service.Read(4,self.Index)
	buffer := byteslib.NewReader(*chunk)
	err := binary.Read(buffer, binary.LittleEndian, &value)
	if err != nil {
		panic("Slice byte conversion to int32 has failled")
	}
	return value
}
func(self Integer)Dot_label()string{
	return strconv.Itoa(int(self.Get()))
}





type Boolean struct{
	Super_service *datamanagment.IOService
	Index int32
	Size int32
}
func (self Boolean)Set(a bool) {
	if a {
		self.Super_service.Write([]byte{1},self.Index)
	} else{
		self.Super_service.Write([]byte{0},self.Index)
	}
}
func (self Boolean)Get() bool{
	chunk := self.Super_service.Read(1,self.Index)
	return (*chunk)[0] != 0
}
func (self Boolean)Dot_label() string{
	if self.Get() {
		return "true"
	}
	return "false"
}








type Character struct{
	Super_service *datamanagment.IOService
	Index int32
	Size int32
}
func (self Character)Set(a string) {
	self.Super_service.Write([]byte{a[0]},self.Index)
}
func (self Character)Get() string{
	chunk := self.Super_service.Read(1,self.Index)
	// dsf :=(*chunk)[0]
	// if (*chunk)[0] < byte(32) {return "-"}
	// fmt.Println((*chunk)[0])
	return string(*chunk)
}
func (self Character)Dot_label() string{
	chunk := self.Super_service.Read(1,self.Index)
	// dsf :=(*chunk)[0]
	if int8((*chunk)[0]) < int8(32) {return "-"}
	return string(*chunk)
}
