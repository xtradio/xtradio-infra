package civocom

import (
	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func DnsSetup(ctx *pulumi.Context, org string) (domain *civo.DnsDomainName, err error) {

	domainName := "ako.sh"
	domain, err = civo.NewDnsDomainName(ctx, domainName, &civo.DnsDomainNameArgs{
		Name: pulumi.StringPtr(domainName),
	}, pulumi.Protect(true))
	if err != nil {
		return
	}

	return
}
