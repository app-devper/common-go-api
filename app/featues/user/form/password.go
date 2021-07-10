package form

type ChangePassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type SetPassword struct {
	Password string `json:"password" binding:"required"`
}
