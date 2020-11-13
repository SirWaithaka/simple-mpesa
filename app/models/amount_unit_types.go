package models

// Shillings is the unit of money used in Kenyan Currency (KES). A shilling
// has 100 cents
type Shillings uint

// ToCents is a utility method that returns the amount as Cents
func (amt Shillings) ToCents() Cents {
	return Cents(amt * 100)
}

// Cents is a type representing a small unit of money. 100 cents
// make up a shilling.
type Cents uint

// ToShillings is a utility method that returns the amount as shillings
func (amt Cents) ToShillings() Shillings {
	return Shillings(amt / 100)
}
