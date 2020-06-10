package cashout

import (
	"log"
	"testing"
)

func TestCashout_SupremeaddToCard(t *testing.T) {

	instance := NewCashout("webhook_url", nil)

	client, err := instance.getHTTPClient()
	if err != nil {
		log.Fatal(err)
	}

	err = instance.SupremeAddToCard(173019, 26482, 76049, client)
	if err != nil {
		log.Fatal(err)
	}

	err = instance.cashoutUs(client, "Nabil Tarik", "obito@gang.moe", "0782646623", "Menfou frérot", "San Diego", "P0M 0B8", "NB", "CANADA", "1234 1234 1234 1234", "02", "2023", "123")
	if err != nil {
		log.Fatal(err)
	}

	/*
		urlParsed, err := url.Parse("https://supremenewyork.com")
		if err != nil {
			log.Fatal(err)
		}

		cookies := client.Jar.Cookies(urlParsed)
		setcookies(page, ".supremenewyork.com", cookies)

		SupremeCashout(page, "Nabil Tarik", "obito@gmail.com", "0782644623", "Jmenfou frerot", "Paris", "7500", "FRANCE", "1234 1234 1234 1234", "04", "2025", "123")
	*/
}
