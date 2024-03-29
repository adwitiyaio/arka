package monitoring

import (
	"os"
	"strconv"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const appNameKey = "APP_NAME"
const newRelicLicenseKey = "NEW_RELIC_LICENSE"
const enableMonitoringKey = "ENABLE_MONITORING"

type newRelicManager struct {
	sm secrets.Manager

	app *newrelic.Application
}

func (c *newRelicManager) initialize() {
	if os.Getenv("CI") == "true" {
		return
	}
	enabledMonitoringConfig := strings.TrimSpace(c.sm.GetValueForKey(enableMonitoringKey))
	enableMonitoring, err := strconv.ParseBool(enabledMonitoringConfig)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to parse new relic configuration")
		return
	}
	if !enableMonitoring {
		logger.Log.Error().Err(err).Msg("monitoring not enabled")
		return
	}
	appName := strings.TrimSpace(c.sm.GetValueForKey(appNameKey))
	license := strings.TrimSpace(c.sm.GetValueForKey(newRelicLicenseKey))
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(license),
		newrelic.ConfigDistributedTracerEnabled(true))
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to initialize new relic")
		return
	}
	c.app = app
}
