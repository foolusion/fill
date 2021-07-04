package split

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
	file   string
	out    string
	width  int
	height int
}

func NewCommand() *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "split",
		Short: "Splits an image into smaller images.",
		Long:  "Splits an image into n x m smaller images.",
		RunE: func(*cobra.Command, []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.file, "file", "", "The file you want to split.")
	cmd.Flags().StringVar(&flags.out, "out", "", "The output directory.")
	cmd.Flags().IntVar(&flags.width, "c", 10, "How many columns of images to produce.")
	cmd.Flags().IntVar(&flags.height, "r", 10, "How many rows of images to produce.")
	return cmd
}

func runE(flags *flagpole) error {
	f, err := os.Open(flags.file)
	if err != nil {
		return err
	}
	i, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	f.Close()
	tileWidth := i.Bounds().Dx() / flags.width
	tileHeight := i.Bounds().Dy() / flags.height

	for x := 0; x < flags.width; x++ {
		for y := 0; y < flags.width; y++ {
			tx := x * tileWidth
			ty := y * tileHeight
			ii := image.NewNRGBA(image.Rect(tx, ty, tx+tileWidth, ty+tileHeight))
			draw.Draw(ii, ii.Bounds(), i, image.Pt(tx, ty), draw.Src)
			f, err = os.Create(filepath.Join(flags.out, fmt.Sprintf("tile_%d_%d.png", x, y)))
			if err != nil {
				return err
			}
			err = png.Encode(f, ii)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
