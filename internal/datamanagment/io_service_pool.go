package datamanagment

import "fmt"


type Disk struct{
	Letter string
	Io_service IOService
}
type IOServicePool struct{
	main_path string
	Pool []Disk
}
func New_IOServicePool(path string)IOServicePool{
	pool := make([]Disk,0,3)
	return IOServicePool{
		main_path: path,
		Pool:      pool,
	}
}
func (self *IOServicePool) Get_service_with(letter string )(*IOService,error){
	for i := 0; i < len(self.Pool); i++ {
		if self.Pool[i].Letter == letter{
			return &self.Pool[i].Io_service,nil
		}
	}
	ioservice,err:=IOService_from(self.main_path+"/"+letter+".dsk")
	if err != nil{return nil,fmt.Errorf("can not open disk with %s letter because it doesnt exist",letter)}
	self.Pool = append(self.Pool, Disk{
		Letter:     letter,
		Io_service: ioservice,
	})
	return &self.Pool[len(self.Pool)-1].Io_service,nil
}
func (self *IOServicePool) Flush_changes()(){
	for i := 0; i < len(self.Pool); i++ {
		if self.Pool[i].Io_service.has_changes{
			self.Pool[i].Io_service.Flush()
		}
	}
}