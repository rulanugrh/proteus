package entity

type ResponseXendit struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		Event      string `json:"event"`
		BusinessID string `json:"businessID"`
		Data       struct {
			ID          string `json:"id"`
			Amount      string `json:"amout"`
			Country     string `json:"ID"`
			Currency    string `json:"IDR"`
			ReferenceID string `json:"reference_id"`
			Status      string `json:"status"`
		}
	}
}
