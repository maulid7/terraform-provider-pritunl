package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPritunlRoute(t *testing.T) {

	t.Run("creates multiple route on a test server", func(t *testing.T) {
		serverName := "tfacc-server4"
		routes := map[string]interface{}{"test1": "1.1.1.1/32", "test2": "2.2.2.2/32"}
		addedRoute :=  map[string]interface{}{"test1": "1.1.1.1/32", "test2": "2.2.2.2/32", "test3": "8.8.8.8/32"}
		removedRoute := map[string]interface{}{"test2": "2.2.2.2/32", "test3": "8.8.8.8/32"}

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				preCheck(t)
			},
			ProviderFactories: providerFactories,
			CheckDestroy:      testResourceDestroy("pritunl_server"),
			Steps: []resource.TestStep{
				{
					Config: testPritunlMultipleRoute(serverName, routes),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						resource.TestCheckResourceAttr("pritunl_route.test1", "network", routes["test1"].(string)),
						resource.TestCheckResourceAttr("pritunl_route.test2", "network", routes["test2"].(string)),
					),
				},
				{
					Config: testPritunlMultipleRoute(serverName, addedRoute),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						resource.TestCheckResourceAttr("pritunl_route.test1", "network", addedRoute["test1"].(string)),
						resource.TestCheckResourceAttr("pritunl_route.test2", "network", addedRoute["test2"].(string)),
						resource.TestCheckResourceAttr("pritunl_route.test3", "network", addedRoute["test3"].(string)),
					),
				},
				{
					Config: testPritunlMultipleRoute(serverName, removedRoute),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						resource.TestCheckResourceAttr("pritunl_route.test2", "network", removedRoute["test2"].(string)),
						resource.TestCheckResourceAttr("pritunl_route.test3", "network", removedRoute["test3"].(string)),
						testResourceDestroy("pritunl_route.test1"),
					),
				},
			},
		})
	})

	t.Run("creates a route on a test server", func(t *testing.T) {
		serverName := "tfacc-server1"
		route := "8.8.8.8/32"

		resource.Test(t, resource.TestCase{
			PreCheck: func() { 
				preCheck(t)
			},
			ProviderFactories: providerFactories,
			CheckDestroy:      testResourceDestroy("pritunl_server"),
			Steps: []resource.TestStep{
				{
					Config: testPritunlRouteSimpleConfig(serverName, route),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						resource.TestCheckResourceAttr("pritunl_route.test", "network", route),
						resource.TestCheckResourceAttrPair("pritunl_route.test", "server_id", "pritunl_server.test", "id"),
					),
				},
				// import test
				pritunlRouteImportStep("pritunl_route.test"),
				{
					Config: testPritunlServerWithoutRouteConfig(serverName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						testResourceDestroy("pritunl_route"),
					),
				},
			},
		})
	})

	t.Run("creates a route on a test server with comment", func(t *testing.T) {
		serverName := "tfacc-server2"
		route := "1.1.1.1/32"
		comment := "aaaa"
		updatedComment := "testing-again"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				preCheck(t)
			},
			ProviderFactories: providerFactories,
			CheckDestroy:      testResourceDestroy("pritunl_server"),
			Steps: []resource.TestStep{
				{
					Config: testPritunlRouteWithComment(serverName, route, comment),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						resource.TestCheckResourceAttr("pritunl_route.test", "network", route),
						resource.TestCheckResourceAttr("pritunl_route.test", "comment", comment),
						resource.TestCheckResourceAttrPair("pritunl_route.test", "server_id", "pritunl_server.test", "id"),
					),
				},
				// import test
				pritunlRouteImportStep("pritunl_route.test"),
				{
					// PreConfig: func() {
					// 	fmt.Println("called")
					// 	resetInitOnce()
					// },
					Config: testPritunlRouteWithComment(serverName, route, updatedComment),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						resource.TestCheckResourceAttr("pritunl_route.test", "network", route),
						resource.TestCheckResourceAttr("pritunl_route.test", "comment", updatedComment),
						resource.TestCheckResourceAttrPair("pritunl_route.test", "server_id", "pritunl_server.test", "id"),
					),
				},
			},
		})
	})

	t.Run("creates a route on a test server with nat attribute", func(t *testing.T) {
		serverName := "tfacc-server3"
		route := "1.1.1.1/32"

		testCase := func(t *testing.T, nat bool) {
			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					preCheck(t)
				},
				ProviderFactories: providerFactories,
				CheckDestroy:      testResourceDestroy("pritunl_server"),
				Steps: []resource.TestStep{
					{
						Config: testPritunlRouteWithNat(serverName, route, nat),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							resource.TestCheckResourceAttr("pritunl_route.test", "network", route),
							resource.TestCheckResourceAttr("pritunl_route.test", "nat", strconv.FormatBool(nat)),
							resource.TestCheckResourceAttrPair("pritunl_route.test", "server_id", "pritunl_server.test", "id"),
						),
					},
					// import test
					pritunlRouteImportStep("pritunl_route.test"),
					{
						Config: testPritunlServerWithoutRouteConfig(serverName),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							testResourceDestroy("pritunl_route"),
						),
					},
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
				PreCheck: func() {
					preCheck(t)
				},
				ProviderFactories: providerFactories,
				CheckDestroy:      testResourceDestroy("pritunl_server"),
				Steps: []resource.TestStep{
					{
						Config: testPritunlRouteSimpleConfig(serverName, route),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							resource.TestCheckResourceAttr("pritunl_route.test", "network", route),
							resource.TestCheckResourceAttr("pritunl_route.test", "nat", "true"),
							resource.TestCheckResourceAttrPair("pritunl_route.test", "server_id", "pritunl_server.test", "id"),
						),
					},
					// import test
					pritunlRouteImportStep("pritunl_route.test"),
					{
						Config: testPritunlServerWithoutRouteConfig(serverName),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
							testResourceDestroy("pritunl_route"),
						),
					},
				},
			})
		})
	})


	t.Run("creates invalid route on a server", func(t *testing.T) {
		serverName := "tfacc-server5"
		invalidNetwork := "10.100.1"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				preCheck(t)
			},
			ProviderFactories: providerFactories,
			CheckDestroy:      testResourceDestroy("pritunl_server"),
			Steps: []resource.TestStep{
				{
					Config:      testPritunlRouteSimpleConfig(serverName, invalidNetwork),
					ExpectError: regexp.MustCompile(fmt.Sprintf("invalid CIDR address: %s", invalidNetwork)),
				},
				{
					Config: testPritunlServerWithoutRouteConfig(serverName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						testResourceDestroy("pritunl_route"),
					),
				},
			},
		})
	})

	t.Run("recreate route on a server", func(t *testing.T) {
		serverName := "tfacc-server6"
		network := "10.100.10.1/32"
		update := "10.100.10.2/32"

		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				preCheck(t)
			},
			ProviderFactories: providerFactories,
			CheckDestroy:      testResourceDestroy("pritunl_server"),
			Steps: []resource.TestStep{
				{
					Config: testPritunlRouteSimpleConfig(serverName, network),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						resource.TestCheckResourceAttr("pritunl_route.test", "network", network),
						resource.TestCheckResourceAttrPair("pritunl_route.test", "server_id", "pritunl_server.test", "id"),
					),
				},
				{
					Config: testPritunlRouteSimpleConfig(serverName, update),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("pritunl_server.test", "name", serverName),
						resource.TestCheckResourceAttr("pritunl_route.test", "network", update),
						resource.TestCheckResourceAttrPair("pritunl_route.test", "server_id", "pritunl_server.test", "id"),
					),
				},
			},
		})
	})
}

func testPritunlServerWithoutRouteConfig(serverName string) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name	= "%[1]s"
		}
	`, serverName)
}

func testPritunlRouteSimpleConfig(serverName string, network string) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name	= "%[1]s"
		}
		resource "pritunl_route" "test" {
			server_id    = pritunl_server.test.id
			network		 = "%[2]s"
		}
	`, serverName, network)
}

func testPritunlRouteWithNat(serverName string, network string, nat bool) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name	= "%[1]s"
		}
		resource "pritunl_route" "test" {
			server_id    = pritunl_server.test.id
			network		 = "%[2]s"
			nat			 = %[3]t
		}
	`, serverName, network, nat)
}

func testPritunlRouteWithComment(serverName string, network string, comment string) string {
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name	= "%[1]s"
		}
		resource "pritunl_route" "test" {
			server_id    = pritunl_server.test.id
			network		 = "%[2]s"
			comment		 = "%[3]s"
		}
	`, serverName, network, comment)
}

func testPritunlMultipleRoute(serverName string, networks map[string]interface{}) string {
	routeResource := ""
	for k, v := range networks {
		routeResource += fmt.Sprintf(`
		resource "pritunl_route" "%[1]s" {
			server_id    = pritunl_server.test.id
			network		 = "%[2]s"
		}`, k, v)
	}
	
	return fmt.Sprintf(`
		resource "pritunl_server" "test" {
			name	= "%[1]s"
		}

		%[2]s`, serverName, routeResource)
}

func testPritunlRouteDestroy(s *terraform.State) error {
	serverId := s.RootModule().Resources["pritunl_server.test"].Primary.Attributes["id"]
	fmt.Println(serverId)
	routeId := s.RootModule().Resources["pritunl_route.test"].Primary.Attributes["id"]

	routes, err := testClient.GetRoutesByServer(serverId)
	if err != nil {
		return err
	}
	for _, route := range routes {
		if route.ID == routeId {
			return fmt.Errorf("a route is not destroyed")
		}
	}
	return nil
}

func testResourceDestroy(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if ok && rs.Primary.ID != "" {
			return fmt.Errorf("Resource %s still exists", resource)
		}
		return nil
	}
}