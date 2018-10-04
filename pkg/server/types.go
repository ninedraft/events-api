package server

import (
	"time"

	"github.com/containerum/kube-client/pkg/model"
	"github.com/gin-gonic/gin"
)

type EventsActions interface {
	GetNamespaceChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetAllNamespacesChanges(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	GetDeploymentChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetNamespaceDeploymentsChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetAllNamespacesDeploymentsChanges(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	GetServiceChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetNamespaceServicesChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetAllNamespacesServicesChanges(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	GetIngressChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetNamespaceIngressesChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetAllNamespacesIngressesChanges(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	GetPVCChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetNamespacePVCsChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetAllNamespacesPVCsChanges(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	GetSecretChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetNamespaceSecretsChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetAllNamespacesSecretsChanges(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	GetConfigMapChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetNamespaceConfigMapsChanges(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetAllNamespacesConfigMapsChanges(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	GetPodEvents(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetNamespacePodsEvents(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetAllNamespacesPodsEvents(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	GetPVCEvents(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetNamespacePVCsEvents(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetAllNamespacesPVCsEvents(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	GetAllNodesEvents(params gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	GetUsersEvents(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)
	GetSystemEvents(_ gin.Params, limit int, startTime time.Time) (*model.EventsList, error)

	AddUserEvent(event model.Event) error
	AddSystemEvent(event model.Event) error
}
