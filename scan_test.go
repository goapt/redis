package redis

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type t_users struct {
	UserId     int       `json:"user_id" db:"user_id"`
	UserName   string    `json:"user_name" db:"user_name"`
	Status     int       `json:"status" db:"status"`
	Timezone   string    `json:"timezone" db:"timezone"`
	Lang       string    `json:"lang" db:"lang"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

type t_orders struct {
	OrderId    int          `json:"order_id" db:"order_id"`
	Body       string       `json:"body" db:"body"`
	ExtAttr    OrderExtAttr `json:"ext_attr" db:"ext_attr"` // 扩展属性存放
	CreateTime time.Time    `json:"create_time" db:"create_time"`
	UpdateTime time.Time    `json:"update_time" db:"update_time"`
}

type OrderExtAttr struct {
	NotifyUrl string `json:"notify_url,omitempty"` // 回调通知url
}

var emptyJSON = json.RawMessage("{}")

func JsonObject(value interface{}) (json.RawMessage, error) {
	var source []byte
	switch t := value.(type) {
	case string:
		source = []byte(t)
	case []byte:
		source = t
	case nil:
		source = emptyJSON
	default:
		return nil, errors.New("incompatible type for json.RawMessage")
	}

	if len(source) == 0 {
		source = emptyJSON
	}

	return source, nil
}

// Scan implements the Scanner interface.
func (n *OrderExtAttr) Scan(value interface{}) error {
	b, err := JsonObject(value)
	if err != nil {
		return err
	}
	// 忽略非对象形式的数据脏数据，比如：数组
	if len(b) > 0 && []byte(b)[0] == '[' {
		b = json.RawMessage("{}")
	}
	return json.Unmarshal(b, n)
}

// Value implements the driver Valuer interface.
func (n OrderExtAttr) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *OrderExtAttr) UnmarshalBinary(data []byte) error {
	// convert data to yours, let's assume its json data
	return json.Unmarshal(data, n)
}

func (n OrderExtAttr) MarshalBinary() ([]byte, error) {
	return json.Marshal(n)
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
		user := &t_users{}
		scanStruct(m, user)
	}
}

func TestScanStructOrder(t *testing.T) {
	order := &t_orders{}

	m := map[string]string{
		"order_id":    "123123",
		"body":        "test",
		"ext_attr":    `{"notify_url":"https://v.com"}`,
		"create_time": "2015-03-18 18:20:28",
		"update_time": "2017-09-20 10:29:59",
	}

	err := scanStruct(m, order)
	assert.NoError(t, err)
	assert.Equal(t, 123123, order.OrderId)
	assert.Equal(t, "2015-03-18 18:20:28", order.CreateTime.Format("2006-01-02 15:04:05"))
	assert.Equal(t, "https://v.com", order.ExtAttr.NotifyUrl)
}

func TestScanStruct(t *testing.T) {
	user := &t_users{}
	err := scanStruct(m, user)

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
	user := &t_users{
		CreateTime: time.Date(2019, 11, 02, 15, 04, 05, 0, time.UTC),
	}
	m := structToMapInterface(user)
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
