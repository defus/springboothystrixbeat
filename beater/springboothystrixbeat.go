package beater

import (
	"fmt"
	"time"
	"net/url"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/defus/springboothystrixbeat/config"
)

// Springboothystrixbeat configuration.
type Springboothystrixbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of springboothystrixbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Springboothystrixbeat{
		done:   make(chan struct{}),
		config: c,
	}

	_, err := url.Parse(c.Host)
        if err != nil {
                logp.Err("Invalid Spring Boot URL: %v", err)
                return nil, err
        }

	logp.Debug("springboothystrixbeat", "Init springboothystrixbeat")
	logp.Debug("springboothystrixbeat", "Period %v\n", bt.config.Period)
	logp.Debug("springboothystrixbeat", "Watch %v", c.Host)

	return bt, nil
}

// Run starts springboothystrixbeat.
func (bt *Springboothystrixbeat) Run(b *beat.Beat) error {
	logp.Info("springboothystrixbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		//Générr l'evenement en recuperant la data de actuator
                bt.ProcessMetricsActuator(b)
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

//		event := beat.Event{
//			Timestamp: time.Now(),
//			Fields: common.MapStr{
//				"type":    b.Info.Name,
//				"counter": counter,
//			},
//		}
//		bt.client.Publish(event)
//		logp.Info("Event sent")
//		counter++
	}
}

// Stop stops springboothystrixbeat.
func (bt *Springboothystrixbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
