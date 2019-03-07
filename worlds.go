package main
import (
    "database/sql"
    _"github.com/mattn/go-sqlite3"
    "fmt"
    "strconv"
    //"os"
    //"context"

)
var AllRooms = make(map[int]*Room)
var Players = make(map[string]Info)
var Zones map[int]*Zone
var dtoi = map[string]int{
    "n": 0,
    "e": 1,
    "w": 2,
    "s": 3,
    "u": 4,
    "d": 5,
}

type Zone struct {
    ID    int 
    Name  string
    Rooms []*Room
}

type Room struct {
    ID          int 
    Zone        *Zone
    Name        string
    Description string
    Exits       [6]Exit
}

type Exit struct {
    To          *Room
    Description string
}

type Player struct {
    Location    *Room
    Outputs     chan OutputEvent
    Name        string

}
type Info struct{
    Name        string
    Salt        string
    Hash        string
}


func ReadRooms(db *sql.DB, Zones map[int]*Zone) map[int]*Room{
    tx, err := db.Begin()
    if err != nil{
        fmt.Println(err)
        tx.Rollback()
    }

    rows, err := tx.Query("SELECT * FROM rooms ORDER BY id")

    if err != nil{
        fmt.Println(err)
        tx.Rollback()
        
    }
    defer rows.Close()
    //AllRooms := make(map[int]*Room)
    for rows.Next(){
        r := &Room{}
        var zone_id int
    
        if err := rows.Scan(
            &r.ID, 
            &zone_id, 
            &r.Name, 
            &r.Description,
        ); 
        err != nil{
            fmt.Println(err)
            tx.Rollback()
            
        }
        
        r.Zone = Zones[zone_id]
        AllRooms[r.ID] = r
    }
    if err != nil{
        fmt.Println(err)
        tx.Rollback()
        
    }
    tx.Commit()
    return AllRooms
}

func ReadZones(db *sql.DB){
    tx, err := db.Begin()
    Zones = make(map[int]*Zone)
    if err != nil{
        fmt.Println(err)
        tx.Rollback()
    }

    rows, err := tx.Query("SELECT * FROM zones ORDER BY id")

    if err != nil{
        fmt.Println(err)
        tx.Rollback()
        
    }
    defer rows.Close()
    for rows.Next(){
        z := &Zone{}
        if err := rows.Scan(
            &z.ID,
            &z.Name,
        ); 
        err != nil{
            fmt.Println(err)
            tx.Rollback()
            
        }
        Zones[z.ID] = z
        //fmt.Println("z: ", z.ID)
    }
    if err != nil{
        fmt.Println(err)
        tx.Rollback()
        
    }
    //return allZones
    tx.Commit()
    
}

func ReadExits(db *sql.DB){
    tx, err := db.Begin()
    if err != nil{
        fmt.Println(err)
        tx.Rollback()
    }

    rows, err := tx.Query("SELECT * FROM exits ORDER BY from_room_id")

    if err != nil{
        fmt.Println(err)
        tx.Rollback()
    }
    defer rows.Close()
    for rows.Next(){
        var e Exit
        var idFrom int
        var idTo int
        var direction string
        var description string
        if err := rows.Scan(
            &idFrom,
            &idTo,
            &direction,
            &description,
        ); 
        err != nil{
            fmt.Println(err)
            tx.Rollback()
            
        }
        e.To = AllRooms[idTo]
        e.Description = description
        //fmt.Println(AllRooms[idFrom])
        AllRooms[idFrom].Exits[dtoi[direction]] = e
    }
    if err != nil{
        fmt.Println(err)
        tx.Rollback()
        
    }
    tx.Commit()
    
}

func ReadPlayers(db *sql.DB)map[string]Info{
    tx, err := db.Begin()
    if err != nil{
        fmt.Println(err)
        tx.Rollback()
    }
    Players = make(map[string]Info)
    p, err := tx.Query("SELECT * FROM players")
    if err != nil{
        fmt.Println(err)
    }
    for p.Next(){
        var id int
        var name string
        var salt string
        var hash string
        p.Scan(&id, &name, &salt, &hash)
        Players[name] = Info{name,salt,hash}
    }
    return Players
    
}

func (p*Player) Printf(format string, a ...interface{}){
    msg := fmt.Sprintf(format, a...)
    p.Outputs <- OutputEvent{p, msg}
}

func makePlayer(tx *sql.Tx, name string, salt string, hash string, numPlayers string){
    id, err := strconv.Atoi(numPlayers)
    if err != nil{
        fmt.Println(err)
    }
    id +=1
    newID := strconv.Itoa(id)
    query:= "INSERT INTO players VALUES("+ newID+ ", " + name + ", "+ salt + ", " + hash

    tx.Exec(query)
    Players[name] = Info{name,salt,hash}

}

