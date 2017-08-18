package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

func dumpWeatherToDatabase(info Response, address, username, password, option string) error {
	s, err := mgo.Dial(fmt.Sprintf("mongodb://%s:%s@%s/%s", username, password, address, option))
	if err != nil {
		return err
	}
	defer s.Close()
	info.Data = fmt.Sprintf("%d/%d/%d", time.Now().Day(), int(time.Now().Month()), time.Now().Year())
	return s.DB("weather").C("history").Insert(info)
}

func dumpWeatherToDatabaseInline(info Response, s *mgo.Session) error {
	info.Data = fmt.Sprintf("%d/%d/%d", time.Now().Day(), int(time.Now().Month()), time.Now().Year())
	return s.DB("weather").C("history").Insert(info)
}

func NewDbSession(address, username, password, option string, enablessl bool) (*mgo.Session, error) {
	if enablessl {
		return DialViaSSL(address, option, username, password)
	} else {
		return mgo.Dial(fmt.Sprintf("mongodb://%s:%s@%s/%s", username, password, address, option))
	}
}

func DialViaSSL(addresses string, dboption string, username string, password string) (*mgo.Session, error) {
	dboptions := strings.Split(dboption, "=")
	if len(dboption) < 2 {
		dboptions = make([]string, 2)
		dboptions[1] = "admin"
	}
	tlsConfig := &tls.Config{}
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{addresses},
		Database: dboptions[1],
		Username: username,
		Password: password,
	}

	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}
	session.EnsureSafe(&mgo.Safe{
		W:     0,
		FSync: false,
	})
	return session, nil
}
