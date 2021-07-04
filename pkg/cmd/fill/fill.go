package fill

import (
	"aponeill.com/fill/pkg/cmd/fill/merge"
	"aponeill.com/fill/pkg/cmd/fill/split"
	"aponeill.com/fill/pkg/cmd/fill/tile"
	"aponeill.com/fill/pkg/cmd/fill/world"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "fill",
		Short: "fill is a tool for painting images with a color.",
	}
	cmd.AddCommand(tile.NewCommand())
	cmd.AddCommand(world.NewCommand())
	cmd.AddCommand(split.NewCommand())
	cmd.AddCommand(merge.NewCommand())
	return cmd
}
