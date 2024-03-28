package datamanagment

import (
	"fmt"
	"os"
)

// type Chunks struct{
// 	index int32
// 	buffer []byte
// }
type IOService struct{
	// fragments []Chunks
	buffer	[]byte
	file_path string
	has_changes bool
}


func IOService_from_bytes(file_path string, bytes []byte) IOService{
	return IOService{
		buffer: bytes,
		file_path :file_path,
	}
}
func IOService_from(file_path string) (IOService,error){
	
	b, err := os.ReadFile(file_path)
	if err != nil {
		return IOService{},err
	}
	return IOService{
		buffer: b,
		file_path :file_path,
	},nil
}
func (self *IOService) Read(amount int32, at int32) *[]byte {
	// fmt.Println(self == nil)
	// fmt.Printf("sssssssssssssssssssssssssssssssssssssss amount %d at %d buffsize %d\n",amount,at,len(self.buffer))
	chunk := self.buffer[at:at+amount]
	// fmt.Println("------------------------------------")
	fmt.Print("")
	return &chunk
}
func (self *IOService) Write(content []byte, at int32) {
	self.has_changes = true
	for n,b := range content{
		self.buffer[at+int32(n)] = b
	}
}

func (self *IOService) Flush(){	
	file,err := os.Create(self.file_path)
	if err != nil {panic(err)}
	file.Write(self.buffer)
	
}
