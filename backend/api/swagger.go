package api

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/zorcal/sbgfit/backend/pkg/git"
)

func swaggerUIHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>SBGFit API Documentation</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *, *:before, *:after {
      box-sizing: inherit;
    }
    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: '/swagger/openapi.yml',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        validatorUrl: null,
        onComplete: function() {
          // Lock the URL input to prevent users from changing it
          const urlInput = document.querySelector('.download-url-input');
          if (urlInput) {
            urlInput.readOnly = true;
            urlInput.style.backgroundColor = '#f5f5f5';
            urlInput.title = 'URL is locked to prevent errors';
          }
        }
      });
    };
  </script>
</body>
</html>`

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})
}

func openAPISpecHandler(log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		gitRoot, err := git.ProjectRoot(ctx)
		if err != nil {
			log.ErrorContext(ctx, "Failed to get project git root path", "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		specPath := filepath.Clean(filepath.Join(gitRoot, "schemas", "openapi.yml"))

		specContent, err := os.ReadFile(specPath)
		if err != nil {
			log.ErrorContext(ctx, "Failed to read "+specPath, "error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write(specContent)
	})
}
