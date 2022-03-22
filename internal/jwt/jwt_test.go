package jwt

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	user2 "github.com/pastoapp/astroid-api/internal/middleware/user"
	"github.com/pastoapp/astroid-api/internal/orbitdb"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hi": "mom",
		})
	})

	return r
}

func TestNewJWTMiddleware(t *testing.T) {
	privateK, _ := rsa.GenerateKey(rand.Reader, 4096)
	pubkBytes := x509.MarshalPKCS1PublicKey(&privateK.PublicKey)
	pubkPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubkBytes,
	})

	PublicKey := string(pubkPEM)

	cancelFunc, err := orbitdb.InitializeOrbitDB("http://localhost:5001", t.TempDir())
	if err != nil {
		t.Fatalf("Error initializing OrbitDB: %v", err)
	}
	defer cancelFunc()
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := setupRouter()

	t.Run("should test the gin instance", func(t *testing.T) {
		w := performRequest(r, "GET", "/")
		if w.Code != http.StatusOK {
			t.Errorf("Response code is %v", w.Code)
		}
	})

	t.Run("should return a valid JWT", func(t *testing.T) {
		user, err := user2.NewUser(PublicKey, false)

		jwtMiddleware, err := NewJWTMiddleware()

		if err != nil {
			t.Errorf("Error creating JWT middleware: %s", err)
		}
		if jwtMiddleware == nil {
			t.Error("Error creating JWT middleware: nil")
		}

		r.POST("/login", jwtMiddleware.LoginHandler)

		r.NoRoute(jwtMiddleware.MiddlewareFunc(), func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			log.Printf("NoRoute claims: %#v\n", claims)
			c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
		})
		auth := r.Group("/auth")
		// Refresh time can be longer than token timeout
		auth.GET("/refresh_token", jwtMiddleware.RefreshHandler)
		auth.Use(jwtMiddleware.MiddlewareFunc())
		{
			auth.GET("/hello", func(c *gin.Context) {
				claims := jwt.ExtractClaims(c)
				//user, _ := c.Get(identityKey)
				c.JSON(200, gin.H{
					"userID": claims[identityKey],
					//"userName": user.(*User).UserName,
					"text": "Hello World.",
				})

			})
		}

		w, err := func(r http.Handler, path string) (*httptest.ResponseRecorder, error) {
			// sign nonce

			nonce, err2 := base64.StdEncoding.DecodeString(user.Nonce)

			if err2 != nil {
				return nil, err2
			}

			signedNonce, err3 := privateK.Sign(rand.Reader, nonce, &rsa.PSSOptions{
				SaltLength: rsa.PSSSaltLengthAuto,
				Hash:       crypto.SHA256,
			})

			if err3 != nil {
				return nil, err3
			}

			var jsonData, err = json.Marshal(gin.H{
				"id":        user.ID.String(),
				"signature": base64.StdEncoding.EncodeToString(signedNonce),
			})

			if err != nil {
				return nil, err
			}

			req, err := http.NewRequest("POST", path, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			return w, err
		}(r, "/login")

		if err != nil {
			t.Errorf("Error making request: %s", err)
		}
		resBody := w.Body.String()
		t.Log(resBody)

		if len(resBody) == 0 {
			t.Errorf("Response body is empty")
		}

		if w.Code != http.StatusOK {
			t.Errorf("Response code is %v", w.Code)
		}

		//w = performRequest(r, "GET", "/auth/refresh_token")
		//if w.Code != http.StatusOK {
		//	t.Errorf("Response code is %v", w.Code)
		//}
		//
		//w = performRequest(r, "GET", "/auth/hello")
		//if w.Code != http.StatusOK {
		//	t.Errorf("Response code is %v", w.Code)
		//}
	})
}
