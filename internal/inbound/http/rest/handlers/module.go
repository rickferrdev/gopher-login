// Module groups inbound adapters.
//
// Inbound adapters are entry points into the application,
// such as HTTP routes, CLI commands, workers, webhooks, or gRPC handlers.
package handlers

import (
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/inbound/http/rest/handlers/user"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"handlers",
	user.Invoke,
)
