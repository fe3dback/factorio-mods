---@param pos TilePosition
---@return string
function TilePositionHash(pos)
    return "x=" .. pos.x .. ";y=" .. pos.y
end