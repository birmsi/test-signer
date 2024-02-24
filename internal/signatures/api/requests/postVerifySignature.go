package requests

type PostVerifySignature struct {
	Jwt       string `json:"jwt"`
	Signature []byte `json:"signature"`
}
