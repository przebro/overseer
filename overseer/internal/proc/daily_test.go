package proc

import (
	"fmt"
	"testing"
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/stretchr/testify/suite"
)

type mockPoolConfig struct {
}

func (m *mockPoolConfig) NewDayProc() types.HourMinTime {
	return types.HourMinTime("00:00")
}
func (m *mockPoolConfig) CleanupCompletedTasks() int {
	return 0
}

type mockManager struct{}

func (M *mockManager) OrderNewTasks() int {
	return 0
}

type DailyTestSuite struct {
	suite.Suite
	daily *DailyExecutor
}

func TestDailyTestSuite(t *testing.T) {
	suite.Run(t, new(DailyTestSuite))
}

func (s *DailyTestSuite) SetupSuite() {

	manager := &mockManager{}
	poolconfig := &mockPoolConfig{}
	s.daily = NewDailyExecutor(manager, poolconfig)
}

func (s *DailyTestSuite) TestDailyProc() {

	del, ord := s.daily.DailyProcedure()
	s.Equal(del, 0)
	s.Equal(ord, 0)
}

func (s *DailyTestSuite) TestCheckDailyProc() {

	s.daily.lastExecutionDate = date.FromDateString("2019-01-01")

	tm := time.Now()
	// 	h, m, _ := tm.Clock()

	result := s.daily.CheckDailyProcedure(tm)
	//should contain test for positive and negative result
	fmt.Println(result)
}
