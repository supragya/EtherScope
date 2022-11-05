package types

type ResumeAPIResponse struct {
	Data struct {
		Height uint64 `json:"block_height"`
	} `json:"data"`
}
