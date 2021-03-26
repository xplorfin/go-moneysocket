package provider

import (
	"log"

	"github.com/xplorfin/moneysocket-go/moneysocket/layer/compat"
	message_base "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/notification"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

const NexusName = "ProviderNexus"

type Nexus struct {
	*base.NexusBase
	RequestReferenceUUID      string
	handleProviderInfoRequest compat.HandleProviderInfoRequest
	ProviderFinishedCb        func(nx nexus.Nexus)
}

func NewProviderNexus(belowNexus nexus.Nexus) *Nexus {
	baseNexus := base.NewBaseNexusBelow(NexusName, belowNexus)
	pn := Nexus{baseNexus, "", nil, nil}
	belowNexus.SetOnBinMessage(pn.OnBinMessage)
	belowNexus.SetOnMessage(pn.OnMessage)
	return &pn
}

func (o *Nexus) IsLayerMessage(message message_base.MoneysocketMessage) bool {
	if message.MessageClass() == message_base.Request {
		return false
	}
	ntfn := message.(notification.MoneysocketNotification)
	return ntfn.RequestType() == message_base.ProviderRequest || ntfn.RequestType() == message_base.PingRequest
}

func (o *Nexus) NotifyProvider() {
	ss := o.SharedSeed()
	providerInfo := o.handleProviderInfoRequest(*ss)
	_ = o.Send(notification.NewNotifyProvider(o.UUID().String(), providerInfo.Details.Payer(), providerInfo.Details.Payee(), providerInfo.Details.Wad, o.RequestReferenceUUID))
}

func (o *Nexus) NotifyProviderNotReady() {
	_ = o.Send(notification.NewNotifyProviderNotReady(o.RequestReferenceUUID))
}

func (o *Nexus) OnMessage(belowNexus nexus.Nexus, msg message_base.MoneysocketMessage) {
	log.Println("provider nexus got message")
	if !o.IsLayerMessage(msg) {
		o.NexusBase.OnMessage(belowNexus, msg)
		return
	}
	ntfn := msg.(notification.MoneysocketNotification)
	if ntfn.RequestType() == message_base.ProviderRequest {
		ss := belowNexus.SharedSeed()
		providerInfo := o.handleProviderInfoRequest(*ss)
		if providerInfo.Details.Ready() {
			o.NotifyProvider()
			o.ProviderFinishedCb(o)
		} else {
			o.NotifyProviderNotReady()
			o.Layer.(compat.SellingLayerInterface).NexusWaitingForApp(ss, o)
		}
	} else if ntfn.RequestType() == message_base.PingRequest {
		o.NotifyPong()
	}

}

func (o *Nexus) NotifyPong() {
	_ = o.Send(notification.NewNotifyPong(o.RequestReferenceUUID))
}

func (o *Nexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	log.Println("provider nexus got raw msg")
	o.NexusBase.OnBinMessage(belowNexus, msg)
}

func (o *Nexus) WaitForConsumer(providerFinishedCb func(nexus2 nexus.Nexus)) {
	o.ProviderFinishedCb = providerFinishedCb
}

func (o *Nexus) NotifyProviderReady() {
	ss := o.SharedSeed()
	providerInfo := o.Layer.(compat.ProviderTransactLayerInterface).HandleProviderInfoRequest(*ss)
	if !providerInfo.Details.Ready() {
		panic("expected provider to be ready")
	}
	_ = o.Send(notification.NewNotifyProvider(o.UUID().String(), providerInfo.Details.Payer(), providerInfo.Details.Payee(), providerInfo.Details.Wad, o.RequestReferenceUUID))
}
func (o *Nexus) ProviderNowReady() {
	o.NotifyProviderReady()
}
