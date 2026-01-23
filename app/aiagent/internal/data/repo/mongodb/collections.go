package mongodb

import (
	"store/pkg/clients/mgz"
)

type Collections struct {
	Resource     *ResourceCollection
	Session      *SessionCollection
	Question     *QuestionCollection
	Answer       *AnswerCollection
	Profile      *ProfileCollection
	Settings     *SettingsCollection
	Item         *ItemCollection
	Survey       *SurveyCollection
	PromptConfig *PromptConfigCollection
	Account      *AccountCollection
}

func NewCollections(config mgz.Config) *Collections {

	database, err := mgz.Database(config)

	if err != nil {
		panic(err)
	}

	return &Collections{
		Resource:     NewResourceCollection(database),
		Session:      NewSessionCollection(database),
		Question:     NewQuestionCollection(database),
		Answer:       NewAnswerCollection(database),
		Profile:      NewProfileCollection(database),
		Settings:     NewSettingsCollection(database),
		Item:         NewItemCollection(database),
		Survey:       NewSurveyCollection(database),
		PromptConfig: NewPromptConfigCollection(database),
		Account:      NewAccountCollection(database),
	}
}
