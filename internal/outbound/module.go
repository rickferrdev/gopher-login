package outbound

import "go.uber.org/fx"

// Module groups outbound adapters.
//
// Outbound adapters are implementations that communicate with external systems,
// such as databases, queues, APIs, caches, file systems, or third-party services.
//
// This template intentionally does not provide any concrete outbound adapter.
// Add your own adapters here when your application needs them.
var Module = fx.Module(
	"outbound",
)
