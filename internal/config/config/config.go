package config

import (
	"errors"
	"flag"
	"fmt"
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

func RemoveIndSlice(data *[]string, s int) []string {
	return append((*data)[:s], (*data)[s+1:]...)
}

type FilePathType struct {
	FileName   string
	Directorys []string
}

var FilePath = &FilePathType{FileName: "short-url-db.json", Directorys: []string{"temp"}}

func (f *FilePathType) String() (ds string) {
	for _, v := range f.Directorys {
		ds += fmt.Sprintf("%s%s", v, "/")
	}
	if ds == "" {
		return f.FileName
	}
	return ds + f.FileName
}

func (f *FilePathType) Set(a string) error {
	if a == "`" || a == "" || a == "``" || a == " " {
		*f = FilePathType{}
		return nil
	}

	*f = FilePathType{}

	argsPath := strings.Split(a, "/")
	for _, v := range argsPath {
		if v == "" || v == " " {
			continue
		}
		f.Directorys = append(f.Directorys, v)
	}
	f.FileName = f.Directorys[len(f.Directorys)-1]
	f.Directorys = RemoveIndSlice(&f.Directorys, len(f.Directorys)-1)
	return nil
}

func (f *FilePathType) Dir() (dirs string) {
	if len(f.Directorys) != 0 {
		for _, dir := range f.Directorys {
			dirs += fmt.Sprintf("%s%s", dir, "/")
		}
		return dirs
	}
	return ""
}

func (f *FilePathType) GetDir() string {
	return ""
}

var (
	ResultURL string
	LevelLog  string
)

func ParseFlags() error {

	_ = flag.Value(Addr)
	_ = flag.Value(FilePath)
	flag.Var(Addr, "a", "Net address host:port")
	flag.StringVar(&ResultURL, "b", "http://localhost:8080", "result url")
	flag.StringVar(&LevelLog, "l", "info", "level logger")
	flag.Var(FilePath, "f", "file path to save json")
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

	if envFilePath, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		err := FilePath.Set(envFilePath)
		if err != nil {
			return err
		}
	}
	return nil

}
