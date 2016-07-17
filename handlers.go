package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type omit bool

//func SendStatus(w http.ResponseWriter, r *http.Request) {
//
//	ms := GetServiceStatus()
//
//	fmt.Println(ms)
//
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusOK)
//	if err := json.NewEncoder(w).Encode(ms); err != nil {
//		panic(err)
//	}
//}
//
//func RestartServices(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	fmt.Println("Restarting Services on ", vars["ServerName"])
//}
//
//func SendLS(w http.ResponseWriter, r *http.Request) {
//	LSS := (ExecLS("ls"))
//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
//	w.WriteHeader(http.StatusOK)
//	if err := json.NewEncoder(w).Encode(LSS); err != nil {
//		panic(err)
//	}
//}

func ListOutage(w http.ResponseWriter, r *http.Request) {
	//mo := ParseOutage(string(CallTest()))
	mo := ParseOutage(string(CallAllOutage()))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(mo)
}

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	fileHandler := http.FileServer(http.Dir("./static/"))
	fmt.Println(*r.URL)
	fileHandler.ServeHTTP(w, r)
}

func ListOutageNodes(w http.ResponseWriter, r *http.Request) {
	//mo := ParseOutage(string(CallTest()))
	mo := ParseOutage(string(CallAllOutage()))

	var mn MultiName
	for i := range mo {
		data := struct {
			*Outage

			OmitLabel       omit `json:"label,omitempty"`
			OmitIPAddress   omit `json:"ipaddress,omitempty"`
			OmitNetworkType omit `json:"networktype,omitempty"`
			OmitMachineType omit `json:"machinetype,omitempty"`
			OmitCommType    omit `json:"commtype,omitempty"`
			OmitDHCPenabled omit `json:"dhcpenabled,omitempty"`
		}{
			Outage: &mo[i],
		}
		mn = append(mn, data)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(mn)
}

func ListSingleNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	io := IsInOutage(ParseOutage(string(CallAllOutage())), vars["ServerName"])
	//io := IsInOutage(string(CallTest()))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(io)

}

func AssignNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	io := IsInOutage(ParseOutage(string(CallAllOutage())), vars["ServerName"])
	if io.InOutage == false {
		a := (CallAssignNode(vars["ServerName"], "NETWORK_IP"))
		check, err := CheckError(string(a))
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error()+" with NETWORK_IP", 500)
		}
		if check == true {
			a := string(CallAssignNode(vars["ServerName"], "PATTERN_OTHER"))
			check, err := CheckError(string(a))
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error()+" with PATTERN_OTHER", 500)
			}
			if check == true {
				a := string(CallAssignNode(vars["ServerName"], "PATTERN_IP_NAME"))
				check, err := CheckError(string(a))
				if err != nil {
					fmt.Println(err)
					http.Error(w, err.Error()+" with PATTERN_IP_NAME", 500)
				}
				if check == true {
					a := string(CallAssignNode(vars["ServerName"], "PATTERN_IP_ADDR"))
					_, err := CheckError(string(a))
					if err != nil {
						fmt.Println(err)
						http.Error(w, err.Error()+" with PATTERN_IP_ADDR", 500)
					}
				}
			}
		}
		if check == false {
			Msg := Msg{Msg: "Operation successfully completed."}
			fmt.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		}
	}
	if io.InOutage == true {
		Msg := Msg{Msg: "Node " + vars["ServerName"] + " is in outage"}
		fmt.Println(Msg.Msg)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(Msg)
	}
}

func DeassignNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	io := IsInOutage(ParseOutage(string(CallAllOutage())), vars["ServerName"])

	if io.InOutage == true {
		a := (CallDeassignNode(vars["ServerName"], "NETWORK_IP"))
		fmt.Println(string(a))
		check, err := CheckError(string(a))
		if err != nil {
			http.Error(w, err.Error()+" with NETWORK_IP", 500)
		}
		if check == true {
			a := string(CallDeassignNode(vars["ServerName"], "PATTERN_OTHER"))
			check, err := CheckError(string(a))
			if err != nil {
				http.Error(w, err.Error()+" with PATTERN_OTHER", 500)
			}
			if check == true {
				a := string(CallDeassignNode(vars["ServerName"], "PATTERN_IP_NAME"))
				check, err := CheckError(string(a))
				if err != nil {
					http.Error(w, err.Error()+" with PATTERN_IP_NAME", 500)
				}
				if check == true {
					a := string(CallDeassignNode(vars["ServerName"], "PATTERN_IP_ADDR"))
					_, err := CheckError(string(a))
					if err != nil {
						fmt.Println(err)
						http.Error(w, err.Error()+" with PATTERN_IP_ADDR", 500)
					}
				}
			}
		}
		if check == false {
			Msg := Msg{Msg: "Operation successfully completed."}
			fmt.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		}
	}
	if io.InOutage == false {
		a := string(CallDeassignNode(vars["ServerName"], "NETWORK_IP"))
		switch {
		case strings.Contains(a, "Operation successfully completed"):
			fmt.Println(a)
			Msg := Msg{Msg: "Node " + vars["ServerName"] + " is not in outage"}
			fmt.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		default:
			fmt.Println(a)
			Msg := Msg{Msg: "Ouchies, Node was not found in outage " + a}
			fmt.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		}
	}
}
