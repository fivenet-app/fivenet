package testdata

import "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"

func ptr(s string) *string {
	return &s
}

func ptrFloat(f float64) *float64 {
	return &f
}

var Users = []*users.User{
	{
		UserId:      1,
		Identifier:  "char1:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4",
		Job:         "ambulance",
		JobGrade:    17,
		Firstname:   "Dr. Amy",
		Lastname:    "Clockwork",
		Sex:         ptr("f"),
		Dateofbirth: "08.04.2003",
		Height:      ptrFloat(182),
	},
	{
		UserId:      2,
		Identifier:  "char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57",
		Job:         "ambulance",
		JobGrade:    20,
		Firstname:   "Philipp",
		Lastname:    "Scott",
		Sex:         ptr("m"),
		Dateofbirth: "01.08.1982",
		Height:      ptrFloat(185),
	},
	{
		UserId:      3,
		Identifier:  "char1:db7e039146d5bf1b6781e7bc1bef31f0bb1298ea",
		Job:         "doj",
		JobGrade:    16,
		Firstname:   "Jonas",
		Lastname:    "Striker",
		Sex:         ptr("m"),
		Dateofbirth: "28.10.1990",
		Height:      ptrFloat(186),
	},
	{
		UserId:      4,
		Identifier:  "char2:fcee377a1fda007a8d2cc764a0a272e04d8c5d57",
		Job:         "police",
		JobGrade:    2,
		Firstname:   "Hannibal",
		Lastname:    "Scott",
		Sex:         ptr("m"),
		Dateofbirth: "15.06.1990",
		Height:      ptrFloat(180),
	},
	{
		UserId:      5,
		Identifier:  "char2:3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4",
		Job:         "unemployed",
		JobGrade:    1,
		Firstname:   "Peter",
		Lastname:    "Hans",
		Sex:         ptr("m"),
		Dateofbirth: "10.02.1991",
		Height:      ptrFloat(178),
	},
}
