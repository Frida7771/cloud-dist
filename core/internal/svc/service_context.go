package svc

import (
	"cloud-disk/core/svc"
	appcfg "cloud-disk/internal/config"
)

// ServiceContext is re-exported for legacy imports under core/internal.
type ServiceContext = svc.ServiceContext

// NewServiceContext proxies to the public svc package constructor.
func NewServiceContext(c appcfg.Config) (*ServiceContext, error) {
	return svc.NewServiceContext(c)
}
