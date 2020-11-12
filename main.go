package main

import (
	"hellofudango/hellofudan"
)

func main() {
	hfm := hellofudan.NewManager([]hellofudan.Student{
		{
			StudentID: "",
			Password:  "",
		},
	})

	hfm.Start()
}
