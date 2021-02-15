package conversion

import (
	"sync"
	"time"
)

func (schedule *Schedule) ScheduleClear(rang int, m *sync.Mutex) {
	for index,day:=range schedule.days {
		if day.date.Before(time.Now().AddDate(0, 0, -rang)) {
			fastRemove(&schedule.days, index)	
		}
	}
}

func fastRemove(slice *[]Day,index int) {
	(*slice)[index]=(*slice)[len(*slice)-1]
	(*slice)[len(*slice)-1]=Day{}
	*slice=(*slice)[:len(*slice)-1]
}