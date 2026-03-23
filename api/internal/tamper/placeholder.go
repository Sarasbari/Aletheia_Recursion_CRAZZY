package tamper

type Region struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func PlaceholderRegions(hashMatch, metadataValid bool) []Region {
	if hashMatch && metadataValid {
		return []Region{}
	}
	return []Region{
		{X: 10, Y: 20},
		{X: 11, Y: 20},
		{X: 12, Y: 21},
	}
}
