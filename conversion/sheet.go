package conversion

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"sgoTgBot-go/parser"
	"sync"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type Schedule struct {
	days []Day
}

type Day struct {
	date    time.Time
	classes []Class
}

type Class struct {
	class   string
	lessons []Lesson
}

type Lesson struct {
	number int
	time   string
	title  string
}

var Classes = []string{"10А", "10Б", "10В", "10Г", "10Д"}

var LessonTime = map[int]string{
	1: "08:00",
	2: "08:50",
	3: "09:50",
	4: "10:50",
	5: "11:50",
	6: "12:40",
	7: "13:30",
	8: "14:30",
	9: "15:30",
}

func (schedule *Schedule) GetSchedule()  {
	var data Schedule
	files, err := ioutil.ReadDir("data/")
	if err != nil {
		log.Fatal(err)
	}

	var names []string
	for _, f := range files {
		date := f.Name()[0:(len(f.Name()) - 5)]
		names = append(names, date)
	}

	for _, day := range names {
		data.days = append(data.days, parseDay(day, Classes))
		os.Remove(fmt.Sprintf("data/%s.xlsx", day))
	}

	*schedule=data

	log.Println(*schedule)
}

func parseDay(name string, classes []string) Day {
	var data Day

	for _, class := range classes {
		data.classes = append(data.classes, parseClass(name, class))
	}

	name += fmt.Sprintf(".%v", time.Now().Year())
	date, _ := time.Parse("02.01.2006", name)
	data.date = date

	return data
}
func parseClass(name, class string) Class {

	var classData Class

	f, err := excelize.OpenFile(fmt.Sprintf("data/%s.xlsx", name))
	if err != nil {
		log.Println(err)
	}

	classData.class = class
	number, _ := strconv.Atoi(class[0 : len(class)-2])
	liter := strings.Split(class, "")[len(class)-2]

	titleBegin, err := f.SearchSheet("Лист1", class)
	if err != nil {

		log.Println(err)
	}
	titleEnd, err := f.SearchSheet("Лист1", fmt.Sprintf("%d%v", number+1, liter))
	if err != nil {

		log.Println(err)
	}
	begin, _ := strconv.Atoi(titleBegin[0][1:])
	end, _ := strconv.Atoi(titleEnd[0][1:])
	for i := begin + 1; i < end; i++ {
		ind := fmt.Sprintf("%c%d", titleBegin[0][0], i)
		less, _ := f.GetCellValue("Лист1", ind)
		if less != "" {
			time := LessonTime[i-begin]
			classData.lessons = append(classData.lessons, Lesson{number: i - begin, title: less, time: time})
		}
	}
	return classData
}

func (schedule Schedule) PrettyPrint(name, day string) string {
	var data string
	date, err := time.Parse("02.01.2006", day)
	if err != nil {
		log.Println(err)
	}

	data += fmt.Sprintln(name)
	data += fmt.Sprintln(day)

	for _, day := range schedule.days {
		if day.date == date {
			for _, class := range day.classes {
				if class.class == name {
					for _, lesson := range class.lessons {
						data += fmt.Sprintf("%d-%s: %s\n", lesson.number, lesson.time, lesson.title)
					}
				}
			}
		}
	}

	return data
}

const layout="02.01.2006"

func (schedule Schedule) AvailableDates() []string {
	var data []string

	for _,day:=range schedule.days {
		data = append(data, day.date.Format(layout))
	}

	return data
}

func (schedule *Schedule) AssignSchedule(m *sync.Mutex) {
	parser.Parse()
	Format()
	m.Lock()
	schedule.GetSchedule()
	m.Unlock()
}