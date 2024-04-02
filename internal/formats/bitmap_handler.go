package formats

import (
	"fmt"
	"project/internal/datamanagment"
	"project/internal/types"
)
type Bitmap struct{
	super_service *datamanagment.IOService
	length int32
	index int32
	space_manager datamanagment.SpaceManager
	free_transformer types.Integer
	free int32
}
func (self *Bitmap) Get_SpaceManager()datamanagment.SpaceManager{return self.space_manager}
func New_Bitmap(super_service *datamanagment.IOService,index int32,length int32,free types.Integer) Bitmap{
	return Bitmap{
		super_service:    super_service,
		length:           length,
		index:            index,
		space_manager:    datamanagment.Empty_SpaceManager_from(length),
		free_transformer: free,
		free:             free.Get(),
	}
}

func (self *Bitmap) Clear(){
	self.free = self.length
	self.free_transformer.Set(self.free)
	bytes_to_write := make([]byte,self.length)
	for i := 0; i < int(self.length); i++ {
		bytes_to_write[i] = 0
	}
	self.super_service.Write(bytes_to_write,self.index)
	self.space_manager = datamanagment.Empty_SpaceManager_from(self.length)
}
func (self *Bitmap) Init_mapping(){
	bitmap := self.super_service.Read(self.length,self.index)
	chunks := make([]datamanagment.Space,0,10)
	current_chunk := datamanagment.New_Space(0,0)
	for n,b := range *bitmap{
		if b == 0{
			current_chunk.Length ++
		}else{
			if current_chunk.Length == 0{
				current_chunk.Index = int32(n) + 1
				continue
			}
			chunks = append(chunks, current_chunk)
			current_chunk = datamanagment.New_Space(int32(n)+1,0)
			//  Chunk{bit_no: int32(n) + 1,length: 0}
		}
	}
	if current_chunk.Length != 0{
		chunks = append(chunks, current_chunk)
	}
	self.space_manager = datamanagment.SpaceManager_from_free_spaces(chunks,self.length)
}

func (self *Bitmap) Best_fit(for_length int32) int32{
	space_indx := self.space_manager.Best_fit(for_length)
	if space_indx == -1{return -1}
	bit_no := self.space_manager.Chunk_no(space_indx).Index
	self.space_manager.Ocupe_space_unchecked(int(space_indx),for_length)
	self.set_in_bitmap(for_length,bit_no,1)
	return bit_no
}
func (self *Bitmap) Worst_fit(for_length int32) int32{
	space_indx := self.space_manager.Worst_fit(for_length)
	if space_indx == -1{return -1}
	bit_no := self.space_manager.Chunk_no(space_indx).Index
	self.space_manager.Ocupe_space_unchecked(int(space_indx),for_length)
	self.set_in_bitmap(for_length,bit_no,1)
	return bit_no
}
func (self *Bitmap) First_fit(for_length int32) int32{
	space_indx := self.space_manager.First_fit(for_length)
	if space_indx == -1{return -1}
	bit_no := self.space_manager.Chunk_no(space_indx).Index
	self.space_manager.Ocupe_space_unchecked(int(space_indx),for_length)
	self.set_in_bitmap(for_length,bit_no,1)
	return bit_no
}
func (self *Bitmap) Erase(for_length int32, at_bit_no int32)error{
	err:=self.space_manager.Free_space(for_length,at_bit_no)
	if err==nil{
		self.set_in_bitmap(for_length,at_bit_no,0)
	}
	return nil
}

func (self *Bitmap) set_in_bitmap(amount int32, at_bit_no int32, data byte){
	// Debug porpouses
	for _,b := range *self.super_service.Read(amount,self.index+at_bit_no){
		if b == data {panic("Setting same bit again")}
	}
	// Debug porpouses
	
	if data == 1{
		self.free -= amount
	}else if data == 0{
		self.free += amount
	}
	self.free_transformer.Set(self.free)
	bytes_to_write := make([]byte,amount)
	for i := 0; i < int(amount); i++ {
		bytes_to_write[i] = data
	}
	self.super_service.Write(bytes_to_write,self.index+at_bit_no)
}


func (self *Bitmap) Log_chunks_state(){
	self.space_manager.Log_chunks_state()
}
func (self Bitmap) Log_bitmap_state(){
	bytes := self.super_service.Read(self.length,self.index)
	fmt.Print("\"")
	for _,b :=range *bytes{
		fmt.Print(int(b))
		fmt.Print(",")
	}
	fmt.Println("\"")
}