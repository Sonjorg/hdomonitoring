package database

import (
	"database/sql"
	"log"
	//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	//_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type RoutingE struct {
	Id        int
	Ipaddress string
	RoutingTable string
	RoutingEntry string
	 //map consisting of routingtables and their routingentries
}
type RoutingT struct {
	Id        int
	Ipaddress string
	Time      string
	RoutingTables string
	 //map consisting of routingtables and their routingentries
}
/*
type RoutingE struct {
	Time      string
	RoutingTable string
	RoutingEntry string
}*/
/*
type RoutingTmp struct {
	Id        int
	Ipaddress string
	Time      string
	RoutingTablesnEntries map[string][]string
	//RoutingEntries []string
}*/
func CreateRoutingSqlite(db * sql.DB) error{
	createRoutingTables := `CREATE TABLE IF NOT EXISTS routingtables (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"ipaddress" TEXT,
		"time" TEXT,
		"routingtable" TEXT
		);`

	createRoutingEntries := `CREATE TABLE IF NOT EXISTS routingentries (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"ipaddress" TEXT,
		"routingtable" TEXT,
		"routingentries" TEXT
		);`

	statement, err := db.Prepare(createRoutingTables) // Prepare SQL Statement
	if err != nil {
		return err
	}
	statement2, err := db.Prepare(createRoutingEntries) // Prepare SQL Statement
	if err != nil {
		return err
	}
	statement.Exec()
	statement2.Exec()
	return nil
}

func StoreRoutingTables(db *sql.DB, ipaddress string, time string, routingTables []string) error{
	log.Println("Inserting tables ...")
	for i := range routingTables {
		insertSQL1 := `INSERT INTO routingtables(ipaddress, time, routingtable) VALUES (?, ?, ?)`
		statement, err := db.Prepare(insertSQL1) // Prepare statement.
													// This is good to avoid SQL injections
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = statement.Exec(ipaddress, time, routingTables[i])
		if err != nil {
			fmt.Println(err)
			return err
		}
}
	return nil
}

func StoreRoutingEntries(db *sql.DB, ipaddress string, routingTable string, routingEntries []string) error{
	log.Println("Inserting entries ...")
	for i := range routingEntries {
		insertSQL1 := `INSERT INTO routingentries(ipaddress, routingtable, routingentries) VALUES (?, ?, ?)`

		statement, err := db.Prepare(insertSQL1) // Prepare statement.
													// This is good to avoid SQL injections
		if err != nil {
			fmt.Println(err)

			return err

		}
		_, err = statement.Exec(ipaddress, routingTable, routingEntries[i])
		if err != nil {
			fmt.Println(err)
			return err
		}
}
	return nil
}

func RoutingEntriesExists(db * sql.DB,ipaddress string) bool {
   // sqlStmt := `SELECT ipaddress FROM routingtables WHERE ipaddress = ?`
	sqlStmt := `SELECT ipaddress FROM routingentries WHERE ipaddress = ?`
    err := db.QueryRow(sqlStmt, ipaddress).Scan()
    if err != nil {
        if err != sql.ErrNoRows {
            // a real error happened! you should change your function return
            // to "(bool, error)" and return "false, err" here
            return false
        }

        return false
    }

    return true
}
func RoutingTablesExists(db * sql.DB,ipaddress string) bool {
	// sqlStmt := `SELECT ipaddress FROM routingtables WHERE ipaddress = ?`
	 sqlStmt := `SELECT ipaddress FROM routingtables WHERE ipaddress = ?`
	 err := db.QueryRow(sqlStmt,ipaddress).Scan()
	 if err != nil {
		 if err != sql.ErrNoRows {
			 // a real error happened! you should change your function return
			 // to "(bool, error)" and return "false, err" here
			 return false
		 }

		 return false
	 }

	 return true
 }

func GetRoutingTables(db *sql.DB,ipaddress string) ([]string, error) {

	//if (routingTablesExists(db,ipaddress)) {
		//row, err := db.Query("SELECT * FROM routingtables")
		row, err := db.Query(`SELECT * FROM routingtables WHERE ipaddress = ?`, ipaddress)
		//row.Scan(ip)
		if err != nil {
			return nil, err
			//fmt.Println(err)
		}

		defer row.Close()
		/*err = row.QueryRow(ipaddress).Scan(&Id, &Ipaddress, &Time, &RoutingTable, &RoutingEntry)
		if err != nil {
      	  log.Println(err)
    	}*/
		var rt []string
		//var data []*RoutingT
		for row.Next() {
			r := &RoutingT{}
				if err := row.Scan(&r.Id, &r.Ipaddress,&r.Time,&r.RoutingTables); err != nil{
					fmt.Println(err)
				}
					//data = append(data, r)
				rt = append(rt, r.RoutingTables)
		}
		return rt ,err
}
func GetRoutingEntries(db *sql.DB,ipaddress string,routingTable string) ([]string, error) {

	//if (routingTablesExists(db,ipaddress)) {
		//row, err := db.Query("SELECT * FROM routingtables")
		row, err := db.Query(`SELECT * FROM routingentries WHERE ipaddress = ?`, ipaddress)
		//row.Scan(ip)
		if err != nil {
			return nil, err
			//fmt.Println(err)
		}

		defer row.Close()
		/*err = row.QueryRow(ipaddress).Scan(&Id, &Ipaddress, &Time, &RoutingTable, &RoutingEntry)
		if err != nil {
      	  log.Println(err)
    	}*/
		var re []string
		//var data []*RoutingT
		for row.Next() {
			r := &RoutingE{}
				if err := row.Scan(&r.Id, &r.Ipaddress,&r.RoutingTable, &r.RoutingEntry); err != nil{
					fmt.Println(err)
				}
				if (r.Ipaddress == ipaddress) {
					//data = append(data, r)
					if (r.RoutingTable == routingTable) {
						re = append(re, r.RoutingEntry)
					}
				}
		}
		return re ,err
}

/*
func main() {

	var sqliteDatabase *sql.DB

				sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db")
				if err != nil {
					fmt.Println(err)
				}
	var s []string
	s = append(s, "1")
	s = append(s, "2")
	s = append(s, "3")


	createRoutingSqlite(sqliteDatabase)
	storeRoutingEntries(sqliteDatabase, "ipadresse", "time","5", s)
	if (routingTablesExists(sqliteDatabase, "ipadresse")) {
		g, err := getRoutingEntries(sqliteDatabase,"ipadresse","5")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(g)

	}
}
*/