package main

import (
	"time"

	"gopkg.in/mgo.v2"
)

const (
	collectionName = "request_analytics"
)

type requestAnalytics struct {
	URL         string
	Method      string
	RequestTime int64
	Day         time.Weekday
	Hour        int
}

type mongo struct {
	sess *mgo.Session
}

func (m mongo) Close() error {
	m.sess.Close()
	return nil
}

func (m mongo) Write(r requestAnalytics) error {
	return m.sess.DB("pusher_tutorial").C(collectionName).Insert(r)
}

func newMongo(addr string) (mongo, error) {
	sess, err := mgo.Dial(addr)
	if err != nil {
		return mongo{}, err
	}

	return mongo{
		sess: sess,
	}, nil
}
