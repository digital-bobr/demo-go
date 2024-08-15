package router

import (
	"encoding/json"
	"fmt"
	"service-status-reporter/pkg/client"
	"service-status-reporter/pkg/model"
	"sync"
)

const (
	healthEndpoint = "/v1/healthcheck/info"
)

// healthCheck response body -- START

type AppInfo struct {
	ServiceName  string   `json:"serviceName"`
	Version      string   `json:"version"`
	ChartVersion string   `json:"chartVersion"`
	GoVersion    string   `json:"goVersion"`
	Codebase     Codebase `json:"codebase"`
	Runtime      Runtime  `json:"runtime"`
}

type Codebase struct {
	Repo    string `json:"repo"`
	Ref     string `json:"ref"`
	Commit  string `json:"commit"`
	Release string `json:"release"`
}

type Runtime struct {
	Environment string `json:"environment"`
	Hostname    string `json:"hostname"`
	Cluster     string `json:"cluster"`
	Namespace   string `json:"namespace"`
}

type HealthResponseBody struct {
	App AppInfo `json:"app"`
}

// healthCheck response body -- END

type ServiceHealth struct {
	Production    Environment `json:"production"`
	Preproduction Environment `json:"preproduction"`
	Nonproduction Environment `json:"nonproduction"`
}

type Environment struct {
	Version string `json:"version"`
	Ref     string `json:"ref"`
	Commit  string `json:"commit"`
	Error   string `json:"error"`
}

func getServiceHealthInfo(s model.Service) ServiceHealth {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var serviceHealth ServiceHealth

	envs := map[string]*Environment{
		s.ProdEnv:    &serviceHealth.Production,
		s.PreprodEnv: &serviceHealth.Preproduction,
		s.NonprodEnv: &serviceHealth.Nonproduction,
	}

	for env, health := range envs {
		wg.Add(1)
		go func(env string, health *Environment) {
			defer wg.Done()
			healthResponseBody, err := getHealthInfoByEnv(env)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				health.Error = err.Error()
			} else {
				health.Version = healthResponseBody.App.Version
				health.Commit = healthResponseBody.App.Codebase.Commit
				health.Ref = healthResponseBody.App.Codebase.Ref
			}
		}(env, health)
	}

	wg.Wait()
	return serviceHealth
}

func getHealthInfoByEnv(baseUrl string) (HealthResponseBody, error) {
	req, err := client.NewRequestBuilder().
		SetMethod(client.GET{}).
		SetURLTemplate(baseUrl+"{endpoint}").
		AddPathParam("endpoint", healthEndpoint).
		AddHeader("Content-Type", "application/json").
		Build()

	if err != nil {
		return HealthResponseBody{}, fmt.Errorf("error building request: %v", err)
	}

	statusCode, body, err := client.SendRequest(req)
	if err != nil {
		return HealthResponseBody{}, fmt.Errorf("error executing request: %v", err)
	}
	if statusCode != 200 {
		return HealthResponseBody{}, fmt.Errorf("unexpected status code: %d, body: %s", statusCode, string(body))
	}
	var healthResponse HealthResponseBody
	if err = json.Unmarshal(body, &healthResponse); err != nil {
		return HealthResponseBody{}, err
	}
	return healthResponse, nil
}
