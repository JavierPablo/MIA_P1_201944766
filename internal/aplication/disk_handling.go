package aplication

import (
	"os"
	"project/internal/utiles"
	"project/internal/types"
	"project/internal/datamanagment"
)
const DISK_CONTAINER_PATH string = "./MIA/P1"
func (self *Aplication) Make_disk(size int32, unite utiles.SizeUnit, fit utiles.FitCriteria) (string,types.MasterBootRecord){
	var file *os.File
	ascii := 65
	path := DISK_CONTAINER_PATH+"/"+string(ascii)
	file, err :=os.Open(path)
	for err == nil{
		file.Close()
		ascii += 1
		if ascii == 91{panic("There are no abc names available for disks")}
		path = DISK_CONTAINER_PATH+"/"+string(ascii)
		file, err = os.Open(path)
	}
	file, err = os.Create(path)
	if err != nil {panic("Problems creating disk even thouhg the name is available")}
	defer file.Close()
	total_bytes := int(size*int32(unite))
	for i := 0; i < total_bytes; i++ {
		_, err = file.Write([]byte{0})
		if err != nil {panic("Problems writing 0s in disk")}
	}
	io_service:=datamanagment.IOService_from(path)	
	mbr := types.CreateMasterBootRecord(&io_service,0)
	mbr.Mbr_tamano().Set(int32(total_bytes))
	mbr.Mbr_fecha_creacion().Set(utiles.Current_Time())
	mbr.Dsk_fit().Set(string(fit))
	for _,partition := range mbr.Mbr_partitions().Spread(){
		partition.Part_start().Set(-1)
	}
	
	return path,mbr
}

func (self *Aplication) Remove_disk(letter string) bool{
	path := DISK_CONTAINER_PATH+"/"+letter
	err := os.Remove(path)
	return err == nil

}

func (self *Aplication) Partition_disk(size int32, letter string, 
	p_name string, unite utiles.SizeUnit, partition_type utiles.PartitionType,
	fit utiles.FitCriteria) {
	path:=DISK_CONTAINER_PATH+"/"+letter
	_, err :=os.Open(path)
	if err != nil{
		return // disk doesnt exist
	}
	io_service:=datamanagment.IOService_from(path)	

	mbr := types.CreateMasterBootRecord(&io_service,0)
	// partition_size := size*int32(unite)
	if partition_type == utiles.Primary{
		for _,partition := range mbr.Mbr_partitions().Spread(){
			if partition.Part_start().Get() != -1{continue}
			
		}
		panic("There are no partition spaces available")

	}else if partition_type == utiles.Extendend{

	}else if partition_type == utiles.Logic{

	}

}