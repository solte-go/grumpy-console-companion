package model

type Thoughts struct {
	Phrase string `bson:"phrase"`
	//TODO answer checker
	Continuable    bool           `bson:"continuable"`
	PossibleAnswer []string       `bson:"answer"`
	Continuation   []Continuation `bson:"continuation"`
	DefaultEndings []string       `bson:"endings"`
}
type Continuation struct {
	March   string `bson:"march"`
	Respond string `bson:"respond"`
}
