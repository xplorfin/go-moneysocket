package local

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer"
	base_moneysocket "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

type LocalNexus struct {
	base.BaseNexus
}

const LocalNexusName = "LocalNexusName"

func NewLocalNexus(belowNexus nexus.Nexus, layer layer.Layer) LocalNexus {
	ln := LocalNexus{
		base.NewBaseNexusFull(LocalNexusName, belowNexus, layer),
	}
	ln.RegisterAboveNexus(belowNexus)
	return ln
}

func (l *LocalNexus) OnMessage(belowNexus nexus.Nexus, msg base_moneysocket.MoneysocketMessage) {
	log.Println("local nexus got msg")
	l.BaseNexus.OnMessage(belowNexus, msg)
}

var _ nexus.Nexus = &LocalNexus{}
