package middlewares

import (
	//"fmt"
	"fmt"
	"net/http"
	"strings"
)

type HPPoptions struct {
	CheckQuery              bool
	CheckBody               bool
	CheckBodyOnlyForContent string
	WhiteList               []string
}

func Hpp(options HPPoptions) func(http.Handler) http.Handler {
	fmt.Println("HPP middlware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("HPP middlware being returned...")
			// For BodyParams
			if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyOnlyForContent) {
				filterBodyParams(r, options.WhiteList)
			}
			// For QueryParam
			if options.CheckQuery && r.URL.Query() != nil {
				filterQeuryParams(r, options.WhiteList)
			}

			next.ServeHTTP(w, r)
			fmt.Println("HPP middlware ends")
		})
	}
}

func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func filterQeuryParams(r *http.Request, whiteList []string) {

	query := r.URL.Query()

	for k, v := range query {
		if len(v) > 1 {
			//query.Set(k, v[0])
			query.Set(k, v[len(v)-1])
		    
		}
		if !iswhiteListed(k, whiteList) {
			query.Del(k)
		}
		
	}
	r.URL.RawQuery = query.Encode()

}

func filterBodyParams(r *http.Request, whiteList []string) {

	//query := r.URL.Query()

	err := r.ParseForm()
	if err != nil {
		return
	}

	for k, v := range r.Form {
		if len(v) > 1 {
			//r.Form.Set(k, v[0])
			r.Form.Set(k, v[len(v)-1])
		}
		
		if !iswhiteListed(k, whiteList) {
			delete(r.Form, k)
		}
	}

}

func iswhiteListed(param string, whiteList []string) bool {

	for _, v := range whiteList {
		if v == param {
		  return true
		}
	}
	return false
}
