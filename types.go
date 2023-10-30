package main

type Config struct {
	LogToStdout bool   `json:"log_to_stdout" yaml:"log_to_stdout"` // log captured metrics to STDOUT
	DBPAth      string `json:"db_path" yaml:"db_path"`             // path to sqlite database
	DropDB      bool   `json:"drop_db" yaml:"drop_db"`             // drop the database when connecting to it
	Frequency   int    `json:"frequency" yaml:"frequency"`         // check frequency, in seconds
}
