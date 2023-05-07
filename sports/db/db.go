package db

import (
	"time"

	"syreclabs.com/go/faker"
)

func (s *sportsRepo) seed() error {
	statement, err := s.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, name TEXT, category TEXT, mascot TEXT, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		statement, err = s.db.Prepare(`INSERT OR IGNORE INTO events(id, name, category, mascot, advertised_start_time) VALUES (?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Team().Name(),
				faker.RandomChoice(category()),
				faker.RandomChoice(mascot()),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
		}
	}

	return err
}

func category() []string {
	return []string{
		"Tennis",
		"VolleyBall",
		"Golf",
		"Table Tennis",
		"Rugby Union",
		"Lacrosse",
		"Rowing",
		"Cycling",
		"Basketball",
		"Cricket",
		"Badminton",
		"Archery",
		"Handball",
		"Surfing",
		"Field Hocky",
		"Swimming",
		"Football",
		"Boxing",
		"Baseball",
		"Ice Hockey",
		"Gymnastics",
		"Fencing",
		"Water Polo",
		"Wrestling",
		"Rugble League",
		"Netball",
	}
}

func mascot() []string {
	return []string{
		"GOBBLEDOK",
		"GARY THE SQUIRREL",
		"DOGMO",
		"WAWA GOOSE",
		"CRAZY EDDIE",
		"POPPIN' FRESH",
		"RESIDENTS EYEBALL",
		"EDDIE FROM Iron Maiden",
		"VINNIE, THE VITNER’S CHIPS MASCOT",
		"BOILER MAN",
		"SONNY FROM COCOA PUFFS",
		"COORS LIGHT CAN",
		"CRIMSON GHOST (MISFITS)",
		"MR. ZIP",
		"MORTON SALT GIRL",
		"PHILLY PHANATIC",
		"SCABBY THE RAT",
		"MR. LITTLE",
		"KOOL AID MAN",
		"COOKIE CAT",
		"FLO",
		"ROWDY RED",
		"MR. CELERY",
		"PEPSI MAN",
		"SIX FLAGS DANCING MAN",
		"TRADEMARK PIG",
		"RICH UNCLE PENNYBAGS",
		"JUDY ROSEN NOSE PICKING LOGO",
		"MYSTERIOUS FISH",
		"DOCTOR CRONUT",
		"TALKING TOM",
		"HIP HOP",
		"LITTLE PENNY",
		"PHOENIX SUNS GORILLA",
		"PINK PANTHER",
		"FIGHTING OKRA",
		"SAUSAGE RACERS",
		"MAVMAN",
		"CLIPPY",
		"BOBBY BLOTZER AS NEWBRIDGE RATMAN",
		"THE GENERAL",
		"UTZ POTATO CHIP GIRL",
		"MICHELIN MAN",
		"JENNIFER LAWRENCE COUGAR",
		"OTTO THE ORANGE",
		"TILLY",
		"ROSCOE BED BUG SNIFFING DOG",
		"YOUPPI!",
		"KING CAKE BABY",
		"NAUGGY",
	}
}
