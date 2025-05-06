package administradorpermisos

import (
	TI "MIA_2S_P2_201513656/ToolsInodos"
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"	
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func MKfile(entrada []string) string{
	respuesta := "Comando mkfile"
	parametrosDesconocidos := false
	var path string
	var cont string	 
	size := 0 		 
	r := false
	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		fmt.Println("ERROR MKFILE: SESION NO INICIADA")
		respuesta += "ERROR MKFILE: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores) == 2{
			 
			if strings.ToLower(valores[0]) == "path"  {
				path = strings.ReplaceAll(valores[1],"\"","")
			 
			}else if strings.ToLower(valores[0]) == "size"{
				 
				var err error
				size, err = strconv.Atoi(valores[1])  
				if err != nil {
					fmt.Println("MKFILE Error: Size solo acepta valores enteros. Ingreso: ", valores[1])
					return "MKFILE Error: Size solo acepta valores enteros. Ingreso: " + valores[1]
				}

				 
				if size < 0 {
					fmt.Println("MKFILE Error: Size solo acepta valores positivos. Ingreso: ", valores[1])
					return "MKFILE Error: Size solo acepta valores positivos. Ingreso: "+ valores[1]
				}
			 
			}else if strings.ToLower(valores[0]) == "cont"{
				cont = strings.ReplaceAll(valores[1], "\"", "")
				_, err := os.Stat(cont)
				if os.IsNotExist(err) {
					fmt.Println("MKFILE Error: El archivo cont no existe")
					respuesta +=  "MKFILE Error: El archivo cont no existe"+ "\n"
					return respuesta  
				}
			}else{
				parametrosDesconocidos = true
			}
		}else if len(valores) == 1{
			if strings.ToLower(valores[0]) == "r"{
				r = true
			}else{
				parametrosDesconocidos = true
			}
		}else{
			parametrosDesconocidos = true
		}

		if parametrosDesconocidos{
			fmt.Println("MKFILE Error: Parametro desconocido: ", valores[0])
			respuesta = "MKFILE Error: Parametro desconocido: "+ valores[0]
			return respuesta  
		}
	}
	

	if path ==""{
		fmt.Println("MKFIEL ERROR NO SE INGRESO PARAMETRO PATH")
		return "MKFIEL ERROR NO SE INGRESO PARAMETRO PATH"
	}

	 
	Disco, err := Herramientas.OpenFile(UsuarioA.PathD)
	if err != nil {
		return "MKFILE ERROR OPEN FILE "+err.Error()+ "\n"
	}

	var mbr Structs.MBR
	 
	if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
		return "MKFILE ERROR READ FILE "+err.Error()+ "\n"
	}

	 
	defer Disco.Close()

	 
	agregar := false
	part := -1  
	for i := 0; i < 4; i++ {		
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == UsuarioA.IdPart {
			part = i
			agregar = true
			break  
		}
	}

	if agregar{
		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			fmt.Println("MKFILE ERROR. Particion sin formato")
			return "MKFILE ERROR. Particion sin formato"+ "\n"
		}

		 
		stepPath := strings.Split(path, "/")
		finRuta := len(stepPath) - 1  
		idInicial := int32(0)
		idActual := int32(0)
		crear := -1
		 
		for i, itemPath := range stepPath[1:finRuta] {
			idActual = TI.BuscarInodo(idInicial, "/"+itemPath, superBloque, Disco)
			 
			if idInicial != idActual {
				idInicial = idActual
			} else {
				crear = i + 1  
				break
			}
		}

		 
		if crear != -1 {
			if r {
				for _, item := range stepPath[crear:finRuta] {
					idInicial = TI.CreaCarpeta(idInicial, item, int64(mbr.Partitions[part].Start), Disco)
					if idInicial == 0 {
						fmt.Println("MKDIR ERROR: No se pudo crear carpeta")
						return "MKFILE ERROR: No se pudo crear carpeta"
					}
				}
			} else {
				fmt.Println("MKDIR ERROR: Carpeta ", stepPath[crear], " no existe. Sin permiso de crear carpetas padre")
				return "MKFILE ERROR: Carpeta "+ stepPath[crear]+ " no existe. Sin permiso de crear carpetas padre"
			}
		}

		 
		idNuevo := TI.BuscarInodo(idInicial, "/"+stepPath[finRuta], superBloque, Disco)
		if idNuevo == idInicial {
			if cont == "" {
				digito := 0
				var content string

				 
				for i := 0; i < size; i++ {
					if digito == 10 {
						digito = 0
					}
					content += strconv.Itoa(digito)
					digito++
				}
				respuesta = crearArchivo(idInicial, stepPath[finRuta], size, content, int64(mbr.Partitions[part].Start), Disco)				
			}else{
				archivoC, err := Herramientas.OpenFile(cont)
				if err != nil {
					return "MKFILE ERROR OPEN FILE "+err.Error()+ "\n"
				}

				 
				content, err := ioutil.ReadFile(cont)
				if err != nil {
					fmt.Println(err)
					return "ERROR MKFILE "+err.Error()
				}
				 
				defer archivoC.Close()
				respuesta = crearArchivo(idInicial, stepPath[finRuta], size, string(content), int64(mbr.Partitions[part].Start), Disco)
			}
		}else{
			fmt.Println("El archivo ya existe")
			return "ERROR: El archivo ya existe"
		}
	}
	return respuesta
}

func crearArchivo(idInodo int32, file string, size int, contenido string, initSuperBloque int64, disco *os.File) string{
	 
	var superB Structs.Superblock
	Herramientas.ReadObject(disco, &superB, initSuperBloque)
	 
	var inodoFile Structs.Inode
	Herramientas.ReadObject(disco, &inodoFile, int64(superB.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

	 
	for i := 0; i < 12; i++{
		idBloque := inodoFile.I_block[i]
		if idBloque != -1 {
			 
			var folderBlock Structs.Folderblock
			Herramientas.ReadObject(disco, &folderBlock, int64(superB.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

			 
			for j := 2; j < 4; j++ {
				apuntador := folderBlock.B_content[j].B_inodo
				 
				if apuntador == -1 {
					 
					copy(folderBlock.B_content[j].B_name[:], file)
					ino := superB.S_first_ino  
					folderBlock.B_content[j].B_inodo = ino
					 
					Herramientas.WriteObject(disco, folderBlock, int64(superB.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

					 
					var newInodo Structs.Inode
					newInodo.I_uid = Structs.UsuarioActual.IdUsr
					newInodo.I_gid = Structs.UsuarioActual.IdGrp
					newInodo.I_size = int32(size)  
					 
					ahora := time.Now()
					date := ahora.Format("02/01/2006 15:04")
					copy(newInodo.I_atime[:], date)
					copy(newInodo.I_ctime[:], date)
					copy(newInodo.I_mtime[:], date)
					copy(newInodo.I_type[:], "1")  
					copy(newInodo.I_perm[:], "664")

					 
					for i := int32(0); i < 15; i++ {
						newInodo.I_block[i] = -1
					}

					 
					fileblock := superB.S_first_blo

					 
					inicio := 0
					fin := 0
					sizeContenido := len(contenido)
					if sizeContenido < 64 {
						fin = len(contenido)
					} else {
						fin = 64
					}

					 
					for i := int32(0); i < 12; i++ {
						newInodo.I_block[i] = fileblock
						 
						data := contenido[inicio:fin]
						var newFileBlock Structs.Fileblock
						copy(newFileBlock.B_content[:], []byte(data))
						 
						Herramientas.WriteObject(disco, newFileBlock, int64(superB.S_block_start+(fileblock*int32(binary.Size(Structs.Fileblock{})))))

						 
						superB.S_free_blocks_count -= 1
						superB.S_first_blo += 1

						 
						Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_block_start+fileblock))

						 
						calculo := len(contenido[fin:])
						if calculo > 64 {
							inicio = fin
							fin += 64
						} else if calculo > 0 {
							inicio = fin
							fin += calculo
						} else {
							 
							break
						}
						 
						fileblock++
					}

					 
					Herramientas.WriteObject(disco, newInodo, int64(superB.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

					 
					superB.S_free_inodes_count -= 1
					superB.S_first_ino += 1
					 
					Herramientas.WriteObject(disco, superB, initSuperBloque)

					 
					Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_inode_start+ino))

					return "Archivo creado exitosamente"
				} 
			} 
		}else{
			 
			 
			block := superB.S_first_blo  
			inodoFile.I_block[i] = block
			 
			Herramientas.WriteObject(disco, &inodoFile, int64(superB.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

			 
			var folderBlock Structs.Folderblock
			bloque := inodoFile.I_block[0]  
			Herramientas.ReadObject(disco, &folderBlock, int64(superB.S_block_start+(bloque*int32(binary.Size(Structs.Folderblock{})))))

			 
			var newFolderBlock1 Structs.Folderblock
			newFolderBlock1.B_content[0].B_inodo = folderBlock.B_content[0].B_inodo  
			copy(newFolderBlock1.B_content[0].B_name[:], ".")
			newFolderBlock1.B_content[1].B_inodo = folderBlock.B_content[1].B_inodo  
			copy(newFolderBlock1.B_content[1].B_name[:], "..")
			ino := superB.S_first_ino                           
			newFolderBlock1.B_content[2].B_inodo = ino          
			copy(newFolderBlock1.B_content[2].B_name[:], file)  
			newFolderBlock1.B_content[3].B_inodo = -1
			 
			Herramientas.WriteObject(disco, newFolderBlock1, int64(superB.S_block_start+(block*int32(binary.Size(Structs.Folderblock{})))))

			 
			Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_block_start+block))

			 
			superB.S_first_blo += 1
			superB.S_free_blocks_count -= 1

			 
			var newInodo Structs.Inode
			newInodo.I_uid = Structs.UsuarioActual.IdUsr
			newInodo.I_gid = Structs.UsuarioActual.IdGrp
			newInodo.I_size = int32(size)  
			 
			ahora := time.Now()
			date := ahora.Format("02/01/2006 15:04")
			copy(newInodo.I_atime[:], date)
			copy(newInodo.I_ctime[:], date)
			copy(newInodo.I_mtime[:], date)
			copy(newInodo.I_type[:], "1")  
			copy(newInodo.I_mtime[:], "664")

			 
			for i := int32(0); i < 15; i++ {
				newInodo.I_block[i] = -1
			}

			 
			fileblock := superB.S_first_blo

			 
			inicio := 0
			fin := 0
			sizeContenido := len(contenido)
			if sizeContenido < 64 {
				fin = len(contenido)
			} else {
				fin = 64
			}

			 
			for i := int32(0); i < 12; i++{
				newInodo.I_block[i] = fileblock
				 
				data := contenido[inicio:fin]
				var newFileBlock Structs.Fileblock
				copy(newFileBlock.B_content[:], []byte(data))
				 
				Herramientas.WriteObject(disco, newFileBlock, int64(superB.S_block_start+(fileblock*int32(binary.Size(Structs.Fileblock{})))))

				 
				superB.S_free_blocks_count -= 1
				superB.S_first_blo += 1

				 
				Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_block_start+fileblock))

				 
				calculo := len(contenido[fin:])
				if calculo > 64 {
					inicio = fin
					fin += 64
				} else if calculo > 0 {
					inicio = fin
					fin += calculo
				} else {
					 
					break
				}
				 
				fileblock++
			}

			 
			Herramientas.WriteObject(disco, newInodo, int64(superB.S_inode_start+(ino*int32(binary.Size(Structs.Inode{})))))

			 
			superB.S_free_inodes_count -= 1
			superB.S_first_ino += 1
			 
			Herramientas.WriteObject(disco, superB, initSuperBloque)

			 
			Herramientas.WriteObject(disco, byte(1), int64(superB.S_bm_inode_start+ino))

			return "Archivo creado exitosamente"
		}
	}

	return "ERROR MKFILE: OCURRIO UN ERROR INESPERADO AL CREAR EL ARCHIVO"
}