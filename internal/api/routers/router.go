package routers

import (
	"net/http"
)

func Router() *http.ServeMux {
	sRouter := studentsRouter()
	tRouter := teachersRouter()

	tRouter.Handle("/", sRouter)
	return tRouter
}
