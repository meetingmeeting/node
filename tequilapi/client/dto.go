package client

// StatusDto holds connection status and session id
type StatusDto struct {
	Status    string `json:"status"`
	SessionId string `json:"sessionId"`
}

// IdentityDto holds identity address
type IdentityDto struct {
	Address string `json:"id"`
}

// IdentityList holds returned list of identities
type IdentityList struct {
	Identities []IdentityDto `json:"identities"`
}
