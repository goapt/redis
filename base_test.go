package redis

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	configs := make(map[string]Config)

	configs["default"] = Config{
		Server: ":6379",
	}

	Connect(configs)
	m.Run()
}

type user struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func TestBaseRedis_HGetAll(t *testing.T) {
	client := NewBaseRedis("default")

	u := &user{
		Id:1,
		Name:"test",
		CreatedAt:time.Now(),
	}

	err := client.HMSet("test",u)

	if err != nil {
		t.Errorf("hmset error:%s",err)
	}

	u2 := &user{}

	err = client.HGetAll("test",u2)

	if err != nil {
		t.Errorf("hgetall error:%s",err)
	}

	fmt.Println(u2)
}

type user2 struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

func TestBaseRedis_HGetAllString(t *testing.T) {
	client := NewBaseRedis("default")

	u := &user2{
		Id:1,
		Name:"test",
		CreatedAt:time.Now().Format("2006-01-02 15:04:05"),
	}

	err := client.HMSet("test",u)

	if err != nil {
		t.Errorf("hmset error:%s",err)
	}

	u2 := &user2{}

	err = client.HGetAll("test",u2)

	if err != nil {
		t.Errorf("hgetall error:%s",err)
	}

	fmt.Println(u2)
}
