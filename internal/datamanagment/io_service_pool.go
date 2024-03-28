package datamanagment


type Disk struct{
	letter string
	io_service IOService
}
type IOServicePool struct{
	main_path string
	pool []Disk
}
func New_IOServicePool(path string)IOServicePool{
	pool := make([]Disk,0,3)
	return IOServicePool{
		main_path: path,
		pool:      pool,
	}
}
func (self *IOServicePool) Get_service_with(letter string )(*IOService,error){
	for i := 0; i < len(self.pool); i++ {
		if self.pool[i].letter == letter{
			return &self.pool[i].io_service,nil
		}
	}
	ioservice,err:=IOService_from(self.main_path+"/"+letter+".dsk")
	if err != nil{return nil,err}
	self.pool = append(self.pool, Disk{
		letter:     letter,
		io_service: ioservice,
	})
	return &self.pool[len(self.pool)-1].io_service,nil
}
func (self *IOServicePool) Flush_changes()(){
	for i := 0; i < len(self.pool); i++ {
		if self.pool[i].io_service.has_changes{
			self.pool[i].io_service.Flush()
		}
	}
}