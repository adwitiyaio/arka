package email

import (
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/adwitiyaio/arka/cloud"
	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

type SesManagerTestSuite struct {
	suite.Suite

	m Manager
}

func TestSesManager(t *testing.T) {
	suite.Run(t, new(SesManagerTestSuite))
}

func (ts *SesManagerTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	cloud.Bootstrap(cloud.ProviderAws)
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
	Bootstrap(ProviderSes)
	ts.m = dependency.GetManager().Get(DependencyEmailManager).(Manager)
}

func (ts SesManagerTestSuite) Test_sesManager_SendEmail() {
	ts.Run("recipient limit", func() {
		var to []string
		for i := 0; i <= 1001; i++ {
			to = append(to, gofakeit.Email())
		}
		options := Options{
			Sender:  gofakeit.Email(),
			Subject: "Integration Testing",
			Html:    "<body>Hello</body>",
			Text:    "Hello",
			To:      to,
			Cc:      []string{gofakeit.Email()},
			Bcc:     []string{gofakeit.Email()},
		}

		res, err := ts.m.SendEmail(options)
		assert.NoError(ts.T(), err)
		assert.Nil(ts.T(), res)
	})

	ts.Run("success", func() {
		options := Options{
			Sender:  gofakeit.Email(),
			Subject: "Integration Testing",
			Html:    "<body>Hello</body>",
			Text:    "Hello",
			To:      []string{gofakeit.Email(), gofakeit.Name()},
			Cc:      []string{gofakeit.Email()},
			Bcc:     []string{gofakeit.Email()},
		}

		res, err := ts.m.SendEmail(options)
		assert.NoError(ts.T(), err)
		assert.Nil(ts.T(), res)
	})

	ts.Run("success - with attachments", func() {
		options := Options{
			Sender:      gofakeit.Email(),
			Subject:     "Integration Testing",
			Html:        "<body>Hello</body>",
			Text:        "Hello",
			To:          []string{gofakeit.Email(), gofakeit.Name()},
			Cc:          []string{gofakeit.Email()},
			Bcc:         []string{gofakeit.Email()},
			Attachments: []string{"./testdata/sample.txt"},
		}

		res, err := ts.m.SendEmail(options)
		assert.NoError(ts.T(), err)
		assert.Nil(ts.T(), res)
	})
}
