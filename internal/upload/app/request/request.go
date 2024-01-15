package request

type UpdateAttachmentRequest struct {
	AttachableType string `json:"attachableType"`
	AttachableID   string `json:"attachableId"`
	EntityUpload   string `json:"entityUpload"`
}
type UpdateAttachmentsByIdsRequest struct {
	AttachmentIds  []string `json:"attachmentIds"`
	AttachableType string   `json:"attachableType"`
	AttachableID   string   `json:"attachableId"`
	EntityUpload   string   `json:"entityUpload"`
}
type DeleteAttachmentsByIdsRequest struct {
	AttachmentIds []string `json:"attachmentIds"`
}
