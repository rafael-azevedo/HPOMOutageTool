package router

import (
	"net/http"

	"github.com/rafael-azevedo/HPOMOutageTool/handlers"
)

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
		handlers.ListOutage,
	},
	Route{
		"ListOutageNodes",
		"GET",
		"/outage/nodes",
		handlers.ListOutageNodes,
	},
	Route{
		"ServeHTML",
		"GET",
		"/",
		handlers.ServeHTML,
	},
	Route{
		"NodeInOutage",
		"GET",
		"/outage/{ServerName}",
		handlers.ListSingleNode,
	},
	Route{
		"AssignNode",
		"GET",
		"/outage/assign/{ServerName}",
		handlers.AssignNode,
	},
	Route{
		"DeassignNode",
		"GET",
		"/outage/deassign/{ServerName}",
		handlers.DeassignNode,
	},
	Route{
		"MultiNode",
		"POST",
		"/outage",
		handlers.ListMultiNode,
	},
}
