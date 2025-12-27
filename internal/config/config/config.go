package config

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
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

type FilePathType struct {
	FilePath string
}

var (
	ResultURL string
	LevelLog  string
)

var DatabaseDSN string

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

var FilePath = &FilePathType{FilePath: "/tmp/short-url-db.json"}

func (f *FilePathType) String() (ds string) {
	return f.FilePath
}

func (f *FilePathType) Set(a string) error {
	if a != "" {
		f.FilePath = a
		return nil
	}
	return errors.New("writing/reading to a file is disabled")
}

func (f *FilePathType) Dir() (dirs string) {

	return filepath.Dir(f.FilePath)
}

func (f *FilePathType) IsEnabled() bool {
	if f.FilePath == "" || len(f.FilePath) < 5 {
		return false
	}
	return true
}

//host=localhost port=5432 user=postgres password=postgres dbname=short_url sslmode=disable

func ParseFlags() error {

	_ = flag.Value(Addr)

	flag.Var(Addr, "a", "Net address host:port")
	flag.StringVar(&ResultURL, "b", "http://localhost:8080", "result url")
	flag.StringVar(&LevelLog, "l", "info", "level logger")
	flag.StringVar(&DatabaseDSN, "d", "", "database adr")
	flag.Var(FilePath, "f", "file path to save json")
	flag.Parse()

	if envRunBaseURL, ok := os.LookupEnv("BASE_URL"); ok {
		ResultURL = envRunBaseURL
	}
	if logLevel, ok := os.LookupEnv("LEVEL_LOG"); ok {
		LevelLog = logLevel
	}
	if envDatabaseDSN, ok := os.LookupEnv("DATABASE_DSN"); ok {
		DatabaseDSN = envDatabaseDSN
	}

	if envRunAddr, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		err := Addr.Set(envRunAddr)
		if err != nil {
			return err
		}
	}

	if envFilePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		err := FilePath.Set(envFilePath)
		if err != nil {
			return err
		}
	}
	return nil

}
