package operatorsource

import (
	"github.com/operator-framework/operator-marketplace/pkg/appregistry"
	"github.com/operator-framework/operator-marketplace/pkg/datastore"
	"github.com/operator-framework/operator-marketplace/pkg/phase"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewPoller returns a new instance of Poller interface.
func NewPoller(client client.Client) Poller {
	poller := &poller{
		datastore: datastore.Cache,
		checker: &checker{
			factory:      appregistry.NewClientFactory(),
			datastore:    datastore.Cache,
			client:       client,
			transitioner: phase.NewTransitioner(),
		},
	}

	return poller
}

// Poller is an interface that wraps the Poll method.
//
// Poll iterates through all available operator source(s) that are in the
// underlying datastore and performs the following action(s):
//   a) It polls the remote registry namespace to check if there are any
//      update(s) available.
//
//   b) If there is an update available then it triggers a purge andrebuild
//      operation for the specified OperatorSource object.
//
// On ay error during each iteration it logs the error encountered and moves
// on ti the next OperatorSource object.
type Poller interface {
	Poll()
}

// poller implements the Poller interface.
type poller struct {
	checker   UpdateChecker
	datastore datastore.Writer
}

func (w *poller) Poll() {
	sources := w.datastore.GetAllOperatorSources()

	for _, source := range sources {
		updated, err := w.checker.Check(source)
		if err != nil {
			log.Errorf("[sync] error checking for updates [%s] - %v", source.Name, err)
			continue
		}

		if !updated {
			continue
		}

		log.Infof("[sync] remote registry has update(s) - purging OperatorSource [%s]", source.Name)
		deleted, err := w.checker.Trigger(source)
		if err != nil {
			log.Errorf("[sync] error updating object [%s] - %v", source.Name, err)
		}

		if deleted {
			log.Infof("[sync] object deleted [%s] - no action taken", source.Name)
		}
	}
}
