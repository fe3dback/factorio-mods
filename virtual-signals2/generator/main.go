package main

import (
	"bytes"
	"fmt"
	"image"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
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
	groupID         string
)

const (
	glyphTypeMath glyphType = iota
	glyphTypeIcon
	glyphTypeLabel
	glyphTypeTinyText
)

const (
	ftOffsetUp = iota
	ftOffsetDown
	ftOffsetRight
	ftScaleSM
	ftUnScaleSM
	ftUnScaleMD
	ftScaleMD
	ftScaleLG
	ftScaleXXL
	ftUseLabelFont
	ftUseIconFont
)

type fineTune struct {
	offsetX      float64 // offset in % of glyph height
	offsetY      float64 // offset in % of glyph height
	textScale    float64 // 1 = no scale
	useLabelFont bool
	useIconFont  bool
}

var fineTunes = map[fineTuneGroupID]fineTune{
	ftOffsetUp:     {offsetY: -20},
	ftOffsetDown:   {offsetY: 5},
	ftOffsetRight:  {offsetX: 5},
	ftScaleSM:      {textScale: 1.1},
	ftUnScaleSM:    {textScale: 0.9},
	ftUnScaleMD:    {textScale: 0.8},
	ftScaleMD:      {textScale: 1.35},
	ftScaleLG:      {textScale: 1.5},
	ftScaleXXL:     {textScale: 1.75},
	ftUseLabelFont: {useLabelFont: true},
	ftUseIconFont:  {useIconFont: true},
}

const (
	iconWidth         = 64
	iconHeight        = 64
	previewCellsWidth = 10
	previewCellsSize  = 64
)

var loadedFonts = make(map[string]font.Face)
var loadedImages = make(map[string]image.Image)
var groupSignals = map[groupID][]signal{}

func main() {
	loadSignalsFromCSV()

	tmp := groupSignals
	_ = tmp

	totalImages := 0
	for groupID := range groups {
		totalImages += len(groupSignals[groupID])
	}

	previewHeightCells := int(math.Ceil(float64(totalImages) / float64(previewCellsWidth)))
	previewHeightCells += len(groups)
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

				if settings.useLabelFont {
					fontName = "gl_type_font_labels.ttf"
				}
				if settings.useIconFont {
					fontName = "gl_type_fontawesome.ttf"
				}
			}

			for _, id := range group.fineTune {
				settings := fineTunes[id]
				if settings.textScale != 0 {
					fontHeight *= settings.textScale
				}
				offsetX += (settings.offsetX * 0.01) * fontHeight
				offsetY += (settings.offsetY * 0.01) * fontHeight

				if settings.useLabelFont {
					fontName = "gl_type_font_labels.ttf"
				}
				if settings.useIconFont {
					fontName = "gl_type_fontawesome.ttf"
				}
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
				dc.SetHexColor("#000044AC")
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

		if previewSellX != 0 {
			previewSellX = 0
			previewSellY++
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
	createLocales()
}

func loadSignalsFromCSV() {
	data, err := os.ReadFile(directory + "in/signals.csv")
	if err != nil {
		panic(fmt.Errorf("failed to load csv: %w", err))
	}

	lines := strings.Split(string(data), "\n")
	for lineNumber, line := range lines {
		// skip header
		if lineNumber == 0 {
			continue
		}

		// skip separators
		if strings.HasPrefix(line, "--") {
			continue
		}

		parts := strings.Split(line, ",")
		for partInd := range parts {
			parts[partInd] = strings.TrimSpace(parts[partInd])
		}

		dataGroupID := groupID(parts[0])
		dataSignalID := parts[1]
		dataContent := parts[2]
		dataFineTuneUnScaleXS := parts[3]
		dataFineTuneScaleXS := parts[4]
		dataFineTuneScaleMD := parts[5]
		dataFineTuneScaleXL := parts[6]
		dataFineTuneMoveDown := parts[7]
		dataFineTuneMoveUp := parts[8]
		dataFineTuneMoveRight := parts[9]
		dataFineTuneUseIconFont := parts[10]
		dataFineTuneUseLabelFont := parts[11]
		dataLocaleEn := parts[12]
		dataLocaleRu := parts[13]

		if strings.HasPrefix(dataContent, "\\u") {
			iconRune, err := strconv.ParseInt(dataContent[2:], 16, 32)
			if err != nil {
				panic(fmt.Errorf("failed to parse unicode '%s' icon at line %d: %w", dataContent, lineNumber, err))
			}

			dataContent = string(rune(iconRune))
		}

		if _, exist := groupSignals[dataGroupID]; !exist {
			groupSignals[dataGroupID] = make([]signal, 0, 64)
		}

		sig := signal{
			group:    dataGroupID,
			name:     dataSignalID,
			localeEn: dataLocaleEn,
			localeRu: dataLocaleRu,
			content:  dataContent,
			fineTune: make([]fineTuneGroupID, 0, 8),
		}

		if dataFineTuneMoveDown == "Y" {
			sig.fineTune = append(sig.fineTune, ftOffsetDown)
		}
		if dataFineTuneMoveUp == "Y" {
			sig.fineTune = append(sig.fineTune, ftOffsetUp)
		}
		if dataFineTuneMoveRight == "Y" {
			sig.fineTune = append(sig.fineTune, ftOffsetRight)
		}

		if dataFineTuneUnScaleXS == "Y" {
			sig.fineTune = append(sig.fineTune, ftUnScaleSM)
		}
		if dataFineTuneScaleXS == "Y" {
			sig.fineTune = append(sig.fineTune, ftScaleSM)
		}
		if dataFineTuneScaleMD == "Y" {
			sig.fineTune = append(sig.fineTune, ftScaleMD)
		}
		if dataFineTuneScaleXL == "Y" {
			sig.fineTune = append(sig.fineTune, ftScaleXXL)
		}

		if dataFineTuneUseIconFont == "Y" {
			sig.fineTune = append(sig.fineTune, ftUseIconFont)
		}
		if dataFineTuneUseLabelFont == "Y" {
			sig.fineTune = append(sig.fineTune, ftUseLabelFont)
		}

		groupSignals[dataGroupID] = append(groupSignals[dataGroupID], sig)
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

func groupIDs() []groupID {
	result := make([]groupID, 0, len(groups))

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
	buf.WriteString(fmt.Sprintf("-- Generated at %s\n\n", time.Now().Format("2006-01-02")))

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
	buf.WriteString(fmt.Sprintf("-- Generated at %s\n\n", time.Now().Format("2006-01-02")))

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
				  order = "b[%s]-[%03d]"
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

func createLocales() {
	locales := []string{"en", "ru"}

	for _, localeID := range locales {
		var buf bytes.Buffer

		// header
		buf.WriteString("[virtual-signal-name]\n")

		// content
		for _, groupID := range groupIDs() {
			for _, sig := range groupSignals[groupID] {
				var value = ""
				var prefix = ""

				switch localeID {
				case "en":
					value = sig.localeEn
					prefix = "Signal"
				case "ru":
					value = sig.localeRu
					prefix = "Сигнал"
				}

				if prefix == "" || value == "" {
					panic("unexpected locale")
				}

				if value == "-" {
					value = sig.localeEn
				}

				buf.WriteString(fmt.Sprintf(
					"signal-vs2-%s=%s %s\n",
					sig.name,
					prefix,
					value,
				))
			}

			// write
			err := os.WriteFile(directoryMod+fmt.Sprintf("locale/%s/auto_gen_signals.cfg", localeID), buf.Bytes(), 0666)
			if err != nil {
				panic(fmt.Errorf("could not create lua signals: %w", err))
			}
		}
	}
}
