package usuario

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func Mkusr(entrada []string) string {
	var respuesta string
	var user string  
	var pass string  
	var grp string   
	Valido := true
	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		Valido = false
		respuesta += "ERROR MKUSR: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	if UsuarioA.Nombre != "root" {
		Valido = false
		 
		respuesta += "ERROR MKGRO: ESTE USUARIO NO CUENTA CON LOS PERMISOS PARA REALIZAR ESTA ACCION"
		return respuesta
	}

	for _, parametro := range entrada[1:] {
		tmp := strings.TrimRight(parametro, " ")
		valores := strings.Split(tmp, "=")

		if len(valores) != 2 {
			 
			respuesta += "ERROR MKGRP, valor desconocido de parametros " + valores[1] + "\n"
			Valido = false
			 
			return respuesta
		}		

		 
		if strings.ToLower(valores[0]) == "grp" {
			grp = strings.ReplaceAll(valores[1],"\"","")
			 
			if len(grp) > 10 {
				Valido = false
				 
				return "ERROR MKGRP: grp debe tener maximo 10 caracteres"
			}
		 
		} else if strings.ToLower(valores[0]) == "user" {
			user = strings.ReplaceAll(valores[1],"\"","")
			tmp1 := strings.TrimRight(user, " ")
			user = tmp1
			 
			if len(user) > 10 {
				Valido = false
				 
				return "ERROR MKGRP: user debe tener maximo 10 caracteres"
			}
		 
		} else if strings.ToLower(valores[0]) == "pass" {
			pass = valores[1]
			 
			if len(pass) > 10 {
				Valido = false
				 
				return "ERROR MKGRP: pass debe tener maximo 10 caracteres"
			}
		 
		} else {
			Valido = false
			 
			 
			return "MKUSR ERROR: Parametro desconocido: " + valores[0] + "\n"
		}
	}

	 
	if pass == "" {
		Valido = false
		 
		return "MKUSR ERROR: FALTO EL PARAMETRO PASS " + "\n"
	}

	if user == "" {
		Valido = false
		 
		return "MKUSR ERROR: FALTO EL PARAMETRO USER " + "\n"
	}

	if grp == "" {
		Valido = false
		 
		return "MKUSR ERROR: FALTO EL PARAMETRO GRP " + "\n"
	}

	if Valido {
		file, err := Herramientas.OpenFile(UsuarioA.PathD)
		if err != nil {
			return "ERROR MKUSR ERROR SB OPEN FILE " + err.Error() + "\n"
		}

		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR MKUSR ERROR SB READ FILE " + err.Error() + "\n"
		}

		 
		defer file.Close()

		 
		AddNewUser := false
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == UsuarioA.IdPart {
				part = i
				AddNewUser = true
				break  
			}
		}

		if AddNewUser {
			var superBloque Structs.Superblock
			errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
			if errREAD != nil {
				 
				return "MKUSR Error. Particion sin formato" + "\n"
			}

			var inodo Structs.Inode
			 
			Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

			 
			var contenido string
			var fileBlock Structs.Fileblock
			var idFb int32  
			for _, item := range inodo.I_block {
				if item != -1 {
					Herramientas.ReadObject(file, &fileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
					contenido += string(fileBlock.B_content[:])
					idFb = item
				}
			}

			lineaID := strings.Split(contenido, "\n")

			ExGrupo := false
			for _, registro := range lineaID[:len(lineaID)-1] {
				datos := strings.Split(registro, ",")
				 
				if len(datos) == 3 {
					if datos[2] == grp {
						ExGrupo = true						
					}
					
				} 
				 
				if len(datos) == 5 {
					if datos[3] == user {
						 
						return "MKUSR ERROR: ESTE USUARIO YA EXISTE"
					}
				}
			}

			if !ExGrupo {
				 
				return "MKURS ERROR, NO EXISTE EL GRUPO, POR FAVOR INGRESE UN GRUPO QUE SI EXISTA"
			}

			 
			 
			id := -1         
			var errId error  
			for i := len(lineaID) - 2; i >= 0; i-- {
				registro := strings.Split(lineaID[i], ",")
				 
				if registro[1] == "U" {
					 
					if registro[0] != "0" {
						 
						id, errId = strconv.Atoi(registro[0])
						if errId != nil {
							 
							return "MKUSR ERROR: NO SE PUDO OBTENER UN NUEVO ID PARA EL NUEVO GRUPO"
						}
						id++
						break
					}
				}
			}

			 
			if id != -1 {
				contenidoActual := string(fileBlock.B_content[:])
				posicionNulo := strings.IndexByte(contenidoActual, 0)
				data := fmt.Sprintf("%d,U,%s,%s,%s\n", id, grp, user, pass)

				 
				if posicionNulo != -1 {
					libre := 64 - (posicionNulo + len(data))
					if libre > 0 {
						copy(fileBlock.B_content[posicionNulo:], []byte(data))
						 
						Herramientas.WriteObject(file, fileBlock, int64(superBloque.S_block_start+(idFb*int32(binary.Size(Structs.Fileblock{})))))
					} else {
						 
						data1 := data[:len(data)+libre]
						 
						copy(fileBlock.B_content[posicionNulo:], []byte(data1))
						Herramientas.WriteObject(file, fileBlock, int64(superBloque.S_block_start+(idFb*int32(binary.Size(Structs.Fileblock{})))))

						 
						guardoInfo := true
						for i, item := range inodo.I_block {
							 
							if item == -1 {
								guardoInfo = false
								 
								inodo.I_block[i] = superBloque.S_first_blo
								 
								superBloque.S_free_blocks_count -= 1
								superBloque.S_first_blo += 1
								data2 := data[len(data)+libre:]
								 
								var newFileBlock Structs.Fileblock
								copy(newFileBlock.B_content[:], []byte(data2))

								 
								 
								Herramientas.WriteObject(file, superBloque, int64(mbr.Partitions[part].Start))

								 
								Herramientas.WriteObject(file, byte(1), int64(superBloque.S_bm_block_start+inodo.I_block[i]))

								 
								Herramientas.WriteObject(file, inodo, int64(superBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

								 
								Herramientas.WriteObject(file, newFileBlock, int64(superBloque.S_block_start+(inodo.I_block[i]*int32(binary.Size(Structs.Fileblock{})))))
								break
							}
						}

						if guardoInfo {
							 
							return "MKUSR ERROR: ESPACIO INSUFICIENTE PARA EL NUEVO USUARIO. "
						}
					}
					 
					respuesta += "Se ha agregado el usuario '"+user+"' al grupo '"+grp+"' exitosamente. "
					for k:=0; k<len(lineaID)-1; k++{
						 
					}
					return respuesta
				}else{
					respuesta += "ERROR MKUSR NO HAY ESPACIO SUFICIENTE"
					 
				}
			}
			
			
			return respuesta
			
		}  
	}

	return respuesta
}
