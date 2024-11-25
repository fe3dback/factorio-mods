local const = require("src.classes")
require("src.entities")
require("src.tiles")

local colorMissConfTitle = { r = 0.95, g = 0, b = 0, a = 0.8}
local colorMissConfText = { r = 0.3, g = 0.6, b = 0.95, a = 1}

---@param context ToolContext
---@return boolean
local function ensureSettingsIsValid(context)
    if table_size(context.settings.errors) == 0 then
        return true
    end

    context.player.print("Auto ground cover is miss configured:", {
        color = colorMissConfTitle,
    })

    for _, err in pairs(context.settings.errors) do
        context.player.print(" - " .. err, {
            color = colorMissConfText,
        })
    end

    return false
end

---@param context ToolContext
local function debugDumpSettings(context)
    -- todo: remove debug log
    --context.player.print("default "..context.settings.defaultCover)
    --for key, group in pairs(context.settings.groups) do
    --    context.player.print("group ["..key.."] = "..group.cover)
    --
    --    for _, building in pairs(group.buildings) do
    --        context.player.print("- building ".. building.entityName.. " (border=" .. tostring(building.borderSize) .. ")")
    --    end
    --end
end

---@param settings GroundCoverSettings
---@return table<string, EntitySettings>        -- entityProto.name (buildingName) -> EntitySettings
local function buildCoverMap(settings)
    local result = {}

    for _, group in pairs(settings.groups) do
        for _, building in pairs(group.buildings) do
            result[building.entityName] = {
                coverName = group.cover,
                coverBorderName = group.coverBorder,
                borderSize = building.borderSize,
            }
        end
    end

    return result
end

---@param surface LuaSurface
---@param tiles LuaTile[]
---@param borderSize number
---@return LuaTile[]
local function calculateBorderMap(surface, tiles, borderSize)
    local minY = 999999999999
    local minX = 999999999999
    local maxX = -99999999999
    local maxY = -99999999999

    for _, tile in pairs(tiles) do
        local pos = tile.position

        if pos.x < minX then
            minX = pos.x
        end
        if pos.y < minY then
            minY = pos.y
        end
        if pos.x > maxX then
            maxX = pos.x
        end
        if pos.y > maxY then
            maxY = pos.y
        end
    end

    local result = {}
    for y = minY-borderSize, maxY+borderSize, 1 do
        for x = minX-borderSize, maxX+borderSize, 1 do
            table.insert(result, surface.get_tile(x, y))
        end
    end

    return result
end

local function isRail(entity)
    if entity == nil then
        return false
    end

    if entity.prototype == nil then
        return false
    end

    if entity.prototype.collision_mask == nil then
        return false
    end

    if entity.prototype.collision_mask.layers == nil then
        return false
    end

    return entity.prototype.collision_mask.layers["rail"]
end

---@param context ToolContext
---@return ResolveResult
local function resolveGroundCoverMap(context)
    ---@type ResolveResult
    local result = {
        mapPrimary = {},
        mapBorders = {},
        ghostTiles = {},
    }

    ---@type EntityWithOwnedTiles[]
    local entities = FindEntitiesOnTiles(context.selectedTiles)

    for _, entityWithTiles in pairs(entities) do
        local ghostTile = GhostTileName(entityWithTiles.entity)
        if ghostTile == nil then
            -- is not ghost tile
            goto nextGhost
        end

        for _, childTile in pairs(entityWithTiles.tiles) do
            local tileHash = TilePositionHash(childTile.position)
            result.ghostTiles[tileHash] = {
                ghostTileName = ghostTile,
                ghostEntity = entityWithTiles.entity,
            }
        end

        ::nextGhost::
    end

    ---@type table<string, EntitySettings>
    -- (entityProto.name (buildingName) -> EntitySettings
    local coverMap = buildCoverMap(context.settings)

    -- apply all buildings into cover map
    for _, entityWithTiles in pairs(entities) do
        local buildingName = BuildingName(entityWithTiles.entity)
        if buildingName == nil then
            goto continue
        end

        -- all rails is different entity, but we assume that user want cover for all rail types
        if isRail(entityWithTiles.entity) then
            -- todo: custom tiles mask for each rail type (+border support)
            buildingName = "straight-rail"
        end

        local coverProps = coverMap[buildingName]
        if coverProps == nil then
            goto continue
        end

        -- add tiles directly below building to primary map
        for _, childTile in pairs(entityWithTiles.tiles) do
            local tileHash = TilePositionHash(childTile.position)
            result.mapPrimary[tileHash] = {
                coverName = coverProps.coverName,
                ownedByBuilding = buildingName,
                tile = childTile,
            }
        end

        -- add tiles around building (border) to border map
        local borderMap = calculateBorderMap(context.surface, entityWithTiles.tiles, coverProps.borderSize)
        for _, borderTile in pairs(borderMap) do
            local tileHash = TilePositionHash(borderTile.position)
            result.mapBorders[tileHash] = {
                coverName = coverProps.coverBorderName,
                ownedByBuilding = buildingName,
                tile = borderTile,
            }
        end

        ::continue::
    end

    return result
end


---@param context ToolContext
function ApplyAutoGroundCover(context)
    if not ensureSettingsIsValid(context) then
        return result
    end

    debugDumpSettings(context)

    local coverMap = resolveGroundCoverMap(context)

    for _, selectedTile in pairs(context.selectedTiles) do
        local tileHash = TilePositionHash(selectedTile.position)

        --- resolve new cover

        ---@type ResolveTileResult
        local cover

        if coverMap.mapPrimary[tileHash] ~= nil then
            cover = coverMap.mapPrimary[tileHash]
        elseif coverMap.mapBorders[tileHash] ~= nil then
            cover = coverMap.mapBorders[tileHash]
        else
            cover = {
                tile = selectedTile,
                coverName = context.settings.defaultCover,
                ownedByBuilding = nil,
            }
        end

        --- check current cover
        if selectedTile.collides_with("water_tile") then
            -- todo: option to landfill water instead of ignore
            goto next
        end

        -- if tile already has this cover - nothing to-do
        if selectedTile.prototype.name == cover.coverName then
            goto next
        end

        -- if tile already has this ghost cover - nothing to-do
        ---@type ResolveResultGhostTile? -- tileProto.name if exist
        local ghostTile = coverMap.ghostTiles[tileHash]

        if ghostTile ~= nil then
            -- already same
            if ghostTile.ghostTileName == cover.coverName then
                goto next
            end

            -- we have another ghost, but it should be changed to new ghost
            -- deconstruct current ghost, and place it to new ghost in next calls
            ghostTile.ghostEntity.order_deconstruction(context.player.force, context.player)
        end

        --- apply new cover (deconstruct if not set)

        if cover.coverName == const.CoverTypeSpecialDestruct then
            selectedTile.order_deconstruction(context.player.force, context.player)
            goto next
        end

        --- apply new cover (add build order)

        context.surface.create_entity({
            name = "tile-ghost",
            inner_name = cover.coverName,
            position = selectedTile.position,
            force = context.player.force,
        })

        ::next::

    end
end