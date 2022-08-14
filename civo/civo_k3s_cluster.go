package civocom

import (
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	util "github.com/xtradio/xtradio-infra/util"
)

func K3sCluster(ctx *pulumi.Context, org string, region string) (k3s *civo.KubernetesCluster, err error) {

	k3sClusterName := org
	k3sClusterNode := "g3.k3s.medium"

	// Firewall and rules
	firewall, err := civo.NewFirewall(ctx, org, &civo.FirewallArgs{
		Name:               pulumi.StringPtr(org),
		Region:             pulumi.StringPtr(region),
		CreateDefaultRules: pulumi.BoolPtr(false),
	})
	if err != nil {
		return
	}

	firewallRules := []civo.FirewallRuleArgs{
		{Direction: pulumi.String("egress"), Action: pulumi.String("allow"), Cidrs: pulumi.StringArray{pulumi.String("0.0.0.0/0")}, Protocol: pulumi.String("tcp"), StartPort: pulumi.String("0"), EndPort: pulumi.String("65535")},
		{Direction: pulumi.String("egress"), Action: pulumi.String("allow"), Cidrs: pulumi.StringArray{pulumi.String("0.0.0.0/0")}, Protocol: pulumi.String("udp"), StartPort: pulumi.String("0"), EndPort: pulumi.String("65535")},
		{Direction: pulumi.String("egress"), Action: pulumi.String("allow"), Cidrs: pulumi.StringArray{pulumi.String("0.0.0.0/0")}, Protocol: pulumi.String("icmp")},
		{Direction: pulumi.String("ingress"), Action: pulumi.String("allow"), Cidrs: pulumi.StringArray{pulumi.String(util.MyPublicIP())}, Protocol: pulumi.String("tcp"), StartPort: pulumi.String("6443"), EndPort: pulumi.String("6443")},
		{Direction: pulumi.String("ingress"), Action: pulumi.String("allow"), Cidrs: pulumi.StringArray{pulumi.String("0.0.0.0/0")}, Protocol: pulumi.String("tcp"), StartPort: pulumi.String("80"), EndPort: pulumi.String("80")},
		{Direction: pulumi.String("ingress"), Action: pulumi.String("allow"), Cidrs: pulumi.StringArray{pulumi.String("0.0.0.0/0")}, Protocol: pulumi.String("tcp"), StartPort: pulumi.String("443"), EndPort: pulumi.String("443")},
	}

	for k, rule := range firewallRules {
		fwRuleName := fmt.Sprintf("%s-%s", org, strconv.Itoa(k))
		_, err = civo.NewFirewallRule(ctx, fwRuleName, &civo.FirewallRuleArgs{
			Direction:  rule.Direction,
			Action:     rule.Action,
			Cidrs:      rule.Cidrs,
			Protocol:   rule.Protocol,
			StartPort:  rule.StartPort,
			EndPort:    rule.EndPort,
			FirewallId: firewall.ID(),
			Region:     pulumi.StringPtr(region),
		})
		if err != nil {
			return
		}
	}

	// _, err = civo.NewFirewallRule(ctx, org, &civo.FirewallRuleArgs{})

	k3s, err = civo.NewKubernetesCluster(ctx, k3sClusterName, &civo.KubernetesClusterArgs{
		Name: pulumi.StringPtr(k3sClusterName),
		Cni:  pulumi.StringPtr("cilium"),
		Pools: civo.KubernetesClusterPoolsArgs{
			NodeCount: pulumi.Int(1),
			Size:      pulumi.String(k3sClusterNode),
		},
		FirewallId: firewall.ID(),
		Region:     pulumi.StringPtr(region), // Use the region code
	})

	if err != nil {
		return
	}

	return
}
