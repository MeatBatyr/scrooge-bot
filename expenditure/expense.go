package expenditure

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ErrManyCategories = errors.New("You can't specify multiple categories")

const NothingForCalc = "There is nothing for calculation"
const moneyRegex = `^\$?([1-9]{1}[0-9]{0,2}(\,[0-9]{3})*(\.[0-9]{0,2})?|[1-9]{1}[0-9]{0,}(\.[0-9]{0,2})?|0(\.[0-9]{0,2})?|(\.[0-9]{1,2})?)$`

type Expense struct {
	Id        int64     //identifier
	Category  string    //canetory
	Amount    float64   //money
	ChatId    int64     //chat identifier
	CreatedAt time.Time //created date
	CreatedBy string    //user
}

func (r *Expense) SetContextData(chatId int64, createdAt string) {
	r.ChatId = chatId
	r.CreatedBy = createdAt
}

func Parse(source string) (r *Expense, err error) {
	category := ""
	money := 0.0
	words := strings.Split(source, " ")

	for _, word := range words {
		if word == "" {
			continue
		}

		if strings.HasPrefix(word, "#") {
			if category != "" {
				return nil, ErrManyCategories
			}

			category = word[1:]
		}

		if validID := regexp.MustCompile(moneyRegex); validID.MatchString(word) {
			var currentPayment float64
			currentPayment, err = strconv.ParseFloat(word, 64)

			if err != nil {
				return nil, err
			}

			money += currentPayment
		}
	}

	if category == "" || money <= 0.0 {
		return nil, nil
	}

	result := &Expense{Category: category, Amount: money, CreatedAt: time.Now().UTC()}
	return result, nil
}

func MakeCalculationResult(expenses []Expense) string {
	text := ""
	results := make(map[string]float64)
	for _, v := range expenses {
		results[v.Category] += v.Amount
	}

	for key := range results {
		text += key + ": " + strconv.FormatFloat(results[key], 'f', -1, 32) + "\n"
	}

	if text == "" {
		return NothingForCalc
	}
	return strings.TrimSuffix(text, "\n")
}
