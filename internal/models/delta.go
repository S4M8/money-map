package models

type Delta struct {
	TotalIncome     float64 `json:"totalIncome"`
	TotalExpenses   float64 `json:"totalExpenses"`
	RemainingAmount float64 `json:"remainingAmount"`
	CoreExpenses    float64 `json:"coreExpenses"`
	ChoiceExpenses  float64 `json:"choiceExpenses"`
}