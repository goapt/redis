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

func (u *user) DbName() string {
	return "default"
}

func (u *user) TableName() string {
	return "user"
}

func (u *user) PK() string {
	return "id"
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
