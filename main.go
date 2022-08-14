package main

import (
	"fmt"

	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	civocom "github.com/xtradio/xtradio-infra/civo"
	k3sCharts "github.com/xtradio/xtradio-infra/k3s"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		org := "xtradio"
		region := "FRA1"

		k3s, err := civocom.K3sCluster(ctx, org, region)
		if err != nil {
			return err
		}

		domain, err := civocom.DnsSetup(ctx, org)
		if err != nil {
			return err
		}

		k3sProvider, err := kubernetes.NewProvider(ctx, "civo", &kubernetes.ProviderArgs{
			Kubeconfig: k3s.Kubeconfig,
		}, pulumi.DependsOn([]pulumi.Resource{k3s}))
		if err != nil {
			return err
		}

		err = k3sCharts.Charts(ctx, k3sProvider) // TODO: Read new load balancer hostname automatically
		if err != nil {
			return err
		}

		lbName := fmt.Sprintf("%s-kube-system-traefik", org)
		loadBalancer, err := civo.GetLoadBalancer(ctx, &civo.GetLoadBalancerArgs{
			Name:   pulumi.StringRef(lbName),
			Region: pulumi.StringRef(region),
		}, nil)
		if err != nil {
			return err
		}

		// Domain which points to the k3s cluster
		_, err = civo.NewDnsDomainRecord(ctx, "wildcard", &civo.DnsDomainRecordArgs{
			DomainId: domain.ID(),
			Value:    pulumi.String(loadBalancer.PublicIp),
			Name:     pulumi.StringPtr("*.k3s"),
			Ttl:      pulumi.Int(600), // Expected ttl to be in the range (600 - 3600)
			Type:     pulumi.String("A"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
