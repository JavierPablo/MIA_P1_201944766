package formats

import (
	"fmt"
	"project/internal/datamanagment"
	"project/internal/types"
	"strconv"
	"strings"
	"time"
)

type JournalingMember string
const (
	Login JournalingMember = "l"
	Unlog JournalingMember = "j"
	Make_group JournalingMember = "1"
	Remove_group JournalingMember = "2"
	Make_user JournalingMember = "3"
	Remove_user JournalingMember = "4"
	Make_file JournalingMember = "5"
	Remove JournalingMember = "6"
	Edit_file JournalingMember = "7"
	Rename_inode JournalingMember = "8"
	Make_dir JournalingMember = "9"
	Copy JournalingMember = "a"
	Move JournalingMember = "b"
	Chown JournalingMember = "c"
	Chgrp JournalingMember = "d"
	Chmod JournalingMember = "e"
)
type JournalingManager struct {
	io_service *datamanagment.IOService
	init_index int32
	length int32
	top_helper DataExtractionHelper
	elm_count int32
	full bool
}

func New_Journaling(io_service *datamanagment.IOService,index int32,length int32)JournalingManager{
	return JournalingManager{
		io_service: io_service,
		init_index: index,
		length:     length,
		top_helper: DataExtractionHelper{},
		elm_count:  0,
		full:       false,
	}
}
func (self *JournalingManager) Generate_dot_rep()string{
	normalize := func(txt string)string{
		str:=""
		for i := 0; i < len(txt); i++ {
			char:=string(txt[i])
			ascii_val:=txt[i]
			if ascii_val >= 32 && ascii_val <= 126{
				switch ascii_val{
				case 34:
					char = `\"`
				case 92:
					char = `\\`
				case 123:
					char = `\{`
				case 125:
					char = `\}`
				case 124:
					char = `\|`
				case 60:
					char = `\<`
				case 62:
					char = `\>`
				}
			}else{
				switch ascii_val{
				case 10:
					char = `\\n`
				case 11:
					char = `\\t`
				default: 
					char = "-"
				}
			}
			str+=char
		}
		return strings.TrimSpace(str)
	}
	builder:=""
	iterator:=self.Instructions_iter()
	for iterator.Has_next(){
		ins:=iterator.Next()
		switch ins.Member_type {
		case Login:
			builder+=fmt.Sprintf(`<TR><TD>Login</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>user</TD>
						<TD>pass</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],ins.Attrs[1])
		case Unlog:
			builder+=fmt.Sprintf(`<TR><TD>Unlog</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD></TD>
					</TR></TABLE></TD>`)
		case Make_group:
			builder+=fmt.Sprintf(`<TR><TD>Mkgroup</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>name</TD>
					</TR><TR>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0])
		case Remove_group:
			builder+=fmt.Sprintf(`<TR><TD>Rmgroup</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>name</TD>
					</TR><TR>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0])
		case Make_user:
			builder+=fmt.Sprintf(`<TR><TD>Mkuser</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>user</TD>
						<TD>pass</TD>
						<TD>grp</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],ins.Attrs[1],ins.Attrs[2])
		case Remove_user:
			builder+=fmt.Sprintf(`<TR><TD>Rmuser</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>user</TD>
					</TR><TR>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0])
		case Make_file:
			builder+=fmt.Sprintf(`<TR><TD>Mkfile</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>recursive</TD>
						<TD>path</TD>
						<TD>content</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],ins.Attrs[1],normalize(ins.Attrs[2]))
		case Remove:
			builder+=fmt.Sprintf(`<TR><TD>Remove</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>path</TD>
					</TR><TR>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0])
		case Edit_file:
			builder+=fmt.Sprintf(`<TR><TD>Edit</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>path</TD>
						<TD>content</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],normalize(ins.Attrs[1]))
		case Rename_inode:
			builder+=fmt.Sprintf(`<TR><TD>Rename</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>path</TD>
						<TD>name</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],ins.Attrs[1])
		case Make_dir:
			builder+=fmt.Sprintf(`<TR><TD>Mkdir</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>recursive</TD>
						<TD>path</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],ins.Attrs[1])
		case Copy:
			builder+=fmt.Sprintf(`<TR><TD>Copy</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>path</TD>
						<TD>path-destino</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],ins.Attrs[1])
		case Move:
			builder+=fmt.Sprintf(`<TR><TD>Move</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>path</TD>
						<TD>path-destino</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],ins.Attrs[1])
		case Chown:
			builder+=fmt.Sprintf(`<TR><TD>Chown</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>recursive</TD>
						<TD>path</TD>
						<TD>user</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],ins.Attrs[1],ins.Attrs[2])
		case Chgrp:
			builder+=fmt.Sprintf(`<TR><TD>Chgrp</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>user</TD>
						<TD>grp</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],ins.Attrs[1])
		case Chmod:
			builder+=fmt.Sprintf(`<TR><TD>Chmod</TD><TD ><TABLE CELLSPACING="0" BORDER="0" CELLBORDER = "1"  >
					<TR>
						<TD>recursive</TD>
						<TD>user</TD>
						<TD>path</TD>
					</TR><TR>
						<TD>%s</TD>
						<TD>%s</TD>
						<TD>%s</TD>
					</TR></TABLE></TD>`,ins.Attrs[0],ins.Attrs[1],ins.Attrs[1])
		default: panic("Corrupted")
		}
		builder+=fmt.Sprintf("<TD>%s</TD></TR>",ins.Attrs[len(ins.Attrs)-1])
	}
	return builder
}
func (self *JournalingManager) Restart_count(){
	self.top_helper = New_DataExtractionHelper(self.io_service,self.init_index)
	self.top_helper.Set_int(0)
}
func (self *JournalingManager) Reset_cursor(){
	self.top_helper.cursor = self.init_index + 4
}
func (self *JournalingManager) Init_mapping(go_to_top bool){
	self.top_helper = New_DataExtractionHelper(self.io_service,self.init_index)
	self.elm_count = self.top_helper.Extract_int()
	if go_to_top{
		self.top_helper.Advance_items(self.elm_count)
	}
}
func (self *JournalingManager) Push_instruction(ins Instruction){
	if self.full {return}
	start:=0
	must_set_recursive:=false
	switch ins.Member_type {
	case Make_file:
		must_set_recursive=true
		start = 1
	case Make_dir:
		must_set_recursive=true
		start = 1
	case Chown:
		must_set_recursive=true
		start = 1
	case Chmod:
		must_set_recursive=true
		start = 1
	}
	full_length:=1 + DATE_LENGTH
	for i := start; i < len(ins.Attrs); i++ {
		full_length += 4 + len(ins.Attrs[i])
	}
	if must_set_recursive{full_length += 1}
	if full_length >=  (int(self.init_index) + int(self.length)) - int(self.top_helper.cursor){
		self.full = true
		return
	}
	self.top_helper.Set_char(string(ins.Member_type))	
	for i := start; i < len(ins.Attrs); i++ {
		self.top_helper.Set_int(int32(len(ins.Attrs[i])))
	}
	if must_set_recursive{
		self.top_helper.Set_char(ins.Attrs[0])
	}
	for i := start; i < len(ins.Attrs); i++ {
		self.top_helper.Set_string(ins.Attrs[i])
	}
	conv:=func(no int)string{
		str:=strconv.Itoa(no)
		if len(str) == 1{
			return "0"+str
		}
		return str
	}
	time := time.Now()
	date:=fmt.Sprintf("%s/%s/%s %s:%s",
	conv(int(time.Day())),
	conv(int(time.Month())),
	conv(int(time.Year())),
	conv(int(time.Hour())),
	conv(int(time.Minute())))
	self.top_helper.Set_string(date)
	self.elm_count+=1
	temp:=self.top_helper.cursor
	self.top_helper.cursor=self.init_index
	self.top_helper.Set_int(self.elm_count)
	self.top_helper.cursor=temp
}
type Instruction struct{
	Member_type JournalingMember
	Attrs []string
}
func New_inst(ins JournalingMember, params []string)Instruction{
	return Instruction{
		Member_type: ins,
		Attrs:       params,
	}
}
func (self *JournalingManager) Instructions_iter()JournalingIter{
	return JournalingIter{
		elm_count: self.elm_count,
		helper:    New_DataExtractionHelper(self.io_service,self.init_index+4),
	}
}
func New_DataExtractionHelper(io *datamanagment.IOService,init_index int32)DataExtractionHelper{
	return DataExtractionHelper{
		io_service:     io,
		cursor:         init_index,
		int_extractor:  types.Integer{
			Super_service: io,
			Index:         init_index,
			Size:          4,
		},
		char_extractor: types.Character{
			Super_service: io,
			Index:         init_index,
			Size:          1,
		},
	}
}
type DataExtractionHelper struct{
	io_service *datamanagment.IOService
	cursor int32
	int_extractor types.Integer
	char_extractor types.Character
}
const DATE_LENGTH = 16
func (self *DataExtractionHelper) Extract_string(length int32)string{
	data:=self.io_service.Read(length,self.cursor)
	self.cursor+=length
	return string(*data)
}
func (self *DataExtractionHelper) Extract_int()int32{
	self.int_extractor.Index = self.cursor
	int_val:=self.int_extractor.Get()
	self.cursor += 4
	return int_val
}
func (self *DataExtractionHelper) Extract_char()string{
	self.char_extractor.Index = self.cursor
	char_val:=self.char_extractor.Get()
	self.cursor += 1
	return char_val
}
func (self *DataExtractionHelper) Advance_items(no_items int32){
	for i := 0; i < int(no_items); i++ {
		ch:=self.Extract_char()
		
		switch ch {
		case string(Login):
			user_length:=self.Extract_int()
			pass_length:=self.Extract_int()
			self.cursor+=user_length+pass_length+DATE_LENGTH
		case string(Unlog):
			self.cursor+=DATE_LENGTH
		case string(Make_group):
			name_length:=self.Extract_int()
			self.cursor+=name_length+DATE_LENGTH
		case string(Remove_group):
			name_length:=self.Extract_int()
			self.cursor+=name_length+DATE_LENGTH
		case string(Make_user):
			user_length:=self.Extract_int()
			pass_length:=self.Extract_int()
			grp_length:=self.Extract_int()
			self.cursor+=user_length+pass_length+grp_length+DATE_LENGTH
		case string(Remove_user):
			user_length:=self.Extract_int()
			self.cursor+=user_length+DATE_LENGTH
		case string(Make_file):
			// recursive length = 1
			path_length:=self.Extract_int()
			content_length:=self.Extract_int()
			self.cursor+=1+path_length+content_length+DATE_LENGTH
		case string(Remove):
			path_length:=self.Extract_int()
			self.cursor+=path_length+DATE_LENGTH
		case string(Edit_file):
			path_length:=self.Extract_int()
			content_length:=self.Extract_int()
			self.cursor+=path_length+content_length+DATE_LENGTH
		case string(Rename_inode):
			path_length:=self.Extract_int()
			name_length:=self.Extract_int()
			self.cursor+=path_length+name_length+DATE_LENGTH
		case string(Make_dir):
			// recursive length = 1
			path_length:=self.Extract_int()
			self.cursor+=1+path_length+DATE_LENGTH
		case string(Copy):
			path_length:=self.Extract_int()
			path_dest_length:=self.Extract_int()
			self.cursor+=path_length+path_dest_length+DATE_LENGTH
		case string(Move):
			path_length:=self.Extract_int()
			path_dest_length:=self.Extract_int()
			self.cursor+=path_length+path_dest_length+DATE_LENGTH
		case string(Chown):
			// recursive length = 1
			path_length:=self.Extract_int()
			user_length:=self.Extract_int()
			self.cursor+=1+path_length+user_length+DATE_LENGTH
		case string(Chgrp):
			user_length:=self.Extract_int()
			grp_length:=self.Extract_int()
			self.cursor+=user_length+grp_length+DATE_LENGTH
		case string(Chmod):
			// recursive length = 1
			path_length:=self.Extract_int()
			ugo_length:=self.Extract_int()
			self.cursor+=1+path_length+ugo_length+DATE_LENGTH
		default: panic("Corrupted")
		}
	}
}
func (self *DataExtractionHelper) Set_string(str string){
	byts:=str[:]
	self.io_service.Write([]byte(byts),self.cursor)
	self.cursor+=int32(len(str))
}
func (self *DataExtractionHelper) Set_int(i int32){
	self.int_extractor.Index = self.cursor
	self.int_extractor.Set(i)
	self.cursor += 4
}
func (self *DataExtractionHelper) Set_char(char string){
	self.char_extractor.Index = self.cursor
	self.char_extractor.Set(char)
	self.cursor += 1
}
type JournalingIter struct{
	elm_count int32
	helper DataExtractionHelper
}
func (self *JournalingIter) Has_next()bool{
	return self.elm_count != 0
	}
func (self *JournalingIter) Next()Instruction{
	if self.elm_count == 0 {panic("No elements left")}
	member:=Login
	attrs:=make([]string, 0,4)
	get:=func(i int32)string{
		return self.helper.Extract_string(i)
	}
	ch:=self.helper.Extract_char()
	switch ch {
	case string(Login):
		member = Login
		user_length:=self.helper.Extract_int()
		pass_length:=self.helper.Extract_int()
		attrs = append(attrs, get(user_length),get(pass_length))
	case string(Unlog):
		member = Unlog
	case string(Make_group):
		member = Make_group
		name_length:=self.helper.Extract_int()
		attrs = append(attrs, get(name_length))
	case string(Remove_group):
		member = Remove_group
		name_length:=self.helper.Extract_int()
		attrs = append(attrs, get(name_length))
	case string(Make_user):
		member = Make_user
		user_length:=self.helper.Extract_int()
		pass_length:=self.helper.Extract_int()
		grp_length:=self.helper.Extract_int()
		attrs = append(attrs, get(user_length),get(pass_length),get(grp_length))
	case string(Remove_user):
		member = Remove_user
		user_length:=self.helper.Extract_int()
		attrs = append(attrs, get(user_length))
	case string(Make_file):
		member = Make_file
		// recursive length = 1
		path_length:=self.helper.Extract_int()
		content_length:=self.helper.Extract_int()
		attrs = append(attrs, self.helper.Extract_char(),get(path_length),get(content_length))
	case string(Remove):
		member = Remove
		path_length:=self.helper.Extract_int()
		attrs = append(attrs, get(path_length))
	case string(Edit_file):
		member = Edit_file
		path_length:=self.helper.Extract_int()
		content_length:=self.helper.Extract_int()
		attrs = append(attrs, get(path_length),get(content_length))
	case string(Rename_inode):
		member = Rename_inode
		path_length:=self.helper.Extract_int()
		name_length:=self.helper.Extract_int()
		attrs = append(attrs, get(path_length),get(name_length))
	case string(Make_dir):
		member = Make_dir
		// recursive length = 1
		path_length:=self.helper.Extract_int()
		attrs = append(attrs, self.helper.Extract_char(),get(path_length))
	case string(Copy):
		member = Copy
		path_length:=self.helper.Extract_int()
		path_dest_length:=self.helper.Extract_int()
		attrs = append(attrs,get(path_length),get(path_dest_length))
	case string(Move):
		member = Move
		path_length:=self.helper.Extract_int()
		path_dest_length:=self.helper.Extract_int()
		attrs = append(attrs,get(path_length),get(path_dest_length))
	case string(Chown):
		member = Chown
		// recursive length = 1
		path_length:=self.helper.Extract_int()
		user_length:=self.helper.Extract_int()
		attrs = append(attrs, self.helper.Extract_char(),get(path_length),get(user_length))
	case string(Chgrp):
		member = Chgrp
		user_length:=self.helper.Extract_int()
		grp_length:=self.helper.Extract_int()
		attrs = append(attrs,get(user_length),get(grp_length))
	case string(Chmod):
		member = Chmod
		// recursive length = 1
		path_length:=self.helper.Extract_int()
		ugo_length:=self.helper.Extract_int()
		attrs = append(attrs, self.helper.Extract_char(), get(path_length),get(ugo_length))
	default: panic(fmt.Sprintf("Corrupted \"%s\"",ch))
	}
	attrs = append(attrs, get(DATE_LENGTH))
	self.elm_count -=1
	return Instruction{
		Member_type: member,
		Attrs:       attrs,
	}
}