package main

import (
	"Filesystem_Golang/filesystem/fs"
	"fmt"
)

func main() {
	// lets load the filesystem from the disk.
	var fsystem fs.FileSystem
	fsystem.LoadFileSystem()

	//TODO: start shell // for now leaving this

	//Enter root directory
//	if !fsystem.EnterRootDir() {
		if err := fsystem.CreateFile("root", "", -1, fs.DirType); err != nil {
			fmt.Println("Unable to enter root directory: ", err)
		}
//	}

//	panic("failing")
	// trying creation request
//	filename := "abc.txt"
//	value := "abcdef"
//	err := fsystem.CreateFile(filename, value, 0, 0)
//
//	fmt.Println("creation request failure", err)

	// Try creating directory!!
	//folderName:= "MyDir"
}
