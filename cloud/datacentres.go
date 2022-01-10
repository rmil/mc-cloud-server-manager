package cloud

import "github.com/hetznercloud/hcloud-go/hcloud"

func (c *Clouder) getDCByName(name string) (hcloud.Datacenter, error) {
	for _, datacentre := range c.dc.Datacentres {
		if datacentre.Name == name {
			return datacentre, nil
		}
	}
	return hcloud.Datacenter{}, ErrDCNotFound
}

func (c *Clouder) getRecommendedDatacentre() (hcloud.Datacenter, error) {
	for _, datacentre := range c.dc.Datacentres {
		if datacentre.ID == c.dc.Recommendation {
			return datacentre, nil
		}
	}
	return hcloud.Datacenter{}, ErrRecommendationNotFound
}
