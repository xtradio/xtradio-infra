package k3s

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func traefik(ns string) (chartArgs helm.ChartArgs) {

	chartArgs = helm.ChartArgs{
		Namespace: pulumi.String(ns),
		Chart:     pulumi.String("traefik"),
		FetchArgs: helm.FetchArgs{
			Repo: pulumi.String("https://helm.traefik.io/traefik"),
		},

		Version: pulumi.String("10.24.0"),
		Values: pulumi.Map{
			"deployment": pulumi.Map{
				"replicas": pulumi.IntPtr(1),
			},
			"image": pulumi.Map{
				"pullPolicy": pulumi.StringPtr("Always"),
			},
			"providers": pulumi.Map{
				"kubernetesCRD": pulumi.Map{
					"allowCrossNamespace": pulumi.BoolPtr(true),
				},
			},
			"ingressRoute": pulumi.Map{
				"dashboard": pulumi.Map{
					"enabled": pulumi.BoolPtr(false),
				},
			},
			"service": pulumi.Map{
				"annotations": pulumi.Map{
					"kubernetes.civo.com/loadbalancer-algorithm":             pulumi.StringPtr("least_connections"),
					"kubernetes.civo.com/loadbalancer-enable-proxy-protocol": pulumi.StringPtr("send-proxy"),
				},
			},
			// 	"websecure": pulumi.Map{
			// 		"hostPort": pulumi.IntPtr(443),
			// 	},
			// },
		},
	}

	return
}
