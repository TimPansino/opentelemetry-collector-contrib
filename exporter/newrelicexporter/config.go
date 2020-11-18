// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package newrelicexporter

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
	"go.opentelemetry.io/collector/config/configmodels"
)

// Config defines configuration options for the New Relic exporter.
type Config struct {
	configmodels.ExporterSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.

	// APIKey is the required authentication credentials for New Relic APIs.
	APIKey string `mapstructure:"apikey"`

	// Timeout is the total amount of time spent attempting a request,
	// including retries, before abandoning and dropping data. Default is 15
	// seconds.
	Timeout time.Duration `mapstructure:"timeout"`

	// CommonAttributes are the attributes to be applied to all telemetry
	// sent to New Relic.
	CommonAttributes map[string]interface{} `mapstructure:"common_attributes"`

	// MetricsURLOverride overrides the metrics endpoint.
	MetricsURLOverride string `mapstructure:"metrics_url_override"`

	// InsecureSkipVerify sets the http transport to InsecureSkipVerify: true.
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`

	// SpansURLOverride overrides the spans endpoint.
	SpansURLOverride string `mapstructure:"spans_url_override"`
}

// HarvestOption sets all relevant Config values when instantiating a New
// Relic Harvester.
func (c Config) HarvestOption(cfg *telemetry.Config) {
	cfg.APIKey = c.APIKey
	cfg.HarvestPeriod = 0 // use collector harvest period.
	cfg.HarvestTimeout = c.Timeout
	cfg.CommonAttributes = c.CommonAttributes
	cfg.Product = product
	cfg.ProductVersion = version
	cfg.MetricsURLOverride = c.MetricsURLOverride
	cfg.SpansURLOverride = c.SpansURLOverride

	if c.InsecureSkipVerify {
		// No error logger is configured by this point and this is
		// primarily for testing use, so a panic is acceptable.

		if cfg.Client == nil {
			cfg.Client = &http.Client{}
			cfg.Client.Transport = http.DefaultTransport.(*http.Transport).Clone()
		} else if cfg.Client.Transport == nil {
			cfg.Client.Transport = http.DefaultTransport.(*http.Transport)
		}

		customTransport := cfg.Client.Transport.(*http.Transport).Clone()
		if customTransport.TLSClientConfig == nil {
			customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		} else {
			customTransport.TLSClientConfig.InsecureSkipVerify = true
		}

		cfg.Client.Transport = customTransport
	}
}
