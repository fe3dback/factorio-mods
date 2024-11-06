data:extend(
{
	{
		type = "shortcut",
		name = "fe3dback__autogc__short_cut",
		action = "lua",
		icon = "__base__/graphics/icons/signal/signal_A.png",
		icon_size = 64,
		small_icon = "__base__/graphics/icons/signal/signal_A.png",
		small_icon_size = 64
	},
    {
        type = "selection-tool",
        name = "fe3dback__autogc__select_tool",
        icon = "__base__/graphics/icons/signal/signal_A.png",
        icon_size = 64,
        subgroup = "tool",
        stack_size = 1,
        stackable = false,
        draw_label_for_cursor_render = true,
        select = {
            border_color = { r = 0, g = 0, b = 1 },
            mode = { "buildable-type", "any-tile" },
            cursor_box_type = "entity",
            entity_filter_mode = "blacklist",
            entity_type_filters = {  }
        },
        alt_select = {
            border_color = { r = 0, g = 0, b = 1 },
            mode = { "buildable-type", "any-tile" },
            cursor_box_type = "entity",
            entity_filter_mode = "blacklist",
            entity_type_filters = {  }
        },
        always_include_tiles = true,
        mouse_cursor = "selection-tool-cursor",
        skip_fog_of_war = false
    }
})