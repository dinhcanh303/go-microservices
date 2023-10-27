package request

type UpdateAttachmentRequest struct {
	AttachableType string `json:"attachable_type"`
	AttachableID   string `json:"attachable_id"`
}
type UpdateAttachmentsByIdsRequest struct {
	AttachmentIds  []string `json:"attachment_ids"`
	AttachableType string   `json:"attachable_type"`
	AttachableID   string   `json:"attachable_id"`
}
type DeleteAttachmentsByIdsRequest struct {
	AttachmentIds []string `json:"attachment_ids"`
}
