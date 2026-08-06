package templates

import "embed"

// FS holds scaffold files used by `dug init`.
//
//go:embed edge.yaml docker-compose.yml gitignore gitkeep README.md
var FS embed.FS
