// Module provides application services.
//
// Services usually contain use case orchestration and depend on ports,
// not directly on frameworks, databases, queues, or external clients
package services

import (
	"github.com/rickferrdev/go-ports-and-adapters-template/internal/core/services/user"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	user.Provide,
)
