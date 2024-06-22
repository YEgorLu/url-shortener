package config

import (
	"flag"
)

var Params struct {
	ServerAddress  string
	ShortUrlPrefix string
}

func init() {
	flag.StringVar(&Params.ServerAddress, "a", "localhost:8080", "Server address like localhost:8080")
	flag.StringVar(&Params.ShortUrlPrefix, "b", "localhost:8080", "Prefix of short url like (localhost:8080)/url_code")
}
