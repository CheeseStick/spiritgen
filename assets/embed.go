package assets

import _ "embed"

//go:embed fonts/NotoSerifKR-Black.ttf
var FontNotoSerifKR []byte

//go:embed images/spirit_pad_design_1.png
var DesignOne []byte

//go:embed images/spirit_pad_design_2.png
var DesignTwo []byte

//go:embed spirit_pad_template.xlsx
var SpiritPadTemplate []byte
