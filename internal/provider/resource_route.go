package provider

import (
	"fmt"
	"context"
	"sync"
	"strings"

	"github.com/maulid7/terraform-provider-pritunl/internal/pritunl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// var (
// 	routesCache []pritunl.Route
// 	once        sync.Once
// )

var resourceMutex sync.RWMutex

// func initRoutesCache(meta interface{}, serverId string) error {
// 	apiClient := meta.(pritunl.Client)
// 	var err error
// 	once.Do(func() {
// 		fmt.Println("CREATE CACHE")
// 		routesCache, err = apiClient.GetRoutesByServer(serverId)
// 		if err != nil {
// 			return
// 		}
// 	})
// 	return err
// }

// func clearRouteCache() {
// 	once = sync.Once{}
// }

// func getRouteFromCache(id string) pritunl.Route {
// 	var matchedRoute pritunl.Route
// 	for _, route := range routesCache {
// 		if route.ID == id {
// 			matchedRoute = route
// 			break
// 		}
// 	}
// 	return matchedRoute
// }

func getRouteFromList(id string, list []pritunl.Route) pritunl.Route {
	var matchedRoute pritunl.Route
	for _, route := range list {
		if route.ID == id {
			matchedRoute = route
			break
		}
	}
	return matchedRoute
}

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		Description: "The route resource allows managing information about a particular route in a Pritunl server.",
		Schema: map[string]*schema.Schema{
			"network": {
				Type: 		  schema.TypeString,
				Required:	  true,
				ForceNew: 	  true,
				Description:  "Network address CIDR to route",
				ValidateFunc: func(i interface{}, s string) ([]string, []error) {
					return validation.IsCIDR(i, s)
				},
			},
			"comment": {
				Type:		 schema.TypeString,
				Required: 	 false,
				Optional:	 true,
				Description: "Comment for the route",
			},
			"nat": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Default:	 true,
				Description: "NAT vpn traffic destined to this network",
			},
			"net_gateway": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Description: "Net Gateway vpn traffic destined to this network",
				Computed:    true,
			},
			"server_id": {
				Type:		 schema.TypeString,
				Required: 	 true,
				ForceNew: 	 true,
				Description: "Server ID to attach this route to",
			},
		},
		CreateContext: resourceCreateRoute,
		ReadContext: resourceReadRoute,
		UpdateContext: resourceUpdateRoute,
		DeleteContext: resourceDeleteRoute,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRouteImport,
		},
	}
}

func resourceCreateRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// fmt.Println("CREATE CALLED")
	apiClient := meta.(pritunl.Client)

	serverId := d.Get("server_id").(string)

	routeData := map[string]interface{}{
		"network": 		d.Get("network"),
		"comment": 		d.Get("comment"),
		"nat":			d.Get("nat"),
		"net_gateway":	d.Get("net_gateway"),
	}

	routePayload := pritunl.ConvertMapToRoute(routeData)
	//fmt.Printf("create: %+v\n", routePayload)
	
	// Get latest server status
	server, err := apiClient.GetServer(d.Get("server_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// Start server if it was ONLINE before and status wasn't changed OR status was changed to ONLINE
	shouldServerBeStarted := server.Status == pritunl.ServerStatusOnline

	// Stop server before applying route change
	err = apiClient.StopServer(d.Get("server_id").(string))
	if err != nil {
		return diag.Errorf("Error on stopping server: %s", err)
	}

	// fmt.Printf("%+v", routePayload)

	route, err := apiClient.AddRouteToServer(serverId, routePayload)
	if err != nil {
		return diag.FromErr(err)
	}

	if shouldServerBeStarted {
		err = apiClient.StartServer(d.Get("server_id").(string))
		if err != nil {
			return diag.Errorf("Error on starting server: %s", err)
		}
	}

	d.SetId(route.ID)
	d.Set("network", route.Network)
	d.Set("comment", route.Comment)
	d.Set("nat", route.Nat)
	d.Set("net_gateway", route.NetGateway)

	return nil
}

func resourceReadRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceMutex.RLock()
	defer resourceMutex.RUnlock()

	// fmt.Println("READ CALLED")
	apiClient := meta.(pritunl.Client)

	// if err := initRoutesCache(apiClient, d.Get("server_id").(string)); err != nil {
	// 	return diag.FromErr(err)
	// }

	// route := getRouteFromCache(d.Id())

	routes, err := apiClient.GetRoutesByServer(d.Get("server_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	route := getRouteFromList(d.Id(), routes)

	d.Set("network", route.Network)
	d.Set("comment", route.Comment)
	d.Set("nat", route.Nat)
	d.Set("net_gateway", route.NetGateway)

	return nil
}

func resourceUpdateRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// fmt.Println("UPDATE CALLED")
	apiClient := meta.(pritunl.Client)
	
	server, err := apiClient.GetServer(d.Get("server_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// if err := initRoutesCache(apiClient, d.Get("server_id").(string)); err != nil {
	// 	return diag.FromErr(err)
	// }

	// route := getRouteFromCache(d.Id())

	routes, err := apiClient.GetRoutesByServer(d.Get("server_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	route := getRouteFromList(d.Id(), routes)

	if v, ok := d.GetOk("network"); ok {
		route.Network = v.(string)
	}

	if d.HasChange("comment") {
		route.Comment = d.Get("comment").(string)
	}

	if d.HasChange("nat") {
		route.Nat = d.Get("nat").(bool)
	}

	if d.HasChange("net_gateway") {
		route.NetGateway = d.Get("net_gateway").(bool)
	}

	// Stop server before applying route change
	err = apiClient.StopServer(d.Get("server_id").(string))
	if err != nil {
		return diag.Errorf("Error on stopping server: %s", err)
	}
	// Start server if it was ONLINE before and status wasn't changed OR status was changed to ONLINE
	shouldServerBeStarted := server.Status == pritunl.ServerStatusOnline

	err = apiClient.UpdateRouteOnServer(d.Get("server_id").(string), route)
	if err != nil {
		// start server in case of error?
		return diag.FromErr(err)
	}

	if shouldServerBeStarted {
		err = apiClient.StartServer(d.Get("server_id").(string))
		if err != nil {
			return diag.Errorf("Error on starting server: %s", err)
		}
	}

	return nil
}

func resourceDeleteRoute(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resourceMutex.Lock()
	defer resourceMutex.Unlock()

	// fmt.Println("DELETE CALLED")
	apiClient := meta.(pritunl.Client)

	server, err := apiClient.GetServer(d.Get("server_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// if err := initRoutesCache(apiClient, d.Get("server_id").(string)); err != nil {
	// 	return diag.FromErr(err)
	// }

	// route := getRouteFromCache(d.Id())

	routes, err := apiClient.GetRoutesByServer(d.Get("server_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	route := getRouteFromList(d.Id(), routes)

	// Stop server before applying route change
	err = apiClient.StopServer(d.Get("server_id").(string))
	if err != nil {
		return diag.Errorf("Error on stopping server: %s", err)
	}
	// Start server if it was ONLINE before and status wasn't changed OR status was changed to ONLINE
	shouldServerBeStarted := server.Status == pritunl.ServerStatusOnline

	err = apiClient.DeleteRouteFromServer(d.Get("server_id").(string), route)
	if err != nil {
		return diag.FromErr(err)
	}

	if shouldServerBeStarted {
		err = apiClient.StartServer(d.Get("server_id").(string))
		if err != nil {
			return diag.Errorf("Error on starting server: %s", err)
		}
	}

	d.SetId("")

	return nil
}

func resourceRouteImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	apiClient := meta.(pritunl.Client)

	attributes := strings.Split(d.Id(), "-")
	if len(attributes) < 2 {
		return nil, fmt.Errorf("invalid format: expected ${serverId}-${routeId}, e.g. 60cd0be07723cf3c9114686c-60cd0be17723cf3c91146873, actual id is %s", d.Id())
	}

	serverId := attributes[0]
	routeId := attributes[1]

	d.SetId(routeId)
	d.Set("server_id", serverId)

	routes, err := apiClient.GetRoutesByServer(d.Get("server_id").(string))
	if err != nil {
		return nil, err
	}

	_ = getRouteFromList(d.Id(), routes)

	return []*schema.ResourceData{d}, nil
}
