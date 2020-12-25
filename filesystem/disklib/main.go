// This file contains operations related to disk
package disklib

import (
	"os"
)

//const blockSize = 4
//const dumpFile="../dump.txt"

// setFileOffset sets the offset of file to read/write the blocks
func setFileOffset(offset int, file *os.File) error {
	whence := 0 //Offset with reference to the origin
	if _, err := file.Seek(int64(offset), whence); err != nil {
		return err
	}
	return nil
}

// ReadBlockFromDisk helps in reding block from dumpFile
func ReadBlockFromDisk(file *os.File, offset int, blockSize int) ([]byte, error) {
	setFileOffset(offset, file)
	data := make([]byte, blockSize)
	if _, err := file.Read(data); err != nil {
		return data, err
	}
	return data, nil
}

// WriteBlockOnDisk helps in writing disk.
func WriteBlockOnDisk(file *os.File, offset int, data []byte) error {
	setFileOffset(offset, file)
	if _, err := file.Write(data); err != nil {
		return err
	}
	return nil
}

// readBlock reads disk block blocknr from the disk pointed to by disk into a buffer pointed to by block.
//func readBlock() ([]byte,error) {
//
//	data := make([]byte, blockSize)
//	file, err := os.Open(dumpFile)
//	if err != nil {
//		return data,err
//	}
//	defer file.Close()
//
//	if _,err := file.Read(data); err != nil {
//			return data,err
//	}
//	return data,nil
//}
//
//// writeBlock writes disk block.
//func writeBlock(data []byte) error {
//	file, err := os.OpenFile(dumpFile,os.O_WRONLY,os.ModeAppend)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	if len(data) > blockSize {
//		return fmt.Errorf("size of data is greater than blocksize")
//	}
//
//	if _,err := file.Write(data); err != nil {
//		return err
//	}
//	return nil
//}

//// setFileOffset sets the offset of file to read/write the blocks
//func setFileOffset(blockNum int, file *os.File) error {
////	defer file.Close()
//	offset:= blockNum * blockSize
//	whence:= 0 //Offset with reference to the origin
//	if _,err := file.Seek(int64(offset),whence); err != nil {
//		return err
//	}
//	return nil
//}
//

//func main() {
//   	file, err := os.Open(dumpFile)
//	if err != nil {
//		fmt.Println("Error in opening file : ", err)
//		return
//	}
//
//	blockNum:= 0
//	if err := setFileOffset(blockNum,file); err != nil {
//			fmt.Println("Unable to set file offset : ",err)
//			return
//	}
//
////	content:= []byte{97,98,99,100}
////	if err := writeBlock(content); err != nil {
////			fmt.Println("Unable to write Block : ",err)
////			return
////	}
//
//	if data,err := readBlock(); err != nil {
//			fmt.Println("Unable to read Block : ",err)
//			return
//	} else {
//			fmt.Println("data read ", string(data))
//	}
//
//
//}
