package fill

import (
	"embed"

	"github.com/spf13/cobra"
	"go.aponeill.com/fill/pkg/cmd/fill/create"
	"go.aponeill.com/fill/pkg/cmd/fill/merge"
	"go.aponeill.com/fill/pkg/cmd/fill/split"
	"go.aponeill.com/fill/pkg/cmd/fill/tile"
	"go.aponeill.com/fill/pkg/cmd/fill/world"
)

func NewCommand(fs embed.FS) *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "fill",
		Short: "fill is a tool for painting images with a color.",
	}
	cmd.AddCommand(tile.NewCommand())
	cmd.AddCommand(world.NewCommand())
	cmd.AddCommand(split.NewCommand())
	cmd.AddCommand(merge.NewCommand())
	cmd.AddCommand(create.NewCommand(fs))
	return cmd
}
