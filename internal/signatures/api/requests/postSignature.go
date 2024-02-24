package requests

type PostSignature struct {
	Jwt       string   `json:"jwt"`
	Questions []string `json:"questions"`
	Answers   []string `json:"answers"`
}
