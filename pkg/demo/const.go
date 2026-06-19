package demo

type demoSeedJob struct {
	Name  string
	Label string
}

type demoSeedJobGrade struct {
	JobName string
	Grade   int32
	Label   string
}

type demoSeedLicense struct {
	Type  string
	Label string
}

type demoSeedCentrumUnit struct {
	Name        string
	Initials    string
	Color       string
	Icon        string
	Description string
	Attributes  *string
	HomePostal  *string
}

type demoSeedLawbook struct {
	ID          int32
	Name        string
	Description string
}

type demoSeedLaw struct {
	ID            int32
	LawbookID     int32
	Name          string
	Description   string
	Hint          *string
	Fine          int32
	DetentionTime int32
	StvoPoints    int32
}

const (
	unemployedJob = "unemployed"
	PoliceJob     = "police"

	ambulanceJob = "ambulance"
	cafeJob      = "cafe"
	dojJob       = "doj"
	fibJob       = "fib"
	mechanicJob  = "mechanic"
	yardiesJob   = "yardies"

	demoLicenseDrive = "drive"

	demoCentrumUnitColor       = "#0008f0"
	demoCentrumUnitIcon        = "MapMarkerIcon"
	demoCentrumUnitDescription = "Streife"

	demoLawDescriptionPrisonSentenceRequired = "Prison sentence required"
	demoLawDescriptionFinePossible           = "Fine possible"
	demoLawDescriptionFine                   = "Fine"
)

var (
	unemployedJobName = unemployedJob

	demoSeedJobs = []demoSeedJob{
		{Name: ambulanceJob, Label: "LSMD"},
		{Name: cafeJob, Label: "Cat Café"},
		{Name: dojJob, Label: "DOJ"},
		{Name: fibJob, Label: "FIB"},
		{Name: mechanicJob, Label: "Nagata Performance"},
		{Name: PoliceJob, Label: "LSPD"},
		{Name: "prisoner", Label: "Inmate"},
		{Name: unemployedJob, Label: "Unemployed"},
		{Name: yardiesJob, Label: "Yardies"},
	}

	demoSeedJobGrades = []demoSeedJobGrade{
		{JobName: ambulanceJob, Grade: 1, Label: "Trainee Paramedic"},
		{JobName: ambulanceJob, Grade: 2, Label: "Paramedic"},
		{JobName: ambulanceJob, Grade: 3, Label: "Emergency Medical Technician"},
		{JobName: ambulanceJob, Grade: 4, Label: "Emergency Medical Assistant"},
		{JobName: ambulanceJob, Grade: 5, Label: "Emergency Medical Specialist"},
		{JobName: ambulanceJob, Grade: 6, Label: "Medical Student"},
		{JobName: ambulanceJob, Grade: 7, Label: "Assistant Doctor"},
		{JobName: ambulanceJob, Grade: 8, Label: "Emergency Specialist"},
		{JobName: ambulanceJob, Grade: 9, Label: "Senior Assistant Doctor"},
		{JobName: ambulanceJob, Grade: 10, Label: "Specialist Doctor"},
		{JobName: ambulanceJob, Grade: 11, Label: "Rescue Specialist"},
		{JobName: ambulanceJob, Grade: 12, Label: "Experienced Specialist Doctor"},
		{JobName: ambulanceJob, Grade: 13, Label: "Technical Rescue Specialist"},
		{JobName: ambulanceJob, Grade: 14, Label: "Senior Doctor"},
		{JobName: ambulanceJob, Grade: 15, Label: "Chief Senior Doctor"},
		{JobName: ambulanceJob, Grade: 16, Label: "Deputy Chief Doctor"},
		{JobName: ambulanceJob, Grade: 17, Label: "Chief Doctor"},
		{JobName: ambulanceJob, Grade: 18, Label: "Deputy Medical Director"},
		{JobName: ambulanceJob, Grade: 19, Label: "Medical Director"},
		{JobName: cafeJob, Grade: 1, Label: "Intern"},
		{JobName: cafeJob, Grade: 2, Label: "Apprentice"},
		{JobName: cafeJob, Grade: 3, Label: "Barista"},
		{JobName: cafeJob, Grade: 4, Label: "Waiter"},
		{JobName: cafeJob, Grade: 5, Label: "Coffee Art Creator"},
		{JobName: cafeJob, Grade: 6, Label: "Bartender"},
		{JobName: cafeJob, Grade: 7, Label: "Bar Manager"},
		{JobName: cafeJob, Grade: 8, Label: "Security Manager"},
		{JobName: cafeJob, Grade: 9, Label: "HR Manager"},
		{JobName: cafeJob, Grade: 10, Label: "Manager"},
		{JobName: cafeJob, Grade: 11, Label: "Deputy Head"},
		{JobName: cafeJob, Grade: 12, Label: "Head"},
		{JobName: dojJob, Grade: 1, Label: "Office Assistant"},
		{JobName: dojJob, Grade: 2, Label: "Clerk"},
		{JobName: dojJob, Grade: 3, Label: "Prosecutor"},
		{JobName: dojJob, Grade: 4, Label: "Senior Prosecutor"},
		{JobName: dojJob, Grade: 5, Label: "Probationary Judge"},
		{JobName: dojJob, Grade: 6, Label: "City Clerk"},
		{JobName: dojJob, Grade: 7, Label: "County Clerk"},
		{JobName: dojJob, Grade: 8, Label: "Assistant District Attorney"},
		{JobName: dojJob, Grade: 9, Label: "District Attorney"},
		{JobName: dojJob, Grade: 10, Label: "State Attorney"},
		{JobName: dojJob, Grade: 11, Label: "Associate Judge"},
		{JobName: dojJob, Grade: 12, Label: "Judge"},
		{JobName: dojJob, Grade: 13, Label: "Senior Judge"},
		{JobName: dojJob, Grade: 14, Label: "United States Attorney"},
		{JobName: dojJob, Grade: 15, Label: "Head Clerk"},
		{JobName: dojJob, Grade: 16, Label: "Deputy Chief Judge"},
		{JobName: dojJob, Grade: 17, Label: "Deputy Attorney General"},
		{JobName: dojJob, Grade: 18, Label: "Chief Judge"},
		{JobName: dojJob, Grade: 19, Label: "Attorney General"},
		{JobName: fibJob, Grade: 1, Label: "Trainee"},
		{JobName: fibJob, Grade: 2, Label: "Junior Agent"},
		{JobName: fibJob, Grade: 3, Label: "Agent"},
		{JobName: fibJob, Grade: 4, Label: "Senior Agent"},
		{JobName: fibJob, Grade: 5, Label: "Special Agent"},
		{JobName: fibJob, Grade: 6, Label: "First Special Agent"},
		{JobName: fibJob, Grade: 7, Label: "Supervisory Special Agent"},
		{JobName: fibJob, Grade: 8, Label: "Deputy Section Chief"},
		{JobName: fibJob, Grade: 9, Label: "Section Chief"},
		{JobName: fibJob, Grade: 10, Label: "Assistant Director"},
		{JobName: fibJob, Grade: 11, Label: "Deputy Director"},
		{JobName: fibJob, Grade: 12, Label: "Director"},
		{JobName: mechanicJob, Grade: 1, Label: "Intern"},
		{JobName: mechanicJob, Grade: 2, Label: "Apprentice"},
		{JobName: mechanicJob, Grade: 3, Label: "Employee"},
		{JobName: mechanicJob, Grade: 4, Label: "Journeyman"},
		{JobName: mechanicJob, Grade: 5, Label: "Senior Journeyman"},
		{JobName: mechanicJob, Grade: 6, Label: "Master Mechanic"},
		{JobName: mechanicJob, Grade: 7, Label: "Tuning Expert"},
		{JobName: mechanicJob, Grade: 8, Label: "Workshop Manager"},
		{JobName: mechanicJob, Grade: 9, Label: "HR Manager"},
		{JobName: mechanicJob, Grade: 10, Label: "Department Head"},
		{JobName: mechanicJob, Grade: 11, Label: "Deputy CEO"},
		{JobName: mechanicJob, Grade: 12, Label: "CEO"},
		{JobName: PoliceJob, Grade: 1, Label: "Cadet"},
		{JobName: PoliceJob, Grade: 2, Label: "Rookie"},
		{JobName: PoliceJob, Grade: 3, Label: "Officer"},
		{JobName: PoliceJob, Grade: 4, Label: "Officer 2"},
		{JobName: PoliceJob, Grade: 5, Label: "Officer 3"},
		{JobName: PoliceJob, Grade: 6, Label: "Senior Officer"},
		{JobName: PoliceJob, Grade: 7, Label: "Sergeant"},
		{JobName: PoliceJob, Grade: 8, Label: "Sergeant 2"},
		{JobName: PoliceJob, Grade: 9, Label: "Staff Sergeant"},
		{JobName: PoliceJob, Grade: 10, Label: "Lieutenant"},
		{JobName: PoliceJob, Grade: 11, Label: "Captain"},
		{JobName: PoliceJob, Grade: 12, Label: "Detective Trainee"},
		{JobName: PoliceJob, Grade: 13, Label: "Junior Detective"},
		{JobName: PoliceJob, Grade: 14, Label: "Detective"},
		{JobName: PoliceJob, Grade: 15, Label: "Senior Detective"},
		{JobName: PoliceJob, Grade: 16, Label: "Commander"},
		{JobName: PoliceJob, Grade: 17, Label: "Deputy Chief of Police"},
		{JobName: PoliceJob, Grade: 18, Label: "Assistant Chief of Police"},
		{JobName: PoliceJob, Grade: 19, Label: "Chief of Police"},
		{JobName: "prisoner", Grade: 1, Label: "Prisoner"},
		{JobName: unemployedJob, Grade: 1, Label: "Unemployed"},
		{JobName: yardiesJob, Grade: 1, Label: "Runner"},
		{JobName: yardiesJob, Grade: 2, Label: "Dealer"},
		{JobName: yardiesJob, Grade: 3, Label: "Bouncer"},
		{JobName: yardiesJob, Grade: 4, Label: "Homie"},
		{JobName: yardiesJob, Grade: 5, Label: "Hustler"},
		{JobName: yardiesJob, Grade: 6, Label: "Real Yardie"},
		{JobName: yardiesJob, Grade: 7, Label: "Original Yardie"},
		{JobName: yardiesJob, Grade: 8, Label: "Shotcaller"},
		{JobName: yardiesJob, Grade: 9, Label: "Master Yardie"},
		{JobName: yardiesJob, Grade: 10, Label: "Hood Watcher"},
		{JobName: yardiesJob, Grade: 11, Label: "Vice Hood Master"},
		{JobName: yardiesJob, Grade: 12, Label: "O.G."},
	}

	demoSeedCentrumUnits = []demoSeedCentrumUnit{
		{
			Name:        "Adam 1-11",
			Initials:    "1-11",
			Color:       demoCentrumUnitColor,
			Icon:        demoCentrumUnitIcon,
			Description: demoCentrumUnitDescription,
		},
		{
			Name:        "Adam 1-12",
			Initials:    "1-12",
			Color:       demoCentrumUnitColor,
			Icon:        demoCentrumUnitIcon,
			Description: demoCentrumUnitDescription,
		},
		{
			Name:        "Adam 1-13",
			Initials:    "1-13",
			Color:       demoCentrumUnitColor,
			Icon:        demoCentrumUnitIcon,
			Description: demoCentrumUnitDescription,
		},
		{
			Name:        "Supervisor",
			Initials:    "SUP",
			Color:       "#800000",
			Icon:        "MapMarkerAccountIcon",
			Description: "Supervisor Unit",
		},
	}

	demoSeedLicenses = []demoSeedLicense{
		{Type: "aircraft", Label: "Flugschein"},
		{Type: "boat", Label: "Boots-Führerschein"},
		{Type: "commercial", Label: "Rohstofflizenz"},
		{Type: "dmv", Label: "Theoretische Fahrprüfung"},
		{Type: demoLicenseDrive, Label: "PKW-Führerschein"},
		{Type: "drive_bike", Label: "Motorrad-Führerschein"},
		{Type: "drive_truck", Label: "LKW-Führerschein"},
		{Type: "weapon", Label: "Waffenschein"},
	}

	demoSeedLawbooks = []demoSeedLawbook{
		{ID: 1, Name: "StGB", Description: "Criminal Code"},
		{ID: 2, Name: "WaffG", Description: "Weapons Act"},
		{ID: 3, Name: "BtMG", Description: "Narcotics Act"},
		{ID: 4, Name: "LuftVO", Description: "Aviation Regulation"},
		{ID: 5, Name: "StVO", Description: "Road Traffic Regulations"},
		{ID: 6, Name: "GewO", Description: "Trade Regulation"},
		{ID: 7, Name: "WirtG", Description: "Economic Law"},
	}

	// Demo law subset based on demo-db/data/01-fivenet-laws.sql.
	// Redundant StVO speed-band rows were intentionally reduced to keep seed data compact.
	demoSeedLaws = []demoSeedLaw{
		{
			ID:            1,
			LawbookID:     1,
			Name:          "§12 Murder",
			Description:   demoLawDescriptionPrisonSentenceRequired + ", not less than 45 detention units!",
			Fine:          0,
			DetentionTime: 45,
			StvoPoints:    0,
		},
		{
			ID:            2,
			LawbookID:     1,
			Name:          "§13 Manslaughter",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 20,
			StvoPoints:    0,
		},
		{
			ID:            4,
			LawbookID:     1,
			Name:          "§14 Bodily Harm",
			Description:   demoLawDescriptionFinePossible,
			Fine:          5000,
			DetentionTime: 15,
			StvoPoints:    0,
		},
		{
			ID:            6,
			LawbookID:     1,
			Name:          "§15 Dangerous Bodily Harm",
			Description:   demoLawDescriptionFinePossible,
			Fine:          20000,
			DetentionTime: 30,
			StvoPoints:    0,
		},
		{
			ID:            10,
			LawbookID:     1,
			Name:          "§19 Theft",
			Description:   demoLawDescriptionFinePossible + ", based on the value of the stolen goods",
			Fine:          15000,
			DetentionTime: 15,
			StvoPoints:    0,
		},
		{
			ID:            12,
			LawbookID:     1,
			Name:          "§21 Particularly Serious Case of Theft",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 20,
			StvoPoints:    0,
		},
		{
			ID:            15,
			LawbookID:     1,
			Name:          "§24 Robbery",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 5,
			StvoPoints:    0,
		},
		{
			ID:            16,
			LawbookID:     1,
			Name:          "§25 Aggravated Robbery",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 25,
			StvoPoints:    0,
		},
		{
			ID:            18,
			LawbookID:     1,
			Name:          "§27 Human Trafficking",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 30,
			StvoPoints:    0,
		},
		{
			ID:            20,
			LawbookID:     1,
			Name:          "§29 Hostage Taking",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 30,
			StvoPoints:    0,
		},
		{
			ID:            22,
			LawbookID:     1,
			Name:          "§30a Para. 1 Threat",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 20,
			StvoPoints:    0,
		},
		{
			ID:            27,
			LawbookID:     1,
			Name:          "§33 Endangering Road Traffic",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 20,
			StvoPoints:    0,
		},
		{
			ID:            28,
			LawbookID:     1,
			Name:          "§34 Illegal Motor Vehicle Racing",
			Description:   demoLawDescriptionPrisonSentenceRequired + ", if applicable license revocation.",
			Fine:          0,
			DetentionTime: 15,
			StvoPoints:    0,
		},
		{
			ID:            35,
			LawbookID:     1,
			Name:          "§38 Leaving the Scene of an Accident",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          5000,
			DetentionTime: 15,
			StvoPoints:    0,
		},
		{
			ID:            36,
			LawbookID:     1,
			Name:          "§39 Driving Without a License",
			Description:   demoLawDescriptionFinePossible,
			Fine:          15000,
			DetentionTime: 10,
			StvoPoints:    0,
		},
		{
			ID:            38,
			LawbookID:     1,
			Name:          "§41 Obstruction of Justice in Office",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 20,
			StvoPoints:    0,
		},
		{
			ID:            41,
			LawbookID:     1,
			Name:          "§45 Property Damage",
			Description:   demoLawDescriptionFinePossible + ", depending on property value",
			Fine:          0,
			DetentionTime: 0,
			StvoPoints:    0,
		},
		{
			ID:            44,
			LawbookID:     1,
			Name:          "§46 Terrorist Offenses",
			Description:   "ONLY BY PROSECUTION/JUDGE",
			Fine:          0,
			DetentionTime: 120,
			StvoPoints:    0,
		},
		{
			ID:            47,
			LawbookID:     1,
			Name:          "§49a Restricted Zones",
			Description:   "Administrative fine possible. Detention time required for repeat offenders; if emergency personnel are obstructed, apply a prison sentence of at least 15 and up to 30 detention units.",
			Fine:          25000,
			DetentionTime: 30,
			StvoPoints:    0,
		},
		{
			ID:            52,
			LawbookID:     1,
			Name:          "§52 Resistance Against Law Enforcement Officers",
			Description:   "Prison sentence required. If the exercise of state authority is unlawful, it is not punishable | In especially serious cases up to 30 detention units, otherwise up to 25!",
			Fine:          0,
			DetentionTime: 30,
			StvoPoints:    0,
		},
		{
			ID:            55,
			LawbookID:     1,
			Name:          "§54a Possession of Illegal Items",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 25,
			StvoPoints:    0,
		},
		{
			ID:            56,
			LawbookID:     1,
			Name:          "§55 Unauthorized Use of a Vehicle",
			Description:   "Administrative fine possible",
			Fine:          0,
			DetentionTime: 25,
			StvoPoints:    0,
		},
		{
			ID:            116,
			LawbookID:     1,
			Name:          "§18 Para. 2 Misuse of Emergency Calls",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 20,
			StvoPoints:    0,
		},
		{
			ID:            117,
			LawbookID:     1,
			Name:          "§49b Military Restricted Zones",
			Description:   "Detention time required for repeat offenders; if emergency personnel are obstructed, apply a prison sentence of at least 30 detention units.",
			Fine:          50000,
			DetentionTime: 30,
			StvoPoints:    0,
		},
		{
			ID:            121,
			LawbookID:     1,
			Name:          "§50c Ban on Face Covering in Motor Vehicles",
			Description:   "Administrative fine possible. Prison sentence required for repeat offenders or if the face covering is worn in a motor vehicle in connection with a criminal offense.",
			Fine:          0,
			DetentionTime: 0,
			StvoPoints:    0,
		},
		{
			ID:            125,
			LawbookID:     1,
			Name:          "§56 Para. 1 & 2  False Accusation ",
			Description:   "Administrative fine possible",
			Fine:          20000,
			DetentionTime: 25,
			StvoPoints:    0,
		},
		{
			ID:            130,
			LawbookID:     1,
			Name:          "§58 Disturbance of Public Peace by Threatening Criminal Offenses",
			Description:   "Administrative fine possible. Prison sentence up to 30 detention units. ",
			Fine:          20000,
			DetentionTime: 30,
			StvoPoints:    0,
		},
		{
			ID:            142,
			LawbookID:     1,
			Name:          "§42a Especially Serious Case of Disturbance of Public Peace",
			Description:   "prison sentence of 40 to 80 detention units required",
			Fine:          0,
			DetentionTime: 80,
			StvoPoints:    0,
		},

		{
			ID:            59,
			LawbookID:     2,
			Name:          "§7 Para. 2 No. 1 Unlawful Weapon Possession",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 30,
			StvoPoints:    0,
		},
		{
			ID:            60,
			LawbookID:     2,
			Name:          "§7 Para. 2 No. 2 Public Carrying of a Firearm",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 30,
			StvoPoints:    0,
		},
		{
			ID:            110,
			LawbookID:     2,
			Name:          "§7 Para. 2 No. 3 Unlawful Discharge of a Firearm",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 30,
			StvoPoints:    0,
		},

		{
			ID:            61,
			LawbookID:     3,
			Name:          "§3 Offenses",
			Description:   "Law applies to: cannabis, cocaine, methamphetamine, lysergide, opium",
			Hint:          new(demoLawDescriptionFinePossible),
			Fine:          0,
			DetentionTime: 20,
			StvoPoints:    0,
		},

		{
			ID:            62,
			LawbookID:     4,
			Name:          "§1 Registration of Aircraft",
			Description:   demoLawDescriptionFine,
			Fine:          20000,
			DetentionTime: 0,
			StvoPoints:    0,
		},
		{
			ID:            63,
			LawbookID:     4,
			Name:          "§2 No-Fly Zones",
			Description:   demoLawDescriptionFine,
			Fine:          30000,
			DetentionTime: 0,
			StvoPoints:    0,
		},

		{
			ID:            67,
			LawbookID:     5,
			Name:          "§1 Basic Rules (General Caution in Road Traffic)",
			Description:   demoLawDescriptionFine,
			Fine:          2200,
			DetentionTime: 0,
			StvoPoints:    0,
		},
		{
			ID:            70,
			LawbookID:     5,
			Name:          "§34 Up to 51-60 km/h speeding within city limits",
			Description:   "3 km/h subtract tolerance, " + demoLawDescriptionFine,
			Fine:          5400,
			DetentionTime: 0,
			StvoPoints:    3,
		},
		{
			ID:            72,
			LawbookID:     5,
			Name:          "§34 From 100 km/h speeding within city limits",
			Description:   "3 km/h subtract tolerance, Prison sentence possible",
			Fine:          12500,
			DetentionTime: 0,
			StvoPoints:    5,
		},
		{
			ID:            78,
			LawbookID:     5,
			Name:          "§12 Stopping and Parking",
			Description:   "Fine, double administrative fine and 1 road traffic point when obstructing emergency personnel/vehicles of PD/FIB/LSMD/DoJ",
			Fine:          1150,
			DetentionTime: 0,
			StvoPoints:    0,
		},
		{
			ID:            85,
			LawbookID:     5,
			Name:          "§21a Seat Belts, Protective Helmets",
			Description:   demoLawDescriptionFine,
			Fine:          2500,
			DetentionTime: 0,
			StvoPoints:    1,
		},
		{
			ID:            109,
			LawbookID:     5,
			Name:          "§33 Administrative Offenses",
			Description:   demoLawDescriptionFine,
			Fine:          2500,
			DetentionTime: 0,
			StvoPoints:    0,
		},

		{
			ID:            91,
			LawbookID:     6,
			Name:          "§7a Illegal Commercial Extraction, Transport, Trade, and Processing of Raw Materials",
			Description:   "License-free allowance: 200 kg per raw material per day. If more raw materials are carried without a license: confiscation of the excess raw materials (the allowance may be retained).",
			Fine:          30000,
			DetentionTime: 20,
			StvoPoints:    0,
		},

		{
			ID:            92,
			LawbookID:     7,
			Name:          "§12 Asset Theft",
			Description:   demoLawDescriptionFinePossible,
			Fine:          17500,
			DetentionTime: 15,
			StvoPoints:    0,
		},
		{
			ID:            95,
			LawbookID:     7,
			Name:          "§15 Tax Evasion",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 25,
			StvoPoints:    0,
		},
		{
			ID:            97,
			LawbookID:     7,
			Name:          "§17 Possession of Illicit Cash",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 20,
			StvoPoints:    0,
		},
		{
			ID:            98,
			LawbookID:     7,
			Name:          "§18 Possession of Counterfeit Money",
			Description:   "Prison sentence and Fine required",
			Fine:          25000,
			DetentionTime: 30,
			StvoPoints:    0,
		},
		{
			ID:            100,
			LawbookID:     7,
			Name:          "§20 Production of Counterfeit Money",
			Description:   "Prison sentence and fine required (2$ per counterfeit dollar)",
			Fine:          2,
			DetentionTime: 40,
			StvoPoints:    0,
		},
		{
			ID:            102,
			LawbookID:     7,
			Name:          "§22 Money Laundering",
			Description:   demoLawDescriptionPrisonSentenceRequired,
			Fine:          0,
			DetentionTime: 25,
			StvoPoints:    0,
		},
	}
)
