package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "db.db?_fk=true")
	if err != nil {
		log.Fatal(err)
	}
	startServer(db)
	// e := Employee {
	// 	Eid:      100,
	// 	Name:     "hemu",
	// 	Email: "hemu@hemu.hemu",
	// 	PhoneNo: "989898998",
	// 	Projects: []Project{Project{Pid: 100, Name: "bla project"}},
	// 	WorkExp:  []WorkExperience{WorkExperience{CompanyName: "adfa"}},
	// 	Skills:    []Skill{Skill{"bullshitting"}},
	// }
	// InsertEmployee(&ctx, e)
	// fmt.Println("inserted, eid is: ", ctx.eid)
	// ee := GetEmployee(&ctx)
	// if ctx.err != nil {
	// 	log.Fatal(ctx.err)
	// }
	// fmt.Print(ee)
}

