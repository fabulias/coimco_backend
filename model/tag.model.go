package model

//This function allow insert tag' resource
func InsertTag(in *Tag) (*Tag, bool) {
	err = dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
