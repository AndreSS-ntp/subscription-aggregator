package config

import "os"

var IP_port = os.Getenv("IP_PORT")
var DB_URL = os.Getenv("DB_URL")
