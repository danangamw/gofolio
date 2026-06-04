package gocms

import "embed"

// WebFS holds all files under the web/ directory embedded at compile time.
// This file must live at the project root (next to web/) because Go's
// //go:embed directive cannot traverse into parent directories.
//
//go:embed web
var WebFS embed.FS
