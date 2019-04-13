package main

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gopkg.in/mgo.v2"
)

const (
	collectionName = "request_analytics"
)

type requestAnalytics struct {
	URL         string       `json:"url"`
	Method      string       `json:"method"`
	RequestTime int64        `json:"request_time"`
	Day         time.Weekday `json:"day"`
	Hour        int          `json:"hour"`
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

func (m mongo) Count() (int, error) {
	return m.sess.DB("pusher_tutorial").C(collectionName).Count()
}

type statsPerRoute struct {
	ID struct {
		Method string `bson:"method" json:"method"`
		URL    string `bson:"url" json:"url"`
	} `bson:"_id" json:"id"`
	NumberOfRequests int `bson:"numberOfRequests" json:"number_of_requests"`
}

func (m mongo) AverageResponseTime() (float64, error) {

	type res struct {
		AverageResponseTime float64 `bson:"averageResponseTime" json:"average_response_time"`
	}

	var ret = []res{}

	var baseMatch = bson.M{
		"$group": bson.M{
			"_id":                 nil,
			"averageResponseTime": bson.M{"$avg": "$requesttime"},
		},
	}

	err := m.sess.DB("pusher_tutorial").C(collectionName).
		Pipe([]bson.M{baseMatch}).All(&ret)
	return ret[0].AverageResponseTime, err
}

func (m mongo) StatsPerRoute() ([]statsPerRoute, error) {

	var ret []statsPerRoute

	var baseMatch = bson.M{
		"$group": bson.M{
			"_id":              bson.M{"url": "$url", "method": "$method"},
			"responseTime":     bson.M{"$avg": "$requesttime"},
			"numberOfRequests": bson.M{"$sum": 1},
		},
	}

	err := m.sess.DB("pusher_tutorial").C(collectionName).
		Pipe([]bson.M{baseMatch}).All(&ret)
	return ret, err
}

type requestsPerDay struct {
	ID               int `bson:"_id" json:"id"`
	NumberOfRequests int `bson:"numberOfRequests" json:"number_of_requests"`
}

func (m mongo) RequestsPerHour() ([]requestsPerDay, error) {

	var ret []requestsPerDay

	var baseMatch = bson.M{
		"$group": bson.M{
			"_id":              "$hour",
			"numberOfRequests": bson.M{"$sum": 1},
		},
	}

	var sort = bson.M{
		"$sort": bson.M{
			"numberOfRequests": 1,
		},
	}

	err := m.sess.DB("pusher_tutorial").C(collectionName).
		Pipe([]bson.M{baseMatch, sort}).All(&ret)
	return ret, err
}

func (m mongo) RequestsPerDay() ([]requestsPerDay, error) {

	var ret []requestsPerDay

	var baseMatch = bson.M{
		"$group": bson.M{
			"_id":              "$day",
			"numberOfRequests": bson.M{"$sum": 1},
		},
	}

	var sort = bson.M{
		"$sort": bson.M{
			"numberOfRequests": 1,
		},
	}

	err := m.sess.DB("pusher_tutorial").C(collectionName).
		Pipe([]bson.M{baseMatch, sort}).All(&ret)
	return ret, err
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
