// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package api

import (
	"errors"
	"log/slog"

	"github.com/cilium/hive/cell"
	"github.com/go-openapi/runtime/middleware"
	"k8s.io/client-go/discovery"

	"github.com/cilium/cilium/api/v1/models"
	"github.com/cilium/cilium/api/v1/operator/server/restapi/operator"
	k8sClient "github.com/cilium/cilium/pkg/k8s/client"
	"github.com/cilium/cilium/pkg/kvstore"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/option"
)

type kvstoreEnabledFunc func() bool
type isOperatorLeadingFunc func() bool

func HealthHandlerCell(
	kvstoreEnabled kvstoreEnabledFunc,
	isOperatorLeading isOperatorLeadingFunc,
) cell.Cell {
	return cell.Module(
		"health-handler",
		"Operator health HTTP handler",

		cell.Provide(func(
			clientset k8sClient.Clientset,
			logger *slog.Logger,
		) operator.GetHealthzHandler {
			if !clientset.IsEnabled() {
				return &healthHandler{
					enabled: false,
				}
			}

			return &healthHandler{
				enabled:           true,
				isOperatorLeading: isOperatorLeading,
				kvstoreEnabled:    kvstoreEnabled,
				discovery:         clientset.Discovery(),
				log:               logger,
			}
		}),
	)
}

type healthHandler struct {
	enabled           bool
	isOperatorLeading isOperatorLeadingFunc
	kvstoreEnabled    kvstoreEnabledFunc
	discovery         discovery.DiscoveryInterface
	log               *slog.Logger
}

func (h *healthHandler) Handle(params operator.GetHealthzParams) middleware.Responder {
	if !h.enabled {
		return operator.NewGetHealthzNotImplemented()
	}

	if err := h.checkStatus(); err != nil {
		h.log.Warn("Health check status", logfields.Error, err)
		return operator.NewGetHealthzInternalServerError().WithPayload(err.Error())
	}

	return operator.NewGetHealthzOK().WithPayload("ok")
}

// checkStatus verifies the connection status to the kvstore and the
// k8s apiserver and returns an error if any of them is unhealthy
func (h *healthHandler) checkStatus() error {
	// We check if we are the leader here because only the leader has
	// access to the kvstore client. Otherwise, the kvstore client check
	// will block. It is safe for a non-leader to skip this check, as the
	// it is the leader's responsibility to report the status of the
	// kvstore client.
	if h.isOperatorLeading() && h.kvstoreEnabled() {
		client := kvstore.Client()
		if client == nil {
			return errors.New("kvstore client not configured")
		}

		status := client.Status()
		if status.State != models.StatusStateOk &&
			// Don't treat warnings as errors when the support for running
			// etcd in pod network is enabled. This is necessary to allow
			// Cilium turning ready even before connecting to the kvstore,
			// and break the chicken-and-egg dependency during startup.
			!(status.State == models.StatusStateWarning && option.Config.KVstorePodNetworkSupport) {
			return errors.New(status.Msg)
		}
	}
	_, err := h.discovery.ServerVersion()
	return err
}
