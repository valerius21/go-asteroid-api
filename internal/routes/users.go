package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pastoapp/astroid-api/internal/middleware/user"
)

type Users struct {
	RGroup *gin.RouterGroup
}

var users Users

func InitUsers(router *gin.Engine) {
	group := router.Group("/users")
	users = Users{
		RGroup: group,
	}
	group.POST("/", users.Create)
	group.GET("/:id", users.Find)
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
		return
	}

	context.JSON(200, u.response(find))
}

//func (u User) findAll(context *gin.Context) {
//	context.JSON(200, gin.H{
//		"user": "dummy",
//	})
//}

func (u Users) Create(context *gin.Context) {
	// get Form data
	publicKey := context.PostForm("publicKey")

	if publicKey == "" {
		context.JSON(400, gin.H{
			"error": "publicKey is required",
		})
		return
	}

	// create user
	newUser, err := user.NewUser(publicKey, false)
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
