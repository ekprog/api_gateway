package forms

type InstanceUpdateForm struct {
	Id       *int32  `json:"id" validate:"required"`
	Folder   *string `json:"folder" validate:""`
	Endpoint *string `json:"endpoint" validate:""`
	IsActive *bool   `json:"is_active" validate:""`
}
