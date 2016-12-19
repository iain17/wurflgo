package main

import (
	"github.com/srinathgs/wurflgo"
	"fmt"
)

func main() {
	repository := wurflgo.Read("wurfl.gob")
	if repository == nil {
		repository = wurflgo.New("wurfl.xml", "product_info,xhtml_ui")
		error := repository.Save("wurfl.gob")
		if error != nil {
			fmt.Println(error)
		}
	}
	fmt.Println("Finished loading")
	device := repository.Match("Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")
	if device != nil {
		fmt.Println(device.Properties)
	} else {
		fmt.Println("Device not found!")
	}
}