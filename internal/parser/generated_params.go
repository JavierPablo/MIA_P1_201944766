package parser

import (
	"fmt"
	"project/internal/utiles"
	"strconv"
)

type MkdiskParam struct {
	Size int32
	Fit  utiles.FitCriteria
	Unit utiles.SizeUnit
}

func (self *Task) Get_MkdiskParam() (MkdiskParam, error) {

	var res0 int32
	var err0 error
	param0, err0 := self.get_value("size")
	if err0 != nil {
		return MkdiskParam{}, err0
	}

	tempres0, err0 := strconv.Atoi(param0)
	if err0 != nil {
		return MkdiskParam{}, err0
	}
	if tempres0 < 0 {
		return MkdiskParam{}, fmt.Errorf("value is negative")
	}
	res0 = int32(tempres0)

	var res1 utiles.FitCriteria
	var err1 error
	param1, err1 := self.get_value("fit")
	if err1 != nil {
		res1 = utiles.First
	} else {

		res1, err1 = utiles.Translate_fit(param1)
		if err1 != nil {
			return MkdiskParam{}, err1
		}

	}

	var res2 utiles.SizeUnit
	var err2 error
	param2, err2 := self.get_value("unit")
	if err2 != nil {
		res2 = utiles.Mb
	} else {

		res2, err2 = utiles.Translate_size_unit(param2)
		if err2 != nil {
			return MkdiskParam{}, err2
		}

	}

	return MkdiskParam{
		Size: res0,
		Fit:  res1,
		Unit: res2,
	}, nil

}

type RmdiskParam struct {
	Driveletter string
}

func (self *Task) Get_RmdiskParam() (RmdiskParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("driveletter")
	if err0 != nil {
		return RmdiskParam{}, err0
	}

	res0 = param0

	return RmdiskParam{
		Driveletter: res0,
	}, nil

}

type FdiskParam struct {
	Size        int32
	Driveletter string
	Name        string
	Unit        utiles.SizeUnit
	Type        utiles.PartitionType
	Fit         utiles.FitCriteria
	Delete      bool
	Add         int32
}

func (self *Task) Get_FdiskParam() (FdiskParam, error) {

	var res0 int32
	var err0 error
	param0, err0 := self.get_value("size")
	if err0 != nil {
		res0 = 0
	} else {

		tempres0, err0 := strconv.Atoi(param0)
		if err0 != nil {
			return FdiskParam{}, err0
		}
		if tempres0 < 0 {
			return FdiskParam{}, fmt.Errorf("value is negative")
		}
		res0 = int32(tempres0)

	}

	var res1 string
	var err1 error
	param1, err1 := self.get_value("driveletter")
	if err1 != nil {
		return FdiskParam{}, err1
	}

	res1 = param1

	var res2 string
	var err2 error
	param2, err2 := self.get_value("name")
	if err2 != nil {
		return FdiskParam{}, err2
	}

	res2 = param2

	var res3 utiles.SizeUnit
	var err3 error
	param3, err3 := self.get_value("unit")
	if err3 != nil {
		res3 = utiles.Kb
	} else {

		res3, err3 = utiles.Translate_size_unit(param3)
		if err3 != nil {
			return FdiskParam{}, err3
		}

	}

	var res4 utiles.PartitionType
	var err4 error
	param4, err4 := self.get_value("type")
	if err4 != nil {
		res4 = utiles.Primary
	} else {

		res4, err4 = utiles.Translate_partition_type(param4)
		if err4 != nil {
			return FdiskParam{}, err4
		}

	}

	var res5 utiles.FitCriteria
	var err5 error
	param5, err5 := self.get_value("fit")
	if err5 != nil {
		res5 = utiles.Worst
	} else {

		res5, err5 = utiles.Translate_fit(param5)
		if err5 != nil {
			return FdiskParam{}, err5
		}

	}

	var res6 bool
	var err6 error
	param6, err6 := self.get_value("delete")
	if err6 != nil {
		res6 = false
	} else {

		if param6 == "" {
		}
		res6 = true

	}

	var res7 int32
	var err7 error
	param7, err7 := self.get_value("add")
	if err7 != nil {
		res7 = 0
	} else {

		tempres7, err7 := strconv.Atoi(param7)
		res7 = int32(tempres7)
		if err7 != nil {
			return FdiskParam{}, err7
		}

	}

	return FdiskParam{
		Size:        res0,
		Driveletter: res1,
		Name:        res2,
		Unit:        res3,
		Type:        res4,
		Fit:         res5,
		Delete:      res6,
		Add:         res7,
	}, nil

}

type PrintParam struct {
	Val string
}

func (self *Task) Get_PrintParam() (PrintParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("val")
	if err0 != nil {
		return PrintParam{}, err0
	}

	res0 = param0

	return PrintParam{
		Val: res0,
	}, nil

}

type MountParam struct {
	Driveletter string
	Name        string
}

func (self *Task) Get_MountParam() (MountParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("driveletter")
	if err0 != nil {
		return MountParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("name")
	if err1 != nil {
		return MountParam{}, err1
	}

	res1 = param1

	return MountParam{
		Driveletter: res0,
		Name:        res1,
	}, nil

}

type UnmountParam struct {
	Id string
}

func (self *Task) Get_UnmountParam() (UnmountParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("id")
	if err0 != nil {
		return UnmountParam{}, err0
	}

	res0 = param0

	return UnmountParam{
		Id: res0,
	}, nil

}

type MkfsParam struct {
	Id   string
	Type bool
	Fs   utiles.Format
}

func (self *Task) Get_MkfsParam() (MkfsParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("id")
	if err0 != nil {
		return MkfsParam{}, err0
	}

	res0 = param0

	var res1 bool
	var err1 error
	param1, err1 := self.get_value("type")
	if err1 != nil {
		res1 = false
	} else {

		if param1 == "" {
		}
		res1 = true

	}

	var res2 utiles.Format
	var err2 error
	param2, err2 := self.get_value("fs")
	if err2 != nil {
		res2 = utiles.Ext2
	} else {

		res2, err2 = utiles.Translate_format_type(param2)
		if err2 != nil {
			return MkfsParam{}, err2
		}

	}

	return MkfsParam{
		Id:   res0,
		Type: res1,
		Fs:   res2,
	}, nil

}

type LoginParam struct {
	User string
	Pass string
	Id   string
}

func (self *Task) Get_LoginParam() (LoginParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("user")
	if err0 != nil {
		return LoginParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("pass")
	if err1 != nil {
		return LoginParam{}, err1
	}

	res1 = param1

	var res2 string
	var err2 error
	param2, err2 := self.get_value("id")
	if err2 != nil {
		return LoginParam{}, err2
	}

	res2 = param2

	return LoginParam{
		User: res0,
		Pass: res1,
		Id:   res2,
	}, nil

}

type LogoutParam struct {
}

func (self *Task) Get_LogoutParam() (LogoutParam, error) {

	return LogoutParam{}, nil

}

type MkgrpParam struct {
	Name string
}

func (self *Task) Get_MkgrpParam() (MkgrpParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("name")
	if err0 != nil {
		return MkgrpParam{}, err0
	}

	res0 = param0

	return MkgrpParam{
		Name: res0,
	}, nil

}

type RmgrpParam struct {
	Name string
}

func (self *Task) Get_RmgrpParam() (RmgrpParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("name")
	if err0 != nil {
		return RmgrpParam{}, err0
	}

	res0 = param0

	return RmgrpParam{
		Name: res0,
	}, nil

}

type MkusrParam struct {
	User string
	Pass string
	Grp  string
}

func (self *Task) Get_MkusrParam() (MkusrParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("user")
	if err0 != nil {
		return MkusrParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("pass")
	if err1 != nil {
		return MkusrParam{}, err1
	}

	res1 = param1

	var res2 string
	var err2 error
	param2, err2 := self.get_value("grp")
	if err2 != nil {
		return MkusrParam{}, err2
	}

	res2 = param2

	return MkusrParam{
		User: res0,
		Pass: res1,
		Grp:  res2,
	}, nil

}

type RmusrParam struct {
	User string
}

func (self *Task) Get_RmusrParam() (RmusrParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("user")
	if err0 != nil {
		return RmusrParam{}, err0
	}

	res0 = param0

	return RmusrParam{
		User: res0,
	}, nil

}

type MkfileParam struct {
	Path      string
	R         bool
	Size      int32
	Cont      string
	Fixedcont string
}

func (self *Task) Get_MkfileParam() (MkfileParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return MkfileParam{}, err0
	}

	res0 = param0

	var res1 bool
	var err1 error
	param1, err1 := self.get_value("r")
	if err1 != nil {
		res1 = false
	} else {

		if param1 == "" {
		}
		res1 = true

	}

	var res2 int32
	var err2 error
	param2, err2 := self.get_value("size")
	if err2 != nil {
		res2 = 0
	} else {

		tempres2, err2 := strconv.Atoi(param2)
		if err2 != nil {
			return MkfileParam{}, err2
		}
		if tempres2 < 0 {
			return MkfileParam{}, fmt.Errorf("value is negative")
		}
		res2 = int32(tempres2)

	}

	var res3 string
	var err3 error
	param3, err3 := self.get_value("cont")
	if err3 != nil {
		res3 = ""
	} else {

		res3 = param3

	}

	var res4 string
	var err4 error
	param4, err4 := self.get_value("fixedcont")
	if err4 != nil {
		res4 = ""
	} else {

		res4 = param4

	}

	return MkfileParam{
		Path:      res0,
		R:         res1,
		Size:      res2,
		Cont:      res3,
		Fixedcont: res4,
	}, nil

}

type RemoveParam struct {
	Path string
}

func (self *Task) Get_RemoveParam() (RemoveParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return RemoveParam{}, err0
	}

	res0 = param0

	return RemoveParam{
		Path: res0,
	}, nil

}

type EditParam struct {
	Path      string
	Cont      string
	Fixedcont string
}

func (self *Task) Get_EditParam() (EditParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return EditParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("cont")
	if err1 != nil {
		return EditParam{}, err1
	}

	res1 = param1

	var res2 string
	var err2 error
	param2, err2 := self.get_value("fixedcont")
	if err2 != nil {
		res2 = ""
	} else {

		res2 = param2

	}

	return EditParam{
		Path:      res0,
		Cont:      res1,
		Fixedcont: res2,
	}, nil

}

type RenameParam struct {
	Path string
	Name string
}

func (self *Task) Get_RenameParam() (RenameParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return RenameParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("name")
	if err1 != nil {
		return RenameParam{}, err1
	}

	res1 = param1

	return RenameParam{
		Path: res0,
		Name: res1,
	}, nil

}

type MkdirParam struct {
	Path string
	R    bool
}

func (self *Task) Get_MkdirParam() (MkdirParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return MkdirParam{}, err0
	}

	res0 = param0

	var res1 bool
	var err1 error
	param1, err1 := self.get_value("r")
	if err1 != nil {
		res1 = false
	} else {

		if param1 == "" {
		}
		res1 = true

	}

	return MkdirParam{
		Path: res0,
		R:    res1,
	}, nil

}

type CopyParam struct {
	Path    string
	Destino string
}

func (self *Task) Get_CopyParam() (CopyParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return CopyParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("destino")
	if err1 != nil {
		return CopyParam{}, err1
	}

	res1 = param1

	return CopyParam{
		Path:    res0,
		Destino: res1,
	}, nil

}

type MoveParam struct {
	Path    string
	Destino string
}

func (self *Task) Get_MoveParam() (MoveParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return MoveParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("destino")
	if err1 != nil {
		return MoveParam{}, err1
	}

	res1 = param1

	return MoveParam{
		Path:    res0,
		Destino: res1,
	}, nil

}

type FindParam struct {
	Path string
	Name string
}

func (self *Task) Get_FindParam() (FindParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return FindParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("name")
	if err1 != nil {
		return FindParam{}, err1
	}

	res1 = param1

	return FindParam{
		Path: res0,
		Name: res1,
	}, nil

}

type ChownParam struct {
	Path string
	User string
	R    bool
}

func (self *Task) Get_ChownParam() (ChownParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return ChownParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("user")
	if err1 != nil {
		return ChownParam{}, err1
	}

	res1 = param1

	var res2 bool
	var err2 error
	param2, err2 := self.get_value("r")
	if err2 != nil {
		res2 = false
	} else {

		if param2 == "" {
		}
		res2 = true

	}

	return ChownParam{
		Path: res0,
		User: res1,
		R:    res2,
	}, nil

}

type ChgrpParam struct {
	User string
	Grp  string
}

func (self *Task) Get_ChgrpParam() (ChgrpParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("user")
	if err0 != nil {
		return ChgrpParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("grp")
	if err1 != nil {
		return ChgrpParam{}, err1
	}

	res1 = param1

	return ChgrpParam{
		User: res0,
		Grp:  res1,
	}, nil

}

type ChmodParam struct {
	Path string
	Ugo  string
	R    bool
}

func (self *Task) Get_ChmodParam() (ChmodParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return ChmodParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("ugo")
	if err1 != nil {
		return ChmodParam{}, err1
	}

	res1 = param1

	var res2 bool
	var err2 error
	param2, err2 := self.get_value("r")
	if err2 != nil {
		res2 = false
	} else {

		if param2 == "" {
		}
		res2 = true

	}

	return ChmodParam{
		Path: res0,
		Ugo:  res1,
		R:    res2,
	}, nil

}

type PauseParam struct {
}

func (self *Task) Get_PauseParam() (PauseParam, error) {

	return PauseParam{}, nil

}

type RecoveryParam struct {
	Id string
}

func (self *Task) Get_RecoveryParam() (RecoveryParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("id")
	if err0 != nil {
		return RecoveryParam{}, err0
	}

	res0 = param0

	return RecoveryParam{
		Id: res0,
	}, nil

}

type LossParam struct {
	Id string
}

func (self *Task) Get_LossParam() (LossParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("id")
	if err0 != nil {
		return LossParam{}, err0
	}

	res0 = param0

	return LossParam{
		Id: res0,
	}, nil

}

type ExecuteParam struct {
	Path string
}

func (self *Task) Get_ExecuteParam() (ExecuteParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("path")
	if err0 != nil {
		return ExecuteParam{}, err0
	}

	res0 = param0

	return ExecuteParam{
		Path: res0,
	}, nil

}

type RepParam struct {
	Name string
	Path string
	Id   string
	Ruta string
}

func (self *Task) Get_RepParam() (RepParam, error) {

	var res0 string
	var err0 error
	param0, err0 := self.get_value("name")
	if err0 != nil {
		return RepParam{}, err0
	}

	res0 = param0

	var res1 string
	var err1 error
	param1, err1 := self.get_value("path")
	if err1 != nil {
		return RepParam{}, err1
	}

	res1 = param1

	var res2 string
	var err2 error
	param2, err2 := self.get_value("id")
	if err2 != nil {
		return RepParam{}, err2
	}

	res2 = param2

	var res3 string
	var err3 error
	param3, err3 := self.get_value("ruta")
	if err3 != nil {
		res3 = ""
	} else {

		res3 = param3

	}

	return RepParam{
		Name: res0,
		Path: res1,
		Id:   res2,
		Ruta: res3,
	}, nil

}
