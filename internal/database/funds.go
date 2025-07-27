package database

import "log"

type Fund struct {
	EmergencyFund float64 `json:"emergencyFund"`
	EducationFund float64 `json:"educationFund"`
	Investments   float64 `json:"investments"`
	Other         float64 `json:"other"`
}

func GetFunds() (Fund, error) {
	var fund Fund
	err := DB.QueryRow("SELECT emergency_fund, education_fund, investments, other FROM funds WHERE user_id = 1").Scan(&fund.EmergencyFund, &fund.EducationFund, &fund.Investments, &fund.Other)
	if err != nil {
		log.Println("Error getting funds:", err)
		return Fund{}, err
	}
	return fund, nil
}

func UpdateFunds(fund Fund) error {
	_, err := DB.Exec("UPDATE funds SET emergency_fund = $1, education_fund = $2, investments = $3, other = $4 WHERE user_id = 1", fund.EmergencyFund, fund.EducationFund, fund.Investments, fund.Other)
	if err != nil {
		log.Println("Error updating funds:", err)
		return err
	}
	return nil
}
