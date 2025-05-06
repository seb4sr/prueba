package usuario

import (
	"MIA_2S_P2_201513656/Structs"
	//"fmt"
)

func Logout() (string, int){
	var respuesta string
	var res int
	if Structs.UsuarioActual.Status {
		Structs.SalirUsuario()
		 
		respuesta += "Se ha cerrado la sesion"
	}else{
		respuesta += "ERROR LOGUT: NO HAY SECION INICIADA"
		res = 1
	}

	return respuesta, res
}