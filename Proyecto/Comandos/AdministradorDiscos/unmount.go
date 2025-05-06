package administradordiscos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	//"fmt"
	"strings"
)

func Unmoun(entrada []string) (string){
	var respuesta string
	var id string

	for _, parametro := range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			 
			respuesta += "ERROR UNMOUNT, valor desconocido de parametros " + valores[1]
			 
			return respuesta
		}

		if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])
		}else{
			 
			return "UNMOUNT Error: Parametro desconocido: " + valores[0]  
		}
	}

	if id!=""{		
		var pathDico string
		var registro int  

		eliminar := false
		 
		for i,montado := range Structs.Montadas{
			if montado.Id == id{
				eliminar = true
				pathDico = montado.PathM
				registro = i
			}
		}

		if eliminar{
			Disco, err := Herramientas.OpenFile(pathDico)
			if err != nil {
				return "ERROR UNMOUNT OPEN FILE "+err.Error()
			}

			var mbr Structs.MBR
			 
			if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
				return "ERROR UNMOUNT READ FILE "+err.Error()
			}

			 
			for i := 0; i < 4; i++ {		
				identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
				if identificador == id {
					//name := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					var unmount Structs.Partition

					 
					mbr.Partitions[i].Id = unmount.Id
					copy(mbr.Partitions[i].Status[:], "I")

					 
					if err := Herramientas.WriteObject(Disco, mbr, 0); err != nil {  
						return "ERROR UNMOUNT "+err.Error()
					}
				 
					break  
				}
			}			

			 
			Structs.Montadas = append(Structs.Montadas[:registro], Structs.Montadas[registro+1:]...)

			//for _,montada := range Structs.Montadas{
				 
				 
			//}	

			for i := 0; i < 4; i++ {
				estado := string(mbr.Partitions[i].Status[:])
				if estado == "A" {
					 
					 
					 
				}
			}
		}else{
			 
			return "ERROR UNMOUNT: ID NO ENCONTRADO"
		}

	}else{
		 
		return "ERROR UNMOUNT NO SE INGRESO PARAMETRO ID"
	}
	return respuesta	
}