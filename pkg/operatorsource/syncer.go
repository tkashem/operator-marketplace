package operatorsource

import (
	"time"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewRegistrySyncer returns a new instance of RegistrySyncer interface.
func NewRegistrySyncer(client client.Client) RegistrySyncer {
	return &registrySyncer{
		poller: NewPoller(client),
	}
}

// RegistrySyncer is an interface that wraps the Sync method.
//
// Sync kicks off the registry sync operation every N (specified in threshold)
// minutes. Sync will stop running once the stop channel is closed.
type RegistrySyncer interface {
	Sync(threshold time.Duration, stop <-chan struct{})
}

// registrySyncer implements RegistrySyncer interface.
type registrySyncer struct {
	poller Poller
}

func (s *registrySyncer) Sync(threshold time.Duration, stop <-chan struct{}) {
	log.Info("[sync] Starting operator source sync loop")
	for {
		select {
		case <-time.After(threshold * time.Minute):
			log.Debug("[sync] Checking for operator source update(s) in remote registry")
			s.poller.Poll()

		case <-stop:
			log.Info("[sync] Ending operator source watch loop")
			return
		}
	}
}
