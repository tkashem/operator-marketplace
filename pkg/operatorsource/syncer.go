package operatorsource

import (
	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
)

type Syncer interface {
	Sync(opsrc *v1alpha1.OperatorSource) error
}

type syncer struct {
}

func (s *syncer) Sync(opsrc *v1alpha1.OperatorSource) error {
	// change the status to purging and let reconciliation begin

	return nil
}

type SyncManager struct {
}

func (sm *SyncManager) Sync() {

}
