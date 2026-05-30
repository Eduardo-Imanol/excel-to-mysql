package routes

import "net/http"

// helloGO es el mensaje que retorna el endpoint de health check.
var helloGO string = "funcionando correctamente api rest excel a mysql"

// HomeHandler maneja la ruta raíz (GET /).
// Se usa como health check para verificar que el servidor está activo.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(helloGO))
}
