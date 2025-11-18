package main

import (
	"log"
	"net/http"

	"github.com/cheetahbyte/centra/internal/api"
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/crypt"
	"github.com/cheetahbyte/centra/internal/helper"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	api.Register(r)

	port := config.GetPort()

	log.Printf("Centra API running on :%s\n", port)

	keyDir := config.GetKeysDir()

	pubKey, err := crypt.EnsureKeys(keyDir)
	if err != nil {
		log.Fatal("Startup failed: ", err)
	}

	helper.PrettyKey(pubKey)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
