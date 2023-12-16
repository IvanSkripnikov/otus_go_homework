package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"app/controllers"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

var (
	routes = []route{
		newRoute("GET", "/", controllers.HelloPageHandler),
		newRoute("GET", "/banners", controllers.BannersHandler),
		newRoute("GET", "/banners/([0-9]+)", controllers.BannerHandler),
		newRoute("GET", "/add_banner_to_slot/([\\S]+)", controllers.AddBannerHandler),
		newRoute("GET", "/remove_banner_from_slot/([\\S]+)", controllers.RemoveBannerFromSlot),
		newRoute("GET", "/get_banner_for_show/([\\S]+)", controllers.GetBannerForShow),
		newRoute("GET", "/event_click/([\\S]+)", controllers.EventClick),

		newRoute("GET", "/tasks", controllers.GetAllHandler),
		newRoute("POST", "/tasks", controllers.CreateHandler),
		newRoute("GET", "/tasks/([0-9]+)", controllers.GetHandler),
		newRoute("PUT", "/tasks/([0-9]+)", controllers.UpdateHandler),
		newRoute("DELETE", "/tasks/([0-9]+)", controllers.DeleteHandler),
	}
)

func initHTTPServer() error {
	http.HandleFunc("/", Serve)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		errMessage := fmt.Sprintf("Can't init HTTP server: %v", err)
		fmt.Println(errMessage)
	}
	return nil
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	found := false
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			found = true
			route.handler(w, r)
		}
	}
	if !found && len(allow) == 0 {
		w.WriteHeader(http.StatusNotFound)
		http.NotFound(w, r)
		return
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func GetHttpHandler() *http.ServeMux {
	httpHandler := http.NewServeMux()

	httpHandler.HandleFunc("/banners/", controllers.BannerHandler)
	httpHandler.HandleFunc("/banners", controllers.BannersHandler)

	return httpHandler
}

func main() {
	if err := initHTTPServer(); err != nil {
		fatalMessage := fmt.Sprintf("Can't init http server, err: %v", err)
		fmt.Println(fatalMessage)
	}
	return
}
