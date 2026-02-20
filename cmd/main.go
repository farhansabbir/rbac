package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/farhansabbir/rbac/lib/controllers"
)

func main() {
	ctrl := controllers.GetController()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
	ctrl.Stop()
	os.Exit(0)
}
