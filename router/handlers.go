package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rafael-azevedo/HPOMOutageTool/utils"
)

type omit bool

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	fileHandler := http.FileServer(http.Dir("./static/"))
	log.Println(*r.URL)
	fileHandler.ServeHTTP(w, r)
}

func ListOutage(w http.ResponseWriter, r *http.Request) {
	mo, _ := utils.List()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(mo)
}

func ListOutageNodes(w http.ResponseWriter, r *http.Request) {
	_, mn := utils.List()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(mn)
}

func ListSingleNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	io := utils.IsInOutage(utils.ParseOutage(string(utils.CallAllOutage())), vars["ServerName"])
	//io := IsInOutage(string(CallTest()))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(io)

}

func ListMultiNode(w http.ResponseWriter, r *http.Request) {

	or, err := utils.ParsePost(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(or); err != nil {
		panic(err)
	}
}

func AssignNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	io := utils.IsInOutage(utils.ParseOutage(string(utils.CallAllOutage())), vars["ServerName"])
	if io.InOutage == false {
		a := (utils.CallAssignNode(vars["ServerName"], "NETWORK_IP"))
		check, err := utils.CheckError(string(a))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error()+" with NETWORK_IP", 500)
		}
		if check == true {
			a := string(utils.CallAssignNode(vars["ServerName"], "PATTERN_OTHER"))
			check, err := utils.CheckError(string(a))
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error()+" with PATTERN_OTHER", 500)
			}
			if check == true {
				a := string(utils.CallAssignNode(vars["ServerName"], "PATTERN_IP_NAME"))
				check, err := utils.CheckError(string(a))
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error()+" with PATTERN_IP_NAME", 500)
				}
				if check == true {
					a := string(utils.CallAssignNode(vars["ServerName"], "PATTERN_IP_ADDR"))
					_, err := utils.CheckError(string(a))
					if err != nil {
						log.Println(err)
						http.Error(w, err.Error()+" with PATTERN_IP_ADDR", 500)
					}
				}
			}
		}
		if check == false {
			Msg := utils.Msg{Msg: "Operation successfully completed."}
			log.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		}
	}
	if io.InOutage == true {
		Msg := utils.Msg{Msg: "Node " + vars["ServerName"] + " is in outage"}
		log.Println(Msg.Msg)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(Msg)
	}
}

func DeassignNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	io := utils.IsInOutage(utils.ParseOutage(string(utils.CallAllOutage())), vars["ServerName"])

	if io.InOutage == true {
		a := (utils.CallDeassignNode(vars["ServerName"], "NETWORK_IP"))
		log.Println(string(a))
		check, err := utils.CheckError(string(a))
		if err != nil {
			http.Error(w, err.Error()+" with NETWORK_IP", 500)
		}
		if check == true {
			a := string(utils.CallDeassignNode(vars["ServerName"], "PATTERN_OTHER"))
			check, err := utils.CheckError(string(a))
			if err != nil {
				http.Error(w, err.Error()+" with PATTERN_OTHER", 500)
			}
			if check == true {
				a := string(utils.CallDeassignNode(vars["ServerName"], "PATTERN_IP_NAME"))
				check, err := utils.CheckError(string(a))
				if err != nil {
					http.Error(w, err.Error()+" with PATTERN_IP_NAME", 500)
				}
				if check == true {
					a := string(utils.CallDeassignNode(vars["ServerName"], "PATTERN_IP_ADDR"))
					_, err := utils.CheckError(string(a))
					if err != nil {
						log.Println(err)
						http.Error(w, err.Error()+" with PATTERN_IP_ADDR", 500)
					}
				}
			}
		}
		if check == false {
			Msg := utils.Msg{Msg: "Operation successfully completed."}
			log.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		}
	}
	if io.InOutage == false {
		a := string(utils.CallDeassignNode(vars["ServerName"], "NETWORK_IP"))
		switch {
		case strings.Contains(a, "Operation successfully completed"):
			log.Println(a)
			Msg := utils.Msg{Msg: "Node " + vars["ServerName"] + " is not in outage"}
			log.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		default:
			log.Println(a)
			Msg := utils.Msg{Msg: "Ouchies, Node was not found in outage " + a}
			log.Println(Msg.Msg)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(Msg)
		}
	}
}
