package dapper

func renderSyncOnce(r Renderer, v Value) {
	done := v.Value.FieldByName("done")

	desc := "<pending>"
	if done.IsValid() && !done.IsZero() {
		desc = "<complete>"
	}

	printWithTypeIfAmbiguous(
		r,
		v,
		"%s",
		desc,
	)
}
