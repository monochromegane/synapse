package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/monochromegane/synapse"
)

type profiler struct {
	db *sql.DB
}

func (p *profiler) Profile(ctx synapse.Context) (synapse.Profile, error) {
	profile := synapse.Profile{}

	var name string
	var birth string
	err := p.db.QueryRow("SELECT name, birth FROM users WHERE id = ?", ctx["user_id"]).Scan(&name, &birth)
	if err != nil {
		return profile, err
	}
	profile["name"] = name
	profile["birth"] = birth
	return profile, nil
}

func (p *profiler) Initialize() error {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/synapse")
	if err != nil {
		return err
	}
	p.db = db
	return nil
}

func (p *profiler) Finalize() error {
	p.db.Close()
	return nil
}

func NewProfiler() synapse.Profiler {
	return &profiler{}
}
