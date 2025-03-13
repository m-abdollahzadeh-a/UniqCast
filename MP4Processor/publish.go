package main

import "github.com/nats-io/nats.go"

func publish(nc *nats.Conn, subject string, msg []byte) error {
	return nc.Publish(subject, msg)
}
