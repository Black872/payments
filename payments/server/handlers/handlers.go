package handlers

import "payments/service"

type Handlers struct {
	Auth *AuthorizationHandlers
}

func NewHandlers(s *service.Service) *Handlers {
	return &Handlers{
		Auth: NewAuthorizationHandlers(s),
	}
}
