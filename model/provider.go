package model

type Provider struct {
	Agent
}

type ProviderRankK struct {
	Name string
	Days float64
}

type ProviderRankPP struct {
	Name  string
	Price uint
}
