package model

import "time"

type ThemePreset struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Name                string    `gorm:"size:64;not null;uniqueIndex" json:"name"`
	BgMode              string    `gorm:"size:16;not null;default:gradient" json:"bg_mode"`
	BgColor             string    `gorm:"size:7" json:"bg_color"`
	BgGradientStart     string    `gorm:"size:7" json:"bg_gradient_start"`
	BgGradientEnd       string    `gorm:"size:7" json:"bg_gradient_end"`
	BgGradientAngle     int       `gorm:"default:135" json:"bg_gradient_angle"`
	BgBlur              bool      `gorm:"default:false" json:"bg_blur"`
	DarkBgMode          string    `gorm:"size:16;not null;default:gradient" json:"dark_bg_mode"`
	DarkBgColor         string    `gorm:"size:7" json:"dark_bg_color"`
	DarkBgGradientStart string    `gorm:"size:7" json:"dark_bg_gradient_start"`
	DarkBgGradientEnd   string    `gorm:"size:7" json:"dark_bg_gradient_end"`
	DarkBgGradientAngle int       `gorm:"default:135" json:"dark_bg_gradient_angle"`
	DarkBgBlur          bool      `gorm:"default:false" json:"dark_bg_blur"`
	PrimaryColor        string    `gorm:"size:7;not null" json:"primary_color"`
	SecondaryColor      string    `gorm:"size:7;not null" json:"secondary_color"`
	AccentColor         string    `gorm:"size:7;not null" json:"accent_color"`
	MutedColor          string    `gorm:"size:7;not null" json:"muted_color"`
	DarkPrimaryColor    string    `gorm:"size:7;not null" json:"dark_primary_color"`
	DarkSecondaryColor  string    `gorm:"size:7;not null" json:"dark_secondary_color"`
	DarkAccentColor     string    `gorm:"size:7;not null" json:"dark_accent_color"`
	DarkMutedColor      string    `gorm:"size:7;not null" json:"dark_muted_color"`
	CardOpacity         int       `gorm:"default:90" json:"card_opacity"`
	CardBlur            int       `gorm:"default:12" json:"card_blur"`
	DarkCardOpacity     int       `gorm:"default:85" json:"dark_card_opacity"`
	DarkCardBlur        int       `gorm:"default:16" json:"dark_card_blur"`
	CreatedAt           time.Time `json:"created_at"`
}

// BuiltinPresets contains the 4 built-in theme presets
var BuiltinPresets = map[string]ThemePreset{
	"default": {
		Name:                "Default",
		BgMode:              "gradient",
		BgColor:             "#f5f5f7",
		BgGradientStart:     "#ffffff",
		BgGradientEnd:       "#e8e8ed",
		BgGradientAngle:     135,
		BgBlur:              false,
		DarkBgMode:          "gradient",
		DarkBgColor:         "#0f1117",
		DarkBgGradientStart: "#1a1c25",
		DarkBgGradientEnd:   "#0f1117",
		DarkBgGradientAngle: 135,
		DarkBgBlur:          false,
		PrimaryColor:        "#7CB99A",
		SecondaryColor:      "#e8e8ed",
		AccentColor:         "#9ABBE0",
		MutedColor:          "#A0AEC0",
		DarkPrimaryColor:    "#8ECFAB",
		DarkSecondaryColor:  "#252833",
		DarkAccentColor:     "#A8C8E8",
		DarkMutedColor:      "#5c6173",
		CardOpacity:         90,
		CardBlur:            12,
		DarkCardOpacity:     85,
		DarkCardBlur:        16,
	},
	"ocean": {
		Name:                "Ocean",
		BgMode:              "gradient",
		BgColor:             "#e0f2f1",
		BgGradientStart:     "#e0f2f1",
		BgGradientEnd:       "#e1f5fe",
		BgGradientAngle:     135,
		BgBlur:              false,
		DarkBgMode:          "gradient",
		DarkBgColor:         "#0d1b2a",
		DarkBgGradientStart: "#0d1b2a",
		DarkBgGradientEnd:   "#1b2838",
		DarkBgGradientAngle: 135,
		DarkBgBlur:          false,
		PrimaryColor:        "#00897B",
		SecondaryColor:      "#e1f5fe",
		AccentColor:         "#0288D1",
		MutedColor:          "#78909C",
		DarkPrimaryColor:    "#4DB6AC",
		DarkSecondaryColor:  "#1b2838",
		DarkAccentColor:     "#4FC3F7",
		DarkMutedColor:      "#546E7A",
		CardOpacity:         90,
		CardBlur:            12,
		DarkCardOpacity:     85,
		DarkCardBlur:        16,
	},
	"spring": {
		Name:                "Spring",
		BgMode:              "gradient",
		BgColor:             "#f1f8e9",
		BgGradientStart:     "#f1f8e9",
		BgGradientEnd:       "#e8f5e9",
		BgGradientAngle:     135,
		BgBlur:              false,
		DarkBgMode:          "gradient",
		DarkBgColor:         "#1b2e1b",
		DarkBgGradientStart: "#1b2e1b",
		DarkBgGradientEnd:   "#152015",
		DarkBgGradientAngle: 135,
		DarkBgBlur:          false,
		PrimaryColor:        "#66BB6A",
		SecondaryColor:      "#e8f5e9",
		AccentColor:         "#AED581",
		MutedColor:          "#8BC34A",
		DarkPrimaryColor:    "#81C784",
		DarkSecondaryColor:  "#1b2e1b",
		DarkAccentColor:     "#C5E1A5",
		DarkMutedColor:      "#558B2F",
		CardOpacity:         90,
		CardBlur:            12,
		DarkCardOpacity:     85,
		DarkCardBlur:        16,
	},
	"sakura": {
		Name:                "Sakura",
		BgMode:              "gradient",
		BgColor:             "#fce4ec",
		BgGradientStart:     "#fce4ec",
		BgGradientEnd:       "#f3e5f5",
		BgGradientAngle:     135,
		BgBlur:              false,
		DarkBgMode:          "gradient",
		DarkBgColor:         "#1f0f18",
		DarkBgGradientStart: "#1f0f18",
		DarkBgGradientEnd:   "#1a1025",
		DarkBgGradientAngle: 135,
		DarkBgBlur:          false,
		PrimaryColor:        "#EC407A",
		SecondaryColor:      "#f3e5f5",
		AccentColor:         "#CE93D8",
		MutedColor:          "#B39DDB",
		DarkPrimaryColor:    "#F06292",
		DarkSecondaryColor:  "#1a1025",
		DarkAccentColor:     "#CE93D8",
		DarkMutedColor:      "#7E57C2",
		CardOpacity:         90,
		CardBlur:            12,
		DarkCardOpacity:     85,
		DarkCardBlur:        16,
	},
}

// ThemePresetFromSettings creates a ThemePreset snapshot from current ActivitySettings
func ThemePresetFromSettings(s ActivitySettings, name string) ThemePreset {
	return ThemePreset{
		Name:                name,
		BgMode:              s.BgMode,
		BgColor:             s.BgColor,
		BgGradientStart:     s.BgGradientStart,
		BgGradientEnd:       s.BgGradientEnd,
		BgGradientAngle:     s.BgGradientAngle,
		BgBlur:              s.BgBlur,
		DarkBgMode:          s.DarkBgMode,
		DarkBgColor:         s.DarkBgColor,
		DarkBgGradientStart: s.DarkBgGradientStart,
		DarkBgGradientEnd:   s.DarkBgGradientEnd,
		DarkBgGradientAngle: s.DarkBgGradientAngle,
		DarkBgBlur:          s.DarkBgBlur,
		PrimaryColor:        s.PrimaryColor,
		SecondaryColor:      s.SecondaryColor,
		AccentColor:         s.AccentColor,
		MutedColor:          s.MutedColor,
		DarkPrimaryColor:    s.DarkPrimaryColor,
		DarkSecondaryColor:  s.DarkSecondaryColor,
		DarkAccentColor:     s.DarkAccentColor,
		DarkMutedColor:      s.DarkMutedColor,
		CardOpacity:         s.CardOpacity,
		CardBlur:            s.CardBlur,
		DarkCardOpacity:     s.DarkCardOpacity,
		DarkCardBlur:        s.DarkCardBlur,
	}
}

// ToSettingsUpdates converts the preset to a map for GORM Updates()
func (p ThemePreset) ToSettingsUpdates() map[string]interface{} {
	return map[string]interface{}{
		"bg_mode":               p.BgMode,
		"bg_color":              p.BgColor,
		"bg_gradient_start":     p.BgGradientStart,
		"bg_gradient_end":       p.BgGradientEnd,
		"bg_gradient_angle":     p.BgGradientAngle,
		"bg_blur":               p.BgBlur,
		"dark_bg_mode":          p.DarkBgMode,
		"dark_bg_color":         p.DarkBgColor,
		"dark_bg_gradient_start": p.DarkBgGradientStart,
		"dark_bg_gradient_end":   p.DarkBgGradientEnd,
		"dark_bg_gradient_angle": p.DarkBgGradientAngle,
		"dark_bg_blur":           p.DarkBgBlur,
		"primary_color":          p.PrimaryColor,
		"secondary_color":        p.SecondaryColor,
		"accent_color":           p.AccentColor,
		"muted_color":            p.MutedColor,
		"dark_primary_color":     p.DarkPrimaryColor,
		"dark_secondary_color":   p.DarkSecondaryColor,
		"dark_accent_color":      p.DarkAccentColor,
		"dark_muted_color":       p.DarkMutedColor,
		"card_opacity":           p.CardOpacity,
		"card_blur":              p.CardBlur,
		"dark_card_opacity":      p.DarkCardOpacity,
		"dark_card_blur":         p.DarkCardBlur,
	}
}
