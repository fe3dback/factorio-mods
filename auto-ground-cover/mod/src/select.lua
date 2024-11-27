local const = require("src.classes")
require("src.settings")
local profiles = require("src.profiles")
require("src.apply")

script.on_event(
    defines.events.on_lua_shortcut,
    function(event)
        if event.prototype_name ~= const.ShortCutID then
            game.print(event.prototype_name)
            game.print(const.ShortCutID)
            return
        end

        local player = game.players[event.player_index]
        if player.clear_cursor() then
            local stack = player.cursor_stack
            if player.cursor_stack and stack.can_set_stack({ name = const.SelectToolID }) then
                stack.set_stack({ name = const.SelectToolID })
            end
        end
    end
)

script.on_event(
    defines.events.on_player_selected_area,
    function(event)
        if event.item and event.item ~= const.SelectToolID then
            return
        end

        if table_size(event.tiles) == 0 then
            return
        end

        local settings

        -- todo: remove switch
        if not const.useNewSettings then
            settings = ReadSettingsFromConstantCombinator(event.surface)
        else
            settings = profiles.intoApplySettings(profiles.defaultProfileIds.tmp)
        end

        ---@type ToolContext
        local context = {
            selectedTiles = event.tiles,
            surface = event.surface,
            player = game.players[event.player_index],
            settings = settings,
            topLeft = event.area.left_top,
            width = event.area.right_bottom.x - event.area.left_top.x,
            height = event.area.right_bottom.y - event.area.left_top.y,
        }

        ApplyAutoGroundCover(context)
    end
)