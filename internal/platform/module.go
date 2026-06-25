package platform

import (
	"github.com/rickferrdev/gopher-login/internal/platform/hasher"
	"github.com/rickferrdev/gopher-login/internal/platform/tokenizer"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"platform",
	tokenizer.Provide,
	hasher.Provide,
)
