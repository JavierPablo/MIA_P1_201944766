package datamanagment

import (
	"fmt"
	// "project/internal/utiles"
)
type SpaceManager struct{
	index int32
	length int32
	free_spaces []Space
}
func (self *SpaceManager) Get_length()int32{return self.length}
type Space struct{
	Index int32
	Length int32
}
func New_Space(indx int32,length int32) Space{return Space{
	Index: indx,
	Length:         length,
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
func (self *Space) Boundary() int32{return self.Index + self.Length}
func (self *Space) Show() string{return fmt.Sprintf("[%d,%d]",self.Index,self.Boundary()-1)}
func (self *Space) Contains(sp Space) SpaceComparation{

	if self.Index <= sp.Index && self.Boundary() >= sp.Boundary(){
		if self.Index == sp.Index{
			if self.Boundary() == sp.Boundary(){return Same} 
			return ColindantIn
		}
		if self.Boundary() == sp.Boundary(){
			return ColindantIn
		}
		return In
	}

	if self.Index == sp.Boundary() || self.Boundary() == sp.Index{
		return ColindantOut
	}

	if self.Index <= sp.Index && self.Boundary() > sp.Index &&
	 self.Boundary() < sp.Boundary(){
		return Partial
	}
	if self.Index > sp.Index && self.Index < sp.Boundary() &&
	 self.Boundary() >= sp.Boundary(){
		return Partial
	}
	if self.Index > sp.Index && self.Boundary() < sp.Boundary(){
		return Cover
	}
	return Out

}
func (self *Space) extend_with(sp Space){
	if self.Index == sp.Boundary(){
		self.Index = sp.Index
		self.Length += sp.Length
		return
	} else if self.Boundary() == sp.Index{
		self.Length += sp.Length
		return
	}
	panic(fmt.Sprintf("Can not extend %s with %s",self.Show(),sp.Show()))
}
func (self *Space) reduce_with(sp Space){
	if self.Index == sp.Index{
		if self.Boundary() == sp.Boundary(){
			panic(fmt.Sprintf("Can not reduce since both are equal %s AND %s",self.Show(),sp.Show()))
		} 
		self.Index = sp.Boundary()
		self.Length -=sp.Length
		return
	}
	if self.Boundary() == sp.Boundary(){
		self.Length -=sp.Length
		return
	}
	panic(fmt.Sprintf("Can not reduce %s with %s",self.Show(),sp.Show()))
}
func (self *Space) Split(sp Space)(Space,Space){
	if self.Index < sp.Index && self.Boundary() > sp.Boundary(){
		first_half := Space{
			Index: self.Index,
			Length: sp.Index - self.Index,
		}
		second_half := Space{
			Index: sp.Boundary(),
			Length: self.Boundary() - sp.Boundary(),
		}
		return first_half, second_half
	}
	panic(fmt.Sprintf("Can not split %s using %s",self.Show(),sp.Show()))
}















func Empty_SpaceManager_from(total_length int32) SpaceManager{
	new_spman := SpaceManager{
		index:       0,
		length:      total_length,
		free_spaces: []Space{},
	}
	return new_spman
}
func SpaceManager_from_free_spaces(free_spaces []Space,total_length int32) SpaceManager{
	new_spman := SpaceManager{
		index:       0,
		length:      total_length,
		free_spaces: free_spaces,
	}
	return new_spman
}
func SpaceManager_from_occuped_spaces(occuped_spaces []Space,total_length int32) (SpaceManager,error){
	new_spman := SpaceManager{
		index:       0,
		length:      total_length,
		free_spaces: []Space{Space{
			Index: 0,
			Length:         total_length,
		}},
	}
	for _, occuped := range occuped_spaces {
		err := new_spman.Ocupe_raw_space(occuped.Length,occuped.Index)
		if err !=nil {return SpaceManager{},err}
	}
	return new_spman,nil
}
func (self *SpaceManager) Chunk_no(i int32) *Space{
	return &self.free_spaces[i]
}
func (self *SpaceManager) Best_fit(for_length int32) int32{
	candidate := 0
	candidate_exist := false
	for i := 0; i < len(self.free_spaces); i++ {
		if for_length <= self.free_spaces[i].Length &&
		 self.free_spaces[candidate].Length >= self.free_spaces[i].Length{
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
		if for_length <= self.free_spaces[i].Length &&
		 self.free_spaces[candidate].Length <= self.free_spaces[i].Length{
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
		if for_length <= self.free_spaces[i].Length {
			return int32(i)
		}
	}
	return -1
}
func (self *SpaceManager) Ocupe_raw_space(for_length int32, at_bit_no int32)error{
	trgt_space := New_Space(at_bit_no,for_length)
	if trgt_space.Boundary() > self.length {
		return fmt.Errorf("The requested space is passing the bondarys of the current arena")
	}
	for i := 0; i < len(self.free_spaces); i++ {
		result := self.free_spaces[i].Contains(trgt_space)
		switch result {
			case ColindantIn:
				self.free_spaces[i].reduce_with(trgt_space)
				return nil
			case In:
				first,last:=self.free_spaces[i].Split(trgt_space)
				self.free_spaces[i] = first
				self.free_spaces = append((self.free_spaces)[:i+1], (self.free_spaces)[i:]...)
				self.free_spaces[i+1] = last
				return nil
			case Same:
				self.free_spaces = append((self.free_spaces)[:i], (self.free_spaces)[i+1:]...)
				return nil
			case Partial:
				
				return fmt.Errorf("%s %s with %s\n%s\n",self.free_spaces[i].Show(),self.free_spaces[i].Contains(trgt_space),trgt_space.Show(),
				"The requested space is invading other unexpected areas, and canot be properly handled")
			case Out:
				continue
			default:
				return fmt.Errorf("Anything else is expected, since all cases are handeled, this means that the data is corrupted")
				
		}	
	}
	return fmt.Errorf("The space requested to ocupe has fail")
	
}
func (self *SpaceManager) Ocupe_space_unchecked(space_no int,for_length int32)int32{
	bit_no := self.free_spaces[space_no].Index
	self.free_spaces[space_no].Index+=for_length
	self.free_spaces[space_no].Length-=for_length
	if self.free_spaces[space_no].Length == 0{
		self.free_spaces = append((self.free_spaces)[:space_no], (self.free_spaces)[space_no+1:]...)
	}
	return bit_no
}

func (self *SpaceManager) Free_space(for_length int32, at_bit_no int32)error{
	trgt_space := New_Space(at_bit_no,for_length)
	closest := len(self.free_spaces)
	if trgt_space.Boundary() > self.length {
		return fmt.Errorf("The requested space is passing the bondarys of the current arena")
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
					}else if result != Out{
						return fmt.Errorf("Abnormal result for %s by %s with case %s",self.free_spaces[i].Show(),trgt_space.Show(),result)
					}
				}

				return nil
			case Out: 
				if trgt_space.Index < self.free_spaces[i].Index{
				closest = i}
				continue
			default:
				return fmt.Errorf("Anything else is expected, since all cases are handeled but got %s in %s by %s, this means that the data is corrupted",result,self.free_spaces[i].Show(),trgt_space.Show())
				

		}	
	}
	self.free_spaces = append(self.free_spaces[:closest+1], self.free_spaces[closest:]...)
	self.free_spaces[closest] = trgt_space
	return nil
}









func (self *SpaceManager) Free_spaces()[]Space{
	return self.free_spaces
}
func (self *SpaceManager) Log_chunks_state(){
	for _,b :=range self.free_spaces{
		// fmt.Printf("{at=%d,len(%d)},",b.relative_index,b.length)
		fmt.Printf("[%d,%d],",b.Index,b.Length+b.Index-1)
	}
	fmt.Println()
}
