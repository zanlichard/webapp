package setup

import (
	"net/http"
	metrics_mux2 "webapp/frame/internal/metrics_mux"
)

// NewServerMux ...
func NewServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux = metrics_mux2.GetElasticMux(mux)
	//mux = metrics_mux.GetPProfMux(mux)
	mux = metrics_mux2.GetPrometheusMux(mux)
	return mux
}
