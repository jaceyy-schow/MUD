package main
import (
	"strings"
	"fmt"

)
var itod = []string{"north", "east", "west", "south", "up", "down"}

func doLook(command string, p *Player){
	cmd := strings.Fields(command)
	if len(cmd) <= 1 {
		ReadRoom(p.Location.ID, AllRooms, p)
	}
	if len(cmd) > 1{
		if strings.HasPrefix(cmd[1], "n"){
			index := dtoi["n"]
			ReadRoom(p.Location.Exits[index].To.ID, AllRooms, p)
		}
		if strings.HasPrefix(cmd[1], "s"){
			index := dtoi["s"]
			ReadRoom(p.Location.Exits[index].To.ID, AllRooms, p)
		}
		if strings.HasPrefix(cmd[1], "e"){
			index := dtoi["e"]
			ReadRoom(p.Location.Exits[index].To.ID, AllRooms, p)
		}
		if strings.HasPrefix(cmd[1], "w"){
			index := dtoi["w"]
			ReadRoom(p.Location.Exits[index].To.ID, AllRooms, p)
		}
		if strings.HasPrefix(cmd[1], "u"){
			index := dtoi["u"]
			ReadRoom(p.Location.Exits[index].To.ID, AllRooms, p)
		}
		if strings.HasPrefix(cmd[1], "d"){
			index := dtoi["d"]
			ReadRoom(p.Location.Exits[index].To.ID, AllRooms, p)
		}
	}
}

func move(input string, p *Player){
	if strings.HasPrefix(input, "n"){
		index := dtoi["n"]
		if p.Location.Exits[index].To == nil{
			p.Printf("You can't move that way")
			ReadDir(p.Location.ID, AllRooms, p)
		}else{
			AnnounceExit(p)
			p.Location = p.Location.Exits[index].To
			AnnounceEnter(p)
			ReadRoom(p.Location.ID, AllRooms, p)
		}
		//p.Printf("I'm moving North")
	}
	if strings.HasPrefix(input, "s"){
		index := dtoi["s"]
		
		if p.Location.Exits[index].To == nil{
			p.Printf("You can't move that way")
			ReadDir(p.Location.ID, AllRooms, p)
		}else{
			AnnounceExit(p)
			p.Location = p.Location.Exits[index].To
			AnnounceEnter(p)
			ReadRoom(p.Location.ID, AllRooms, p)
		}
		//fmt.Println("I'm moving south")
	}
	if strings.HasPrefix(input, "e"){
		index := dtoi["e"]
		
		if p.Location.Exits[index].To == nil{
			p.Printf("You can't move that way")
			ReadDir(p.Location.ID, AllRooms, p)
		}else{
			AnnounceExit(p)
			p.Location = p.Location.Exits[index].To
			AnnounceEnter(p)
			ReadRoom(p.Location.ID, AllRooms, p)
		}
		//fmt.Println("I'm moving east")
	}
	if strings.HasPrefix(input, "w"){
		index := dtoi["w"]
		
		if p.Location.Exits[index].To == nil{
			p.Printf("You can't move that way")
			ReadDir(p.Location.ID, AllRooms, p)
		}else{
			AnnounceExit(p)
			p.Location = p.Location.Exits[index].To
			AnnounceEnter(p)
			ReadRoom(p.Location.ID, AllRooms, p)
		}
		//fmt.Println("I'm moving west")
	}
	if strings.HasPrefix(input, "u"){
		index := dtoi["u"]
		if p.Location.Exits[index].To == nil{
			p.Printf("You can't move that way")
			ReadDir(p.Location.ID, AllRooms, p)
		}else{
			AnnounceExit(p)
			p.Location = p.Location.Exits[index].To
			AnnounceEnter(p)
			ReadRoom(p.Location.ID, AllRooms, p)
		}
		//fmt.Println("I'm Flying")
	}
	if strings.HasPrefix(input, "d"){
		index := dtoi["d"]
		
		if p.Location.Exits[index].To == nil{
			p.Printf("You can't move that way")
			ReadDir(p.Location.ID, AllRooms, p)
		}else{
			AnnounceExit(p)
			p.Location = p.Location.Exits[index].To
			AnnounceEnter(p)
			ReadRoom(p.Location.ID, AllRooms, p)
		}
		//fmt.Println("Digging a hole underground")
	}
}
func recall(s string, p *Player){
	p.Printf("With The click of your heels three times you begin to whisper. 'There's no place like home.' Nothing Happens, you begin the long journey back. Why must I have to walk everywhere oh great one.\n")
	p.Location = AllRooms[3001]
	ReadRoom(p.Location.ID, AllRooms, p)
}

func smile(s string, p *Player){
	for _,item := range ActivePlayers{
		if item.Location == p.Location{
			p.Printf(p.Name + ": releases a glistening smoulder to the world.\n")
		}
	}
	ReadDir(p.Location.ID, AllRooms, p)
}

func cry(s string, p *Player){
	for _,item := range ActivePlayers{
		if item.Location == p.Location{
			p.Printf(p.Name + ": begins to sob undontrolably. Why must you make my life this hard my god.\n")
		}
	}
	ReadDir(p.Location.ID, AllRooms, p)
}
func eat(s string, p *Player){
	p.Printf("You whip out the weird meat you found on the road earlier. Not sure what this is but I'm sure it will taste fine with some garlic.\n")
	ReadDir(p.Location.ID, AllRooms, p)
}
func pray(s string, p *Player){
	p.Printf("Oh god, please make something good come my way today\n")
	ReadDir(p.Location.ID, AllRooms, p)
}
func curse(s string, p *Player){
	for _,item := range ActivePlayers{
		if item.Location == p.Location{
			p.Printf(p.Name + ": Well gosh diddly darn. Why musn't I be able to use my native language here. This place fricken sucks\n")
		}
	}
	ReadDir(p.Location.ID, AllRooms, p)
}
func quit(_ string, p *Player){
	p.Printf("Goodbye!\n")
	close(p.Outputs)
	fmt.Println("Closed channel for ", p.Name )
	p.Outputs = nil
	
	delete(ActivePlayers, p.Name)

	for _,item := range ActivePlayers{
		if item.Location == p.Location{
			item.Printf(p.Name + " has left the game\n")
		}
	}
}

func gossip(message string, p *Player){
	for _,item := range ActivePlayers{
		item.Printf(p.Name + " gossips: " + strings.TrimPrefix(message, "gossip")+ "\n")
	}
}

func say(message string, p *Player) {
	for _, item := range ActivePlayers {
		if item.Location == p.Location {
			item.Printf(p.Name + " says:" + strings.TrimPrefix(message, "say") + "\n")
		}
	}
}
func tell(message string, p *Player) {
	field := strings.Fields(message)
	for _, item := range ActivePlayers {
		if item.Name == field[1] {
			item.Printf(p.Name + " tells you:" + strings.TrimPrefix(message, "tell "+item.Name) + "\n")
			return
		}
	}
}

func shove(message string, p *Player){
	field := strings.Fields(message)
	for _,item := range ActivePlayers {
		if item.Name == field[1]{
			item.Printf(p.Name + " shoved you pretty freaking hard. What a jerk.\n")
		}
	}
}

func shout(message string, p *Player) {
	for _, item := range ActivePlayers {
		if item.Location.Zone == p.Location.Zone {
			item.Printf(p.Name + " shouts:" + strings.TrimPrefix(message, "shout") + "\n")
		}
	}
}

func where(message string, p *Player) {
	if len(ActivePlayers) == 0{
		p.Printf("Nobody else is online.")
	}

	for _, item := range ActivePlayers {
		if item != p && item.Location.Zone == p.Location.Zone {
			p.Printf(item.Name + ":" + item.Location.Name + "\n")
		}
	}

}
func AnnounceEnter(p *Player){
	for _,item := range ActivePlayers{
		if item.Location == p.Location && item != p{
			item.Printf(p.Name + " has entered the room.\n")
		}
	}
}
func AnnounceExit(p *Player){
	for _,item := range ActivePlayers{
		if item.Location == p.Location && item != p{
			item.Printf(p.Name + " has exited the room.\n")
		}
	}
}
func help(message string, p *Player){
	p.Printf("You can say directions to get which directions you can move to.")
	p.Printf("Try using (n, s, e, w, u, d) to move")
	p.Printf("gossip speaks to all active players")
	p.Printf("shout yells to everyone in your zone")
	p.Printf("say will go to everyone in your room")
	p.Printf("Where will tell you where all active players are located")
	p.Printf("If you say 'tell' and someones name your message will go directly to them")
	p.Printf("Try saying cry, smile, eat, pray, or curse for fun!")
	p.Printf("Saying recall will take you back to the temple and quit will exit the game")
}
func directions(message string, p *Player){
	ReadDir(p.Location.ID, AllRooms, p)
}


func HandleCommand(s string, p *Player){
	//command := strings.Fields(s)
	if  verb, ok := Commands[s]; ok {
		//do something here
		verb(s, p)
	} else{
	p.Printf("Huh? Try Again")
	}
}

func AddCommand(command string, action func(string, *Player)){
	for index, _ := range command{
		//fmt.Println(command[:index+1])
		new_command := command[:index +1]
		Commands[new_command] = action
	}
}
func AddAllCommands(){
	AddCommand("recall", recall)
	AddCommand("smile", smile)
	AddCommand("cry", cry)
	AddCommand("eat", eat)
	AddCommand("pray", pray)
	AddCommand("curse", curse)
	AddCommand("shove", shove)
	AddCommand("gossip", gossip)
	AddCommand("say", say)
	AddCommand("tell", tell)
	AddCommand("shout", shout)
	AddCommand("where", where)
	AddCommand("help", help)
	AddCommand("directions", directions)
	AddCommand("look north", doLook)
	AddCommand("look south", doLook)
	AddCommand("look east", doLook)
	AddCommand("look west", doLook)
	AddCommand("look up", doLook)
	AddCommand("look down", doLook)
	AddCommand("north", move)
	AddCommand("south", move)
	AddCommand("east", move)
	AddCommand("west", move)
	AddCommand("up", move)
	AddCommand("down", move)
	AddCommand("quit", quit)
	
}

func ReadRoom(id int, AllRooms map[int]*Room, p *Player){
	p.Printf(AllRooms[id].Name + "\n")
	p.Printf(AllRooms[id].Description)
	p.Printf("\nDirections: [ ")
	for index, element := range AllRooms[id].Exits{
		if element.To != nil{
			p.Printf(itod[index] + " ")
		}
	}
	p.Printf(" ] \n")
	if len(ActivePlayers) > 1{
		for _,item := range ActivePlayers{
			if item.Location == AllRooms[id] && item != p{
				if len(ActivePlayers)<3{
					p.Printf(item.Name)
				}else{
					p.Printf(item.Name + ", ")
				}
			}
		}
		p.Printf(" is in the room with you.\n")
	}
}

func ReadDir(id int, AllRooms map[int]*Room, p *Player){
	p.Printf("\nDirections: [ ")
	for index, element := range AllRooms[id].Exits{
		if element.To != nil{
			p.Printf(itod[index] + " ")
		}
	}
	p.Printf(" ] \n")

}

