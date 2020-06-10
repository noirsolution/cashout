package cashout

import (
	"log"
	"testing"
)

func TestCashout_ShopifyAddToCart(t *testing.T) {
	instance := NewCashout("webhook_url", nil)

	client, err := instance.getHTTPClient()
	if err != nil {
		log.Fatal(err)
	}

	err = instance.ShopifyAddToCart("https://stussy.co.uk", "sophnet-s-s-standard-big-b-d-shirt-indigo", 31458014658671, client)
	if err != nil {
		log.Fatal(err)
	}

	finalURL, authenticityToken, gatewayID, priceTarget, err := instance.ShopifyStartCashout("https://stussy.co.uk", "Nabil", "Tarik", "oklmzer@gmail.com", "Jmen fou frérot", "7500", "Paris", "FRANCE", "0649381931", client)
	if err != nil {
		log.Fatal(err)
	}

	err = instance.ShopifyPayment(finalURL, authenticityToken, gatewayID, priceTarget, "Nabil Tarik", "1234123412341234", 03, 20, "123", client)
	if err != nil {
		log.Fatal(err)
	}
}
