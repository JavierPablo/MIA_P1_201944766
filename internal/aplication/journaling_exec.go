package aplication

// func (self *Aplication) Execute_for_journaling{
// 	switch self.helper.Extract_char() {
// 	case string(Login):
// 		user_length:=self.helper.Extract_int()
// 		pass_length:=self.helper.Extract_int()
// 		attrs = append(attrs, get(user_length),get(pass_length))
// 	case string(Unlog):
// 	case string(Make_group):
// 		name_length:=self.helper.Extract_int()
// 		attrs = append(attrs, get(name_length))
// 	case string(Remove_group):
// 		name_length:=self.helper.Extract_int()
// 		attrs = append(attrs, get(name_length))
// 	case string(Make_user):
// 		user_length:=self.helper.Extract_int()
// 		pass_length:=self.helper.Extract_int()
// 		grp_length:=self.helper.Extract_int()
// 		attrs = append(attrs, get(user_length),get(pass_length),get(grp_length))
// 	case string(Remove_user):
// 		user_length:=self.helper.Extract_int()
// 		attrs = append(attrs, get(user_length))
// 	case string(Make_file):
// 		// recursive length = 1
// 		path_length:=self.helper.Extract_int()
// 		content_length:=self.helper.Extract_int()
// 		attrs = append(attrs, self.helper.Extract_char(),get(path_length),get(content_length))
// 	case string(Remove):
// 		path_length:=self.helper.Extract_int()
// 		attrs = append(attrs, get(path_length))
// 	case string(Edit_file):
// 		path_length:=self.helper.Extract_int()
// 		content_length:=self.helper.Extract_int()
// 		attrs = append(attrs, get(path_length),get(content_length))
// 	case string(Rename_inode):
// 		path_length:=self.helper.Extract_int()
// 		name_length:=self.helper.Extract_int()
// 		attrs = append(attrs, get(path_length),get(name_length))
// 	case string(Make_dir):
// 		// recursive length = 1
// 		path_length:=self.helper.Extract_int()
// 		attrs = append(attrs, self.helper.Extract_char(),get(path_length))
// 	case string(Copy):
// 		path_length:=self.helper.Extract_int()
// 		path_dest_length:=self.helper.Extract_int()
// 		attrs = append(attrs,get(path_length),get(path_dest_length))
// 	case string(Move):
// 		path_length:=self.helper.Extract_int()
// 		path_dest_length:=self.helper.Extract_int()
// 		attrs = append(attrs,get(path_length),get(path_dest_length))
// 	case string(Chown):
// 		// recursive length = 1
// 		path_length:=self.helper.Extract_int()
// 		user_length:=self.helper.Extract_int()
// 		attrs = append(attrs, self.helper.Extract_char(),get(path_length),get(user_length))
// 	case string(Chgrp):
// 		user_length:=self.helper.Extract_int()
// 		grp_length:=self.helper.Extract_int()
// 		attrs = append(attrs,get(user_length),get(grp_length))
// 	case string(Chmod):
// 		// recursive length = 1
// 		path_length:=self.helper.Extract_int()
// 		ugo_length:=self.helper.Extract_int()
// 		attrs = append(attrs, self.helper.Extract_char(), get(path_length),get(ugo_length))
// 	default: panic("Corrupted")
// 	}
// }