package parser

import "strings"


type CatParam struct {
	Paths []string
}

func (self *Task) Get_CatParam() (CatParam, error) {
	paths:=make([]string, 0,len(self.Flags))
	for i := 0; i < len(self.Flags); i++ {
		value := self.Flags[i].Value
		if strings.HasPrefix(value,"\""){
			value =  value[1 : len(value)-1]
		}
		paths = append(paths, value)
	}
	return CatParam{
		Paths: paths,
	}, nil

}


