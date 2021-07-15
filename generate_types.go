package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-30/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/kwoodson/azurevmtypes/pkg/azurevms"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/resourceskus"
)

type MyInt struct {
	auth autorest.Authorizer
}

func convertSkuToMap(skus []compute.ResourceSku) map[string]compute.ResourceSku {
	mapSkus := map[string]compute.ResourceSku{}
	for _, s := range skus {
		mapSkus[to.String(s.Name)] = s
	}
	return mapSkus
}

func main() {
	// create a network service client and call the
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		panic(err)
	}
	m := MyInt{
		auth: authorizer,
	}

	cli := compute.NewResourceSkusClient(m.SubscriptionID())
	cli.Authorizer = m.Authorizer()

	results, err := cli.List(context.TODO(), "")
	if err != nil {
		panic(err)
	}

	mapSkus := convertSkuToMap(results.Values())

	cache, err := resourceskus.GetCache(m, "east-us")
	if err != nil {
		panic(err)
	}
	keys := make([]string, 0)
	ctx := context.TODO()
	for _, machineType := range azurevms.InstanceTypes {
		keys = append(keys, machineType.InstanceType)
		s, err := cache.Get(ctx, machineType.InstanceType, resourceskus.VirtualMachines)
		if err != nil {
			// resource sku with name 'Standard_D2as_v3' and category 'virtualMachines' not found in location 'east-us'
			if strings.Contains(err.Error(), "not found in location") {
				fmt.Println("Setting to false for type", machineType.InstanceType)
				continue
			} else {
				panic(err)
			}
		}
		if _, ok := mapSkus[to.String(s.Name)]; !ok {
			continue
		}
		machineType.AcceleratedNetworking = s.HasCapability(resourceskus.AcceleratedNetworking)
		azurevms.InstanceTypes[machineType.InstanceType] = machineType
	}

	// sort results
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	instType := `
	"{{.InstanceType}}": {
		InstanceType:          "{{.InstanceType}}",
		VCPU:                  {{.VCPU}},
		MemoryMb:              {{.MemoryMb}},
		GPU:                   {{.GPU}},
		AcceleratedNetworking: {{.AcceleratedNetworking}},
	},`
	f, err := os.Create("azure_instance_types.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	t, err := template.New("instancetypes").Parse(instType)
	if err != nil {
		panic(err)
	}

	for _, k := range keys {
		t.Execute(f, azurevms.InstanceTypes[k])
	}
}

/*
SubscriptionID() string
	ClientID() string
	ClientSecret() string
	CloudEnvironment() string
	TenantID() string
	BaseURI() string
	Authorizer() autorest.Authorizer
	HashKey() string
*/
func (m MyInt) SubscriptionID() string {
	return os.Getenv("AZURE_SUBSCRIPTION_ID")
}
func (m MyInt) ClientSecret() string {
	return ""
}
func (m MyInt) CloudEnvironment() string {
	return "AzurePublicCloud"
}
func (m MyInt) TenantID() string {
	return ""
}

func (m MyInt) ClientID() string {
	return ""
}
func (m MyInt) BaseURI() string {
	return "https://management.azure.com"
}
func (m MyInt) Authorizer() autorest.Authorizer {
	return m.auth
}
func (m MyInt) HashKey() string {
	return ""
}
