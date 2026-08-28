package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"gesture-nebula-console/api"
	"gesture-nebula-console/domain"
	"gesture-nebula-console/notifier"
	"gesture-nebula-console/service"
	"gesture-nebula-console/store"
)

func main() {
	cfg := LoadConfig()
	path := flag.String("db", cfg.DatabasePath, "bolt database path")
	addr := flag.String("addr", cfg.Address, "http listen address")
	flag.Parse()
	if len(flag.Args()) > 0 {
		os.Exit(commandMain())
	}
	cfg.DatabasePath, cfg.Address = *path, *addr
	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}
	s, err := store.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	if cfg.Seed {
		seedUser(s)
	}
	svc := service.New(s)
	primary := &notifier.MemorySender{}
	router := notifier.NewRouter(primary)
	_ = router.AddRule("urgent", primary)
	n := notifier.New(s, router)
	server := api.NewServer(svc, n)
	log.Printf("gesture nebula listening on %s", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		log.Print(err)
	}
}

func seedUser(s *store.Store) { _ = s.SaveUser(defaultUser()) }
func defaultUser() domain.User {
	return domain.User{ID: "operator", Name: "Console Operator", Role: "operator", Active: true}
}
func dataPath() string {
	if value := os.Getenv("NEBULA_DB"); value != "" {
		return value
	}
	return "nebula.db"
}
