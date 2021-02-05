package taskdef

import "testing"

func TestBuilderFromTemplate(t *testing.T) {

	builder, builder2 := &DummyTaskBuilder{}, &DummyTaskBuilder{}

	sd := SchedulingData{FromTime: "10:30", AllowPastSub: true}
	intd := []InTicketData{{Name: "TICKET01", Odate: "ODATE"}}
	outtd := []OutTicketData{{Action: "ADD", Name: "TICKET01", Odate: "ODATE"}}

	task, _ := builder.WithBase("testgroup", "testname", "testdescription").WithRetention(5).WithConfirm().Build()

	task2, _ := builder2.FromTemplate(task).WithSchedule(sd).WithInTicekts(intd, InTicketAND).WithOutTickets(outtd).Build()

	n, g, d := task2.GetInfo()

	if n != "testname" || g != "testgroup" || d != "testdescription" {
		a, b, c := task.GetInfo()
		t.Error("Unexpected values:", n, g, d, "expected:", a, b, c)
	}

	if task.Retention() != 5 {
		t.Error("Unexpected value:", task2.Retention(), "expected:", task.Retention())
	}

	from, to := task2.TimeSpan()
	allow := task2.AllowPast()

	if from != "10:30" || to != "" || allow != true {
		t.Error("Unexpected values, expected:", "10:30", "", true, "actual:", from, to, allow)
	}

	if task2.TicketsIn()[0].Name != intd[0].Name && task2.TicketsIn()[0].Odate != intd[0].Odate {
		t.Error("Unexpected values, expected:", intd[0].Name, intd[0].Odate, "actual:", task2.TicketsIn()[0].Name, task2.TicketsIn()[0].Odate)

	}

	if task2.TicketsOut()[0].Name != outtd[0].Name && task2.TicketsOut()[0].Odate != outtd[0].Odate {
		t.Error("Unexpected values, expected:", outtd[0].Name, outtd[0].Odate, "actual:", task2.TicketsOut()[0].Name, task2.TicketsOut()[0].Odate)

	}

}
