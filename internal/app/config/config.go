package config

import (
	"flag"
	"os"
)

var Params struct {
	ServerAddress  string
	ShortUrlPrefix string
}

func init() {
	flag.StringVar(&Params.ServerAddress, "a", "localhost:8080", "Server address like localhost:8080")
	flag.StringVar(&Params.ShortUrlPrefix, "b", "http://localhost:8080", "Prefix of short url like (localhost:8080)/url_code")

	for _, v := range []struct {
		from string
		to   *string
	}{
		{"SERVER_ADDRESS", &Params.ServerAddress},
		{"BASE_URL", &Params.ShortUrlPrefix},
	} {
		if envValue := os.Getenv(v.from); envValue != "" {
			*v.to = envValue
		}
	}
}
