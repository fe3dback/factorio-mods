require("src.classes")

---@param tiles LuaTile[]
---@return table<MapPosition, TileData>
function MapEntitiesIntoTiles(tiles)
    result = {}

    for _, tile in pairs(tiles) do
        above = tile.surface.find_entities({tile.position, {x=tile.position.x+1, y=tile.position.y+1}})

        if table_size(above) == 0 then
            result[tile.position] = {
                tile = tile
            }
            goto continueMain
        end

        for _, ent in pairs(above) do
            if ent.prototype.type == "tile-ghost" then
                goto continue
            end

            result[tile.position] = {
                ent = ent, 
                tile = tile
            }

            ::continue::
        end

        ::continueMain::
    end

    return result
end

---@param sig string
---@return string?
function GetTileName(sig)
    local item_proto = prototypes.item[sig]
    
    if item_proto and item_proto.place_as_tile_result then
        return item_proto.place_as_tile_result.result.name
    end

    return nil
end


---@param at TilePosition
---@param sig string
---@param surface LuaSurface
---@param force string|integer|LuaForce
function CreateTileGhostFromSignal(at, sig, surface, force)
    local tile_name = GetTileName(sig)
    if tile_name == nil then
        return
    end
    
    surface.create_entity({
        name = "tile-ghost",
        inner_name = tile_name,
        position = at,
        force = force,
    })
end

---@param position TilePosition
---@return string
function TilePositionHash(position)
    return position.x .. ";" .. position.y
end