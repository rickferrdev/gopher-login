// Module provides application services.
//
// Services usually contain use case orchestration and depend on ports,
// not directly on frameworks, databases, queues, or external clients
package services

import (
	"github.com/rickferrdev/gopher-login/internal/core/services/user/auth"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	auth.Provide,
)
