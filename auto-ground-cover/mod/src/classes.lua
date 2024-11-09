const__DESTRUCT_COVER = "fe3dback/auto-ground-cover__destruct_cover__"

do
    ---@class TileData
    ---@field tile LuaTile
    ---@field ent LuaEntity?
end

do
    ---@class GroundCoverSettings
    ---@field defaultCover string            -- coverName (or special const)
    ---@field groups table<string, string[]> -- coverName -> entityNames[]
    ---@field errors string[]
end