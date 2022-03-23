package main

import (
	jwt2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/pastoapp/astroid-api/internal/jwt"
	"github.com/pastoapp/astroid-api/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
	odb "github.com/pastoapp/astroid-api/internal/orbitdb"
)

var (
	ipfsURL       = "http://localhost:5001"
	orbitDbDir    = "./data/orbitdb"
	defaultStores = []string{"users", "notes"}
)

func main() {
	// main database context

	// create a new orbitdb instance
	cancelODB, err := odb.InitializeOrbitDB(ipfsURL, orbitDbDir)
	defer cancelODB() // cancel the orbitdb context

	// gin server
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/keys", routes.KeyGen)
	r.POST("/sign", routes.Sign)

	r.GET("/users/:id", routes.DefaultUsers.Find)
	r.POST("/users", routes.DefaultUsers.Create)

	authMiddleware, err := jwt.NewJWTMiddleware()

	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt2.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	auth := r.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", func(c *gin.Context) {
			claims := jwt2.ExtractClaims(c)
			//user, _ := c.Get(identityKey)
			c.JSON(200, gin.H{
				"userID": claims["_id"],
				//"userName": user.(*User).UserName,
				"text": "Hello World.",
			})

		})
	}

	err = r.Run(":3000")

	if err != nil {
		log.Panicf("Error starting server: %v", err)
	}
}
