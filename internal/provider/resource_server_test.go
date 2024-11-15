package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPritunlServer(t *testing.T) {

	t.Run("creates a server with default configuration", func(t *testing.T) {
		serverName := "tfacc-server1"

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { preCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testPritunlServerDestroy,
			Steps: []resource.TestStep{
				{
					Config: testPritunlServerSimpleConfig(serverName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
					),
				},
				// import test
				importStep("pritunl_server.test"),
			},
		})
	})

	t.Run("creates a server with sso_auth attribute", func(t *testing.T) {
		serverName := "tfacc-server1"

		testCase := func(t *testing.T, ssoAuth bool) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config: testPritunlServerConfigWithSsoAuth(serverName, ssoAuth),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							resource.TestCheckResourceAttr("pritunl_server.test", "sso_auth", strconv.FormatBool(ssoAuth)),
						),
					},
					// import test
					importStep("pritunl_server.test"),
				},
			})
		}

		t.Run("with enabled option", func(t *testing.T) {
			testCase(t, true)
		})

		t.Run("with disabled option", func(t *testing.T) {
			testCase(t, false)
		})

		t.Run("without an option", func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config: testPritunlServerSimpleConfig(serverName),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							resource.TestCheckResourceAttr("pritunl_server.test", "sso_auth", "false"),
						),
					},
					// import test
					importStep("pritunl_server.test"),
				},
			})
		})
	})

	t.Run("creates a server with device_auth attribute", func(t *testing.T) {
		serverName := "tfacc-server1"

		testCase := func(t *testing.T, deviceAuth bool) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config: testPritunlServerConfigWithDeviceAuth(serverName, deviceAuth),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							resource.TestCheckResourceAttr("pritunl_server.test", "device_auth", strconv.FormatBool(deviceAuth)),
						),
					},
					// import test
					importStep("pritunl_server.test"),
				},
			})
		}

		t.Run("with enabled option", func(t *testing.T) {
			testCase(t, true)
		})

		t.Run("with disabled option", func(t *testing.T) {
			testCase(t, false)
		})

		t.Run("without an option", func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config: testPritunlServerSimpleConfig(serverName),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							resource.TestCheckResourceAttr("pritunl_server.test", "device_auth", "false"),
						),
					},
					// import test
					importStep("pritunl_server.test"),
				},
			})
		})
	})

	t.Run("creates a server with dynamic_firewall attribute", func(t *testing.T) {
		serverName := "tfacc-server1"

		testCase := func(t *testing.T, dynamicFirewall bool) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config: testPritunlServerConfigWithDynamicFirewall(serverName, dynamicFirewall),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							resource.TestCheckResourceAttr("pritunl_server.test", "dynamic_firewall", strconv.FormatBool(dynamicFirewall)),
						),
					},
					// import test
					importStep("pritunl_server.test"),
				},
			})
		}

		t.Run("with enabled option", func(t *testing.T) {
			testCase(t, true)
		})

		t.Run("with disabled option", func(t *testing.T) {
			testCase(t, false)
		})

		t.Run("without an option", func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config: testPritunlServerSimpleConfig(serverName),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							resource.TestCheckResourceAttr("pritunl_server.test", "dynamic_firewall", "false"),
						),
					},
					// import test
					importStep("pritunl_server.test"),
				},
			})
		})
	})

	t.Run("creates a server with an attached organization", func(t *testing.T) {
		serverName := "tfacc-server1"
		orgName := "tfacc-org1"

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { preCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testPritunlServerDestroy,
			Steps: []resource.TestStep{
				{
					Config: testPritunlServerConfigWithAttachedOrganization(serverName, orgName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						resource.TestCheckResourceAttr("pritunl_organization.test", "name", orgName),

						func(s *terraform.State) error {
							attachedOrganizationId := s.RootModule().Resources["pritunl_server.test"].Primary.Attributes["organization_ids.0"]
							organizationId := s.RootModule().Resources["pritunl_organization.test"].Primary.Attributes["id"]
							if attachedOrganizationId != organizationId {
								return fmt.Errorf("organization_id is invalid or empty")
							}
							return nil
						},
					),
				},
				// import test
				importStep("pritunl_server.test"),
			},
		})
	})

	t.Run("creates a server with a few attached organizations", func(t *testing.T) {
		serverName := "tfacc-server1"
		org1Name := "tfacc-org1"
		org2Name := "tfacc-org2"

		expectedOrganizationIds := make(map[string]struct{})

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { preCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testPritunlServerDestroy,
			Steps: []resource.TestStep{
				{
					Config: testPritunlServerConfigWithAFewAttachedOrganizations(serverName, org1Name, org2Name),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						resource.TestCheckResourceAttr("pritunl_organization.test", "name", org1Name),
						resource.TestCheckResourceAttr("pritunl_organization.test2", "name", org2Name),

						func(s *terraform.State) error {
							attachedOrganization1Id := s.RootModule().Resources["pritunl_server.test"].Primary.Attributes["organization_ids.0"]
							attachedOrganization2Id := s.RootModule().Resources["pritunl_server.test"].Primary.Attributes["organization_ids.1"]
							organization1Id := s.RootModule().Resources["pritunl_organization.test"].Primary.Attributes["id"]
							organization2Id := s.RootModule().Resources["pritunl_organization.test2"].Primary.Attributes["id"]
							expectedOrganizationIds = map[string]struct{}{
								organization1Id: {},
								organization2Id: {},
							}

							if attachedOrganization1Id == attachedOrganization2Id {
								return fmt.Errorf("first and seconds attached organization_id is the same")
							}

							if _, ok := expectedOrganizationIds[attachedOrganization1Id]; !ok {
								return fmt.Errorf("attached organization_id %s doesn't contain in expected organizations list", attachedOrganization1Id)
							}

							if _, ok := expectedOrganizationIds[attachedOrganization2Id]; !ok {
								return fmt.Errorf("attached organization_id %s doesn't contain in expected organizations list", attachedOrganization2Id)
							}

							return nil
						},
					),
				},
				// import test (custom import that ignores order of organization IDs)
				{
					ResourceName: "pritunl_server.test",
					ImportStateCheck: func(states []*terraform.InstanceState) error {
						importedOrganization1Id := states[0].Attributes["organization_ids.0"]
						importedOrganization2Id := states[0].Attributes["organization_ids.1"]

						if _, ok := expectedOrganizationIds[importedOrganization1Id]; !ok {
							return fmt.Errorf("imported organization_id %s doesn't contain in expected organizations list", importedOrganization1Id)
						}

						if _, ok := expectedOrganizationIds[importedOrganization2Id]; !ok {
							return fmt.Errorf("imported organization_id %s doesn't contain in expected organizations list", importedOrganization2Id)
						}

						return nil
					},
					ImportState:       true,
					ImportStateVerify: false,
				},
			},
		})
	})

	t.Run("creates a server with error", func(t *testing.T) {
		t.Run("due to an invalid network", func(t *testing.T) {
			serverName := "tfacc-server1"
			port := 11111
			missedSubnetNetwork := "10.100.0.2"
			invalidNetwork := "10.100.0"

			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config:      testGetServerConfigWithNetworkAndPort(serverName, missedSubnetNetwork, port),
						ExpectError: regexp.MustCompile(fmt.Sprintf("invalid CIDR address: %s", missedSubnetNetwork)),
					},
					{
						Config:      testGetServerConfigWithNetworkAndPort(serverName, invalidNetwork, port),
						ExpectError: regexp.MustCompile(fmt.Sprintf("invalid CIDR address: %s", invalidNetwork)),
					},
				},
			})
		})

		t.Run("due to an unsupported network", func(t *testing.T) {
			serverName := "tfacc-server1"
			port := 11111
			unsupportedNetwork := "172.14.68.0/24"
			supportedNetwork := "172.16.68.0/24"

			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config:      testGetServerConfigWithNetworkAndPort(serverName, unsupportedNetwork, port),
						ExpectError: regexp.MustCompile(fmt.Sprintf("provided subnet %s does not belong to expected subnets 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16", unsupportedNetwork)),
					},
					{
						Config: testGetServerConfigWithNetworkAndPort(serverName, supportedNetwork, port),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							resource.TestCheckResourceAttr("pritunl_server.test", "network", supportedNetwork),
						),
					},
				},
			})
		})

		t.Run("due to an invalid bind_address", func(t *testing.T) {
			serverName := "tfacc-server1"
			network := "172.16.68.0/24"
			port := 11111
			invalidBindAddress := "10.100.0.1/24"
			correctBindAddress := "10.100.0.1"

			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config:      testGetServerConfigWithBindAddress(serverName, network, invalidBindAddress, port),
						ExpectError: regexp.MustCompile(fmt.Sprintf("expected bind_address to contain a valid IP, got: %s", invalidBindAddress)),
					},
					{
						Config: testGetServerConfigWithBindAddress(serverName, network, correctBindAddress, port),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							resource.TestCheckResourceAttr("pritunl_server.test", "bind_address", correctBindAddress),
						),
					},
				},
			})
		})
	})

	t.Run("creates a server with groups attribute", func(t *testing.T) {
		serverName := "tfacc-server1"

		t.Run("with correct group name", func(t *testing.T) {
			correctGroupName := "Group-Has-No-Spaces"
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config: testPritunlServerConfigWithGroups(serverName, correctGroupName),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),

							func(s *terraform.State) error {
								groupName := s.RootModule().Resources["pritunl_server.test"].Primary.Attributes["groups.0"]
								if groupName != correctGroupName {
									return fmt.Errorf("group name mismatch")
								}

								return nil
							},
						),
					},
					// import test
					importStep("pritunl_server.test"),
				},
			})
		})

		t.Run("with invalid group name", func(t *testing.T) {
			invalidGroupName := "Group Name With Spaces"
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { preCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testPritunlServerDestroy,
				Steps: []resource.TestStep{
					{
						Config:      testPritunlServerConfigWithGroups(serverName, invalidGroupName),
						ExpectError: regexp.MustCompile(fmt.Sprintf("%s contains spaces", invalidGroupName)),
					},
				},
			})
		})
	})
}

func testPritunlServerSimpleConfig(name string) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name    = "%[1]s"
		}
	`, name)
}

func testPritunlServerConfigWithSsoAuth(name string, ssoAuth bool) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name     = "%[1]s"
			sso_auth = %[2]v
		}
	`, name, ssoAuth)
}

func testPritunlServerConfigWithDeviceAuth(name string, deviceAuth bool) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name     = "%[1]s"
			device_auth = %[2]v
		}
	`, name, deviceAuth)
}

func testPritunlServerConfigWithDynamicFirewall(name string, dynamicFirewall bool) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name     = "%[1]s"
			dynamic_firewall = %[2]v
		}
	`, name, dynamicFirewall)
}

func testPritunlServerConfigWithAttachedOrganization(name, organizationName string) string {
	return fmt.Sprintf(`
		resource "pritunl_organization" "test" {
			name    = "%[2]s"
		}

		resource "pritunl_server" "test" {
			name    = "%[1]s"
			organization_ids = [
				pritunl_organization.test.id
			]
		}
	`, name, organizationName)
}

func testPritunlServerConfigWithAFewAttachedOrganizations(name, organization1Name, organization2Name string) string {
	return fmt.Sprintf(`
		resource "pritunl_organization" "test" {
			name    = "%[2]s"
		}

		resource "pritunl_organization" "test2" {
			name    = "%[3]s"
		}

		resource "pritunl_server" "test" {
			name    = "%[1]s"
			organization_ids = [
				pritunl_organization.test.id,
				pritunl_organization.test2.id
			]
		}
	`, name, organization1Name, organization2Name)
}

func testGetServerConfigWithNetworkAndPort(name, network string, port int) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name    = "%[1]s"
			network  =  "%[2]s"
			port     = %[3]d
			protocol = "tcp"
		}
	`, name, network, port)
}

func testGetServerConfigWithBindAddress(name, network, bindAddress string, port int) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name    		= "%[1]s"
			network  		= "%[2]s"
			bind_address  	= "%[3]s"
			port     		= %[4]d
			protocol 		= "tcp"
		}
	`, name, network, bindAddress, port)
}

func testPritunlServerConfigWithGroups(name string, groupName string) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name    = "%[1]s"
			groups    = ["%[2]s"]
		}
	`, name, groupName)
}

func testPritunlServerDestroy(s *terraform.State) error {
	serverId := s.RootModule().Resources["pritunl_server.test"].Primary.Attributes["id"]

	servers, err := testClient.GetServers()
	if err != nil {
		return err
	}
	for _, server := range servers {
		if server.ID == serverId {
			return fmt.Errorf("a server is not destroyed")
		}
	}
	return nil
}
