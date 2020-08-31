package main

import "errors"

var (
	// ErrNotExist means no wanted resource
	ErrNotExist = errors.New("Resource does not exist")
)
