package main

import (
	"fmt"
	"image"
	"math"
	"sort"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
)

const (
	directory = "/home/neo/code/fe3dback/factorio-mods/virtual-signals2/generator/"
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

type render struct {
	group    group
	fileName string
	content  string
	fineTune []fineTuneGroupID
}

type groupMeta struct {
	sortOrder      int
	fileNamePrefix string
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
		fileNamePrefix: "cmn",
		backgroundFile: "sig_red",
		glyphType:      glyphTypeLabel,
		fineTune: []fineTuneGroupID{
			fnScaleSM,
		},
	},
	groupCommon2: {
		sortOrder:      2,
		fileNamePrefix: "cmn2",
		backgroundFile: "sig_green",
		glyphType:      glyphTypeLabel,
		fineTune: []fineTuneGroupID{
			fnScaleSM,
		},
	},
	groupCommon3: {
		sortOrder:      3,
		fileNamePrefix: "cmn3",
		backgroundFile: "sig_lime",
		glyphType:      glyphTypeLabel,
		fineTune: []fineTuneGroupID{
			fnScaleSM,
		},
	},
	groupMath: {
		sortOrder:      10,
		fileNamePrefix: "math",
		backgroundFile: "sig_blue",
		glyphType:      glyphTypeMath,
	},
	groupGreek: {
		sortOrder:      20,
		fileNamePrefix: "greek",
		backgroundFile: "sig_purple",
		glyphType:      glyphTypeMath,
	},
	groupIcons: {
		sortOrder:      30,
		fileNamePrefix: "fa",
		backgroundFile: "sig_light_purple",
		glyphType:      glyphTypeIcon,
	},
}

var renders = map[group][]render{
	groupMath: {
		// expressions
		{fileName: "eq", content: "=", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{fileName: "neq", content: "≠", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{fileName: "aeq", content: "≈", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{fileName: "gt", content: ">", fineTune: []fineTuneGroupID{
			fnScaleMD,
			ftOffsetRightXS,
		}},
		{fileName: "lt", content: "<", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{fileName: "ge", content: "≥", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{fileName: "le", content: "≤", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{fileName: "plus", content: "+", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{fileName: "minus", content: "-", fineTune: []fineTuneGroupID{
			ftOffsetUpLG,
			ftScaleXXL,
		}},
		{fileName: "plus_minus", content: "±", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{fileName: "percent", content: "%", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
		}},
		// geo
		{fileName: "geo_angle", content: "∠", fineTune: []fineTuneGroupID{
			fnScaleMD,
		}},
		{fileName: "geo_perpendicular", content: "⊥", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
		}},
		{fileName: "geo_parallel", content: "∥", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
		}},
		{fileName: "geo_similar", content: "~", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		// algo
		{fileName: "algo_lemni", content: "∞", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{fileName: "algo_func", content: "f", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{fileName: "algo_delta", content: "∆", fineTune: []fineTuneGroupID{
			ftOffsetDownXS,
			fnScaleMD,
		}},
		{fileName: "algo_sigma", content: "∑"},
		{fileName: "algo_golden", content: "φ"},
		{fileName: "algo_pi", content: "π"},
		{fileName: "algo_integral", content: "∫"},
	},
	groupGreek: {
		{fileName: "alpha", content: "α"},
		{fileName: "beta", content: "β"},
		{fileName: "gamma", content: "γ"},
		{fileName: "delta", content: "δ"},
		{fileName: "epsilon", content: "ε"},
		{fileName: "zeta", content: "ζ"},
		{fileName: "eta", content: "η"},
		{fileName: "theta", content: "θ"},
		{fileName: "iota", content: "ι"},
		{fileName: "kappa", content: "κ"},
		{fileName: "lambda", content: "λ"},
		{fileName: "mu", content: "μ"},
		{fileName: "nu", content: "ν"},
		{fileName: "xi", content: "ξ"},
		{fileName: "omicron", content: "ο"},
		{fileName: "pi", content: "π"},
		{fileName: "rho", content: "ρ"},
		{fileName: "sigma", content: "σ"},
		{fileName: "tau", content: "τ"},
		{fileName: "upsilon", content: "υ"},
		{fileName: "phi", content: "φ"},
		{fileName: "chi", content: "χ"},
		{fileName: "psi", content: "ψ"},
		{fileName: "omega", content: "ω"},
	},
	groupCommon: {
		{fileName: "on", content: "ON"},
		{fileName: "off", content: "OFF"},
		{fileName: "min", content: "MIN"},
		{fileName: "max", content: "MAX"},
		{fileName: "in", content: "IN"},
		{fileName: "out", content: "OUT"},
	},
	groupCommon2: {
		{fileName: "size", content: "SIZE"},
		{fileName: "buffer", content: "BUF"},
		{fileName: "limit", content: "LIM"},
		{fileName: "capacity", content: "CAP"},
		{fileName: "nil", content: "NIL"},
		{fileName: "fuck", content: "FUCK", fineTune: []fineTuneGroupID{fnUnScaleSM}},
	},
	groupCommon3: {
		{fileName: "size_small", content: "SM"},
		{fileName: "size_medium", content: "MD"},
		{fileName: "size_large", content: "LG"},
		{fileName: "size_extra_large", content: "XL"},
		{fileName: "size_xxl", content: "XXL"},
		{fileName: "sev_info", content: "INF"},
		{fileName: "sev_warn", content: "WRN", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{fileName: "sev_err", content: "ERR"},
		{fileName: "sev_critical", content: "CRT"},
		{fileName: "alert", content: "ALRT", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{fileName: "trigger", content: "TRIG", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{fileName: "if", content: "IF"},
		{fileName: "else", content: "ELSE", fineTune: []fineTuneGroupID{fnUnScaleSM}},
		{fileName: "action", content: "ACT"},
		{fileName: "var_x", content: "[X]"},
		{fileName: "var_y", content: "[Y]"},
		{fileName: "var_z", content: "[Z]"},
		{fileName: "var_w", content: "[W]"},
		{fileName: "color_r", content: "(R)"},
		{fileName: "color_g", content: "(G)"},
		{fileName: "color_b", content: "(B)"},
		{fileName: "led1", content: "L1"},
		{fileName: "led2", content: "L2"},
		{fileName: "led3", content: "L3"},
		{fileName: "led4", content: "L4"},
		{fileName: "led5", content: "L5"},
		{fileName: "led6", content: "L6"},
		{fileName: "led7", content: "L7"},
		{fileName: "led8", content: "L8"},
		{fileName: "led9", content: "L9"},
		{fileName: "led0", content: "L0"},
	},
	groupIcons: {
		{fileName: "bell_on", content: "\uF0F3"},
		{fileName: "bell_off", content: "\uF1F6"},
		{fileName: "question", content: "?"},
		{fileName: "exclamation", content: "!"},
		{fileName: "plug", content: "\uF1E6"},
		{fileName: "puzzle", content: "\uF12E"},
		{fileName: "left", content: "\uF30A"},
		{fileName: "right", content: "\uF30B"},
		{fileName: "up", content: "\uF30C"},
		{fileName: "down", content: "\uF309"},
		{fileName: "location", content: "\uF3C5"},
		{fileName: "location", content: "\uF3C5"},
		{fileName: "power_exc", content: "\uE55D"},
		{fileName: "power_ok", content: "\uE55C"},
		{fileName: "power_bolt", content: "\uE55B"},
		{fileName: "power_bolt", content: "\uE55B"},
		{fileName: "wireless", content: "\uF1EB"},
		{fileName: "signal", content: "\uF012"},
		{fileName: "broadcast", content: "\uF519"},
		{fileName: "satellite", content: "\uF7C0"},
		{fileName: "satellite2", content: "\uF7BF"},
		{fileName: "micro", content: "\uF2DB"},
		{fileName: "env_bolt", content: "\uF0E7"},
		{fileName: "env_fire", content: "\uF06D"},
		{fileName: "env_sun", content: "\uF185"},
		{fileName: "env_water", content: "\uF773"},
		{fileName: "env_leaf", content: "\uF06C"},
		{fileName: "env_wind", content: "\uF72E"},
		{fileName: "env_explosion", content: "\uE4E9"},
		{fileName: "env_atom", content: "\uF5D2"},
		{fileName: "env_solar", content: "\uF5BA"},
		{fileName: "env_tree", content: "\uF1BB"},
		{fileName: "science", content: "\uF0C3"},
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
		totalImages += len(renders[groupID])
	}

	previewHeightCells := int(math.Ceil(float64(totalImages) / float64(previewCellsWidth)))
	previewCanvas := gg.NewContext(previewCellsWidth*previewCellsSize, previewHeightCells*previewCellsSize)

	centerX, centerY := iconWidth/2, iconHeight/2
	previewSellX := 0
	previewSellY := 0

	for _, groupID := range groupIDs() {
		meta := groups[groupID]

		for _, rnd := range renders[groupID] {
			background := loadBackground(meta.backgroundFile)
			dc := gg.NewContextForImage(background)

			offsetX, offsetY, fontHeight := float64(0), float64(0), float64(1)
			fontHeightBig := float64(iconWidth) / 2
			fontHeightMd := float64(iconWidth) / 2.75
			fontHeightSm := float64(iconWidth) / 10
			shadowScale := 1.03
			shadowSize := float64(1)

			var fontName string

			switch meta.glyphType {
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

			for _, id := range rnd.fineTune {
				settings := fineTunes[id]
				if settings.textScale != 0 {
					fontHeight *= settings.textScale
				}
				offsetX += (settings.offsetX * 0.01) * fontHeight
				offsetY += (settings.offsetY * 0.01) * fontHeight
			}

			for _, id := range meta.fineTune {
				settings := fineTunes[id]
				if settings.textScale != 0 {
					fontHeight *= settings.textScale
				}
				offsetX += (settings.offsetX * 0.01) * fontHeight
				offsetY += (settings.offsetY * 0.01) * fontHeight
			}

			// shadow
			if meta.glyphType == glyphTypeTinyText {
				textFont := loadFont(fontName, fontHeight*shadowScale)
				dc.SetFontFace(textFont)
				dc.SetHexColor("#00000015")
				dc.DrawStringAnchored(rnd.content, float64(centerX)+offsetX-shadowSize, float64(centerY)+offsetY-shadowSize, 0.5, 0.5)
				dc.DrawStringAnchored(rnd.content, float64(centerX)+offsetX+shadowSize, float64(centerY)+offsetY-shadowSize, 0.5, 0.5)
				dc.DrawStringAnchored(rnd.content, float64(centerX)+offsetX+shadowSize, float64(centerY)+offsetY+shadowSize, 0.5, 0.5)
				dc.DrawStringAnchored(rnd.content, float64(centerX)+offsetX-shadowSize, float64(centerY)+offsetY+shadowSize, 0.5, 0.5)
			}

			// text
			textFont := loadFont(fontName, fontHeight)
			dc.SetFontFace(textFont)
			if meta.glyphType == glyphTypeIcon {
				dc.SetHexColor("#00000099")
			} else {
				dc.SetHexColor("#000000ff")
			}
			dc.DrawStringAnchored(rnd.content, float64(centerX)+offsetX, float64(centerY)+offsetY, 0.5, 0.5)

			// save
			mip := applyMipmaps(dc)
			err := mip.SavePNG(directory + "out/" + meta.fileNamePrefix + "_" + rnd.fileName + ".png")
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
