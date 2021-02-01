package redis

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testServer = "127.0.0.1:6379"

func init() {
	if os.Getenv("DRONE") == "true" {
		testServer = "redis:6379"
	}
}

func TestConnectSentinel(t *testing.T) {
	if testing.Short() {
		t.Skip("not test")
	}

	configs := make(map[string]Config)

	configs["default"] = Config{
		MasterName:    "verystar",
		SentinelAddrs: "10.64.144.101:26379,10.64.144.102:26379,10.64.144.103:26379",
		Password:      "test",
	}

	Connect(configs)
	client := Client()
	assert.NotEqual(t, client, nil, "Redis connect error")

	ret := client.Set("test", "test", time.Second*100)
	assert.NoError(t, ret.Err())
}

func TestConnect(t *testing.T) {
	configs := make(map[string]Config)

	configs["default"] = Config{
		Server: testServer,
	}

	Connect(configs)
	client := Client()
	assert.NotEqual(t, client, nil, "Redis connect error")
}

func TestConnectError(t *testing.T) {

	assert.Panics(t, func() {
		configs := make(map[string]Config)

		configs["default"] = Config{
			Server: "111.111.111.111",
		}

		Connect(configs)
		client := Client()
		assert.NotEqual(t, client, nil, "Redis connect error")
	})

}
