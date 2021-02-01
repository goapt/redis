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
		Server:       testServer,
		MaxRetries:   3,
		DialTimeout:  2,
		ReadTimeout:  3,
		WriteTimeout: 4,
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

	{
		err := client.HSet("test", "a", "3")
		if err != nil {
			t.Errorf("hget error:%s", err)
		}
	}

	{
		v, err := client.HGet("test", "a")
		if err != nil {
			t.Errorf("hget error:%s", err)
		}

		if v != "3" {
			t.Errorf("hmset map[string]string want get:%s but get:%s", "1", v)
		}
	}

	{
		err := client.HDel("test", "a")
		assert.NoError(t, err)
	}

	{
		has, err := client.HExists("test", "a")
		assert.NoError(t, err)
		assert.False(t, has)
	}

	{
		v, err := client.HIncrBy("test", "b", 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), v)
	}
}

func TestIsNil(t *testing.T) {
	client := NewRedis(Client("default"))

	_, err := client.Get("nil_test")
	assert.True(t, IsNil(err))
}

func TestRedis_Exists(t *testing.T) {
	client := NewRedis(Client("default"))

	has := client.Exists("no_test")
	assert.False(t, has)
}

func TestRedis_List(t *testing.T) {
	client := NewRedis(Client("default"))

	err := client.Del("test_list")
	assert.NoError(t, err)

	has, err := client.LPush("test_list", 1, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), has)

	{
		v, err := client.LPop("test_list")
		assert.NoError(t, err)
		assert.Equal(t, "3", v)
	}

	{
		v, err := client.RPush("test_list", 4)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), v)
	}

	{
		v, err := client.RPop("test_list")
		assert.NoError(t, err)
		assert.Equal(t, "4", v)
	}

}

func TestRedis_Set(t *testing.T) {
	client := NewRedis(Client("default"))

	err := client.Del("test_set")
	assert.NoError(t, err)

	err = client.Set("test_set", "10")
	assert.NoError(t, err)
	{
		v, err := client.Get("test_set")
		assert.NoError(t, err)
		assert.Equal(t, v, "10")
	}

	{
		v, err := client.Incr("test_set")
		assert.NoError(t, err)
		assert.Equal(t, int64(11), v)
	}

	{
		v, err := client.Decr("test_set")
		assert.NoError(t, err)
		assert.Equal(t, int64(10), v)
	}

	{
		v, err := client.IncrBy("test_set", 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(12), v)
	}

	{
		v, err := client.DecrBy("test_set", 3)
		assert.NoError(t, err)
		assert.Equal(t, int64(9), v)
	}

	{
		v, err := client.Expire("test_set", 10*time.Second)
		assert.NoError(t, err)
		assert.True(t, v)
	}

	{
		v := client.TTL("test_set")
		assert.NoError(t, err)
		assert.True(t, v.Seconds() > 8)
	}

}

func TestRedis_SAdd(t *testing.T) {
	client := NewRedis(Client("default"))

	err := client.Del("test_sadd")
	assert.NoError(t, err)

	{
		v, err := client.SAdd("test_sadd", 1, 2, 3)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), v)
	}
}

func TestRedis_SetEX(t *testing.T) {
	client := NewRedisWithName("default")

	err := client.Del("test_set_ex")
	assert.NoError(t, err)

	{
		err := client.SetEX("test_set_ex", "1", 10*time.Second)
		assert.NoError(t, err)
	}

	{
		v := client.TTL("test_set_ex")
		assert.NoError(t, err)
		assert.True(t, v.Seconds() > 8)
	}
}
