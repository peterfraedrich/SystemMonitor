package main

import (
	"os"

	"github.com/akamensky/argparse"
)

var CONFIG *Config

func parseArgs() (*Config, error) {
	parser := argparse.NewParser("SystemMonitor", "Monitors system stuff and writes it to a sqlite DB for troubleshooting issues")
	p := parser.String("p", "path", &argparse.Options{Required: false, Default: "/usr/share/sysmon.db", Help: "Path to the SQLIte database file"})
	d := parser.Flag("d", "drop", &argparse.Options{Required: false, Help: "Drop the database when restarting the app; useful for dev/debugging SystemMonitor"})
	f := parser.Int("f", "frequency", &argparse.Options{Required: false, Default: 60, Help: "Frequency (in seconds) to poll for metrics and info"})
	s := parser.Flag("s", "stdout", &argparse.Options{Required: false, Default: false, Help: "Log everything to stdout"})
	err := parser.Parse(os.Args)
	if err != nil {
		return nil, err
	}
	return &Config{
		DBPAth:      *p,
		DropDB:      *d,
		Frequency:   *f,
		LogToStdout: *s,
	}, nil
}

func main() {
	var err error
	CONFIG, err = parseArgs()
	if err != nil {
		panic(err)
	}
	db, err := OpenDB(CONFIG.DBPAth, CONFIG.DropDB)
	if err != nil {
		panic(err)
	}
	svc := NewService(db)
	svc.Start()
}
