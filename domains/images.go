package domains

type SignedUrlResponse struct {
	SignedUrl string `json:"signed_url"`
}

type ImageUploadedRequest struct {
	Bucket string `json:"bucket"`
	Path   string `json:"path"`
}
