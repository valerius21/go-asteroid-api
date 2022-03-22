package jwt

import (
	"encoding/base64"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/pastoapp/astroid-api/internal/middleware/user"
	"log"
	"time"
)

// TODO: implement jwt
func init() {
	_, err := NewJWTMiddleware()

	if err != nil {
		log.Fatalf("[JWT] %v\n", err)
	}
}

type Login struct {
	ID        string `json:"id" form:"id" binding:"required"`
	Signature string `json:"signature" form:"signature" binding:"required"`
}

var identityKey = "id"

func NewJWTMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "main",
		Key:        []byte("secret key"), // TODO: change to env variable
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login Login

			if err := c.ShouldBindJSON(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			uid := login.ID
			sign, err := base64.StdEncoding.DecodeString(login.Signature)

			if uid == "" {
				return "", jwt.ErrFailedAuthentication
			}
			if err != nil {
				log.Fatalln(err)
				return "", jwt.ErrFailedAuthentication
			}

			log.Printf("[JWT] Authenticating user: %s\n", uid)

			usr, err := user.Find(uid)
			if err != nil {
				return nil, err
			}

			err = usr.VerifyUser(string(sign))

			if err != nil {
				err2 := usr.RefreshNonce()
				if err2 != nil {
					log.Fatalln(err2)
				}
				return nil, err
			}

			err = usr.RefreshNonce()
			if err != nil {
				return nil, err
			}

			return gin.H{
				"id": uid,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			panic("authorization not implemented")
			return false
		},
		IdentityKey: identityKey,
		IdentityHandler: func(context *gin.Context) interface{} {
			panic("identity handler not implemented")
		},
		// the JWT middleware will call this function if an authorization succeeds. It's in the JWT payload
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			return jwt.MapClaims{"hI": "MOm"}
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TimeFunc:   time.Now,
		CookieName: "asteroid",
	})
}
