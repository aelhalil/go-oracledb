/*
** Copyright (c) 2026 Oracle and/or its affiliates.
**
** The Universal Permissive License (UPL), Version 1.0
**
** Subject to the condition set forth below, permission is hereby granted to any
** person obtaining a copy of this software, associated documentation and/or data
** (collectively the "Software"), free of charge and under any and all copyright
** rights in the Software, and any and all patent rights owned or freely
** licensable by each licensor hereunder covering either (i) the unmodified
** Software as contributed to or provided by such licensor, or (ii) the Larger
** Works (as defined below), to deal in both
**
** (a) the Software, and
** (b) any piece of software and/or hardware listed in the lrgrwrks.txt file if
** one is included with the Software (each a "Larger Work" to which the Software
** is contributed by such licensors),
**
** without restriction, including without limitation the rights to copy, create
** derivative works of, display, perform, and distribute the Software and make,
** use, sell, offer for sale, import, export, have made, and have sold the
** Software and the Larger Work(s), and to sublicense the foregoing rights on
** either these or other terms.
**
** This license is subject to the following condition:
** The above copyright notice and either this complete permission notice or at
** a minimum a reference to the UPL must be included in all copies or
** substantial portions of the Software.
**
** THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
** IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
** FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
** AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
** LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
** OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
** SOFTWARE.
 */

package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	err := InitConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "InitConfig failed: %v\n", err)
		os.Exit(1)
	} else {
		os.Exit(m.Run())
	}
}

var testCases = []struct {
	name       string
	categories string
	exclusive  bool
	f          func(t *testing.T)
}{}

func TestCategoryExecutor(t *testing.T) {
	var regularCases, exclusiveCases []struct {
		name       string
		categories string
		exclusive  bool
		f          func(t *testing.T)
	}

	for _, c := range testCases {
		cats := strings.Split(c.categories, ",")
		for _, p := range cats {
			if strings.Compare(strings.TrimSpace(p), TestCategory) == 0 {
				if c.exclusive {
					exclusiveCases = append(exclusiveCases, c)
				} else {
					regularCases = append(regularCases, c)
				}
				break
			}
		}
	}

	if len(regularCases) > 0 {
		t.Run("parallel", func(t *testing.T) {
			t.Parallel()
			for _, c := range regularCases {
				t.Run(c.name, c.f)
			}
		})
	}

	for _, c := range exclusiveCases {
		t.Run(c.name, c.f)
	}
}

// TestConfig structure that represents a testing configuration
// A configuration contains any information neede to connect to a database
type TestConfig struct {
	ConfigName      string `json:"config_name"`
	DatabaseVersion int    `json:"database_version"`
	Enabled         bool

	Driver struct {
		Name string
	}

	Database struct {
		ServiceName  string
		SIDName      string `json:",omitempty"`
		InstanceName string `json:",omitempty"`
		Port         int16
		Host         string
		Protocol     string
		ServerType   string `json:",omitempty"` // dedicated/shared
	}

	Credentials struct {
		Username  string
		Password  string
		LogonMode string
	}

	Security struct {
		WalletLocation      string `json:"wallet_location,omitempty"`
		SslServerDnMatch    string `json:"ssl_server_dn_match,omitempty"`
		SslServerCertDn     string `json:"ssl_server_cert_dn,omitempty"`
		SslAllowWeakDnMatch string `json:"ssl_allow_weak_dn_match,omitempty"`
	}

	ConnectionProperties struct {
		StrictNullValueHandling string `json:"oracle.go.StrictNullValueHandling"`
	}
}

// _assignStringIfNeeded assign "from" value to src if from is a valid string
func _assignStringIfNeeded(src *string, from string) {
	if len(strings.TrimSpace(from)) > 0 {
		*src = from
	}
}

// _assignIntIfNeeded assign "from" value to src if from is a valid int
func _assignIntIfNeeded(src *int16, from int16) {
	if from >= 0 {
		*src = from
	}
}

// Clone clones a test config
func (t *TestConfig) Clone() *TestConfig {
	newOne := &TestConfig{}
	newOne.Driver.Name = t.Driver.Name
	newOne.Database.ServiceName = t.Database.ServiceName
	newOne.Database.SIDName = t.Database.SIDName
	newOne.Database.InstanceName = t.Database.InstanceName
	newOne.Database.Host = t.Database.Host
	newOne.Database.Port = t.Database.Port
	newOne.Database.Protocol = t.Database.Protocol
	newOne.Database.ServerType = t.Database.ServerType

	newOne.Credentials.Username = t.Credentials.Username
	newOne.Credentials.Password = t.Credentials.Password
	newOne.Credentials.LogonMode = t.Credentials.LogonMode

	newOne.Security = t.Security

	return newOne
}

// MergeWith merges config "from" with value currently assigned
//
//	returned the merged config
func (t *TestConfig) MergeWith(from *TestConfig) {

	_assignStringIfNeeded(&(t.Driver.Name), from.Driver.Name)

	_assignStringIfNeeded(&(t.Database.ServiceName), from.Database.ServiceName)
	_assignStringIfNeeded(&(t.Database.SIDName), from.Database.SIDName)
	_assignStringIfNeeded(&(t.Database.InstanceName), from.Database.InstanceName)
	_assignStringIfNeeded(&(t.Database.ServerType), from.Database.ServerType)
	_assignIntIfNeeded(&(t.Database.Port), from.Database.Port)
	_assignStringIfNeeded(&(t.Database.Host), from.Database.Protocol)
	_assignStringIfNeeded(&(t.Database.Protocol), from.Database.Protocol)

	_assignStringIfNeeded(&(t.Credentials.Username), from.Credentials.Username)
	_assignStringIfNeeded(&(t.Credentials.Password), from.Credentials.Password)
	_assignStringIfNeeded(&(t.Credentials.LogonMode), from.Credentials.LogonMode)

	_assignStringIfNeeded(&(t.Security.WalletLocation), from.Security.WalletLocation)
	_assignStringIfNeeded(&(t.Security.SslServerDnMatch), from.Security.SslServerDnMatch)
	_assignStringIfNeeded(&(t.Security.SslServerCertDn), from.Security.SslServerCertDn)
	_assignStringIfNeeded(&(t.Security.SslAllowWeakDnMatch), from.Security.SslAllowWeakDnMatch)
}

// GetConnectionString Build a connection string from a test config
func (t *TestConfig) GetConnectionDSN() string {
	dsn := t.GetConnectionStringWithProperties(nil)
	s := strings.SplitN(dsn, "@", 2)
	return s[1]
}

// GetConnectionString Build a connection string from a test config
func (t *TestConfig) GetConnectionString() string {
	return t.GetConnectionStringWithProperties(nil)
}

// GetConnectionStringWithProperties Build a connection string from a test config and some properties
func (t *TestConfig) GetConnectionStringWithProperties(properties map[string]string) string {
	var b strings.Builder
	if properties != nil {
		for k, v := range properties {
			b.WriteString(fmt.Sprintf("(%s=%s)", k, v))
		}
	}
	var res = fmt.Sprintf("%s/%s@(description=%s(address=(protocol=%s)(host=%s)(port=%d))(connect_data=",
		t.Credentials.Username,
		t.Credentials.Password,
		b.String(),
		t.Database.Protocol,
		t.Database.Host,
		t.Database.Port)

	var resC strings.Builder
	resC.WriteString(res)
	if len(t.Database.ServiceName) > 0 {
		resC.WriteString(fmt.Sprintf("(service_name=%s)", t.Database.ServiceName))
		if len(t.Database.InstanceName) > 0 {
			resC.WriteString(fmt.Sprintf("(instance_name=%s)", t.Database.InstanceName))
		}
	} else {
		if len(t.Database.SIDName) > 0 {
			resC.WriteString(fmt.Sprintf("(sid=%s)", t.Database.SIDName))
		}
	}

	if len(t.Database.ServerType) > 0 {
		resC.WriteString(fmt.Sprintf("(server=%s)", t.Database.ServerType))
	}

	resC.WriteString(")") // close connect_data

	// Add security if present
	if t.Security.WalletLocation != "" || t.Security.SslServerDnMatch != "" ||
		t.Security.SslServerCertDn != "" || t.Security.SslAllowWeakDnMatch != "" {
		resC.WriteString("(security=")
		if t.Security.WalletLocation != "" {
			resC.WriteString(fmt.Sprintf("(wallet_location=%s)", t.Security.WalletLocation))
		}
		if t.Security.SslServerDnMatch != "" {
			resC.WriteString(fmt.Sprintf("(ssl_server_dn_match=%s)", t.Security.SslServerDnMatch))
		}
		if t.Security.SslAllowWeakDnMatch != "" {
			resC.WriteString(fmt.Sprintf("(ssl_allow_weak_dn_match=%s)", t.Security.SslAllowWeakDnMatch))
		}
		if t.Security.SslServerCertDn != "" {
			resC.WriteString(fmt.Sprintf("(ssl_server_cert_dn=\"%s\")", t.Security.SslServerCertDn))
		}
		resC.WriteString(")")
	}

	resC.WriteString(")") // close description

	if len(t.Credentials.LogonMode) > 0 {
		resC.WriteString(fmt.Sprintf("?oracle.go.Credentials.logonMode=%s", t.Credentials.LogonMode))
	}
	return resC.String()
}

// GetConnectionStringWithMergedConfig Build a connection string from a test config after merging with a given config
func (t *TestConfig) GetConnectionStringWithMergedConfig(config *TestConfig) string {

	_c := t.Clone()
	_c.MergeWith(config)
	return _c.GetConnectionStringWithProperties(nil)

}

// TestingEnvironment Holds driver configuration for tests
type TestingEnvironment struct {
	// Testing configuration array parsec from YAML file
	driverConfigs []TestConfig
}

// DefaultTestConfig Default reference to TestEnvironement
// That should not be that way but we need
// a way to pass config to sub package.
// that will be removed after refactoring
var DefaultTestConfig *TestConfig

// NewTestingEnvironment creates a new environment for given file
// On failure, error is returned
func NewTestingEnvironment(fileName string) (TestingEnvironment, error) {

	var driverConfigs []TestConfig

	// load YAML file
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return TestingEnvironment{},
			fmt.Errorf("specified configuration file %s do not exists", fileName)
	}
	f, err := os.Open(fileName)
	if err != nil {
		return TestingEnvironment{},
			fmt.Errorf("unable to open configuration %s: %v",
				fileName,
				err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			//ignored
		}
	}(f)

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&driverConfigs)
	if err != nil {
		return TestingEnvironment{},
			fmt.Errorf("unable to read configuration %s: %w", fileName, err)
	}

	return TestingEnvironment{
		driverConfigs: driverConfigs,
	}, nil
}

// GetConfig gets a configuration by name.
// Returns null if the configuration is not found
func (e *TestingEnvironment) GetConfig(name string) (*TestConfig, error) {
	if e.driverConfigs == nil {
		return nil, fmt.Errorf("attempt to get a configuration but not configuration available")
	}
	for _, config := range e.driverConfigs {
		if config.ConfigName == name {
			return &config, nil
		}
	}
	return nil, fmt.Errorf("no configuration %s found", name)

}

// Test configuration flag. This flag gives configuration file path
//
//	is mandatory to run the test
var configFileName string

// Test configuration flag. This flag will specify which
// configuration to use.
var configName string

var TestEnvironement TestingEnvironment

// TestingConfig Usable by tests, may be nil if not flag provided
var TestingConfig *TestConfig

// TestCategory category of tests to be un
var TestCategory string

// InitConfig init the environment configuration
// Main task is to parse dirver config flags and populate the
// default configuration
func InitConfig() error {

	flag.StringVar(&configFileName, "driver.config.filename", "", "testing config name")
	flag.StringVar(&configName, "driver.config.name", "", "testing config name")

	flag.StringVar(&TestCategory, "test.category", "", "testing category, can be unitary, functional, performance, robustness")

	if !flag.Parsed() {
		flag.Parse()
	}
	if len(configFileName) != 0 {
		env, err := NewTestingEnvironment(configFileName)
		if err != nil {
			return fmt.Errorf("cannot get test environment : %w", err)
		}
		TestEnvironement = env

		if len(configName) != 0 {
			TestingConfig, _ = TestEnvironement.GetConfig(configName)
			// Keep DefaultTestConfig in sync for legacy driver API
			DefaultTestConfig = TestingConfig
		}
	}
	return nil
}
