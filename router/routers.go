package router

import (
	"log"
	"net/http"
	"os"
	"xyzeshop/auth"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*gin.Context)
}

type routes struct {
	router *gin.Engine //embedded  struct
}

type Routes []Route

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Contol-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Contol-Allow-Credentials", "true")
		// c.Writer.Header().Set("Access-Contol-Allow-Methods", "GET, POST, PUT, DELETE")
		// c.Writer.Header().Set("Access-Contol-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			// c.AbortWithStatus(http.StatusNoContent)
			// return
		}
		//c.Next()
	}
}

// group xyz health routes
func (r routes) XyzHealthCheck(g *gin.RouterGroup) {
	orderRouteGrouping := g.Group("/xyz")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range healthCheckRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"result": "Specify a valid http method with route"})
			})
		}
	}
}

// grouping user
func (r routes) XyzUser(g *gin.RouterGroup) {
	orderRouteGrouping := g.Group("/xyz")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range healthCheckRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"result": "Specify a valid http method with route"})
			})
		}
	}
}

// grouping product
func (r routes) XyzProduct(g *gin.RouterGroup) {
	orderRouteGrouping := g.Group("/xyz")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range healthCheckRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"result": "Specify a valid http method with route"})
			})
		}
	}
}

// grouping product global routes
func (r routes) XyzGlobalProductRoutes(g *gin.RouterGroup) {
	orderRouteGrouping := g.Group("/xyz-product")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range healthCheckRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"result": "Specify a valid http method with route"})
			})
		}
	}
}

// grouping user
func (r routes) XyzAuthUser(g *gin.RouterGroup) {
	orderRouteGrouping := g.Group("/xyz")
	orderRouteGrouping.Use(CORSMiddleware())
	for _, route := range healthCheckRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"result": "Specify a valid http method with route"})
			})
		}
	}
}

// append routes with version
func GustRoutes() {
	r := routes{
		router: gin.Default(),
	}
	v1 := r.router.Group(os.Getenv("API_VERSION"))
	r.XyzHealthCheck(v1)
	r.XyzUser(v1)
	r.XyzGlobalProductRoutes(v1)

	//for auth
	v1.Use(auth.Auth())
	r.XyzProduct(v1)
	r.XyzAuthUser(v1)

	err := r.router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Printf("server failed to run :%v\n", err)
	}
}
