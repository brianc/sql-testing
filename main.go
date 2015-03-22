package main

import (
	_ "database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
	"log"
	"strings"
	//"time"
)

var schema = `
CREATE TABLE IF NOT EXISTS peeps (id serial primary key, name text)
`

type Model struct {
	ID uint64 `json:id`
	//Created time.Time
}

type Person struct {
	Model
	Name string
}

func NewPerson(id uint64, name string) Person {
	person := Person{}
	person.ID = id
	person.Name = name
	return person
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := sqlx.Connect("postgres", "postgres://localhost")

	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	db.SetMaxOpenConns(20)
	db.MustExec(schema)

	person := Person{}
	person.Name = "Brian"
	_, err = db.NamedExec("INSERT INTO peeps(name) VALUES(:name)", person)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Databse connected: %T", person)
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM peeps")

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Number of rows: %d", count)

	peeps := []Person{}

	err = db.Select(&peeps, "SELECT * FROM peeps")
	if err != nil {
		log.Fatal(err)
	}

	for _, person := range peeps {
		log.Printf("Person: %d: %s", person.ID, person.Name)
	}
}
