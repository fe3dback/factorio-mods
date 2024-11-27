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
        on_gui_value_changed = {},
        on_gui_elem_changed = {},
    }
end

local function ensurePlayerStorageDropped(playerName)
    if eventStorage.players[playerName] == nil then
        return
    end

    eventStorage.players[playerName] = nil
end

local function triggerEventIn(eventsTable, event)
    local player = game.get_player(event.player_index)
    local playerName = playerUtils.extractPlayerName(player)

    ensurePlayerStorageExist(playerName)
    local playerStorage = eventStorage.players[playerName][eventsTable]

    for elemName, callback in pairs(playerStorage) do
        if event.element.name == elemName then
            callback(event)
            return
        end
    end
end

script.on_event(
    defines.events.on_gui_text_changed,
    function(event)
        triggerEventIn("on_gui_text_changed", event)
    end
)

script.on_event(
    defines.events.on_gui_value_changed,
    function(event)
        triggerEventIn("on_gui_value_changed", event)
    end
)

script.on_event(
    defines.events.on_gui_elem_changed,
    function(event)
        triggerEventIn("on_gui_elem_changed", event)
    end
)

return {
    onValueChanged = function(playerOrIndex, guiElemName, func)
        local playerName = playerUtils.extractPlayerName(playerOrIndex)
        ensurePlayerStorageExist(playerName)
        eventStorage.players[playerName].on_gui_value_changed[guiElemName] = func
    end,

    onElemChoose = function(playerOrIndex, guiElemName, func)
        local playerName = playerUtils.extractPlayerName(playerOrIndex)
        ensurePlayerStorageExist(playerName)
        eventStorage.players[playerName].on_gui_elem_changed[guiElemName] = func
    end,

    resetEvents = function(playerOrIndex)
        local playerName = playerUtils.extractPlayerName(playerOrIndex)
        ensurePlayerStorageDropped(playerName)
    end
}