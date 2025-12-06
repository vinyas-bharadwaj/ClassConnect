package routers

import (
	"net/http"
)

func Router() *http.ServeMux {
	eRouter := execsRouter()
	sRouter := studentsRouter()
	tRouter := teachersRouter()

	sRouter.Handle("/", eRouter)
	tRouter.Handle("/", sRouter)
	return tRouter
}
