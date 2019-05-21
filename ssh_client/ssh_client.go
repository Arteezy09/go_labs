package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"strings"
)

func connect(client *ssh.Client, adr, login, password string) (*ssh.Client, bool) {
	config := &ssh.ClientConfig{
		User:            login,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", adr, config)
	if err != nil {
		fmt.Println("Connect error!", err)
		return client, false
	}
	return client, true
}

func makeSessionWithTerminal(client *ssh.Client, session *ssh.Session) (*ssh.Session, bool) {
	session, err := client.NewSession()
	if err != nil {
		fmt.Println("NewSession error!", err)
		return session, false
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	fileDescriptor := int(os.Stdin.Fd())

	if terminal.IsTerminal(fileDescriptor) {
		originalState, err := terminal.MakeRaw(fileDescriptor)
		if err != nil {
			fmt.Println("terminal.MakeRaw error!", err)
			return session, false
		}
		defer terminal.Restore(fileDescriptor, originalState)

		termWidth, termHeight, err := terminal.GetSize(fileDescriptor)

		err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
		if err != nil {
			fmt.Println("RequestPty error!", err)
			return session, false
		}
	}
	return session, true
}

type Config struct {
	Server   string
	Username string
	Key      string
}

func sendCommands(config Config, cmdsfile string) {
	connected := false
	var client *ssh.Client
	client, connected = connect(client, config.Server, config.Username, config.Key)
	if !connected {
		return
	}
	defer client.Close()
	fmt.Println("Connected to " + config.Server + "!")
	var session *ssh.Session
	session, sessOk := makeSessionWithTerminal(client, session)
	session.Stdin, _ = os.OpenFile(cmdsfile, os.O_RDONLY, 0777)
	//session.Stdin = strings.NewReader(commandToAllSevers)
	defer session.Close()
	if !sessOk {
		fmt.Println("Terminal error!")
		return
	}
	err := session.Shell()
	if err != nil {
		fmt.Println("Shell error!", err)
		return
	}
	session.Wait()
	fmt.Println("Disconnected from " + config.Server + "!")
}

func test_gorutine() {
	/*	file, err := os.Create("file_for_test_gorutine.txt")

		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		file.WriteString("mkdir " + strconv.Itoa(time.Now().Nanosecond()/1000) + "\r\nexit")
		file.Close()*/

	connected := false
	var client *ssh.Client
	client, connected = connect(client, "lab.posevin.com:9907", "-", "-")
	if !connected {
		return
	}
	defer client.Close()
	//fmt.Println("Connected to " + config.Server + "!")
	var session *ssh.Session
	session, sessOk := makeSessionWithTerminal(client, session)
	//session.Stdin, _ = os.OpenFile("file_for_test_gorutine.txt", os.O_RDONLY, 0777)
	session.Stdin = strings.NewReader("mkdir ttt ");
	//session.Stdin = strings.NewReader(commandToAllSevers)
	defer session.Close()
	if !sessOk {
		fmt.Println("Terminal error!")
		return
	}
	err_ := session.Shell()
	if err_ != nil {
		fmt.Println("Shell error!", err_)
		return
	}
	session.Wait()
	//fmt.Println("Disconnected from " + config.Server + "!")
}

func main() {
	close := false
	var commandToClient string
	for !close {
		var client *ssh.Client
		fmt.Print("mySSH> ")
		fmt.Fscan(os.Stdin, &commandToClient)
		switch commandToClient {
		case "test_gorutine":
			{
				for i := 0; i < 1; i++ {
					go test_gorutine()
				}
			}
		case "connect":
			{
				var adr, login, password string
				fmt.Fscan(os.Stdin, &adr)
				fmt.Print("Enter login: ")
				fmt.Fscan(os.Stdin, &login)
				fmt.Print("Enter password: ")
				fmt.Fscan(os.Stdin, &password)
				connected := false
				client, connected = connect(client, adr, login, password)
				if !connected {
					break
				}
				defer client.Close()
				fmt.Println("Connected!")
				var session *ssh.Session
				session, sessOk := makeSessionWithTerminal(client, session)
				defer session.Close()
				if !sessOk {
					fmt.Println("Terminal error!")
					break
				}
				err := session.Shell()
				if err != nil {
					fmt.Println("Shell error!", err)
					break
				}
				session.Wait()
				fmt.Println("Disconnected!")
			}
		case "multiconnect":
			{
				var infofile string
				fmt.Fscan(os.Stdin, &infofile)
				var cmdsfile string
				fmt.Fscan(os.Stdin, &cmdsfile)
				configs := make([]Config, 0)

				configFile, err := ioutil.ReadFile(infofile)

				if err != nil {
					fmt.Println("Read infofile error!")
				}

				configLines := strings.Split(string(configFile), "\n")

				for i := 0; i < len(configLines); i++ {

					if configLines[i] != "" {

						configLine := strings.Split(string(configLines[i]), " ")

						newConfig := Config{Server: configLine[0], Username: configLine[1], Key: configLine[2]}
						configs = append(configs, newConfig)
					}
				}

				for _, config := range configs {
					go sendCommands(config, cmdsfile)
				}
			}
		case "close":
			{
				close = true
				break
			}
		}
	}
	fmt.Println("Client closed!")
	return
}
