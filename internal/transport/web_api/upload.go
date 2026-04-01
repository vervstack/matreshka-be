package web_api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
)

func (s *Server) UploadConfig(c *gin.Context) {
	grpcReq := &api.SaveConfig_Request{
		ConfigName: c.Param(configNameParam),
		Format:     extractFormat(c),
		Version:    extractVersion(c),
	}

	var err error
	grpcReq.Config, err = c.GetRawData()
	if err != nil {
		c.String(http.StatusBadRequest, "Error reading request body")
		return
	}

	_, err = s.grpcApiServer.SaveConfig(c, grpcReq)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}

func extractVersion(c *gin.Context) *string {
	version, isPresented := c.GetQuery(versionParam)
	if isPresented {
		return &version
	}

	return nil
}

func extractFormat(ctx *gin.Context) api.Format {
	switch ctx.Query(formatParam) {
	case api.Format_env.String():
		return api.Format_env
	default:
		return api.Format_yaml
	}
}
