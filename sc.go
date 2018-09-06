package tinnitus

import (
	"github.com/scgolang/sc"
)

type SCClient struct {
	*sc.Client
	group *sc.GroupNode
	ids   chan chan int32
}

var SC *SCClient

func InitSuperCollider() {
	client, err := sc.NewClient("udp", "0.0.0.0:0", Config.SC.Host, Config.SC.Timeout)

	if err != nil {
		panic(err)
	}

	group, err := client.AddDefaultGroup()

	if err != nil {
		panic(err)
	}

	SC = &SCClient{Client: client, group: group}
}
