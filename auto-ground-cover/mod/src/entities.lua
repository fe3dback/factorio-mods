---@param tile LuaEntity
---@return boolean
local function isUsefulEntity(entity)
    return BuildingName(entity) or GhostTileName(entity)
end

---@param tiles LuaTile[]
---@return EntityWithOwnedTiles[]
function FindEntitiesOnTiles(tiles)
    ---@type table<number, EntityWithOwnedTiles>
    local entityById = {}

    for _, tile in pairs(tiles) do
        colliderTL = {x=tile.position.x+0.01, y=tile.position.y+0.01}
        colliderBR = {x=tile.position.x+0.99, y=tile.position.y+0.99}

        above = tile.surface.find_entities({colliderTL, colliderBR})

        for _, entityAbove in pairs(above) do
            if not isUsefulEntity(entityAbove) then
                goto continue
            end

            local entId = entityAbove.unit_number
            if entId == nil then
                goto continue
            end

            if entityById[entId] == nil then
                entityById[entId] = {
                    entity = entityAbove,
                    tiles = {}
                }
            end

            table.insert(entityById[entId].tiles, tile)
            ::continue::
        end
    end

    return entityById
end

--- return entity building prototype name
--- if this entity is building (or ghost of the building), or nil otherwise
---@param entity LuaEntity
---@return string?
function BuildingName(entity)
    -- we interesting in real buildings
    if entity.prototype.is_building then
        return entity.prototype.name
    end

    if entity.type ~= "entity-ghost" then
        return nil
    end

    -- also in ghost buildings
    if entity.ghost_prototype.object_name == 'LuaEntityPrototype' and entity.ghost_prototype.is_building then
        return entity.ghost_prototype.name
    end

    -- otherwise is not building
    return nil
end

--- return tile.proto.name from entity tile-ghost or nil otherwise
---@param entity LuaEntity
---@return string?
function GhostTileName(entity)
    if entity.type ~= "tile-ghost" then
        return nil
    end

    if entity.ghost_prototype.object_name ~= 'LuaTilePrototype' then
        return nil
    end

    return entity.ghost_prototype.name
end