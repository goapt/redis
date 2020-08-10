package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	configs := make(map[string]Config)

	configs["default"] = Config{
		Server: testServer,
	}

	Connect(configs)
	m.Run()
}

type user struct {
	Id        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func TestRedis_SetNX(t *testing.T) {
	client := NewRedis(Client("default"))

	has, err := client.SetNX("setnxlock2", "1", 10*time.Second)

	if err != nil {
		t.Errorf("setnx error:%s", err)
	}

	if !has {
		t.Errorf("setnx must true")
	}

	has, err = client.SetNX("setnxlock2", "1", 10*time.Second)

	if err != nil {
		t.Errorf("setnx error:%s", err)
	}

	if has {
		t.Errorf("setnx must false")
	}
}

func TestRedis_HGetAll(t *testing.T) {
	client := NewRedis(Client("default"))

	testTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2012-10-24 07:06:00", time.Local)
	u := &user{
		Id:        1,
		Name:      "test",
		CreatedAt: testTime,
	}

	err := client.HMSet("test", u)

	if err != nil {
		t.Errorf("hmset error:%s", err)
	}

	u2 := &user{}

	err = client.HGetAll("test", u2)

	if err != nil {
		t.Errorf("hgetall error:%s", err)
	}

	if u2.CreatedAt.String() != testTime.String() {
		t.Errorf("hgetall time error, want get %s,but get %s", testTime, u2.CreatedAt)
	}
}

func TestRedis_HGetAllMap(t *testing.T) {
	client := NewRedis(Client("default"))
	data, err := client.HGetAllMap("nodata")

	if err == nil {
		t.Error("hgetall empty data must return error")
	}

	if len(data) != 0 {
		t.Errorf("hgetall error data length must 0, %s", err)
	}
}

type user2 struct {
	Id        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

func TestRedis_HGetAllString(t *testing.T) {
	client := NewRedis(Client("default"))
	testTime := "2012-10-24 07:06:00"
	u := &user2{
		Id:        1,
		Name:      "test",
		CreatedAt: "2012-10-24 07:06:00",
	}

	err := client.HMSet("test", u)

	if err != nil {
		t.Errorf("hmset error:%s", err)
	}

	u2 := &user2{}

	err = client.HGetAll("test", u2)

	if err != nil {
		t.Errorf("hgetall error:%s", err)
	}

	if u2.CreatedAt != testTime {
		t.Errorf("hgetall time error, want get %s,but get %s", testTime, u2.CreatedAt)
	}

	fmt.Println(u2)
}

func TestRedis_HMSet(t *testing.T) {
	client := NewRedis(Client("default"))

	err := client.HMSet("test", map[string]string{
		"a": "1",
		"b": "2",
	})

	if err != nil {
		t.Errorf("hmset error:%s", err)
	}

	v, err := client.HGet("test", "a")
	if err != nil {
		t.Errorf("hget error:%s", err)
	}

	if v != "1" {
		t.Errorf("hmset map[string]string want get:%s but get:%s", "1", v)
	}
}

func TestIsNil(t *testing.T) {
	client := NewRedis(Client("default"))

	_, err := client.Get("nil_test")
	assert.True(t, IsNil(err))
}
