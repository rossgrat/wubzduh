package web

import "net/http"

func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		s.logger.Info("request",
			"remote_addr", r.RemoteAddr,
			"method", r.Method,
			"uri", r.RequestURI,
		)
	})
}
