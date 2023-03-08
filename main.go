package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Ifconfigme struct {
	IPAddr     string `json:"ip_addr"`
	RemoteHost string `json:"remote_host"`
	UserAgent  string `json:"user_agent"`
	Port       int    `json:"port"`
	Language   string `json:"language"`
	Method     string `json:"method"`
	Encoding   string `json:"encoding"`
	Mime       string `json:"mime"`
	Via        string `json:"via"`
	Forwarded  string `json:"forwarded"`
}

func main() {
	tracer.Start()
	defer tracer.Stop()
	// Create a traced mux router
	mux := httptrace.NewServeMux()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	e.GET("/hello", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/ifconfig", func(c echo.Context) error {
		request, err := http.NewRequest("GET", "http://ifconfig.me/all.json", nil)
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Accept-Language", "en_US")
		client := &http.Client{}
		resp, err := client.Do(request)
		if err == nil {
			e.Logger.Debug("ifconfig: ", resp.Body)
			fmt.Printf("ifconfig: %s, ----, %s", resp.Header, resp.Body)

			defer resp.Body.Close()
			body, err2 := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
			if err2 != nil {
				e.Logger.Fatal(err)
			}
			var v Ifconfigme
			err2 = json.Unmarshal(body, &v)
			if err2 != nil {
				e.Logger.Fatal(err)
			}
			e.Logger.Debug(v.IPAddr)

			return c.JSON(http.StatusOK, struct{ IPAddr string }{IPAddr: string(v.IPAddr)})
		} else {
			e.Logger.Error("ERROR! ", err)
			return c.JSON(http.StatusInternalServerError, struct{ Error string }{Error: "500"})
		}

	})

	e.GET("/gsping1", func(c echo.Context) error {
		request, err := http.NewRequest("GET", "http://gsping2.default:8080/gsping2", nil)
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Accept-Language", "en_US")
		client := &http.Client{}
		resp, err := client.Do(request)
		type Result struct {
			Result string `json:"Result"`
		}
		if err == nil {
			e.Logger.Debug("gsping1: ", resp.Body)
			defer resp.Body.Close()
			body, err2 := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
			if err2 != nil {
				e.Logger.Fatal(err)
			}
			var v Result
			err2 = json.Unmarshal(body, &v)
			if err2 != nil {
				e.Logger.Fatal(err)
			}
			e.Logger.Debug(v.Result)

			return c.JSON(http.StatusOK, struct{ IPAddr string }{IPAddr: string(v.Result)})
		} else {
			e.Logger.Error("ERROR! ", err)
			return c.JSON(http.StatusInternalServerError, struct{ Error string }{Error: "500"})
		}

	})

	e.GET("/gsping2", func(c echo.Context) error {
		request, err := http.NewRequest("GET", "http://gsping3.default:8080/return123", nil)
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Accept-Language", "en_US")
		client := &http.Client{}
		resp, err := client.Do(request)
		type Result struct {
			Result string `json:"Result"`
		}
		if err == nil {
			e.Logger.Debug("gsping2: ", resp.Body)
			defer resp.Body.Close()
			body, err2 := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
			if err2 != nil {
				e.Logger.Fatal(err)
			}
			var v Result
			err2 = json.Unmarshal(body, &v)
			if err2 != nil {
				e.Logger.Fatal(err)
			}
			e.Logger.Debug(v.Result)

			return c.JSON(http.StatusOK, struct{ Result string }{Result: string(v.Result)})
		} else {
			e.Logger.Error("ERROR! ", err)
			return c.JSON(http.StatusInternalServerError, struct{ Error string }{Error: "500"})
		}

	})

	e.GET("/return123", func(c echo.Context) error {
		time.Sleep(time.Duration(1000) * time.Millisecond)
		return c.JSON(http.StatusOK, struct{ Result string }{Result: "123"})

	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
