package digitalocean

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pearkes/digitalocean"
)

func resourceDigitalOceanDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanDomainCreate,
		Read:   resourceDigitalOceanDomainRead,
		Delete: resourceDigitalOceanDomainDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDigitalOceanDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*digitalocean.Client)

	// Build up our creation options
	opts := &digitalocean.CreateDomain{
		Name:      d.Get("name").(string),
		IPAddress: d.Get("ip_address").(string),
	}

	log.Printf("[DEBUG] Domain create configuration: %#v", opts)
	name, err := client.CreateDomain(opts)
	if err != nil {
		return fmt.Errorf("Error creating Domain: %s", err)
	}

	d.SetId(name)
	log.Printf("[INFO] Domain Name: %s", name)

	return resourceDigitalOceanDomainRead(d, meta)
}

func resourceDigitalOceanDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*digitalocean.Client)

	domain, err := client.RetrieveDomain(d.Id())
	if err != nil {
		// If the domain is somehow already destroyed, mark as
		// successfully gone
		if strings.Contains(err.Error(), "404 Not Found") {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving domain: %s", err)
	}

	d.Set("name", domain.Name)

	return nil
}

func resourceDigitalOceanDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*digitalocean.Client)

	log.Printf("[INFO] Deleting Domain: %s", d.Id())
	err := client.DestroyDomain(d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting Domain: %s", err)
	}

	d.SetId("")
	return nil
}
