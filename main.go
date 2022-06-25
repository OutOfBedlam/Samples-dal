package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alecthomas/kong"
	"github.com/upper/db/v4/adapter/postgresql"
)

func main() {
	var cli struct {
		CheckConn struct {
			Url string `short:"U" default:"" help:""`
		} `cmd:"" help:"connect test"`
		CreateSchema struct{} `cmd:"" help:"create tables"`
		Collections  struct{} `cmd:"" help:"list collections"`
		Select       struct{} `cmd:"" help:"select table"`
	}

	cmd := kong.Parse(&cli)
	switch cmd.Command() {
	case "check-conn":
		// -U "postgres://my-user:my-secret@127.0.0.1:50000/my-db"
		settings, err := Settings(cli.CheckConn.Url)
		if err != nil {
			log.Fatal("Settings: ", err)
		}

		sess, err := postgresql.Open(settings)
		if err != nil {
			log.Fatal("Open: ", err)
		}
		defer sess.Close()

		driver := sess.Driver().(*sql.DB)
		t1 := time.Now()
		err = driver.Ping()
		if err != nil {
			log.Fatal("Ping: ", err)
		}
		fmt.Printf("Ping: %s\n", time.Since(t1))

		fmt.Printf("DSN: %q\n", settings)
	case "create-schema":
		var err error
		settings, _ := Settings("")
		sess, _ := postgresql.Open(settings)
		defer sess.Close()
		driver := sess.Driver().(*sql.DB)
		_, err = driver.Exec(`DROP TABLE IF EXISTS "birthday";`)
		if err != nil {
			log.Fatalf("drop table: %q\n", err)
		}

		_, err = driver.Exec(`CREATE TABLE "birthday" (name VARCHAR(50), born TIMESTAMP);`)
		if err != nil {
			log.Fatalf("create table: %q\n", err)
		}

		birthdayCollection := sess.Collection("birthday")
		err = birthdayCollection.Truncate()
		if err != nil {
			log.Fatalf("Truncate(): %q\n", err)
		}
		birthdayCollection.Insert(Birthday{
			Name: "Timothee Hal Chalamet",
			Born: time.Date(1995, time.December, 27, 0, 0, 0, 0, time.Local),
		})

		birthdayCollection.Insert(Birthday{
			Name: "Zendaya Maree Stoermer",
			Born: time.Date(1996, time.September, 1, 0, 0, 0, 0, time.Local),
		})

		birthdayCollection.Insert(Birthday{
			Name: "Rebecca Ferguson",
			Born: time.Date(1983, time.October, 19, 0, 0, 0, 0, time.Local),
		})

		res := birthdayCollection.Find()

		var birthdays []Birthday
		err = res.All(&birthdays)
		if err != nil {
			log.Fatalf("select(): %q\n", err)
		}
		for _, b := range birthdays {
			fmt.Printf("%s was born in %s\n", b.Name, b.Born.Format("January 2, 2006"))
		}
	case "select":
		settings, _ := Settings("")
		sess, _ := postgresql.Open(settings)
		defer sess.Close()

		birthdayCollection := sess.Collection("birthday")
		res := birthdayCollection.Find()

		var birthdays []Birthday
		if err := res.All(&birthdays); err != nil {
			log.Fatalf("select(): %q\n", err)
		}
		for _, b := range birthdays {
			fmt.Printf("%s was born in %s\n", b.Name, b.Born.Format("January 2, 2006"))
		}
	case "collections":
		settings, _ := Settings("")
		sess, _ := postgresql.Open(settings)
		defer sess.Close()

		collections, err := sess.Collections()
		if err != nil {
			log.Fatalf("collections(): %q\n", err)
		}
		for _, c := range collections {
			fmt.Printf("%s\n", c.Name())
		}
	default:
		panic(cmd.Command())
	}
}

func Settings(url string) (*postgresql.ConnectionURL, error) {
	var settings *postgresql.ConnectionURL
	var err error
	if len(url) > 0 {
		settings, err = postgresql.ParseURL(url)
	} else {
		settings = &postgresql.ConnectionURL{
			Host:     "127.0.0.1:50000",
			Database: "my-db",
			User:     "my-user",
			Password: "my-secret",
		}
	}
	return settings, err
}

type Birthday struct {
	Name string    `db:"name"`
	Born time.Time `db:"born"`
}

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

type Employee struct {
	ID     uint64 `db:"id,omitempty"`
	Person `db:",inline"`
}

type Author struct {
	ID     uint64 `db:"id,omitempty"`
	Person `db:",inline"`
}
