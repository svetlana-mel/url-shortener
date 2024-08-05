package response

import "net/http"

type ResponseBase struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func RenderError(w http.ResponseWriter, r *http.Request, err error) {

}
