package main

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func Test_parseAndAppendStatement(t *testing.T) {
	type args struct {
		reference     string
		accountNumber string
		description   string
		startBalance  string
		mutation      string
		endBalance    string
		statements    []Statement
	}
	var statementsExpectCorrect []Statement
	var statementsExpectDuplicateReference []Statement
	var statementsExpectBadEndBalance []Statement

	var statementsAr []Statement
	startBal, _ := decimal.NewFromString("10.00")
	mutation, _ := decimal.NewFromString("10.00")
	endBal, _ := decimal.NewFromString("20.00")
	badEndBal, _ := decimal.NewFromString("10.00")

	statementsExpectCorrect = append(statementsExpectCorrect, Statement{
		Reference:     1234,
		AccountNumber: "1234",
		Description:   "test",
		StartBalance:  startBal,
		Mutation:      mutation,
		EndBalance:    endBal,
		State:         CORRECT,
	})

	statementsExpectBadEndBalance = append(statementsExpectBadEndBalance, Statement{
		Reference:     1234,
		AccountNumber: "1234",
		Description:   "test",
		StartBalance:  startBal,
		Mutation:      mutation,
		EndBalance:    badEndBal,
		State:         BADENDBALANCE,
	})

	statementsExpectDuplicateReference = append(statementsExpectDuplicateReference, Statement{
		Reference:     1234,
		AccountNumber: "1234",
		Description:   "test",
		StartBalance:  startBal,
		Mutation:      mutation,
		EndBalance:    endBal,
		State:         DUPLICATEREFERENCE,
	})
	statementsExpectDuplicateReference = append(statementsExpectDuplicateReference, Statement{
		Reference:     1234,
		AccountNumber: "1234",
		Description:   "test",
		StartBalance:  startBal,
		Mutation:      mutation,
		EndBalance:    endBal,
		State:         DUPLICATEREFERENCE,
	})

	tests := []struct {
		name string
		args args
		want []Statement
	}{
		// TODO: Add test cases.

		{"Good Result", args{
			reference:     "1234",
			accountNumber: "1234",
			description:   "test",
			startBalance:  "10.00",
			mutation:      "10.00",
			endBalance:    "20.00",
			statements:    statementsAr},
			statementsExpectCorrect},

		{"Duplicate Reference", args{
			reference:     "1234",
			accountNumber: "1234",
			description:   "test",
			startBalance:  "10.00",
			mutation:      "10.00",
			endBalance:    "20.00",
			statements:    statementsExpectCorrect},
			statementsExpectDuplicateReference},

		{"Bad end balance", args{
			reference:     "1234",
			accountNumber: "1234",
			description:   "test",
			startBalance:  "10.00",
			mutation:      "10.00",
			endBalance:    "10.00",
			statements:    statementsAr},
			statementsExpectBadEndBalance},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAndAppendStatement(tt.args.reference, tt.args.accountNumber, tt.args.description, tt.args.startBalance, tt.args.mutation, tt.args.endBalance, tt.args.statements); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAndAppendStatement() = %v, want %v", got, tt.want)
			}
		})
	}
}
