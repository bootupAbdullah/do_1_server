package internal

import (
      "encoding/json"
      "net/http"
      "net/http/httptest"
      "testing"	
)


func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()


	HealthHandler(w, req)

	if w.Code != http.StatusOK {
	  t.Errorf("expected status 200, got %d", w.Code)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)

	if body["status"] != "ok" {
	  t.Errorf("expected status ok, got %s", body["status"])
	}

}

