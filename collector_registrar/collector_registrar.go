package collector_registrar

import (
	"encoding/json"

	"github.com/apcera/nats"
	"github.com/cloudfoundry-incubator/metricz"
	"github.com/cloudfoundry/yagnats"
)

type CollectorRegistrar interface {
	RegisterWithCollector(metricz.Component) error
}

type natsCollectorRegistrar struct {
	natsClient yagnats.NATSConn
}

func New(natsClient yagnats.NATSConn) CollectorRegistrar {
	return &natsCollectorRegistrar{
		natsClient: natsClient,
	}
}

func (registrar *natsCollectorRegistrar) RegisterWithCollector(component metricz.Component) error {
	message := NewAnnounceComponentMessage(component)

	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = registrar.natsClient.Subscribe(
		DiscoverComponentMessageSubject,
		func(msg *nats.Msg) {
			registrar.natsClient.Publish(msg.Reply, messageJson)
		},
	)
	if err != nil {
		return err
	}

	err = registrar.natsClient.Publish(
		AnnounceComponentMessageSubject,
		messageJson,
	)
	if err != nil {
		return err
	}

	return nil
}
