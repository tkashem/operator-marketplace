package operatorsource

import (
	"github.com/operator-framework/operator-marketplace/pkg/apis/marketplace/v1alpha1"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
)

type Watcher interface {
	Add(opsrc *v1alpha1.OperatorSource) error
	Remove(opsrcUID types.UID)
	Watch()
}

type watcher struct {
	sources map[types.UID]types.NamespacedName
	checker UpdateChecker
}

func (w *watcher) Add(opsrc *v1alpha1.OperatorSource) {
	w.sources[opsrc.GetUID()] = types.NamespacedName{
		Namespace: opsrc.GetNamespace(),
		Name:      opsrc.GetName(),
	}
}

func (w *watcher) Remove(opsrcUID types.UID) {
	delete(w.sources, opsrcUID)
}

func (w *watcher) Watch() {
	log.Info("Checking for update(s)")
	for opsrcUID, namespacedName := range w.sources {
		updated, err := w.checker.Check(opsrcUID)
		if err != nil {
			log.Errorf("[poller] error checking for updates [%s] - %v", namespacedName, err)
			continue
		}

		if !updated {
			continue
		}

		deleted, err := w.checker.Trigger(namespacedName)
		if err != nil {
			log.Errorf("[poller] error updating object [%s] - %v", namespacedName, err)
		}

		if deleted {
			log.Infof("[poller] object deleted [%s] - removing from watch list", namespacedName)
		}
	}
}
