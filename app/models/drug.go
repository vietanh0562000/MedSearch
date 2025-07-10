package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Drug struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Name             string             `bson:"name"`
	Slug             string             `bson:"slug"`
	Description      string             `bson:"description"`
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
	MAH              string             `bson:"mah,omitempty"`
	DosageForm       string             `bson:"dosage_form"`
	Packaging        string             `bson:"packaging,omitempty"`
	Storage          string             `bson:"storage"`
	Notes            string             `bson:"notes"`
	Overdose         string             `bson:"overdose,omitempty"`
	Price            string             `bson:"price,omitempty"`
	Image            string             `bson:"image,omitempty"`
	Retailer         string             `bson:"retailer"`
	UpdateDate       primitive.DateTime `bson:"update_date"`
}
