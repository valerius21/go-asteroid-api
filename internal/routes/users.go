package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pastoapp/astroid-api/internal/middleware/user"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"log"
)

func init() {
	log.SetPrefix("[routes/users] ")
}

type Users struct {
	DB     *orbitdb.Database
	RGroup *gin.RouterGroup
}

type publicKey struct {
	PublicKey string `json:"publicKey"`
}

var users Users

func InitUsers(router *gin.Engine, db *orbitdb.Database) *Users {
	group := router.Group("/users")
	users = Users{
		DB:     db,
		RGroup: group,
	}
	group.POST("/", users.Create)
	group.GET("/:id", users.Find)

	return &users
}

func (u Users) Find(context *gin.Context) {

	id := context.Param("id")

	if id == "" {
		context.JSON(400, gin.H{
			"error": "id is required",
		})
		return
	}

	find, err := user.Find(id)

	if err != nil {
		context.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(200, u.response(find))
}

func (u Users) Create(context *gin.Context) {
	// get Form data
	var pk publicKey

	err := context.ShouldBindJSON(&pk)
	if err != nil {
		context.JSON(400, gin.H{
			"error":   err.Error(),
			"message": "Invalid JSON",
		})
		return
	}

	if pk.PublicKey == "" {
		context.JSON(400, gin.H{
			"error": "publicKey is required",
		})
		return
	}

	// create user
	newUser, err := user.NewUser(pk.PublicKey, false)
	if err != nil {
		context.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(200, u.response(newUser))
}

func (_ Users) response(u *user.User) gin.H {
	return gin.H{
		"_id":       u.ID.String(),
		"publicKey": u.PublicKey,
		"nonce":     u.Nonce,
		"createdAt": u.CreatedAt,
		"updatedAt": u.UpdatedAt,
	}
}

//func (u User) update(context *gin.Context) {
//
//}

//func (u User) requestToken(context *gin.Context) {
//
//}
