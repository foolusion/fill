package tile

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"go.aponeill.com/fill/pkg/fill"
	"github.com/spf13/cobra"
)

type flagpole struct {
	file     string
	outFile  string
	color    []int
	position []int
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "tile",
		Short: "Fills a tile with color.",
		Long:  "fills a tile with color from the given position.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.file, "file", "example.png", "The file to fill.")
	cmd.Flags().StringVar(&flags.outFile, "out", "", "optional output file")
	cmd.Flags().IntSliceVar(&flags.color, "color", []int{255, 0, 0, 255}, "The color to fill with. Alpha is optional")
	cmd.Flags().IntSliceVar(&flags.position, "position", []int{0, 0}, "position to fill the image from.")
	return cmd
}

var colorNames = []string{"red", "green", "blue", "alpha"}

func runE(flags *flagpole) error {
	if len(flags.color) < 3 || len(flags.color) > 4 {
		return fmt.Errorf("not enough values for color")
	}
	toColor := color.RGBA{}
	for i, v := range flags.color {
		if v < 0 || v > 255 {
			return fmt.Errorf("%s must be between 0 and 255", colorNames[i])
		}
		switch i {
		case 0:
			toColor.R = uint8(v)
		case 1:
			toColor.G = uint8(v)
		case 2:
			toColor.B = uint8(v)
		case 3:
			toColor.A = uint8(v)
		}
	}
	if len(flags.position) != 2 {
		return fmt.Errorf("position must contain 2 parts")
	}
	p := fill.Position{X: flags.position[0], Y: flags.position[1]}

	f, err := os.Open(flags.file)
	if err != nil {
		return err
	}
	i, imType, err := image.Decode(f)
	if err != nil {
		return err
	}
	f.Close()
	rgba := image.NewRGBA(image.Rect(0, 0, i.Bounds().Dx(), i.Bounds().Dy()))
	draw.Draw(rgba, rgba.Bounds(), i, i.Bounds().Min, draw.Src)
	fromColor := rgba.At(p.X, p.Y)
	if fill.EqualRGB(fromColor, toColor) {
		return nil
	}
	fill.TileFill(rgba, []fill.Position{p}, fromColor, toColor)
	outFile := flags.outFile
	if outFile == "" {
		outFile = flags.file
	}
	f, err = os.Create(outFile)
	if err != nil {
		return err
	}
	switch imType {
	case "png":
		err = png.Encode(f, rgba)
	case "jpeg":
		err = jpeg.Encode(f, rgba, nil)
	case "gif":
		err = gif.Encode(f, rgba, nil)
	}
	return err
}
