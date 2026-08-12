/*
Vendor extensions to the measurand and unit enums.

Chargers in the field ship sampled values the 1.6 spec does not define -- VMOTO
reports its tariff totals as measurand "elecAmt" in unit "RMB", for instance.
Because sampled values are validated when the message is parsed, a single
unknown string fails the entire enclosing Call: a StopTransaction carrying such
a value is answered with a CallError, the handler never runs, and the station
never sees an acknowledgement for a session it has already ended.

Rather than loosen validation for everyone, a consumer declares the values its
hardware actually sends and only those are accepted in addition to the spec's.
Register during init, before any message is parsed.
*/
package types

import "sync"

var (
	vendorMu             sync.RWMutex
	vendorMeasurands     = map[Measurand]struct{}{}
	vendorUnitsOfMeasure = map[UnitOfMeasure]struct{}{}
)

// RegisterMeasurand accepts a vendor-specific measurand in sampled values,
// alongside the ones OCPP 1.6 defines.
func RegisterMeasurand(measurands ...Measurand) {
	vendorMu.Lock()
	defer vendorMu.Unlock()
	for _, m := range measurands {
		vendorMeasurands[m] = struct{}{}
	}
}

// RegisterUnitOfMeasure accepts a vendor-specific unit in sampled values,
// alongside the ones OCPP 1.6 defines.
func RegisterUnitOfMeasure(units ...UnitOfMeasure) {
	vendorMu.Lock()
	defer vendorMu.Unlock()
	for _, u := range units {
		vendorUnitsOfMeasure[u] = struct{}{}
	}
}

func isVendorMeasurand(measurand Measurand) bool {
	vendorMu.RLock()
	defer vendorMu.RUnlock()
	_, ok := vendorMeasurands[measurand]
	return ok
}

func isVendorUnitOfMeasure(unit UnitOfMeasure) bool {
	vendorMu.RLock()
	defer vendorMu.RUnlock()
	_, ok := vendorUnitsOfMeasure[unit]
	return ok
}
