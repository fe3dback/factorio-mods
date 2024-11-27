local const = require("src.classes")

do
    ---@class SettingsProfile
    ---@field name           string
    ---@field defaultCover   string?    -- TilePrototype.Name (or null)
    ---@field sections SettingsProfileSection[]
end

do
    ---@class SettingsProfileSection
    ---@field id           number
    ---@field borderRadius number
    ---@field coverMain    string? -- TilePrototype.Name (or null)
    ---@field coverBorder  string? -- TilePrototype.Name (or null)
    ---@field slots        SettingsProfileSlot[]
end

do
    ---@class SettingsProfileSlot
    ---@field gridX    number
    ---@field gridY    number
    ---@field building string? -- ItemPrototype.Name (or null)
end

local constDefaultProfileIDs = {
    tmp = "__tmp"
}

---@type SettingsProfile
local profileDefaultTmp = {
    name = constDefaultProfileIDs.tmp,
    defaultCover = "concrete",
    sections = {
        {
            id = "1",
            borderRadius = 1,
            coverMain = "stone-path",
            slots = {
                {gridX=1, gridY=1, building="transport-belt"},
                {gridX=2, gridY=1, building="fast-transport-belt"},
                {gridX=3, gridY=1, building="express-transport-belt"},
                {gridX=4, gridY=1, building="turbo-transport-belt"},
                {gridX=1, gridY=2, building="rail"},
            }
        },
        {
            id = "2",
            borderRadius = 0,
            coverMain = "refined-hazard-concrete-left",
            slots = {
                {gridX=1, gridY=1, building="underground-belt"},
                {gridX=2, gridY=1, building="fast-underground-belt"},
                {gridX=3, gridY=1, building="express-underground-belt"},
                {gridX=4, gridY=1, building="turbo-underground-belt"},
                {gridX=1, gridY=2, building="splitter"},
                {gridX=2, gridY=2, building="fast-splitter"},
                {gridX=3, gridY=2, building="express-splitter"},
                {gridX=4, gridY=2, building="turbo-splitter"},
                {gridX=1, gridY=3, building="accumulator"},
                {gridX=2, gridY=3, building="substation"},
                {gridX=3, gridY=3, building="lightning-rod"},
                {gridX=4, gridY=3, building="lightning-collector"},
                {gridX=5, gridY=3, building="radar"},
                {gridX=6, gridY=3, building="small-lamp"},
            }
        },
        {
            id = "3",
            borderRadius = 1,
            coverMain = "refined-concrete",
            coverBorder = "stone-path",
            slots = {
                {gridX=1, gridY=1, building="steel-chest"},
                {gridX=2, gridY=1, building="passive-provider-chest"},
                {gridX=3, gridY=1, building="storage-chest"},
                {gridX=4, gridY=1, building="buffer-chest"},
                {gridX=5, gridY=1, building="requester-chest"},
                {gridX=6, gridY=1, building="active-provider-chest"},
                {gridX=1, gridY=2, building="assembling-machine-1"},
                {gridX=2, gridY=2, building="assembling-machine-2"},
                {gridX=3, gridY=2, building="assembling-machine-3"},
            }
        },
    }
}

local profiles = {}
profiles[constDefaultProfileIDs.tmp] = profileDefaultTmp

---@param name string
---@return SettingsProfile?
local function getProfile(name)
    return profiles[name]
end

---@param tileProtoName string?
local function coverOrDestruct(tileProtoName)
    if tileProtoName == nil then
        return const.CoverTypeSpecialDestruct
    end

    return tileProtoName
end

---@param tileProtoName string?
---@param defaultProtoName string?
local function coverOrDefault(tileProtoName, defaultProtoName)
    if tileProtoName == nil then
        return defaultProtoName
    end

    return tileProtoName
end

---@param profileName string
---@return GroundCoverSettings
local function intoApplySettings(profileName)
    -- todo: remove GroundCoverSettings and use profiles in apply code

    ---@type GroundCoverSettings
    local settings = {
        errors = {},
        groups = {},
    }

    local profile = getProfile(profileName)
    if profile == nil then
        table.insert(settings.errors, "not found profile "..profileName)
        return settings
    end

    settings.defaultCover = coverOrDestruct(profile.defaultCover)

    for _, section in ipairs(profile.sections) do
        ---@type GroundCoverSettingsGroup
        local group = {}

        group.cover = coverOrDestruct(section.coverMain)
        group.coverBorder = coverOrDefault(section.coverBorder, group.cover)
        group.buildings = {}

        for _, slot in ipairs(section.slots) do
            ---@type GroundCoverSettingsBuilding
            local building = {}

            building.borderSize = section.borderRadius
            building.entityName = slot.building

            table.insert(group.buildings, building)
        end

        table.insert(settings.groups, group)
    end

    return settings
end

return {
    defaultProfileIds = constDefaultProfileIDs,
    getProfile = getProfile,
    intoApplySettings = intoApplySettings,
}