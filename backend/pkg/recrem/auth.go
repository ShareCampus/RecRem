package recrem

import (
	"github.com/ShareCampus/RecRem/backend/pkg/utils/httpjson"
	"net/http"
)

func (s *server) hello(w http.ResponseWriter, r *http.Request) {
	httpjson.WriteJson(w, "hello RecRem")
	return
}
