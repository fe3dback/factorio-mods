package main

const (
	groupMath = iota
	groupGreek
	groupIcons
	groupCommon
	groupCommon2
	groupCommon3
	groupCommon4
)

type groupData struct {
	sortOrder      int
	name           string
	backgroundFile string
	glyphType      glyphType
	fineTune       []fineTuneGroupID
}

var groups = map[group]groupData{
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
		backgroundFile: "sig_gold",
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
	groupCommon4: {
		sortOrder:      50,
		name:           "cmn_digits",
		backgroundFile: "sig_light_blue",
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
