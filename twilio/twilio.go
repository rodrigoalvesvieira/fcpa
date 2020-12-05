package twilio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"operatordesk/envars"
)

// SendMessage receives a phone number
func SendMessage(phoneNumber string) {
	developmentMode := true

	// callToBuyEuro := "Vale a pena comprar Euro agora por que a balança está para o lucro."
	var amount int64 = 300

	callToBuyCrypto := fmt.Sprintf("Agora é hora de comprar Bitcoin se o Euro tiver sido comprado a R$%s",
		strconv.FormatInt(amount, 10))

	if developmentMode {
	} else {
		// Set account keys & information
		accountSid := envars.GetTwilioAccountSID()
		authToken := envars.GetTwilioAuthToken()
		urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

		appPhoneNumber := envars.GetTwilioPhoneNumber()

		// Pack up the data for our message
		msgData := url.Values{}
		msgData.Set("To", phoneNumber)
		msgData.Set("From", appPhoneNumber)
		msgData.Set("Body", callToBuyCrypto)
		msgDataReader := *strings.NewReader(msgData.Encode())

		// Create HTTP request client
		client := &http.Client{}
		req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
		req.SetBasicAuth(accountSid, authToken)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Make HTTP POST request and return message SID
		resp, _ := client.Do(req)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var data map[string]interface{}
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&data)

			if err == nil {
				fmt.Println(data["sid"])
			}
		} else {
			fmt.Println(resp.Status)
		}
	}
}
