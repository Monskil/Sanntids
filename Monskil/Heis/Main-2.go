package main

import "./Driver"

func main() {
	Driver.Elev_init()
	for{
		Driver.Floor_tracking()
	}
}