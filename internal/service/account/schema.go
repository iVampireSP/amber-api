package account

type UnitReduceRequest struct {
	Amount string `json:"amount"`
	Reason string `json:"reason"`
	Unit   string `json:"unit"`
	UserId string `json:"user_id"`
}

type CanBillUnitRequest struct {
	UserId string `json:"user_id"`
	Unit   string `json:"unit"`
}

type CanBill struct {
	CanBill bool `json:"can_bill"`
}
