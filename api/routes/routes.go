package routes

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"aletheia-api/internal/blockchain"
	"aletheia-api/internal/device"
	"aletheia-api/internal/fraud"
	"aletheia-api/internal/hash"
	"aletheia-api/internal/merkle"
	"aletheia-api/internal/models"
	"aletheia-api/internal/signature"
	"aletheia-api/internal/storage"
	urlproof "aletheia-api/internal/url"
	"aletheia-api/internal/verify"
	"aletheia-api/internal/video"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	repo       storage.MetadataRepository
	ipfs       storage.IPFSClient
	bigchain   storage.BigchainClient
	blockchain blockchain.Client
	verifySvc  *verify.Service
	videoSvc   *video.Service
	signer     signature.Service
	fraud      *fraud.Detector
	urlFetcher *urlproof.Service
	logger     *log.Logger
}

type verifyURLRequest struct {
	ImageURL         string `json:"imageUrl"`
	CaptureTimestamp string `json:"captureTimestamp"`
	Location         string `json:"location"`
	DeviceInfo       string `json:"deviceInfo"`
	DevicePublicKey  string `json:"devicePublicKey"`
	DeviceSignature  string `json:"deviceSignature"`
}

func Register(
	r *gin.Engine,
	repo storage.MetadataRepository,
	ipfs storage.IPFSClient,
	bigchain storage.BigchainClient,
	blockchainClient blockchain.Client,
	logger *log.Logger,
) {
	h := &Handler{
		repo:       repo,
		ipfs:       ipfs,
		bigchain:   bigchain,
		blockchain: blockchainClient,
		signer:     signature.NewSecp256k1Service(),
		verifySvc:  verify.NewService(repo, signature.NewSecp256k1Service()),
		videoSvc:   video.NewService(repo, ipfs, bigchain, blockchainClient),
		fraud:      fraud.NewDetector(10*time.Minute, 20),
		urlFetcher: urlproof.NewService(12*time.Second, 20<<20),
		logger:     logger,
	}

	v1 := r.Group("/api/v1")
	v1.POST("/images/upload", h.uploadImage)
	v1.POST("/verify", h.verifyImage)
	v1.POST("/verify/url", h.verifyImageURL)
	v1.POST("/video/upload", h.uploadVideo)
	v1.POST("/video/verify", h.verifyVideo)
}

func (h *Handler) uploadImage(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required as multipart/form-data field 'image'"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, 20<<20))
	if err != nil {
		h.logger.Printf("read upload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read uploaded image"})
		return
	}
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uploaded image is empty"})
		return
	}

	sha := hash.SHA256Hex(data)
	devicePublicKey := strings.TrimSpace(c.PostForm("devicePublicKey"))
	if devicePublicKey == "" {
		devicePublicKey = strings.TrimSpace(c.PostForm("device_public_key"))
	}
	deviceSignature := strings.TrimSpace(c.PostForm("deviceSignature"))
	if deviceSignature == "" {
		deviceSignature = strings.TrimSpace(c.PostForm("device_signature"))
	}
	captureTimestampStr := strings.TrimSpace(c.PostForm("captureTimestamp"))
	if captureTimestampStr == "" {
		captureTimestampStr = strings.TrimSpace(c.PostForm("capture_timestamp"))
	}
	parentHash := strings.TrimSpace(c.PostForm("parentHash"))
	if parentHash == "" {
		parentHash = strings.TrimSpace(c.PostForm("parent_hash"))
	}
	location := strings.TrimSpace(c.PostForm("location"))
	deviceInfo := strings.TrimSpace(c.PostForm("deviceInfo"))
	if deviceInfo == "" {
		deviceInfo = strings.TrimSpace(c.PostForm("device_info"))
	}
	uploaderID := strings.TrimSpace(c.PostForm("uploaderId"))
	if uploaderID == "" {
		uploaderID = strings.TrimSpace(c.PostForm("uploader_id"))
	}

	captureTimestamp, err := device.ParseCaptureTimestamp(captureTimestampStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid capture timestamp (RFC3339 required)"})
		return
	}
	if err := device.ValidateUploadProof(devicePublicKey, deviceSignature, captureTimestamp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sigValid, err := h.signer.VerifyHashHex(sha, deviceSignature, devicePublicKey)
	if err != nil || !sigValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device signature validation failed"})
		return
	}

	existingByHash, err := h.repo.FindBySHA256(c.Request.Context(), sha)
	if err != nil {
		h.logger.Printf("replay precheck: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to perform replay protection"})
		return
	}
	actor := uploaderID
	if actor == "" {
		actor = devicePublicKey
	}
	spam := h.fraud.RegisterUploadAndDetectSpam(actor, time.Now().UTC())
	fraudFlags, fraudWarnings := fraud.EvaluateUpload(existingByHash, captureTimestamp, location, deviceInfo, spam)

	if existingByHash != nil {
		c.JSON(http.StatusOK, models.UploadResponse{
			ProofID:           existingByHash.ProofID,
			SHA256:            existingByHash.SHA256,
			PHash:             existingByHash.PHash,
			MerkleRoot:        existingByHash.MerkleRoot,
			Signature:         existingByHash.Signature,
			PublicKey:         existingByHash.PublicKey,
			CaptureTimestamp:  existingByHash.CaptureTimestamp.Format(time.RFC3339),
			IPFSCID:           existingByHash.IPFSCID,
			TxHash:            existingByHash.TxHash,
			BlockNumber:       existingByHash.BlockNumber,
			Status:            "DUPLICATE",
			Flags:             fraudFlags,
			Warnings:          fraudWarnings,
		})
		return
	}

	ph, err := hash.ComputePHash(data)
	if err != nil {
		h.logger.Printf("compute phash: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "image decode failed"})
		return
	}

	merkleRoot, _, err := merkle.BuildMerkleRoot(data, 16)
	if err != nil {
		h.logger.Printf("build merkle: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to compute merkle root"})
		return
	}

	cid, err := h.ipfs.Upload(c.Request.Context(), fileHeader.Filename, data)
	if err != nil {
		h.logger.Printf("ipfs upload: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to upload to ipfs"})
		return
	}

	record := models.ProofRecord{
		ProofID:          uuid.NewString(),
		SHA256:           sha,
		PHash:            ph,
		MerkleRoot:       merkleRoot,
		Signature:        deviceSignature,
		PublicKey:        devicePublicKey,
		ParentHash:       parentHash,
		Location:         location,
		DeviceInfo:       deviceInfo,
		UploaderID:       uploaderID,
		CaptureTimestamp: captureTimestamp,
		IPFSCID:          cid,
		Timestamp:        time.Now().UTC(),
		Status:           "PROCESSING",
	}

	receipt, err := h.blockchain.AnchorProof(c.Request.Context(), record)
	if err != nil {
		h.logger.Printf("anchor polygon: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to anchor on polygon"})
		return
	}
	record.TxHash = receipt.TxHash
	record.BlockNumber = receipt.BlockNumber
	record.Timestamp = receipt.Timestamp
	record.Status = "ANCHORED"

	if err := h.repo.Save(c.Request.Context(), record); err != nil {
		h.logger.Printf("save metadata repo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store metadata"})
		return
	}

	if err := h.bigchain.SaveMetadata(c.Request.Context(), record); err != nil {
		h.logger.Printf("save metadata bigchaindb: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to persist metadata in bigchaindb"})
		return
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		ProofID:           record.ProofID,
		SHA256:            record.SHA256,
		PHash:             record.PHash,
		MerkleRoot:        record.MerkleRoot,
		Signature:         record.Signature,
		PublicKey:         record.PublicKey,
		CaptureTimestamp:  record.CaptureTimestamp.Format(time.RFC3339),
		IPFSCID:           record.IPFSCID,
		TxHash:            record.TxHash,
		BlockNumber:       record.BlockNumber,
		Status:            record.Status,
		Flags:             fraudFlags,
		Warnings:          fraudWarnings,
	})
}

func (h *Handler) verifyImage(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required as multipart/form-data field 'image'"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, 20<<20))
	if err != nil {
		h.logger.Printf("read verify image: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read uploaded image"})
		return
	}
	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uploaded image is empty"})
		return
	}

	verifyPublicKey := strings.TrimSpace(c.PostForm("devicePublicKey"))
	if verifyPublicKey == "" {
		verifyPublicKey = strings.TrimSpace(c.PostForm("device_public_key"))
	}
	verifySignature := strings.TrimSpace(c.PostForm("deviceSignature"))
	if verifySignature == "" {
		verifySignature = strings.TrimSpace(c.PostForm("device_signature"))
	}
	verifyCaptureTimestamp := time.Time{}
	captureTsRaw := strings.TrimSpace(c.PostForm("captureTimestamp"))
	if captureTsRaw == "" {
		captureTsRaw = strings.TrimSpace(c.PostForm("capture_timestamp"))
	}
	if captureTsRaw != "" {
		parsed, parseErr := device.ParseCaptureTimestamp(captureTsRaw)
		if parseErr == nil {
			verifyCaptureTimestamp = parsed
		}
	}

	verifyLocation := strings.TrimSpace(c.PostForm("location"))
	verifyDeviceInfo := strings.TrimSpace(c.PostForm("deviceInfo"))
	if verifyDeviceInfo == "" {
		verifyDeviceInfo = strings.TrimSpace(c.PostForm("device_info"))
	}

	resp, err := h.verifySvc.VerifyImage(c.Request.Context(), data, models.VerifyOptions{
		SourceType:       "UPLOAD",
		CaptureTimestamp: verifyCaptureTimestamp,
		Location:         verifyLocation,
		DeviceInfo:       verifyDeviceInfo,
		DevicePublicKey:  verifyPublicKey,
		DeviceSignature:  verifySignature,
	})
	if err != nil {
		h.logger.Printf("verify image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "verification failed"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) verifyImageURL(c *gin.Context) {
	var req verifyURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if strings.TrimSpace(req.ImageURL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "imageUrl is required"})
		return
	}

	imageBytes, err := h.urlFetcher.FetchImageBytes(c.Request.Context(), req.ImageURL)
	if err != nil {
		h.logger.Printf("verify url fetch image: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	captureTimestamp := time.Time{}
	if strings.TrimSpace(req.CaptureTimestamp) != "" {
		parsed, parseErr := device.ParseCaptureTimestamp(req.CaptureTimestamp)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid captureTimestamp; expected RFC3339"})
			return
		}
		captureTimestamp = parsed
	}

	resp, err := h.verifySvc.VerifyImage(c.Request.Context(), imageBytes, models.VerifyOptions{
		SourceType:       "URL",
		CaptureTimestamp: captureTimestamp,
		Location:         req.Location,
		DeviceInfo:       req.DeviceInfo,
		DevicePublicKey:  req.DevicePublicKey,
		DeviceSignature:  req.DeviceSignature,
	})
	if err != nil {
		h.logger.Printf("verify url image: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "verification failed"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) uploadVideo(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video is required as multipart/form-data field 'video'"})
		return
	}
	defer file.Close()

	resp, err := h.videoSvc.Upload(c.Request.Context(), file, fileHeader.Filename)
	if err != nil {
		h.logger.Printf("upload video: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) verifyVideo(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video is required as multipart/form-data field 'video'"})
		return
	}
	defer file.Close()

	resp, err := h.videoSvc.Verify(c.Request.Context(), file, fileHeader.Filename)
	if err != nil {
		h.logger.Printf("verify video: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
