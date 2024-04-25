package tax

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type ReqTax struct {
	TotalIncome float64 `json:"totalIncome"`
	Wht         float64 `json:"wht"`
	Allowances  []Allowance
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type ResTaxLevel struct {
	Tax       float64 `json:"tax"`
	TaxRefund float64 `json:"taxRefund"`
	TaxLevel  []TaxLevel
}

type ReqAmount struct {
	Amount float64 `json:"amount"`
}

type ResPersonalDeduction struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type ResCsvTax struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
	TaxRefund   float64 `json:"taxRefund"`
}

type ResAllCsv struct {
	Taxes []ResCsvTax `json:"taxes"`
}

type DB struct {
	ID             int     `postgres:"id"`
	Minimum_salary float64 `postgres:"minimum_salary"`
	Maximum_salary float64 `postgres:"maximum_salary|NULL"`
	Rate           float64 `postgres:"rate"`
	Created_at     string  `postgres:"created_at"`
}

type DbDeduction struct {
	ID             int     `postgres:"id"`
	Type           string  `postgres:"deducation_type"`
	Minimum_amount float64 `postgres:"minimum_amount"`
	Maximum_amount float64 `postgres:"maximum_amount"`
	Amount         float64 `postgres:"amount"`
	Created_at     string  `postgres:"created_at"`
	Updated_at     string  `postgres:"updated_at"`
}

type Tax struct {
	info InfoTax
}

type Err struct {
	Message string `json:"message"`
}

type InfoTax interface {
	GetTax() ([]DB, error)
	GetTaxDeducationByType(deducation_type string) (DbDeduction, error)
	SetTaxDeducationByType(deducation_type string, amount float64) error
}

func New(info InfoTax) Tax {
	return Tax{
		info: info,
	}
}
