package helper

import (
	"crunchgarage/restaurant-food-delivery/config"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"context"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
)

type File struct {
	File multipart.File `json:"file,omitempty" validate:"required"`
}

var USE_CLOUDINARY = false

func SingleImageUpload(c *gin.Context, avatar, bucket_storage_folder, subDir string) (string, error) {
	if USE_CLOUDINARY {
		return SingleImageUploadToCloudinary(c, avatar, bucket_storage_folder, subDir)
	} else {
		return SingleImageUploadToLocal(c, avatar, bucket_storage_folder, subDir)
	}
}

func SingleImageUploadToCloudinary(c *gin.Context, avatar, bucket_storage_folder, subDir string) (string, error) {

	file, err := c.FormFile(avatar)

	if err != nil {
		SendErrorPayload(c, http.StatusBadRequest, err)
		return "", err
	}

	folder := filepath.Join(bucket_storage_folder, subDir)

	result, err := CloudinaryUpload(file, folder, file.Filename)

	if err != nil {
		SendErrorPayload(c, http.StatusBadRequest, err)
		return "", err
	}

	return result, err

}

func CloudinaryUpload(media_url interface{}, bucket_storage_folder string, file_name string) (string, error) {

	cld, _ := cloudinary.NewFromParams(config.EnvCloudName(),
		config.EnvCloudAPIKey(),
		config.EnvCloudAPISecret())

	if err != nil {
		return "", err
	}

	var ctx = context.Background()

	uploadResult, err := cld.Upload.Upload(
		ctx,
		media_url,
		uploader.UploadParams{Folder: bucket_storage_folder, PublicID: file_name})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil

}

func SingleImageUploadToLocal(c *gin.Context, avatar string, bucket_storage_folder string, subDir string) (string, error) {
	file, err := c.FormFile(avatar)

	if err != nil {
		return "", err
	}

	filename := filepath.Base(file.Filename)
	savePath := filepath.Join("./uploads/", bucket_storage_folder, subDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
		return "", err
	}

	return savePath, nil
}
