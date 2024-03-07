package ext2

import (
	// "fmt"
	"project/internal/datamanagment"
)


type Section struct{
	super_service *datamanagment.IOService	
	index int32
	length int32
	sub_size int32
}
func New_section(super_service *datamanagment.IOService,
	index int32,length int32,sub_size int32) Section{
return Section{
	super_service: super_service,
	index:         index,
	length:        length,
	sub_size:      sub_size,
}
	}
func (self *Section) Bit_no_for(index int32)int32{
	// fmt.Printf("DEBUG: init index for block section %d\n",self.index)
	// fmt.Printf("DEBUG: length of block %d\n",self.sub_size)
	return (index - self.index)/self.sub_size
}
func (self *Section) Index_for(bit_no int32)int32{
	return self.index + (bit_no * self.sub_size)
}
func (self *Section) Clear(){
	bytes_to_write := make([]byte,self.length)
	for i := 0; i < int(self.length); i++ {
		bytes_to_write[i] = 0
	}
	self.super_service.Write(bytes_to_write,self.index)
}