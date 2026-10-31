package api

import "github.com/router-for-me/CLIProxyAPI/v8/internal/api/handlers/management"

// ManagementHandler shares the configured OAuth persistence hooks with the pool.
func (s *Server) ManagementHandler() *management.Handler {
	return s.mgmt
}
