package types
import "project/internal/datamanagment"
import (
	byteslib "bytes"
	"encoding/binary"
)
type Integer struct{
	super_service *datamanagment.IOService
	index int32
	Size int32
}
func (self Integer)Set(obj int32) {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(obj))
	self.super_service.Write(bytes,self.index)
}
func (self Integer)Get() int32{
	var value int32
	chunk := self.super_service.Read(4,self.index)
	buffer := byteslib.NewReader(*chunk)
	err := binary.Read(buffer, binary.LittleEndian, &value)
	if err != nil {
		panic("Slice byte conversion to int32 has failled")
	}
	return value
}
type Boolean struct{
	super_service *datamanagment.IOService
	index int32
	Size int32
}
func (self Boolean)Set(a bool) {
	if a {
		self.super_service.Write([]byte{1},self.index)
	} else{
		self.super_service.Write([]byte{0},self.index)
	}
}
func (self Boolean)Get() bool{
	chunk := self.super_service.Read(1,self.index)
	return (*chunk)[0] != 0
}

type Character struct{
	super_service *datamanagment.IOService
	index int32
	Size int32
}
func (self Character)Set(a string) {
	self.super_service.Write([]byte{a[0]},self.index)
}
func (self Character)Get() string{
	chunk := self.super_service.Read(1,self.index)
	return string(*chunk)
}
