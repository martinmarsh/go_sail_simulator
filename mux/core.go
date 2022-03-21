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
	go ais(&channels)

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
	course := 100.2
	data := make(map[string]string)
	sg := -1.0

	for{
		data["hdm"] = fmt.Sprintf("%03.1f°M", course)
		nm.Update(data)
		hdm, _ := nm.WriteSentence("hc", "hdm")
		(*channels)["to_compass"] <- hdm
		//(*channels)["to_2000"] <- hdm
		fmt.Println(hdm)
		time.Sleep(100 * time.Millisecond)
		if course > 105{
			sg = -1.0
		}else if course < 95 {
			sg = 1.0
		} 
		course += 0.1*sg

	}
}

func ais(channels *map[string](chan string)){


	for{
		
		(*channels)["to_2000"] <- "!AIVDM,1,1,,A,13aEOK?P00PD2wVMdLDRhgvL289?,0*26"
		time.Sleep(10 * time.Millisecond)
		(*channels)["to_2000"] <- "!AIVDM,1,1,,B,16S`2cPP00a3UF6EKT@2:?vOr0S2,0*00"
		time.Sleep(10 * time.Millisecond)
		(*channels)["to_2000"] <- "!AIVDM,2,1,9,B,53nFBv01SJ<thHp6220H4heHTf2222222222221?50:454o<`9QSlUDp,0*09"
		time.Sleep(10 * time.Millisecond)
		(*channels)["to_2000"] <- "!AIVDM,2,2,9,B,888888888888880,2*2E"
		time.Sleep(10 * time.Millisecond)
	}
}

