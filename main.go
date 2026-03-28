package main

import (
	"fmt"
	"log"

	"github.com/mpdev25/pokedexcli/blog_aggregator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	fmt.Printf("Read config %+v\n", c)
	err = c.SetUser("Mike")
	if err != nil {
		log.Fatalf("Couldn't set user name: %v", err)
	}
	c, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	fmt.Printf("Read config again: %+v\n", c)
}

//currentConfig, err := config.Read()
//if err != nil {
//	fmt.Println("Error reading config:", err)
//	return
//}
//err = currentConfig.SetUser("Mike")
//if err != nil {
//	fmt.Println("Error setting user:", err)
//	return
//}
//updatedConfig, err := config.Read()
//if err != nil {
//	fmt.Println(" Error reading config:", err)
//	return
//}
//fmt.Println(updatedConfig)

//}
