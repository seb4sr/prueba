package administradorpermisos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	TI "MIA_2S_P2_201513656/ToolsInodos"
	"encoding/binary"
	//"fmt"
	"path/filepath"
	"strings"
)

func Rename(entrada []string) string {
	respuesta := ""
	var path string	 
	var name string	 

	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		respuesta += "ERROR RENAME: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			 
			respuesta += "ERROR RENAME, valor desconocido de parametros " + valores[1]
			 
			return respuesta
		}

		 
		if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1],"\"","")	
		 
		}else if strings.ToLower(valores[0]) == "name" {
			 
			name = strings.ReplaceAll(valores[1], "\"", "")
			 
			name = strings.TrimSpace(name)
		
		 
		} else {
			 
			respuesta += "ERROR RENAME: Parametro desconocido: "+ valores[0]
			return respuesta  
		}
	}

	if path==""{
		 
		return "ERROR RENAME FALTA PAREMETRO PATH"
	}

	if name==""{
		 
		return "ERROR RENAME FALTA PAREMETRO NAME"
	}

	 
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "CAR ERROR OPEN FILE " + err.Error() + "\n"
	}

	var mbr Structs.MBR
	 
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "CAR ERROR READ FILE " + err.Error() + "\n"
	}

	 
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

	if buscar{
		 
		var superBloque Structs.Superblock

		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			 
			return "CAT ERROR. Particion sin formato" + "\n"
		}

		var inodo Structs.Inode
		var folderBlock Structs.Folderblock
		ArchivoEncontrado := true
		 
		 
		carpeta := filepath.Dir(path)
		tmp := strings.Split(path, "/")
		nombre := tmp[len(tmp)-1]

		 
		idInodo := int32(0)
		 
		if carpeta != "/"{
			idInodo = TI.BuscarInodo(0, carpeta, superBloque, Disco)
		}

		 
		 
		Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
		 
		rename := true
		for _, idBlock := range inodo.I_block {
			if idBlock != -1{
				Herramientas.ReadObject(Disco, &folderBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Folderblock{})))))				
				for k := 2; k < 4; k++ {
					apuntador := folderBlock.B_content[k].B_inodo
					if apuntador != -1 {
						pathActual := Structs.GetB_name(string(folderBlock.B_content[k].B_name[:]))
						if pathActual == name{
							rename = false
							return "ERROR RENAME EL NOMBRE "+name+" YA EXISTE"
						}						
					}
				}
			}
		}
		
		 
		for _, idBlock := range inodo.I_block {
			if idBlock != -1{
				Herramientas.ReadObject(Disco, &folderBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Folderblock{})))))				
				 
				for k := 2; k < 4; k++ {
					apuntador := folderBlock.B_content[k].B_inodo
					if apuntador != -1 {
						pathActual := Structs.GetB_name(string(folderBlock.B_content[k].B_name[:]))
						if pathActual == nombre && rename{
							copy(folderBlock.B_content[k].B_name[:], name)
							ArchivoEncontrado = false
							 
							Herramientas.WriteObject(Disco, folderBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Folderblock{})))))
						}						
					}
				}
			}
		}
		 
		defer Disco.Close()
		if ArchivoEncontrado {
			 
			return "ERROR RENAME: EL ARCHIVO O CARPETA EN PATH NO EXISTE"
		}else{
			 
			return "El nombre del archivo fue modificado con exito "			
		}
	}

	
	return respuesta
}