// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clients

import (
	"errors"

	"github.com/openclarity/apiclarity/backend/pkg/modules/internal/core"
	"github.com/openclarity/apiclarity/backend/pkg/modules/internal/fuzzer/config"
	"github.com/openclarity/apiclarity/backend/pkg/modules/internal/fuzzer/logging"
)

type Client interface {
	TriggerFuzzingJob(apiID int64, endpoint string, securityItem string) error
}

//nolint: ireturn,nolintlint
func NewClient(moduleConfig *config.Config, accessor core.BackendAccessor) (Client, error) {
	if moduleConfig.GetDeploymentType() == config.DeploymentTypeKubernetes {
		client, err := NewKubernetesClient(moduleConfig, accessor)
		if err != nil {
			logging.Logf("[Fuzzer] Error, can't create Kubernetes client, err=(%v)", err)
			return nil, err
		}
		logging.Logf("[Fuzzer] Docker client creation, ok")
		return client, nil
	} else if moduleConfig.GetDeploymentType() == config.DeploymentTypeDocker {
		client, err := NewDockerClient(moduleConfig)
		if err != nil {
			logging.Logf("[Fuzzer] Error, can't create Docker client, err=(%v)", err)
			return nil, err
		}
		logging.Logf("[Fuzzer] Docker client creation, ok")
		return client, nil
	} else if moduleConfig.GetDeploymentType() == config.DeploymentTypeFake {
		client, err := NewFakeClient(moduleConfig)
		if err != nil {
			logging.Logf("[Fuzzer] Error, can't create Fake client, err=(%v)", err)
			return nil, err
		}
		logging.Logf("[Fuzzer] Docker Fake creation, ok")
		return client, nil
	}

	// ... Else, not supported
	logging.Errorf("[Fuzzer] unsupported DEPLOYMENT_TYPE = (%v)", moduleConfig.GetDeploymentType())
	return nil, errors.New("not supported")
}