package utils

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"
)

func CleanUp(arg any) (r string) {
	var m string

	switch v := arg.(type) {
	case []byte:
		m = string(v)
	}

	r = strings.ReplaceAll(m, `"`, "")
	return r
}

func IsEmpty(args ...string) bool {
	for _, arg := range args {
		if len(arg) != 0 {
			return false
		}
	}

	return true
}

func ResponseImage(ctx *fasthttp.RequestCtx, buf *bytes.Buffer) {
	ctx.Response.Header.Set("Cache-Control", "max-age=5")
	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(buf.Bytes())))
	ctx.SetContentType("image/webp")
	ctx.Write(buf.Bytes())
}
