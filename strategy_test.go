package smp

import (
	"encoding/json"
	"testing"

	"github.com/myfantasy/mft"
)

// strategy
const StartegyTestName = "StartegyTest"

type StartegyTest struct {
	A int     `json:"a"`
	B float64 `json:"b"`
}

func (st *StartegyTest) Type() string {
	return StartegyTestName
}
func (st *StartegyTest) Marshal() json.RawMessage {
	b, er0 := json.Marshal(st)
	if er0 != nil {
		panic(er0)
	}

	return b
}
func (st *StartegyTest) Step(s StepParams) {
	st.A++
	st.B--
}

func (st *StartegyTest) Status() Status {
	return nil
}
func (st *StartegyTest) Logs() Logs {
	return nil
}

type TestStrategyStorageStruct struct {
	Stor *StrategyStorage `json:"storage"`
}

func TestStrategyStorage(t *testing.T) {
	DefaultStrategyGenerator.Add(StartegyTestName,
		func() (s Strategy) { return &StartegyTest{} },
		func(data json.RawMessage) (s Strategy, err *mft.Error) {
			s = &StartegyTest{}
			er0 := json.Unmarshal(data, s)
			if er0 != nil {
				return nil, mft.ErrorNew("fail inmarshal", er0)
			}
			return s, nil
		})

	tss := TestStrategyStorageStruct{
		Stor: StrategyStorageCreate(),
	}

	s, err := DefaultStrategyGenerator.Create(StartegyTestName)

	if err != nil {
		t.Fatal(err)
	}

	s.Step(nil)

	tss.Stor.Starategy["abc"] = s

	b, er0 := json.Marshal(tss)

	if er0 != nil {
		t.Fatal(er0)
	}

	tss2 := TestStrategyStorageStruct{}

	er0 = json.Unmarshal(b, &tss2)
	if er0 != nil {
		t.Fatal(er0)
	}

	if len(tss2.Stor.Starategy) != 1 {
		t.Fatalf("Unmarshal should return 1 object (current %v)", len(tss2.Stor.Starategy))
	}

	s2, ok := tss2.Stor.Starategy["abc"]
	if !ok {
		t.Fatalf("Unmarshal should contains `abc` startegy")
	}

	sStruct, ok := s2.(*StartegyTest)
	if !ok {
		t.Fatalf("Unmarshal should `abc` be *StartegyTest startegy")
	}

	if sStruct.A != 1 {
		t.Fatalf("Unmarshal should `abc.A`.(*StartegyTest) == 1 (current %v)", sStruct.A)
	}
}
