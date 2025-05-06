package administradorpermisos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	ToolsInodos "MIA_2S_P2_201513656/ToolsInodos"
	"encoding/binary"
	//"fmt"
	"os"
	"strings"
)

func Edit(entrada []string) string {
	respuesta := ""
	var path string			 
	var contenido string	 

	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		respuesta += "ERROR EDIT: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _, parametro := range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			 
			respuesta += "ERROR EDIT, valor desconocido de parametros " + valores[1]
			 
			return respuesta
		}

		 
		if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1],"\"","")	
		 
		}else if strings.ToLower(valores[0]) == "contenido" {
			 
			contenido = strings.ReplaceAll(valores[1], "\"", "")
			_, err := os.Stat(contenido)
				if os.IsNotExist(err) {
					 
					respuesta +=  "MKFILE Error: El archivo cont no existe"+ "\n"
					return respuesta  
				}
		
		 
		} else {
			 
			respuesta += "ERROR EDIT: Parametro desconocido: "+ valores[0]
			return respuesta  
		}
	}

	if path==""{
		 
		return "ERROR EDIT FALTA PAREMETRO PATH"
	}

	if contenido==""{
		 
		return "ERROR EDIT FALTA PAREMETRO CONTENIDO"
	}

	 
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "EDIT ERROR OPEN FILE " + err.Error() + "\n"
	}

	var mbr Structs.MBR
	 
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "EDIT ERROR READ FILE " + err.Error() + "\n"
	}
	
	 
	editar := false
	part := -1  
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			editar = true
			break  
		}
	}

	if editar{
		var fileBlock Structs.Fileblock
		var superBloque Structs.Superblock

		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			 
			return "EDIT ERROR. Particion sin formato" + "\n"
		}

		 
		idInodo := ToolsInodos.BuscarInodo(0, path, superBloque, Disco)
		var inodo Structs.Inode
		Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

		 
		if inodo.I_uid == UsuarioA.IdUsr || UsuarioA.Nombre=="root"{
			var oldContenido string					
			 
			for _, idBlock := range inodo.I_block {
				if idBlock != -1 {
					Herramientas.ReadObject(Disco, &fileBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
					oldContenido = string(fileBlock.B_content[:])										
				}
			}
			Editar(contenido, oldContenido, len(contenido), len(oldContenido), idInodo, int64(mbr.Partitions[part].Start), Disco)
			respuesta += "\n"
		}else{
			respuesta += "ERROR EDIT: No tiene permisos para visualizar el archivo \n"
		}
		
	}
	return respuesta
}

func Editar(NewCont string, OldCont string, NewSize int, OldSize int, idInodo int32, initSuperBloque int64, disco *os.File) string{
	var respuesta string
	return respuesta
}