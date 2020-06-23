package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/shishkebaber/message-api/server"
	"github.com/shishkebaber/message-api/storage"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var bindAddress = os.Getenv("BIND_ADDRESS")
	stor := storage.NewMStorage()
	router := mux.NewRouter().StrictSlash(true)
	s := server.NewServer(stor, router)
	httpSrv := &http.Server{
		Handler: s,
		Addr:    bindAddress,
	}
	go func() {
		s.Logger.Info("Server starting")
		err := httpSrv.ListenAndServe()
		if err != nil {
			s.Logger.Error("Error during starting the server", err)
			os.Exit(1)
		}
	}()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	signal.Notify(exitChan, os.Kill)
	sig := <-exitChan
	s.Logger.Println("Got signal: ", sig)

	s.Logger.Info("Server stopping")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	httpSrv.Shutdown(ctx)
}
