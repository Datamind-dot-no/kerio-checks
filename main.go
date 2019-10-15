package main

import (
	"github.com/datamind-dot-no/kerio-checks/config"

	"github.com/datamind-dot-no/kerio-checks/app/notifications"
	"github.com/datamind-dot-no/kerio-checks/app/qcheck"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	notifications := notifications.New(conf)
	qchk := qcheck.New(conf, notifications)

	err = qchk.CheckQ()
	if err != nil {
		panic(err)
	}
}
