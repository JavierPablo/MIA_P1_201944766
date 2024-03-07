package datamanagment

import (
	"os"
)

type Chunks struct{
	index int32
	buffer []byte
}
type IOService struct{
	// fragments []Chunks
	buffer	[]byte
	file_path string
}
func IOService_from(file_path string) IOService{
	
	b, err := os.ReadFile(file_path)
	if err != nil {
		panic(err)
	}
	return IOService{
		buffer: b,
		file_path :file_path,
	}
}
func (self *IOService) Read(amount int32, at int32) *[]byte {
	chunk := self.buffer[at:at+amount]
	return &chunk
}
func (self *IOService) Write(content []byte, at int32) {
	for n,b := range content{
		self.buffer[at+int32(n)] = b
	}
}

func (self *IOService) Flush(){	
	file,err := os.Create(self.file_path)
	if err != nil {panic(err)}
	file.Write(self.buffer)
	
}
