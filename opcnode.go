package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Name interface {
}

type ErrorMsg struct {
	Error error `json:"error"`
}

type Msg struct {
	Msg string `json:"msg"`
}
type Outage struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	IPAddress   string `json:"ipaddress"`
	NetworkType string `json:"networktype"`
	MachineType string `json:"machinetype"`
	CommType    string `json:"commtype"`
	DHCPenabled string `json:"dhcpenabled"`
}

type InOutage struct {
	InOutage bool `json:"outage"`
}

type MultiName []Name

type MultiOutage []Outage

func CallAllOutage() []byte {
	cmdName := "/opt/OV/bin/opcnode"
	args := []string{"-list_nodes", "group_name=outage"}
	cmd := exec.Command(cmdName, args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	return out
}

func CallTest() []byte {
	cmdName := "cat"
	args := []string{"test.txt"}
	cmd := exec.Command(cmdName, args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	return out
}

func ParseOutage(text string) MultiOutage {
	r := strings.NewReplacer("==", "",
		"=", ":", "List of all Nodes in the HPOM database:", "", "Operation successfully completed.", "")
	text = r.Replace(text)
	lines := strings.Split(text, "\n")
	var mo MultiOutage
	//nameNumber := 0
	itemNumber := 0
	var fieldslice []string
	for _, line := range lines {
		if len(line) > 0 {
			sl := strings.Split(line, ":")
			if itemNumber <= 6 {
				field := strings.TrimSpace(sl[1])
				fieldslice = append(fieldslice, field)
				itemNumber++
				if itemNumber > 6 {
					o := Outage{fieldslice[0], fieldslice[1], fieldslice[2], fieldslice[3], fieldslice[4], fieldslice[5], fieldslice[6]}
					itemNumber = 0
					fieldslice = nil
					mo = append(mo, o)
				}
			}
		}
	}
	return mo
}

func IsInOutage(mo MultiOutage, name string) InOutage {
	var o InOutage
	for i := range mo {

		if strings.Contains(mo[i].Name, name) {
			o.InOutage = true
		}
	}
	return o
}

func CallAssignNode(name string, net string) []byte {
	cmdName := "/opt/OV/bin/opcnode"
	nodeName := "node_name=" + name
	netType := "net_type=" + net
	args := []string{"-assign_node", "-list_nodes", "group_name=outage", nodeName, netType}
	cmd := exec.Command(cmdName, args...)
	out, _ := cmd.Output()
	return out
}

func CallDeassignNode(name string, net string) []byte {
	cmdName := "/opt/OV/bin/opcnode"
	nodeName := "node_name=" + name
	netType := "net_type=" + net
	args := []string{"-deassign_node", "-group_name=outage", "group_name=outage", nodeName, netType}
	cmd := exec.Command(cmdName, args...)
	out, _ := cmd.Output()
	return out
}

func CheckError(text string) (bool, error) {
	switch {
	case strings.Contains(text, "failed"):
		return true, errors.New(text)
	case strings.Contains(text, "Illegal"):
		return true, errors.New(text)
	case strings.Contains(text, "Operation successfully completed."):
		fmt.Println(text)
		return false, nil
	default:
		return true, errors.New(text)
	}
}
