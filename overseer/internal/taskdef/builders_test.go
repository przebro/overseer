package taskdef

import (
	"testing"
)

func TestBuilderFromTemplate(t *testing.T) {

	builder, builder2 := NewBuilder(), NewBuilder()

	sd := SchedulingData{FromTime: "10:30", OrderType: OrderingManual}
	intd := []InTicketData{{Name: "TICKET01", Odate: "ODATE"}}
	outtd := []OutTicketData{{Action: "ADD", Name: "TICKET01", Odate: "ODATE"}}

	task, err := builder.WithBase("testgroup", "testname", "testdescription").WithSchedule(SchedulingData{OrderType: OrderingManual}).WithConfirm().Build()
	if err != nil {
		t.Error(err)
	}

	task2, err := builder2.FromTemplate(task).WithSchedule(sd).WithInTicekts(intd, "AND").WithOutTickets(outtd).Build()
	if err != nil {
		t.Error(err)
	}

	n, g, d := task2.GetInfo()

	if n != "testname" || g != "testgroup" || d != "testdescription" {
		a, b, c := task.GetInfo()
		t.Error("Unexpected values:", n, g, d, "expected:", a, b, c)
	}

	from, to := task2.TimeSpan()

	if from != "10:30" || to != "" {
		t.Error("Unexpected values, expected:", "10:30", "", true, "actual:", from, to)
	}

	if task2.OutTickets[0].Name != intd[0].Name && task2.InTickets[0].Odate != intd[0].Odate {
		t.Error("Unexpected values, expected:", intd[0].Name, intd[0].Odate, "actual:", task2.InTickets[0].Name, task2.InTickets[0].Odate)

	}

	if task2.OutTickets[0].Name != outtd[0].Name && task2.InTickets[0].Odate != outtd[0].Odate {
		t.Error("Unexpected values, expected:", outtd[0].Name, outtd[0].Odate, "actual:", task2.OutTickets[0].Name, task2.OutTickets[0].Odate)

	}

}

func TestBuilderFromTemplate_Tickets(t *testing.T) {

	builder, builder2 := NewBuilder(), NewBuilder()

	intd := []InTicketData{{Name: "TICKET01", Odate: "ODATE"}}
	outtd := []OutTicketData{{Action: "ADD", Name: "TICKET01", Odate: "ODATE"}}

	task, err := builder.WithBase("testgroup", "testname", "testdescription").WithSchedule(SchedulingData{OrderType: OrderingManual}).
		WithInTicekts(intd, "").
		WithOutTickets(outtd).Build()
	if err != nil {
		t.Error(err)
	}

	expected, _ := builder2.FromTemplate(task).Build()
	if len(expected.InTickets) == 0 {
		t.Error("unexpected result:", 0, "expected:", 1)
	}
	if len(expected.OutTickets) == 0 {
		t.Error("unexpected result:", 0, "expected:", 1)
	}
}

func TestBuilder_WithCyclic(t *testing.T) {

	builder := NewBuilder()

	task, err := builder.WithBase("testgroup", "testname", "testdescription").WithSchedule(SchedulingData{OrderType: OrderingManual}).
		WithCyclic(CyclicTaskData{IsCycle: true, MaxRuns: 10, TimeInterval: 5}).Build()
	if err != nil {
		t.Error(err)
	}

	actual := task.Cyclic
	if actual.IsCycle != true && actual.MaxRuns != 10 && actual.TimeInterval != 5 {
		t.Error("unexpected result:", actual)
	}
}
