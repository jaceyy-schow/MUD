package main

import (
	"database/sql"
    _"github.com/mattn/go-sqlite3"
	"bufio"
	"fmt"
	"os"
	"net"
	"strings"
	"encoding/base64"
	crand "crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"strconv"
	"golang.org/x/crypto/pbkdf2"
	//"context"
)
const(
	CONN_HOST = "localhost"
	CONN_PORT = "9001"
	CONN_TYPE = "tcp"
)
type InputEvent struct{
    player *Player
    command []string
    Login   bool
}
type OutputEvent struct{
    player *Player
    Text    string
}
var Commands map[string]func(string, *Player)
var ActivePlayers map[string]*Player

func main() {
	db, err := sql.Open("sqlite3", "./world.db")
    if err != nil {
        fmt.Println(err)
	}
	readWorlds(db)
	//scanCommands(player)
	in := make(chan InputEvent)
	go listenForConnections(in)
	for event := range in{
		if event.player.Outputs != nil{
			if event.Login{
				createPlayer, err := db.Begin()
				makePlayer(createPlayer, event.command[0], event.command[1], event.command[2], event.command[3])
				if err != nil{
					_=createPlayer.Rollback()
					fmt.Println("failed to add to db")
				}
				createPlayer.Commit()
				fmt.Println("New Player", event.player.Name)
			}else {
				Commands[event.command[0]](strings.Join(event.command, " "), event.player)
			}
		}
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
func listenForConnections(in chan InputEvent){
	ActivePlayers = make(map[string]*Player)
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	checkError(err)
	fmt.Println("Listening on " + CONN_HOST + ":"+ CONN_PORT)
	for{
		conn, err := l.Accept()
		checkError(err)
		go func(){
			fmt.Println("New Conn Found: ", conn.RemoteAddr())
			scanner := bufio.NewScanner(conn)
			var name, pass string

			fmt.Fprintf(conn, "Username: \n>")
			for scanner.Scan(){
				name = strings.Trim(scanner.Text()," ")
				if name == ""{
					fmt.Fprintf(conn, "Username cannot be empty")
					conn.Close()
					fmt.Println("Connection closed for ", conn.RemoteAddr())
					return
				}
				break
			}
			fmt.Fprintf(conn, "Password: \n>")
			for scanner.Scan(){
				pass = scanner.Text()
				if len(pass) < 6{
					fmt.Fprintf(conn, "Password must contain more than 6 characters")
					conn.Close()
					fmt.Println("Connection closed for ", conn.RemoteAddr())
					return
				}
				break
			}
			var makeSalt bool
			if item, ok := Players[name]; ok{
				s, err := base64.StdEncoding.DecodeString(item.Salt)
				checkError(err)
				hash2, err2 := base64.StdEncoding.DecodeString(item.Hash)
				checkError(err2)
				enteredPass := pbkdf2.Key([]byte(pass), s, 64*1024, 32, sha256.New)
				if subtle.ConstantTimeCompare(hash2, enteredPass) != 1{
					fmt.Fprintln(conn, "Incorrect Password.")
					fmt.Println("Incorrect login for", name)
					conn.Close()
					fmt.Println("Connection closed for ", conn.RemoteAddr())
					return
					
				} else{
					fmt.Fprintln(conn, "Welcome back " + name)
					makeSalt = false
				}
			} else{
				fmt.Fprintf(conn, "Welcome New Player!\n")
				makeSalt = true
			}
			output := make(chan OutputEvent, 100)
			p := Player{AllRooms[3001], output, name,}
			AnnounceEnter(&p)
			if makeSalt {
				salt := make([]byte, 32)
				_, err := crand.Read(salt)
				checkError(err)
				salt64 := base64.StdEncoding.EncodeToString(salt)

				hash1 := pbkdf2.Key([]byte(pass), salt, 64*1024, 32, sha256.New)

				hash164 := base64.StdEncoding.EncodeToString(hash1)

				in <- InputEvent{&p, []string{name, salt64, hash164, strconv.Itoa(len(Players))}, true}
			}
			if item, ok := ActivePlayers[name]; ok{
				p.Location = item.Location

				fmt.Println(p.Name, "has joined from a new location: ", conn.RemoteAddr())
				close(item.Outputs)

				item.Outputs = nil
			}
			ActivePlayers[name] = &p
			go handleRequest(in, conn, &p)
		}()
	}
}
func handleRequest(in chan InputEvent, conn net.Conn, p *Player){
	scanner := bufio.NewScanner(conn)
	fmt.Println("Joined game from ", conn.RemoteAddr())
	ReadRoom(p.Location.ID, AllRooms, p)
	go func(){
		for event := range p.Outputs{
			fmt.Fprintf(conn, event.Text)
		}
		conn.Close()
		fmt.Println("Connection closed at ", conn.RemoteAddr())
		return
	}()
	for scanner.Scan(){
		line := scanner.Text()
		lineList := strings.Fields(line)
		if len(lineList) == 0{
			fmt.Fprintln(conn, "You must enter a command")
		} else if _, ok := Commands[lineList[0]]; ok == true{
			inputEvent := InputEvent{p, lineList, false}
			in <- inputEvent
		} else{
			p.Printf("Huh?\n")
		}
	}
	fmt.Println("(Goroutine closing) End of Scanner for ", conn.RemoteAddr())
	exitgame := InputEvent{p, strings.Fields("quit"), false}
	in <- exitgame
	return
}
func readWorlds(db *sql.DB){
	ReadZones(db)
	ReadRooms(db, Zones)
	ReadExits(db)
	ReadPlayers(db)
	Commands = make(map[string]func(string, *Player))
	AddAllCommands()
}
/*func scanCommands(p *Player){
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		//fmt.Println(scanner.Text()) // Println will add back the final '\n'
		s := scanner.Text()
		HandleCommand(s, p)
		//fmt.Printf("Fields are: %q", strings.Fields(s))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}*/



