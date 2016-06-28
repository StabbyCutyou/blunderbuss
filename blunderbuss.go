// Package main is the entrypoint for blunderbuss
// This was generated via go-servicescaffolder Version 0.0.1 on 2016-06-08 22:55:25 -0400
package main

import (
	"log"
	"os"

	"github.com/StabbyCutyou/blunderbuss/boot"
	"github.com/StabbyCutyou/blunderbuss/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const Version = "0.0.1"

func main() {
	// TODO incorporate extracted logging solution?
	log.SetOutput(os.Stdout)
	log.Printf("Blunderbuss v%s starting up...", Version)

	bp, err := boot.Boot()
	if err != nil {
		log.Fatal(err)
	}

	// This will only return if something interrupts it
	log.Fatal(bp.HTTPServer.Listen())
}

func openDB(cfg *config.Config) (*sqlx.DB, error) {
	return sqlx.Open("postgres", cfg.DBConnString)
}
