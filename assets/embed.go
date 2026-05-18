// Package assets exposes embedded binary assets via //go:embed.
// Variables are grouped by product to indicate scope; shared assets
// have no product prefix.
package assets

import _ "embed"

// ────────── Shared ──────────

//go:embed fonts/NotoSerifKR-Black.ttf
var FontNotoSerifKR []byte

// ────────── Spirit Tablet (영가위패) ──────────

//go:embed images/spirit_pad_design_1.png
var SpiritTabletDesignOne []byte

//go:embed images/spirit_pad_design_2.png
var SpiritTabletDesignTwo []byte

//go:embed spirit_pad_template.xlsx
var SpiritTabletTemplate []byte
