package main

const (
	groupCommon       = "cmn"
	groupMath         = "math"
	groupCommon2      = "cmn2"
	groupCommon3      = "cmn3"
	groupGreek        = "greek"
	groupIcons        = "fa"
	groupCommonDigits = "cmn_digits"
)

type group struct {
	sortOrder      int
	name           string
	backgroundFile string
	glyphType      glyphType
	fineTune       []fineTuneGroupID
}

type signal struct {
	group    groupID
	name     string
	localeEn string
	localeRu string
	content  string
	fineTune []fineTuneGroupID
}

var groups = map[groupID]group{
	groupCommon: {
		sortOrder:      1,
		name:           groupCommon,
		backgroundFile: "sig_red",
		glyphType:      glyphTypeLabel,
		fineTune: []fineTuneGroupID{
			ftScaleSM,
		},
	},
	groupCommon2: {
		sortOrder:      2,
		name:           groupCommon2,
		backgroundFile: "sig_gold",
		glyphType:      glyphTypeLabel,
		fineTune: []fineTuneGroupID{
			ftScaleSM,
		},
	},
	groupCommon3: {
		sortOrder:      3,
		name:           groupCommon3,
		backgroundFile: "sig_lime",
		glyphType:      glyphTypeLabel,
		fineTune: []fineTuneGroupID{
			ftScaleSM,
		},
	},
	groupMath: {
		sortOrder:      10,
		name:           groupMath,
		backgroundFile: "sig_blue",
		glyphType:      glyphTypeMath,
	},
	groupGreek: {
		sortOrder:      20,
		name:           groupGreek,
		backgroundFile: "sig_purple",
		glyphType:      glyphTypeMath,
	},
	groupIcons: {
		sortOrder:      30,
		name:           groupIcons,
		backgroundFile: "sig_light_purple",
		glyphType:      glyphTypeIcon,
	},
	groupCommonDigits: {
		sortOrder:      50,
		name:           groupCommonDigits,
		backgroundFile: "sig_light_blue",
		glyphType:      glyphTypeLabel,
		fineTune: []fineTuneGroupID{
			ftScaleSM,
		},
	},
}
