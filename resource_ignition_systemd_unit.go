package ignition

import (
	"github.com/coreos/ignition/config/types"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSystemdUnit() *schema.Resource {
	return &schema.Resource{
		Create: resourceSystemdUnitCreate,
		Delete: resourceSystemdUnitDelete,
		Exists: resourceSystemdUnitExists,
		Read:   resourceSystemdUnitRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"mask": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dropin": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"content": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceSystemdUnitCreate(d *schema.ResourceData, meta interface{}) error {
	id, err := buildSystemdUnit(d, meta.(*cache))
	if err != nil {
		return err
	}

	d.SetId(id)
	return nil
}

func resourceSystemdUnitDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourceSystemdUnitExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id, err := buildSystemdUnit(d, meta.(*cache))
	if err != nil {
		return false, err
	}

	return id == d.Id(), nil
}

func resourceSystemdUnitRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func buildSystemdUnit(d *schema.ResourceData, c *cache) (string, error) {
	var dropins []types.SystemdUnitDropIn
	for _, raw := range d.Get("dropin").([]interface{}) {
		value := raw.(map[string]interface{})

		if err := validateUnitContent(value["content"].(string)); err != nil {
			return "", err
		}

		dropins = append(dropins, types.SystemdUnitDropIn{
			Name:     types.SystemdUnitDropInName(value["name"].(string)),
			Contents: value["content"].(string),
		})
	}

	if err := validateUnitContent(d.Get("content").(string)); err != nil {
		if err != errEmptyUnit || (err == errEmptyUnit && len(dropins) == 0) {
			return "", err
		}
	}

	return c.addSystemdUnit(&types.SystemdUnit{
		Name:     types.SystemdUnitName(d.Get("name").(string)),
		Contents: d.Get("content").(string),
		Enable:   d.Get("enable").(bool),
		Mask:     d.Get("mask").(bool),
		DropIns:  dropins,
	}), nil
}
