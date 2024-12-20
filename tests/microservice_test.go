package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestMicroservice is the entry point for the microservice test suite.
// service and api tests are included in this suite, but declared as a suite messages.
func TestMicroservice(t *testing.T) {
	suite.Run(t, new(MicroserviceSuite))
}
