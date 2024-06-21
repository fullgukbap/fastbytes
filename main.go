package main

import (
	"bytes"
	"encoding/json"
	"fastimage/dto"
	"fastimage/storage"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"

	"log"
	"os"
	"slices"

	"github.com/chai2010/webp"
	"github.com/fasthttp/router"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
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

	// 빈 파일인지 검사하기
	emptySlice := make([]byte, 0)
	if slices.Equal(postBody, emptySlice) {
		ctx.Error("Empty Body", fasthttp.StatusBadRequest)
		return
	}

	// 이미지 파일인가?
	if !filetype.IsImage(postBody) {
		ctx.Error("It's not an image file", fasthttp.StatusBadRequest)
		return
	}

	// uuid를 생성하여 새로운 이미지 생성
	uuid := uuid.New().String()

	// 파일 유형 검사 및 파일 확장자 추출
	kind, err := filetype.Match(postBody)
	if err != nil || kind == filetype.Unknown {
		ctx.Error("Unsupported file type", fasthttp.StatusUnsupportedMediaType)
		return
	}

	// database에 확장자 저장
	err = storage.DB().SaveExtendsion(uuid, kind.Extension)
	if err != nil {
		ctx.Error("Not save extension", fasthttp.StatusInternalServerError)
		return
	}

	err = os.WriteFile(fmt.Sprintf("./storage/images/%s.%s", uuid, kind.Extension), postBody, 0644)
	if err != nil {
		ctx.Error("Unable to save file", fasthttp.StatusInternalServerError)
		return
	}

	// uuid 반환
	json.NewEncoder(ctx).Encode(&dto.CreateImageResponse{
		UUID: uuid,
	})
}

// 뭔가 이상한게 자꾸 끼어 있음 query paras에 그래서 깨끗하게 청소하는 작업을 ㅁ나ㅡㄹ자

func findImage(ctx *fasthttp.RequestCtx) {
	uuid := ctx.UserValue("uuid").(string)
	w := string(ctx.QueryArgs().Peek("webp"))
	width := string(ctx.QueryArgs().Peek("width"))
	height := string(ctx.QueryArgs().Peek("height"))

	w = strings.ReplaceAll(w, `"`, "")
	width = strings.ReplaceAll(width, `"`, "")
	height = strings.ReplaceAll(height, `"`, "")

	// 이미지 가져오고
	ext := storage.DB().FindExtension(uuid)
	file, err := os.Open(fmt.Sprintf("./storage/images/%s.%s", uuid, ext)) // if err != nil {
	if err != nil {
		ctx.Error("Failed to read image", fasthttp.StatusInternalServerError)
		return
	}
	defer file.Close()

	var myImg image.Image

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	myImg = img

	if width != "" && height != "" {
		bounds := img.Bounds()
		changeWidth := bounds.Max.X
		changeHeight := bounds.Max.Y

		// TODO: height, width 최대값 정하기
		if width != "" {
			w, err := strconv.Atoi(width)
			if err != nil {
				panic(err)
			}
			changeWidth = w
		}

		if height != "" {
			h, _ := strconv.Atoi(height)
			changeHeight = h
		}

		newImage := resize.Resize(uint(changeWidth), uint(changeHeight), img, resize.Lanczos3)
		myImg = newImage
	}

	if w == "true" {
		var buf bytes.Buffer
		options := &webp.Options{Quality: 100}
		err = webp.Encode(&buf, myImg, options)
		if err != nil {
			panic(err)
		}
		ctx.SetContentType("image/webp")
		ctx.Write(buf.Bytes())
		return
	}

	var buf bytes.Buffer
	switch ext {
	case "jpg":
		fallthrough
	case "jpeg":
		ctx.SetContentType("image/jpeg")
		err := jpeg.Encode(&buf, myImg, &jpeg.Options{Quality: 100})
		if err != nil {
			ctx.Error("Failed to encode image", fasthttp.StatusInternalServerError)
			return
		}
	case "png":
		ctx.SetContentType("image/png")
		err := png.Encode(&buf, myImg)
		if err != nil {
			ctx.Error("Failed to encode image", fasthttp.StatusInternalServerError)
			return
		}
	}

	ctx.Write(buf.Bytes())
}

func main() {
	defer storage.DB().Close()

	r := router.New()
	r.POST("/api/v1/images", createImage)
	r.GET("/api/v1/images/{uuid}", findImage)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
