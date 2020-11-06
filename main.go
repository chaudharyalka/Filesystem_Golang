// This program creates Shell with basic functionality.
//commands tested with parameters - pwd,cd,ls,mkdir,rmdir
//Exculsive case cd: cd is not there because its the functionality of shell its not there in /bin

package main 

import (
	"os"
	"fmt"
	"bufio"
	"os/exec"
	"strings"
)

func main() {
	cmd:= make(chan string)
	go InputReader(cmd)
	go InputExecuter(cmd)
	//Block forever
	select {}
}

//InputReader this function reads user command from shell passes it to the executor
func InputReader(cmd chan string) {
	reader:= bufio.NewReader(os.Stdin)
	for {
		usrCmd, _ := reader.ReadString('\n')
		cmd <-  strings.TrimSuffix(usrCmd, "\n")
	}
}

//InputExecuter this function excutes the user command.
func InputExecuter(cmd chan string) {
	for {
		args := strings.Split(<-cmd," ")
		if len(args) > 0 {
			switch args[0] {
				case "cd":
					//if without argument then go the home dir 
					if len(args) > 1 {
						os.Chdir(args[1])
					} else {
						os.Chdir("~")
						
					}

				default:
					// Set the correct output device.
					input := exec.Command(args[0],args[1:]...)
				   	input.Stderr = os.Stderr
					input.Stdout = os.Stdout

					if err:= input.Run(); err != nil {
						fmt.Println("Error: ",err)
					}
				}
			}
		}
}
