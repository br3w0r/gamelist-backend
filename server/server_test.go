package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/br3w0r/gamelist-backend/entity"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

type requestOptions struct {
	Token                 string
	Body                  string
	EqualStatus           int
	ResponseStructPointer interface{} // Must be a pointer
}

func genericRequest(server *gin.Engine, method string, url string, options requestOptions) func() {
	return func() {
		buf := bytes.NewBuffer([]byte(options.Body))
		req, err := http.NewRequest(method, url, buf)
		So(err, ShouldBeNil)

		req.Header.Add("Content-Type", "application/json")
		if options.Token != "" {
			req.Header.Add("Authorization", "Bearer "+options.Token)
		}

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		So(w.Code, ShouldEqual, options.EqualStatus)

		if options.ResponseStructPointer != nil {
			switch method {
			case "POST", "PUT", "PATCH":
				err = json.Unmarshal(w.Body.Bytes(), options.ResponseStructPointer)
				So(err, ShouldBeNil)
			}
		}
	}
}

func TestServer(t *testing.T) {
	t.Log("TestServer requires a working scrapper instance to run")

	if os.Getenv("CI") != "" {
		t.Skip("TestServer doesn't support CI environment")
	}

	// Initializing server with test db
	options := ServerOptions{
		Production:         true,
		ServeStatic:        false,
		ForceScrape:        true,
		ScraperAsync:       false,
		DatabaseDist:       "./test.db",
		ScraperGRPCAddress: "localhost",
		StressTest:         false,
		SilentMode:         true,
	}

	server := NewServer(options)

	defer os.Remove("./test.db")

	Convey("Create profile should return ok status", t,
		genericRequest(server, "POST", "http://localhost/api/v0/profiles",
			requestOptions{
				Body: `
					{
						"nickname": "test",
						"email": "test@mail.com",
						"password": "123456"
					}
				`,
				EqualStatus: http.StatusOK,
			}))

	var tokenPair entity.TokenPair
	Convey("Aquire tokens should give refresh and authentication token", t,
		genericRequest(server, "POST", "http://localhost/api/v0/aquire-tokens", requestOptions{
			Body: `
				{
					"nickname": "test",
					"password": "123456"
				}
			`,
			EqualStatus:           http.StatusOK,
			ResponseStructPointer: &tokenPair,
		}))

	Convey("Refresh token pair should return new refresh and authorization token pair", t,
		genericRequest(server, "POST", "http://localhost/api/v0/refresh-tokens", requestOptions{
			Body: fmt.Sprintf(`
				{
					"refresh_token": "%s"
				}
			`, tokenPair.RefreshToken),
			EqualStatus:           http.StatusOK,
			ResponseStructPointer: &tokenPair,
		}))

	Convey("Authorization should fail for given requests", t, func() {
		Convey("Without authorization header", genericRequest(server, "POST", "http://localhost/api/v0/games/all",
			requestOptions{
				EqualStatus: http.StatusUnauthorized,
			},
		))
		Convey("With empty authorization header", func() {
			buf := bytes.NewBuffer([]byte(""))
			req, err := http.NewRequest("POST", "http://localhost/api/v0/games/all", buf)
			So(err, ShouldBeNil)

			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "")

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})
		Convey("With wrong authorization prefix", func() {
			buf := bytes.NewBuffer([]byte(""))
			req, err := http.NewRequest("POST", "http://localhost/api/v0/games/all", buf)
			So(err, ShouldBeNil)

			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "JWT "+tokenPair.Token)

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusUnauthorized)
		})
		Convey("With wrong authorization token", genericRequest(server, "POST", "http://localhost/api/v0/games/all",
			requestOptions{
				EqualStatus: http.StatusUnauthorized,
				Token:       "abcd",
			},
		))
	})
}
