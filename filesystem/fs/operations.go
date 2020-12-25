package fs

import (
	"Filesystem_Golang/filesystem/disklib"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// LoadFileSystem loads content from dumpfile into memory.
func (fs *FileSystem) LoadFileSystem() {
	// open dump file
	file, err := os.Open(DumpFile)
	if err != nil {
		fmt.Println("Error in opening file while loading the filesystem: ", err)
		os.Exit(1)
	}
	defer file.Close()

	// first load the fst table
	if fs.fileSystemTable == nil {
		fs.fileSystemTable = make(map[string]int)
	}
	if err := fs.fileSystemTable.loadFST(file); err != nil {
		fmt.Println("Unable to load FST ", err)
		os.Exit(1)
	}

	// Load Free inodes list
	if err := fs.nextFreeInode.loadNextFreeInode(file); err != nil {
		fmt.Println("Unable to load FST ", err)
		os.Exit(1)
	}

	// Load Free data block list
	if err := fs.nextFreeDataBlock.loadNextFreeDataBlock(file); err != nil {
		fmt.Println("Unable to load FST ", err)
		os.Exit(1)
	}

	fmt.Println("Successfully loaded filesystem from dump file")

}

func (fst fsTable) loadFST(file *os.File) error {
	data, err := disklib.ReadBlockFromDisk(file, OffsetFST, SizeFST)
	if err != nil {
		return err
	}

	idx := 0
	for idx < SizeFST {
		value := strings.TrimSpace(string(data[idx : idx+3]))
		fileName := strings.TrimSpace(string(data[idx+3 : idx+13]))
		if fileName != "" {
			fst[fileName], _ = strconv.Atoi(value)
		}
		idx = idx + CellFSTSize
	}
	fmt.Println("after loading ", fst)

	return nil
}

func (list *inodeList) loadNextFreeInode(file *os.File) error {
	data, err := disklib.ReadBlockFromDisk(file, OffsetInodeTable, SizeInodeTable)
	if err != nil {
		return err
	}

	idx := 0
	for _, val := range data {
		if int(val) == 0 {
			*list = append(*list, idx)
		}
		idx++
	}

	return nil
}

func (list *dataBlockList) loadNextFreeDataBlock(file *os.File) error {
	data, err := disklib.ReadBlockFromDisk(file, OffsetDataTable, SizeDataTable)
	if err != nil {
		return err
	}

	idx := 0
	for _, val := range data {
		if int(val) == 0 {
			*list = append(*list, idx)
		}
		idx++
	}

	return nil
}

// EnterRootDir will check root dir exists or not.
func (fs *FileSystem) EnterRootDir() bool {
	fmt.Println("fstable", fs.fileSystemTable)
	_, ok := fs.fileSystemTable["root"]
	return ok
}

// CreateFile is used to create files and dir.
// TODO can use channels which will upadte the dump file.
func (fs *FileSystem) CreateFile(fileName string, data string, parentInodeNum int, fileType int) error {

	// Validation of the arguments
	// TODO same name file in the directory.
	if err := fs.validateCreationRequest(fileName); err != nil {
		fmt.Println("Error: Creation request fails while validating : ", err)
		return err
	}
	dataBlockRequired := int(math.Ceil(float64(len(data) / DataBlockSize)))
	// Check resources available or not
	if err := resourceAvailable(fs, dataBlockRequired); err != nil {
		fmt.Println("Error: Creation request fails while check availabilty of resource : ", err)
		return err
	}

	fmt.Println("filename", fileName, "datablockrequired", dataBlockRequired)
	// Get Parent Inode
	parInode, err := getInodeInfo(parentInodeNum)
	if err != nil {
		fmt.Println("Unable to get parent inode ", err)
		return err
	}

	// check parent inode has space to accomodate new file/ directory inside it.
	// here 4 is used because 1 for comma and 3 bytes representing inode number.
	if len(parInode) < (InodeBlockSize - 4) {
		return fmt.Errorf("Parent inode doesn't have space left to accomodate new file in it")
	}

	// Allocate an inode and intialise
	if dataBlockRequired != 0 {
		dataBlockRequired++
	}
	inode := inode{fs.nextFreeInode[0], fileType, parentInodeNum, fs.nextFreeDataBlock[:dataBlockRequired]}

	fmt.Println("inode", inode)
	// Update fst with new inode entries.
	fs.UpdateFst(inode)

	// Add entry in FST in memory
	fs.fileSystemTable[fileName] = inode.inodeNum

	parentInode := parseInode(parInode)
	parentInode.dataList = append(parentInode.dataList, inode.inodeNum)

	// Update the dumpFile with the file content.
	if err := UpdateDumpFile(inode, data, fileName, parentInode, parentInodeNum); err != nil {
		fmt.Println("unable to update the disk : ", err)
		return err
	}

	// TODO : After successfull creation of file, update the directory data block accordingly..

	fmt.Println("successful updation in disk", inode)

	return nil
}

// UpdateFst will update in memory struct whenever some file is created or deleted.
func (fs *FileSystem) UpdateFst(file inode) error {
	fs.nextFreeInode = fs.nextFreeInode[1:]
	fs.nextFreeDataBlock = fs.nextFreeDataBlock[len(file.dataList):]
	return nil
}

func (fs *FileSystem) validateCreationRequest(fileName string) error {
	if len(fileName) > FileNameSize {
		return fmt.Errorf("File name length exceeds upperlimit. we only support file names of length %d", FileNameSize)
	}

	//TODO: Check fileName already exists in folder or not!!!
	// Get all inodes under parent folder ; then check for these inodes corresponding filename ; matches or not!!
	return nil
}

func resourceAvailable(fs *FileSystem, dataBlockRequired int) error {
	if len(fs.nextFreeInode) == 0 {
		return fmt.Errorf("Insufficient resource: out of inode block")
	}

	if len(fs.nextFreeDataBlock) < dataBlockRequired {
		return fmt.Errorf("INsufficient resource: out of data blocks")
	}
	return nil
}

// This will give parent inode.
func getInodeInfo(inodeNum int) ([]byte, error) {
	var inode []byte

	// open dump file
	file, err := os.Open(DumpFile)
	if err != nil {
		fmt.Println("Error in opening dump file: ", err)
		return inode, err
	}
	defer file.Close()

	offset := OffsetInodeBlock + (inodeNum * InodeBlockSize)
	data, err := disklib.ReadBlockFromDisk(file, offset, InodeBlockSize)
	if err != nil {
		return inode, err
	}

	return data, err

}

//TODO: test all the following with read operations.
// UpdateDumpFile will update the disk whenever there is a creation or deletion of file.
func UpdateDumpFile(fileInode inode, data string, fileName string, parentInode inode, parentInodeNum int) error {
	// open dump file
	file, err := os.OpenFile(DumpFile, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Error in opening file while loading the filesystem: ", err)
		os.Exit(1)
	}
	defer file.Close()

	offset := OffsetInodeTable + fileInode.inodeNum
	// update inode bitmap
	if err := disklib.WriteBlockOnDisk(file, offset, []byte(strconv.Itoa(SetBit))); err != nil {
		fmt.Println("Error: unable to update inode bitmap", err)
		return err
	}

	// update inode block
	offset = OffsetInodeBlock + (fileInode.inodeNum * InodeBlockSize)
	inodeInfo := frameInodeString(fileInode)
	if err := disklib.WriteBlockOnDisk(file, offset, []byte(inodeInfo)); err != nil {
		fmt.Println("Error: unable to update inode block", err)
		return err
	}

	//update parent inode block
	if parentInodeNum != -1 { // skipping it for root directory
		offset = OffsetInodeBlock + (parentInodeNum * InodeBlockSize)
		inodeInfo = frameInodeString(parentInode)
		if err := disklib.WriteBlockOnDisk(file, offset, []byte(inodeInfo)); err != nil {
			fmt.Println("Error: unable to update parent inode block", err)
			return err
		}
	}

	// update data bitmap and fill the data block as well ...
	offset = OffsetDataTable
	listByte := []byte(data)
	count := 0
	for _, val := range fileInode.dataList {
		// update inode bitmap
		if err := disklib.WriteBlockOnDisk(file, offset+val, []byte(strconv.Itoa(SetBit))); err != nil {
			fmt.Println("Error: unable to update inode bitmap", err)
			return err
		}

		// update data block in disk
		if err := disklib.WriteBlockOnDisk(file, OffsetDataBlock+(val*DataBlockSize), []byte(listByte[count:count+DataBlockSize])); err != nil {
			fmt.Println("Error: unable to update inode bitmap", err)
			return err
		}
		count += DataBlockSize
	}

	// put fst in disk
	offset = OffsetFST + (fileInode.inodeNum * CellFSTSize)
	str := fmt.Sprintf("%3d", fileInode.inodeNum) + fmt.Sprintf("%10s", fileName)
	if err := disklib.WriteBlockOnDisk(file, offset, []byte(str)); err != nil {
		fmt.Println("Error: unable to update inode bitmap", err)
		return err
	}
	return nil

}

func parseInode(byteData []byte) inode {
	var node inode
	str := string(byteData)
	list := strings.Split(str, ",")
	idx := 0
	for idx < len(list) {
		elem, _ := strconv.Atoi(list[idx])
		switch idx {
		case 0:
			node.inodeNum = elem
		case 1:
			node.inodeType = elem
		case 2:
			node.parentInode = elem
		default:
			node.dataList = append(node.dataList, elem)
		}
		idx++
	}

	return node
}

func frameInodeString(fileInode inode) string {
	var str string
	str = str + strconv.Itoa(fileInode.inodeNum) + "," + strconv.Itoa(fileInode.inodeType) + "," + strconv.Itoa(fileInode.parentInode) + "," + strings.Trim(strings.Replace(fmt.Sprint(fileInode.dataList), " ", ",", -1), "[]")
	fmt.Println(str)
	return str
}
