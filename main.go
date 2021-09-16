package main

import (
	"fmt"
	"log"

	"github.com/rgynn/pensionera/pkg/config"
	"github.com/rgynn/pensionera/pkg/symbol"
)

func main() {

	config, err := config.NewFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	svc, err := symbol.NewService(config.Count, config.APIKey, config.SymbolNames)
	if err != nil {
		log.Fatal(err)
	}
	defer svc.Close()

	for _, name := range config.SymbolNames {
		fmt.Printf("Waiting for atleast %d results to calculate SMA for: %s\n", config.Count, name)
	}

	if err := svc.Run(); err != nil {
		log.Fatal(err)
	}
}
