package tests

import (
	"testing"

	"github.com/G0tem/go-servise-auth/internal/config"
)

func failOnError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf("%s: %s", msg, err)
	}
}

func TestRabbitMqConnect(t *testing.T) {
	cfg := config.LoadConfig()

	var err error
	failOnError(t, err, "Failed to declare a queue")
	if cfg.HttpPort == 0 {
		t.Errorf("Wrong http port")
	}

}
