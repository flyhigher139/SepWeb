package file

import (
	"github.com/igevin/sepweb/pkg/context"
	"github.com/igevin/sepweb/pkg/handler"
	"net/http"
	"path/filepath"
)

type Downloader struct {
	Dir string
}

func (d *Downloader) Handle() handler.Handle {
	return func(ctx *context.Context) {
		req, _ := ctx.QueryValue("file").ToString()
		path := filepath.Join(d.Dir, filepath.Clean(req))
		fn := filepath.Base(path)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}
