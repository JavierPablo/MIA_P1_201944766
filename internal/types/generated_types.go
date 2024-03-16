package types
import "project/internal/datamanagment"
import "fmt"

type Time struct{
        Super_service *datamanagment.IOService
        Index int32
        Size int32
    }
func CreateTime(Super_service *datamanagment.IOService,Index int32) Time{
        return Time{
            Super_service: Super_service,
            Index: Index,
            Size:24,
        }
    }
type TimeHolder struct{
        Hour int32
Minute int32
Second int32
Day int32
Month int32
Year int32

    }
func (self Time)Get() TimeHolder{
        return TimeHolder{
            Hour: self.Hour().Get(),
Minute: self.Minute().Get(),
Second: self.Second().Get(),
Day: self.Day().Get(),
Month: self.Month().Get(),
Year: self.Year().Get(),

        }
    }
func (self Time)Set(trgt TimeHolder) {
        self.Hour().Set(trgt.Hour)
self.Minute().Set(trgt.Minute)
self.Second().Set(trgt.Second)
self.Day().Set(trgt.Day)
self.Month().Set(trgt.Month)
self.Year().Set(trgt.Year)

    }
func (self Time) Hour() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:0 + self.Index,
                Size:4,
            }
        }
func (self Time) Minute() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:4 + self.Index,
                Size:4,
            }
        }
func (self Time) Second() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:8 + self.Index,
                Size:4,
            }
        }
func (self Time) Day() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:12 + self.Index,
                Size:4,
            }
        }
func (self Time) Month() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:16 + self.Index,
                Size:4,
            }
        }
func (self Time) Year() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:20 + self.Index,
                Size:4,
            }
        }
func (self Time)Dot_label()string {
        return fmt.Sprintf("<TABLE BGCOLOR=\" #45B3F1\" BORDER=\"1\"  COLOR=\"BLACK\"> <TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Hour </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Minute </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Second </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Day </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Month </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Year </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR> </TABLE>",self.Hour().Dot_label(),self.Minute().Dot_label(),self.Second().Dot_label(),self.Day().Dot_label(),self.Month().Dot_label(),self.Year().Dot_label())
    }
func CreateArrayCharacter4(Super_service *datamanagment.IOService,Index int32) ArrayCharacter4{
        return ArrayCharacter4{
            Super_service: Super_service,
            Index: Index,
            Size:4,
        }
    }
type ArrayCharacter4 struct{
        Super_service *datamanagment.IOService
	    Index int32
	    Size int32
    }
func (self ArrayCharacter4)Get() [4]string{
        var array [4]string
        proxyes := self.Spread()
        var i int32= 0
        for i < 4 {
            array[i] = proxyes[i].Get()
            i++
        }
        return array
    }
func (self ArrayCharacter4)Set(trgt [4]string) {
        proxyes := self.Spread()
    	var i int32= 0
	    for i < 4 {
            proxyes[i].Set(trgt[i])
            i++
	    }
    }
func (self ArrayCharacter4) Spread() [4]Character {
        var array [4]Character
	    var i int32= 0
	    for i < 4 {
            array[i] = Character{
                Super_service:self.Super_service,
                Index: self.Index + (i*1),
                Size:1,
            }
	    	i++
    	}
	return array
    }
func (self ArrayCharacter4) No(i int32) Character {
        return Character{
                Super_service:self.Super_service,
                Index: self.Index + (i*1),
                Size:1,
            }
    }
func (self ArrayCharacter4)Dot_label()string {
        return fmt.Sprintf("<TABLE BORDER=\"0\"> <TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 0 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 1 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 2 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 3 </TD></TR> <TR><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD></TR></TABLE>",self.No(0).Dot_label(),self.No(1).Dot_label(),self.No(2).Dot_label(),self.No(3).Dot_label())
    }
func CreateArrayCharacter3(Super_service *datamanagment.IOService,Index int32) ArrayCharacter3{
        return ArrayCharacter3{
            Super_service: Super_service,
            Index: Index,
            Size:3,
        }
    }
type ArrayCharacter3 struct{
        Super_service *datamanagment.IOService
	    Index int32
	    Size int32
    }
func (self ArrayCharacter3)Get() [3]string{
        var array [3]string
        proxyes := self.Spread()
        var i int32= 0
        for i < 3 {
            array[i] = proxyes[i].Get()
            i++
        }
        return array
    }
func (self ArrayCharacter3)Set(trgt [3]string) {
        proxyes := self.Spread()
    	var i int32= 0
	    for i < 3 {
            proxyes[i].Set(trgt[i])
            i++
	    }
    }
func (self ArrayCharacter3) Spread() [3]Character {
        var array [3]Character
	    var i int32= 0
	    for i < 3 {
            array[i] = Character{
                Super_service:self.Super_service,
                Index: self.Index + (i*1),
                Size:1,
            }
	    	i++
    	}
	return array
    }
func (self ArrayCharacter3) No(i int32) Character {
        return Character{
                Super_service:self.Super_service,
                Index: self.Index + (i*1),
                Size:1,
            }
    }
func (self ArrayCharacter3)Dot_label()string {
        return fmt.Sprintf("<TABLE BORDER=\"0\"> <TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 0 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 1 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 2 </TD></TR> <TR><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD></TR></TABLE>",self.No(0).Dot_label(),self.No(1).Dot_label(),self.No(2).Dot_label())
    }
func CreateArrayCharacter12(Super_service *datamanagment.IOService,Index int32) ArrayCharacter12{
        return ArrayCharacter12{
            Super_service: Super_service,
            Index: Index,
            Size:12,
        }
    }
type ArrayCharacter12 struct{
        Super_service *datamanagment.IOService
	    Index int32
	    Size int32
    }
func (self ArrayCharacter12)Get() [12]string{
        var array [12]string
        proxyes := self.Spread()
        var i int32= 0
        for i < 12 {
            array[i] = proxyes[i].Get()
            i++
        }
        return array
    }
func (self ArrayCharacter12)Set(trgt [12]string) {
        proxyes := self.Spread()
    	var i int32= 0
	    for i < 12 {
            proxyes[i].Set(trgt[i])
            i++
	    }
    }
func (self ArrayCharacter12) Spread() [12]Character {
        var array [12]Character
	    var i int32= 0
	    for i < 12 {
            array[i] = Character{
                Super_service:self.Super_service,
                Index: self.Index + (i*1),
                Size:1,
            }
	    	i++
    	}
	return array
    }
func (self ArrayCharacter12) No(i int32) Character {
        return Character{
                Super_service:self.Super_service,
                Index: self.Index + (i*1),
                Size:1,
            }
    }
func (self ArrayCharacter12)Dot_label()string {
        return fmt.Sprintf("<TABLE BORDER=\"0\"> <TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 0 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 1 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 2 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 3 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 4 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 5 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 6 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 7 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 8 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 9 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 10 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 11 </TD></TR> <TR><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD></TR></TABLE>",self.No(0).Dot_label(),self.No(1).Dot_label(),self.No(2).Dot_label(),self.No(3).Dot_label(),self.No(4).Dot_label(),self.No(5).Dot_label(),self.No(6).Dot_label(),self.No(7).Dot_label(),self.No(8).Dot_label(),self.No(9).Dot_label(),self.No(10).Dot_label(),self.No(11).Dot_label())
    }
func CreateArrayCharacter16(Super_service *datamanagment.IOService,Index int32) ArrayCharacter16{
        return ArrayCharacter16{
            Super_service: Super_service,
            Index: Index,
            Size:16,
        }
    }
type ArrayCharacter16 struct{
        Super_service *datamanagment.IOService
	    Index int32
	    Size int32
    }
func (self ArrayCharacter16)Get() [16]string{
        var array [16]string
        proxyes := self.Spread()
        var i int32= 0
        for i < 16 {
            array[i] = proxyes[i].Get()
            i++
        }
        return array
    }
func (self ArrayCharacter16)Set(trgt [16]string) {
        proxyes := self.Spread()
    	var i int32= 0
	    for i < 16 {
            proxyes[i].Set(trgt[i])
            i++
	    }
    }
func (self ArrayCharacter16) Spread() [16]Character {
        var array [16]Character
	    var i int32= 0
	    for i < 16 {
            array[i] = Character{
                Super_service:self.Super_service,
                Index: self.Index + (i*1),
                Size:1,
            }
	    	i++
    	}
	return array
    }
func (self ArrayCharacter16) No(i int32) Character {
        return Character{
                Super_service:self.Super_service,
                Index: self.Index + (i*1),
                Size:1,
            }
    }
func (self ArrayCharacter16)Dot_label()string {
        return fmt.Sprintf("<TABLE BORDER=\"0\"> <TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 0 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 1 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 2 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 3 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 4 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 5 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 6 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 7 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 8 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 9 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 10 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 11 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 12 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 13 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 14 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 15 </TD></TR> <TR><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD></TR></TABLE>",self.No(0).Dot_label(),self.No(1).Dot_label(),self.No(2).Dot_label(),self.No(3).Dot_label(),self.No(4).Dot_label(),self.No(5).Dot_label(),self.No(6).Dot_label(),self.No(7).Dot_label(),self.No(8).Dot_label(),self.No(9).Dot_label(),self.No(10).Dot_label(),self.No(11).Dot_label(),self.No(12).Dot_label(),self.No(13).Dot_label(),self.No(14).Dot_label(),self.No(15).Dot_label())
    }
func CreateArrayCharacter64(Super_service *datamanagment.IOService,Index int32) ArrayCharacter64{
        return ArrayCharacter64{
            Super_service: Super_service,
            Index: Index,
            Size:64,
        }
    }
type ArrayCharacter64 struct{
        Super_service *datamanagment.IOService
	    Index int32
	    Size int32
    }
func (self ArrayCharacter64)Get() [64]string{
        var array [64]string
        proxyes := self.Spread()
        var i int32= 0
        for i < 64 {
            array[i] = proxyes[i].Get()
            i++
        }
        return array
    }
func (self ArrayCharacter64)Set(trgt [64]string) {
        proxyes := self.Spread()
    	var i int32= 0
	    for i < 64 {
            proxyes[i].Set(trgt[i])
            i++
	    }
    }
func (self ArrayCharacter64) Spread() [64]Character {
        var array [64]Character
	    var i int32= 0
	    for i < 64 {
            array[i] = Character{
                Super_service:self.Super_service,
                Index: self.Index + (i*1),
                Size:1,
            }
	    	i++
    	}
	return array
    }
func (self ArrayCharacter64) No(i int32) Character {
        return Character{
                Super_service:self.Super_service,
                Index: self.Index + (i*1),
                Size:1,
            }
    }
func (self ArrayCharacter64)Dot_label()string {
        return fmt.Sprintf("<TABLE BORDER=\"0\"> <TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 0 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 1 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 2 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 3 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 4 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 5 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 6 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 7 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 8 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 9 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 10 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 11 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 12 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 13 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 14 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 15 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 16 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 17 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 18 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 19 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 20 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 21 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 22 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 23 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 24 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 25 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 26 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 27 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 28 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 29 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 30 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 31 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 32 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 33 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 34 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 35 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 36 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 37 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 38 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 39 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 40 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 41 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 42 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 43 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 44 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 45 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 46 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 47 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 48 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 49 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 50 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 51 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 52 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 53 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 54 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 55 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 56 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 57 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 58 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 59 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 60 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 61 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 62 </TD><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 63 </TD></TR> <TR><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD><TD BORDER=\"0\"> %s </TD></TR></TABLE>",self.No(0).Dot_label(),self.No(1).Dot_label(),self.No(2).Dot_label(),self.No(3).Dot_label(),self.No(4).Dot_label(),self.No(5).Dot_label(),self.No(6).Dot_label(),self.No(7).Dot_label(),self.No(8).Dot_label(),self.No(9).Dot_label(),self.No(10).Dot_label(),self.No(11).Dot_label(),self.No(12).Dot_label(),self.No(13).Dot_label(),self.No(14).Dot_label(),self.No(15).Dot_label(),self.No(16).Dot_label(),self.No(17).Dot_label(),self.No(18).Dot_label(),self.No(19).Dot_label(),self.No(20).Dot_label(),self.No(21).Dot_label(),self.No(22).Dot_label(),self.No(23).Dot_label(),self.No(24).Dot_label(),self.No(25).Dot_label(),self.No(26).Dot_label(),self.No(27).Dot_label(),self.No(28).Dot_label(),self.No(29).Dot_label(),self.No(30).Dot_label(),self.No(31).Dot_label(),self.No(32).Dot_label(),self.No(33).Dot_label(),self.No(34).Dot_label(),self.No(35).Dot_label(),self.No(36).Dot_label(),self.No(37).Dot_label(),self.No(38).Dot_label(),self.No(39).Dot_label(),self.No(40).Dot_label(),self.No(41).Dot_label(),self.No(42).Dot_label(),self.No(43).Dot_label(),self.No(44).Dot_label(),self.No(45).Dot_label(),self.No(46).Dot_label(),self.No(47).Dot_label(),self.No(48).Dot_label(),self.No(49).Dot_label(),self.No(50).Dot_label(),self.No(51).Dot_label(),self.No(52).Dot_label(),self.No(53).Dot_label(),self.No(54).Dot_label(),self.No(55).Dot_label(),self.No(56).Dot_label(),self.No(57).Dot_label(),self.No(58).Dot_label(),self.No(59).Dot_label(),self.No(60).Dot_label(),self.No(61).Dot_label(),self.No(62).Dot_label(),self.No(63).Dot_label())
    }
type Partition struct{
        Super_service *datamanagment.IOService
        Index int32
        Size int32
    }
func CreatePartition(Super_service *datamanagment.IOService,Index int32) Partition{
        return Partition{
            Super_service: Super_service,
            Index: Index,
            Size:35,
        }
    }
type PartitionHolder struct{
        Part_status string
Part_type string
Part_fit string
Part_start int32
Part_s int32
Part_name [16]string
Part_correlative int32
Part_id [4]string

    }
func (self Partition)Get() PartitionHolder{
        return PartitionHolder{
            Part_status: self.Part_status().Get(),
Part_type: self.Part_type().Get(),
Part_fit: self.Part_fit().Get(),
Part_start: self.Part_start().Get(),
Part_s: self.Part_s().Get(),
Part_name: self.Part_name().Get(),
Part_correlative: self.Part_correlative().Get(),
Part_id: self.Part_id().Get(),

        }
    }
func (self Partition)Set(trgt PartitionHolder) {
        self.Part_status().Set(trgt.Part_status)
self.Part_type().Set(trgt.Part_type)
self.Part_fit().Set(trgt.Part_fit)
self.Part_start().Set(trgt.Part_start)
self.Part_s().Set(trgt.Part_s)
self.Part_name().Set(trgt.Part_name)
self.Part_correlative().Set(trgt.Part_correlative)
self.Part_id().Set(trgt.Part_id)

    }
func (self Partition) Part_status() Character{
            return Character{
                Super_service:self.Super_service,
                Index:0 + self.Index,
                Size:1,
            }
        }
func (self Partition) Part_type() Character{
            return Character{
                Super_service:self.Super_service,
                Index:1 + self.Index,
                Size:1,
            }
        }
func (self Partition) Part_fit() Character{
            return Character{
                Super_service:self.Super_service,
                Index:2 + self.Index,
                Size:1,
            }
        }
func (self Partition) Part_start() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:3 + self.Index,
                Size:4,
            }
        }
func (self Partition) Part_s() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:7 + self.Index,
                Size:4,
            }
        }
func (self Partition) Part_name() ArrayCharacter16{
            return ArrayCharacter16{
                Super_service:self.Super_service,
                Index:11 + self.Index,
                Size:16,
            }
        }
func (self Partition) Part_correlative() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:27 + self.Index,
                Size:4,
            }
        }
func (self Partition) Part_id() ArrayCharacter4{
            return ArrayCharacter4{
                Super_service:self.Super_service,
                Index:31 + self.Index,
                Size:4,
            }
        }
func (self Partition)Dot_label()string {
        return fmt.Sprintf("<TABLE BGCOLOR=\" #45B3F1\" BORDER=\"1\"  COLOR=\"BLACK\"> <TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_status </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_type </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_fit </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_start </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_s </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_name </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_correlative </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_id </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR> </TABLE>",self.Part_status().Dot_label(),self.Part_type().Dot_label(),self.Part_fit().Dot_label(),self.Part_start().Dot_label(),self.Part_s().Dot_label(),self.Part_name().Dot_label(),self.Part_correlative().Dot_label(),self.Part_id().Dot_label())
    }
func CreateArrayPartition4(Super_service *datamanagment.IOService,Index int32) ArrayPartition4{
        return ArrayPartition4{
            Super_service: Super_service,
            Index: Index,
            Size:140,
        }
    }
type ArrayPartition4 struct{
        Super_service *datamanagment.IOService
	    Index int32
	    Size int32
    }
func (self ArrayPartition4)Get() [4]PartitionHolder{
        var array [4]PartitionHolder
        proxyes := self.Spread()
        var i int32= 0
        for i < 4 {
            array[i] = proxyes[i].Get()
            i++
        }
        return array
    }
func (self ArrayPartition4)Set(trgt [4]PartitionHolder) {
        proxyes := self.Spread()
    	var i int32= 0
	    for i < 4 {
            proxyes[i].Set(trgt[i])
            i++
	    }
    }
func (self ArrayPartition4) Spread() [4]Partition {
        var array [4]Partition
	    var i int32= 0
	    for i < 4 {
            array[i] = Partition{
                Super_service:self.Super_service,
                Index: self.Index + (i*35),
                Size:35,
            }
	    	i++
    	}
	return array
    }
func (self ArrayPartition4) No(i int32) Partition {
        return Partition{
                Super_service:self.Super_service,
                Index: self.Index + (i*35),
                Size:35,
            }
    }
func (self ArrayPartition4)Dot_label()string {
                return fmt.Sprintf("<TABLE BORDER=\"0\"><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 0 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 1 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 2 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 3 </TD><TD BORDER=\"0\"> %s </TD></TR></TABLE>",self.No(0).Dot_label(),self.No(1).Dot_label(),self.No(2).Dot_label(),self.No(3).Dot_label())
            }
type MasterBootRecord struct{
        Super_service *datamanagment.IOService
        Index int32
        Size int32
    }
func CreateMasterBootRecord(Super_service *datamanagment.IOService,Index int32) MasterBootRecord{
        return MasterBootRecord{
            Super_service: Super_service,
            Index: Index,
            Size:173,
        }
    }
type MasterBootRecordHolder struct{
        Mbr_tamano int32
Mbr_fecha_creacion TimeHolder
Mbr_dsk_signature int32
Dsk_fit string
Mbr_partitions [4]PartitionHolder

    }
func (self MasterBootRecord)Get() MasterBootRecordHolder{
        return MasterBootRecordHolder{
            Mbr_tamano: self.Mbr_tamano().Get(),
Mbr_fecha_creacion: self.Mbr_fecha_creacion().Get(),
Mbr_dsk_signature: self.Mbr_dsk_signature().Get(),
Dsk_fit: self.Dsk_fit().Get(),
Mbr_partitions: self.Mbr_partitions().Get(),

        }
    }
func (self MasterBootRecord)Set(trgt MasterBootRecordHolder) {
        self.Mbr_tamano().Set(trgt.Mbr_tamano)
self.Mbr_fecha_creacion().Set(trgt.Mbr_fecha_creacion)
self.Mbr_dsk_signature().Set(trgt.Mbr_dsk_signature)
self.Dsk_fit().Set(trgt.Dsk_fit)
self.Mbr_partitions().Set(trgt.Mbr_partitions)

    }
func (self MasterBootRecord) Mbr_tamano() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:0 + self.Index,
                Size:4,
            }
        }
func (self MasterBootRecord) Mbr_fecha_creacion() Time{
            return Time{
                Super_service:self.Super_service,
                Index:4 + self.Index,
                Size:24,
            }
        }
func (self MasterBootRecord) Mbr_dsk_signature() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:28 + self.Index,
                Size:4,
            }
        }
func (self MasterBootRecord) Dsk_fit() Character{
            return Character{
                Super_service:self.Super_service,
                Index:32 + self.Index,
                Size:1,
            }
        }
func (self MasterBootRecord) Mbr_partitions() ArrayPartition4{
            return ArrayPartition4{
                Super_service:self.Super_service,
                Index:33 + self.Index,
                Size:140,
            }
        }
func (self MasterBootRecord)Dot_label()string {
        return fmt.Sprintf("<TABLE BGCOLOR=\" #45B3F1\" BORDER=\"1\"  COLOR=\"BLACK\"> <TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Mbr_tamano </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Mbr_fecha_creacion </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Mbr_dsk_signature </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Dsk_fit </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Mbr_partitions </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR> </TABLE>",self.Mbr_tamano().Dot_label(),self.Mbr_fecha_creacion().Dot_label(),self.Mbr_dsk_signature().Dot_label(),self.Dsk_fit().Dot_label(),self.Mbr_partitions().Dot_label())
    }
type ExtendedBootRecord struct{
        Super_service *datamanagment.IOService
        Index int32
        Size int32
    }
func CreateExtendedBootRecord(Super_service *datamanagment.IOService,Index int32) ExtendedBootRecord{
        return ExtendedBootRecord{
            Super_service: Super_service,
            Index: Index,
            Size:30,
        }
    }
type ExtendedBootRecordHolder struct{
        Part_mount string
Part_fit string
Part_start int32
Part_s int32
Part_next int32
Part_name [16]string

    }
func (self ExtendedBootRecord)Get() ExtendedBootRecordHolder{
        return ExtendedBootRecordHolder{
            Part_mount: self.Part_mount().Get(),
Part_fit: self.Part_fit().Get(),
Part_start: self.Part_start().Get(),
Part_s: self.Part_s().Get(),
Part_next: self.Part_next().Get(),
Part_name: self.Part_name().Get(),

        }
    }
func (self ExtendedBootRecord)Set(trgt ExtendedBootRecordHolder) {
        self.Part_mount().Set(trgt.Part_mount)
self.Part_fit().Set(trgt.Part_fit)
self.Part_start().Set(trgt.Part_start)
self.Part_s().Set(trgt.Part_s)
self.Part_next().Set(trgt.Part_next)
self.Part_name().Set(trgt.Part_name)

    }
func (self ExtendedBootRecord) Part_mount() Character{
            return Character{
                Super_service:self.Super_service,
                Index:0 + self.Index,
                Size:1,
            }
        }
func (self ExtendedBootRecord) Part_fit() Character{
            return Character{
                Super_service:self.Super_service,
                Index:1 + self.Index,
                Size:1,
            }
        }
func (self ExtendedBootRecord) Part_start() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:2 + self.Index,
                Size:4,
            }
        }
func (self ExtendedBootRecord) Part_s() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:6 + self.Index,
                Size:4,
            }
        }
func (self ExtendedBootRecord) Part_next() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:10 + self.Index,
                Size:4,
            }
        }
func (self ExtendedBootRecord) Part_name() ArrayCharacter16{
            return ArrayCharacter16{
                Super_service:self.Super_service,
                Index:14 + self.Index,
                Size:16,
            }
        }
func (self ExtendedBootRecord)Dot_label()string {
        return fmt.Sprintf("<TABLE BGCOLOR=\" #45B3F1\" BORDER=\"1\"  COLOR=\"BLACK\"> <TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_mount </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_fit </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_start </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_s </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_next </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  Part_name </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR> </TABLE>",self.Part_mount().Dot_label(),self.Part_fit().Dot_label(),self.Part_start().Dot_label(),self.Part_s().Dot_label(),self.Part_next().Dot_label(),self.Part_name().Dot_label())
    }
type SuperBlock struct{
        Super_service *datamanagment.IOService
        Index int32
        Size int32
    }
func CreateSuperBlock(Super_service *datamanagment.IOService,Index int32) SuperBlock{
        return SuperBlock{
            Super_service: Super_service,
            Index: Index,
            Size:108,
        }
    }
type SuperBlockHolder struct{
        S_filesystem_type int32
S_inodes_count int32
S_blocks_count int32
S_free_blocks_count int32
S_free_inodes_count int32
S_mtime TimeHolder
S_umtime TimeHolder
S_mnt_count int32
S_magic int32
S_inode_s int32
S_block_s int32
S_firts_ino int32
S_first_blo int32
S_bm_inode_start int32
S_bm_block_start int32
S_inode_start int32
S_block_start int32

    }
func (self SuperBlock)Get() SuperBlockHolder{
        return SuperBlockHolder{
            S_filesystem_type: self.S_filesystem_type().Get(),
S_inodes_count: self.S_inodes_count().Get(),
S_blocks_count: self.S_blocks_count().Get(),
S_free_blocks_count: self.S_free_blocks_count().Get(),
S_free_inodes_count: self.S_free_inodes_count().Get(),
S_mtime: self.S_mtime().Get(),
S_umtime: self.S_umtime().Get(),
S_mnt_count: self.S_mnt_count().Get(),
S_magic: self.S_magic().Get(),
S_inode_s: self.S_inode_s().Get(),
S_block_s: self.S_block_s().Get(),
S_firts_ino: self.S_firts_ino().Get(),
S_first_blo: self.S_first_blo().Get(),
S_bm_inode_start: self.S_bm_inode_start().Get(),
S_bm_block_start: self.S_bm_block_start().Get(),
S_inode_start: self.S_inode_start().Get(),
S_block_start: self.S_block_start().Get(),

        }
    }
func (self SuperBlock)Set(trgt SuperBlockHolder) {
        self.S_filesystem_type().Set(trgt.S_filesystem_type)
self.S_inodes_count().Set(trgt.S_inodes_count)
self.S_blocks_count().Set(trgt.S_blocks_count)
self.S_free_blocks_count().Set(trgt.S_free_blocks_count)
self.S_free_inodes_count().Set(trgt.S_free_inodes_count)
self.S_mtime().Set(trgt.S_mtime)
self.S_umtime().Set(trgt.S_umtime)
self.S_mnt_count().Set(trgt.S_mnt_count)
self.S_magic().Set(trgt.S_magic)
self.S_inode_s().Set(trgt.S_inode_s)
self.S_block_s().Set(trgt.S_block_s)
self.S_firts_ino().Set(trgt.S_firts_ino)
self.S_first_blo().Set(trgt.S_first_blo)
self.S_bm_inode_start().Set(trgt.S_bm_inode_start)
self.S_bm_block_start().Set(trgt.S_bm_block_start)
self.S_inode_start().Set(trgt.S_inode_start)
self.S_block_start().Set(trgt.S_block_start)

    }
func (self SuperBlock) S_filesystem_type() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:0 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_inodes_count() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:4 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_blocks_count() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:8 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_free_blocks_count() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:12 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_free_inodes_count() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:16 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_mtime() Time{
            return Time{
                Super_service:self.Super_service,
                Index:20 + self.Index,
                Size:24,
            }
        }
func (self SuperBlock) S_umtime() Time{
            return Time{
                Super_service:self.Super_service,
                Index:44 + self.Index,
                Size:24,
            }
        }
func (self SuperBlock) S_mnt_count() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:68 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_magic() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:72 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_inode_s() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:76 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_block_s() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:80 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_firts_ino() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:84 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_first_blo() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:88 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_bm_inode_start() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:92 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_bm_block_start() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:96 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_inode_start() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:100 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock) S_block_start() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:104 + self.Index,
                Size:4,
            }
        }
func (self SuperBlock)Dot_label()string {
        return fmt.Sprintf("<TABLE BGCOLOR=\" #45B3F1\" BORDER=\"1\"  COLOR=\"BLACK\"> <TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_filesystem_type </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_inodes_count </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_blocks_count </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_free_blocks_count </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_free_inodes_count </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_mtime </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_umtime </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_mnt_count </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_magic </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_inode_s </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_block_s </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_firts_ino </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_first_blo </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_bm_inode_start </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_bm_block_start </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_inode_start </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  S_block_start </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR> </TABLE>",self.S_filesystem_type().Dot_label(),self.S_inodes_count().Dot_label(),self.S_blocks_count().Dot_label(),self.S_free_blocks_count().Dot_label(),self.S_free_inodes_count().Dot_label(),self.S_mtime().Dot_label(),self.S_umtime().Dot_label(),self.S_mnt_count().Dot_label(),self.S_magic().Dot_label(),self.S_inode_s().Dot_label(),self.S_block_s().Dot_label(),self.S_firts_ino().Dot_label(),self.S_first_blo().Dot_label(),self.S_bm_inode_start().Dot_label(),self.S_bm_block_start().Dot_label(),self.S_inode_start().Dot_label(),self.S_block_start().Dot_label())
    }
func CreateArrayInteger16(Super_service *datamanagment.IOService,Index int32) ArrayInteger16{
        return ArrayInteger16{
            Super_service: Super_service,
            Index: Index,
            Size:64,
        }
    }
type ArrayInteger16 struct{
        Super_service *datamanagment.IOService
	    Index int32
	    Size int32
    }
func (self ArrayInteger16)Get() [16]int32{
        var array [16]int32
        proxyes := self.Spread()
        var i int32= 0
        for i < 16 {
            array[i] = proxyes[i].Get()
            i++
        }
        return array
    }
func (self ArrayInteger16)Set(trgt [16]int32) {
        proxyes := self.Spread()
    	var i int32= 0
	    for i < 16 {
            proxyes[i].Set(trgt[i])
            i++
	    }
    }
func (self ArrayInteger16) Spread() [16]Integer {
        var array [16]Integer
	    var i int32= 0
	    for i < 16 {
            array[i] = Integer{
                Super_service:self.Super_service,
                Index: self.Index + (i*4),
                Size:4,
            }
	    	i++
    	}
	return array
    }
func (self ArrayInteger16) No(i int32) Integer {
        return Integer{
                Super_service:self.Super_service,
                Index: self.Index + (i*4),
                Size:4,
            }
    }
func (self ArrayInteger16)Dot_label()string {
                return fmt.Sprintf("<TABLE BORDER=\"0\"><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 0 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 1 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 2 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 3 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 4 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 5 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 6 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 7 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 8 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 9 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 10 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 11 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 12 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 13 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 14 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 15 </TD><TD BORDER=\"0\"> %s </TD></TR></TABLE>",self.No(0).Dot_label(),self.No(1).Dot_label(),self.No(2).Dot_label(),self.No(3).Dot_label(),self.No(4).Dot_label(),self.No(5).Dot_label(),self.No(6).Dot_label(),self.No(7).Dot_label(),self.No(8).Dot_label(),self.No(9).Dot_label(),self.No(10).Dot_label(),self.No(11).Dot_label(),self.No(12).Dot_label(),self.No(13).Dot_label(),self.No(14).Dot_label(),self.No(15).Dot_label())
            }
type IndexNode struct{
        Super_service *datamanagment.IOService
        Index int32
        Size int32
    }
func CreateIndexNode(Super_service *datamanagment.IOService,Index int32) IndexNode{
        return IndexNode{
            Super_service: Super_service,
            Index: Index,
            Size:152,
        }
    }
type IndexNodeHolder struct{
        I_uid int32
I_gid int32
I_s int32
I_atime TimeHolder
I_ctime TimeHolder
I_mtime TimeHolder
I_block [16]int32
I_type string
I_perm [3]string

    }
func (self IndexNode)Get() IndexNodeHolder{
        return IndexNodeHolder{
            I_uid: self.I_uid().Get(),
I_gid: self.I_gid().Get(),
I_s: self.I_s().Get(),
I_atime: self.I_atime().Get(),
I_ctime: self.I_ctime().Get(),
I_mtime: self.I_mtime().Get(),
I_block: self.I_block().Get(),
I_type: self.I_type().Get(),
I_perm: self.I_perm().Get(),

        }
    }
func (self IndexNode)Set(trgt IndexNodeHolder) {
        self.I_uid().Set(trgt.I_uid)
self.I_gid().Set(trgt.I_gid)
self.I_s().Set(trgt.I_s)
self.I_atime().Set(trgt.I_atime)
self.I_ctime().Set(trgt.I_ctime)
self.I_mtime().Set(trgt.I_mtime)
self.I_block().Set(trgt.I_block)
self.I_type().Set(trgt.I_type)
self.I_perm().Set(trgt.I_perm)

    }
func (self IndexNode) I_uid() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:0 + self.Index,
                Size:4,
            }
        }
func (self IndexNode) I_gid() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:4 + self.Index,
                Size:4,
            }
        }
func (self IndexNode) I_s() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:8 + self.Index,
                Size:4,
            }
        }
func (self IndexNode) I_atime() Time{
            return Time{
                Super_service:self.Super_service,
                Index:12 + self.Index,
                Size:24,
            }
        }
func (self IndexNode) I_ctime() Time{
            return Time{
                Super_service:self.Super_service,
                Index:36 + self.Index,
                Size:24,
            }
        }
func (self IndexNode) I_mtime() Time{
            return Time{
                Super_service:self.Super_service,
                Index:60 + self.Index,
                Size:24,
            }
        }
func (self IndexNode) I_block() ArrayInteger16{
            return ArrayInteger16{
                Super_service:self.Super_service,
                Index:84 + self.Index,
                Size:64,
            }
        }
func (self IndexNode) I_type() Character{
            return Character{
                Super_service:self.Super_service,
                Index:148 + self.Index,
                Size:1,
            }
        }
func (self IndexNode) I_perm() ArrayCharacter3{
            return ArrayCharacter3{
                Super_service:self.Super_service,
                Index:149 + self.Index,
                Size:3,
            }
        }
func (self IndexNode)Dot_label()string {
        return fmt.Sprintf("<TABLE BGCOLOR=\" #45B3F1\" BORDER=\"1\"  COLOR=\"BLACK\"> <TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  I_uid </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  I_gid </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  I_s </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  I_atime </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  I_ctime </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  I_mtime </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  I_block </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  I_type </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  I_perm </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR> </TABLE>",self.I_uid().Dot_label(),self.I_gid().Dot_label(),self.I_s().Dot_label(),self.I_atime().Dot_label(),self.I_ctime().Dot_label(),self.I_mtime().Dot_label(),self.I_block().Dot_label(),self.I_type().Dot_label(),self.I_perm().Dot_label())
    }
type Content struct{
        Super_service *datamanagment.IOService
        Index int32
        Size int32
    }
func CreateContent(Super_service *datamanagment.IOService,Index int32) Content{
        return Content{
            Super_service: Super_service,
            Index: Index,
            Size:16,
        }
    }
type ContentHolder struct{
        B_name [12]string
B_inodo int32

    }
func (self Content)Get() ContentHolder{
        return ContentHolder{
            B_name: self.B_name().Get(),
B_inodo: self.B_inodo().Get(),

        }
    }
func (self Content)Set(trgt ContentHolder) {
        self.B_name().Set(trgt.B_name)
self.B_inodo().Set(trgt.B_inodo)

    }
func (self Content) B_name() ArrayCharacter12{
            return ArrayCharacter12{
                Super_service:self.Super_service,
                Index:0 + self.Index,
                Size:12,
            }
        }
func (self Content) B_inodo() Integer{
            return Integer{
                Super_service:self.Super_service,
                Index:12 + self.Index,
                Size:4,
            }
        }
func (self Content)Dot_label()string {
        return fmt.Sprintf("<TABLE BGCOLOR=\" #45B3F1\" BORDER=\"1\"  COLOR=\"BLACK\"> <TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  B_name </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR><TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  B_inodo </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR> </TABLE>",self.B_name().Dot_label(),self.B_inodo().Dot_label())
    }
func CreateArrayContent4(Super_service *datamanagment.IOService,Index int32) ArrayContent4{
        return ArrayContent4{
            Super_service: Super_service,
            Index: Index,
            Size:64,
        }
    }
type ArrayContent4 struct{
        Super_service *datamanagment.IOService
	    Index int32
	    Size int32
    }
func (self ArrayContent4)Get() [4]ContentHolder{
        var array [4]ContentHolder
        proxyes := self.Spread()
        var i int32= 0
        for i < 4 {
            array[i] = proxyes[i].Get()
            i++
        }
        return array
    }
func (self ArrayContent4)Set(trgt [4]ContentHolder) {
        proxyes := self.Spread()
    	var i int32= 0
	    for i < 4 {
            proxyes[i].Set(trgt[i])
            i++
	    }
    }
func (self ArrayContent4) Spread() [4]Content {
        var array [4]Content
	    var i int32= 0
	    for i < 4 {
            array[i] = Content{
                Super_service:self.Super_service,
                Index: self.Index + (i*16),
                Size:16,
            }
	    	i++
    	}
	return array
    }
func (self ArrayContent4) No(i int32) Content {
        return Content{
                Super_service:self.Super_service,
                Index: self.Index + (i*16),
                Size:16,
            }
    }
func (self ArrayContent4)Dot_label()string {
                return fmt.Sprintf("<TABLE BORDER=\"0\"><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 0 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 1 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 2 </TD><TD BORDER=\"0\"> %s </TD></TR><TR><TD BORDER=\"1\" BGCOLOR=\"WHITE\"> 3 </TD><TD BORDER=\"0\"> %s </TD></TR></TABLE>",self.No(0).Dot_label(),self.No(1).Dot_label(),self.No(2).Dot_label(),self.No(3).Dot_label())
            }
type DirectoryBlock struct{
        Super_service *datamanagment.IOService
        Index int32
        Size int32
    }
func CreateDirectoryBlock(Super_service *datamanagment.IOService,Index int32) DirectoryBlock{
        return DirectoryBlock{
            Super_service: Super_service,
            Index: Index,
            Size:64,
        }
    }
type DirectoryBlockHolder struct{
        B_content [4]ContentHolder

    }
func (self DirectoryBlock)Get() DirectoryBlockHolder{
        return DirectoryBlockHolder{
            B_content: self.B_content().Get(),

        }
    }
func (self DirectoryBlock)Set(trgt DirectoryBlockHolder) {
        self.B_content().Set(trgt.B_content)

    }
func (self DirectoryBlock) B_content() ArrayContent4{
            return ArrayContent4{
                Super_service:self.Super_service,
                Index:0 + self.Index,
                Size:64,
            }
        }
func (self DirectoryBlock)Dot_label()string {
        return fmt.Sprintf("<TABLE BGCOLOR=\" #45B3F1\" BORDER=\"1\"  COLOR=\"BLACK\"> <TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  B_content </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR> </TABLE>",self.B_content().Dot_label())
    }
type FileBlock struct{
        Super_service *datamanagment.IOService
        Index int32
        Size int32
    }
func CreateFileBlock(Super_service *datamanagment.IOService,Index int32) FileBlock{
        return FileBlock{
            Super_service: Super_service,
            Index: Index,
            Size:64,
        }
    }
type FileBlockHolder struct{
        B_content [64]string

    }
func (self FileBlock)Get() FileBlockHolder{
        return FileBlockHolder{
            B_content: self.B_content().Get(),

        }
    }
func (self FileBlock)Set(trgt FileBlockHolder) {
        self.B_content().Set(trgt.B_content)

    }
func (self FileBlock) B_content() ArrayCharacter64{
            return ArrayCharacter64{
                Super_service:self.Super_service,
                Index:0 + self.Index,
                Size:64,
            }
        }
func (self FileBlock)Dot_label()string {
        return fmt.Sprintf("<TABLE BGCOLOR=\" #45B3F1\" BORDER=\"1\"  COLOR=\"BLACK\"> <TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  B_content </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR> </TABLE>",self.B_content().Dot_label())
    }
type PointerBlock struct{
        Super_service *datamanagment.IOService
        Index int32
        Size int32
    }
func CreatePointerBlock(Super_service *datamanagment.IOService,Index int32) PointerBlock{
        return PointerBlock{
            Super_service: Super_service,
            Index: Index,
            Size:64,
        }
    }
type PointerBlockHolder struct{
        B_pointers [16]int32

    }
func (self PointerBlock)Get() PointerBlockHolder{
        return PointerBlockHolder{
            B_pointers: self.B_pointers().Get(),

        }
    }
func (self PointerBlock)Set(trgt PointerBlockHolder) {
        self.B_pointers().Set(trgt.B_pointers)

    }
func (self PointerBlock) B_pointers() ArrayInteger16{
            return ArrayInteger16{
                Super_service:self.Super_service,
                Index:0 + self.Index,
                Size:64,
            }
        }
func (self PointerBlock)Dot_label()string {
        return fmt.Sprintf("<TABLE BGCOLOR=\" #45B3F1\" BORDER=\"1\"  COLOR=\"BLACK\"> <TR><TD BORDER=\"1\" COLOR=\" #135F8A\">  B_pointers </TD><TD BORDER=\"1\" COLOR=\" #135F8A\"> %s </TD></TR> </TABLE>",self.B_pointers().Dot_label())
    }
