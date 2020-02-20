package gin_middleware_cas

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/shenshouer/cas"
	"net/http"
	"net/url"
)

const (
	CAS_URL = "http://1.1.1.1"
)

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, _ := url.Parse(CAS_URL)
		client := cas.NewClient(&cas.Options{URL: u})
		handler := client.HandleFunc(func(writer http.ResponseWriter, request *http.Request) {
			// DO NOTHING
		})
		handler.ServeHTTP(c.Writer, c.Request)
		if !cas.IsAuthenticated(c.Request) {
			client.RedirectToLogin(c.Writer, c.Request)
			// Aborting ther running of  other function
			c.Abort()
		} else {
			session := sessions.Default(c)
			u := session.Get("username")
			if u == nil {
				gob.Register(cas.UserAttributes{})
				// information from cas
				attributes := cas.Attributes(c.Request)
				username := cas.Username(c.Request)
				// set sessions
				if attributes != nil {
					title := attributes["title"][0]
					company := attributes["company"][0]
					department := attributes["department"][0]
					displayName := attributes["displayName"][0]
					session.Set("username", username)
					session.Set("title", title)
					session.Set("company", company)
					session.Set("department", department)
					session.Set("displayName", displayName)
					err := session.Save()
					if err != nil {
						fmt.Printf("error: %v", err)
					}
				}
			}
		}
		c.Next()
	}

}
