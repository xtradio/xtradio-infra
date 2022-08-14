package k3s

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func sealedSecrets(ns string) (chartArgs helm.ChartArgs) {

	chartArgs = helm.ChartArgs{
		Namespace: pulumi.String(ns),
		Chart:     pulumi.String("sealed-secrets"),
		FetchArgs: helm.FetchArgs{
			Repo: pulumi.String("https://bitnami-labs.github.io/sealed-secrets"),
		},

		Version: pulumi.String("2.6.0"),
		Values: pulumi.Map{
			"fullNameOverride": pulumi.StringPtr("sealed-secrets-controller"),
		},
	}

	return
}
