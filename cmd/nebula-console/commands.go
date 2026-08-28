package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gesture-nebula-console/domain"
	"gesture-nebula-console/store"
)

func printRecords(w io.Writer, records []domain.Record) error {
	return json.NewEncoder(w).Encode(records)
}

func listCommand(path string, w io.Writer) error {
	s, err := store.Open(path)
	if err != nil {
		return err
	}
	defer s.Close()
	records, err := s.ListRecords()
	if err != nil {
		return err
	}
	return printRecords(w, records)
}

func healthCommand(w io.Writer) { _, _ = fmt.Fprintln(w, "gesture-nebula ok") }

func runCommand(args []string, w io.Writer) error {
	if len(args) == 0 {
		healthCommand(w)
		return nil
	}
	if args[0] == "list" {
		path := dataPath()
		if len(args) > 1 {
			path = args[1]
		}
		return listCommand(path, w)
	}
	if args[0] == "health" {
		healthCommand(w)
		return nil
	}
	if args[0] == "backup" {
		path := dataPath()
		if len(args) > 1 {
			path = args[1]
		}
		s, err := store.Open(path)
		if err != nil {
			return err
		}
		defer s.Close()
		backup, err := s.ExportBackup()
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(backup)
	}
	return fmt.Errorf("unknown command %s", args[0])
}

func commandMain() int {
	if err := runCommand(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
