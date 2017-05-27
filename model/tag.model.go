package model

import "strconv"

//This function allow obtain tags' resource.
func GetTag(limit, offset string) ([]Tag, string) {
	var tags []Tag
	var count int64
	//Here obtain total length of table.
	count, err = dbmap.SelectInt("select count(*) from tag")
	checkErr(err, countFailed)
	//Here obtain the tags previously selected.
	_, err = dbmap.Select(&tags, "select * from tag limit $1 offset $2", limit, offset)
	checkErr(err, selectFailed)
	return tags, strconv.Itoa(int(count))
}

//This function allow obtain tag' resource for his id.
func GetTag(rut string) (Tag, error) {
	var tag Tag
	tag.Rut = rut
	err := dbmap.SelectOne(&tag, "select * from tag where rut=$1", tag.Rut)
	checkErr(err, selectOneFailed)
	return tag, err
}

//This function allow insert tag' resource
func InsertTag(in *Tag) (*Tag, bool) {
	err = dbmap.Insert(in)
	if err != nil {
		return in, false
	} else {
		return in, true
	}
}
