package cashout

import (
	"log"
	"testing"
)

func TestCaptchaClient_SolveReCaptcha(t *testing.T) {
	capClient := New("x")

	answer, err := capClient.SolveReCaptcha("https://www.supremenewyork.com/", "6LeWwRkUAAAAAOBsau7KpuC9AV-6J8mhw4AjC3Xz")

	if err != nil {
		log.Fatal(err)
	}

	log.Print(answer)
}
