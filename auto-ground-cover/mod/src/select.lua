require("src.classes")
require("src.resolver")
require("src.utils")

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

    local player = game.players[event.player_index]
    local settings = GetSettingsFromConstCombinator(event.tiles[1].surface)

    mappedTiles = MapEntitiesIntoTiles(event.tiles)

    for _, found in pairs(mappedTiles) do
        newCover = ResolveCoverageNameProtoName(settings, found)

        if newCover == const__DESTRUCT_COVER then
            found.tile.order_deconstruction(player.force, player)
            goto continue
        end

        if found.tile.prototype.name == newCover then
            goto continue
        end

        for _, ghost in pairs(found.tile.get_tile_ghosts()) do
            if ghost.prototype.type == "tile-ghost" and ghost.prototype.name == newCover then
                goto continue
            end
        end

        if found.ent ~= nil then
            found.ent.surface.create_entity({
                name = "tile-ghost",
                inner_name = newCover,
                position = found.tile.position,
                force = player.force,
            })
        end

        ::continue::
    end
end

script.on_event(defines.events.on_lua_shortcut, on_lua_shortcut)
script.on_event(defines.events.on_player_selected_area, on_player_selected_area)