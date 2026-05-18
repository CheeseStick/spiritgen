// Package assets exposes embedded binary assets via //go:embed.
// Variables are grouped by product to indicate scope; shared assets
// have no product prefix.
package assets

import _ "embed"

// ────────── Shared ──────────

//go:embed fonts/NotoSerifKR-Black.ttf
var FontNotoSerifKR []byte

//go:embed fonts/NotoSerifKR-Medium.ttf
var FontNotoSerifKRMedium []byte

//go:embed fonts/NotoSerifKR-Light.ttf
var FontNotoSerifKRLight []byte

// ────────── Spirit Tablet (영가위패) ──────────

//go:embed images/spirit_pad_design_1.png
var SpiritTabletDesignOne []byte

//go:embed images/spirit_pad_design_2.png
var SpiritTabletDesignTwo []byte

//go:embed spirit_pad_template.xlsx
var SpiritTabletTemplate []byte

// ────────── Lantern Tablet (연등/인등) ──────────

//go:embed images/lantern_tablet_design_1.png
var LanternTabletDesignOne []byte

//go:embed images/lantern_tablet_design_2.png
var LanternTabletDesignTwo []byte

//go:embed lantern_tablet_template.xlsx
var LanternTabletTemplate []byte
