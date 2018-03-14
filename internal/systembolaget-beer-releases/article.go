package search

import "time"

type Article struct {
	Nr                    int       `json:"Nr"`
	ArticleID             int       `json:"ArticleID"`
	ArticleNumber         int       `json:"ArticleNumber"`
	Name                  string    `json:"Name"`
	SecondaryName         string    `json:"SecondaryName"`
	PriceIncludingVAT     float32   `json:"PriceIncludingVAT"`
	VolumeInMl            float32   `json:"VolumeInMl"`
	PricePerLitre         float32   `json:"PricePerLitre"`
	SalesStart            time.Time `json:"SalesStart"`
	Expired               bool      `json:"Expired"`
	ArticleGroup          string    `json:"ArticleGroup"`
	ArticleType           string    `json:"ArticleType"`
	ArticleStyle          string    `json:"ArticleStyle"`
	Packaging             string    `json:"Packaging"`
	Seal                  string    `json:"Seal"`
	Origin                string    `json:"Origin"`
	OriginCountry         string    `json:"OriginCountry"`
	Producer              string    `json:"Producer"`
	Supplier              string    `json:"Supplier"`
	Vintage               string    `json:"Vintage"`
	AlcoholPercentage     float64   `json:"AlcoholPercentage"`
	Selection             string    `json:"Selection"`
	SelectionText         string    `json:"SelectionText"`
	Organic               bool      `json:"Organic"`
	Ethical               bool      `json:"Ethical"`
	Koscher               bool      `json:"Koscher"`
	IngredientDescription string    `json:"IngredientDescription"`
}
