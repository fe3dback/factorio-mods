local tables = {
    gui = "gui"
}

---@param playerOrIndex LuaPlayer|number
---@return number
local function extractPlayerName(playerOrIndex)
    if type(playerOrIndex) == "table" or type(playerOrIndex) == "userdata" then
        return playerOrIndex.name
    end

    local player = game.players[playerOrIndex] -- index can be uint or playerName (both work ok)
    if player ~= nil then
        return player.name
    end

    return "unknown-player-index-"..tostring(playerOrIndex)
end

local function playerSet(playerOrIndex, tableName, key, data)
    local playerName = extractPlayerName(playerOrIndex)

    if storage.players == nil then
        storage.players = {}
    end

    if storage.players[playerName] == nil then
        storage.players[playerName] = {}
    end

    if storage.players[playerName][tableName] == nil then
        storage.players[playerName][tableName] = {}
    end

    storage.players[playerName][tableName][key] = data
end

local function playerGet(playerOrIndex, tableName, key)
    local playerName = extractPlayerName(playerOrIndex)

    local players = storage.players
    if players == nil then
        return nil
    end

    local player = players[playerName]
    if player == nil then
        return nil
    end

    local dataTable = player[tableName]
    if dataTable == nil then
        return nil
    end

    return dataTable[key]
end

local function playerGetOrDefault(playerOrIndex, tableName, key, default)
    local playerInx = extractPlayerName(playerOrIndex)
    local data = playerGet(playerInx, tableName, key)
    if data ~= nil then
        return data
    end

    return default
end

return {
    player = {
        gui = {
            ---@param playerOrIndex LuaPlayer|number
            set = function(playerOrIndex, key, data)
                playerSet(playerOrIndex, tables.gui, key, data)
            end,

            ---@param playerOrIndex LuaPlayer|number
            get = function(playerOrIndex, key)
                return playerGet(playerOrIndex, tables.gui, key)
            end,

            ---@param playerOrIndex LuaPlayer|number
            getOr = function(playerOrIndex, key, default)
                return playerGetOrDefault(playerOrIndex, tables.gui, key, default)
            end,
        }
    },
}