package main

import (
	"bufio"
	"io/ioutil"
	"time"

	//"io"
	"math/rand"
	"os"
	"strconv"

	//"fmt"      // пакет для форматированного ввода вывода
	"log" // пакет для логирования

	"github.com/jlaffaye/ftp"
	//"github.com/RealJK/rss-parser-go"
	//"strings"  // пакет для работы с  UTF-8 строками
)

func tryGetCurrentDir(serverConnection *ftp.ServerConn) string {
	currentDir, err := serverConnection.CurrentDir()
	if err != nil {
		return "error"
	} else {
		return currentDir
	}
}
func tryMakeDir(serverConnection *ftp.ServerConn, way string) {
	err := serverConnection.MakeDir(way)
	if err != nil {
		log.Fatal(err)
	}
}
func tryMakeRandomDir(serverConnection *ftp.ServerConn) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	tryMakeDir(serverConnection,
		"/yetAnotherFolder"+strconv.Itoa(r1.Int()))
}

func proceedOperations(serverConnection *ftp.ServerConn) {
	log.Print(tryGetCurrentDir(serverConnection))

	file, err := os.Open("fileToServer.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	reader := bufio.NewReader(file)
	err = serverConnection.Stor("fileFromClient.txt", reader)
	if err != nil {
		log.Fatal(err)
	}

	fileReader, fileReadError := serverConnection.Retr("fileToClient.txt")
	if fileReadError != nil {
		log.Fatal(fileReadError)
	} else {
		buf, err2 := ioutil.ReadAll(fileReader)

		fileSaveError := ioutil.WriteFile("fileFromServer.txt", buf, 0777)
		if fileSaveError != nil {
			log.Println(fileSaveError)
		}

		if err2 != nil {
			log.Fatal(err2)
		}
		fileReader.Close()
	}
}

func main() {
	time.Sleep(1000 * time.Millisecond)
	serverConnection, err := ftp.Dial("lab.posevin.com:9008")
	if err != nil {
		log.Fatal(err)
	} else {
		loginErr := serverConnection.Login("admin", "123456")
		if err != nil {
			log.Fatal(loginErr)
		} else {
			log.Print(tryGetCurrentDir(serverConnection))
			tryMakeRandomDir(serverConnection)
			proceedOperations(serverConnection)
		}
	}
}
