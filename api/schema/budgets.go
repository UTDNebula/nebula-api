package schema

// Must also be updated in api-tools/parser/budgets.go

// Table utils
type Table[T any] struct {
	Name  string   `bson:"name" json:"name"`
	Rows  []Row[T] `bson:"rows" json:"rows"`
	Total T        `bson:"total" json:"total"`
}
type Table2[T any] struct {
	Name  string     `bson:"name" json:"name"`
	Rows  []Table[T] `bson:"rows" json:"rows"`
	Total T          `bson:"total" json:"total"`
}
type Row[T any] struct {
	Label string `bson:"label" json:"label"`
	Value T      `bson:"value" json:"value"`
}

// Top level container for both the planned budget and the actual annual financial report
type Budget struct {
	Id                    string                `bson:"_id" json:"_id"`
	OperatingBudget       OperatingBudget       `bson:"operating_budget" json:"operating_budget"`
	AnnualFinancialReport AnnualFinancialReport `bson:"annual_financial_report" json:"annual_financial_report"`
	Notes                 string                `bson:"notes" json:"notes"`
}

// Operating Budget Structs
type OperatingBudget struct {
	OperatingRevenues             *Table[float64]                  `bson:"operating_revenues" json:"operating_revenues"`
	OperatingExpenses             *Table[float64]                  `bson:"operating_expenses" json:"operating_expenses"`
	BudgetedNonoperatingRevenues  *Table[float64]                  `bson:"budgeted_nonoperating_revenues" json:"budgeted_nonoperating_revenues"`
	BudgetedTuitionAndStudentFees *Table2[float64]                 `bson:"budgeted_tuition_and_student_fees" json:"budgeted_tuition_and_student_fees"`
	BudgetBySchool                *Table2[float64]                 `bson:"budget_by_school" json:"budget_by_school"`
	AuxiliaryExpenses             *Table2[AuxiliaryExpensesValues] `bson:"auxiliary_expenses" json:"auxiliary_expenses"`
}
type AuxiliaryExpensesValues struct {
	EstimatedIncome  float64 `bson:"estimated_income" json:"estimated_income"`
	BudgetedExpenses float64 `bson:"budgeted_expenses" json:"budgeted_expenses"`
	DebtService      float64 `bson:"debt_service" json:"debt_service"`
	Other            float64 `bson:"other" json:"other"`
	ExcessIncome     float64 `bson:"excess_income" json:"excess_income"`
	BeginningBalance float64 `bson:"beginning_balance" json:"beginning_balance"`
	EndingBalance    float64 `bson:"ending_balance" json:"ending_balance"`
}

// Annual Financial Report (AFR) Structs
type AnnualFinancialReport struct {
	OperatingRevenues            *Table[float64] `bson:"operating_revenues" json:"operating_revenues"`
	OperatingExpenses            *Table[float64] `bson:"operating_expenses" json:"operating_expenses"`
	BudgetedNonoperatingRevenues *Table[float64] `bson:"budgeted_nonoperating_revenues" json:"budgeted_nonoperating_revenues"`
	BeginningNetPosition         *float64        `bson:"beginning_net_position" json:"beginning_net_position"`
	EndingNetPosition            *float64        `bson:"ending_net_position" json:"ending_net_position"`
}
