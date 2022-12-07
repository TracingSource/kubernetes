package view

import "sync"

var (
	exportersMu sync.RWMutex // guards exporters
	exporters   = make(map[Exporter]struct{})
)

// Exporter exports the collected records as view data.
//
// The ExportView method should return quickly; if an
// Exporter takes a significant amount of time to
// process a Data, that work should be done on another goroutine.
//
// It is safe to assume that ExportView will not be called concurrently from
// multiple goroutines.
//
// The Data should not be modified.
type Exporter interface {
	ExportView(viewData *Data)
}

// RegisterExporter registers an exporter.
// Collected data will be reported via all the
// registered exporters. Once you no longer
// want data to be exported, invoke UnregisterExporter
// with the previously registered exporter.
//
// Binaries can register exporters, libraries shouldn't register exporters.
func RegisterExporter(e Exporter) {
	exportersMu.Lock()
	defer exportersMu.Unlock()

	exporters[e] = struct{}{}
}

// UnregisterExporter unregisters an exporter.
func UnregisterExporter(e Exporter) {
	exportersMu.Lock()
	defer exportersMu.Unlock()

	delete(exporters, e)
}
