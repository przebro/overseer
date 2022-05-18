package taskdef

import (
	"testing"
)

func TestBuilderFromTemplate(t *testing.T) {

	builder, builder2 := &DummyTaskBuilder{}, &DummyTaskBuilder{}

	sd := SchedulingData{FromTime: "10:30", OrderType: OrderingManual}
	intd := []InTicketData{{Name: "TICKET01", Odate: "ODATE"}}
	outtd := []OutTicketData{{Action: "ADD", Name: "TICKET01", Odate: "ODATE"}}

	task, err := builder.WithBase("testgroup", "testname", "testdescription").WithSchedule(SchedulingData{OrderType: OrderingManual}).WithRetention(5).WithConfirm().Build()
	if err != nil {
		t.Error(err)
	}

	task2, err := builder2.FromTemplate(task).WithSchedule(sd).WithInTicekts(intd, InTicketAND, "").WithOutTickets(outtd).Build()
	if err != nil {
		t.Error(err)
	}

	n, g, d := task2.GetInfo()

	if n != "testname" || g != "testgroup" || d != "testdescription" {
		a, b, c := task.GetInfo()
		t.Error("Unexpected values:", n, g, d, "expected:", a, b, c)
	}

	if task.Retention() != 5 {
		t.Error("Unexpected value:", task2.Retention(), "expected:", task.Retention())
	}

	from, to := task2.TimeSpan()

	if from != "10:30" || to != "" {
		t.Error("Unexpected values, expected:", "10:30", "", true, "actual:", from, to)
	}

	if task2.TicketsIn()[0].Name != intd[0].Name && task2.TicketsIn()[0].Odate != intd[0].Odate {
		t.Error("Unexpected values, expected:", intd[0].Name, intd[0].Odate, "actual:", task2.TicketsIn()[0].Name, task2.TicketsIn()[0].Odate)

	}

	if task2.TicketsOut()[0].Name != outtd[0].Name && task2.TicketsOut()[0].Odate != outtd[0].Odate {
		t.Error("Unexpected values, expected:", outtd[0].Name, outtd[0].Odate, "actual:", task2.TicketsOut()[0].Name, task2.TicketsOut()[0].Odate)

	}

}

func TestBuilderFromTemplate_Tickets(t *testing.T) {

	builder, builder2 := &DummyTaskBuilder{}, &DummyTaskBuilder{}

	intd := []InTicketData{{Name: "TICKET01", Odate: "ODATE"}}
	outtd := []OutTicketData{{Action: "ADD", Name: "TICKET01", Odate: "ODATE"}}

	task, err := builder.WithBase("testgroup", "testname", "testdescription").WithSchedule(SchedulingData{OrderType: OrderingManual}).
		WithInTicekts(intd, InTicketAND, "").
		WithOutTickets(outtd).Build()
	if err != nil {
		t.Error(err)
	}

	expected, _ := builder2.FromTemplate(task).Build()
	if len(expected.TicketsIn()) == 0 {
		t.Error("unexpected result:", 0, "expected:", 1)
	}
	if len(expected.TicketsOut()) == 0 {
		t.Error("unexpected result:", 0, "expected:", 1)
	}
}

func TestBuilder_WithCyclic(t *testing.T) {

	builder := &DummyTaskBuilder{}

	task, err := builder.WithBase("testgroup", "testname", "testdescription").WithSchedule(SchedulingData{OrderType: OrderingManual}).
		WithCyclic(CyclicTaskData{IsCycle: true, MaxRuns: 10, TimeInterval: 5}).Build()
	if err != nil {
		t.Error(err)
	}

	actual := task.Cyclic()
	if actual.IsCycle != true && actual.MaxRuns != 10 && actual.TimeInterval != 5 {
		t.Error("unexpected result:", actual)
	}
}
