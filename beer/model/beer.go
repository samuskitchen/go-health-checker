// Package model defines domain structures and their representations for the Beer entity.
// Contains both the internal beer definition in the database and its HTTP response format.
package model

import "time"

// Beers represents the Beer entity as stored in the database.
// Each field is tagged to map to the corresponding column.
type Beers struct {
	ID        uint      `db:"id"`         // Unique identifier for the beer
	Name      string    `db:"name"`       // Name of the beer
	Brewery   string    `db:"brewery"`    // Producer brewery
	Country   string    `db:"country"`    // Country of origin
	Price     float64   `db:"price"`      // Unit price of the beer in original currency
	Currency  string    `db:"currency"`   // ISO currency code (e.g., USD)
	CreatedAt time.Time `db:"created_at"` // Timestamp of beer creation
	UpdatedAt time.Time `db:"updated_at"` // Timestamp of last beer update
}

// BeersResponse defines the structure sent to the client in HTTP responses.
// Omits empty fields in JSON (omitempty).
type BeersResponse struct {
	ID        uint      `json:"id,omitempty"`         // Unique beer identifier
	Name      string    `json:"name,omitempty"`       // Beer name
	Brewery   string    `json:"brewery,omitempty"`    // Producer brewery
	Country   string    `json:"country,omitempty"`    // Country of origin
	Price     float64   `json:"price,omitempty"`      // Unit price
	Currency  string    `json:"currency,omitempty"`   // Currency code
	CreatedAt time.Time `json:"created_at,omitempty"` // Creation date
	UpdatedAt time.Time `json:"updated_at,omitempty"` // Last update date
}

// ToBeersResponse transforms the internal Beers model to its HTTP response representation.
// Returns a BeersResponse object with publicly exposed fields.
func (b *Beers) ToBeersResponse() BeersResponse {
	return BeersResponse{
		ID:        b.ID,
		Name:      b.Name,
		Brewery:   b.Brewery,
		Country:   b.Country,
		Price:     b.Price,
		Currency:  b.Currency,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
