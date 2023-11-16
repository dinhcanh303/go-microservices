package mail

import (
	"testing"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestSendMail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	err := utils.LoadFileEnvOnLocal()
	require.NoError(t, err)
	config, err := configs.NewConfigMail()
	require.NoError(t, err)
	require.NotEmpty(t, config)

	senderMail := NewEmailSender(config)
	subject := "Email test"
	content := `<h1>Hello world</h1>
	<p>This is a test message from <a href="https://github.com/dinhcanh303">Foden Ngo</a></p>`

	to := []string{"dinhcanhng303@gmail.com"}
	// attachFiles := []string{"../../README.md"}
	err = senderMail.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}
