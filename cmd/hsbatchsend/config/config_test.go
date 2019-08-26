package config

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {

	// read config
	err := Init("testConfig1", "./")
	assert.Equal(t, err, nil)

	strValue := GetString("log.logFile", "")
	fmt.Printf("strValue=%s\n", strValue)
}
