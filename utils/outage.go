package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Host string `json:"host"`
}

type OutageRequest struct {
	Username     string    `json:"username"`
	ChangeTicket string    `json:"changeticket"`
	IP           string    `json:"ip"`
	RequestTime  time.Time `json:"requesttime"`
	ExpTime      time.Time `json:"exptime"`
	ServerList   []Server  `json:"serverlist"`
}

func GetIP(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func SetExpTime() (requestTime time.Time, expTime time.Time) {
	requestTime = time.Now()
	expTime = requestTime.AddDate(0, 0, 60)
	return
}

func ParsePost(r *http.Request) (OutageRequest, error) {
	var or OutageRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return or, err
	}
	if err := r.Body.Close(); err != nil {
		return or, err
	}
	if err := json.Unmarshal(body, &or); err != nil {
		return or, err
	}
	or.IP = GetIP(r)
	or.RequestTime, or.ExpTime = SetExpTime()
	return or, nil
}
