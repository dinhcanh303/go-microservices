package request

type UpdateAttachmentRequest struct {
	AttachableType string `json:"attachableType"`
	AttachableID   string `json:"attachableId"`
	EntityUploadID string `json:"entityUploadId"`
}
type UpdateAttachmentsByIdsRequest struct {
	AttachmentIds  []string `json:"attachmentIds"`
	AttachableType string   `json:"attachableType"`
	AttachableID   string   `json:"attachableId"`
	EntityUploadID string   `json:"entityUploadId"`
}
type DeleteAttachmentsByIdsRequest struct {
	AttachmentIds []string `json:"attachmentIds"`
}
