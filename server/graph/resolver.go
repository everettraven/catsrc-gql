package graph

import (
	"github.com/operator-framework/operator-registry/alpha/declcfg"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	catalog *declcfg.DeclarativeConfig
}

func NewResolver(declarativeConfig *declcfg.DeclarativeConfig) *Resolver {
	return &Resolver{catalog: declarativeConfig}
}
