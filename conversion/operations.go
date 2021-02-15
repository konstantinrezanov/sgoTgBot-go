package conversion

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

func Format() {
	path:=downloadDir()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}

	cur,_:=os.Getwd()
	for _, name := range names {
		if strings.Contains(name, ".xlsx") {
			os.Rename(path+"/"+name, cur+"/data"+"/"+strings.Split(name, " ")[0]+".xlsx")
		}
	}
}
func downloadDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := usr.HomeDir + "/Downloads"

	return path
}