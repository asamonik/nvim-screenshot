package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.design/x/clipboard"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/neovim/go-client/nvim/plugin"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
)

type EvalResponse struct {
	BufferContent []string `msgpack:",array"`
	Syntax        []string `msgpack:",array"`
	FgColor       string
	BgColor       string
}

func parseColor(color_string string, default_color int) (color.RGBA, error) {

	color_rgba := color.RGBA{0x33, 0x33, 0x33, 0xff}

	_, err := fmt.Sscanf(
		color_string,
		"#%02x%02x%02x",
		&color_rgba.R,
		&color_rgba.G,
		&color_rgba.B)

	if err != nil {

		log.Println(err)

		color_rgba = color.RGBA{
			uint8(default_color >> 24 & 0xff),
			uint8(default_color >> 16 & 0xff),
			uint8(default_color >> 8 & 0xff),
			uint8(default_color & 0xff),
		}
	}

	return color_rgba, err
}

func generateImage(code []string, syntax []string, fontSize float64, fgcol string, bgcol string) (*image.RGBA, error) {

	loadedFont, err := truetype.Parse(gomono.TTF)
	if err != nil {
		return nil, err
	}

	fgColor, _ := parseColor(fgcol, 0xffffffff)
	bgColor, _ := parseColor(bgcol, 0x00000000)

	fg := image.NewUniform(fgColor)
	bg := image.NewUniform(bgColor)
	lineNo := image.NewUniform(fgColor)

	imageheight := int(fontSize*1.5)*len(code) + int(0.5*fontSize)

	rgba := image.NewRGBA(image.Rect(0, 0, 2500, imageheight))

	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{0, 0}, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(loadedFont)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetHinting(font.HintingNone)
	c.SetClip(rgba.Bounds())

	for i, s := range code {

		pt := freetype.Pt(10, int(fontSize*1.5)*(i+1))

		// line numbers
		c.SetSrc(lineNo)
		c.DrawString(strconv.Itoa(i+1), pt)

		if len(syntax[i]) == 0 {
			continue
		}

		// line
		pt.X = c.PointToFixed(100)

		colors := strings.Split(syntax[i], ";")

		var line []string
		end, counter := 0, 0
		for chridx, character := range s {

			if chridx >= end && counter < len(colors) {

				line = strings.Split(colors[counter], ":")

				counter++
				if counter < len(colors) {
					end, _ = strconv.Atoi(strings.Split(colors[counter], ":")[0])
				}
			}

			fgColor, _ = parseColor(
				line[1],
				0xffffffff)

			fg = image.NewUniform(fgColor)
			c.SetSrc(fg)

			text := strings.ReplaceAll(string(character), "\t", "    ")
			pt, err = c.DrawString(text, pt)

			if err != nil {
				log.Println(err)
			}
		}

	}

	return rgba, nil
}

func screenshot(args []string, start_end []int, eval *EvalResponse) {

	log.Println("Taking screenshot")

	start := start_end[0] - 1
	end := start_end[1]

	code := eval.BufferContent[start:end]
	syntax := eval.Syntax[start:end]

	image, err := generateImage(
		code,
		syntax,
		30,
		eval.FgColor, eval.BgColor)

	if err != nil {
		log.Println(err)
	}

	ctc := len(args) == 0
	buf := new(bytes.Buffer)
	err = png.Encode(buf, image)

	if err != nil {
		log.Println(err)
	}

	if ctc {
		log.Println("Copying image to clipboard")
		clipboard.Write(clipboard.FmtImage, buf.Bytes())
		log.Println("Copied image to clipboard")
	} else {
		if _, err := os.Stat(args[0]); !os.IsNotExist(err) {

			timestamp := time.Now().Unix()
			filename := args[0] + "/vimscreenshot_" + strconv.FormatInt(timestamp, 10) + ".png"

			log.Println("Saving to " + filename)

			f, err := os.Create(filename)
			if err != nil {
				log.Println(err)
			}
			f.Write(buf.Bytes())
			f.Close()

			log.Println("Saved image as " + filename)
		} else {
			log.Println("invalid path")
		}
	}
}

func main() {

	err := clipboard.Init()
	if err != nil {
		log.Println(err)
	}

	plugin.Main(func(p *plugin.Plugin) error {
		p.HandleCommand(&plugin.CommandOptions{
			Name:  "Screenshot",
			NArgs: "?",
			Range: "",
			Eval:  `getline(1, '$'),GetSyntax(),synIDattr(hlID("Normal"), "fg"),synIDattr(hlID("Normal"),"bg")]`},
			screenshot)

		return nil
	})
}
