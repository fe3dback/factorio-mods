require("src.classes")

---@param data TileData
---@param settings GroundCoverSettings
---@return string
function ResolveCoverageNameProtoName(settings, data)
    if data.ent ~= nil then
        for coverName, entNames in pairs(settings.groups) do
            for _, entName in pairs(entNames) do
                if data.ent.prototype.name == entName then
                    return coverName
                end
            end
        end
    end

    return settings.defaultCover
end

---@param sig SignalFilter
---@return boolean
local function isSignalWithTileCoverage(sig)
    if sig.type ~= "item" then
        return false
    end

    ---@type LuaItemPrototype
    local proto = prototypes.item[sig.name]
    
    if not (proto ~= nil and proto.place_as_tile_result) then
        return false
    end

    return true
end

---@param sig SignalFilter
---@return boolean
local function isSignalWithEntity(sig)
    if sig.type ~= "item" then
        return false
    end

    if isSignalWithTileCoverage(sig) then
        return false
    end

    return true
end

---@param surface LuaSurface
---@return GroundCoverSettings
function GetSettingsFromConstCombinator(surface)
    combinators = surface.find_entities_filtered{
        type = 'constant-combinator'
    }

    ---@type GroundCoverSettings
    settings = {
        defaultCover = const__DESTRUCT_COVER,
        groups = {},
        errors = {}
    }

    for _, combinator in pairs(combinators) do
        if combinator.combinator_description ~= "AUTO-GROUND-COVER" then
            goto continue
        end

        ---@type LuaConstantCombinatorControlBehavior
        local behavior = combinator.get_control_behavior()

        if not (behavior ~= nil and behavior.enabled) then
            goto continue
        end

        for _, section in pairs(behavior.sections) do
            logSectionName = "Section [" .. section.index .. "] (" .. section.object_name .. ")"

            if not (section.active and section.is_manual) then
                goto next_section
            end

            local currentCoverage = nil
            local currentEntities = {}

            for ind, slot in pairs(section.filters) do
                -- we want only table entity signals
                if not (slot.value ~= nil and slot.value.name) then
                    goto next_slot
                end

                ---@type SignalFilter
                local signal = slot.value
                local isCoverage = isSignalWithTileCoverage(signal)
                local isEntity = isSignalWithEntity(signal)

                logSlotName = "Slot [" .. ind .. "] (" .. slot.value.name .. ")"

                if ind == 1 then
                    if not isCoverage then
                        table.insert(settings.errors, logSectionName .. " " .. logSlotName .. " | First Slot of each section MUST contain signal of ground cover tile")
                        goto next_section
                    end

                    if section.filters_count == 1 then
                        -- if this is section with one slot (ground coverage)
                        -- this is default ground coverage
                        settings.defaultCover = signal.name
                        goto next_section
                    end

                    currentCoverage = signal.name
                else
                    if not isEntity then
                        -- invalid section
                        -- we want entity signal in all other slots
                        table.insert(settings.errors, logSectionName .. " " .. logSlotName .. " | Slot MUST contain placeable building entity signal, but is not")
                        goto next_section
                    end

                    -- ok, this is valid group
                    table.insert(currentEntities, signal.name)
                end

                ::next_slot::
            end

            if currentCoverage ~= nil and table_size(currentEntities) > 0 then
                if settings.groups[currentCoverage] == nil then
                    settings.groups[currentCoverage] = {}
                end

                for _, ent in pairs(currentEntities) do
                    table.insert(settings.groups[currentCoverage], ent)
                end
            end

            ::next_section::
        end

        ::continue::
    end

    return settings
end
