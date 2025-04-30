package req

type UpdateLabelReq struct {
	Name   string `json:"name"`
	Labels []Item `json:"labels"`
}
