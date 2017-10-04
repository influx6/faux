package actions

// CreateFile defines a struct for containing details of file created operation.
type CreateFile struct {
	FileName string `json:"file_name"`
	Dir      string `json:"dir"`
	RootDir  string `json:"root_dir"`
	Mode     int    `json:"mode"`
}

// MkDirectory defines a struct for containing details of dir created operation.
type MkDirectory struct {
	Dir     string `json:"dir"`
	RootDir string `json:"root_dir"`
	Mode    int    `json:"mode"`
}

// CreateUser defines a type for user creation details.
type CreateUser struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
