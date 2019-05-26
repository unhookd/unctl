package config

var Current Config
var CurrentProvider ProviderInterface

type EndpointsTable map[string]string

type ContextsTable map[string][]string

type NotificationsTable map[string]string
type NotificationsLookup []NotificationsTable

type TargetTable struct {
	Release       string
	Repo          string
	Namespace     string
	Cluster       string
	Branch        string
	Chart         string
	Version       string
	Notifications NotificationsLookup
}

type ProjectTable map[string]TargetTable
type DeploymentsTable map[string]ProjectTable

type Config struct {
	Contexts    ContextsTable
	Endpoints   EndpointsTable
	Deployments DeploymentsTable
}

func LoadConfig() {
	Current = CurrentProvider.GetConfig()
}
