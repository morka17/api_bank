package mail

import (
	"testing"

	"github.com/morka17/shiny_bank/v1/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestSendEmailWithGmail(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	config, err := utils.LoadConfig("../../")
	assert.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
		<h1>Hello world</h1>
		<p>This is a test message from <a href="https://github.com/morka17/api_bank> API Bank</a></p>	
	`
	to := []string{config.EmailSenderAddress}
	attachFiles :=[]string{"../../readme.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	assert.NoError(t, err)
}