package beer

const (
	errNameIsRequired = "name is required"
	errInvalidType    = "invalid beer type"
	errInvalidStyle   = "invalid beer style"
)

type Beer struct {
	ID    int64     `json:"id,omitempty"`
	Name  string    `json:"name"`
	Type  BeerType  `json:"type"`
	Style BeerStyle `json:"style"`
}

func (b *Beer) Validate() ([]string, bool) {
	var validationErrors []string

	if b.Name == "" {
		validationErrors = append(validationErrors, errNameIsRequired)
	}

	if b.Type.String() == "Unknown" {
		validationErrors = append(validationErrors, errInvalidType)
	}

	if b.Style.String() == "Unknown" {
		validationErrors = append(validationErrors, errInvalidStyle)
	}

	if len(validationErrors) > 0 {
		return validationErrors, true
	}

	return nil, false
}

// https://www.thebeerstore.ca/beer-101/beer-types/
type BeerType int

const (
	TypeAle   = 1
	TypeLager = 2
	TypeMalt  = 3
	TypeStout = 4
)

func (t BeerType) String() string {
	switch t {
	case TypeAle:
		return "Ale"
	case TypeLager:
		return "Lager"
	case TypeMalt:
		return "Malt"
	case TypeStout:
		return "Stout"
	}
	return "Unknown"
}

type BeerStyle int

const (
	StyleAmber = iota + 1
	StyleBlonde
	StyleBrown
	StyleCream
	StyleDark
	StylePale
	StyleStrong
	StyleWheat
	StyleRed
	StyleIPA
	StyleLime
	StylePilsner
	StyleGolden
	StyleFruit
	StyleHoney
)

func (s BeerStyle) String() string {
	switch s {
	case StyleAmber:
		return "Amber"
	case StyleBlonde:
		return "Blonde"
	case StyleBrown:
		return "Brown"
	case StyleCream:
		return "Cream"
	case StyleDark:
		return "Dark"
	case StylePale:
		return "Pale"
	case StyleStrong:
		return "Strong"
	case StyleWheat:
		return "Wheat"
	case StyleRed:
		return "Red"
	case StyleIPA:
		return "IPA"
	case StyleLime:
		return "Lime"
	case StylePilsner:
		return "Pilsner"
	case StyleGolden:
		return "Golden"
	case StyleFruit:
		return "Fruit"
	case StyleHoney:
		return "Honey"
	}
	return "Unknown"
}
