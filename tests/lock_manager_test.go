package dynalock_tests

import (
	"github.com/kmdeveloping/dynalock"
	"github.com/stretchr/testify/suite"
)

type lock_manager_UnitTestSuite struct {
	suite.Suite

	client *dynalock.Client
}

func (s *lock_manager_UnitTestSuite) SetupSuite() {

}
