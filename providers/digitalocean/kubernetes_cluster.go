// Copyright 2019 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package digitalocean

import (
	"context"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/digitalocean/godo"
)

type KubernetesClusterGenerator struct {
	DigitalOceanService
}

func (g KubernetesClusterGenerator) listKubernetesClusters(ctx context.Context, client *godo.Client) ([]*godo.KubernetesCluster, error) {
	list := []*godo.KubernetesCluster{}

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		clusters, resp, err := client.Kubernetes.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		for _, cluster := range clusters {
			list = append(list, cluster)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return list, nil
}

func (g KubernetesClusterGenerator) createResources(clusterList []*godo.KubernetesCluster) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, cluster := range clusterList {
		resources = append(resources, terraform_utils.NewSimpleResource(
			cluster.ID,
			cluster.Name,
			"digitalocean_kubernetes_cluster",
			"digitalocean",
			[]string{}))
	}
	return resources
}

func (g *KubernetesClusterGenerator) InitResources() error {
	client := g.generateClient()
	output, err := g.listKubernetesClusters(context.TODO(), client)
	if err != nil {
		return err
	}
	g.Resources = g.createResources(output)
	return nil
}
