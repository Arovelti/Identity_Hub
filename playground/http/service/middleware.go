package service

import (
	"fmt"
	"net/http"
)

func (s *Service) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Add("X-Basic-Auth", `Basic realm="Give username and password"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message": "No basic auth present"}`))
			return
		}

		if err := s.isAuthorised(username, password); err != nil {
			w.Header().Add("X-Basic-Auth", `Basic realm="Give username and password"`)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("{invalid credentials: %v}", err)))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Service) isAuthorised(username, password string) error {
	p, err := s.Repo.GetByUsername(username)
	if err != nil {
		return err
	}

	if p.Password != password {
		return fmt.Errorf("invalid password for user %s", username)
	}

	return nil
}
