package router

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	swaggerUIPageContent = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<meta name="description" content="SwaggerUI"/>
	<title>SwaggerUI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@latest/swagger-ui.css" />
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@latest/index.css" />
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@latest/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@latest/favicon-16x16.png" sizes="16x16" />
</head>
<body>
<div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@latest/swagger-ui-bundle.js" charset="UTF-8" crossorigin> </script>
    <script src="https://unpkg.com/swagger-ui-dist@latest/swagger-ui-standalone-preset.js" charset="UTF-8" crossorigin> </script>
    <script src="https://unpkg.com/swagger-ui-dist@latest/swagger-initializer.js" charset="UTF-8" crossorigin> </script>
<script>
	window.onload = () => {
		window.ui = SwaggerUIBundle({
			url:    '/api.json',
			dom_id: '#swagger-ui',
		});
	};
</script>
</body>
</html>
`
)

func init() {
	s := g.Server()
	//跨域处理
	s.Use(ghttp.MiddlewareCORS)
	swaggerUI(s)
}

func swaggerUI(s *ghttp.Server) {
	//API SwaggerUI
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/swagger-ui", func(r *ghttp.Request) {
			r.Response.Write(swaggerUIPageContent)
		})
	})
}
