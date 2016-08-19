package handlers

import (
	"encoding/json"
	"io"

	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rafael-azevedo/HPOMOutageTool/hputils"
)

type omit bool

//func SendStatus(w http.ResponseWriter, r *http.Request) {
//
//	ms := GetServiceStatus()
//
//	log.Println(ms)
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
//	log.Println("Restarting Services on ", vars["ServerName"])
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
type Server struct {
	Host string `json:"host"`
}

type OutageRequest struct {
	Username     string   `json:"username"`
	ChangeTicket string   `json:"changeticket"`
	IP           string   `json:"ip"`
	TimeIn       int      `json:"timein"`
	TimeOut      int      `json:"timeout"`
	ServerList   []Server `json:"serverlist"`
}

func ListOutage(w http.ResponseWriter, r *http.Request) {
	//mo := ParseOutage(string(CallTest()))
	mo := hputils.ParseOutage(string(hputils.CallAllOutage()))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(mo)
}

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	fileHandler := http.FileServer(http.Dir("./static/"))
	log.Println(*r.URL)
	fileHandler.ServeHTTP(w, r)
}

func ListOutageNodes(w http.ResponseWriter, r *http.Request) {
	//mo := ParseOutage(string(CallTest()))
	mo := hputils.ParseOutage(string(hputils.CallAllOutage()))

	var mn hputils.MultiName
	for i := range mo {
		data := struct {
			*hputils.Outage

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
	io := hputils.IsInOutage(hputils.ParseOutage(string(hputils.CallAllOutage())), vars["ServerName"])
	//io := IsInOutage(string(CallTest()))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(io)

}

func ListMultiNode(w http.ResponseWriter, r *http.Request) {
	var or OutageRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &or); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	log.Println(r.RemoteAddr)
	log.Println(r.Header.Get("X-FORWARDED-FOR"))
	log.Println(r.Header.Get("x-forwarded-for"))
	log.Println(r.Header.Get("X-Forwarded-For"))
	ips := []string{r.Header.Get("X-FORWARDED-FOR"), r.Header.Get("x-forwarded-for"), r.Header.Get("X-Forwarded-For")}
	log.Println(ips)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	//for i := range or.ServerList {
	//	log.Println(or.ServerList[i].Host)
	//}
	if err := json.NewEncoder(w).Encode(or); err != nil {
		panic(err)
	}
}

func AssignNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	io := hputils.IsInOutage(hputils.ParseOutage(string(hputils.CallAllOutage())), vars["ServerName"])
	if io.InOutage == false {
		a := (hputils.CallAssignNode(vars["ServerName"], "NETWORK_IP"))
		check, err := hputils.CheckError(string(a))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error()+" with NETWORK_IP", 500)
		}
		if check == true {
			a := string(hputils.CallAssignNode(vars["ServerName"], "PATTERN_OTHER"))
			check, err := hputils.CheckError(string(a))
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error()+" with PATTERN_OTHER", 500)
			}
			if check == true {
				a := string(hputils.CallAssignNode(vars["ServerName"], "PATTERN_IP_NAME"))
				check, err := hputils.CheckError(string(a))
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error()+" with PATTERN_IP_NAME", 500)
				}
				if check == true {
					a := string(hputils.CallAssignNode(vars["ServerName"], "PATTERN_IP_ADDR"))
					_, err := hputils.CheckError(string(a))
					if err != nil {
						log.Println(err)
						http.Error(w, err.Error()+" with PATTERN_IP_ADDR", 500)
					}
				}
			}
		}
		if check == false {
			Msg := hputils.Msg{Msg: "Operation successfully completed."}
			log.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		}
	}
	if io.InOutage == true {
		Msg := hputils.Msg{Msg: "Node " + vars["ServerName"] + " is in outage"}
		log.Println(Msg.Msg)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(Msg)
	}
}

func DeassignNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	io := hputils.IsInOutage(hputils.ParseOutage(string(hputils.CallAllOutage())), vars["ServerName"])

	if io.InOutage == true {
		a := (hputils.CallDeassignNode(vars["ServerName"], "NETWORK_IP"))
		log.Println(string(a))
		check, err := hputils.CheckError(string(a))
		if err != nil {
			http.Error(w, err.Error()+" with NETWORK_IP", 500)
		}
		if check == true {
			a := string(hputils.CallDeassignNode(vars["ServerName"], "PATTERN_OTHER"))
			check, err := hputils.CheckError(string(a))
			if err != nil {
				http.Error(w, err.Error()+" with PATTERN_OTHER", 500)
			}
			if check == true {
				a := string(hputils.CallDeassignNode(vars["ServerName"], "PATTERN_IP_NAME"))
				check, err := hputils.CheckError(string(a))
				if err != nil {
					http.Error(w, err.Error()+" with PATTERN_IP_NAME", 500)
				}
				if check == true {
					a := string(hputils.CallDeassignNode(vars["ServerName"], "PATTERN_IP_ADDR"))
					_, err := hputils.CheckError(string(a))
					if err != nil {
						log.Println(err)
						http.Error(w, err.Error()+" with PATTERN_IP_ADDR", 500)
					}
				}
			}
		}
		if check == false {
			Msg := hputils.Msg{Msg: "Operation successfully completed."}
			log.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		}
	}
	if io.InOutage == false {
		a := string(hputils.CallDeassignNode(vars["ServerName"], "NETWORK_IP"))
		switch {
		case strings.Contains(a, "Operation successfully completed"):
			log.Println(a)
			Msg := hputils.Msg{Msg: "Node " + vars["ServerName"] + " is not in outage"}
			log.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		default:
			log.Println(a)
			Msg := hputils.Msg{Msg: "Ouchies, Node was not found in outage " + a}
			log.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		}
	}
}
