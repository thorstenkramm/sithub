// Package assets provides embedded frontend assets.
package assets

import "embed"

// Web contains the embedded frontend assets.
//
//go:embed web/* web/*/*
var Web embed.FS
