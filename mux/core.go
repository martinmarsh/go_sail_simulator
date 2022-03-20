/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"time"
	"go_sail_simulator/nmea"


)

type ConfigData struct {
	Index    map[string]([]string)
	TypeList map[string]([]string)
	Values   map[string]map[string]([]string)
}

func Execute(config *ConfigData) {
	
	channels := make(map[string](chan string))
	fmt.Println("Go sail simulator execute")
	channels["command"] = make(chan string, 2)
	
	for name, param := range config.Index {
		for _, value := range param {
			if value == "outputs" {
				for _, chanName := range config.Values[name][value] {
					if _, ok := channels[chanName]; !ok {
						channels[chanName] = make(chan string, 30)
					}
				}
			}
			if value == "input" {
				for _, chanName := range config.Values[name][value] {
					if _, ok := channels[chanName]; !ok {
						channels[chanName] = make(chan string, 30)
					}
				}
			}
		}
	}

	
	for processType, names := range config.TypeList {
		fmt.Println(processType, names)
		for _, name := range names {
			switch processType {
			case "serial":
				serialProcess(name, config.Values[name], &channels)
			case "udp_client":
				udpClientProcess(name, config.Values[name], &channels)
			case "keyboard":
				keyBoardProcess(name, config.Values[name], &channels)
			case "ships_log":
				shipsLogProcess(name, config.Values[name], &channels)
			case "auto-helm":
				autoHelmProcess(name, config.Values[name], &channels)
					
			}
		}
	}

	go compass(&channels)
	
	for {
		command := <-(channels["command"])
		fmt.Printf("Command '%s' received\n", command)
		switch command {
		
		case "8":
			for i:=0 ; i <= 1000 ; i++{
				fmt.Println("Sending to 2000")
				channels["to_compass"] <- "$HCHDM,172.5,M*28"
				time.Sleep(100 * time.Millisecond)
				
			}	
		
		}
	}
}

func compass(channels *map[string](chan string)){
	nm := nmea.Sentences.MakeHandle()
	nm.Parse("$HCHDM,172.5,M*28")
	nm.Parse("$GPRMC,110910.59,A,5047.3986,N,00054.6007,W,0.08,0.19,150920,0.24,W,D,V*75")
	fmt.Println(nm.GetMap())
	dmap := nm.GetMap()
	dmap["hdm"] = "123.5°M"
	nm.Update(dmap)
	fmt.Println(dmap)
	hdm, _ := nm.WriteSentence("hc", "hdm")
	fmt.Println(hdm)
	fmt.Println(nm.GetMap())


	for{
		(*channels)["to_compass"] <- "$HCHDM,172.5,M*28"
		(*channels)["to_2000"] <- "$HCHDM,172.5,M*28"
		time.Sleep(100 * time.Millisecond)
	}
}
