package sender

import "fmt"

type MailService struct {

}

func NewMailService() *MailService {
	return &MailService{}
}

// Sender sending mail
func (s *MailService) Send(msg interface{}) error {
	fmt.Println("sending...", msg)
	return nil
}
