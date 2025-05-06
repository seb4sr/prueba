package administradordiscos

import (
	//"fmt"
	"os"
	"strings"
)

func Rmdisk(entrada []string) string{
	 
	tmp := strings.TrimRight(entrada[1]," ")
	valores := strings.Split(tmp,"=")
	var path string

	if len(valores)!=2{
		 
		return "ERROR RMDISK, valor desconocido de parametros "+ valores[1]
	}else{		
		path = strings.ReplaceAll(valores[1],"\"","")
	}

	 
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		 
		return "RMDISK Error: El disco "+ path + " no existe"
	}

	 
	err2 := os.Remove(path)
	if err2 != nil {
		 
		return "RMDISK Error: No pudo removerse el disco "
	}
	 

	disco := strings.Split(path,"/")
	return "Disco " + disco[len(disco)-1] + " eliminado "
}