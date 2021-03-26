package local

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	base2 "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	nexusHelper "github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

const JoinedLocalNexusName = "JoinedLocalNexus"

type JoinedLocalNexus struct {
	*base.NexusBase
	outgoingNexus compat.RevokableNexus
	incomingNexus compat.RevokableNexus

	// TODO event listeners
}

func NewJoinedLocalNexus() *JoinedLocalNexus {
	bn := base.NewBaseNexus(JoinedLocalNexusName)
	return &JoinedLocalNexus{
		NexusBase:     bn,
		outgoingNexus: nil,
		incomingNexus: nil,
	}
}

func (j *JoinedLocalNexus) SendFromOutgoing(msg base2.MoneysocketMessage) {
	log.Printf("from outgoing %s", msg)
	j.incomingNexus.OnMessage(j, msg)
}

func (j *JoinedLocalNexus) SendFromIncoming(msg base2.MoneysocketMessage) {
	log.Printf("from incoming %s", msg)
	j.outgoingNexus.OnMessage(j, msg)
}

func (j *JoinedLocalNexus) SendBinFromIncoming(msg []byte) {
	log.Printf("raw from incoming: %d", len(msg))
	j.outgoingNexus.OnBinMessage(j, msg)
}

func (j *JoinedLocalNexus) SetIncomingNexus(incomingNexus compat.RevokableNexus) {
	j.incomingNexus = incomingNexus
}

func (j *JoinedLocalNexus) SetOutgoingNexus(outgoingNexus compat.RevokableNexus) {
	j.outgoingNexus = outgoingNexus
}

func (j *JoinedLocalNexus) InitiateClose() {
	j.incomingNexus.InitiateClose()
	j.outgoingNexus.InitiateClose()
}

var _ nexusHelper.Nexus = &JoinedLocalNexus{}
