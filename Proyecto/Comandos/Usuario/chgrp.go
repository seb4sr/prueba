package usuario

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	"encoding/binary"
	//"fmt"
	"strings"
)

func Chgrp(entrada []string) string{
	var respuesta string
	var user string
	var grp string
	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		respuesta += "ERROR MKUSR: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	if UsuarioA.Nombre != "root" {
		 
		respuesta += "ERROR MKGRO: ESTE USUARIO NO CUENTA CON LOS PERMISOS PARA REALIZAR ESTA ACCION"
		return respuesta
	}

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			 
			respuesta += "ERROR MKGRP, valor desconocido de parametros " + valores[1] + "\n"
			 
			return respuesta
		}

		 
		if strings.ToLower(valores[0]) == "grp" {
			grp = (valores[1])
			 
			if len(grp) > 10 {
				 
				return "ERROR MKGRP: grp debe tener maximo 10 caracteres"
			}
			 
		} else if strings.ToLower(valores[0]) == "user" {
			user = valores[1]
			 
			if len(user) > 10 {
				 
				return "ERROR MKGRP: user debe tener maximo 10 caracteres"
			}
		 
		} else {
			 
			 
			return "CHGRP ERROR: Parametro desconocido: " + valores[0] + "\n"
		}
	}

	 
	if user == "" {
		 
		return "MKUSR ERROR: FALTO EL PARAMETRO USER " + "\n"
	}

	if grp == "" {
		 
		return "MKUSR ERROR: FALTO EL PARAMETRO GRP " + "\n"
	}	

	file, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "RMUSR ERRORSB OPEN FILE "+err.Error()+ "\n"
	}

	var mbr Structs.MBR
	 
	if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
		return "RMUSR ERRORSB READ FILE "+err.Error()+ "\n"
	}

	 
	defer file.Close()

	 
	continuar := false
	part := -1  
	for i := 0; i < 4; i++ {		
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			continuar = true
			break  
		}
	}

	if continuar{
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
		 
		NOExGrupo := true
		for _, registro := range lineaID[:len(lineaID)-1] {
			datos := strings.Split(registro, ",")
			 
			if len(datos) == 3 {
				if datos[2] == grp {
					NOExGrupo = false
					break
				}
			}
		}

		if NOExGrupo {
			 
			return "CHGRP ERROR, NO EXISTE EL GRUPO, POR FAVOR INGRESE UN GRUPO QUE SI EXISTA"
		}

		chGro := false
		for k:=0; k<len(lineaID); k++{
			datos := strings.Split(lineaID[k], ",")
			if len(datos) ==5{
				if datos[3] == user{
					if datos[0] != "0"{
						chGro = true
						datos[2] = grp
						lineaID[k] = datos[0] + "," + datos[1] + "," + datos[2]+ "," + datos[3]+ "," + datos[4]
					}else{
						 
						return "ERROR RMUSR, EL USUARIO '"+user+"' fue eliminado "
					}
				}
			}
		}

		if chGro{
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

			 
			respuesta += "El usuario '"+user+"' fue cambiado al grupo '"+grp+"' exitosamente"
			for k:=0; k<len(lineaID)-1; k++{
				 
			}
			return respuesta
		}
	}else{
		 
		return "ERROR CHGRP: OCURRIO UN ERROR INESPERADO"
	}
	
	return respuesta
}