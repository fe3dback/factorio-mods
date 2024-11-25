local const = require("src.classes")
local dataApi = require("src.storage")
local hudEvents = require("src.hud_events")


---@param parent LuaGuiElement
local function hud_add_section(player, parent, id)
    local sectionId = parent.name .. "-" .. id

    -- header

    local headerFrame = parent.add{
        name=sectionId.."-id",
        type="frame",
        style="subheader_frame"
    }

    local headerContent = headerFrame.add{
        name=sectionId.."-header",
        type="flow",
        direction="horizontal",
    }
    headerContent.style.horizontally_stretchable = true

    --local headerTurnOnSwitcher = headerContent.add{
    --    name=sectionId.."-switcher",
    --    type="checkbox",
    --    caption="group name",
    --    state=true,
    --    style="subheader_caption_checkbox"
    --}
    --
    --local headerDeleteButton = headerContent.add{
    --    name=sectionId.."-delete-btn",
    --    type="button",
    --    caption="group name",
    --    state=true,
    --    style="tool_button_red"
    --}

    local headerCoverButton = headerContent.add{
        name=sectionId.."-select-cover",
        type="choose-elem-button",
        caption="Select ground cover",
        elem_type="tile",
        elem_filters={
            {filter="item-to-place", mode="and"},
            {filter="minable", mode="and"}
        },
        style="slot_button"
    }

    local headerBorderCoverButton = headerContent.add{
        name=sectionId.."-select-border-cover",
        type="choose-elem-button",
        caption="Select border ground cover",
        elem_type="tile",
        elem_filters={
            {filter="item-to-place", mode="and"},
            {filter="minable", mode="and"}
        },
        style="slot_button"
    }

    local headerBorderRadiusSlider = headerContent.add{
        name=sectionId.."-radius-slider",
        type="slider",
        caption="Border radius",
        minimum_value=0,
        maximum_value=4,
        value=0,
        style="notched_slider"
    }

    local headerBorderRadiusValue = headerContent.add{
        name=sectionId.."-radius-value",
        type="textfield",
        caption="Border radius",
        numeric=true,
        style="slider_value_textfield",
        text=tostring(0)
    }
    headerBorderRadiusValue.enabled = false

    hudEvents.onValueChanged(player, headerBorderRadiusSlider.name, function(ev)
        headerBorderRadiusValue.text = tostring(ev.element.slider_value)
    end)

    -- content

    local slotsContent = parent.add{
        name=sectionId.."-table",
        type="table",
        column_count=8,
    }

    for slotY = 1, 4 do
        for slotX = 1, 8 do
            local slotId = tostring(slotX) .. "-" .. tostring(slotY)
            local slot = slotsContent.add{
                name=sectionId.."-slot-"..slotId,
                type="choose-elem-button",
                caption="Select building",
                elem_type="item",
                elem_filters={
                    {filter="place-result", elem_filters={
                        {filter="building", mode="and"}
                    }},
                },
                style="slot_button"
            }
        end
    end

    --parent.add{
    --    name=const.hud.frames.settings.widgets.btnTest,
    --    type="choose-elem-button",
    --    elem_type="tile",
    --    elem_filters={
    --        {filter="item-to-place", mode="and"},
    --        {filter="minable", mode="and"}
    --    },
    --}
end

---@param player LuaPlayer
local function hud_create(player)
    local hud = player.gui.screen
    if hud[const.hud.frames.settings.id] ~= nil then
        -- already exist
        return
    end

    local frameID = const.hud.frames.settings.id
    local frame = player.gui.screen.add{
        name=frameID,
        type="frame",
        caption="Cover Settings",
        auto_center=true,
        style=const.hud.frames.settings.style.frame,
    }
    local frameX = dataApi.player.gui.getOr(player, const.hud.frames.settings.id.."-x", 80)
    local frameY = dataApi.player.gui.getOr(player, const.hud.frames.settings.id.."-y", 240)

    frame.location = {x=frameX, y=frameY}
    frame.bring_to_front()

    local containerFrame = frame.add{
        name=frameID.."-container",
        type="frame",
        style="entity_frame",
    }

    local containerContent = containerFrame.add{
        name=const.hud.frames.settings.widgets.content,
        type="flow",
        direction="vertical"
    }

    hud_add_section(player, containerContent, "1")
    hud_add_section(player, containerContent, "2")
end

---@param element LuaGuiElement?
local function hud_destroy_elem_recursive(element)
    if element == nil then
        return
    end

    for _, child in pairs(element.children_names) do
        hud_destroy_elem_recursive(element[child])
    end

    element.destroy()
end

---@param player LuaPlayer
local function hud_destroy(player)
    local hud = player.gui.screen
    hud_destroy_elem_recursive(hud[const.hud.frames.settings.id])
    hudEvents.resetEvents(player)
end

script.on_event(
    defines.events.on_player_cursor_stack_changed,
    function(event)
        local player = game.get_player(event.player_index)
        local cursor = player.cursor_stack

        --if not (cursor and cursor.valid and cursor.valid_for_read) then
        --    hud_destroy(player)
        --    return
        --end
        --
        --if cursor.name == const.SelectToolID then
        --    hud_create(player)
        --end
    end
)

script.on_event(
    defines.events.on_gui_location_changed,
    function(event)
        if event.element == nil then
            return
        end

        if event.element.name ~= const.hud.frames.settings.id then
            return
        end

        dataApi.player.gui.set(event.player_index, const.hud.frames.settings.id.."-x", event.element.location.x)
        dataApi.player.gui.set(event.player_index, const.hud.frames.settings.id.."-y", event.element.location.y)
    end
)
