package main

import (
	"bytes"
	"encoding/json"
	"fastimage/dto"
	"fastimage/storage"
	"fastimage/utils"
	"strconv"

	"log"

	"github.com/chai2010/webp"
	"github.com/fasthttp/router"
	"github.com/nfnt/resize"
	"github.com/valyala/fasthttp"
)

func createImage(ctx *fasthttp.RequestCtx) {
	// Content-Type이 application/octet-stream 인가?
	contentType := string(ctx.Request.Header.ContentType())
	if contentType != "application/octet-stream" {
		ctx.Error("Invalid Content-Type", fasthttp.StatusUnsupportedMediaType)
		return
	}

	postBody := ctx.PostBody()
	uuid, err := storage.SaveImage(postBody)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	// uuid 반환
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(&dto.CreateImageResponse{
		UUID: uuid,
	})
}

func findImage(ctx *fasthttp.RequestCtx) {
	uuid := ctx.UserValue("uuid").(string)
	width := utils.CleanUp(ctx.QueryArgs().Peek("width"))
	height := utils.CleanUp(ctx.QueryArgs().Peek("height"))

	image, err := storage.FindImage(uuid)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	if !utils.IsEmpty(width, height) {
		bounds := image.Bounds()
		changeWidth := bounds.Max.X
		changeHeight := bounds.Max.Y

		if width != "" {
			w, _ := strconv.Atoi(width)
			changeWidth = w
		}

		if height != "" {
			h, _ := strconv.Atoi(height)
			changeHeight = h
		}

		newImage := resize.Resize(uint(changeWidth), uint(changeHeight), image, resize.Lanczos3)
		image = newImage
	}

	var buf bytes.Buffer
	webp.Encode(&buf, image, &webp.Options{Quality: 100})
	ctx.SetContentType("image/webp")
	ctx.Write(buf.Bytes())

}

func main() {
	defer storage.DB().Close()

	r := router.New()
	r.POST("/api/v1/images", createImage)
	r.GET("/api/v1/images/{uuid}", findImage)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
