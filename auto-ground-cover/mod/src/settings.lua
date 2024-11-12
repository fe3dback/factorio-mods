require("src.classes")

---@param surface LuaSurface
---@return LuaConstantCombinatorControlBehavior?
local function findMyCombinator(surface)
    local combinators = surface.find_entities_filtered{
        type = 'constant-combinator'
    }

    for _, combinator in pairs(combinators) do
        if combinator.combinator_description == const__COMBINATOR_SEARCH_DESCRIPTION then
            ---@type LuaConstantCombinatorControlBehavior?
            local behavior = combinator.get_control_behavior()

            if behavior ~= nil and behavior.enabled then
                return behavior
            end
        end
    end
end

---@param signal SignalFilter
---@return DetectedSlogType
local function detectSignalType(signal)
    if signal.type ~= "item" then
        return { signalType = const__SIGNAL_TYPE__OTHER }
    end

    ---@type LuaItemPrototype
    local itemProto = prototypes.item[signal.name]
    if itemProto == nil then
        -- unreachable place (if sig.type is item, this is not possible)
        return { signalType = const__SIGNAL_TYPE__OTHER }
    end

    -- item signals can be split into [entities(buildings, other..), tiles and others]

    -- [for tiles]
    if itemProto.place_as_tile_result ~= nil then
        ---@type LuaTilePrototype
        local tileProto = itemProto.place_as_tile_result.result

        return {
            signalType = const__SIGNAL_TYPE__TILE_PROTO,
            resolvedName = tileProto.name,
        }
    -- [for entities(buildings)]
    elseif itemProto.place_result ~= nil then
        ---@type LuaEntityPrototype
        local entityProto = itemProto.place_result

        if entityProto.is_building then
            return {
                signalType = const__SIGNAL_TYPE__BUILDING_PROTO,
                resolvedName = entityProto.name,
            }
        end
    end

    -- [for other]
    return { signalType = const__SIGNAL_TYPE__OTHER }
end

---@param surface LuaSurface
---@return GroundCoverSettings
function ReadSettingsFromConstantCombinator(surface)
    ---@type GroundCoverSettings
    local settings = {
        defaultCover = const__DESTRUCT_COVER,
        errors = {},
        groups = {},
    }

    ---@type table<string, string>
    local alreadySet = {}

    local defaultCoverIsFound = false

    local combinator = findMyCombinator(surface)
    if combinator == nil then
        table.insert(settings.errors, "Not found enabled \"const combinator\" with label \"AUTO-GROUND-COVER\". You need to setup cover rules via signals, see mod readme for details.")
        return settings
    end

    for _, section in pairs(combinator.sections) do
        if not (section.active and section.is_manual) then
            goto next_section
        end

        local logSectionName = "Section [" .. section.index .. "]"

        ---@type string
        local sectionCover -- TilePrototype.name

        ---@type GroundCoverSettingsBuilding[]
        local sectionBuildings = {}

        for ind, slot in pairs(section.filters) do
            local signalQuantity = 0

            if not (slot.value ~= nil and slot.value.name) then
                goto next_slot
            end

            if slot.min ~= nil then
                signalQuantity = slot.min
            end

            local logSlotName = "Slot [" .. ind .. "] (" .. slot.value.name .. ")"
            local detectedSlot = detectSignalType(slot.value)

            if detectedSlot.signalType == const__SIGNAL_TYPE__OTHER then
                table.insert(settings.errors, logSectionName .. " " .. logSlotName .. " | unexpected signal type. Only Tile and building signals is expected")
                goto next_section
            end

            -- tile or entity prototype name (depend on signalType)
            local prototypeName = detectedSlot.resolvedName

            if ind == 1 then
                if detectedSlot.signalType ~= const__SIGNAL_TYPE__TILE_PROTO then
                    table.insert(settings.errors, logSectionName .. " " .. logSlotName .. " | each section must contain TILE cover signal at first slot")
                    goto next_section
                end

                -- if no other signals exist in this section, is default cover
                if section.filters_count == 1 then
                    if not defaultCoverIsFound then
                        defaultCoverIsFound = true
                        settings.defaultCover = prototypeName
                        goto next_section
                    end

                    -- we already have another default cover
                    table.insert(settings.errors, logSectionName .. " " .. logSlotName .. " | have duplicated default cover (only one section with single slot must exist)")
                    goto next_section
                end

                -- otherwise is section cover
                sectionCover = prototypeName
                goto next_slot
            else
                if detectedSlot.signalType == const__SIGNAL_TYPE__TILE_PROTO then
                    table.insert(settings.errors, logSectionName .. " " .. logSlotName .. " | TILE cover signal must be only in first section slot")
                    goto next_section
                end

                if detectedSlot.signalType == const__SIGNAL_TYPE__BUILDING_PROTO then
                    if alreadySet[prototypeName] then
                        table.insert(settings.errors, logSectionName .. " " .. logSlotName .. " | already exist in other section - " .. alreadySet[prototypeName])
                        goto next_section
                    end

                    table.insert(sectionBuildings, {
                        entityName = prototypeName,
                        borderSize = signalQuantity,
                    })
                    alreadySet[prototypeName] = logSectionName .. " " .. logSlotName
                    goto next_slot
                end
            end

            ::next_slot::
        end

        if sectionCover ~= nil and table_size(sectionBuildings) > 0 then
            if settings.groups[sectionCover] == nil then
                settings.groups[sectionCover] = {
                    cover = sectionCover,
                    buildings = {}
                }
            end

            for _, building in pairs(sectionBuildings) do
                table.insert(settings.groups[sectionCover].buildings, building)
            end
        end

        ::next_section::
    end

    return settings
end