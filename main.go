package main

import (
	"fmt"

	"github.com/mpdev25/pokedexcli/blog_aggregator/internal/config"
)

func main() {
	currentConfig, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	err = currentConfig.SetUser("Mike")
	if err != nil {
		fmt.Println("Error setting user:", err)
		return
	}
	updatedConfig, err := config.Read()
	if err != nil {
		fmt.Println(" Error reading config:", err)
		return
	}
	fmt.Println(updatedConfig)

}
