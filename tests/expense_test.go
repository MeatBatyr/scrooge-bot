package tests

import (
	"scroogebot/expenditure"
	"testing"
)

func TestMakeCalculationResult_many_records_for_category(t *testing.T) {
	in := []expenditure.Expense{{Category: "one", Amount: 5.5}, {Category: "one", Amount: 4.5}}
	want := "one: 10"

	got := expenditure.MakeCalculationResult(in)

	if got != want {
		t.Errorf("MakeCalculationResult(%v) == %q, want %q", in, got, want)
	}
}

func TestMakeCalculationResult_many_categories(t *testing.T) {
	in := []expenditure.Expense{{Category: "one", Amount: 5.5}, {Category: "two", Amount: 4.5}}
	want := "one: 5.5\ntwo: 4.5"

	got := expenditure.MakeCalculationResult(in)

	if got != want {
		t.Errorf("MakeCalculationResult(%v) == %q, want %q", in, got, want)
	}
}

func TestMakeCalculationResult_no_any_records(t *testing.T) {
	in := []expenditure.Expense{}

	got := expenditure.MakeCalculationResult(in)

	if got != expenditure.NothingForCalc {
		t.Errorf("MakeCalculationResult(%v) == %q, want %q", in, got, expenditure.NothingForCalc)
	}
}

func TestParse_correct_maney_value(t *testing.T) {
	cases := []struct {
		in   string
		want expenditure.Expense
	}{
		{in: "42 #one", want: expenditure.Expense{Category: "one", Amount: 42}},
		{in: "4.2 #one", want: expenditure.Expense{Category: "one", Amount: 4.2}},
		{in: "4.20 #one", want: expenditure.Expense{Category: "one", Amount: 4.2}},
	}

	for _, v := range cases {
		got, err := expenditure.Parse(v.in)

		if err != nil || got.Amount != v.want.Amount || got.Category != "one" {
			t.Errorf("Parse(%q): got Amaunt %v; Category:%v, want %v; category: one", v.in, got.Amount, got.Category, v.want.Amount)
		}
	}
}

func TestParse_no_category(t *testing.T) {
	in := "345 one"
	got, err := expenditure.Parse(in)
	if err != nil || got != nil {
		t.Errorf("Parse(%q): got Amaunt %v; Category:%v, want nil", in, got.Amount, got.Category)
	}
}

func TestParse_many_categories(t *testing.T) {
	in := "345 #one #two"
	got, err := expenditure.Parse(in)

	if err.Error() != expenditure.ErrManyCategories.Error() || got != nil {
		t.Errorf("Parse(%q): got Amaunt %v; Category:%v, want nil", in, got.Amount, got.Category)
	}
}

func TestParse_no_amaunt(t *testing.T) {
	in := "#one "
	got, err := expenditure.Parse(in)
	if err != nil || got != nil {
		t.Errorf("Parse(%q): got Amaunt %v; Category:%v, want nil", in, got.Amount, got.Category)
	}
}

func TestParse_many_amaunts(t *testing.T) {
	in := "#one 4.5 5.5"
	want := 10
	got, err := expenditure.Parse(in)
	if err != nil || got.Amount != float64(want) {
		t.Errorf("Parse(%q): got Amaunt %v; want %v", in, got.Amount, want)
	}
}
