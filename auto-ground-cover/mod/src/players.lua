return {
    ---@param playerOrIndex LuaPlayer|number
    ---@return number
    extractPlayerName = function(playerOrIndex)
        if type(playerOrIndex) == "table" or type(playerOrIndex) == "userdata" then
            return playerOrIndex.name
        end

        local player = game.players[playerOrIndex] -- index can be uint or playerName (both work ok)
        if player ~= nil then
            return player.name
        end

        return "unknown-player-index-"..tostring(playerOrIndex)
    end
}