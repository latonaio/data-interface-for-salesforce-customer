package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/latonaio/salesforce-data-models"
	"github.com/latonaio/aion-core/pkg/log"
)

func HandleAccount(metadata map[string]interface{}) error {
	customers, err := models.MetadataToCustomers(metadata)
	if err != nil {
		return fmt.Errorf("failed to convert customers: %v", err)
	}
	for _, customer := range customers {
		if customer.SfCustomerID == nil {
			continue
		}
		c, err := models.CustomerByID(*customer.SfCustomerID)
		if err != nil {
			log.Printf("failed to get customer: %v", err)
			continue
		}
		if customer.Birthday != nil {
			age,err := calcAge(customer.Birthday.Time)
			if err != nil {
				log.Printf("failed to calculate customer's age: %v", err)
				continue
			}
			customer.Age = &age
		}
		if c != nil {
			log.Printf("update customer: %s\n", *customer.SfCustomerID)
			if err := customer.Update(); err != nil {
				log.Printf("failed to update customer: %v", err)
				continue
			}
		} else {
			log.Printf("register customer: %s\n", *customer.SfCustomerID)
			if err := customer.Register(); err != nil {
				log.Printf("failed to register customer: %v", err)
			}
		}
	}
	return nil
}


func calcAge(t time.Time) (int, error) {
	dateFormatOnlyNumber := "20060102"
	now := time.Now().Format(dateFormatOnlyNumber)
	birthday := t.Format(dateFormatOnlyNumber)

	// 日付文字列をそのまま数値化
	nowInt, err := strconv.Atoi(now)
	if err != nil {
		return 0, err
	}
	birthdayInt, err := strconv.Atoi(birthday)
	if err != nil {
		return 0, err
	}

	// (今日の日付 - 誕生日) / 10000 = 年齢
	age := (nowInt - birthdayInt) / 10000
	return age, nil
}