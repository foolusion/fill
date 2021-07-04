package merge

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type flagpole struct {
	path   string
	out    string
	width  int
	height int
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "merge",
		Short: "Merges many images into one image.",
		Long:  "Merges n x m smaller images into one image.",
		RunE: func(*cobra.Command, []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.path, "path", "", "The path you want to merge.")
	cmd.Flags().StringVar(&flags.out, "out", "", "The output file.")
	cmd.Flags().IntVar(&flags.width, "c", 10, "How many columns of images to merge.")
	cmd.Flags().IntVar(&flags.height, "r", 10, "How many rows of images to merge.")
	return cmd
}

func runE(flags *flagpole) error {
	f, err := os.Open(filepath.Join(flags.path, "tile_0_0.png"))
	if err != nil {
		return err
	}
	i, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	f.Close()
	imgWidth := i.Bounds().Dx() * flags.width
	imgHeight := i.Bounds().Dy() * flags.height
	tileWidth := i.Bounds().Dx()
	tileHeight := i.Bounds().Dy()

	out := image.NewNRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for x := 0; x < flags.width; x++ {
		for y := 0; y < flags.width; y++ {
			file := filepath.Join(flags.path, fmt.Sprintf("tile_%d_%d.png", x, y))
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			i, _, err := image.Decode(f)
			if err != nil {
				return err
			}
			f.Close()
			tx := x * tileWidth
			ty := y * tileHeight
			draw.Draw(out, image.Rect(tx, ty, tx+tileWidth, ty+tileHeight), i, image.Point{}, draw.Src)
		}
	}
	f, err = os.Create(flags.out)
	if err != nil {
		return err
	}
	err = png.Encode(f, out)
	if err != nil {
		return err
	}
	return f.Close()
}
