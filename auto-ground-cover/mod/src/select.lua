require("src.classes")
require("src.settings")
require("src.apply")



function on_lua_shortcut(event)
    if event.prototype_name ~= 'fe3dback__autogc__short_cut' then
        return
    end

    local player = game.players[event.player_index]
    if player.clear_cursor() then
        local stack = player.cursor_stack
        if player.cursor_stack and stack.can_set_stack({ name = 'fe3dback__autogc__select_tool' }) then
            stack.set_stack({ name = 'fe3dback__autogc__select_tool' })
        end
    end
end

function on_player_selected_area(event)
    if event.item and event.item ~= 'fe3dback__autogc__select_tool' then
        return
    end

    if table_size(event.tiles) == 0 then
        return
    end

    ---@type ToolContext
    local context = {
        selectedTiles = event.tiles,
        surface = event.surface,
        player = game.players[event.player_index],
        settings = ReadSettingsFromConstantCombinator(event.surface),
        topLeft = event.area.left_top,
        width = event.area.right_bottom.x - event.area.left_top.x,
        height = event.area.right_bottom.y - event.area.left_top.y,
    }

    ApplyAutoGroundCover(context)
end

script.on_event(defines.events.on_lua_shortcut, on_lua_shortcut)
script.on_event(defines.events.on_player_selected_area, on_player_selected_area)