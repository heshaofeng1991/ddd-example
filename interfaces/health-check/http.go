/*
	@Author  johnny
	@Author  heshaofeng1991@gmail.com
	@Version v1.0.0
	@File    http
	@Date    2022/4/21 12:18
	@Desc
*/

package interfaces

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type HTTPServer struct{}

func NewHTTPServer() HTTPServer {
	return HTTPServer{}
}

func (h HTTPServer) Get(w http.ResponseWriter, r *http.Request) {
	render.Respond(
		w, r, map[string]interface{}{
			"status":    "ok",
			"service":   "ddd-johnny",
			"timestamp": time.Now(),
		},
	)
}
