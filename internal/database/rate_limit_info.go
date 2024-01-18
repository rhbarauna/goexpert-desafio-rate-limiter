package database

type RateLimitInfo struct {
	TtlLimit int `json:"ttl_limit"`
	ReqLimit int `json:"req_limit"`
	Cooldown int `json:"cooldown"`
}
