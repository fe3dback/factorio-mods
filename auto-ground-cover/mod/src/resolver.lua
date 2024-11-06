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

-- ---@param sig SignalFilter
-- ---@return boolean
-- local function isSignalWithTileCoverage(sig)
--     if sig.type ~= "item" then
--         return false
--     end

--     ---@type LuaEntityPrototype
--     local proto = prototypes.entity[sig.name]
--     if proto == nil then
--         log("sig " .. sig.name .. " not exist in tile table:")

--         for key in pairs(prototypes.entity) do
--             log(" - " .. key)
--         end

--         return false
--     end

--     log("sig:"..sig.name.." = "..proto.name)
    
--     return true
-- end

-- ---@param sig SignalFilter
-- ---@return boolean
-- local function isSignalWithEntity(sig)
--     if sig.type ~= "item" then
--         return false
--     end

--     if isSignalWithTileCoverage(sig) then
--         return false
--     end

--     return true
-- end

---@param surface LuaSurface
---@return GroundCoverSettings
function GetSettingsFromConstCombinator(surface)
    combinators = surface.find_entities_filtered{
        type = 'constant-combinator'
    }

    ---@type GroundCoverSettings
    settings = {
        defaultCover = const__DESTRUCT_COVER,
        groups = {}
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
            log("begin section")

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
                -- local isCoverage = isSignalWithTileCoverage(signal)
                -- local isEntity = isSignalWithEntity(signal)

                -- log("sig ".. signal.name .. "| isCover=" .. tostring(isCoverage) .. ", isEntity=" .. tostring(isEntity))

                if ind == 1 then
                    -- if not isCoverage then
                    --     -- invalid section
                    --     -- we want ground coverage signal at first slot
                    --     goto next_section
                    -- end

                    if section.filters_count == 1 then
                        -- if this is section with one slot (ground coverage)
                        -- this is default ground coverage
                        settings.defaultCover = signal.name
                        goto next_section
                    end

                    currentCoverage = signal.name
                else
                    -- if not isEntity then
                    --     -- invalid section
                    --     -- we want entity signal in all other slots
                    --     goto next_section
                    -- end

                    -- ok, this is valid group
                    table.insert(currentEntities, signal.name)
                end

                ::next_slot::
            end

            if currentCoverage ~= nil and table_size(currentEntities) > 0 then
                settings.groups[currentCoverage] = currentEntities
            end

            ::next_section::
        end

        ::continue::
    end

    return settings
end
