package main

import (
	"errors"
	"strings"
)

const (
	nameServiceSeparator  = ":"
	nameServicesSeparator = ","
)

type hostServicePair struct {
	Host    string
	Service string
}

func newHostServicePair(str string) (hostServicePair, error) {
	if strings.Index(str, nameServiceSeparator) == -1 {
		return hostServicePair{}, errors.New("Invalid format")
	}
	slices := strings.Split(str, nameServiceSeparator)
	return hostServicePair{
		Host:    slices[0],
		Service: slices[1],
	}, nil
}

func newHostServicePairs(str string) ([]hostServicePair, error) {
	if strings.Index(str, nameServicesSeparator) == -1 {
		// Process as single entry
		p, e := newHostServicePair(str)
		return []hostServicePair{p}, e
	}
	slices := strings.Split(str, nameServicesSeparator)
	hsp := make([]hostServicePair, len(slices))
	var e error
	for i := range slices {
		hsp[i], e = newHostServicePair(strings.TrimSpace(slices[i]))
		if e != nil {
			return []hostServicePair{}, e
		}
	}
	return hsp, nil
}
