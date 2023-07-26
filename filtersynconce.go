package dapper

func renderSyncOnce(r Renderer, v Value) {
	done := v.Value.FieldByName("done")

	desc := "<unknown state>"
	if done, ok := asUint(done); ok {
		if done != 0 {
			desc = "<complete>"
		} else {
			desc = "<pending>"
		}
	}

	printWithTypeIfAmbiguous(
		r,
		v,
		"%s",
		desc,
	)
}
