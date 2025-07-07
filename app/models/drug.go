package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Drug struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `bson:"name"`
	Slug             string             `bson:"slug"`
	Type             string             `bson:"type"`
	Category         string             `bson:"category"`
	RegistedNumber   string             `bson:"registed_number"`
	Ingredients      []string           `bson:"ingredients"`
	Excipient        []string           `bson:"excipient,omitempty"`
	Indication       string             `bson:"indication"`
	Contraindication string             `bson:"contraindication"`
	Uses             string             `bson:"uses"`
	Dosage           string             `bson:"dosage"`
	Administration   string             `bson:"administration"`
	SideEffects      string             `bson:"side_effects"`
	DrugInteractions string             `bson:"drug_interactions,omitempty"`
	Manufacturer     string             `bson:"manufacturer"`
	MAH              string             `bson:"mah_ky,omitempty"`
	DosageForm       string             `bson:"dosage_form"`
	Packaging        string             `bson:"packaging,omitempty"`
	Storage          string             `bson:"storage"`
	Notes            string             `bson:"notes"`
	QuaLieu          string             `bson:"qua_lieu,omitempty"`
	Price            int                `bson:"gia,omitempty"`
	Image            string             `bson:"hinh_anh,omitempty"`
	Retailer         string             `bson:"retailer"`
	UpdateDate       primitive.DateTime `bson:"update_date"`
}
