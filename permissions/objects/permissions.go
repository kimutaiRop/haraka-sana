package objects

type CreatePosition struct {
	Name string `json:"name"`
}

type CreatePositionPermission struct {
	Position   int  `json:"position"`
	Permission int  `json:"permission"`
	Active     bool `json:"active"`
}
