package datamanagment

import (
	"fmt"
	"os"
)

type ByteChunk struct{
	space Space
	buffer []byte
}
type IOService struct{
	fragments []ByteChunk
	// buffer	[]byte
	file_path string
	has_changes bool
	cont int
}



func IOService_from(file_path string) (IOService,error){
	_, err := os.Stat(file_path)
    if err != nil {return IOService{},err}
	
	return IOService{
		fragments:   []ByteChunk{},
		file_path:   file_path,
		has_changes: false,
	},nil
}
func (self *IOService) request_bytes(amount int32, at int32) []byte {
	file_dscrptr, err := os.Open(self.file_path)
	if err != nil {
		file_dscrptr.Close()
		panic(err)
	}
	buff:= make([]byte,amount)
	_,err = file_dscrptr.ReadAt(buff,int64(at))
	if err != nil {
		file_dscrptr.Close()
		panic(err)
	}
	file_dscrptr.Close()
	return buff
}

func (self *IOService) get_right_chunk(read_space Space,start_at int) *ByteChunk {
	// fmt.Printf("requested %s\n",read_space.Show())
	for i := start_at; i < len(self.fragments); i++ {
		fragment:=&self.fragments[i]
		sapce:=&fragment.space
		switch sapce.Contains(read_space){
		case ColindantOut:
			// fmt.Println("colindant")
			for j := i+1; j < len(self.fragments); j++ {
				fragment2:=&self.fragments[j]
				space2 := &fragment2.space
				switch space2.Contains(read_space){
				case Partial, ColindantOut:
					// fmt.Println("INcolindant partial,colindanout")

					amount_required:=space2.Index - read_space.Index
					new_data:=self.request_bytes(amount_required,read_space.Index)
					sapce.Length += amount_required + space2.Length
					fragment.buffer = append(fragment.buffer, new_data...)
					fragment.buffer = append(fragment.buffer, fragment2.buffer...)
					self.fragments = append(self.fragments[:j],self.fragments[j+1:]... )
					return fragment
				case Out:
					// fmt.Println("INcolindant out")

					new_data:=self.request_bytes(read_space.Length,read_space.Index)
					// fmt.Printf("new data length %d\n",len(new_data))
					// fmt.Printf("buff length before %d\n",len(fragment.buffer))
					fragment.buffer = append(fragment.buffer, new_data...)
					// fmt.Printf("buff length after %d\n",len(fragment.buffer))
					fragment.space.extend_with(read_space)
					return fragment
				case Cover:
					// fmt.Println("INcolindant cover")

					amount_required:=space2.Index - read_space.Index
					new_data:=self.request_bytes(amount_required,read_space.Index)
					sapce.Length += amount_required + space2.Length
					fragment.buffer = append(fragment.buffer, new_data...)
					fragment.buffer = append(fragment.buffer, fragment2.buffer...)
					self.fragments = append(self.fragments[:j],self.fragments[j+1:]... )
					read_space.Index = space2.Boundary()
					read_space.Length -= amount_required - space2.Length
					continue
				case Same,In,ColindantIn:
					panic("Should have matched any branch above")
				default: panic("Should have matched any branch above")
				}
			}
			new_data:=self.request_bytes(read_space.Length,read_space.Index)
			fragment.buffer = append(fragment.buffer, new_data...)
			fragment.space.extend_with(read_space)
			return fragment
		case Out:
			// fmt.Println("out")

			continue
		case Cover:
			// fmt.Println("cover")

			amount_required := fragment.space.Index - read_space.Index
			new_data:=self.request_bytes(amount_required,read_space.Index)
			fragment.buffer = append(new_data, fragment.buffer...)
			fragment.space.Index = read_space.Index
			fragment.space.Length += amount_required

			read_space.Index = fragment.space.Boundary()
			read_space.Length -= read_space.Length
			return self.get_right_chunk(read_space,i)
			
		case Partial:
			// fmt.Println("partial")

			if fragment.space.Index > read_space.Index{
				// fmt.Println("partial case 1")
				amount_required := fragment.space.Index - read_space.Index
				new_data:=self.request_bytes(amount_required,read_space.Index)
				fragment.buffer = append(new_data, fragment.buffer...)
				fragment.space.Index = read_space.Index
				fragment.space.Length +=  amount_required
				return fragment
			}else{
				// fmt.Println("partial case 2")
				amount_diff := fragment.space.Boundary() - read_space.Index
				read_space.Index = fragment.space.Boundary()
				read_space.Length -= amount_diff
				return self.get_right_chunk(read_space,i)
			}
		case ColindantIn,In,Same:
			// fmt.Println("colindantin,in,same")

			return fragment
		}
	}
	// fmt.Println("new")

	new_data:=self.request_bytes(read_space.Length,read_space.Index)
	self.fragments = append(self.fragments, ByteChunk{
		space:  read_space,
		buffer: new_data,
	})
	return &self.fragments[len(self.fragments)-1]
}
func (self *IOService) Read(amount int32, at int32) *[]byte {
	var read_space Space
	if amount < 1024 {
		read_space = New_Space(at,1024)
	}else{
		read_space = New_Space(at,amount)
	}
	// fmt.Printf("FOR READING ++++++++++++++++++++++++++++++++++ with size = %d\n",amount)
	chunk:=self.get_right_chunk(read_space,0)
	// fmt.Printf("FOR READ ++++++++++++++++++++++++++++++++++ got = %s length = %d with buff = %d\n",chunk.space.Show(),chunk.space.Length,len(chunk.buffer))

	new_at:=at-chunk.space.Index
	chunk_to_return := (chunk.buffer)[new_at:new_at+amount]
	if self.cont != len(self.fragments){
		// self.Log_chunks()
		self.cont = len(self.fragments)
	}
	return &chunk_to_return
}
func (self *IOService) Write(content []byte, at int32) {
	self.has_changes = true
	var read_space Space
	if len(content) < 1024 {
		read_space = New_Space(at,1024)
	}else{
		read_space = New_Space(at,int32(len(content)))
	}
	// fmt.Printf("FOR WRITING ++++++++++++++++++++++++++++++++++ with size = %d  at = %d\n",len(content),at)
	chunk:=self.get_right_chunk(read_space,0)
	// fmt.Printf("FOR WRITING ++++++++++++++++++++++++++++++++++ got = %s length = %d with buff = %d\n",chunk.space.Show(),chunk.space.Length,len(chunk.buffer))
	new_at:=at-chunk.space.Index
	for n,b := range content{
		chunk.buffer[new_at+int32(n)] = b
	}
	if self.cont != len(self.fragments){
		// self.Log_chunks()
		self.cont = len(self.fragments)
	}
}

func (self *IOService) Flush(){	
	file_dscrptr, err := os.OpenFile(self.file_path,os.O_WRONLY,os.ModeAppend)
	if err != nil {
		file_dscrptr.Close()
		panic(err)
	}
	for i := 0; i < len(self.fragments); i++ {
		fragment:=self.fragments[i]
		_,err = file_dscrptr.WriteAt(fragment.buffer,int64(fragment.space.Index))
	   if err != nil {
		   file_dscrptr.Close()
		   panic(err)
	   }
	}
	file_dscrptr.Close()
}
func (self *IOService) Log_chunks(){
	fmt.Printf("{ ")
	for i := 0; i < len(self.fragments); i++ {
		fmt.Printf("%s,",self.fragments[i].space.Show())
	}
	fmt.Printf(" }\n")
}
