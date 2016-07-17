package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	//Route{
	//	"RemoveFromOutage",
	//	"GET",
	//	"/outage/{ServerName}",
	//	SendStatus,
	//},
	//Route{
	//	"PutINOutage",
	//	"GET",
	//	"/outage/{ServerName}",
	//	RestartServices,
	//},
	Route{
		"ListOutage",
		"GET",
		"/outage",
		ListOutage,
	},
	Route{
		"ListOutageNodes",
		"GET",
		"/outage/nodes",
		ListOutageNodes,
	},
	Route{
		"ServeHTML",
		"GET",
		"/",
		ServeHTML,
	},
	Route{
		"NodeInOutage",
		"GET",
		"/outage/{ServerName}",
		ListSingleNode,
	},
	Route{
		"AssignNode",
		"GET",
		"/outage/assign/{ServerName}",
		AssignNode,
	},
	Route{
		"DeassignNode",
		"GET",
		"/outage/deassign/{ServerName}",
		DeassignNode,
	},
}
