package file

import (
	"errors"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"reapp/pkg/http/response"
)

type FileHandler struct {
	svc *Service
}

func NewFileHandler(svc *Service) *FileHandler { return &FileHandler{svc: svc} }

// Upload handles multipart form file upload
func (h *FileHandler) Upload(ctx *gin.Context) {
	f, err := ctx.FormFile("file")
	if err != nil {
		response.AsJSON(ctx, nil, errors.New("file is required"))
		return
	}

	keyOrURL, err := h.svc.SaveUpload(ctx.Request.Context(), "/files/", f)
	if err != nil {
		response.AsJSON(ctx, nil, err)
		return
	}

	response.AsJSON(ctx, gin.H{
		"path":      keyOrURL,
		"full_path": ctx.Request.Host + "/storages" + keyOrURL,
	}, nil)
}

// UploadImage handles multipart form image upload with optional WebP conversion
func (h *FileHandler) UploadImage(ctx *gin.Context) {
	f, err := ctx.FormFile("image")
	if err != nil {
		response.AsJSON(ctx, nil, errors.New("image is required"))
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
		response.AsJSON(ctx, nil, errors.New("unsupported image format"))
		return
	}

	// Open the file to check content type
	file, err := f.Open()
	if err != nil {
		response.AsJSON(ctx, nil, errors.New("failed to open image"))
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		response.AsJSON(ctx, nil, errors.New("failed to read image"))
		return
	}
	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		response.AsJSON(ctx, nil, errors.New("file is not an image"))
		return
	}

	// Reset file read pointer for service call
	file.Close()
	file, err = f.Open()
	if err != nil {
		response.AsJSON(ctx, nil, errors.New("failed to reopen image"))
		return
	}
	defer file.Close()

	webpFlag := ctx.Query("webp") == "true" || ctx.PostForm("webp") == "true"

	keyOrURL, err := h.svc.SaveImageUpload(ctx.Request.Context(), "/images/", f, webpFlag)
	if err != nil {
		response.AsJSON(ctx, nil, err)
		return
	}

	response.AsJSON(ctx, gin.H{
		"path":      keyOrURL,
		"full_path": ctx.Request.Host + "/storages" + keyOrURL,
	}, nil)
}

// ServeFile serves local files or redirects to S3 URL
func (h *FileHandler) ServeFile(ctx *gin.Context) {
	p := ctx.Param("path")
	p = strings.TrimPrefix(p, "/")
	// If storage can provide URL, use it
	if url, err := h.svc.URL(ctx.Request.Context(), p); err == nil && strings.HasPrefix(url, "http") {
		ctx.Redirect(http.StatusFound, url)
		return
	}

	data, ct, err := h.svc.Read(ctx.Request.Context(), p)
	if err != nil {
		response.AsJSON(ctx, nil, errors.New("not found"))
		return
	}

	// set content type and serve
	if ct == "" {
		ct = "application/octet-stream"
	}
	ctx.Data(http.StatusOK, ct, data)
}

// Delete remove file
func (h *FileHandler) Delete(ctx *gin.Context) {
	p := ctx.Param("path")
	p = strings.TrimPrefix(p, "/")
	if err := h.svc.Delete(ctx.Request.Context(), p); err != nil {
		response.AsJSON(ctx, nil, err)
		return
	}
	response.AsJSON(ctx, gin.H{"deleted": path.Base(p)}, nil)
}
