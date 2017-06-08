package model

//This function allow insert tag' resource
func InsertTag(in *Tag) (*Tag, bool) {
	err = Dbmap.Create(in).Error
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}

func GetTag(id uint) (Tag, error) {
	var tag Tag
	tag.ID = id
	err = Dbmap.First(&tag, tag.ID).Error
	checkErr(err, selectOneFailed)
	return tag, err
}
