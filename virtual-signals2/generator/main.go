package main

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"os"
	"sort"
	"time"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
)

const (
	directory    = "./"
	directoryMod = directory + "../mod/"
)

type (
	glyphType       int
	fineTuneGroupID int
	group           int
)

const (
	glyphTypeMath glyphType = iota
	glyphTypeIcon
	glyphTypeLabel
	glyphTypeTinyText
)

const (
	ftOffsetUpXS = iota
	ftOffsetUpMD
	ftOffsetUpLG
	ftOffsetDownXS
	fnScaleSM
	fnUnScaleSM
	fnUnScaleMD
	fnScaleMD
	ftScaleLG
	ftScaleXXL
	ftOffsetRightXS
)

const (
	groupMath = iota
	groupGreek
	groupIcons
	groupCommon
	groupCommon2
	groupCommon3
)

type signal struct {
	group    group
	name     string
	locName  string
	content  string
	fineTune []fineTuneGroupID
}

type groupMeta struct {
	sortOrder      int
	name           string
	backgroundFile string
	glyphType      glyphType
	fineTune       []fineTuneGroupID
}

type fineTune struct {
	offsetX   float64 // offset in % of glyph height
	offsetY   float64 // offset in % of glyph height
	textScale float64 // 1 = no scale
}

var fineTunes = map[fineTuneGroupID]fineTune{
	ftOffsetUpMD:    {offsetY: -10},
	ftOffsetUpLG:    {offsetY: -20},
	ftOffsetUpXS:    {offsetY: -5},
	ftOffsetDownXS:  {offsetY: 5},
	ftOffsetRightXS: {offsetX: 5},
	fnScaleSM:       {textScale: 1.1},
	fnUnScaleSM:     {textScale: 0.9},
	fnUnScaleMD:     {textScale: 0.8},
	fnScaleMD:       {textScale: 1.35},
	ftScaleLG:       {textScale: 1.5},
	ftScaleXXL:      {textScale: 1.75},
}

var groups = map[group]groupMeta{
	groupCommon: {
		sortOrder:      0,
		name:           "cmn",
		backgroundFile: "sig_red",
		glyphType:      glyphTypeLabel,
		fineTune: []fineTuneGroupID{
			fnScaleSM,
		},
	},
	groupCommon2: {
		sortOrder:      2,
		name:           "cmn2",
		backgroundFile: "sig_green",
		glyphType:      glyphTypeLabel,
		fineTune: []fineTuneGroupID{
			fnScaleSM,
		},
	},
	groupCommon3: {
		sortOrder:      3,
		name:           "cmn3",
		backgroundFile: "sig_lime",
		glyphType:      glyphTypeLabel,
		fineTune: []fineTuneGroupID{
			fnScaleSM,
		},
	},
	groupMath: {
		sortOrder:      10,
		name:           "math",
		backgroundFile: "sig_blue",
		glyphType:      glyphTypeMath,
	},
	groupGreek: {
		sortOrder:      20,
		name:           "greek",
		backgroundFile: "sig_purple",
		glyphType:      glyphTypeMath,
	},
	groupIcons: {
		sortOrder:      30,
		name:           "fa",
		backgroundFile: "sig_light_purple",
		glyphType:      glyphTypeIcon,
	},
}

var groupSignals = map[group][]signal{
	groupMath: {
		// expressions
		{name: "eq", locName: "Equal", content: "=", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "neq", locName: "Not Equal", content: "≠", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "aeq", locName: "Approximate Equal", content: "≈", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "gt", locName: "Greater Then", content: ">", fineTune: []fineTuneGroupID{
			fnScaleMD,
			ftOffsetRightXS,
		}},
		{name: "lt", locName: "Less Then", content: "<", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{name: "ge", locName: "Greater or Equal", content: "≥", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{name: "le", locName: "Less or Equal", content: "≤", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{name: "plus", locName: "Plus", content: "+", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "minus", locName: "Minus", content: "-", fineTune: []fineTuneGroupID{
			ftOffsetUpLG,
			ftScaleXXL,
		}},
		{name: "plus_minus", locName: "Plus-Minus", content: "±", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{name: "percent", locName: "Percent", content: "%", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
		}},
		// geo
		{name: "geo_angle", locName: "Angle", content: "∠", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{name: "geo_perpendicular", locName: "Perpendicular", content: "⊥", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
		}},
		{name: "geo_parallel", locName: "Parallel", content: "∥", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
		}},
		{name: "geo_similar", locName: "Similar", content: "~", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		// algo
		{name: "algo_lemni", locName: "INF", content: "∞", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "algo_func", locName: "Function", content: "f", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "algo_delta", locName: "Delta", content: "∆", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "algo_sigma", locName: "Sigma", content: "∑"},
		{name: "algo_golden", locName: "Golden", content: "φ"},
		{name: "algo_pi", locName: "PI", content: "π"},
		{name: "algo_integral", locName: "Integral", content: "∫"},
	},
	groupGreek: {
		{name: "alpha", locName: "Alpha", content: "α"},
		{name: "beta", locName: "Beta", content: "β"},
		{name: "gamma", locName: "Gamma", content: "γ"},
		{name: "delta", locName: "Delta", content: "δ"},
		{name: "epsilon", locName: "Epsilon", content: "ε"},
		{name: "zeta", locName: "Zeta", content: "ζ"},
		{name: "eta", locName: "Eta", content: "η"},
		{name: "theta", locName: "Theta", content: "θ"},
		{name: "iota", locName: "Iota", content: "ι"},
		{name: "kappa", locName: "Kappa", content: "κ"},
		{name: "lambda", locName: "Lambda", content: "λ"},
		{name: "mu", locName: "Mu", content: "μ"},
		{name: "nu", locName: "Nu", content: "ν"},
		{name: "xi", locName: "Xi", content: "ξ"},
		{name: "omicron", locName: "Omicron", content: "ο"},
		{name: "pi", locName: "Pi", content: "π"},
		{name: "rho", locName: "Rho", content: "ρ"},
		{name: "sigma", locName: "Sigma", content: "σ"},
		{name: "tau", locName: "Tau", content: "τ"},
		{name: "upsilon", locName: "Upsilon", content: "υ"},
		{name: "phi", locName: "Phi", content: "φ"},
		{name: "chi", locName: "Chi", content: "χ"},
		{name: "psi", locName: "Psi", content: "ψ"},
		{name: "omega", locName: "Omega", content: "ω"},
	},
	groupCommon: {
		{name: "on", locName: "ON", content: "ON"},
		{name: "off", locName: "OFF", content: "OFF"},
		{name: "min", locName: "MIN", content: "MIN"},
		{name: "max", locName: "MAX", content: "MAX"},
		{name: "in", locName: "IN", content: "IN"},
		{name: "out", locName: "OUT", content: "OUT"},
	},
	groupCommon2: {
		{name: "size", locName: "Size", content: "SIZE"},
		{name: "buffer", locName: "Buffer", content: "BUF"},
		{name: "limit", locName: "Limit", content: "LIM"},
		{name: "capacity", locName: "Capacity", content: "CAP"},
		{name: "nil", locName: "Nil", content: "NIL"},
		{name: "fuck", locName: "Fuck", content: "FUCK", fineTune: []fineTuneGroupID{fnUnScaleSM}},
	},
	groupCommon3: {
		{name: "size_small", locName: "Small", content: "SM"},
		{name: "size_medium", locName: "Medium", content: "MD"},
		{name: "size_large", locName: "Large", content: "LG"},
		{name: "size_extra_large", locName: "Extra Large", content: "XL"},
		{name: "size_xxl", locName: "XXL", content: "XXL"},
		{name: "sev_info", locName: "Info", content: "INF"},
		{name: "sev_warn", locName: "Warning", content: "WRN", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{name: "sev_err", locName: "Error", content: "ERR"},
		{name: "sev_critical", locName: "Critical", content: "CRT"},
		{name: "alert", locName: "Alert", content: "ALRT", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{name: "trigger", locName: "Trigger", content: "TRIG", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{name: "if", locName: "IF", content: "IF"},
		{name: "else", locName: "ELSE", content: "ELSE", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{name: "action", locName: "ACT", content: "ACT"},
		{name: "var_x", locName: "X", content: "[X]"},
		{name: "var_y", locName: "Y", content: "[Y]"},
		{name: "var_z", locName: "Z", content: "[Z]"},
		{name: "var_w", locName: "W", content: "[W]"},
		{name: "color_r", locName: "R", content: "(R)"},
		{name: "color_g", locName: "G", content: "(G)"},
		{name: "color_b", locName: "B", content: "(B)"},
		{name: "led1", locName: "L1", content: "L1"},
		{name: "led2", locName: "L2", content: "L2"},
		{name: "led3", locName: "L3", content: "L3"},
		{name: "led4", locName: "L4", content: "L4"},
		{name: "led5", locName: "L5", content: "L5"},
		{name: "led6", locName: "L6", content: "L6"},
		{name: "led7", locName: "L7", content: "L7"},
		{name: "led8", locName: "L8", content: "L8"},
		{name: "led9", locName: "L9", content: "L9"},
		{name: "led0", locName: "L0", content: "L0"},
	},
	groupIcons: {
		{name: "bell_on", locName: "Bell On", content: "\uF0F3"},
		{name: "bell_off", locName: "Bell Off", content: "\uF1F6"},
		{name: "question", locName: "Question", content: "?"},
		{name: "exclamation", locName: "Exclamation", content: "!"},
		{name: "plug", locName: "Plug", content: "\uF1E6"},
		{name: "puzzle", locName: "Puzzle", content: "\uF12E"},
		{name: "left", locName: "Left", content: "\uF30A"},
		{name: "right", locName: "Right", content: "\uF30B"},
		{name: "up", locName: "Up", content: "\uF30C"},
		{name: "down", locName: "Down", content: "\uF309"},
		{name: "location", locName: "Location", content: "\uF3C5"},
		{name: "power_exc", locName: "Power Warning", content: "\uE55D"},
		{name: "power_ok", locName: "Power OK", content: "\uE55C"},
		{name: "power_bolt", locName: "Power Bolt", content: "\uE55B"},
		{name: "wireless", locName: "Wireless", content: "\uF1EB"},
		{name: "signal", locName: "Some Signal", content: "\uF012"},
		{name: "broadcast", locName: "Broadcast", content: "\uF519"},
		{name: "satellite", locName: "Satellite", content: "\uF7C0"},
		{name: "satellite2", locName: "Satellite 2", content: "\uF7BF"},
		{name: "micro", locName: "Microchip", content: "\uF2DB"},
		{name: "env_bolt", locName: "Bolt", content: "\uF0E7"},
		{name: "env_fire", locName: "Fire", content: "\uF06D"},
		{name: "env_sun", locName: "Sun", content: "\uF185"},
		{name: "env_water", locName: "Water", content: "\uF773"},
		{name: "env_leaf", locName: "Leaf", content: "\uF06C"},
		{name: "env_wind", locName: "Wind", content: "\uF72E"},
		{name: "env_explosion", locName: "Explosion", content: "\uE4E9"},
		{name: "env_atom", locName: "Atom", content: "\uF5D2"},
		{name: "env_solar", locName: "Solar", content: "\uF5BA"},
		{name: "env_tree", locName: "Tree", content: "\uF1BB"},
		{name: "science", locName: "Science", content: "\uF0C3"},
	},
}

const (
	iconWidth         = 64
	iconHeight        = 64
	previewCellsWidth = 10
	previewCellsSize  = 64
)

var loadedFonts = make(map[string]font.Face)
var loadedImages = make(map[string]image.Image)

func main() {
	totalImages := 0
	for groupID := range groups {
		totalImages += len(groupSignals[groupID])
	}

	previewHeightCells := int(math.Ceil(float64(totalImages) / float64(previewCellsWidth)))
	previewCanvas := gg.NewContext(previewCellsWidth*previewCellsSize, previewHeightCells*previewCellsSize)

	centerX, centerY := iconWidth/2, iconHeight/2
	previewSellX := 0
	previewSellY := 0

	for _, groupID := range groupIDs() {
		group := groups[groupID]

		for _, sig := range groupSignals[groupID] {
			background := loadBackground(group.backgroundFile)
			dc := gg.NewContextForImage(background)

			offsetX, offsetY, fontHeight := float64(0), float64(0), float64(1)
			fontHeightBig := float64(iconWidth) / 2
			fontHeightMd := float64(iconWidth) / 2.75
			fontHeightSm := float64(iconWidth) / 10
			shadowScale := 1.03
			shadowSize := float64(1)

			var fontName string

			switch group.glyphType {
			case glyphTypeMath:
				offsetY -= float64(iconHeight / 10)
				fontName = "gl_type_font_math.ttf"
				fontHeight = fontHeightBig
			case glyphTypeIcon:
				offsetY -= float64(iconHeight / 10)
				fontName = "gl_type_fontawesome.ttf"
				fontHeight = fontHeightBig
			case glyphTypeLabel:
				offsetY -= float64(iconHeight / 20)
				fontName = "gl_type_font_labels.ttf"
				fontHeight = fontHeightMd
			case glyphTypeTinyText:
				fontName = "gl_type_font_labels.ttf"
				fontHeight = fontHeightSm
			}

			for _, id := range sig.fineTune {
				settings := fineTunes[id]
				if settings.textScale != 0 {
					fontHeight *= settings.textScale
				}
				offsetX += (settings.offsetX * 0.01) * fontHeight
				offsetY += (settings.offsetY * 0.01) * fontHeight
			}

			for _, id := range group.fineTune {
				settings := fineTunes[id]
				if settings.textScale != 0 {
					fontHeight *= settings.textScale
				}
				offsetX += (settings.offsetX * 0.01) * fontHeight
				offsetY += (settings.offsetY * 0.01) * fontHeight
			}

			// shadow
			if group.glyphType == glyphTypeTinyText {
				textFont := loadFont(fontName, fontHeight*shadowScale)
				dc.SetFontFace(textFont)
				dc.SetHexColor("#00000015")
				dc.DrawStringAnchored(sig.content, float64(centerX)+offsetX-shadowSize, float64(centerY)+offsetY-shadowSize, 0.5, 0.5)
				dc.DrawStringAnchored(sig.content, float64(centerX)+offsetX+shadowSize, float64(centerY)+offsetY-shadowSize, 0.5, 0.5)
				dc.DrawStringAnchored(sig.content, float64(centerX)+offsetX+shadowSize, float64(centerY)+offsetY+shadowSize, 0.5, 0.5)
				dc.DrawStringAnchored(sig.content, float64(centerX)+offsetX-shadowSize, float64(centerY)+offsetY+shadowSize, 0.5, 0.5)
			}

			// text
			textFont := loadFont(fontName, fontHeight)
			dc.SetFontFace(textFont)
			if group.glyphType == glyphTypeIcon {
				dc.SetHexColor("#00000099")
			} else {
				dc.SetHexColor("#000000ff")
			}
			dc.DrawStringAnchored(sig.content, float64(centerX)+offsetX, float64(centerY)+offsetY, 0.5, 0.5)

			// save
			mip := applyMipmaps(dc)
			err := mip.SavePNG(directoryMod + "graphics/signal/" + group.name + "_" + sig.name + ".png")
			if err != nil {
				panic(fmt.Errorf("failed to save png: %w", err))
			}

			// add to preview
			finalIcon := dc.Image()
			previewCanvas.DrawImage(
				resize.Resize(previewCellsSize, previewCellsSize, finalIcon, resize.Bilinear),
				previewSellX*previewCellsSize,
				previewSellY*previewCellsSize,
			)

			previewSellX++
			if previewSellX >= previewCellsWidth {
				previewSellX = 0
				previewSellY++
			}
		}
	}

	// export preview
	err := previewCanvas.SavePNG(directory + "preview.png")
	if err != nil {
		panic(fmt.Errorf("failed to save png: %w", err))
	}

	// gen lua files
	createLuaGroups()
	createLuaSignals()
	createDefaultLocale()
}

func applyMipmaps(icon *gg.Context) *gg.Context {
	mipWidth := iconWidth + (iconWidth / 2) + (iconWidth / 4) + (iconWidth / 8)
	mip := gg.NewContext(mipWidth, iconHeight)

	iconImage := icon.Image()

	// mip level 0
	mip.DrawImage(iconImage, 0, 0)

	// mip level 1
	mip.DrawImage(resize.Resize(iconWidth/2, iconHeight/2, iconImage, resize.Bilinear), iconWidth, 0)

	// mip level 2
	mip.DrawImage(resize.Resize(iconWidth/4, iconHeight/4, iconImage, resize.Bilinear), iconWidth+(iconWidth/2), 0)

	// mip level 3
	mip.DrawImage(resize.Resize(iconWidth/8, iconHeight/8, iconImage, resize.Bilinear), iconWidth+(iconWidth/2)+(iconWidth/4), 0)

	return mip
}

func groupIDs() []group {
	result := make([]group, 0, len(groups))

	for g := range groups {
		result = append(result, g)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return groups[result[i]].sortOrder <= groups[result[j]].sortOrder
	})

	return result
}

func loadBackground(name string) image.Image {
	if f, ok := loadedImages[name]; ok {
		return f
	}

	img, err := gg.LoadPNG(directory + fmt.Sprintf("in/%s.png", name))
	if err != nil {
		panic(fmt.Errorf("could not load png image: %v", err))
	}

	loadedImages[name] = img
	return img
}

func loadFont(name string, lineHeight float64) font.Face {
	id := fmt.Sprintf("font_%s_%.2f", name, lineHeight)
	if f, ok := loadedFonts[id]; ok {
		return f
	}

	loaded, err := gg.LoadFontFace(directory+fmt.Sprintf("in/%s", name), lineHeight)
	if err != nil {
		panic(fmt.Errorf("could not load font: %v", err))
	}

	loadedFonts[id] = loaded
	return loaded
}

func createLuaGroups() {
	var buf bytes.Buffer

	buf.WriteString("-- Auto generated, do not edit.\n")
	buf.WriteString(fmt.Sprintf("-- Generated at %s\n\n", time.Now()))

	// header
	buf.WriteString("data:extend({")

	// content
	for _, groupID := range groupIDs() {
		meta := groups[groupID]

		buf.WriteString(fmt.Sprintf(`
		{
		  type = "item-subgroup",
		  name = "virtual-signal-vs2-%s",
		  group = "signals",
		  order = "%s"
		},`,
			meta.name,
			fmt.Sprintf("a_vs2[%d]", meta.sortOrder)))
	}

	// footer
	buf.WriteString("\n})")

	// write
	err := os.WriteFile(directoryMod+"groups.lua", buf.Bytes(), 0666)
	if err != nil {
		panic(fmt.Errorf("could not create lua groups: %w", err))
	}
}

func createLuaSignals() {
	var buf bytes.Buffer

	buf.WriteString("-- Auto generated, do not edit.\n")
	buf.WriteString(fmt.Sprintf("-- Generated at %s\n\n", time.Now()))

	// header
	buf.WriteString("data:extend({")

	// content
	for _, groupID := range groupIDs() {
		group := groups[groupID]
		for ind, sig := range groupSignals[groupID] {
			buf.WriteString(fmt.Sprintf(`
				{
				  type = "virtual-signal",
				  name = "signal-vs2-%s",
				  icon = "__virtual-signals2__/graphics/signal/%s_%s.png",
				  subgroup = "virtual-signal-vs2-%s",
				  order = "b[%s]-[%d]"
				},`,
				sig.name,
				group.name,
				sig.name,
				group.name,
				group.name,
				ind))
		}
	}

	// footer
	buf.WriteString("\n})")

	// write
	err := os.WriteFile(directoryMod+"signals.lua", buf.Bytes(), 0666)
	if err != nil {
		panic(fmt.Errorf("could not create lua signals: %w", err))
	}
}

func createDefaultLocale() {
	var buf bytes.Buffer

	// header
	buf.WriteString("[virtual-signal-name]\n")

	// content
	for _, groupID := range groupIDs() {
		for _, sig := range groupSignals[groupID] {
			buf.WriteString(fmt.Sprintf(
				"signal-vs2-%s=Signal %s\n",
				sig.name,
				sig.locName,
			))
		}
	}

	// write
	err := os.WriteFile(directoryMod+"locale/en/auto_gen_signals.cfg", buf.Bytes(), 0666)
	if err != nil {
		panic(fmt.Errorf("could not create lua signals: %w", err))
	}
}
