package fs

// Following constants are used across the packages.
const (
	DumpFile        string = "dump.txt"
	MaxNumFile      int    = 1000   // Maximum number of files that can exist in the system
	MaxNumDataBlock int    = 100000 // Data blocks in the system
	DataBlockSize   int    = 4      // Data block size in bytes
	InodeBlockSize  int    = 100    // Inode block size in bytes
	FileNameSize    int    = 10     // Maximum file name length supported

	CellFSTSize int = FileNameSize + 3 // here 3 reperents memory used to represent MaxNumFile

	SizeFST        int = CellFSTSize * MaxNumFile // Total space(in bytes) taken by FST
	SizeInodeTable int = 1 * MaxNumFile           // 1 byte is used to represent state of an inode.
	SizeDataTable  int = 1 * MaxNumDataBlock      // 1 byte is used to represent state of an datablock.

	OffsetFST        int = 0
	OffsetInodeTable int = SizeFST
	OffsetDataTable  int = OffsetInodeTable + SizeInodeTable
	OffsetInodeBlock int = OffsetDataTable + SizeDataTable
	OffsetDataBlock  int = OffsetInodeBlock + (InodeBlockSize * MaxNumFile)

	SetBit   int = 1
	UnsetBit int = 0

	FileType int = 0
	DirType  int = 1
)

type dataBlock string
type fsTable map[string]int
type inodeList []int
type dataBlockList []int

type inode struct {
	inodeNum    int   // Unique number across fs
	inodeType   int   // 0 means file , 1 means dir
	parentInode int   // inode number of its parent.
	dataList    []int // for inode_type 0 this will contain block numbers of file otherwise this will contain inode number of file inside this dir.
}

// Inode size distribution:
// inodeNum = 3 | type = 1 | parent inode = 3 | datablocks = 93

// FileSystem struct stores content of dumpFile.
type FileSystem struct {
	fileSystemTable   fsTable       // will contain mapping of filename and inode.
	nextFreeInode     inodeList     // will contain inode number which is free.
	nextFreeDataBlock dataBlockList // will contain data block number which is free.
}
