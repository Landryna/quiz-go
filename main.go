package main

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"time"
)

func main() {
	actors := newActorGroup()
	{
		server := fasthttp.Server{
			Name:               "Admin",
			IdleTimeout:        0,
			MaxRequestsPerConn: 0,
			TCPKeepalive:       false,
			TCPKeepalivePeriod: 0,
			MaxRequestBodySize: 0,
			ReduceMemoryUsage:  false,
			CloseOnShutdown:    true,
		}
		api := newAdminAPI(server)
		rf := func() error {
			return api.Run()
		}
		cf := func() error {
			return api.Close()
		}
		actor := newActor("Admin API", rf, cf)
		actors.add(actor)
	}

	{
		rf := func() error {
			fmt.Println("I'm started 2")
			time.Sleep(360 * time.Second)
			fmt.Println("I'm finished 2")
			return errors.New("ERROR! 2")
		}
		cf := func() error {
			fmt.Println("canceling func 2")
			return nil
		}
		actor := newActor("Test 2", rf, cf)
		actors.add(actor)
	}
	fmt.Println(actors.show())
	actors.run()
}
