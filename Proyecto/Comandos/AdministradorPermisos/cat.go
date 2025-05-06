package administradorpermisos

import (
	toolsinodos "MIA_2S_P2_201513656/ToolsInodos"
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"	
	"encoding/binary"
	//"fmt"
	"strconv"
	"strings"
)

func Cat(entrada []string) string {
	respuesta := ""
	var filen []string

	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		respuesta += "ERROR CAT: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			 
			respuesta += "ERROR CAT, valor desconocido de parametros " + valores[1] + "\n"
			 
			return respuesta
		}
		fileN := valores[0][:4]  

		 
		if strings.ToLower(fileN) == "file" {
			numero := strings.Split(strings.ToLower(valores[0]), "file")
			_, errId := strconv.Atoi(numero[1])
			if errId != nil {
				 
				return "CAT ERROR: No se pudo obtener un numero de fichero"
			}
			 
			tmp1 := strings.ReplaceAll(valores[1], "\"", "")
			filen = append(filen, tmp1)
		 
		} else {
			 
			 
			return "CAT ERROR: Parametro desconocido: " + valores[0] + "\n"
		}
	}

	 
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "CAR ERROR OPEN FILE " + err.Error() + "\n"
	}

	var mbr Structs.MBR
	 
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "CAR ERROR READ FILE " + err.Error() + "\n"
	}

	 
	defer Disco.Close()

	 
	buscar := false
	part := -1  
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			buscar = true
			break  
		}
	}

	if buscar {
		var contenido string
		var fileBlock Structs.Fileblock
		var superBloque Structs.Superblock

		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			 
			return "CAT ERROR. Particion sin formato" + "\n"
		}

		 
		for _, item := range filen {
			 
			idInodo := toolsinodos.BuscarInodo(0, item, superBloque, Disco)
			var inodo Structs.Inode
			
			 
			if idInodo > 0 {
				contenido += "\nContenido del archivo: '"+item+"':\n"
				Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
				
				 
				if inodo.I_uid == UsuarioA.IdUsr || UsuarioA.Nombre=="root"{
					
					 
					for _, idBlock := range inodo.I_block {
						if idBlock != -1 {
							Herramientas.ReadObject(Disco, &fileBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
							tmpConvertir := Herramientas.EliminartIlegibles(string(fileBlock.B_content[:]))
							contenido += tmpConvertir					
						}
					}
					contenido += "\n"
				}else{
					contenido += "ERROR CAT: No tiene permisos para visualizar el archivo " + item +"\n"
				}
				
			} else {
				contenido += "\nCAT ERROR: No se encontro el archivo " + item +"\n"
			}
		}
		respuesta += contenido
		 
	}

	return respuesta
}
