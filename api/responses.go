package api

import (
	"encoding/json"
	"net/http"
)

func createResponse(statusCode int, w http.ResponseWriter, contentType string, response interface{}) (err error) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", contentType)
	var jsonResp []byte
	if contentType == "application/json" {
		jsonResp, err = json.Marshal(response)
	} else {
		jsonResp = []byte(response.(string))
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write(jsonResp)
	}
	return err
}

func Ok(w http.ResponseWriter, data interface{}, meta interface{}) {
	tmp := map[string]interface{}{
		"meta": meta,
		"data": data,
	}
	_ = createResponse(http.StatusOK, w, "application/json", tmp)
}

func Err401Unauthorized(w http.ResponseWriter, err string) {
	_ = createResponse(http.StatusUnauthorized, w, "application/text", err)
}

func Err400BR(w http.ResponseWriter, err string) {
	_ = createResponse(http.StatusBadRequest, w, "application/text", err)
}

func Err500ISE(w http.ResponseWriter, err string) {
	_ = createResponse(http.StatusInternalServerError, w, "application/text", err)
}
