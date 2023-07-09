package main

import (
	"fmt"
	"github.com/kataras/pg"
	"golang.org/x/net/context"
	"log"
	"os"
	"time"
)

type Event struct {
	ID          string    `json:"id" pg:"type=uuid,primary"`
	Title       string    `json:"title" pg:"type=varchar(255),unique"`
	Description *string   `json:"description" pg:"type=varchar(4096),null"`
	StartAt     time.Time `json:"startAt" pg:"type=TIMESTAMPTZ"`
	EndAt       time.Time `json:"endAt" pg:"type=TIMESTAMPTZ"`
	CrtdAt      time.Time `json:"crtdAt" pg:"type=TIMESTAMPTZ,default=clock_timestamp()"`
	UpdAt       time.Time `json:"updAt" pg:"type=TIMESTAMPTZ,default=clock_timestamp()"`
}

type ConnectionString struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Schema   string `yaml:"Schema"`
	DBName   string `yaml:"DBName"`
	SSLMode  string `yaml:"SSLMode"`
}

func main() {
	// Create Schema instance.
	schema := pg.NewSchema()
	// First argument is the table name, second is the struct entity.
	schema.MustRegister("events", Event{})

	opts := ConnectionString{
		Host:     "half-canary-4432.g8z.cockroachlabs.cloud",
		Port:     26257,
		User:     os.Getenv("COCKROACH_USER"),
		Password: os.Getenv("COCKROACH_PASSWORD"),
		DBName:   "godevops",
		Schema:   "public",
		SSLMode:  "verify-full",
	}
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		opts.User, opts.Password, opts.Host, opts.Port, opts.DBName, opts.SSLMode)
	db, err := pg.Open(context.Background(), schema, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.CreateSchema(context.Background()); err != nil {
		log.Fatal(err)
	}

	eventRepo := pg.NewRepository[Event](db)
	var newEvent = Event{}
	newEvent.ID = "101227e1-50c6-443c-9e54-69ee60ae1762"
	newEvent.Title = "Mastering Concurrency in Go"
	newEvent.EndAt = time.Now()
	except := []string{"Title", "EndAt", "UpdAt", "CrtdAt", "ID"}
	_, err = eventRepo.UpdateExceptColumns(context.Background(), except, newEvent)
	if err != nil {
		log.Fatal(err)
	}
}
