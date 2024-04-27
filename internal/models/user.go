package models

type User struct {
	ID       int      `json:"id"`
	Segments []string `json:"segments"`
}

type UpdateUserParams struct {
	ID             int
	AddSegments    []string
	DeleteSegments []string
}
