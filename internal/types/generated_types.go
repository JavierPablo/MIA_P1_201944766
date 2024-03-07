package types
import "project/internal/datamanagment"

type Time struct{
        super_service *datamanagment.IOService
        index int32
        Size int32
    }
func CreateTime(super_service *datamanagment.IOService,index int32) Time{
        return Time{
            super_service: super_service,
            index: index,
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
                super_service:self.super_service,
                index:0 + self.index,
                Size:4,
            }
        }
func (self Time) Minute() Integer{
            return Integer{
                super_service:self.super_service,
                index:4 + self.index,
                Size:4,
            }
        }
func (self Time) Second() Integer{
            return Integer{
                super_service:self.super_service,
                index:8 + self.index,
                Size:4,
            }
        }
func (self Time) Day() Integer{
            return Integer{
                super_service:self.super_service,
                index:12 + self.index,
                Size:4,
            }
        }
func (self Time) Month() Integer{
            return Integer{
                super_service:self.super_service,
                index:16 + self.index,
                Size:4,
            }
        }
func (self Time) Year() Integer{
            return Integer{
                super_service:self.super_service,
                index:20 + self.index,
                Size:4,
            }
        }
func CreateArrayCharacter4(super_service *datamanagment.IOService,index int32) ArrayCharacter4{
        return ArrayCharacter4{
            super_service: super_service,
            index: index,
            Size:4,
        }
    }
type ArrayCharacter4 struct{
        super_service *datamanagment.IOService
	    index int32
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
                super_service:self.super_service,
                index: self.index + (i*1),
                Size:1,
            }
	    	i++
    	}
	return array
    }
func (self ArrayCharacter4) No(i int32) Character {
        return Character{
                super_service:self.super_service,
                index: self.index + (i*1),
                Size:1,
            }
    }
func CreateArrayCharacter3(super_service *datamanagment.IOService,index int32) ArrayCharacter3{
        return ArrayCharacter3{
            super_service: super_service,
            index: index,
            Size:3,
        }
    }
type ArrayCharacter3 struct{
        super_service *datamanagment.IOService
	    index int32
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
                super_service:self.super_service,
                index: self.index + (i*1),
                Size:1,
            }
	    	i++
    	}
	return array
    }
func (self ArrayCharacter3) No(i int32) Character {
        return Character{
                super_service:self.super_service,
                index: self.index + (i*1),
                Size:1,
            }
    }
func CreateArrayCharacter12(super_service *datamanagment.IOService,index int32) ArrayCharacter12{
        return ArrayCharacter12{
            super_service: super_service,
            index: index,
            Size:12,
        }
    }
type ArrayCharacter12 struct{
        super_service *datamanagment.IOService
	    index int32
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
                super_service:self.super_service,
                index: self.index + (i*1),
                Size:1,
            }
	    	i++
    	}
	return array
    }
func (self ArrayCharacter12) No(i int32) Character {
        return Character{
                super_service:self.super_service,
                index: self.index + (i*1),
                Size:1,
            }
    }
func CreateArrayCharacter16(super_service *datamanagment.IOService,index int32) ArrayCharacter16{
        return ArrayCharacter16{
            super_service: super_service,
            index: index,
            Size:16,
        }
    }
type ArrayCharacter16 struct{
        super_service *datamanagment.IOService
	    index int32
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
                super_service:self.super_service,
                index: self.index + (i*1),
                Size:1,
            }
	    	i++
    	}
	return array
    }
func (self ArrayCharacter16) No(i int32) Character {
        return Character{
                super_service:self.super_service,
                index: self.index + (i*1),
                Size:1,
            }
    }
func CreateArrayCharacter64(super_service *datamanagment.IOService,index int32) ArrayCharacter64{
        return ArrayCharacter64{
            super_service: super_service,
            index: index,
            Size:64,
        }
    }
type ArrayCharacter64 struct{
        super_service *datamanagment.IOService
	    index int32
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
                super_service:self.super_service,
                index: self.index + (i*1),
                Size:1,
            }
	    	i++
    	}
	return array
    }
func (self ArrayCharacter64) No(i int32) Character {
        return Character{
                super_service:self.super_service,
                index: self.index + (i*1),
                Size:1,
            }
    }
type Partition struct{
        super_service *datamanagment.IOService
        index int32
        Size int32
    }
func CreatePartition(super_service *datamanagment.IOService,index int32) Partition{
        return Partition{
            super_service: super_service,
            index: index,
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
                super_service:self.super_service,
                index:0 + self.index,
                Size:1,
            }
        }
func (self Partition) Part_type() Character{
            return Character{
                super_service:self.super_service,
                index:1 + self.index,
                Size:1,
            }
        }
func (self Partition) Part_fit() Character{
            return Character{
                super_service:self.super_service,
                index:2 + self.index,
                Size:1,
            }
        }
func (self Partition) Part_start() Integer{
            return Integer{
                super_service:self.super_service,
                index:3 + self.index,
                Size:4,
            }
        }
func (self Partition) Part_s() Integer{
            return Integer{
                super_service:self.super_service,
                index:7 + self.index,
                Size:4,
            }
        }
func (self Partition) Part_name() ArrayCharacter16{
            return ArrayCharacter16{
                super_service:self.super_service,
                index:11 + self.index,
                Size:16,
            }
        }
func (self Partition) Part_correlative() Integer{
            return Integer{
                super_service:self.super_service,
                index:27 + self.index,
                Size:4,
            }
        }
func (self Partition) Part_id() ArrayCharacter4{
            return ArrayCharacter4{
                super_service:self.super_service,
                index:31 + self.index,
                Size:4,
            }
        }
func CreateArrayPartition4(super_service *datamanagment.IOService,index int32) ArrayPartition4{
        return ArrayPartition4{
            super_service: super_service,
            index: index,
            Size:140,
        }
    }
type ArrayPartition4 struct{
        super_service *datamanagment.IOService
	    index int32
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
                super_service:self.super_service,
                index: self.index + (i*35),
                Size:35,
            }
	    	i++
    	}
	return array
    }
func (self ArrayPartition4) No(i int32) Partition {
        return Partition{
                super_service:self.super_service,
                index: self.index + (i*35),
                Size:35,
            }
    }
type MasterBootRecord struct{
        super_service *datamanagment.IOService
        index int32
        Size int32
    }
func CreateMasterBootRecord(super_service *datamanagment.IOService,index int32) MasterBootRecord{
        return MasterBootRecord{
            super_service: super_service,
            index: index,
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
                super_service:self.super_service,
                index:0 + self.index,
                Size:4,
            }
        }
func (self MasterBootRecord) Mbr_fecha_creacion() Time{
            return Time{
                super_service:self.super_service,
                index:4 + self.index,
                Size:24,
            }
        }
func (self MasterBootRecord) Mbr_dsk_signature() Integer{
            return Integer{
                super_service:self.super_service,
                index:28 + self.index,
                Size:4,
            }
        }
func (self MasterBootRecord) Dsk_fit() Character{
            return Character{
                super_service:self.super_service,
                index:32 + self.index,
                Size:1,
            }
        }
func (self MasterBootRecord) Mbr_partitions() ArrayPartition4{
            return ArrayPartition4{
                super_service:self.super_service,
                index:33 + self.index,
                Size:140,
            }
        }
type ExtendedBootRecord struct{
        super_service *datamanagment.IOService
        index int32
        Size int32
    }
func CreateExtendedBootRecord(super_service *datamanagment.IOService,index int32) ExtendedBootRecord{
        return ExtendedBootRecord{
            super_service: super_service,
            index: index,
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
                super_service:self.super_service,
                index:0 + self.index,
                Size:1,
            }
        }
func (self ExtendedBootRecord) Part_fit() Character{
            return Character{
                super_service:self.super_service,
                index:1 + self.index,
                Size:1,
            }
        }
func (self ExtendedBootRecord) Part_start() Integer{
            return Integer{
                super_service:self.super_service,
                index:2 + self.index,
                Size:4,
            }
        }
func (self ExtendedBootRecord) Part_s() Integer{
            return Integer{
                super_service:self.super_service,
                index:6 + self.index,
                Size:4,
            }
        }
func (self ExtendedBootRecord) Part_next() Integer{
            return Integer{
                super_service:self.super_service,
                index:10 + self.index,
                Size:4,
            }
        }
func (self ExtendedBootRecord) Part_name() ArrayCharacter16{
            return ArrayCharacter16{
                super_service:self.super_service,
                index:14 + self.index,
                Size:16,
            }
        }
type SuperBlock struct{
        super_service *datamanagment.IOService
        index int32
        Size int32
    }
func CreateSuperBlock(super_service *datamanagment.IOService,index int32) SuperBlock{
        return SuperBlock{
            super_service: super_service,
            index: index,
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
                super_service:self.super_service,
                index:0 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_inodes_count() Integer{
            return Integer{
                super_service:self.super_service,
                index:4 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_blocks_count() Integer{
            return Integer{
                super_service:self.super_service,
                index:8 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_free_blocks_count() Integer{
            return Integer{
                super_service:self.super_service,
                index:12 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_free_inodes_count() Integer{
            return Integer{
                super_service:self.super_service,
                index:16 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_mtime() Time{
            return Time{
                super_service:self.super_service,
                index:20 + self.index,
                Size:24,
            }
        }
func (self SuperBlock) S_umtime() Time{
            return Time{
                super_service:self.super_service,
                index:44 + self.index,
                Size:24,
            }
        }
func (self SuperBlock) S_mnt_count() Integer{
            return Integer{
                super_service:self.super_service,
                index:68 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_magic() Integer{
            return Integer{
                super_service:self.super_service,
                index:72 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_inode_s() Integer{
            return Integer{
                super_service:self.super_service,
                index:76 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_block_s() Integer{
            return Integer{
                super_service:self.super_service,
                index:80 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_firts_ino() Integer{
            return Integer{
                super_service:self.super_service,
                index:84 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_first_blo() Integer{
            return Integer{
                super_service:self.super_service,
                index:88 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_bm_inode_start() Integer{
            return Integer{
                super_service:self.super_service,
                index:92 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_bm_block_start() Integer{
            return Integer{
                super_service:self.super_service,
                index:96 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_inode_start() Integer{
            return Integer{
                super_service:self.super_service,
                index:100 + self.index,
                Size:4,
            }
        }
func (self SuperBlock) S_block_start() Integer{
            return Integer{
                super_service:self.super_service,
                index:104 + self.index,
                Size:4,
            }
        }
func CreateArrayInteger16(super_service *datamanagment.IOService,index int32) ArrayInteger16{
        return ArrayInteger16{
            super_service: super_service,
            index: index,
            Size:64,
        }
    }
type ArrayInteger16 struct{
        super_service *datamanagment.IOService
	    index int32
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
                super_service:self.super_service,
                index: self.index + (i*4),
                Size:4,
            }
	    	i++
    	}
	return array
    }
func (self ArrayInteger16) No(i int32) Integer {
        return Integer{
                super_service:self.super_service,
                index: self.index + (i*4),
                Size:4,
            }
    }
type IndexNode struct{
        super_service *datamanagment.IOService
        index int32
        Size int32
    }
func CreateIndexNode(super_service *datamanagment.IOService,index int32) IndexNode{
        return IndexNode{
            super_service: super_service,
            index: index,
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
                super_service:self.super_service,
                index:0 + self.index,
                Size:4,
            }
        }
func (self IndexNode) I_gid() Integer{
            return Integer{
                super_service:self.super_service,
                index:4 + self.index,
                Size:4,
            }
        }
func (self IndexNode) I_s() Integer{
            return Integer{
                super_service:self.super_service,
                index:8 + self.index,
                Size:4,
            }
        }
func (self IndexNode) I_atime() Time{
            return Time{
                super_service:self.super_service,
                index:12 + self.index,
                Size:24,
            }
        }
func (self IndexNode) I_ctime() Time{
            return Time{
                super_service:self.super_service,
                index:36 + self.index,
                Size:24,
            }
        }
func (self IndexNode) I_mtime() Time{
            return Time{
                super_service:self.super_service,
                index:60 + self.index,
                Size:24,
            }
        }
func (self IndexNode) I_block() ArrayInteger16{
            return ArrayInteger16{
                super_service:self.super_service,
                index:84 + self.index,
                Size:64,
            }
        }
func (self IndexNode) I_type() Character{
            return Character{
                super_service:self.super_service,
                index:148 + self.index,
                Size:1,
            }
        }
func (self IndexNode) I_perm() ArrayCharacter3{
            return ArrayCharacter3{
                super_service:self.super_service,
                index:149 + self.index,
                Size:3,
            }
        }
type Content struct{
        super_service *datamanagment.IOService
        index int32
        Size int32
    }
func CreateContent(super_service *datamanagment.IOService,index int32) Content{
        return Content{
            super_service: super_service,
            index: index,
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
                super_service:self.super_service,
                index:0 + self.index,
                Size:12,
            }
        }
func (self Content) B_inodo() Integer{
            return Integer{
                super_service:self.super_service,
                index:12 + self.index,
                Size:4,
            }
        }
func CreateArrayContent4(super_service *datamanagment.IOService,index int32) ArrayContent4{
        return ArrayContent4{
            super_service: super_service,
            index: index,
            Size:64,
        }
    }
type ArrayContent4 struct{
        super_service *datamanagment.IOService
	    index int32
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
                super_service:self.super_service,
                index: self.index + (i*16),
                Size:16,
            }
	    	i++
    	}
	return array
    }
func (self ArrayContent4) No(i int32) Content {
        return Content{
                super_service:self.super_service,
                index: self.index + (i*16),
                Size:16,
            }
    }
type DirectoryBlock struct{
        super_service *datamanagment.IOService
        index int32
        Size int32
    }
func CreateDirectoryBlock(super_service *datamanagment.IOService,index int32) DirectoryBlock{
        return DirectoryBlock{
            super_service: super_service,
            index: index,
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
                super_service:self.super_service,
                index:0 + self.index,
                Size:64,
            }
        }
type FileBlock struct{
        super_service *datamanagment.IOService
        index int32
        Size int32
    }
func CreateFileBlock(super_service *datamanagment.IOService,index int32) FileBlock{
        return FileBlock{
            super_service: super_service,
            index: index,
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
                super_service:self.super_service,
                index:0 + self.index,
                Size:64,
            }
        }
type PointerBlock struct{
        super_service *datamanagment.IOService
        index int32
        Size int32
    }
func CreatePointerBlock(super_service *datamanagment.IOService,index int32) PointerBlock{
        return PointerBlock{
            super_service: super_service,
            index: index,
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
                super_service:self.super_service,
                index:0 + self.index,
                Size:64,
            }
        }
