package main

import (
	"io/ioutil"
	"net/http"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Gothmog"
	app.Action = Main
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "rules",
			Usage:    "Path to a JSON file containing exported Balrog Rules",
			Required: true,
		},
	}
	app.RunAndExitOnError()
}

func Main(c *cli.Context) error {
	rawRules, err := ioutil.ReadFile(c.String("rules"))
	if err != nil {
		return err
	}
	loadedRules, err := parseRules(rawRules)
	if err != nil {
		return err
	}
	gothmogHandler := &GothmogHandler{
		rules: &loadedRules,
	}

	mux := http.NewServeMux()
	// We don't have any fixed endpoints; each update request will send data in a specific format
	// that gothmogHandler will parse into usable data.
	mux.Handle("/", gothmogHandler)

	server := &http.Server{
		Addr:    "127.0.0.1:8888",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
