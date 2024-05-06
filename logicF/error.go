package logicF

import (
	"net/http"
	"text/template"
)

type Status struct {
	Code int
	Msg  string
}

// CustomResponseWriter wraps http.ResponseWriter to track if headers have been written.
type CustomResponseWriter struct {
	http.ResponseWriter
	headersWritten bool
}

// WriteHeader marks the headers as written.
func (w *CustomResponseWriter) WriteHeader(code int) {
	w.headersWritten = true
	w.ResponseWriter.WriteHeader(code)
}

// Write ensures headers are written before the body.
func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	if !w.headersWritten {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}

func Error(writer http.ResponseWriter, statusCode int) { // fonction gerant l'affichage de la page error
	var msg string
	switch statusCode {
	// if statusCode égale a http.StatusNotFound,
	case http.StatusNotFound:
		// attribution de "Not Found" a la variable msg
		msg = "Not Found"
	case http.StatusBadRequest:
		msg = "Bad request"
	default:
		msg = "Internal Server Error"
	}
	// initialise la variable a t la valeur de l'emplacement error.tmpl
	t, err := template.ParseFiles("./webpage/error.html")
	if err != nil {
		panic(err)
	}
	// ecris dans la page reponse la valeut=r du statusCode
	writer.WriteHeader(statusCode)
	// execute le writer et la struct status contenant le statusCode et le msg
	t.Execute(writer, Status{statusCode, msg})
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {
	rw := &CustomResponseWriter{ResponseWriter: w}
	if !rw.headersWritten {
		rw.WriteHeader(status)
	}
	erreur := erreur{}
	w.WriteHeader(status)
	tmpl, err := template.ParseFiles("./templates/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if status == http.StatusNotFound {
		erreur.Text = "404 NOT FOUND"
		err = tmpl.Execute(w, erreur)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if status == http.StatusBadRequest {
		erreur.Text = "BAD REQUEST"
		err = tmpl.Execute(w, erreur)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if status == http.StatusInternalServerError {
		erreur.Text = "INTERNAL SERVER ERROR"
		err = tmpl.Execute(w, erreur)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
