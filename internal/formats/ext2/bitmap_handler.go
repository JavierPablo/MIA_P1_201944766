package ext2

import (
	"fmt"
	"project/internal/datamanagment"
	"project/internal/types"
)
type Bitmap struct{
	super_service *datamanagment.IOService
	index int32
	length int32
	free_transformer types.Integer
	free int32
	chunks []Chunk
}
func New_Bitmap(super_service *datamanagment.IOService,index int32,length int32,free types.Integer) Bitmap{
	return Bitmap{
		super_service: super_service,
		index:         index,
		length:        length,
		free:          free.Get(),
		free_transformer:          free,
		chunks:        []Chunk{},
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
}
func (self *Bitmap) Init_mapping(){
	bitmap := self.super_service.Read(self.length,self.index)
	chunks := make([]Chunk,0,10)
	current_chunk := Chunk{bit_no: 0,length: 0}
	for n,b := range *bitmap{
		if b == 0{
			current_chunk.length ++
		}else{
			if current_chunk.length == 0{
				current_chunk.bit_no = int32(n) + 1
				continue
			}
			chunks = append(chunks, current_chunk)
			current_chunk = Chunk{bit_no: int32(n) + 1,length: 0}
		}
	}
	if current_chunk.length != 0{
		chunks = append(chunks, current_chunk)
	}
	self.chunks = chunks
}

type Chunk struct{
	bit_no int32
	length int32
}
func (self *Bitmap) Best_fit(for_length int32) int32{
	candidate := 0
	candidate_exist := false
	for i := 0; i < len(self.chunks); i++ {
		if for_length <= self.chunks[i].length &&
		 self.chunks[candidate].length >= self.chunks[i].length{
			candidate = i
			candidate_exist = true
		}
	}
	if candidate_exist{
		return self.ocupe_chunk(candidate,for_length)
	}
	return -1
}
func (self *Bitmap) Worst_fit(for_length int32) int32{
	candidate := 0
	candidate_exist := false
	for i := 0; i < len(self.chunks); i++ {
		if for_length <= self.chunks[i].length &&
		 self.chunks[candidate].length <= self.chunks[i].length{
			candidate = i
			candidate_exist = true
		}
	}
	if candidate_exist{
		return self.ocupe_chunk(candidate,for_length)
	}
	return -1
}
func (self *Bitmap) First_fit(for_length int32) int32{
	for i := 0; i < len(self.chunks); i++ {
		if for_length <= self.chunks[i].length {
			return self.ocupe_chunk(i,for_length)
		}
	}
	return -1
}
func (self *Bitmap) ocupe_chunk(index int,for_length int32)int32{
	self.set_in_bitmap(for_length,self.chunks[index].bit_no,1)
	bit_no := self.chunks[index].bit_no
	self.chunks[index].bit_no+=for_length
	self.chunks[index].length-=for_length
	if self.chunks[index].length == 0{
		self.chunks = append((self.chunks)[:index], (self.chunks)[index+1:]...)
	}
	return bit_no
}

func (self *Bitmap) Erase(for_length int32, at_bit_no int32){
	for i := 0; i < len(self.chunks); i++ {
		if self.chunks[i].bit_no + self.chunks[i].length == at_bit_no {
			self.chunks[i].length += for_length			
			self.set_in_bitmap(for_length,at_bit_no,0)
			return
			}else if at_bit_no + for_length == self.chunks[i].bit_no{
				self.chunks[i].length += for_length			
				self.chunks[i].bit_no = at_bit_no 
				self.set_in_bitmap(for_length,at_bit_no,0)
				return
			}
	}

	closest_chunk_index := 0
	for i,chunk := range self.chunks{
		if at_bit_no < chunk.bit_no{
			closest_chunk_index = i
			break
		}
	}
	self.chunks = append(self.chunks[:closest_chunk_index+1], self.chunks[closest_chunk_index:]...)
	self.chunks[closest_chunk_index] = Chunk{
		bit_no: at_bit_no,
		length: for_length,
	}
	self.set_in_bitmap(for_length,at_bit_no,0)
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
	for _,b :=range self.chunks{
		fmt.Print("{")
		fmt.Print(b.bit_no)
		fmt.Print(",")
		fmt.Print(b.length)
		fmt.Print("},")
	}
	fmt.Println()
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