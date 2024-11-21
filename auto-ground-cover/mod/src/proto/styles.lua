local const = require("src.classes")

styles = {}
styles.setup = function()
    local default_gui = data.raw["gui-style"].default

    -- Settings window
    -- =============================================

    local style = const.hud.frames.settings.style

    default_gui[style.Frame] = {
        type = "frame_style",
        minimal_width = 300,
        minimal_height = 500,
        top_padding = 4,
        right_padding = 4,
        bottom_padding = 4,
        left_padding = 4,
        use_header_filler = true,
        drag_by_title = true,
    }
end

return styles