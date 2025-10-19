package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
)

type NetAddr struct {
	Host string
	Port int
}

var Addr = &NetAddr{
	Host: "localhost",
	Port: 8080,
}

func (a *NetAddr) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddr) Set(adr string) error {
	argValues := strings.Split(adr, ":")

	if len(argValues) != 2 {
		return errors.New("address not validate")
	}

	port, err := strconv.Atoi(argValues[1])
	if err != nil {
		return errors.New("address not validate")
	}
	a.Port = port
	a.Host = argValues[0]
	return nil
}

var (
	ResultURL string
	LevelLog  string
)

func ParseFlags() error {

	_ = flag.Value(Addr)
	flag.Var(Addr, "a", "Net address host:port")
	flag.StringVar(&ResultURL, "b", "http://localhost:8080", "result url")
	flag.StringVar(&LevelLog, "l", "info", "level logger")
	flag.Parse()

	if envRunBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		ResultURL = envRunBaseURL
	}
	if logLevel, ok := os.LookupEnv("LEVEL_LOG"); ok {
		LevelLog = logLevel
	}

	if envRunAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		err := Addr.Set(envRunAddr)
		if err != nil {
			return err
		}
	}
	return nil

}
