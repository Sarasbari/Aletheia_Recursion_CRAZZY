package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/user/aletheia-api/docs"
	"github.com/user/aletheia-api/internal/verify"
)

type Handler struct {
	verifyService *verify.Service
}

func NewHandler(vs *verify.Service) *Handler {
	return &Handler{verifyService: vs}
}

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/images/upload", h.UploadImage)
		v1.POST("/verify", h.VerifyImage)
	}
}

// UploadImage godoc
// @Summary      Upload an image and generate cryptographic proof
// @Description  Accepts an image via multipart/form-data, generates SHA-256, pHash, and Merkle Root, then stores metadata on IPFS, BigchainDB, and Polygon.
// @Tags         images
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file  true  "Image to upload"
// @Success      200    {object}  models.ImageProof
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /images/upload [post]
func (h *Handler) UploadImage(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}
	defer file.Close()

	proof, err := h.verifyService.ProcessUpload(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, proof)
}

// VerifyImage godoc
// @Summary      Verify an image against existing proofs
// @Description  Recomputes hashes for the provided image and compares them with stored data to determine authenticity and trust score.
// @Tags         verification
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file  true  "Image to verify"
// @Success      200    {object}  models.VerificationResult
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /verify [post]
func (h *Handler) VerifyImage(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}
	defer file.Close()

	result, err := h.verifyService.VerifyImage(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
