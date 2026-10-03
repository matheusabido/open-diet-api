package model

import "slices"

type Sex string

const (
	SexMale   Sex = "male"
	SexFemale Sex = "female"
)

var Sexes = []Sex{SexMale, SexFemale}

func (s Sex) Valid() bool {
	return slices.Contains(Sexes, s)
}

type ActivityLevel string

const (
	ActivityLevelSedentary ActivityLevel = "sedentary"
	ActivityLevelLight     ActivityLevel = "light"
	ActivityLevelModerate  ActivityLevel = "moderate"
	ActivityLevelHeavy     ActivityLevel = "heavy"
	ActivityLevelAthlete   ActivityLevel = "athlete"
)

var ActivityLevels = []ActivityLevel{
	ActivityLevelSedentary,
	ActivityLevelLight,
	ActivityLevelModerate,
	ActivityLevelHeavy,
	ActivityLevelAthlete,
}

func (a ActivityLevel) Valid() bool {
	return slices.Contains(ActivityLevels, a)
}

type Objective string

const (
	ObjectiveCut      Objective = "cut"
	ObjectiveMaintain Objective = "maintain"
	ObjectiveBulk     Objective = "bulk"
)

var Objectives = []Objective{ObjectiveCut, ObjectiveMaintain, ObjectiveBulk}

func (o Objective) Valid() bool {
	return slices.Contains(Objectives, o)
}
