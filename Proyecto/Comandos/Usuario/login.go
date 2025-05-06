package usuario

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	"encoding/binary"
	//"fmt"
	"strconv"
	"strings"
)

func Login(entrada []string) (string, int){
	var respuesta string
	var user string  
	var pass string  
	var id string    
	Valido := true
	var pathDico string

	if Structs.UsuarioActual.Status {
		Valido = false
		return "LOGIN ERROR: Ya existe una sesion iniciada, cierre sesion para iniciar otra",0
	}

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			 
			respuesta += "ERROR LOGIN, valor desconocido de parametros " + valores[1]+ "\n"
			Valido = false
			 
			return respuesta, 5
		}

		 
		if strings.ToLower(valores[0]) == "id" {
			id = strings.ToUpper(valores[1])

		 
		} else if strings.ToLower(valores[0]) == "user" {
			user = valores[1]

		 
		} else if strings.ToLower(valores[0]) == "pass" {
			pass = valores[1]

		 
		} else {
			 
			Valido = false
			 
			return "LOGIN ERROR: Parametro desconocido: "+valores[0] + "\n",5
		}
	}

	 
	if id != ""{
		 
		for _,montado := range Structs.Montadas{
			if montado.Id == id{				
				pathDico = montado.PathM
			}
		}
		if pathDico == ""{
			Valido = false
			return "ERROR LOGIN: ID NO ENCONTRADO"+ "\n",5
		}
	}else{
		 
		Valido = false
		return "LOGIN ERROR: FALTO EL PARAMETRO ID "+ "\n",5
	}

	if pass==""{
		 
		Valido = false
		return "LOGIN ERROR: FALTO EL PARAMETRO PASS "+ "\n",5
	}

	if user==""{
		 
		Valido = false
		return "LOGIN ERROR: FALTO EL PARAMETRO USER "+ "\n",5
	}

	if Valido{
		file, err := Herramientas.OpenFile(pathDico)
		if err != nil {
			return "ERROR LOGIN OPEN FILE "+err.Error()+ "\n",5
		}

		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(file, &mbr, 0); err != nil {
			return "ERROR LOGIN READ FILE "+err.Error()+ "\n",5
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
			 
			return "REP Error. Particion sin formato"+ "\n", 1
		}

		var inodo Structs.Inode		
		 
		Herramientas.ReadObject(file, &inodo, int64(superBloque.S_inode_start + int32(binary.Size(Structs.Inode{}))))
		
		 
		var contenido string
		var fileBLock Structs.Fileblock
		for _, item := range inodo.I_block {
			if item != -1 {
				Herramientas.ReadObject(file, &fileBLock, int64(superBloque.S_block_start+(item*int32(binary.Size(Structs.Fileblock{})))))
				contenido += string(fileBLock.B_content[:])
			}
		}

		 
		linea := strings.Split(contenido, "\n")
		 

		logeado := false
		for _,reglon := range linea{
			Usuario := strings.Split(reglon, ",")
			
			 
			if len(Usuario) == 5{
				if Usuario[0]!="0"{
					if Usuario[3] == user{
						if Usuario[4]== pass{
							Structs.UsuarioActual.IdPart = id
							Structs.UsuarioActual.Nombre = user
							Structs.UsuarioActual.Status = true
							Structs.UsuarioActual.PathD = pathDico
							Add_idUsr(Usuario[0])
							logeado = true		
							Search_IdGrp(linea, Usuario[2])					
						}else{
							 
							return "ERROR LOGIN: LA CONTRASEÑA ES INCORRECTA", 3
						}
						break
					}
				}
			}
		}

		if logeado{
			respuesta = "EL ususario '"+ user +"' ha iniciado sesion exitosamente! \n"
			 
			return respuesta,-1
		}else{
			 
			respuesta += "ERROR AL INTENTAR INGRESAR, NO SE ENCONTRO EL USUARIO \n"
			respuesta+= "POR FAVOR INGRESE LOS DATOS CORRECTOS \n"
			return respuesta, 4
		}
	}

	return respuesta,5
}

func Add_idUsr(id string) string{
	idU, errId := strconv.Atoi(id)
	if errId != nil {
		 
		return "LOGIN ERROR: Error desconcocido con el idUsr"
	}
	Structs.UsuarioActual.IdUsr = int32(idU)
	return ""
}

func Search_IdGrp(lineaID []string, grupo string) {
	for _, registro := range lineaID[:len(lineaID)-1] {
		datos := strings.Split(registro, ",")
		if len(datos) == 3 {
			if datos[2] == grupo {
				 
				id, errId := strconv.Atoi(datos[0])
				if errId != nil {
					 
					return
				}
				Structs.UsuarioActual.IdGrp = int32(id)
				return
			}
		}
	}
}