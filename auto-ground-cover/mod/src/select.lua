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

    if table_size(settings.errors) > 0 then
        player.print("Auto ground cover is miss configured:")

        for _, err in pairs(settings.errors) do
            player.print(" - " .. err)
        end
    end

    -- player.print("default "..settings.defaultCover)
    -- for coverName, entities in pairs(settings.groups) do
    --     player.print("group "..coverName)
    --     for _, entName in pairs(entities) do
    --         player.print("- ent "..entName)
    --     end
    -- end

    mappedTiles = MapEntitiesIntoTiles(event.tiles)
    placedTiles = {}

    --- place tile above buildings
    for _, found in pairs(mappedTiles) do
        if found.ent == nil then
            goto continue
        end

        local at = found.tile.position
        placedTiles[TilePositionHash(at)] = true
        newCover = ResolveCoverageNameProtoName(settings, found)
        newCoverTileName = GetTileName(newCover)

        if newCover == const__DESTRUCT_COVER then
            found.tile.order_deconstruction(player.force, player)
            goto continue
        end

        if found.tile.name == newCover or found.tile.name == newCoverTileName then
            goto continue
        end

        for _, ghost in pairs(found.tile.get_tile_ghosts()) do
            if ghost.prototype.type == "tile-ghost" and (ghost.prototype.name == newCover or ghost.prototype.name == newCoverTileName) then
                goto continue
            end
        end

        CreateTileGhostFromSignal(at, newCover, found.ent.surface, player.force)
        placedTiles[TilePositionHash(at)] = true

        ::continue::
    end

    --- place default cover to all non-occupied places
    newCover = settings.defaultCover
    newCoverTileName = GetTileName(newCover)

    for _, tile in pairs(event.tiles) do
        local at = tile.position

        if placedTiles[TilePositionHash(at)] then
            goto continue
        end

        if newCover == const__DESTRUCT_COVER then
            tile.order_deconstruction(player.force, player)
            goto continue
        end

        if tile.name == newCover or tile.name == newCoverTileName then
            goto continue
        end

        for _, ghost in pairs(tile.get_tile_ghosts()) do
            if ghost.prototype.type ~= "tile-ghost" then
                goto next_ghost
            end

            ---@param LuaTilePrototype
            local ghostProto = ghost.ghost_prototype
            local ghostProtoName = ghostProto.name

            if ghostProtoName == newCover or ghostProtoName == newCoverTileName then
                goto continue
            end

            ::next_ghost::
        end

        CreateTileGhostFromSignal(at, newCover, tile.surface, player.force)

        ::continue::
    end
end

script.on_event(defines.events.on_lua_shortcut, on_lua_shortcut)
script.on_event(defines.events.on_player_selected_area, on_player_selected_area)