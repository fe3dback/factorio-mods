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
		{name: "eq", content: "=", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "neq", content: "≠", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "aeq", content: "≈", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "gt", content: ">", fineTune: []fineTuneGroupID{
			fnScaleMD,
			ftOffsetRightXS,
		}},
		{name: "lt", content: "<", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{name: "ge", content: "≥", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{name: "le", content: "≤", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{name: "plus", content: "+", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "minus", content: "-", fineTune: []fineTuneGroupID{
			ftOffsetUpLG,
			ftScaleXXL,
		}},
		{name: "plus_minus", content: "±", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{name: "percent", content: "%", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
		}},
		// geo
		{name: "geo_angle", content: "∠", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{name: "geo_perpendicular", content: "⊥", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
		}},
		{name: "geo_parallel", content: "∥", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
		}},
		{name: "geo_similar", content: "~", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		// algo
		{name: "algo_lemni", content: "∞", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "algo_func", content: "f", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "algo_delta", content: "∆", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{name: "algo_sigma", content: "∑"},
		{name: "algo_golden", content: "φ"},
		{name: "algo_pi", content: "π"},
		{name: "algo_integral", content: "∫"},
	},
	groupGreek: {
		{name: "alpha", content: "α"},
		{name: "beta", content: "β"},
		{name: "gamma", content: "γ"},
		{name: "delta", content: "δ"},
		{name: "epsilon", content: "ε"},
		{name: "zeta", content: "ζ"},
		{name: "eta", content: "η"},
		{name: "theta", content: "θ"},
		{name: "iota", content: "ι"},
		{name: "kappa", content: "κ"},
		{name: "lambda", content: "λ"},
		{name: "mu", content: "μ"},
		{name: "nu", content: "ν"},
		{name: "xi", content: "ξ"},
		{name: "omicron", content: "ο"},
		{name: "pi", content: "π"},
		{name: "rho", content: "ρ"},
		{name: "sigma", content: "σ"},
		{name: "tau", content: "τ"},
		{name: "upsilon", content: "υ"},
		{name: "phi", content: "φ"},
		{name: "chi", content: "χ"},
		{name: "psi", content: "ψ"},
		{name: "omega", content: "ω"},
	},
	groupCommon: {
		{name: "on", content: "ON"},
		{name: "off", content: "OFF"},
		{name: "min", content: "MIN"},
		{name: "max", content: "MAX"},
		{name: "in", content: "IN"},
		{name: "out", content: "OUT"},
	},
	groupCommon2: {
		{name: "size", content: "SIZE"},
		{name: "buffer", content: "BUF"},
		{name: "limit", content: "LIM"},
		{name: "capacity", content: "CAP"},
		{name: "nil", content: "NIL"},
		{name: "fuck", content: "FUCK", fineTune: []fineTuneGroupID{fnUnScaleSM}},
	},
	groupCommon3: {
		{name: "size_small", content: "SM"},
		{name: "size_medium", content: "MD"},
		{name: "size_large", content: "LG"},
		{name: "size_extra_large", content: "XL"},
		{name: "size_xxl", content: "XXL"},
		{name: "sev_info", content: "INF"},
		{name: "sev_warn", content: "WRN", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{name: "sev_err", content: "ERR"},
		{name: "sev_critical", content: "CRT"},
		{name: "alert", content: "ALRT", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{name: "trigger", content: "TRIG", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{name: "if", content: "IF"},
		{name: "else", content: "ELSE", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{name: "action", content: "ACT"},
		{name: "var_x", content: "[X]"},
		{name: "var_y", content: "[Y]"},
		{name: "var_z", content: "[Z]"},
		{name: "var_w", content: "[W]"},
		{name: "color_r", content: "(R)"},
		{name: "color_g", content: "(G)"},
		{name: "color_b", content: "(B)"},
		{name: "led1", content: "L1"},
		{name: "led2", content: "L2"},
		{name: "led3", content: "L3"},
		{name: "led4", content: "L4"},
		{name: "led5", content: "L5"},
		{name: "led6", content: "L6"},
		{name: "led7", content: "L7"},
		{name: "led8", content: "L8"},
		{name: "led9", content: "L9"},
		{name: "led0", content: "L0"},
	},
	groupIcons: {
		{name: "bell_on", content: "\uF0F3"},
		{name: "bell_off", content: "\uF1F6"},
		{name: "question", content: "?"},
		{name: "exclamation", content: "!"},
		{name: "plug", content: "\uF1E6"},
		{name: "puzzle", content: "\uF12E"},
		{name: "left", content: "\uF30A"},
		{name: "right", content: "\uF30B"},
		{name: "up", content: "\uF30C"},
		{name: "down", content: "\uF309"},
		{name: "location", content: "\uF3C5"},
		{name: "location", content: "\uF3C5"},
		{name: "power_exc", content: "\uE55D"},
		{name: "power_ok", content: "\uE55C"},
		{name: "power_bolt", content: "\uE55B"},
		{name: "power_bolt", content: "\uE55B"},
		{name: "wireless", content: "\uF1EB"},
		{name: "signal", content: "\uF012"},
		{name: "broadcast", content: "\uF519"},
		{name: "satellite", content: "\uF7C0"},
		{name: "satellite2", content: "\uF7BF"},
		{name: "micro", content: "\uF2DB"},
		{name: "env_bolt", content: "\uF0E7"},
		{name: "env_fire", content: "\uF06D"},
		{name: "env_sun", content: "\uF185"},
		{name: "env_water", content: "\uF773"},
		{name: "env_leaf", content: "\uF06C"},
		{name: "env_wind", content: "\uF72E"},
		{name: "env_explosion", content: "\uE4E9"},
		{name: "env_atom", content: "\uF5D2"},
		{name: "env_solar", content: "\uF5BA"},
		{name: "env_tree", content: "\uF1BB"},
		{name: "science", content: "\uF0C3"},
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
