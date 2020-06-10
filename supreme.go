package cashout

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func (c *Cashout) SupremeAddToCard(id, style, size int, client *http.Client) error {
	shopURL := fmt.Sprintf("https://www.supremenewyork.com/shop/%v", id)

	reqCSRF, err := http.NewRequest("GET", fmt.Sprintf("https://www.supremenewyork.com/shop/%v", id), nil)
	if err != nil {
		return err
	}
	reqCSRF.Header.Set("Authority", "www.supremenewyork.com")
	reqCSRF.Header.Set("Upgrade-Insecure-Requests", "1")
	reqCSRF.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")
	reqCSRF.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	reqCSRF.Header.Set("Sec-Fetch-Site", "none")
	reqCSRF.Header.Set("Sec-Fetch-Mode", "navigate")
	reqCSRF.Header.Set("Sec-Fetch-User", "?1")
	reqCSRF.Header.Set("Sec-Fetch-Dest", "document")
	reqCSRF.Header.Set("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")

	respCSRF, err := client.Do(reqCSRF)
	if err != nil {
		return err
	}
	defer respCSRF.Body.Close()

	bodyCSRF, err := ioutil.ReadAll(respCSRF.Body)
	if err != nil {
		return err
	}

	token := strings.Split(strings.Split(string(bodyCSRF), `<meta name="csrf-token" content="`)[1], "\"")[0]
	body := strings.NewReader(fmt.Sprintf(`utf8=%%E2%%9C%%93&style=%v&size=%v&commit=add+to+basket`, style, size))

	req, err := http.NewRequest("POST", shopURL+"/add", body)
	if err != nil {
		log.Print("NIGGA?")
		log.Fatal(err)
	}

	req.Header.Set("Authority", "www.supremenewyork.com")
	req.Header.Set("Accept", "*/*;q=0.5, text/javascript, application/javascript, application/ecmascript, application/x-ecmascript")
	req.Header.Set("X-Csrf-Token", token)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://www.supremenewyork.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.supremenewyork.com/shop/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")

	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	err = resp.Body.Close()

	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Print(string(body))
		return err
	}

	return nil
}

func (c *Cashout) cashoutUs(client *http.Client, billingName, email, tel, address, city, zip, state, country, cnb, month, year, ccv string) error {
	reqCSRF, err := http.NewRequest("GET", "https://www.supremenewyork.com/checkout", nil)
	if err != nil {
		return err
	}
	reqCSRF.Header.Set("Authority", "www.supremenewyork.com")
	reqCSRF.Header.Set("Upgrade-Insecure-Requests", "1")
	reqCSRF.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")
	reqCSRF.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	reqCSRF.Header.Set("Sec-Fetch-Site", "none")
	reqCSRF.Header.Set("Sec-Fetch-Mode", "navigate")
	reqCSRF.Header.Set("Sec-Fetch-User", "?1")
	reqCSRF.Header.Set("Sec-Fetch-Dest", "document")
	reqCSRF.Header.Set("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")

	respCSRF, err := client.Do(reqCSRF)
	if err != nil {
		return err
	}
	defer respCSRF.Body.Close()

	bodyCSRF, err := ioutil.ReadAll(respCSRF.Body)
	if err != nil {
		return err
	}
	token := strings.Split(strings.Split(string(bodyCSRF), `<meta name="csrf-token" content="`)[1], "\"")[0]
	now := time.Now()
	secs := now.Unix()

	captchaResponse, err := requestCaptcha("6LeWwRkUAAAAAOBsau7KpuC9AV-6J8mhw4AjC3Xz", "supremenewyork.com")
	if err != nil {
		log.Fatal(err)
	}

	stringBody := `utf8=%E2%9C%93&authenticity_token=` + token + `&current_time=` + string(secs) + `&order%5Bbilling_name%5D=` + billingName + `&cerear=RMCEAR&order%5Bemail%5D=` + email + `&order%5Btel%5D=` + tel + `&order%5Bbilling_address%5D=` + address + `&order%5Bbilling_address_2%5D=&order%5Bbilling_zip%5D=` + zip + `&order%5Bbilling_city%5D=` + city + `&order%5Bbilling_state%5D=` + state + `&order%5Bbilling_country%5D=` + country + `&same_as_billing_address=1&store_credit_id=&riearmxa=` + cnb + `&credit_card%5Bmonth%5D=` + month + `&credit_card%5Byear%5D=` + year + `&credit_card%5Bmeknk%5D=` + ccv + `&credit_card%5Bvval%5D=&order%5Bterms%5D=0&order%5Bterms%5D=1&g-recaptcha-response=` + captchaResponse

	body := strings.NewReader(stringBody)
	req, err := http.NewRequest("POST", "https://www.supremenewyork.com/checkout.json", body)
	if err != nil {
		return err
	}

	req.Header.Set("Authority", "www.supremenewyork.com")
	req.Header.Set("Accept", "*/ /*")
	req.Header.Set("X-Csrf-Token", token)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://www.supremenewyork.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.supremenewyork.com/checkout")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fr;q=0.8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

/*
func (c *Cashout) SupremeCashout(page *rod.Page, billingName, email, tel, address, city, zip, country, cnb, month, year, ccv string) {
	log.Print("in cashout")

	page.Navigate("https://www.supremenewyork.com/checkout").Timeout(time.Minute)
	log.Print("loaded...")
	page.Element("#order_billing_name").Input(billingName)
	page.Element("#order_email").Input(billingName)
	page.Element("#order_tel").Input(billingName)
	page.Element("#bo").Input(billingName)
	page.Element("#order_billing_city").Input(billingName)
	page.Element("#order_billing_zip").Input(billingName)
	page.Element("#order_billing_country").Input(billingName)
	page.Element("#cnb").Input(billingName)
	page.Element("#credit_card_month").Input(billingName)
	page.Element("#credit_card_year").Input(billingName)
	page.Element("#vval").Input(billingName)
	page.Element("#cart-cc > fieldset > p > label > div > ins").Click()
	page.Element("#pay > input").Click()

	frame01 := page.Timeout(3 * time.Minute).ElementX("/html/body/div[3]/div[2]/iframe").Frame()

	log.Print(frame01.GetWindow().Get("width").Int())
	log.Print(frame01.GetWindow().Get("height").Int())
	page.Window(0, 0, frame01.GetWindow().Get("width").Int(), frame01.GetWindow().Get("height").Int())
	el := page.Timeout(3 * time.Minute).ElementX("/html/body/div[3]/div[2]/iframe").ScrollIntoView()
	el.WaitInvisible()
	return
}
*/
