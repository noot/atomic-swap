package main

import (
	"context"
	"testing"
)

func TestDaemon_Wait(_ *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := &daemon{
		ctx:    ctx,
		cancel: cancel,
	}

	go d.wait()
}
