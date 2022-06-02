package docs

import "embed"

//go:embed index.html swagger.yaml
var Files embed.FS
