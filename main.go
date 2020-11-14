package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"hellofudango/hellofudan"
)

const filename = "./accounts.txt"

func main() {
	stuList := readStudentAccounts()
	if len(stuList) == 0 {
		log.Println("Didn't find account infomation, please input")
		stuList = append(stuList, inputAccount())
	}

	// hfm := hellofudan.NewManager(stuList)
	// hfm.Start()
	fmt.Println(stuList)
}

func readStudentAccounts() (stuList []hellofudan.Student) {
	fd, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	defer fd.Close()

	rd := bufio.NewReader(fd)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			break
		}
		info := strings.Split(line, " ")
		stuList = append(stuList, hellofudan.Student{
			StudentID: info[0],
			Password:  info[1],
		})
	}

	return
}

func inputAccount() hellofudan.Student {
	var account, password string
	for {
		fmt.Printf("👤Account : ")
		fmt.Scanf("%s", &account)
		if len(account) > 0 {
			break
		}
	}
	for {
		fmt.Printf("🔑Password: ")
		fmt.Scanf("%s", &password)
		if len(password) > 0 {
			break
		}
	}

	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	wd := bufio.NewWriter(fd)
	defer func () {
		wd.Flush()
		fmt.Println("Account has been saved in accounts.txt, PLEASE KEEP IT SAFE.")
		fmt.Println("You can also add other people's account information in the same format by appending it to the file.")
	}

	_, err = wd.WriteString(fmt.Sprintf("%s %s\n", account, password))
	if err != nil {
		panic(err)
	}

	return hellofudan.Student{
		StudentID: account,
		Password:  password,
	}
}
