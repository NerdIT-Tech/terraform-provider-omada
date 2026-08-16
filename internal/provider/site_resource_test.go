package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSiteResource(t *testing.T) {
	resourceName := "omada_site.test"
	nameCreate := fmt.Sprintf("tf-acc-test-%s", "create")
	nameUpdate := fmt.Sprintf("tf-acc-test-%s", "update")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing.
			{
				Config: testAccSiteResourceConfig(nameCreate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameCreate),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// ImportState testing.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing.
			{
				Config: testAccSiteResourceConfig(nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", nameUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing happens automatically via CheckDestroy-less cleanup.
		},
	})
}

func testAccSiteResourceConfig(name string) string {
	return fmt.Sprintf(`
resource "omada_site" "test" {
  name = %[1]q
}
`, name)
}
