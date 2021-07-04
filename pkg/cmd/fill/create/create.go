package create

import (
	"embed"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type flagpole struct {
	dir string
}

func NewCommand(fs embed.FS) *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "create",
		Short: "Create static files for dsa challenge.",
		Long:  "Creates the example.png, part_1/dsa_challenge.png and part_2/tile*.png files.",
		RunE: func(*cobra.Command, []string) error {
			return runE(flags, fs)
		},
	}
	cmd.Flags().StringVar(&flags.dir, "dir", ".", "The directory to output files.")
	return cmd
}

func runE(flags *flagpole, fs embed.FS) error {
	part1Dir := filepath.Join(flags.dir, "part_1")
	err := os.MkdirAll(part1Dir, 0755)
	if err != nil {
		return err
	}
	part2Dir := filepath.Join(flags.dir, "part_2")
	err = os.MkdirAll(part2Dir, 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(flags.dir, "example.png"))
	if err != nil {
		return err
	}
	ef, err := fs.Open("res/example.png")
	if err != nil {
		return err
	}
	_, err = io.Copy(f, ef)
	if err != nil {
		return err
	}
	ef.Close()
	f.Close()

	f, err = os.Create(filepath.Join(part1Dir, "dsa_challenge.png"))
	if err != nil {
		return err
	}
	p1f, err := fs.Open("res/dsa_challenge.png")
	if err != nil {
		return err
	}
	_, err = io.Copy(f, p1f)
	if err != nil {
		return err
	}
	p1f.Close()
	f.Close()

	f, err = os.Create(filepath.Join(part2Dir, "dsa_challenge_2.png"))
	if err != nil {
		return err
	}
	p2f, err := fs.Open("res/dsa_challenge_2.png")
	if err != nil {
		return err
	}
	_, err = io.Copy(f, p2f)
	if err != nil {
		return err
	}
	p2f.Close()
	f.Close()

	return nil
}
