package main

import (
	"hellofudango/hellofudan"
)

func main() {
	hf := hellofudan.NewManager([]hellofudan.Student{
		{
			StudentID: "",
			Password:  "",
		},
	})

	hf.Start()
}
