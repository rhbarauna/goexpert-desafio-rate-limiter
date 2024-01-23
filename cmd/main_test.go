package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

func TestLoadTest_IP_WebServer(t *testing.T) {
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/",
	})

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics

	rate := vegeta.Rate{Freq: 5000, Per: time.Second}
	duration := 29 * time.Second

	for res := range attacker.Attack(targeter, rate, duration, "Test Load Test") {
		metrics.Add(res)
	}

	metrics.Close()
}

func TestLoadTest_TOKEN_WebServer(t *testing.T) {

	header := http.Header{}
	header.Add("API_KEY", "tkn_123")

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/",
		Header: header,
	})

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics

	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 29 * time.Second

	for res := range attacker.Attack(targeter, rate, duration, "Test Load Test") {
		metrics.Add(res)
	}

	metrics.Close()
}

func TestMain(m *testing.M) {
	// Iniciar servidor web em uma goroutine
	go main()

	// Aguardar um curto período para garantir que o servidor esteja pronto
	time.Sleep(1 * time.Second)

	// Executar os testes e sair
	exitCode := m.Run()

	// Parar o servidor após os testes
	os.Exit(exitCode)
}
