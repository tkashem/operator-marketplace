package v1alpha1

import "strings"

// The following list is the set of phases an OperatorSource object can be in
// while it is going through reconciliation process.
//
// Initial --> Validating --> Downloading --> Configuring --> Succeeded
//
// On reconciliation error, given OperatorSource object is transitioned
// into "Failed" phase.
// On successful reconciliation, given OperatorSource object
// is transitioned into "Succeeded" phase.
const (
	// This phase applies to when an OperatorSource object has been created
	// and the Phase attribute is empty.
	OperatorSourcePhaseInitial = ""

	// In this phase we validate the given OperatorSource object.
	OperatorSourcePhaseValidating = "Validating"

	// In this phase, we connect to the specified registry, download
	// available manifest(s) and save them to underlying datastore.
	OperatorSourcePhaseDownloading = "Downloading"

	// In this phase, we make sure that a corresponding
	// CatalogSourceConfig object is created.
	OperatorSourcePhaseConfiguring = "Configuring"

	// This phase indicates that the given OperatorSource object has been
	// successfully reconciled.
	OperatorSourcePhaseSucceeded = "Succeeded"

	// This phase indicates that reconciliation of the given
	// OperatorSource object has failed.
	OperatorSourcePhaseFailed = "Failed"
)

var (
	// Default descriptive message associated with each phase.
	operatorSourcePhaseMessages = map[string]string{
		OperatorSourcePhaseValidating:  "Scheduled for validation",
		OperatorSourcePhaseDownloading: "Scheduled for download of operator manifest(s)",
		OperatorSourcePhaseConfiguring: "Scheduled for configuration",
		OperatorSourcePhaseSucceeded:   "The object has been successfully reconciled",
		OperatorSourcePhaseFailed:      "Reconciliation has failed",
	}
)

// GetOperatorSourcePhaseMessage returns the default message associated with a
// particular phase.
func GetOperatorSourcePhaseMessage(phaseName string) string {
	return operatorSourcePhaseMessages[phaseName]
}

// IsEqual returns true if the Spec specified in this is the same as the other.
// Otherwise, the function returns false.
//
// The function performs a case insensitive comparison of corresponding
// attributes.
//
// If the Spec specified in other is nil then the function returns false.
func (this *OperatorSourceSpec) IsEqual(other *OperatorSourceSpec) bool {
	if other == nil {
		return false
	}

	if strings.EqualFold(this.Endpoint, other.Endpoint) &&
		strings.EqualFold(this.RegistryNamespace, other.RegistryNamespace) &&
		strings.EqualFold(this.Type, other.Type) {
		return true
	}

	return false
}
