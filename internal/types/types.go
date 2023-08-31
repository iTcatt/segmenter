package types

type User struct {
	Id       int
	Segments []string
}

type CreateRequest struct {
	Segments []string `json: segments`
}

type UpdateUser struct {
	Id     int      `json: id`
	Add    []string `json: add`
	Delete []string `json: delete`
}
