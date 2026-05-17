package dapper

func renderSyncOnce(r Renderer, v Value) {
	done := v.Value.FieldByName("done")

	desc := "<unknown state>"
	if done.IsValid() {
		if done.IsZero() {
			desc = "<pending>"
		} else {
			desc = "<complete>"
		}
	}

	printWithTypeIfAmbiguous(
		r,
		v,
		"%s",
		desc,
	)
}
