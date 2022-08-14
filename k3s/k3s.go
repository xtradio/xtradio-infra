package k3s

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
)

type chartData struct {
	name   string
	config helm.ChartArgs
}

func Charts(ctx *pulumi.Context, provider *kubernetes.Provider) (err error) {

	o2pNS, err := corev1.NewNamespace(ctx, "oauth2-proxy", &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("oauth2-proxy"),
		},
	}, pulumi.Provider(provider))
	if err != nil {
		return
	}

	charts := []chartData{
		{name: "oauth2-proxy", config: Oauth2Proxy(o2pNS.Metadata.Name())},
		{name: "traefik", config: traefik("kube-system")},
		{name: "sealed-secrets", config: sealedSecrets("kube-system")},
	}

	for _, chart := range charts {
		_, err := helm.NewChart(ctx, chart.name, chart.config)
		if err != nil {
			return err
		}
	}

	return
}
