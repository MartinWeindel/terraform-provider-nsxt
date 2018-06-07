/* Copyright © 2018 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/go-vmware-nsxt"
	"net/http"
	"testing"
)

func TestAccResourceNsxtLbSourceIpPersistenceProfile_basic(t *testing.T) {
	name := "test-nsx-persistence-profile"
	updatedName := fmt.Sprintf("%s-update", name)
	testResourceName := "nsxt_lb_source_ip_persistence_profile.test"
	timeout := "100"
	updatedTimeout := "200"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLbSourceIpPersistenceProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLbSourceIpPersistenceProfileBasicTemplate(name, timeout),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLbSourceIpPersistenceProfileExists(name, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "test description"),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", "true"),
					resource.TestCheckResourceAttr(testResourceName, "ha_persistence_mirroring", "true"),
					resource.TestCheckResourceAttr(testResourceName, "purge_when_full", "false"),
					resource.TestCheckResourceAttr(testResourceName, "timeout", timeout),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
			{
				Config: testAccNSXLbSourceIpPersistenceProfileBasicTemplate(updatedName, updatedTimeout),
				Check: resource.ComposeTestCheckFunc(
					testAccNSXLbSourceIpPersistenceProfileExists(updatedName, testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "test description"),
					resource.TestCheckResourceAttr(testResourceName, "persistence_shared", "true"),
					resource.TestCheckResourceAttr(testResourceName, "ha_persistence_mirroring", "true"),
					resource.TestCheckResourceAttr(testResourceName, "purge_when_full", "false"),
					resource.TestCheckResourceAttr(testResourceName, "timeout", updatedTimeout),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceNsxtLbSourceIpPersistenceProfile_importBasic(t *testing.T) {
	name := "test-nsx-persistence-profile"
	testResourceName := "nsxt_lb_source_ip_persistence_profile.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNSXLbSourceIpPersistenceProfileCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNSXLbSourceIpPersistenceProfileCreateTemplateTrivial(name),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNSXLbSourceIpPersistenceProfileExists(displayName string, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		nsxClient := testAccProvider.Meta().(*nsxt.APIClient)
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("NSX LB source ip persistence profile resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		if resourceID == "" {
			return fmt.Errorf("NSX LB source ip persistence profile resource ID not set in resources ")
		}

		profile, responseCode, err := nsxClient.ServicesApi.ReadLoadBalancerSourceIpPersistenceProfile(nsxClient.Context, resourceID)
		if err != nil {
			return fmt.Errorf("Error while retrieving LB source ip persistence profile with ID %s. Error: %v", resourceID, err)
		}

		if responseCode.StatusCode != http.StatusOK {
			return fmt.Errorf("Error while checking if LB source ip persistence profile %s exists. HTTP return code was %d", resourceID, responseCode.StatusCode)
		}

		if displayName == profile.DisplayName {
			return nil
		}
		return fmt.Errorf("NSX LB source ip persistence profile %s wasn't found", displayName)
	}
}

func testAccNSXLbSourceIpPersistenceProfileCheckDestroy(state *terraform.State, displayName string) error {
	nsxClient := testAccProvider.Meta().(*nsxt.APIClient)
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_lb_source_ip_persistence_profile" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		profile, responseCode, err := nsxClient.ServicesApi.ReadLoadBalancerSourceIpPersistenceProfile(nsxClient.Context, resourceID)
		if err != nil {
			if responseCode.StatusCode != http.StatusOK {
				return nil
			}
			return fmt.Errorf("Error while retrieving LB source ip persistence profile with ID %s. Error: %v", resourceID, err)
		}

		if displayName == profile.DisplayName {
			return fmt.Errorf("NSX LB source ip persistence profile %s still exists", displayName)
		}
	}
	return nil
}

func testAccNSXLbSourceIpPersistenceProfileBasicTemplate(name string, timeout string) string {
	return fmt.Sprintf(`
resource "nsxt_lb_source_ip_persistence_profile" "test" {
  display_name             = "%s"
  description             = "test description"
  persistence_shared       = "true"
  ha_persistence_mirroring = "true"
  purge_when_full          = "false"
  timeout                  = "%s"

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}
`, name, timeout)
}

func testAccNSXLbSourceIpPersistenceProfileCreateTemplateTrivial(name string) string {
	return fmt.Sprintf(`
resource "nsxt_lb_source_ip_persistence_profile" "test" {
  display_name = "%s"
  description  = "test description"
}
`, name)
}
