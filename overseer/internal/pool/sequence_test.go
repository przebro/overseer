package pool

// func TestSequence(t *testing.T) {

// 	var seq SequenceGenerator
// 	var err error
// 	if _, err = NewSequenceGenerator(provider); err == nil {
// 		t.Error("unexpected result:", err)
// 	}

// 	if seq, err = NewSequenceGenerator(provider); err != nil {
// 		t.Error("unexpected result:", err)
// 	}

// 	seqgen := seq.(*sequenceGenerator)

// 	value := seqgen.Next()
// 	fmt.Println(value)

// 	cHi := seqgen.hi
// 	seqgen.lo = 3843
// 	seq.Next()
// 	if seqgen.lo != 0 && seqgen.hi != cHi+1 {
// 		t.Error("unexpected result")
// 	}

// 	seqgen.hi = seqMax
// 	seqgen.lo = seqMax
// 	seq.Next()
// 	if seqgen.lo != 0 && seqgen.hi != 0 {
// 		t.Error("unexpected result")
// 	}

// }

// func TestInvalidValues(t *testing.T) {

// 	var err error

// 	col, _ := provider.GetCollection(context.Background(), "sequence")
// 	model := sequenceModel{Key: sequenceName, Value: -1}
// 	col.Update(context.Background(), &model)

// 	if _, err = NewSequenceGenerator(provider); err == nil {
// 		t.Error("unexpected result:", err)
// 	}

// 	model = sequenceModel{Key: sequenceName, Value: 3844}
// 	col.Update(context.Background(), &model)

// 	if _, err = NewSequenceGenerator(provider); err == nil {
// 		t.Error("unexpected result:", err)
// 	}
// }
