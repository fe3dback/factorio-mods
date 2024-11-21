local const = require("src.classes")
local dataApi = require("src.storage")

---@param parent LuaGuiElement
local function hud_add_section(parent, id)
    local sectionId = parent.name .. "-" .. id

    local sectionContent = parent.add{
        name=sectionId.."-content",
        type="flow",
        direction="vertical",
    }

    local section = sectionContent.add{
        name=sectionId.."-id",
        type="frame",
        style="subheader_frame"
    }

    local header = section.add{
        name=sectionId.."-header",
        type="flow",
        direction="horizontal",
    }
    header.style.horizontally_stretchable = true

    local headerTurnOnSwitcher = header.add{
        name=sectionId.."-switcher",
        type="checkbox",
        caption="group name",
        state=true,
        style="subheader_caption_checkbox"
    }

    local headerDeleteButton = header.add{
        name=sectionId.."-delete-btn",
        type="button",
        caption="group name",
        state=true,
        style="tool_button_red"
    }

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

    local frame = player.gui.screen.add{
        name=const.hud.frames.settings.id,
        type="frame",
        caption="Cover Settings",
        auto_center=true,
        style=const.hud.frames.settings.style.Frame,
    }
    local frameX = dataApi.player.gui.getOr(player, const.hud.frames.settings.id.."-x", 80)
    local frameY = dataApi.player.gui.getOr(player, const.hud.frames.settings.id.."-y", 240)

    frame.location = {x=frameX, y=frameY}
    frame.bring_to_front()

    local containerContent = frame.add{
        name=const.hud.frames.settings.widgets.containerContent,
        type="frame",
        style="inside_shallow_frame_with_padding",
    }

    local containerContentFlow = containerContent.add{
        name=const.hud.frames.settings.widgets.containerContent,
        type="flow",
        direction="vertical",
        style="two_module_spacing_vertical_flow"
    }

    hud_add_section(containerContentFlow, "1")
    hud_add_section(containerContentFlow, "2")
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
