package handlers_utils

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"stories/database"
	"strings"

	"stories/helpers"
)

type Route struct {
	Parameters map[string][]string
	Handler    http.Handler
	Methods    []string
}

type Router struct {
	rootNode *RouteNode
}

var notFoundError = helpers.NewHttpError(http.StatusNotFound, "not found")

func NewRouter() Router {
	var router = Router{rootNode: &RouteNode{}}
	// Setup the main thing
	http.Handle("/", router)
	return router
}

func (router Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	/*var node *RouteNode
	var last *RouteNode*/
	var headers http.Header
	var err error
	var route Route
	// Parse the route
	headers = writer.Header()
	// Add CORS headers (right now because they must be there always)
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	headers.Add("Access-Control-Allow-Credentials", "true")
	// Parse the route now
	route, err = router.parseUrl(request.URL, request.Method)
	if err != nil {
		if methodNotAllowedError, ok := err.(helpers.MethodNotAllowedError); ok {
			headers.Add("Access-Control-Allow-Methods", strings.Join(methodNotAllowedError.Methods, ", "))
		}
		RespondWithError(writer, err)
		return
	}
	headers.Add("Access-Control-Allow-Methods", strings.Join(route.Methods, ", "))
	if request.Method == http.MethodOptions {
		writer.WriteHeader(http.StatusNoContent)
	} else {
		var handler http.Handler
		if len(route.Parameters) > 0 {
			request.Form = route.Parameters
		}
		// Get the resolved handler
		handler = route.Handler
		// Call the serve method
		handler.ServeHTTP(writer, request)
	}
}

func (router Router) Add(root string, endpoint string, handler NoAuthorizedHTTPHandler, method string, q database.IQuery) {
	var node *RouteNode
	var fragments []string
	var err error
	fragments = strings.Split(path.Join(strings.Trim(root, "/"), endpoint), "/")
	node = router.rootNode
	for _, fragment := range fragments {
		node, err = node.AddFragment(fragment)
		if err != nil {
			panic(err)
		}
	}
	// Also add the method as a matcher such that
	// we can have multiple methods with the same
	// path
	node, err = node.AddFragment(method)
	if err != nil {
		panic(err)
	}
	// Set handler
	node.handler = handler

}

func (router Router) AddNoAuthenticated(root string, endpoint string, handler NoAuthorizedHTTPHandler, method string, q database.IQuery) {
	router.Add(root, endpoint, NewRequestNoAuthorizer(handler), method, q)
}

func (router Router) AddAuthenticated(root string, endpoint string, handler AuthorizedHTTPHandler, method string, q database.IQuery) {
	router.Add(root, endpoint, NewRequestAuthorizer(handler), method, q)
}

func (router Router) parseUrl(currentUrl *url.URL, method string) (Route, error) {
	var fragments []string
	var node *RouteNode
	var methodNode *RouteNode
	var route Route
	fragments = strings.Split(strings.Trim(currentUrl.Path, "/"), "/")
	node = router.rootNode
	route.Parameters = make(map[string][]string)
	// Find the handler because we need it anyway
	//
	// Also, collect parameters as we do so
	for _, fragment := range fragments {
		node = node.FindMatchingChild(fragment)
		if node == nil {
			return Route{}, notFoundError
		} else {
			if node.pattern != "" {
				// FIXME: perhaps we could use the actual type
				//        by specifying it in the url pattern
				route.Parameters[node.pattern] = []string{fragment}
			}
		}
	}
	route.Methods = node.GetMethods()
	if method != http.MethodOptions {
		methodNode = node.FindMatchingChild(method)
		if methodNode == nil {
			return Route{}, helpers.NewMethodNotAllowedError(node.GetMethods())
		}
		if methodNode.handler == nil {
			return Route{}, fmt.Errorf("handler is not defined for this endpoint")
		}
		route.Handler = methodNode.handler
	}
	return route, nil
}
