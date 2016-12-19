package main

import (
	"github.com/iain17/wurflgo"
	"fmt"
)

func main() {
	repository := wurflgo.Read("wurfl.gob")
	if repository == nil {
		repository = wurflgo.New("wurfl.xml", "product_info")
		error := repository.Save("wurfl.gob")
		if error != nil {
			fmt.Println(error)
		}
	}
	fmt.Println("Finished loading")
	device := repository.Match("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.95 Safari/537.36")
	if device != nil {
		fmt.Println(device.Properties)
		fmt.Println(device.Children)
	} else {
		fmt.Println("Device not found!")
	}
}