package replies

type UpdateUser struct {
	ID     int               `json:"id"`
	Add    map[string]string `json:"add"`
	Delete map[string]string `json:"delete"`
}

type GetUser struct {
	ID       int      `json:"id"`
	Segments []string `json:"segments"`
}
