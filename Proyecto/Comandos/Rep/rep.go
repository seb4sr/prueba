package rep

import (
	toolsinodos "MIA_2S_P2_201513656/ToolsInodos"
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"	
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Rep(entrada []string) string{
	var respuesta string
	var name string  
	var path string  
	var id string    
	var rutaFile string	 
	Valido := true 

	for _, parametro := range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			 
			 
			return "ERROR REP, valor desconocido de parametros "+valores[1]
		}

		if strings.ToLower(valores[0]) == "name" {
			name = strings.ToLower(valores[1])
		} else if strings.ToLower(valores[0]) == "path" {
			path = strings.ReplaceAll(valores[1], "\"", "")
		} else if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])
		} else if strings.ToLower(valores[0]) == "path_file_ls" {
			rutaFile = strings.ReplaceAll(valores[1], "\"", "")
		} else {
			 
			respuesta+="REP Error: Parametro desconocido: " + valores[0]
			Valido = false
			break  
		}
	}

	if Valido{
		if name != "" && id != "" && path != "" {			
			switch name{
			case "mbr": 
				 
				respuesta+= Rmbr(path, id)
			case "disk": 
				 
				respuesta+= disk(path, id)
			case "bm_inode": 
				 
				respuesta += BM_inode(path, id)
			case "bm_block": 
				 
				respuesta += BM_Bloque(path, id)
			case "sb": 
				 
				respuesta += superBloque(path, id)
			case "file": 
				 
				respuesta += FILE(path, id, rutaFile)
			case "ls": 
				respuesta += LS(path, id, rutaFile)
				 
			case "journal": 
				journal(path, id)
				 
			case "tree": 
				tree(path,id)
				 
			default:
				 
				respuesta+="REP Error: Reporte "+ name+" desconocido"
			}
		}else{
			 
			respuesta+= "REP Error: Faltan parametros"
		}
	}
	return respuesta
}

 
func Rmbr (path string, id string) string{
	var Respuesta string
	var pathDico string
	Valido := false

	 
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			Valido = true
		}
	}

	if Valido{
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]
		
		tmp = strings.Split(pathDico, "/")
		NOmbreDis := strings.Split(tmp[len(tmp)-1], ".")[0]
		
		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			Respuesta += "ERROR REP MBR Open "+ err.Error()		
		}

		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			Respuesta += "ERROR REP MBR Read "+ err.Error()		
		}

		 
		defer file.Close()

		 
		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += " <tr>\n  <td bgcolor='SlateBlue' COLSPAN=\"2\"> Reporte MBR </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_tamano </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", mbr.MbrSize)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#AFA1D1'> mbr_fecha_creacion </td> \n  <td bgcolor='#AFA1D1'> %s </td> \n </tr> \n", string(mbr.FechaC[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> mbr_disk_signature </td> \n  <td bgcolor='Azure'> %d </td> \n </tr>  \n", mbr.Id)
		cad += Structs.RepGraphviz(mbr, file)
		cad += "</table> > ]\n}"

		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
		Respuesta += "Reporte de MBR del disco "+NOmbreDis+" creado con el nombre "+nombre+".png"
	}else{
		Respuesta += "ERROR: EL ID INGRESADO NO EXISTE"
	}

	
	return Respuesta
}

 
func disk(path string, id string)string{
	var Respuesta string
	var pathDico string
	Valido := false

	 
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			Valido = true
		}
	}

	if Valido{
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]
		
		tmp = strings.Split(pathDico, "/")
		NOmbreDis := strings.Split(tmp[len(tmp)-1], ".")[0]
		
		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			Respuesta += "ERROR REP DISK Open "+ err.Error()	
			return Respuesta	
		}

		var TempMBR Structs.MBR
		 
		if err := Herramientas.ReadObject(file, &TempMBR, 0); err != nil {
			Respuesta += "ERROR REP READ Open "+ err.Error()
			return Respuesta	
		}

		defer file.Close()

		 
		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n<tr> \n"
		cad += " <td bgcolor='SlateBlue'  ROWSPAN='3'> MBR </td>\n"
		cad += Structs.RepDiskGraphviz(TempMBR, file)
		cad += "\n</table> > ]\n}"

		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"

		 

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
		Respuesta += "Reporte Disk del disco "+NOmbreDis+" creado con el nombre "+nombre+".png"
	}else{
		Respuesta += "ERROR: EL ID INGRESADO NO EXISTE"
	}
	
	return Respuesta
}

 
func superBloque (path string, id string) string{
	var respuesta string
	var pathDico string
	reportar := false

	 
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	 
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar{
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()
		}

		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()
		}

		 
		defer file.Close()

		 
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				reportar = true
				part = i
				break  
			}
		}

		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += " <tr>\n  <td bgcolor='darkgreen' COLSPAN=\"2\"> <font color='white'> Reporte SUPERBLOQUE </font> </td> \n </tr> \n"
		cad += Structs.RepSB(mbr.Partitions[part], file)
		cad += "</table> > ]\n}"

		 
		carpeta := filepath.Dir(path)
		rutaReporte := carpeta + "/" + nombre + ".dot"
		respuesta += "Reporte BM_Bloque " + nombre +" creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	}

	return respuesta
}

 
func BM_inode(path string, id string) string{
	var respuesta string
	var pathDico string
	reportar := false

	 
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	 
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar{
		 
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()
		}

		 
		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()
		}

		 
		defer file.Close()

		 
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				reportar = true
				part = i
				break  
			}
		}

		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			 
			return "REP Error. Particion sin formato"
		}

		cad := ""
		inicio := superBloque.S_bm_inode_start
		fin := superBloque.S_bm_block_start
		count := 1  

		 
		var bm Structs.Bite

		for i := inicio; i < fin; i++ {
			 
			Herramientas.ReadObject(file, &bm, int64(i))

			if bm.Val[0] == 0 {
				cad += string("0 ")
			} else {
				cad += "1 "
			}

			if count == 20 {
				cad += "\n"
				count = 0
			}

			count++
		}

		 
		carpeta := filepath.Dir(path) 
		rutaReporte := carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, cad)
		respuesta += "Reporte BM Inode " + nombre +" creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco
	}

	return respuesta
}

 
func BM_Bloque(path string, id string) string{
	var respuesta string
	var pathDico string
	reportar := false

	 
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	 
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar{
		 
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()
		}

		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()
		}

		 
		defer file.Close()

		 
		part := -1
		for i := 0; i < 4; i++ {		
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				part = i
				break  
			}
		}

		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			 
			return "REP Error. Particion sin formato"
		}

		cad := ""
		inicio := superBloque.S_bm_block_start
		fin := superBloque.S_inode_start
		count := 1  

		 
		var bm Structs.Bite

		for i := inicio; i < fin; i++ {
			 
			Herramientas.ReadObject(file, &bm, int64(i))

			if bm.Val[0] == 0 {
				cad += string("0 ")
			} else {
				cad += "1 "
			}

			if count == 20 {
				cad += "\n"
				count = 0
			}

			count++
		}


		 
		carpeta := filepath.Dir(path) 
		rutaReporte := carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, cad)		
		respuesta += "Reporte BM_Bloque " + nombre +" creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco
	}
	return respuesta
}

 
func FILE(path string, id string, rutaFile string)string{
	var respuesta string
	var pathDico string
	var contenido string
	reportar := false

	 
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	 
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar{
		 
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		Disco, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP OPEN FILE "+err.Error()
		}

		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
			return "ERROR REP READ FILE "+err.Error()
		}

		 
		defer Disco.Close()

		 
		part := -1
		for i := 0; i < 4; i++ {		
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				part = i
				break  
			}
		}
		
		var superBloque Structs.Superblock
		var fileBlock Structs.Fileblock
		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			 
			return "REP Error. Particion sin formato"
		}

		 
		idInodo := toolsinodos.BuscarInodo(0, rutaFile, superBloque, Disco)
		var inodo Structs.Inode

		 
		if idInodo > 0 {
			contenido += "Contenido del archivo: '"+rutaFile+"'\n"
			Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
			 
			for _, idBlock := range inodo.I_block {
				if idBlock != -1 {
					Herramientas.ReadObject(Disco, &fileBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
					tmpConvertir := Herramientas.EliminartIlegibles(string(fileBlock.B_content[:]))
					contenido += tmpConvertir				
				}
			}

			contenido += "\n"
			
		} else {
			 
			return "REP ERROR: No se encontro el archivo " + rutaFile
		}

		 
		carpeta := filepath.Dir(path) 
		rutaReporte := carpeta + "/" + nombre + ".txt"
		Herramientas.Reporte(rutaReporte, contenido)
		respuesta += "Reporte BM_Bloque " + nombre +" creado \n"
		respuesta += "Pertenece al disco: " + nombreDisco
	}
	return respuesta
}

 
func LS(path string, id string, rutaFile string)string{
	var respuesta string
	var contenido string
	var pathDico string
	reportar := false

	 
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	 
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar{
		Color := "BlueViolet"	
		contenido = "digraph {\nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n<tr>\n\t<td bgcolor='"+Color+"'>PERMISOS</td>\n\t<td bgcolor='"+Color+"'> USUARIO </td>\n\t<td bgcolor='"+Color+"'> GRUPO </td>\n\t<td bgcolor='"+Color+"'> SIZE </td>\n\t<td bgcolor='"+Color+"'> FECHA/HORA </td> \n\t<td bgcolor='"+Color+"'> NOMBRE </td>\n\t<td bgcolor='"+Color+"'> TIPO </td>\n </tr>"
		
		 
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		Disco, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP OPEN FILE "+err.Error()
		}

		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(Disco, &mbr, 0); err != nil {
			return "ERROR REP READ FILE "+err.Error()
		}

		 
		defer Disco.Close()

		 
		part := -1
		for i := 0; i < 4; i++ {		
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				part = i
				break  
			}
		}
		
		 
		var superBloque Structs.Superblock
		errREAD := Herramientas.ReadObject(Disco, &superBloque, int64(mbr.Partitions[part].Start))
		if errREAD != nil {
			 
			return "CAT ERROR. Particion sin formato" + "\n"
		}


		var FstInodo Structs.Inode		
		 
		Herramientas.ReadObject(Disco, &FstInodo, int64(superBloque.S_inode_start + int32(binary.Size(Structs.Inode{}))))
			

		var contUs string
		var FistfileBlock Structs.Fileblock
		for _, item := range FstInodo.I_block {
			if item != -1 {
				Herramientas.ReadObject(Disco, &FistfileBlock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
				contUs += string(FistfileBlock.B_content[:])
			}
		}
		lineaID := strings.Split(contUs, "\n")
		

		idInodo := toolsinodos.BuscarInodo(0, rutaFile, superBloque, Disco)
		var inodo Structs.Inode

		if idInodo > 0 {
			Herramientas.ReadObject(Disco, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
			var folderBlock Structs.Folderblock
			for _, idBlock := range inodo.I_block {
				if idBlock != -1 {
					Herramientas.ReadObject(Disco, &folderBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Folderblock{})))))					
					for k := 2; k < 4; k++ {
						apuntador := folderBlock.B_content[k].B_inodo
						if apuntador != -1 {
							pathActual := Structs.GetB_name(string(folderBlock.B_content[k].B_name[:]))
							
							contenido += InodoLs(pathActual, lineaID, apuntador , superBloque, Disco)
						}
					}					
				}
			}
			
			
		}else{
			respuesta = "REP ERROR NO SE ENCONTRO LA PATH INGRESADA"
		}

		contenido += "\n</table> > ]\n}"
		cad := Herramientas.EliminartIlegibles(contenido)

		 
		carpeta := filepath.Dir(path) 
		rutaReporte := carpeta + "/" + nombre + ".dot"
		Herramientas.Reporte(rutaReporte, contenido)
		respuesta += "Reporte BM_Bloque " + nombre +" creado \n"
		respuesta += "Pertenece al disco: " + nombreDisco
		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	}	
	return respuesta	
}

			 
func InodoLs(name string,lineaID []string,  idInodo int32, superBloque Structs.Superblock, file *os.File)string{
	var contenido string

	 
	var inodo Structs.Inode
	Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))
	
	 
	usuario:= ""
	grupo:=""							
	for m:=0; m<len(lineaID); m++{
		datos := strings.Split(lineaID[m], ",")
		if len(datos) == 5 {	
			us := fmt.Sprintf("%d",inodo.I_uid)													
			if us== datos[0]{
				usuario = datos[3]
			}		
		}
		if len(datos) == 3 {	
			gr := fmt.Sprintf("%d",inodo.I_gid)									
			if gr== (datos[0]){
				grupo = datos[2]
			}		
		}

	}
	
	Color := "Pink"
	tipoArchivo := "Archivo"
	var permisos string	
	
	 
	 
	 
	 
	 
	 
	 
	 
	 
	for i:=0; i<3; i++{	
		if string(inodo.I_perm[i])=="0"{ 
			permisos+="---"
		}else if string(inodo.I_perm[i])=="1"{ 
			permisos+="--x"
		}else if string(inodo.I_perm[i])=="2"{ 
			permisos+="-w-"
		}else if string(inodo.I_perm[i])=="3"{ 
			permisos+="-wx"
		}else if string(inodo.I_perm[i])=="4"{ 
			permisos+="r--"
		}else if string(inodo.I_perm[i])=="5"{ 
			permisos+="r-x"
		}else if string(inodo.I_perm[i])=="6"{ 
			permisos+="rw-"
		}else if string(inodo.I_perm[i])=="7"{ 
			permisos+="rwx"
		}
	}

	if string(inodo.I_type[:]) == "0"{
		Color = "Violet"
		tipoArchivo = "Carpeta"	
	}
	
	permisos = "rw-rw-r--"	
	contenido += "\n  <tr>"
	contenido += "\n\t <td bgcolor='"+Color+"'> "+ permisos +"</td>"
	contenido += fmt.Sprintf("\n\t <td bgcolor='%s'> %s</td>",Color,usuario)
	contenido += fmt.Sprintf("\n\t <td bgcolor='%s'> %s</td>",Color,grupo)
	contenido += fmt.Sprintf("\n\t <td bgcolor='%s'> %d</td>", Color, inodo.I_size)
	contenido += fmt.Sprintf("\n\t <td bgcolor='%s'> %s </td> ", Color, string(inodo.I_ctime[:]))
	contenido += "\n\t <td bgcolor='"+Color+"'> "+ name +"</td>"
	contenido += "\n\t <td bgcolor='"+Color+"'> "+ tipoArchivo +"</td>"
	contenido += "\n  </tr>"
	 
	return contenido
}

 
func journal(path string, id string) string{
	var respuesta string
	var pathDico string
	reportar := false

	 
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	 
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	 
	if reportar {
		 
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()
		}

		 
		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()
		}

		 
		defer file.Close()

		 
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				reportar = true
				part = i
				break  
			}
		}

		cad := "digraph { \nnode [ shape=none ] \nTablaReportNodo [ label = < <table border=\"1\"> \n"
		cad += Structs.RepJournal(mbr.Partitions[part], file)
		cad += "</table> > ]\n}"

		 
		carpeta := filepath.Dir(path) 
		rutaReporte := carpeta + "/" + nombre + ".dot"
		Herramientas.Reporte(rutaReporte, cad)
		respuesta += "Reporte BM Inode " + nombre +" creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	} else {
		 
	}
	return respuesta
}

 
func tree(path string, id string) string{
	var respuesta string
	var pathDico string
	reportar := false

	 
	for _,montado := range Structs.Montadas{
		if montado.Id == id{
			pathDico = montado.PathM
			reportar = true
		}
	}

	 
	if pathDico == ""{
		reportar = false
		return "ERROR REP: ID NO ENCONTRADO"
	}

	if reportar {
		 
		tmp := strings.Split(path, "/")
		nombre := strings.Split(tmp[len(tmp)-1], ".")[0]

		tmp2 := strings.Split(pathDico, "/")
		nombreDisco := strings.Split(tmp2[len(tmp2)-1], ".")[0]

		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR REP SB OPEN FILE "+err.Error()
		}

		 
		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR REP SB READ FILE "+err.Error()
		}

		 
		defer file.Close()

		 
		part := -1
		for i := 0; i < 4; i++ {
			identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
			if identificador == id {
				reportar = true
				part = i
				break  
			}
		}

		var superBloque Structs.Superblock
		err = Herramientas.ReadObject(file, &superBloque, int64(mbr.Partitions[part].Start))
		if err != nil {
			 
			return "REP Error. Particion sin formato"
		}

		var Inode0 Structs.Inode
		Herramientas.ReadObject(file, &Inode0, int64(superBloque.S_inode_start))

		cad := "digraph { \n graph [pad=0.5, nodesep=0.5, ranksep=1] \n node [ shape=plaintext ] \n rankdir=LR \n"

		 
		cad += "\n Inodo0 [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n"
		cad += "    <tr> <td bgcolor='skyblue' colspan=\"2\" port='P0'> Inodo 0 </td> </tr> \n"

		for i := 0; i < 12; i++ {
			cad += fmt.Sprintf("    <tr> <td> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode0.I_block[i])
		}
		 
		for i := 12; i < 15; i++ {
			cad += fmt.Sprintf("    <tr> <td bgcolor='pink'> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode0.I_block[i])
		}
		cad += "   </table> \n  > \n ]; \n"
		 

		 
		for i := 0; i < 15; i++ {
			bloque := Inode0.I_block[i]
			if bloque != -1 {
				 
				cad += treeBlock(bloque, string(Inode0.I_type[:]), 0, i+1, superBloque, file)
			}
		}
		 
		 
		 
		cad += "\n}"

		 
		carpeta := filepath.Dir(path) 
		rutaReporte := carpeta + "/" + nombre + ".dot"
		Herramientas.Reporte(rutaReporte, cad)
		respuesta += "Reporte tree " + nombre +" creado \n"
		respuesta += " Pertenece al disco: " + nombreDisco

		Herramientas.RepGraphizMBR(rutaReporte, cad, nombre)
	} else {
		 
	}
	return respuesta
}

 
 
func treeBlock(idBloque int32, tipo string, idPadre int32, p int, superBloque Structs.Superblock, file *os.File) string {
	cad := fmt.Sprintf("\n Bloque%d [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n", idBloque)

	if tipo == "0" {
		 
		var folderBlock Structs.Folderblock
		Herramientas.ReadObject(file, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))

		 
		cad += fmt.Sprintf("    <tr> <td bgcolor='orchid' colspan=\"2\" port='P0'> Bloque %d </td> </tr> \n", idBloque)
		cad += fmt.Sprintf("    <tr> <td> . </td> <td port='P1'> %d </td> </tr> \n", folderBlock.B_content[0].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> .. </td> <td port='P2'> %d </td> </tr> \n", folderBlock.B_content[1].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> %s </td> <td port='P3'> %d </td> </tr> \n", Structs.GetB_name(string(folderBlock.B_content[2].B_name[:])), folderBlock.B_content[2].B_inodo)
		cad += fmt.Sprintf("    <tr> <td> %s </td> <td port='P4'> %d </td> </tr> \n", Structs.GetB_name(string(folderBlock.B_content[3].B_name[:])), folderBlock.B_content[3].B_inodo)
		cad += "   </table> \n  > \n ]; \n"
		 
		cad += fmt.Sprintf("\n Inodo%d:P%d -> Bloque%d:P0; \n", idPadre, p, idBloque)  
		 
		for i := 2; i < 4; i++ {
			inodo := folderBlock.B_content[i].B_inodo
			if inodo != -1 {
				 
				cad += treeInodo(inodo, idBloque, i+1, superBloque, file)
			}
		}
	} else {
		 
		var fileBlock Structs.Fileblock
		Herramientas.ReadObject(file, &fileBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Fileblock{})))))
		 
		cad += fmt.Sprintf("    <tr> <td bgcolor='#ffff99' port='P0'> Bloque %d </td> </tr> \n", idBloque)
		cad += fmt.Sprintf("    <tr> <td> %s </td> </tr> \n", Structs.GetB_content(string(fileBlock.B_content[:])))
		cad += "   </table> \n  > \n ]; \n"
		 
		cad += fmt.Sprintf("\n Inodo%d:P%d -> Bloque%d:P0; \n", idPadre, p, idBloque)  
	}

	return cad
}

 
 
func treeInodo(idInodo int32, idPadre int32, p int, superBloque Structs.Superblock, file *os.File) string {
	 
	var Inode Structs.Inode
	Herramientas.ReadObject(file, &Inode, int64(superBloque.S_inode_start+(idInodo*int32(binary.Size(Structs.Inode{})))))

	 
	cad := fmt.Sprintf("\n Inodo%d [ \n  label = < \n   <table border=\"0\" cellborder=\"1\" cellspacing=\"0\"> \n", idInodo)
	 
	if string(Inode.I_type[:]) == "0" {
		cad += fmt.Sprintf("    <tr> <td bgcolor='skyblue' colspan=\"2\" port='P0'> Inodo %d </td> </tr> \n", idInodo)
	} else {
		cad += fmt.Sprintf("    <tr> <td bgcolor='#7FC97F' colspan=\"2\" port='P0'> Inodo %d </td> </tr> \n", idInodo)
	}

	 
	for i := 0; i < 12; i++ {
		cad += fmt.Sprintf("    <tr> <td> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode.I_block[i])
	}
	 
	for i := 12; i < 15; i++ {
		cad += fmt.Sprintf("    <tr> <td bgcolor='pink'> AD%d </td> <td port='P%d'> %d </td> </tr> \n", i+1, i+1, Inode.I_block[i])
	}
	cad += "   </table> \n  > \n ]; \n"
	 

	 
	cad += fmt.Sprintf("\n Bloque%d:P%d -> Inodo%d:P0; \n", idPadre, p, idInodo)  

	 
	for i := 0; i < 15; i++ {
		bloque := Inode.I_block[i]
		if bloque != -1 {
			 
			cad += treeBlock(bloque, string(Inode.I_type[:]), idInodo, i+1, superBloque, file)
		}
	}

	return cad
}