package administradorpermisos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	ToolsInodos "MIA_2S_P2_201513656/ToolsInodos"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func Copy(entrada []string) string {
	respuesta := ""
	var path string			 
	var destino string	 

	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		respuesta += "ERROR COPY: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _, parametro := range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			 
			respuesta += "ERROR COPY, valor desconocido de parametros " + valores[1]
			 
			return respuesta
		}

		 
		if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1],"\"","")	
		 
		}else if strings.ToLower(valores[0]) == "destino" {
			 
			destino = strings.ReplaceAll(valores[1], "\"", "")
		
		 
		} else {
			 
			respuesta += "ERROR COPY: Parametro desconocido: "+ valores[0]
			return respuesta  
		}
	}

	if path==""{
		 
		return "ERROR COPY FALTA PAREMETRO PATH"
	}

	if destino==""{
		 
		return "ERROR COPY FALTA PAREMETRO NAME"
	}

	 
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "COPY ERROR OPEN FILE " + err.Error() + "\n"
	}

	var mbr Structs.MBR
	 
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "COPY ERROR READ FILE " + err.Error() + "\n"
	}
	
	 
	copy := false
	part := -1  
	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			copy = true
			break  
		}
	}

	if copy{
		var superBloque Structs.Superblock

		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			 
			return "COPY ERROR. Particion sin formato" + "\n"
		}

		 
		idInodoDestino := ToolsInodos.BuscarInodo(0, destino, superBloque, Disco)
		var inodoDestino Structs.Inode
		Herramientas.ReadObject(Disco, &inodoDestino, int64(superBloque.S_inode_start+(idInodoDestino*int32(binary.Size(Structs.Inode{})))))

		 
		if inodoDestino.I_uid == UsuarioA.IdUsr || UsuarioA.Nombre=="root" || inodoDestino.I_gid == UsuarioA.IdGrp{					
			 
			idNewInodo := ToolsInodos.BuscarInodo(0, path, superBloque, Disco)
			var NewInodo Structs.Inode
			Herramientas.ReadObject(Disco, &NewInodo, int64(superBloque.S_inode_start+(idNewInodo*int32(binary.Size(Structs.Inode{})))))

			var fileBlock Structs.Fileblock
			for _, idBlock := range NewInodo.I_block {
				if idBlock != -1 {
					Herramientas.ReadObject(Disco, &fileBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
					tmpConvertir := Herramientas.EliminartIlegibles(string(fileBlock.B_content[:]))
					fmt.Println(idBlock," ",tmpConvertir)	
				}
			}

		}else{
			respuesta += "ERROR COPY: No tiene permisos para copiar archivos a esta carpeta \n"
		}
		
	}
	return respuesta
}

func Copiar(idNodoDestino int32, idNodoCopiar int32, initSuperBloque int64, disco *os.File){}