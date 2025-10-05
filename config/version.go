package config

import (
	"fmt"
)

var (
	appName = "monitor-etl"
	version = "dev"
	date    = "unknown"
)

func ShowVersion() {
	fmt.Printf("%s %s, built at %s\n", appName, version, date)
}
