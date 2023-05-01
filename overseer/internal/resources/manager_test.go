package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ResourceManagerTestSuite struct {
	suite.Suite
}

func TestResourceManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ResourceManagerTestSuite))
}

func (s *ResourceManagerTestSuite) SetupSuite() {

}

func (s *ResourceManagerTestSuite) TestSingleTest() {
	fmt.Println("single test")
}

func TestBuildExpr(t *testing.T) {

	result := buildExpr("")
	if result != `[\w\-]*|^$` {
		t.Error("unexpected result")
	}

	result = buildExpr("*")
	if result != `^[\w\-]*$` {
		t.Error("unexpected result")
	}

	result = buildExpr("?")
	if result != `^[\w\-]{1}$` {
		t.Error("unexpected result")
	}

	result = buildExpr("AB?X")
	if result != `^AB[\w\-]{1}X$` {
		t.Error("unexpected result:", result)
	}

	result = buildExpr("AB?*")
	if result != `^AB[\w\-]{1}[\w\-]*$` {
		t.Error("unexpected result:", result)
	}
}
func TestBuildDateExpr(t *testing.T) {

	result := buildDateExpr("")
	if result != `[\d]*|^$` {
		t.Error("unexpected result")
	}

	result = buildDateExpr("*")
	if result != `^[\d]*$` {
		t.Error("unexpected result")
	}

	result = buildDateExpr("?")
	if result != `^[\d]{1}$` {
		t.Error("unexpected result")
	}

	result = buildDateExpr("2021051?")
	if result != `^2021051[\d]{1}$` {
		t.Error("unexpected result")
	}

}
