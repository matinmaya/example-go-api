package filestorage

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"

	"reapp/pkg/filesystem"
	"reapp/pkg/http/response"
)

type FileHandler struct {
	svc *Service
}

func NewFileHandler(svc *Service) *FileHandler { return &FileHandler{svc: svc} }

func (h *FileHandler) Upload(ctx *gin.Context) {
	f, err := ctx.FormFile("file")
	if err != nil {
		response.JSON(ctx, nil, errors.New("file is required"))
		return
	}

	keyOrURL, err := h.svc.SaveUpload(ctx.Request.Context(), "/files/", f)
	if err != nil {
		response.JSON(ctx, nil, err)
		return
	}

	response.JSON(ctx, gin.H{
		"path":      keyOrURL,
		"full_path": ctx.Request.Host + "/" + filesystem.TrimPath(filesystem.PrefixRoutePath()) + keyOrURL,
	}, nil)
}

func (h *FileHandler) UploadImage(ctx *gin.Context) {
	f, err := ctx.FormFile("image")
	if err != nil {
		response.JSON(ctx, nil, errors.New("image is required"))
		return
	}

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	ext := strings.ToLower(path.Ext(f.Filename))
	if !allowedExts[ext] {
		response.JSON(ctx, nil, errors.New("unsupported image format"))
		return
	}

	file, err := f.Open()
	if err != nil {
		response.JSON(ctx, nil, errors.New("failed to open image"))
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		response.JSON(ctx, nil, errors.New("failed to read image"))
		return
	}
	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		response.JSON(ctx, nil, errors.New("file is not an image"))
		return
	}

	// Reset file read pointer for service call
	file.Close()
	file, err = f.Open()
	if err != nil {
		response.JSON(ctx, nil, errors.New("failed to reopen image"))
		return
	}
	defer file.Close()

	webpFlag := ctx.Query("webp") == "true" || ctx.PostForm("webp") == "true"

	keyOrURL, err := h.svc.SaveImageUpload(ctx.Request.Context(), "/images/", f, webpFlag)
	if err != nil {
		response.JSON(ctx, nil, err)
		return
	}

	response.JSON(ctx, gin.H{
		"path":      keyOrURL,
		"full_path": ctx.Request.Host + "/" + filesystem.TrimPath(filesystem.PrefixRoutePath()) + keyOrURL,
	}, nil)
}

func (h *FileHandler) ServeFile(ctx *gin.Context) {
	p := strings.TrimPrefix(ctx.Param("path"), "/")

	if url, err := h.svc.URL(ctx.Request.Context(), p); err == nil && strings.HasPrefix(url, "http") {
		ctx.Redirect(http.StatusFound, url)
		return
	}

	data, ct, err := h.svc.Read(ctx.Request.Context(), p)
	if err != nil {
		log.Printf("%s", err.Error())
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	wStr := ctx.Query("w")
	hStr := ctx.Query("h")
	fillFlag := ctx.Query("fill") == "true"

	var numberOnly = regexp.MustCompile(`^\d+$`)
	if wStr != "" && !numberOnly.MatchString(wStr) {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid width parameter"))
		return
	}

	if hStr != "" && !numberOnly.MatchString(hStr) {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid height parameter"))
		return
	}

	if wStr != "" {
		wVal, _ := strconv.Atoi(wStr)
		if wVal < 16 {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("width must be at least 16"))
			return
		}
	}
	if hStr != "" {
		hVal, _ := strconv.Atoi(hStr)
		if hVal < 16 {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("height must be at least 16"))
			return
		}
	}

	if (wStr != "" || hStr != "") && strings.HasPrefix(ct, "image/") {
		fillFlag = fillFlag && wStr != "" && hStr != ""
		ext := filepath.Ext(p)
		name := strings.TrimSuffix(filepath.Base(p), ext)
		fillSuffix := "0"
		if fillFlag {
			fillSuffix = "1"
		}

		cachePath := ""
		if filesystem.IsValidStoragePath(CacheRootPath()) {
			cachePath = filepath.Join(CacheRootPath(), filepath.Dir(p), fmt.Sprintf("%s_w%s_h%s_f%s%s", name, wStr, hStr, fillSuffix, ext))
			if cachedData, err := os.ReadFile(cachePath); err == nil {
				ctx.Data(http.StatusOK, ct, cachedData)
				return
			}
		}

		img, _, err := image.Decode(bytes.NewReader(data))
		if err == nil {
			var width, height int
			if wStr != "" {
				width, _ = strconv.Atoi(wStr)
			}
			if hStr != "" {
				height, _ = strconv.Atoi(hStr)
			}

			var resized image.Image
			if fillFlag {
				resized = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
			} else {
				resized = imaging.Resize(img, width, height, imaging.Lanczos)
			}

			var buf bytes.Buffer
			switch ct {
			case "image/jpeg":
				if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 85}); err != nil {
					log.Printf("image converting %s: %s", ct, err.Error())
				}
			case "image/png":
				if err := png.Encode(&buf, resized); err != nil {
					log.Printf("image converting %s: %s", ct, err.Error())
				}
			case "image/webp":
				if err := webp.Encode(&buf, resized, &webp.Options{Lossless: false, Quality: 85}); err != nil {
					log.Printf("image converting %s: %s", ct, err.Error())
				}
			default:
				if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 85}); err != nil {
					log.Printf("image converting: %s", err.Error())
				}
				ct = "image/jpeg"
			}
			data = buf.Bytes()
		}

		if filesystem.IsValidStoragePath(CacheRootPath()) {
			if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
				log.Printf("write cache image: %s", err.Error())
			}
			if err := os.WriteFile(cachePath, data, 0644); err != nil {
				log.Printf("write cache image: %s", err.Error())
			}
		}
	}

	if ct == "" {
		ct = "application/octet-stream"
	}
	ctx.Data(http.StatusOK, ct, data)
}

func (h *FileHandler) Delete(ctx *gin.Context) {
	p := ctx.Param("path")
	p = strings.TrimPrefix(p, "/")
	if err := h.svc.Delete(ctx.Request.Context(), p); err != nil {
		response.JSON(ctx, nil, err)
		return
	}
	response.JSON(ctx, gin.H{"deleted": path.Base(p)}, nil)
}
