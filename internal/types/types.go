package types

type Form struct {
	Fio      string `json:"Fio"`
	Tel      string `json:"Tel"`
	Email    string `json:"Email"`
	Date     string `json:"Date"`
	Gender   string `json:"Gender"`
	Favlangs []int  `json:"Favlangs"`
	Bio      string `json:"Bio"`
}

type FormErrors struct {
	Fio      string `json:"Fio"`
	Tel      string `json:"Tel"`
	Email    string `json:"Email"`
	Date     string `json:"Date"`
	Gender   string `json:"Gender"`
	Favlangs string `json:"Favlangs"`
	Bio      string `json:"Bio"`
	Familiar string `json:"Familiar"`
}
