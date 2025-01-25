package dtos

type ActivityTypeDTO struct{
	Name	string
	Calory	int
}

var ActivityValues = []ActivityTypeDTO{
    {Name: "Walking", Calory: 4},
    {Name: "Yoga", Calory: 1},
    {Name: "Stretching", Calory: 1},
    {Name: "Cycling", Calory: 1},
    {Name: "Swimming", Calory: 1},
    {Name: "Dancing", Calory: 1},
    {Name: "Hiking", Calory: 1},
    {Name: "Running", Calory: 1},
    {Name: "HIIT", Calory: 1},
    {Name: "JumpRope", Calory: 1},
}

func GetCaloryByActivityName(name string) (int, bool) {
    for _, activity := range ActivityValues {
        if activity.Name == name {
            return activity.Calory, true
        }
    }
    return 0, false
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}