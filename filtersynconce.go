package dapper

// func renderSyncOnce(
// 	w io.Writer,
// 	v Value,
// 	p FilterPrinter,
// ) error {
// 	done := v.Value.FieldByName("done")

// 	s := "<unknown state>"
// 	if done, ok := asUint(done); ok {
// 		if done != 0 {
// 			s = "<complete>"
// 		} else {
// 			s = "<pending>"
// 		}
// 	}

// 	return formatWithTypeName(p, w, v, s)
// }
