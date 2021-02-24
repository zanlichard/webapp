package setup

import (
	"net/http"
	"webapp/internal/metrics_mux"
)

// NewServerMux ...
func NewServerMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux = metrics_mux.GetElasticMux(mux)
	//mux = metrics_mux.GetPProfMux(mux)
	mux = metrics_mux.GetPrometheusMux(mux)
	return mux
}
