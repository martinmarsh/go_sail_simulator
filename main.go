/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/
package main

import (
	"go_sail_simulator/cmd"
	"go_sail_simulator/nmea"

) 

func main() {
	nmea.Setup()
	cmd.Execute()
}
