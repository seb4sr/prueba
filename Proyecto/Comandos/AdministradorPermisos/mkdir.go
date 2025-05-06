package administradorpermisos

import (
	ToolsInodos "MIA_2S_P2_201513656/ToolsInodos"
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	//"fmt"
	"strings"
)

func Mkdir(entrada []string) string{
	respuesta := "Comando mkdir"	
	var path string
	p := false
	UsuarioA := Structs.UsuarioActual

	if !UsuarioA.Status {
		 
		respuesta += "ERROR MKFILE: NO HAY SECION INICIADA" + "\n"
		respuesta += "POR FAVOR INICIAR SESION PARA CONTINUAR" + "\n"
		return respuesta
	}

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if strings.ToLower(valores[0]) == "path" {
			if len(valores)!=2{
				 
				respuesta += "ERROR MKDIR, valor desconocido de parametros " + valores[1]
				 
				return respuesta
			}			
			path = strings.ReplaceAll(valores[1],"\"","")
		} else if strings.ToLower(valores[0]) == "r" {
			if len(tmp) != 1 {
				 
				return "MKDIR Error: Valor desconocido del parametro "+ valores[0]
			}
			p = true

			 
		} else {
			 
			return "MKFILE ERROR: Parametro desconocido: " + valores[0]
		}
	}

	if path ==""{
		 
		return "MKDIR ERROR NO SE INGRESO PARAMETRO PATH"
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
			 
			return "MKFILE ERROR. Particion sin formato"+ "\n"
		}

		 
		stepPath := strings.Split(path, "/")
		idInicial := int32(0)
		idActual := int32(0)
		crear := -1
		for i, itemPath := range stepPath[1:] {
			idActual = ToolsInodos.BuscarInodo(idInicial, "/"+itemPath, superBloque, Disco)
			if idInicial != idActual {
				idInicial = idActual
			} else {
				crear = i + 1  
				break
			}
		}

		 
		if crear != -1 {
			if crear == len(stepPath)-1 {
				ToolsInodos.CreaCarpeta(idInicial, stepPath[crear], int64(mbr.Partitions[part].Start), Disco)
			} else {
				if p {
					for _, item := range stepPath[crear:] {
						idInicial = ToolsInodos.CreaCarpeta(idInicial, item, int64(mbr.Partitions[part].Start), Disco)
						if idInicial == 0 {
							 
							return "MKDIR ERROR: No se pudo crear carpeta"
						}
					}
				} else {
					 
				}
			}
			return "Carpeta(s) creada"
		}else{
			 
			return "MKDIR ERROR: LA CARPETA YA EXISTE"
		}
	}
	return respuesta
}