module github.com/kwoodson/azurevmtypes

go 1.15

require (
	github.com/Azure/azure-sdk-for-go v55.6.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.19 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.8
	sigs.k8s.io/cluster-api-provider-azure v0.0.0-20210709011253-2bb61139dd4d

)

replace (
	sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v0.4.0
	// sigs.k8s.io/cluster-api-provider-azure => sigs.k8s.io/cluster-api-provider-azure v0.0.0-20210709011253-2bb61139dd4d
	sigs.k8s.io/cluster-api/test => sigs.k8s.io/cluster-api/test v0.4.0

)
