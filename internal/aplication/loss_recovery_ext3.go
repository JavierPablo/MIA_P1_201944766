package aplication

import (
	"fmt"
	"project/internal/formats"
	"project/internal/parser"
	"project/internal/types"
	"project/internal/utiles"
)

func (self *Aplication) Recovery_ext3(id string) ([]*parser.Task,error){
	for i := 0; i < len(self.mounted_partitions); i++ {
		// fmt.Printf("%s == %s",(self.mounted_partitions[i]).id,id)
		if (self.mounted_partitions[i]).id == id{
			mounted := &self.mounted_partitions[i]
			io_service := mounted.io
			var super_block_index int32
			var fit utiles.FitCriteria
			var err error
			switch self.active_partition.part_type {
			case utiles.Logic:
				ebr := types.CreateExtendedBootRecord(io_service, self.active_partition.index)
				super_block_index = ebr.Part_start().Get()
				fit,err = utiles.Translate_fit(ebr.Part_fit().Get())
				if err!=nil{return nil,err}

			case utiles.Primary:
				partition := types.CreatePartition(io_service, self.active_partition.index)
				super_block_index = partition.Part_start().Get()
				fit,err = utiles.Translate_fit(partition.Part_fit().Get())
				if err!=nil{return nil,err}
			}
			format:=formats.Get_FormatEXT3_for_heal(io_service,super_block_index,fit)
			if format == nil{
				return nil,fmt.Errorf("can not apply recovery command to non ext3 formated partition")
			}			
			tasks:=make([]*parser.Task,0,10)
			iter:=format.Get_journaling().Instructions_iter()
			for iter.Has_next(){
				inst:=iter.Next()
				var task parser.Task
				switch inst.Member_type {
				case formats.Login:
					task=parser.Task{Command: "login",Flags:[]*parser.Flag{
							&parser.Flag{Key:"user",Value: inst.Attrs[0]},
							&parser.Flag{Key:"pass",Value: inst.Attrs[1]},
							&parser.Flag{Key:"id",Value: id},
					}}
				case formats.Unlog:
					task=parser.Task{Command: "logout",Flags:[]*parser.Flag{
					}}
				case formats.Make_group:
					task=parser.Task{Command: "mkgrp",Flags:[]*parser.Flag{
						&parser.Flag{Key:"name",Value: inst.Attrs[0]},
					}}
				case formats.Remove_group:
					task=parser.Task{Command: "rmgrp",Flags:[]*parser.Flag{
						&parser.Flag{Key:"name",Value: inst.Attrs[0]},
					}}
				case formats.Make_user:
					task=parser.Task{Command: "mkusr",Flags:[]*parser.Flag{
						&parser.Flag{Key:"user",Value: inst.Attrs[0]},
						&parser.Flag{Key:"pass",Value: inst.Attrs[1]},
						&parser.Flag{Key:"grp",Value: inst.Attrs[2]},
					}}
				case formats.Remove_user:
					task=parser.Task{Command: "rmusr",Flags:[]*parser.Flag{
						&parser.Flag{Key:"user",Value: inst.Attrs[0]},
					}}
				case formats.Make_file:
					if inst.Attrs[0] == "Y"{
						task=parser.Task{Command: "mkfile",Flags:[]*parser.Flag{
							&parser.Flag{Key:"r",Value: ""},
							&parser.Flag{Key:"path",Value: inst.Attrs[1]},
							&parser.Flag{Key:"fixedcont",Value: inst.Attrs[2]},
						}}
					}else if inst.Attrs[0] == "N"{
						task=parser.Task{Command: "mkfile",Flags:[]*parser.Flag{
							&parser.Flag{Key:"path",Value: inst.Attrs[1]},
							&parser.Flag{Key:"fixedcont",Value: inst.Attrs[2]},
						}}
					}else{panic("unexpected value of rescursive")}
				case formats.Remove:
					task=parser.Task{Command: "remove",Flags:[]*parser.Flag{
						&parser.Flag{Key:"path",Value: inst.Attrs[0]},
					}}
				case formats.Edit_file:
					task=parser.Task{Command: "edit",Flags:[]*parser.Flag{
						&parser.Flag{Key:"path",Value: inst.Attrs[0]},
						&parser.Flag{Key:"cont",Value: ""},
						&parser.Flag{Key:"fixedcont",Value: inst.Attrs[1]},
					}}
				case formats.Rename_inode:
					task=parser.Task{Command: "rename",Flags:[]*parser.Flag{
						&parser.Flag{Key:"path",Value: inst.Attrs[0]},
						&parser.Flag{Key:"name",Value: inst.Attrs[1]},
					}}
				case formats.Make_dir:
					if inst.Attrs[0] == "Y"{
						task=parser.Task{Command: "mkdir",Flags:[]*parser.Flag{
							&parser.Flag{Key:"r",Value: ""},
							&parser.Flag{Key:"path",Value: inst.Attrs[1]},
						}}
					}else if inst.Attrs[0] == "N"{
						task=parser.Task{Command: "mkdir",Flags:[]*parser.Flag{
							&parser.Flag{Key:"path",Value: inst.Attrs[1]},
						}}
					}else{panic("unexpected value of rescursive")}
				case formats.Copy:
					task=parser.Task{Command: "copy",Flags:[]*parser.Flag{
						&parser.Flag{Key:"path",Value: inst.Attrs[0]},
						&parser.Flag{Key:"destino",Value: inst.Attrs[1]},
					}}
				case formats.Move:
					task=parser.Task{Command: "move",Flags:[]*parser.Flag{
						&parser.Flag{Key:"path",Value: inst.Attrs[0]},
						&parser.Flag{Key:"destino",Value: inst.Attrs[1]},
					}}
				case formats.Chown:
					if inst.Attrs[0] == "Y"{
						task=parser.Task{Command: "chown",Flags:[]*parser.Flag{
							&parser.Flag{Key:"r",Value: ""},
							&parser.Flag{Key:"path",Value: inst.Attrs[1]},
							&parser.Flag{Key:"user",Value: inst.Attrs[2]},
						}}
					}else if inst.Attrs[0] == "N"{
						task=parser.Task{Command: "chown",Flags:[]*parser.Flag{
							&parser.Flag{Key:"path",Value: inst.Attrs[1]},
							&parser.Flag{Key:"user",Value: inst.Attrs[2]},
						}}
					}else{panic("unexpected value of rescursive")}
				case formats.Chgrp:
					task=parser.Task{Command: "chgrp",Flags:[]*parser.Flag{
						&parser.Flag{Key:"user",Value: inst.Attrs[0]},
						&parser.Flag{Key:"grp",Value: inst.Attrs[1]},
					}}
				case formats.Chmod:
					if inst.Attrs[0] == "Y"{
						task=parser.Task{Command: "chmod",Flags:[]*parser.Flag{
							&parser.Flag{Key:"r",Value: ""},
							&parser.Flag{Key:"path",Value: inst.Attrs[1]},
							&parser.Flag{Key:"ugo",Value: inst.Attrs[2]},
						}}
					}else if inst.Attrs[0] == "N"{
						task=parser.Task{Command: "chmod",Flags:[]*parser.Flag{
							&parser.Flag{Key:"path",Value: inst.Attrs[1]},
							&parser.Flag{Key:"ugo",Value: inst.Attrs[2]},
						}}
					}else{panic("unexpected value of rescursive")}
				default: panic("Corrupted")
				}
				tasks = append(tasks, &task)

			}
			return tasks,nil
		}

	}
	return nil,fmt.Errorf("There is no mounted partition with that name")
}
func (self *Aplication) Loss_ext3(id string) error{
	for i := 0; i < len(self.mounted_partitions); i++ {
		// fmt.Printf("%s == %s",(self.mounted_partitions[i]).id,id)
		if (self.mounted_partitions[i]).id == id{
			mounted := &self.mounted_partitions[i]
			io_service := mounted.io
			var super_block_index int32
			var fit utiles.FitCriteria
			var err error
			switch self.mounted_partitions[i].part_type{
			case utiles.Logic:
				ebr := types.CreateExtendedBootRecord(io_service, self.active_partition.index)
				super_block_index = ebr.Part_start().Get()
				fit,err = utiles.Translate_fit(ebr.Part_fit().Get())
				if err!=nil{return err}
		
			case utiles.Primary:
				partition := types.CreatePartition(io_service, self.active_partition.index)
				super_block_index = partition.Part_start().Get()
				fit,err = utiles.Translate_fit(partition.Part_fit().Get())
				if err!=nil{return err}
			}
			mounted.has_session=false
			mounted.session = SessionManager{}
			format:=formats.Recover_Format(io_service,super_block_index,fit)
			if !format.Has_journaling(){
				return fmt.Errorf("can not apply loss command to ext2 formated partition")
			}
			sp_blck:=types.CreateSuperBlock(io_service,super_block_index)
			sp_blck.S_firts_ino().Set(-1)

			format.Block_bitmap.Clear()
			format.Inodes_bitmap.Clear()
			format.Block_section.Clear()
			format.Inodes_section.Clear()
			return nil
		}

	}
	return fmt.Errorf("There is no mounted partition with that name")
}
