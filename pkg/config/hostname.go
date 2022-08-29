package config

import (
	"bytes"
	"math/rand"
)

func Hostname() string {

	hostnames := [8]string{"cat", "world", "coffee", "girl", "man", "book", "pinguin", "moon"}
	hostStatus := [16]string{"blue", "red", "green", "yellow", "big", "small", "tall", "thin", "round", "square", "triangular", "weird", "noisy", "silent", "soft", "irregular"}

	var buf bytes.Buffer

	rnd8 := rand.Intn(8)
	rnd16 := rand.Intn(16)
	buf.WriteString(hostnames[rnd8])
	buf.WriteString("-")
	buf.WriteString(hostStatus[rnd16])
	return buf.String()
}
