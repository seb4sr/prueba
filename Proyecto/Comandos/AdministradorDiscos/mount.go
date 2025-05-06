package administradordiscos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	//"fmt"
	"os"
	"strconv"
	"strings"
)

 
func Mount(entrada []string) (string){
	var respuesta string
	var name string	 
	var pathE string	 
	Valido := true

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			 
			respuesta += "ERROR MOUNT, valor desconocido de parametros " + valores[1]
			 
			return respuesta
		}

		 
		if strings.ToLower(valores[0]) == "path" {
			pathE = strings.ReplaceAll(valores[1],"\"","")			
			_, err := os.Stat(pathE)
			if os.IsNotExist(err) {
				 
				respuesta += "ERROR MOUNT: El disco no existe"
				Valido = false
				return respuesta  
			}
		 
		} else if strings.ToLower(valores[0]) == "name" {
			 
			name = strings.ReplaceAll(valores[1], "\"", "")
			 
			name = strings.TrimSpace(name)
		
		 
		} else {
			 
			respuesta += "ERROR MOUNT: Parametro desconocido: "+ valores[0]
			return respuesta  
		}
	}

	if Valido{
		if pathE != ""{
			if name != ""{
				 
				disco, err := Herramientas.OpenFile(pathE)
				if err != nil {
					respuesta += "ERROR NO SE PUEDE LEER EL DISCO " + err.Error()+ "\n"
					return  respuesta
				}	

				 
				var mbr Structs.MBR
				 
				if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
					respuesta += "ERROR Read " + err.Error()+ "\n"
					return  respuesta
				}
				
				 
				defer disco.Close()

				montar := true  
				reportar := false
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name{
						montar = false
						if string(mbr.Partitions[i].Type[:]) != "E" {
							if string(mbr.Partitions[i].Status[:]) != "A" {
								var id string 							
								var nuevaLetra byte = 'A' 
								contador := 1
								modificada := false															

								 
								for k:=0; k < len(Structs.Pmontaje); k++{
									if Structs.Pmontaje[k].MPath == pathE{
										 
										Structs.Pmontaje[k].Cont = Structs.Pmontaje[k].Cont + 1
										contador = int(Structs.Pmontaje[k].Cont)										
										nuevaLetra = Structs.Pmontaje[k].Letter
										modificada = true	
										break 
									}
								}

								if !modificada{
									if len(Structs.Pmontaje) > 0{
										nuevaLetra = Structs.Pmontaje[len(Structs.Pmontaje)-1].Letter +1
									}
									Structs.AddPathM(pathE, nuevaLetra, 1)
								}

								id = "56"+strconv.Itoa(contador)+string(nuevaLetra)  
								 
								 
								Structs.AddMontadas(id, pathE)

								 
								copy(mbr.Partitions[i].Status[:], "A")
								copy(mbr.Partitions[i].Id[:], id)
								mbr.Partitions[i].Correlative = int32(contador)

								 
								if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {  
									respuesta += "Error "
									return "Error "+err.Error()
								}
								reportar = true

								respuesta+="Particion con nombre "+ name+ " montada correctamente. ID: "+id
								 
							}else{
								 
								respuesta += "ERROR MOUNT. ESTA PARTICION YA FUE MONTADA PREVIAMENTE"
								return respuesta
							}
						}else{
							 
							respuesta += "ERROR MOUNT. No se puede montar una particion extendida"
							return respuesta	
						}
					}
				}

				if montar {
					 
					 
					respuesta += "ERROR MOUNT. NO SE ENCONTRO LA PARTICION " + name
					respuesta += "\nNO SE PUDO MONTAR LA PARICION \n"
					return respuesta
				}

				if reportar {
					partMontadas :="\n\nLISTA DE PARTICIONES MONTADAS EN EL DISCO\n"
					for i := 0; i < 4; i++ {
						estado := string(mbr.Partitions[i].Status[:])
						if estado == "A" {
							tmpMontadas:= "Particion: " + strconv.Itoa(i) + ", name: " +string(mbr.Partitions[i].Name[:]) + ", status: "+string(mbr.Partitions[i].Status[:])+", id: "+string(mbr.Partitions[i].Id[:])+", tipo: "+string(mbr.Partitions[i].Type[:])+", correlativo: "+ strconv.Itoa(int(mbr.Partitions[i].Correlative)) + ", fit: "+string(mbr.Partitions[i].Fit[:])+ ", start: "+strconv.Itoa(int(mbr.Partitions[i].Start))+ ", size: "+strconv.Itoa(int(mbr.Partitions[i].Size))
							partMontadas += Herramientas.EliminartIlegibles(tmpMontadas)+"\n"
						}
					}

					partMontadas +="\n\n\tPARTICIONES MONTADAS\n"
					for _,montada := range Structs.Montadas{
						partMontadas += "Id "+ string(montada.Id)+ ", Disco: "+ montada.PathM+"\n"
					}					
					respuesta += partMontadas
					 
				}
			}else{
				 
				respuesta += "ERROR: FALTA NAME  EN MOUNT"			
			}
		}else{
			 
			respuesta += "ERROR: FALTA PATH EN MOUNT"	
		}
	}

	return respuesta
	
}