package parser

import (
	"io/ioutil"
	"strings"
	"time"
	"log"
	"os/user"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"os/exec"
)


func Parse() {
	defer exec.Command("killall", "chrome")
	login,pass:=getCred()
	u := launcher.New().
        Set("headless").
        Delete("--headless").
        MustLaunch()
	page:=rod.New().ControlURL(u).MustConnect().MustPage("https://sg.lyceum130.ru")
	time.Sleep(5*time.Second)
	page.MustElement("select#schools").MustSelect("МАОУ Лицей №130")
	page.MustElement("input[name=\"UN\"]").MustInput(login)
	page.MustElement("input[name=\"PW\"]").MustInput(pass)
	page.MustElement("a.button-login").MustClick()
	time.Sleep(1*time.Second)
	if page.MustHas("button[title=\"Продолжить\"]") {
		page.MustElement("button[title=\"Продолжить\"]").MustClick()
	}
	log.Println("Logged in...")
	adver,_:=page.MustWaitLoad().Elements("div.advertisement")

	for _,n:=range adver {
		if strings.Contains(n.MustElement("h3").MustText(), "Изменения в расписании") {
			log.Println(n.MustElement("h3").MustText())
			n.MustElement("div.fieldset").MustElement("a").MustClick()
		}
	}
	time.Sleep(2*time.Second)
	err:=page.Close()
	if err!=nil {
		log.Println(err)
	}
}

func getCred() (string, string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := usr.HomeDir + "/.config/sgobot/cred.config"
	content,err:=ioutil.ReadFile(path)
	if err!=nil {
		log.Fatal(err)
	}

	return strings.Split(string(content), " ")[0], strings.Split(string(content), " ")[1]
}