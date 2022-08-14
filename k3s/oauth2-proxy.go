package k3s

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Oauth2Proxy(ns pulumi.StringPtrOutput) (chartArgs helm.ChartArgs) {

	chartArgs = helm.ChartArgs{
		Namespace: ns.Elem(),
		Chart:     pulumi.String("oauth2-proxy"),
		FetchArgs: helm.FetchArgs{
			Repo: pulumi.String("https://oauth2-proxy.github.io/manifests"),
		},

		Version: pulumi.String("6.2.2"),
		Values: pulumi.Map{
			"image": pulumi.Map{
				"pullPolicy": pulumi.StringPtr("Always"),
			},
		},
	}

	return
}
