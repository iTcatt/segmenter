package requests

type CreateSegments struct {
	Segments []string `json:"segments"`
}

type CreateUsers struct {
	Users []int `json:"users"`
}

type UpdateUser struct {
	Id     int      `json:"id"`
	Add    []string `json:"add"`
	Delete []string `json:"delete"`
}
