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

