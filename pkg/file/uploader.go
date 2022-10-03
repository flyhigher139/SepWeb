package file

import (
	"github.com/igevin/sepweb/pkg/context"
	"github.com/igevin/sepweb/pkg/handler"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type Uploader struct {
	FileField string
	// 计算目标路径
	DstPathFunc func(fh *multipart.FileHeader) string
}

func (u *Uploader) Handle() handler.Handle {
	return func(ctx *context.Context) {
		src, srcHeader, err := ctx.Req.FormFile(u.FileField)
		if err != nil {
			u.uploadFailedForNoData(ctx, err)
			return
		}
		defer src.Close()
		dst, err := os.OpenFile(u.DstPathFunc(srcHeader),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			u.failToUpload(ctx, err)
			return
		}
		defer dst.Close()
		_, err = io.CopyBuffer(dst, src, nil)
		if err != nil {
			u.failToUpload(ctx, err)
			return
		}
		ctx.RespData = []byte("上传成功")
	}
}

func (u *Uploader) HandleFunc(ctx *context.Context) {
	src, srcHeader, err := ctx.Req.FormFile(u.FileField)
	if err != nil {
		u.failToUpload(ctx, err)
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(u.DstPathFunc(srcHeader),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		u.failToUpload(ctx, err)
		return
	}
	defer dst.Close()
	_, err = io.CopyBuffer(dst, src, nil)
	if err != nil {
		u.failToUpload(ctx, err)
		return
	}
	ctx.RespData = []byte("上传成功")
}

func (u *Uploader) uploadFailedForNoData(ctx *context.Context, err error) {
	ctx.RespStatusCode = http.StatusBadRequest
	ctx.RespData = []byte("上传失败，未找到数据")
	log.Fatalln(err)
}

func (u *Uploader) failToUpload(ctx *context.Context, err error) {
	ctx.RespStatusCode = http.StatusInternalServerError
	ctx.RespData = []byte("上传失败")
	log.Fatalln(err)
}
