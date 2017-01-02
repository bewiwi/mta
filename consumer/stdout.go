package consumer

import "github.com/bewiwi/mta/models"

func RunStdout() {
	consume(func(ca models.CheckAnswer) {
		ca.Print()
	})
}
