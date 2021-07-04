package world

import (
	"fmt"
	"image/color"

	"github.com/spf13/cobra"
	"go.aponeill.com/fill/pkg/fill"
)

type flagpole struct {
	path         string
	color        []int
	filePosition []int
	tilePosition []int
	numWorkers   int
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "world",
		Short: "Fills all tiles in a world with color.",
		Long:  "Fills a tile with color from the given position and continues filling it's neighboring tiles.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.path, "path", ".", "The path to tile_{x}_{y}.png files to fill.")
	cmd.Flags().IntSliceVar(&flags.color, "color", []int{255, 0, 0, 255}, "The color to fill with. Alpha is optional")
	cmd.Flags().IntSliceVar(&flags.tilePosition, "tp", []int{0, 0}, "Tile position to fill the image from.")
	cmd.Flags().IntSliceVar(&flags.filePosition, "fp", []int{0, 0}, "File position to fill the world from.")
	cmd.Flags().IntVar(&flags.numWorkers, "workers", 10, "The number of concurrent workers to spawn.")
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
	if len(flags.filePosition) != 2 {
		return fmt.Errorf("file position must contain 2 parts")
	}
	fp := fill.Position{X: flags.filePosition[0], Y: flags.filePosition[1]}
	if len(flags.tilePosition) != 2 {
		return fmt.Errorf("tile position must contain 2 parts")
	}
	tp := fill.Position{X: flags.tilePosition[0], Y: flags.tilePosition[1]}
	return fill.WorldFill(flags.path, fp, tp, toColor, flags.numWorkers)
}
