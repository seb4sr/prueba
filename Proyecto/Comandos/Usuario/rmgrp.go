package usuario

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	"encoding/binary"
	//"fmt"
	"strings"
)

func Rmgrp(entrada []string) string {
	var respuesta string
	var name string
	UsuarioA := Structs.UsuarioActual
	
	if !UsuarioA.Status {
		respuesta += "ERROR RMGRP: NO HAY SECION INICIADA"+ "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR"+ "\n"
		return respuesta
	}

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			 
			respuesta += "ERROR RMGRP, valor desconocido de parametros " + valores[1]+ "\n"
			 
			return respuesta
		}

		 
		if strings.ToLower(valores[0]) == "name" {
			name = (valores[1])
			 
			if len(name) > 10 {
				 
				return "ERROR RMGRP: name debe tener maximo 10 caracteres"
			}
		 
		} else {
			 
			 
			return "RMGRP ERROR: Parametro desconocido: "+valores[0] + "\n"
		}
	}

	if UsuarioA.Nombre == "root"{
		file, err := Herramientas.OpenFile(UsuarioA.PathD)
		if err != nil {
			return "RMGRP ERRORSB OPEN FILE "+err.Error()+ "\n"
		}

		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "RMGRP ERRORSB READ FILE "+err.Error()+ "\n"
		}

		 
		defer file.Close()

		 
		delete := false
		part := -1  
		for i := 0; i < 4; i++ {		
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == UsuarioA.IdPart {
				part = i
				delete = true
				break  
			}
		}

		if delete{
			var superBloque Structs.Superblock
			errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
			if errREAD != nil {
				 
				return "RMGRP ERROR. Particion sin formato"+ "\n"
			}

			var inodo Structs.Inode		
			 
			Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start + int32(binary.Size(Structs.Inode{}))))
			
			 
			var contenido string
			var fileBlock Structs.Fileblock
			for _, item := range inodo.I_block {
				if item != -1 {
					Herramientas.ReadObject(file, &fileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
					contenido += string(fileBlock.B_content[:])
				}
			}
			
			lineaID := strings.Split(contenido, "\n")
			modificarUs := false
			for k:=0; k<len(lineaID); k++{
				datos := strings.Split(lineaID[k], ",")
				if len(datos) == 3 {
					if datos[2] == name {
						 
						if datos[0] != "0"{
							modificarUs = true
							datos[0]="0"
							lineaID[k] = datos[0] + "," + datos[1] + "," + datos[2]
						}else{
							 
							return "ERROR RMGRP ESTE GRUPO YA FUE ELIMINADO PREVIAMENTE"
						}						
					}
				}
			}

			if modificarUs{
				 
				for k:=0; k<len(lineaID); k++{
					datos := strings.Split(lineaID[k], ",")
					if len(datos) ==5{
						if datos[2] == name{
							if datos[0] != "0"{
								datos[0]="0"
								lineaID[k] = datos[0] + "," + datos[1] + "," + datos[2]+ "," + datos[3]+ "," + datos[4]
							}
						}
					}
				}
				

				mod := ""
					for _, reg := range lineaID {
						mod += reg + "\n"
					}

				inicio := 0
				var fin int
				if len(mod) > 64 {
					 
					fin = 64
				} else {
					 
					fin = len(mod)
				}

				for _, newItem := range inodo.I_block{
					if newItem != -1 {
						 
						data := mod[inicio:fin]
						 
						var newFileBlock Structs.Fileblock
						copy(newFileBlock.B_content[:], []byte(data))
						Herramientas.WriteObject(file, newFileBlock, int64(superBloque.S_block_start+(newItem*int32(binary.Size(Structs.Fileblock{})))))
						 
						inicio = fin
						calculo := len(mod[fin:])  
						 
						if calculo > 64 {
							fin += 64
						} else {
							fin += calculo
						}
					}
				}

				 
				respuesta += "El grupo '"+name+"' fue eliminado con extiso"
				for k:=0; k<len(lineaID)-1; k++{
					 
				}
				return respuesta
			}			
		}

	}else{
		 
		respuesta += "RMGRP ERROR: ESTE USUARIO NO CUENTA CON LOS PERMISOS PARA REALIZAR ESTA ACCION"
	}
	
	return respuesta
}