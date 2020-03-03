package redis

import (
	"fmt"
	"testing"
	"time"
)

type tUsers struct {
	UserId     int       `json:"user_id" db:"user_id"`
	UserName   string    `json:"user_name" db:"user_name"`
	Status     int       `json:"status" db:"status"`
	Timezone   string    `json:"timezone" db:"timezone"`
	Lang       string    `json:"lang" db:"lang"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

var m = map[string]string{
	"user_id":     "4",
	"user_name":   "test",
	"status":      "1",
	"timezone":    "Asia/Shanghai",
	"lang":        "zh-CN",
	"create_time": "2015-03-18 18:20:28",
	"update_time": "2017-09-20 10:29:59",
}

func BenchmarkScanStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := &tUsers{}
		ScanStruct(m, user)
	}
}

func TestScanStruct(t *testing.T) {
	user := &tUsers{}
	err := ScanStruct(m, user)

	if err != nil {
		t.Error(err)
	}

	if user.UserId != 4 {
		t.Error("Parse user_id error:", user.UserId)
	}

	if user.CreateTime.Format("2006-01-02 15:04:05") != "2015-03-18 18:20:28" {
		t.Error("Parse create_time error:", user.CreateTime)
	}

}

func Test_structToMapInterface(t *testing.T) {
	user := &tUsers{
		CreateTime: time.Date(2019, 11, 02, 15, 04, 05, 0, time.UTC),
	}
	m := structToMapInterface(user)

	fmt.Println(user, m)
	var now2 string

	if m2, has := m["create_time"]; has {
		if m3, ok := m2.(string); ok {
			now2 = m3
		}
	}

	if now2 != "2019-11-02 15:04:05" {
		t.Error("struct to map time parse error")
	}
}
