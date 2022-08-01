package main

import (
	"flag"
	jwt2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator/v10"
	"github.com/pastoapp/astroid-api/internal/jwt"
	"github.com/pastoapp/astroid-api/internal/middleware/note"
	"github.com/pastoapp/astroid-api/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	odb "github.com/pastoapp/astroid-api/internal/orbitdb"
)

var (
	// ipfsURL holds the value of the IPFS HTTP API Endpoint of the desired node. Defaults to a local node.
	ipfsURL = "http://localhost:5001"
	// orbitDbDir is the path which holds the local state of the applications OrbitDB.
	orbitDbDir = "./data/orbitdb"
)

// init runs at the initialization of the module; meaning that it runs before main.
func init() {
	flag.StringVar(&ipfsURL, "ipfs-url", "http://127.0.0.1:5001", "IPFS HTTP API Endpoint")
	flag.StringVar(&orbitDbDir, "data-dir", "./data/orbitdb", "Data Storage Folder")

	flag.Parse()

	// validate flags
	v := validator.New()

	err := v.Var(ipfsURL, "url") // validates if the ipfsURL matches an HTTP schema.
	if err != nil {
		panic("no valid IPFS HTTP API Endpoint is set")
		return
	}

	// verify that the orbitDbDir exists
	_, err = os.Stat(orbitDbDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("[cmd/main] OrbitDB directory does not exist. Creating...")
			err := os.MkdirAll(orbitDbDir, 750)
			if err != nil {
				panic("[cmd/main] Failed to create OrbitDB directory")
				return
			}
		} else {
			panic(err)
			return
		}
	}
}

// main serves as the entry point into the application.
func main() {

	// main database context

	// create a new OrbitDB instance
	cancelODB, err := odb.InitializeOrbitDB(ipfsURL, orbitDbDir)
	defer cancelODB() // cancel the OrbitDB context

	// gin is used as the HTTP server. The following initializes it.
	r := gin.Default()

	// In PoC, we set the Cross-Origin policies to allow all.
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	// Testing endpoint, to see if the server is live.
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Attaching cryptographic methods to routes.
	r.GET("/keys", routes.KeyGen)
	r.POST("/sign", routes.Sign)

	// Attaching the user model to routes
	r.GET("/users/:id", routes.DefaultUsers.Find)
	r.POST("/users", routes.DefaultUsers.Create)

	// Attaching the nodes model to routes
	r.GET("/notes/", note.FindAll)
	r.POST("/notes/", note.Create)

	// Setting up the authentication middleware
	authMiddleware, err := jwt.NewJWTMiddleware()

	// Attaching the middleware for logins to the /login route
	r.POST("/login", authMiddleware.LoginHandler)

	// Catching not-defined routes with a proper response
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt2.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	// Create a route-group for everything authentication related.
	auth := r.Group("/auth")

	// Refresh time can be longer than token timeout.
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// Test/Mock function for authenticated users to see if they are truly authenticated.
		auth.GET("/hello", func(c *gin.Context) {
			claims := jwt2.ExtractClaims(c)
			//user, _ := c.Get(identityKey)
			c.JSON(200, gin.H{
				"userID": claims[jwt2.IdentityKey],
				//"userName": user.(*User).UserName,
				"text": "Hello World.",
			})

		})
	}

	// Open the Gin-Web-Server on port 3000
	err = r.Run(":3000")

	// Otherwise...
	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}
}
