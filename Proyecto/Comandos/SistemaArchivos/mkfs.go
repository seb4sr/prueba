package sistemaarchivos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	"encoding/binary"
	//"fmt"
	"os"
	"strings"
	"time"
)

func MKfs(entrada []string) (string){
	var respuesta string
	var id string  
	fs := "2fs"    
	Valido := true
	var pathDico string

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			 
			respuesta += "ERROR MKFS, valor desconocido de parametros "+valores[1]
			return respuesta
		}

		if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])
		} else if strings.ToLower(valores[0]) == "type" {
			if strings.ToLower(valores[1]) != "full" {
				 
				respuesta += "MKFS Error. Valor de -type desconocido. "
				Valido = false
				return respuesta
			}
		} else if strings.ToLower(valores[0]) == "fs" {
			if strings.ToLower(valores[1]) == "3fs" {
				fs = "3fs"
			} else if strings.ToLower(valores[1]) != "2fs" {
				 
				Valido = false
				break
			}
	
		 
		} else {
			 
			Valido = false
			return "MKFS Error: Parametro desconocido: " + valores[0]  
		}
		
	}

	if id!=""{
		 
		for _,montado := range Structs.Montadas{
			if montado.Id == id{
				pathDico = montado.PathM
			}
		}
		if pathDico == ""{
			respuesta += "ERROR MKFS ID INCORRECTO"
			 
			Valido = false
			return respuesta
		}
	}else{
		respuesta+= "ERROR MKFS NO SE INGRESO ID"
		Valido = false
		return respuesta
	}

	if Valido{
		 
		Disco, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			respuesta += "ERROR MKFS MBR Open "+ err.Error()	
			return respuesta	
		}

		 
		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
			respuesta += "ERROR MKFS MBR Read "+ err.Error()	
			return respuesta	
		}

		 
		defer Disco.Close()

		 
		formatear := true  
		for i := 0; i < 4; i++{
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id{
				formatear = false  
				
				 
				var newSuperBloque Structs.Superblock
				Herramientas.ReadObject(Disco, &newSuperBloque, int64(mbr.Partitions[i].Start))

				 
				 
				 

				numerador := int(mbr.Partitions[i].Size) - binary.Size(Structs.Superblock{})
				denominador := 4 + binary.Size(Structs.Inode{}) + 3*binary.Size(Structs.Fileblock{})
				if fs == "3fs" {
					denominador += binary.Size(Structs.Journaling{})
				}
				n := int32(numerador / denominador)  

				 
				newSuperBloque.S_blocks_count = int32(3 * n)       
				newSuperBloque.S_free_blocks_count = int32(3 * n)  

				newSuperBloque.S_inodes_count = n       
				newSuperBloque.S_free_inodes_count = n  

				newSuperBloque.S_inode_size = int32(binary.Size(Structs.Inode{}))
				newSuperBloque.S_block_size = int32(binary.Size(Structs.Fileblock{}))

				 
				ahora := time.Now()
				copy(newSuperBloque.S_mtime[:], ahora.Format("02/01/2006 15:04"))
				 
				copy(newSuperBloque.S_umtime[:], ahora.Format("02/01/2006 15:04"))
				newSuperBloque.S_mnt_count += 1  
				newSuperBloque.S_magic = 0xEF53

				if fs == "2fs" {
					respuesta = crearEXT2(n, mbr.Partitions[i], newSuperBloque, ahora.Format("02/01/2006 15:04"), Disco)
				} else {
					respuesta = crearEXT3(n, mbr.Partitions[i], newSuperBloque, ahora.Format("02/01/2006 15:04"), Disco)
				}
				 

				 
				if Structs.UsuarioActual.Status {
					var new Structs.UserInfo
					Structs.UsuarioActual = new
				}
			}
		}

		if formatear {
			 
			 
			respuesta += "MKFS Error. No se pudo formatear la particion con id "+ id
			respuesta += "MKFS Error. No existe el id"
		}
	}

	return respuesta
}

func crearEXT2(n int32, particion Structs.Partition, newSuperBloque Structs.Superblock, date string, file *os.File) string{
	var respuesta string
	 
	 

	 
	 

	 
	newSuperBloque.S_filesystem_type = 2  
	 
	newSuperBloque.S_bm_inode_start = particion.Start + int32(binary.Size(Structs.Superblock{}))
	 
	newSuperBloque.S_bm_block_start = newSuperBloque.S_bm_inode_start + n
	 
	newSuperBloque.S_inode_start = newSuperBloque.S_bm_block_start + 3*n
	 
	newSuperBloque.S_block_start = newSuperBloque.S_inode_start + n*int32(binary.Size(Structs.Inode{}))

	 
	 
	newSuperBloque.S_free_inodes_count -= 2
	newSuperBloque.S_free_blocks_count -= 2

	 
	 
	newSuperBloque.S_first_ino = int32(2)
	 
	 
	newSuperBloque.S_first_blo = int32(2)

	 
	bmInodeData := make([]byte, n)
	bmInodeErr := Herramientas.WriteObject(file, bmInodeData, int64(newSuperBloque.S_bm_inode_start))
	if bmInodeErr != nil {
		 
		respuesta += "MKFS Error: " + bmInodeErr.Error()
		return respuesta
	}

	 
	bmBlockData := make([]byte, 3*n)
	bmBlockErr := Herramientas.WriteObject(file, bmBlockData, int64(newSuperBloque.S_bm_block_start))
	if bmBlockErr != nil {
		 
		respuesta += "MKFS Error: " + bmBlockErr.Error()
		return respuesta
	}

	 
	var newInode Structs.Inode
	for i := 0; i < 15; i++ {
		newInode.I_block[i] = -1
	}

	 
	for i := int32(0); i < n; i++ {
		err := Herramientas.WriteObject(file, newInode, int64(newSuperBloque.S_inode_start+i*int32(binary.Size(Structs.Inode{}))))
		if err != nil {
			 
			return "MKFS Error: "+ err.Error()
		}
	}

	 
	fileBlocks := make([]Structs.Fileblock, 3*n)  
	fileBlocksErr := Herramientas.WriteObject(file, fileBlocks, int64(newSuperBloque.S_bm_block_start))
	if fileBlocksErr != nil {
		 
		return "MKFS Error: "+ fileBlocksErr.Error()
	}

	 
	var Inode0 Structs.Inode
	Inode0.I_uid = 1
	Inode0.I_gid = 1
	Inode0.I_size = 0  
	 
	copy(Inode0.I_atime[:], date)
	copy(Inode0.I_ctime[:], date)
	copy(Inode0.I_mtime[:], date)
	copy(Inode0.I_type[:], "0")  
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0  

	var folderBlock0 Structs.Folderblock  
	folderBlock0.B_content[0].B_inodo = 0
	copy(folderBlock0.B_content[0].B_name[:], ".")
	folderBlock0.B_content[1].B_inodo = 0
	copy(folderBlock0.B_content[1].B_name[:], "..")
	folderBlock0.B_content[2].B_inodo = 1
	copy(folderBlock0.B_content[2].B_name[:], "users.txt")
	folderBlock0.B_content[3].B_inodo = -1

	 
	var Inode1 Structs.Inode
	Inode1.I_uid = 1
	Inode1.I_gid = 1
	Inode1.I_size = int32(binary.Size(Structs.Folderblock{}))
	copy(Inode1.I_atime[:], date)
	copy(Inode1.I_ctime[:], date)
	copy(Inode1.I_mtime[:], date)
	copy(Inode1.I_type[:], "1")  
	copy(Inode0.I_perm[:], "664")
	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}
	 
	Inode1.I_block[0] = 1
	data := "1,G,root\n1,U,root,root,123\n" 
	var fileBlock1 Structs.Fileblock  
	copy(fileBlock1.B_content[:], []byte(data))
	 

	 
	 
	
	 
	 

	 

	 
	Herramientas.WriteObject(file, newSuperBloque, int64(particion.Start))

	 
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_inode_start))
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_inode_start+1))  

	 
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_block_start))
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_block_start+1))

	 
	 
	Herramientas.WriteObject(file, Inode0, int64(newSuperBloque.S_inode_start))
	 
	Herramientas.WriteObject(file, Inode1, int64(newSuperBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

	 
	 
	Herramientas.WriteObject(file, folderBlock0, int64(newSuperBloque.S_block_start))
	 
	Herramientas.WriteObject(file, fileBlock1, int64(newSuperBloque.S_block_start+int32(binary.Size(Structs.Fileblock{}))))
	 

	 
	return "La particion con nombre "+string(particion.Name[:])+" fue formateada con exito con EXT2"
}

func crearEXT3(n int32, particion Structs.Partition, newSuperBloque Structs.Superblock, date string, file *os.File) string{
	 
	 
	 
	 
	 
	 
	 

	 
	newSuperBloque.S_filesystem_type = 3  
	 
	newSuperBloque.S_bm_inode_start = particion.Start + int32(binary.Size(Structs.Superblock{})) + int32(binary.Size(Structs.Journaling{}))
	 
	newSuperBloque.S_bm_block_start = newSuperBloque.S_bm_inode_start + n
	 
	newSuperBloque.S_inode_start = newSuperBloque.S_bm_block_start + 3*n
	 
	newSuperBloque.S_block_start = newSuperBloque.S_inode_start + n*int32(binary.Size(Structs.Inode{}))

	 
	 
	newSuperBloque.S_free_inodes_count -= 2
	newSuperBloque.S_free_blocks_count -= 2

	 
	 
	newSuperBloque.S_first_ino = int32(2)
	 
	 
	newSuperBloque.S_first_blo = int32(2)

	 
	var newJournal Structs.Journaling
	newJournal.Ultimo = 0
	newJournal.Size = int32(binary.Size(Structs.Journaling{}))

	 
	dataJ := newJournal.Contenido[0]
	copy(dataJ.Operation[:], "mkdir")
	copy(dataJ.Path[:], "/")
	copy(dataJ.Content[:], "-")
	copy(dataJ.Date[:], date)
	newJournal.Contenido[0] = dataJ

	 
	bmInodeData := make([]byte, n)
	bmInodeErr := Herramientas.WriteObject(file, bmInodeData, int64(newSuperBloque.S_bm_inode_start))
	if bmInodeErr != nil {
		 
		return "MKFS Error: "+ bmInodeErr.Error()
	}

	 
	bmBlockData := make([]byte, 3*n)
	bmBlockErr := Herramientas.WriteObject(file, bmBlockData, int64(newSuperBloque.S_bm_block_start))
	if bmBlockErr != nil {
		 
		return "MKFS Error: "+ bmBlockErr.Error()
	}

	 
	var newInode Structs.Inode
	for i := 0; i < 15; i++ {
		newInode.I_block[i] = -1
	}

	 
	for i := int32(0); i < n; i++ {
		err := Herramientas.WriteObject(file, newInode, int64(newSuperBloque.S_inode_start+i*int32(binary.Size(Structs.Inode{}))))
		if err != nil {
			 
			return "MKFS Error: " + err.Error()
		}
	}

	 
	fileBlocks := make([]Structs.Fileblock, 3*n)  
	fileBlocksErr := Herramientas.WriteObject(file, fileBlocks, int64(newSuperBloque.S_bm_block_start))
	if fileBlocksErr != nil {
		 
		return "MKFS Error: " + fileBlocksErr.Error()
	}

	 
	var Inode0 Structs.Inode
	Inode0.I_uid = 1
	Inode0.I_gid = 1
	Inode0.I_size = 0  
	 
	copy(Inode0.I_atime[:], date)
	copy(Inode0.I_ctime[:], date)
	copy(Inode0.I_mtime[:], date)
	copy(Inode0.I_type[:], "0")  
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0  

	 
	 
	 
	 
	 

	var folderBlock0 Structs.Folderblock  
	folderBlock0.B_content[0].B_inodo = 0
	copy(folderBlock0.B_content[0].B_name[:], ".")
	folderBlock0.B_content[1].B_inodo = 0
	copy(folderBlock0.B_content[1].B_name[:], "..")
	folderBlock0.B_content[2].B_inodo = 1
	copy(folderBlock0.B_content[2].B_name[:], "users.txt")
	folderBlock0.B_content[3].B_inodo = -1

	 
	var Inode1 Structs.Inode
	Inode1.I_uid = 1
	Inode1.I_gid = 1
	Inode1.I_size = int32(binary.Size(Structs.Folderblock{}))
	copy(Inode1.I_atime[:], date)
	copy(Inode1.I_ctime[:], date)
	copy(Inode1.I_mtime[:], date)
	copy(Inode1.I_type[:], "1")  
	copy(Inode0.I_perm[:], "664")
	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}
	 
	Inode1.I_block[0] = 1
	data := "1,G,root\n1,U,root,root,123\n"
	var fileBlock1 Structs.Fileblock  
	copy(fileBlock1.B_content[:], []byte(data))

	 
	 
	 
	 
	 

	 
	 

	 
	Herramientas.WriteObject(file, newSuperBloque, int64(particion.Start))

	 
	Herramientas.WriteObject(file, newJournal, int64(particion.Start+int32(binary.Size(Structs.Superblock{}))))

	 
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_inode_start))
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_inode_start+1))  

	 
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_block_start))
	Herramientas.WriteObject(file, byte(1), int64(newSuperBloque.S_bm_block_start+1))

	 
	 
	Herramientas.WriteObject(file, Inode0, int64(newSuperBloque.S_inode_start))
	 
	Herramientas.WriteObject(file, Inode1, int64(newSuperBloque.S_inode_start+int32(binary.Size(Structs.Inode{}))))

	 
	 
	Herramientas.WriteObject(file, folderBlock0, int64(newSuperBloque.S_block_start))
	 
	Herramientas.WriteObject(file, fileBlock1, int64(newSuperBloque.S_block_start+int32(binary.Size(Structs.Fileblock{}))))
	 
	 
	return "La particion con nombre "+ Herramientas.EliminartIlegibles(string(particion.Name[:])) +" fue formateada con exito con EXT3"
}