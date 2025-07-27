package api

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"money-map/internal/database"
	"money-map/internal/models"

	"github.com/gorilla/mux"
)

func GetIncome(w http.ResponseWriter, r *http.Request) {
	month, year := getMonthYear(r)
	income, err := database.GetIncome(month, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if income == nil {
		income = []models.Transaction{}
	}
	json.NewEncoder(w).Encode(income)
}

func GetExpenses(w http.ResponseWriter, r *http.Request) {
	month, year := getMonthYear(r)
	expenses, err := database.GetExpenses(month, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if expenses == nil {
		expenses = []models.Transaction{}
	}
	json.NewEncoder(w).Encode(expenses)
}

type TransactionRequest struct {
	Date     string  `json:"date"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
}

func AddIncome(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	transaction := models.Transaction{
		Date:     date,
		Name:     req.Name,
		Amount:   req.Amount,
		Category: "Income",
		Type:     "Income",
	}

	if err := database.AddIncome(transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DeleteIncome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid income ID", http.StatusBadRequest)
		return
	}

	if err := database.DeleteIncome(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AddExpense(w http.ResponseWriter, r *http.Request) {
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	transaction := models.Transaction{
		Date:     date,
		Name:     req.Name,
		Amount:   req.Amount,
		Category: req.Category,
		Type:     "Expense",
	}

	if err := database.AddExpense(transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid expense ID", http.StatusBadRequest)
		return
	}

	if err := database.DeleteExpense(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetFunds(w http.ResponseWriter, r *http.Request) {
	funds, err := database.GetFunds()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(funds)
}

func UpdateFunds(w http.ResponseWriter, r *http.Request) {
	var funds database.Fund
	if err := json.NewDecoder(r.Body).Decode(&funds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.UpdateFunds(funds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetDelta(w http.ResponseWriter, r *http.Request) {
	month, year := getMonthYear(r)
	delta, err := database.GetDelta(month, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(delta)
}

func UploadCSV(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	existingExpenses, err := database.GetExpenses(0, 0)
	if err != nil {
		http.Error(w, "Failed to get existing expenses", http.StatusInternalServerError)
		return
	}

	categorizationMap := make(map[string]string)
	for _, expense := range existingExpenses {
		categorizationMap[expense.Name] = expense.Category
	}

	reader := csv.NewReader(file)
	// Skip header
	if _, err := reader.Read(); err != nil {
		http.Error(w, "Failed to read header from CSV", http.StatusInternalServerError)
		return
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Failed to read record from CSV", http.StatusInternalServerError)
			return
		}

		if len(record) < 8 {
			log.Printf("Skipping malformed row: %#v", record)
			continue
		}

		date, err := time.Parse("01/02/06", record[2])
		if err != nil {
			log.Printf("Skipping row with invalid date: %v", err)
			continue 
		}

		amountStr := strings.Replace(record[7], "$", "", -1)
		amountStr = strings.Replace(amountStr, ",", "", -1)
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			log.Printf("Skipping row with invalid amount: %v", err)
			continue 
		}

		transaction := models.Transaction{
			Date: date,
			Name: record[3],
		}

		if amount > 0 {
			transaction.Amount = amount
			transaction.Type = "Income"
			if err := database.AddIncome(transaction); err != nil {
				log.Printf("Failed to add income from CSV: %v", err)
				http.Error(w, "Failed to add income from CSV", http.StatusInternalServerError)
				return
			}
		} else {
			transaction.Amount = -amount 
			transaction.Type = "Expense"

			if category, ok := categorizationMap[transaction.Name]; ok {
				transaction.Category = category
			} else {
				transaction.Category = mapCsvCategory(record[5])
			}

			if err := database.AddExpense(transaction); err != nil {
				log.Printf("Failed to add expense from CSV: %v", err)
				http.Error(w, "Failed to add expense from CSV", http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func getMonthYear(r *http.Request) (int, int) {
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	if monthStr == "" || yearStr == "" {
		now := time.Now()
		return int(now.Month()), now.Year()
	}

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	return month, year
}

func mapCsvCategory(csvCategory string) string {
	csvCategory = strings.ToLower(csvCategory)

	coreCategories := []string{"groceries", "credit card payments", "pharmacy", "doctor", "hospital", "utilities", "rent", "mortgage", "insurance"}
	choiceCategories := []string{"atm/cash withdrawals", "restaurants/dining", "general merchandise", "shopping", "entertainment"}

	for _, cat := range coreCategories {
		if strings.Contains(csvCategory, cat) {
			return "Core"
		}
	}

	for _, cat := range choiceCategories {
		if strings.Contains(csvCategory, cat) {
			return "Choice"
		}
	}

	return "Choice"
}
