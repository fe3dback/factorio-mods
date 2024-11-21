local prefix = "fe3dback__autogc__"

local hudSettings = prefix .. "hud/frame/settings/"
local hudSettingsWidgets = prefix .. "hud/frame/settings/widgets/"
local hudSettingsStyles = prefix .. "hud/frame/settings/styles/"

local const = {
    -- tools
    ShortCutID = prefix .. "shortcut",
    SelectToolID = prefix .. "select_tool",

    -- combinator search id
    CombinatorSearchDescription = "AUTO-GROUND-COVER",

    -- special cover
    CoverTypeSpecialDestruct = prefix .. "cover-type/destruct",

    -- signal types
    SignalTypeTileProto = prefix .. "signal-type/tile-proto",
    SignalTypeBuildingProto = prefix .. "signal-type/building-proto",
    SignalTypeOther = prefix .. "signal-type/other",

    -- hud
    hud = {
        frames = {
            settings = {
                id = hudSettings .. "id",

                widgets = {
                    containerContent = hudSettingsWidgets .. "container-content",
                    containerContentFlow = hudSettingsWidgets .. "container-content-flow",
                    btnTest = hudSettingsWidgets .. "test-btn",
                },

                style = {
                    Frame = hudSettingsStyles .. "frame"
                }
            },
        },
    },
}

do
    ---@class DetectedSlogType
    ---@field signalType string   -- const "SIGNAL_TYPE"
    ---@field resolvedName string -- for item-tile signals, this is TilePrototype.Name
end

do
    ---@class GroundCoverSettings
    ---@field defaultCover string                            -- TilePrototype.Name (or special const)
    ---@field groups table<string, GroundCoverSettingsGroup> -- TilePrototype.Name -> group
    ---@field errors string[]                                -- tool should not work if not empty
end

do
    ---@class GroundCoverSettingsGroup
    ---@field cover string                             -- TilePrototype.Name
    ---@field buildings GroundCoverSettingsBuilding[]  -- array of building settings
end

do
    ---@class GroundCoverSettingsBuilding
    ---@field entityName string                        -- EntityPrototype.Name
    ---@field borderSize number                        -- 0=zero, 1= +one tile in every direction (2x2 -> 4x4)
end

do
    ---@class ToolContext
    ---@field selectedTiles LuaTile[]
    ---@field surface LuaSurface
    ---@field player LuaPlayer
    ---@field settings GroundCoverSettings
    ---@field topLeft TilePosition
    ---@field width number
    ---@field height number
end

do
    ---@class EntitySettings
    ---@field coverName string              -- TilePrototype.Name (or special const)
    ---@field borderSize number             -- 0=zero, 1= +one tile in every direction (2x2 -> 4x4)
end

do
    ---@class ResolveResult
    ---@field mapPrimary table<string, ResolveTileResult>      -- tileHash -> ResolveTileResult
    ---@field mapBorders table<string, ResolveTileResult>      -- tileHash -> ResolveTileResult
    ---@field ghostTiles table<string, ResolveResultGhostTile> -- tileHash -> ResolveResultGhostTile
end

do
    ---@class ResolveResultGhostTile
    ---@field ghostTileName string    -- ghostTile.ghost_prototype(.tile).name
    ---@field ghostEntity LuaEntity
end

do
    ---@class ResolveTileResult
    ---@field tile LuaTile              -- tile ref
    ---@field coverName string          -- coverName (tileProto.name)
    ---@field ownedByBuilding string    -- name of building that want`s to place this tile cover (entityProto.name)
end

do
    ---@class EntityWithOwnedTiles
    ---@field entity LuaEntity
    ---@field tiles LuaTile[]
end

return const