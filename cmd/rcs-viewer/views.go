package main

import (
	"image/color"

	"github.com/blackchip-org/retro-cs/rcs"
	"github.com/blackchip-org/retro-cs/rcs/namco"
	"github.com/blackchip-org/retro-cs/system/c64"
	"github.com/blackchip-org/retro-cs/system/galaga"
	"github.com/blackchip-org/retro-cs/system/pacman"
	"github.com/veandco/go-sdl2/sdl"
)

var views = map[string]view{
	"c64:chars": view{
		system: "c64",
		roms:   c64.SystemROM,
		render: func(r *sdl.Renderer, d map[string][]byte) (rcs.TileSheet, error) {
			return c64.CharGen(r, d["chargen"])
		},
	},
	"c64:colors": view{
		system: "c64",
		render: func(r *sdl.Renderer, _ map[string][]byte) (rcs.TileSheet, error) {
			palettes := [][]color.RGBA{c64.Palette}
			return rcs.NewColorSheet(r, palettes)
		},
	},
	"galaga:sprites": view{
		system: "galaga",
		roms:   galaga.ROM["galaga"],
		render: func(r *sdl.Renderer, d map[string][]byte) (rcs.TileSheet, error) {
			return namco.NewTileSheet(r, d["sprite"],
				galaga.VideoConfig.SpriteLayout, namco.ViewerPalette)
		},
	},
	"galaga:tiles": view{
		system: "galaga",
		roms:   galaga.ROM["galaga"],
		render: func(r *sdl.Renderer, d map[string][]byte) (rcs.TileSheet, error) {
			return namco.NewTileSheet(r, d["tile"],
				galaga.VideoConfig.TileLayout, namco.ViewerPalette)
		},
	},
	"mspacman:sprites": view{
		system: "mspacman",
		roms:   pacman.ROM["mspacman"],
		render: func(r *sdl.Renderer, d map[string][]byte) (rcs.TileSheet, error) {
			return namco.NewTileSheet(r, d["sprite"],
				pacman.VideoConfig.SpriteLayout, namco.ViewerPalette)
		},
	},
	"mspacman:tiles": view{
		system: "mspacman",
		roms:   pacman.ROM["mspacman"],
		render: func(r *sdl.Renderer, d map[string][]byte) (rcs.TileSheet, error) {
			return namco.NewTileSheet(r, d["tile"],
				pacman.VideoConfig.TileLayout, namco.ViewerPalette)
		},
	},
	"pacman:colors": view{
		system: "pacman",
		roms:   pacman.ROM["pacman"],
		render: func(r *sdl.Renderer, d map[string][]byte) (rcs.TileSheet, error) {
			config := pacman.VideoConfig
			colors := namco.ColorTable(config, d["color"])
			return rcs.NewColorSheet(r, [][]color.RGBA{colors})
		},
	},
	"pacman:palettes": view{
		system: "pacman",
		roms:   pacman.ROM["pacman"],
		render: func(r *sdl.Renderer, d map[string][]byte) (rcs.TileSheet, error) {
			config := pacman.VideoConfig
			colors := namco.ColorTable(config, d["color"])
			palettes := namco.PaletteTable(config, d["palette"], colors)
			return rcs.NewColorSheet(r, palettes)
		},
	},
	"pacman:sprites": view{
		system: "pacman",
		roms:   pacman.ROM["pacman"],
		render: func(r *sdl.Renderer, d map[string][]byte) (rcs.TileSheet, error) {
			return namco.NewTileSheet(r, d["sprite"],
				pacman.VideoConfig.SpriteLayout, namco.ViewerPalette)
		},
	},
	"pacman:tiles": view{
		system: "pacman",
		roms:   pacman.ROM["pacman"],
		render: func(r *sdl.Renderer, d map[string][]byte) (rcs.TileSheet, error) {
			return namco.NewTileSheet(r, d["tile"],
				pacman.VideoConfig.TileLayout, namco.ViewerPalette)
		},
	},
}
