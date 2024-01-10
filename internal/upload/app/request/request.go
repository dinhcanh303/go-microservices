package request

type UpdateAttachmentRequest struct {
	AttachableType string `json:"attachableType"`
	AttachableID   string `json:"attachableId"`
}
type UpdateAttachmentsByIdsRequest struct {
	AttachmentIds  []string `json:"attachmentIds"`
	AttachableType string   `json:"attachableType"`
	AttachableID   string   `json:"attachableId"`
}
type DeleteAttachmentsByIdsRequest struct {
	AttachmentIds []string `json:"attachmentIds"`
}
