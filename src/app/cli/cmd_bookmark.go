package cli

import (
	"klog/app"
)

type Bookmark struct {
	File string `arg optional type:"existingfile" name:"file" help:".klg source file"`
}

func (args *Bookmark) Run(ctx app.Context) error {
	if args.File == "" {
		file, err := ctx.Bookmark()
		if err != nil {
			return err
		}
		ctx.Print("Current bookmark: " + file.Path + "\n")
		return nil
	}
	err := ctx.SetBookmark(args.File)
	if err != nil {
		return err
	}
	ctx.Print("Bookmarked file " + args.File + "\n")
	return nil
}
