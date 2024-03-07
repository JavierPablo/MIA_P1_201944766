package utiles

import (
	"fmt"
)
type SpaceManager struct{
	index int32
	length int32
	free_spaces []Space
}
type Space struct{
	index int32
	length int32
}
func New_Space(indx int32,length int32) Space{return Space{
	index: indx,
	length:         length,
}}
type SpaceComparation string
const (
	In  SpaceComparation = "In"
	Out SpaceComparation = "Out"
	Partial SpaceComparation = "Partial"
	ColindantOut SpaceComparation = "ColindantOut"
	ColindantIn SpaceComparation = "ColindantIn"
	Same SpaceComparation = "Same"
	Cover SpaceComparation = "Cover"
)
func (self *Space) Boundary() int32{return self.index + self.length}
func (self *Space) Show() string{return fmt.Sprintf("[%d,%d]",self.index,self.Boundary()-1)}
func (self *Space) Contains(sp Space) SpaceComparation{

	if self.index <= sp.index && self.Boundary() >= sp.Boundary(){
		if self.index == sp.index{
			if self.Boundary() == sp.Boundary(){return Same} 
			return ColindantIn
		}
		if self.Boundary() == sp.Boundary(){
			return ColindantIn
		}
		return In
	}

	if self.index == sp.Boundary() || self.Boundary() == sp.index{
		return ColindantOut
	}

	if self.index <= sp.index && self.Boundary() > sp.index &&
	 self.Boundary() < sp.Boundary(){
		return Partial
	}
	if self.index > sp.index && self.index < sp.Boundary() &&
	 self.Boundary() >= sp.Boundary(){
		return Partial
	}
	if self.index > sp.index && self.Boundary() < sp.Boundary(){
		return Cover
	}
	return Out

}
func (self *Space) extend_with(sp Space){
	if self.index == sp.Boundary(){
		self.index = sp.index
		self.length += sp.length
		return
	} else if self.Boundary() == sp.index{
		self.length += sp.length
		return
	}
	panic(fmt.Sprintf("Can not extend %s with %s",self.Show(),sp.Show()))
}
func (self *Space) reduce_with(sp Space){
	if self.index == sp.index{
		if self.Boundary() == sp.Boundary(){
			panic(fmt.Sprintf("Can not reduce since both are equal %s AND %s",self.Show(),sp.Show()))
		} 
		self.length -=sp.length
	}
	if self.Boundary() == sp.Boundary(){
		self.index = sp.index
		self.length -=sp.length
	}
	panic(fmt.Sprintf("Can not reduce %s with %s",self.Show(),sp.Show()))
}
func (self *Space) Split(sp Space)(Space,Space){
	if self.index < sp.index && self.Boundary() > sp.Boundary(){
		first_half := Space{
			index: self.index,
			length: sp.index - self.index,
		}
		second_half := Space{
			index: sp.Boundary(),
			length: self.Boundary() - sp.Boundary(),
		}
		return first_half, second_half
	}
	panic(fmt.Sprintf("Can not split %s using %s",self.Show(),sp.Show()))
}















func Spacemanager_from_free_spaces(free_spaces []Space,total_length int32) SpaceManager{
	new_spman := SpaceManager{
		index:       0,
		length:      total_length,
		free_spaces: free_spaces,
	}
	return new_spman
}
func Spacemanager_from_occuped_spaces(occuped_spaces []Space,total_length int32) SpaceManager{
	new_spman := SpaceManager{
		index:       0,
		length:      total_length,
		free_spaces: []Space{Space{
			index: 0,
			length:         total_length,
		}},
	}
	for _, occuped := range occuped_spaces {
		succed := new_spman.Ocupe_raw_space(occuped.length,occuped.index)
		if !succed{
			panic("")
		}
	}
	return new_spman
}
func (self *SpaceManager) Chunk_no(i int32) *Space{
	return &self.free_spaces[i]
}
func (self *SpaceManager) Best_fit(for_length int32) int32{
	candidate := 0
	candidate_exist := false
	for i := 0; i < len(self.free_spaces); i++ {
		if for_length <= self.free_spaces[i].length &&
		 self.free_spaces[candidate].length >= self.free_spaces[i].length{
			candidate = i
			candidate_exist = true
		}
	}
	if candidate_exist{
		return int32(candidate)
	}
	return -1
}
func (self *SpaceManager) Worst_fit(for_length int32) int32{
	candidate := 0
	candidate_exist := false
	for i := 0; i < len(self.free_spaces); i++ {
		if for_length <= self.free_spaces[i].length &&
		 self.free_spaces[candidate].length <= self.free_spaces[i].length{
			candidate = i
			candidate_exist = true
		}
	}
	if candidate_exist{
		return int32(candidate)
	}
	return -1
}
func (self *SpaceManager) First_fit(for_length int32) int32{
	for i := 0; i < len(self.free_spaces); i++ {
		if for_length <= self.free_spaces[i].length {
			return int32(i)
		}
	}
	return -1
}
func (self *SpaceManager) Ocupe_raw_space(for_length int32, at_bit_no int32)bool{
	trgt_space := New_Space(at_bit_no,for_length)
	if trgt_space.Boundary() > self.length {
		panic("The requested space is passing the bondarys of the current arena")
	}
	for i := 0; i < len(self.free_spaces); i++ {
		result := self.free_spaces[i].Contains(trgt_space)
		switch result {
			case ColindantIn:
				self.free_spaces[i].reduce_with(trgt_space)
				return true
			case In:
				first,last:=self.free_spaces[i].Split(trgt_space)
				self.free_spaces[i] = first
				self.free_spaces = append((self.free_spaces)[:i+1], (self.free_spaces)[i:]...)
				self.free_spaces[i+1] = last
				return true
			case Same:
				self.free_spaces = append((self.free_spaces)[:i], (self.free_spaces)[i+1:]...)
				return true
			case Partial:
				fmt.Printf("%s %s with %s\n",self.free_spaces[i].Show(),self.free_spaces[i].Contains(trgt_space),trgt_space.Show())
				panic("The requested space is invading other unexpected areas, and canot be properly handled")
			case Out:
				continue
			default:
				panic("Anything else is expected, since all cases are handeled, this means that the data is corrupted")
		}	
	}
	panic("The space requested to ocupe has fail")
}
func (self *SpaceManager) Ocupe_space_unchecked(space_no int,for_length int32)int32{
	bit_no := self.free_spaces[space_no].index
	self.free_spaces[space_no].index+=for_length
	self.free_spaces[space_no].length-=for_length
	if self.free_spaces[space_no].length == 0{
		self.free_spaces = append((self.free_spaces)[:space_no], (self.free_spaces)[space_no+1:]...)
	}
	return bit_no
}

func (self *SpaceManager) Free_space(for_length int32, at_bit_no int32)bool{
	trgt_space := New_Space(at_bit_no,for_length)
	closest := len(self.free_spaces)
	if trgt_space.Boundary() > self.length {
		panic("The requested space is passing the bondarys of the current arena")
	}
	for i := 0; i < len(self.free_spaces); i++ {
		result := self.free_spaces[i].Contains(trgt_space)
		switch result {
			case ColindantOut:
				self.free_spaces[i].extend_with(trgt_space)
				if i+1 != len(self.free_spaces){
					result = self.free_spaces[i].Contains(self.free_spaces[i+1])
					if result == ColindantOut{
						self.free_spaces[i].extend_with(self.free_spaces[i+1])
						self.free_spaces = append((self.free_spaces)[:i+1], (self.free_spaces)[i+2:]...)
					}else if result != Out{panic("Abnormal result")}
				}
				return true
			case Out: 
				if trgt_space.index < self.free_spaces[i].index{
				closest = i}
				continue
			default:
				panic(fmt.Sprintf("Anything else is expected, since all cases are handeled but got %s in %s by %s, this means that the data is corrupted",result,self.free_spaces[i].Show(),trgt_space.Show()))
		}	
	}
	self.free_spaces = append(self.free_spaces[:closest+1], self.free_spaces[closest:]...)
	self.free_spaces[closest] = trgt_space
	return true

	// panic("The space requested to ocupe has fail")
	// for i := 0; i < len(self.free_spaces); i++ {
	// 	if self.free_spaces[i].index + self.free_spaces[i].length == at_bit_no {
	// 		self.free_spaces[i].length += for_length			
	// 		return true
	// 	}else if at_bit_no + for_length == self.free_spaces[i].index{
	// 		self.free_spaces[i].length += for_length			
	// 		self.free_spaces[i].index = at_bit_no 
	// 		return true
	// 	}
	// }

	// closest_chunk_index := len(self.free_spaces) 
	// for i,chunk := range self.free_spaces{
	// 	if at_bit_no < chunk.index{
	// 		closest_chunk_index = i
	// 		if at_bit_no + for_length < chunk.index{break}
	// 		panic(fmt.Sprintf("Space collition detected at %d inside [%d,%d] bounds, caused by requested space liberation [%d,%d]",
	// 		at_bit_no + for_length,
	// 		chunk.index,
	// 		chunk.length + chunk.index,
	// 		at_bit_no,for_length+at_bit_no))
	// 	}else if at_bit_no == chunk.index{
	// 		panic(fmt.Sprintf("Space collition detected at %d inside [%d,%d] bounds, caused by requested space liberation [%d,%d]",
	// 		at_bit_no,
	// 		chunk.index,
	// 		chunk.length + chunk.index,
	// 		at_bit_no,for_length+at_bit_no))
	// 	}
	// }
	// if at_bit_no + for_length > self.length{
	// 	panic(fmt.Sprintf("Attempting to erase space out of arena bondaries (overflow at %d for max limit %d",
	// 	at_bit_no + for_length,self.length))
	// }

	// self.free_spaces = append(self.free_spaces[:closest_chunk_index+1], self.free_spaces[closest_chunk_index:]...)
	// self.free_spaces[closest_chunk_index] = Space{
	// 	index: at_bit_no,    
	// 	length: for_length,	
	// }
	// return true
}

func (self *SpaceManager) Log_chunks_state(){
	for _,b :=range self.free_spaces{
		// fmt.Printf("{at=%d,len(%d)},",b.relative_index,b.length)
		fmt.Printf("[%d,%d],",b.index,b.length+b.index-1)
	}
	fmt.Println()
}