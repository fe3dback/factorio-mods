local playerUtils = require("src.players")

local eventStorage = {
    players = {}
}

local function ensurePlayerStorageExist(playerName)
    if eventStorage.players[playerName] ~= nil then
        return
    end

    eventStorage.players[playerName] = {
        on_gui_text_changed = {},
        on_gui_value_changed = {}
    }
end

local function ensurePlayerStorageDropped(playerName)
    if eventStorage.players[playerName] == nil then
        return
    end

    eventStorage.players[playerName] = nil
end

script.on_event(
    defines.events.on_gui_value_changed,
    function(event)
        local player = game.get_player(event.player_index)
        local playerName = playerUtils.extractPlayerName(player)

        ensurePlayerStorageExist(playerName)
        local playerStorage = eventStorage.players[playerName].on_gui_value_changed

        for elemName, callback in pairs(playerStorage) do
            if event.element.name == elemName then
                callback(event)
                return
            end
        end
    end
)

return {
    onValueChanged = function(playerOrIndex, guiElemName, func)
        local playerName = playerUtils.extractPlayerName(playerOrIndex)
        ensurePlayerStorageExist(playerName)
        eventStorage.players[playerName].on_gui_value_changed[guiElemName] = func
    end,

    resetEvents = function(playerOrIndex)
        local playerName = playerUtils.extractPlayerName(playerOrIndex)
        ensurePlayerStorageDropped(playerName)
    end
}