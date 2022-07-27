package ExternalService

import (
	"encoding/json"
	"fmt"
	"github.com/botaevg/gophermart/internal/repositories"
	"io"
	"log"
	"net/http"
)

type ExternalService struct {
	storage              repositories.Storage
	accrualSystemAddress string
}

type OrderES struct {
	order   string `json:"order"`
	status  string `json:"status"`
	accrual uint   `json:"accrual"`
}

func NewES(storage repositories.Storage, accrualSystemAddress string) ExternalService {
	return ExternalService{
		storage:              storage,
		accrualSystemAddress: accrualSystemAddress,
	}
}

func (e ExternalService) AccrualPoints(orderID uint) {
	client := http.Client{}

	URL := fmt.Sprintf("http://localhost%s/api/orders/%d", e.accrualSystemAddress, orderID)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Print(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(resp.Status)
	if resp.Status == "200 OK" {
		respBody, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()

		var Order OrderES
		err = json.Unmarshal(respBody, &Order)
		if err != nil {
			log.Print(err)
		}
		log.Print(Order)
	}
}
