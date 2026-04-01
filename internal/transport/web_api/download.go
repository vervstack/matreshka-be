package web_api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	api "go.vervstack.ru/matreshka/internal/api/server/matreshka_api"
)

const downloadBrowserViewHTML = `
<!DOCTYPE html>
<html>
<head>
  <title>%s</title>
  <meta charset="UTF-8">
	<link rel="icon" href="https://s3-api.redsock.ru/verv/matreshka/icon.png" type="image/x-icon">
</head>
<body>
  <pre>%s</pre>
</body>
</html>`

const downloadBrowserParam = "browser"

func (s *Server) GetConfig(c *gin.Context) {
	apiReq := &api.GetConfig_Request{
		ConfigName: c.Param(configNameParam),
		Version:    extractVersion(c),
		Format:     extractFormat(c),
	}

	browserView := c.Param(downloadBrowserParam) == "true" ||
		strings.Contains(
			strings.ToLower(c.Request.Header.Get("User-Agent")),
			"browser")
	apiResp, err := s.grpcApiServer.GetConfig(c, apiReq)
	if err != nil {
		return
	}

	out := apiResp.Config
	if browserView {
		out = []byte(fmt.Sprintf(downloadBrowserViewHTML, apiReq.ConfigName, apiResp.Config))
	}

	contentType := "text/html"
	if !browserView {
		contentType = responseContentTypeByFormat(apiReq.Format)
	}

	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, "text/html", out)
}

func responseContentTypeByFormat(format api.Format) string {
	switch format {
	case api.Format_yaml:
		return "application/yaml"
	default:
		return "text/plain"
	}
}
