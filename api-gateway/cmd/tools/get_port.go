package main

import (
	"fmt"
	"github.com/semho/hotel-booking/api-gateway/internal/config"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		os.Exit(1)
	}
	fmt.Print(cfg.HTTP.Port)
}
