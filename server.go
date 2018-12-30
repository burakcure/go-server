package main

/*
Here is a global server which can be used in all applications with a little modification.
Default packet type is:
XXXXXX,YYYYYY,ZZZZ,....,YYYY, (String)
You can store command on XXXXX with its properties

*/
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

var wu []User //Variable which contains inside of the JSON file.
var log []string

//A struct which contains users.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Mapno    int    `json:"mapno"`
}

func handleCommands(c net.Conn, id string, password string) { //Users commit commands in this function.

	for {
		//Users commits command in this place.

		netData, err := bufio.NewReader(c).ReadString('\n') //Waits for command.
		if err != nil {
			fmt.Println(err)
			return
		}
		commandNotParsed := strings.TrimSpace(string(netData))
		command := strings.Split(commandNotParsed, ",")                                                                                                 //Parses input.
		log = append(log, (string(time.Now().Format(time.RFC850) + "	|	" + c.RemoteAddr().String() + "		|	" + id + "	:	" + commandNotParsed + "\r\n"))) //Taking log of command.
		fmt.Printf(command[0])
		if command[0] == "log" {
			//Command using example
			c.Write([]byte(string("You have already logged in.\n")))
			fmt.Printf("Serving to %s\n", c.RemoteAddr().String())
		} else if command[0] == "reg" {

			c.Write([]byte(string("You have already logged in.\n")))

		} else {

			c.Write([]byte(string("Undefined command.\n")))
		}

	}

}
func handleConnection(c net.Conn) {
	var logged = false
	fmt.Printf("Serving %s\n", c.RemoteAddr().String()) //Server prints out the ip of the connection.

	idpw, err := bufio.NewReader(c).ReadString('\n') //First of all it reads id and password to login and take the commands.
	if err != nil {
		fmt.Println(err)
		return
	}

	idpw = strings.TrimSpace(string(idpw))
	temporaryarray := strings.Split(idpw, ",") // Splits the received data so we can inspect the command.
	fmt.Println(temporaryarray[1])
	log = append(log, (string(time.Now().Format(time.RFC850) + "	|	" + c.RemoteAddr().String() + "		|			:	" + idpw + "\r\n")))
	//Read packet if user wants to login or register
	if temporaryarray[0] == "log" { //If our parsed received datas first element is log, it takes id and password and checks.

		for i := range wu {
			//Checks id and password if its available.
			if wu[i].Username == temporaryarray[1] { //If username is available.
				if wu[i].Password == temporaryarray[2] { //If password is available.
					fmt.Println("Logged in as ", temporaryarray[1])
					//User logged in

					logged = true

					c.Write([]byte("Logged in as " + temporaryarray[1] + " \n"))
					handleCommands(c, temporaryarray[1], temporaryarray[2])

				} else {

					fmt.Println(c.RemoteAddr().String(), " tried to breach.")
					c.Write([]byte(string("Wrong password")))
					c.Close()
					break
				}

			}
			//If no username found.

		}
		if logged == false {
			fmt.Println(temporaryarray[1], " no username found.") //If no username founds.
			c.Write([]byte(string("Wrong username")))             //Sends Wrong username to client.

			c.Close() //Closes the socket.}
		}
	} else if temporaryarray[0] == "reg" { //
		var usercheck bool

		usercheck = false
		for i := range wu { //Checks if username is already taken(Non case sensitive)
			if wu[i].Username == temporaryarray[1] {
				c.Write([]byte("Username is not available.\n"))
				usercheck = true
				c.Close()
				break
			}
		}
		if usercheck == false { //If username is available.
			var tempuser User //Creates temporary array which will contain users information to save.

			tempuser.Username = temporaryarray[1] //Adds information of user to tempuser User array.
			tempuser.Password = temporaryarray[2]
			tempuser.Mapno = 2
			fmt.Println(temporaryarray)
			wu = append(wu, tempuser)
			c.Write([]byte("Register successful \n"))
			handleCommands(c, temporaryarray[1], temporaryarray[2])
			//Need to warp into loop that we take commands from user. Instead of that it closes
		}

	} else { //If command is not log or reg Server sends Undefined command and closes the connection.
		c.Write([]byte("Undefined command \n"))
		c.Close()
	}

}
func main() {

	jsonFile, err := os.Open("users.json") //Reads database.
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close() //Closes the file after finishing read.

	byteValue, _ := ioutil.ReadAll(jsonFile) // Reads the bytes values of file.

	json.Unmarshal(byteValue, &wu) //Unmarshal the byte array.
	ServerStart()
}

func ServerStart() {

	PORT := ":7001"                    //Port which server waits for connection.
	l, err := net.Listen("tcp4", PORT) //Starts listening.
	if err != nil {
		fmt.Println(err)
		return
	}

	go SaveDatabase() //Saves database for every 5 seconds.
	//Connection accept
	for {
		c, err := l.Accept() //Accepts the connections.
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c) //Handles the connections.
	}

}
func SaveDatabase() {
	//This function saves our changed data to our JSON file.
	var logStringsave string

	for {

		time.Sleep(5000 * time.Millisecond) //Saves every 5 seconds to reduce the waste of resources
		if len(wu) > 0 {                    //Checks if wu is not empty and if its not save it.
			deliver, _ := json.Marshal(wu)
			ioutil.WriteFile("users.json", deliver, 0644) //Stores users as JSON at users.json
		}
		if len(log) > 0 { //Checks if wu is not empty and if its not save it.
			logStringsave = strings.Join(log, "\n")

			f, _ := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) //Opens logs.txt to write logs. If logs.txt is not exist, creates a new one.

			f.Write([]byte(logStringsave))
			f.Close()
			log = log[:0]
		}
	}

}
