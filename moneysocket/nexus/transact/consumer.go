package transact

import (
	moneysocket_message "github.com/xplorfin/moneysocket-go/moneysocket/message/base"
	"github.com/xplorfin/moneysocket-go/moneysocket/message/request"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus"
	"github.com/xplorfin/moneysocket-go/moneysocket/nexus/base"
)

// TODO handle invoices
type OnInvoice func(transactNexus nexus.Nexus, invoice string, requestReferenceUuid string)

// TODO handl epreimages
type OnPreimage func(transactNexus nexus.Nexus, preimage string, requestReferenceUuid string)

type OnProviderInfo func(consumerTransactNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage)

type ConsumerTrackNexus struct {
	nexus.Nexus
	// invoice event handler
	onInvoice OnInvoice
	// preimage event handler
	onPreimage OnPreimage
	// aclled on provder info
	onProviderInfo OnProviderInfo
}

const ConsumerTrackNexusName = "ConsumerTrackNexus"

func NewConsumerTransactNexus(belowNexus nexus.Nexus) ConsumerTrackNexus {
	c := ConsumerTrackNexus{
		base.NewBaseNexusBelow(ConsumerTrackNexusName, belowNexus),
		nil,
		nil,
		nil,
	}
	belowNexus.SetOnBinMessage(c.OnBinMessage)
	belowNexus.SetOnMessage(c.OnMessage)
	return c
}

func (c ConsumerTrackNexus) HandleLayerNotification() {
	panic("method not yet implemented")
}

func (c ConsumerTrackNexus) IsLayerMessage(msg []byte) bool {
	// TODO handle (this will be transitioned to ana ctual message handler)
	panic("method not yet implemented")
}

func (c ConsumerTrackNexus) OnMessage(belowNexus nexus.Nexus, message moneysocket_message.MoneysocketMessage) {
	panic("method not yet implemented")
}

func (c ConsumerTrackNexus) OnBinMessage(belowNexus nexus.Nexus, msg []byte) {
	// DO nothing
}

// call on invoice function
func (c *ConsumerTrackNexus) OnInvoice(transactNexus nexus.Nexus, invoice string, requestReferenceUuid string) {
	if c.onInvoice != nil {
		c.onInvoice(transactNexus, invoice, requestReferenceUuid)
	}
}

// set function to be called when on invoice is called
func (c *ConsumerTrackNexus) SetOnInvoice(invoice OnInvoice) {
	c.onInvoice = invoice
}

// call on preimage function
func (c *ConsumerTrackNexus) OnPreImage(transactNexus nexus.Nexus, preimage string, requestReferenceUuid string) {
	if c.onPreimage != nil {
		c.onPreimage(transactNexus, preimage, requestReferenceUuid)
	}
}

// set function to be called when onPreImage is called
func (c *ConsumerTrackNexus) SetOnPreimage(preimage OnPreimage) {
	c.onPreimage = preimage
}

// set function to be called when OnProviderInfo is called
func (c *ConsumerTrackNexus) SetOnProviderInfo(info OnProviderInfo) {
	c.onProviderInfo = info
}

func (c *ConsumerTrackNexus) OnProviderInfo(consumerTransactNexus nexus.Nexus, msg moneysocket_message.MoneysocketMessage) {
	if c.onProviderInfo != nil {
		c.onProviderInfo(consumerTransactNexus, msg)
	}
}

func (c ConsumerTrackNexus) RequestInvoice(msats int64, description string) (uuid string) {
	ri := request.NewRequestInvoice(msats)
	c.Send(ri)
	return ri.Uuid()
}

func (c ConsumerTrackNexus) RequestPay(bolt11 string) (uuid string) {
	rp := request.NewRequestPay(bolt11)
	c.Send(rp)
	return rp.Uuid()
}

var _ nexus.Nexus = &ConsumerTrackNexus{}
