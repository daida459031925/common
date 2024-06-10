package netService

import (
	"fmt"
	"github.com/daida459031925/common/net/netService"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestService(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	server := netService.NewServer(":8080", mux)

	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 等待信号
	<-sigChan

	// 停止服务器
	if err := server.Stop(); err != nil {
		log.Printf("Failed to stop server: %v\n", err)
		os.Exit(1)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}
