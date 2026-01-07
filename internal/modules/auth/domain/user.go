package domain

import (
	"github.com/rise-and-shine/pkg/pg"
)

type User struct {
	pg.BaseModel

	ID       int64 `json:"id"`
	PersonID int64 `json:"person_id"`
	IsActive bool  `json:"is_active"`
}

type Person struct {
	pg.BaseModel

	ID         int64  `json:"id"`
	PIN        string `json:"pin"`
	BirthDate  string `json:"birth_date"` // YYYY-MM-DD
	BirthPlace string `json:"birth_place"`
}

// {
//     "pinfl": 32505906610021,
//     "birthDate": "25.05.1990",
//     "birthCountryId": 182,
//     "citizenshipId": 182,
//     "genderCode": 1,
//     "liveStatus": 1,
//     "firstName": "SUNNATULLA",
//     "lastName": "YULDASHEV",
//     "middleName": "XUDOYAROVICH",
//     "birthPlace": "QUMQO‘RG‘ON TUMANI",
//     "birthCountry": "УЗБЕКИСТАН",
//     "citizenship": "УЗБЕКИСТАН",
//     "nationality": "УЗБЕК/УЗБЕЧКА",
//     "document": {
//         "document": "AD1157683",
//         "givePlace": "АЛМАЗАРСКИЙ РУВД ГОРОДА ТАШКЕНТА",
//         "givePlaceId": 26401,
//         "beginDate": "09.03.2022",
//         "endDate": "08.03.2032",
//         "type": "IDMS_RECV_MVD_IDCARD_CITIZEN",
//         "status": 2
//     },
//     "address": {
//         "type": "PERMANENT",
//         "cadastr": "10:04:09:01:04:5042/01",
//         "country": "O‘ZBEKISTON",
//         "region": "TOSHKENT SHAHRI",
//         "district": "YASHNOBOD TUMANI",
//         "address": "Jarboshi MFY, A.Serikbaev (BUSh) kuchasi, 8/148 B-uy",
//         "registrationDate": 1727308800000,
//         "dateFrom": null,
//         "dateTill": null,
//         "regionId": 10,
//         "districtId": 739001080
//     },
//     "photo": null,
//     "photoHashId": "EPwXm37JXD8K"
// }
