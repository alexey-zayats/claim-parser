package parser

import (
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/alexey-zayats/claim-parser/internal/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var excelEpoch = time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)

// ExcelDateToDate ...
func ExcelDateToDate(excelDate string) time.Time {
	var days, _ = strconv.ParseFloat(excelDate, 64)
	return excelEpoch.Add(time.Second * time.Duration(days*86400))
}

/*
	    А    Б    В    Г    Д    Е    Ё    Ж    З    И    Й    К    Л    М    Н    О    П    Р    С    Т    У    Ф    Х    Ц    Ч    Ш    Щ    Ь    Ы    Ъ    Э    Ю    Я"
	[1040 1041 1042 1043 1044 1045 1025 1046 1047 1048 1049 1050 1051 1052 1053 1054 1055 1056 1057 1058 1059 1060 1061 1062 1063 1064 1065 1068 1067 1066 1069 1070 1071]
	  a  b  c   d   e   f   g   h   i   j   k   l   m   n   o   p   q   r   s   t   u   v   w   x   w  z
	[97 98 99 100 101 102 103 104 105 106 107 108 109 110 111 112 113 114 115 116 117 118 119 120 119 122]
	  A  B  C  D  E  F  G  H  I  J  K  L  M  N  O  P  Q  R  S  T  U  V  W  X  Y  Z
	[65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81 82 83 84 85 86 87 88 89 90]

*/
var table = map[int32]int32{
	// A -> А
	97: 1040,
	65: 1040,
	// C -> C
	99: 1057,
	67: 1057,
	// E -> Е
	101: 1045,
	69:  1045,
	// H -> Н
	72:  1053,
	104: 1053,
	// g -> Д
	103: 1044,
	// q - Д
	113: 1044,
	// K -> К
	107: 1050,
	75:  1050,
	// O -> О
	111: 1054,
	79:  1054,
	// p - Р
	80:  1056,
	112: 1056,
	// M - М
	77: 1052,
	// T - Т
	84: 1058,
	// Y - У
	89: 1059,
	// x - Х
	120: 1061,
	88:  1061,
	// n - П
	110: 1055}

// NormalizeCarNumber ...
func NormalizeCarNumber(number string) string {
	rnum := []rune(number)

	for i := range rnum {
		if i < 7 && rnum[i] > 57 && rnum[i] < 1040 {
			rnum[i] = table[rnum[i]]
		}
	}

	return util.TrimNumber(string(rnum))
}

var (
	spaceRe = regexp.MustCompile(`\s+`)
	nanRe   = regexp.MustCompile(`\D`)
)

// ParseFIO ...s
func ParseFIO(s string) (model.FIO, bool) {

	data := spaceRe.Split(strings.TrimSpace(s), -1)

	success := false
	fio := model.FIO{}

	if len(data) >= 3 {
		success = true
		fio.Lastname = data[0]
		fio.Firstname = data[1]
		fio.Patronymic = data[2]
	} else if len(data) == 2 {
		fio.Lastname = data[0]
		fio.Firstname = data[1]
	} else if len(data) == 1 {
		fio.Lastname = data[0]
	}

	return fio, success
}

// ParseInt64 ...
func ParseInt64(s string) (int64, bool) {
	s = nanRe.ReplaceAllString(s, "")
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}
