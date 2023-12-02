package main

import (
	"app/controllers"
	"fmt"
	"net/http"
	"regexp"
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

var routes = []route{
	newRoute("GET", "/", controllers.HelloPage),
	newRoute("GET", "/tasks", controllers.GetAllHandler),
	newRoute("GET", "/banners", controllers.GetAllBanners),
	newRoute("POST", "/tasks", controllers.CreateHandler),
	newRoute("GET", "/tasks/([0-9]+)", controllers.GetHandler),
	newRoute("PUT", "/tasks/([0-9]+)", controllers.UpdateHandler),
	newRoute("DELETE", "/tasks/([0-9]+)", controllers.DeleteHandler),
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			//ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			//route.handler(w, r.WithContext(ctx))
			route.handler(w, r)

		}
	}
	/*if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)*/
}

func main() {
	if err := initHTTPServer(); err != nil {
		fatalMessage := fmt.Sprintf("Can't init http server, err: %v", err)
		fmt.Println(fatalMessage)
	}
	return
}
