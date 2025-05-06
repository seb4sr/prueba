package main

import (
	AD "MIA_2S_P2_201513656/Comandos/AdministradorDiscos"
	AP "MIA_2S_P2_201513656/Comandos/AdministradorPermisos"
	REP "MIA_2S_P2_201513656/Comandos/Rep"
	SA "MIA_2S_P2_201513656/Comandos/SistemaArchivos"
	USR "MIA_2S_P2_201513656/Comandos/Usuario"
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	ToolsInodos "MIA_2S_P2_201513656/ToolsInodos"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/cors"
)

type Entrada struct {
	Text string `json:"text"`
}

type Login struct {
	User string `json:"usuario"`
	Pass string `json:"password"`
	Id   string `json:"id"`
}

type StatusResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type RespuestaDisk struct{
	Discos []string//Contiene los Discos de la carpeta
	DiskPart string//es el id del disco donde de inicio sesion
}

type RespuestaPart struct{
	Partciones []string//partciones del Disco donde se incio sesion
	IdParticion string//id de la particion donde se inicio sesion
}

func main()  {
	//EndPoint 
	//metodos de uso
	http.HandleFunc("/analizar", getCadenaAnalizar)
	http.HandleFunc("/discos", getDiscos)
	http.HandleFunc("/particiones", getParticiones)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/explorer", getContenido)
	http.HandleFunc("/contenido", getContenidoR)
	http.HandleFunc("/file", getFile)
	http.HandleFunc("/back", getBack)	

	// Configurar CORS con opciones predeterminadas
	//Permisos para enviar y recir informacion
	c := cors.Default()

	// Configurar el manejador HTTP con CORS
	handler := c.Handler(http.DefaultServeMux)

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor escuchando en http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}

func getCadenaAnalizar(w http.ResponseWriter, r *http.Request) {
	var respuesta string
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")
	
	var status StatusResponse
	//verificar que sea un post
	if r.Method == http.MethodPost {
		var entrada Entrada
		if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
			http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
			status = StatusResponse{Message: "Error al decodificar JSON", Type: "unsucces"}
			json.NewEncoder(w).Encode(status)
			return
		}
		
		//creo un lector de bufer para el archivo
		lector := bufio.NewScanner(strings.NewReader(entrada.Text))
		//leer el archivo linea por linea
		for lector.Scan() {
			//Elimina los saltos de linea
			if lector.Text() != ""{
				//Divido por # para ignorar todo lo que este a la derecha del mismo
				linea := strings.Split(lector.Text(), "#") //lector.Text() retorna la linea leida
				if len(linea[0]) != 0 {
					 
					 
					respuesta += "==========================================================================================\n"
					respuesta += "Comando en ejecucion: " + linea[0] + "\n"
					respuesta += Analizar(linea[0])  + "\n"
				}	
				//Comentarios			
				if len(linea) > 1 && linea[1] != "" {
					 
					respuesta += "#"+linea[1] +"\n"
				}
			}
			
		}

		// 
		w.WriteHeader(http.StatusOK)

		status = StatusResponse{Message: respuesta, Type: "succes"}
		json.NewEncoder(w).Encode(status)

	} else {
		//http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		status = StatusResponse{Message: "Metodo no permitido", Type: "unsucces"}
		json.NewEncoder(w).Encode(status)
	}
}

func Analizar(entrada string)string{
	tmp := strings.TrimRight(entrada," ")
	//Recibe una linea y la descompone entre el comando y sus parametros
	parametros:= strings.Split(tmp, " -")

	// *============================* ADMINISTRACION DE DISCOS *============================*
	//mkdisk -size=5 -unit=M -path=Calificacion_MIA/Discos/Disco1.mia
	//mkdisk -size=5 -unit=M -path="Calificacion_MIA/Discos/Disco_1.mia"
	if strings.ToLower(parametros[0])=="mkdisk"{
		if len(parametros)>1{	
			AD.Mkdisk(parametros)				
			return AD.Mkdisk(parametros)
		}else{
			 
			return "ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK"
		}

	//rmdisk -path="/home/mis discos/Disco4.mia"
	}else if strings.ToLower(parametros[0])=="rmdisk"{
		if len(parametros)>1{	
			return AD.Rmdisk(parametros)		
		}else{
			 
			return  "ERROR EN RMDISK, FALTAN PARAMETROS EN MKDISK"
		}

	//fdisk -type=P -unit=b -name=Part1 -size=10485760 -path=Calificacion_MIA/Discos/Disco1.mia
	//fdisk -add=-23760 -path=Calificacion_MIA/Discos/Disco1.mia -name=Part1 -size=10485760
	}else if strings.ToLower(parametros[0])=="fdisk"{
		if len(parametros)>1{	
			return AD.Fdisk(parametros)		
		}else{
			 
			return  "ERROR EN FDISK, FALTAN PARAMETROS EN MKDISK"
		}

	//mount -path=/home/Disco3.mia -name=Part1
	}else if strings.ToLower(parametros[0])=="mount"{
		if len(parametros)>1{	
			return AD.Mount(parametros)			
		}else{
			 
			return "ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK"
		}

	//unmount -id=561A
	}else if strings.ToLower(parametros[0])=="unmount"{
		if len(parametros)>1{	
			return AD.Unmoun(parametros)		
		}else{
			 
			return  "ERROR EN MKDISK, FALTAN PARAMETROS EN MKDISK"
		}
	// *===================* ADMINISTRACION DE SISTEMA DE ARCHIVOS *======================*
	//ejm: mkfs -type=full -id=341A	-fs=3fs
	}else if strings.ToLower(parametros[0])=="mkfs"{		
		if len(parametros)>1{			
			return SA.MKfs(parametros)
		}else{
			 
			return "ERROR EN MKFS, FALTAN PARAMETROS"
		}
	
	// *===================* ADMINISTRACION DE USUARIOS Y CARPETAS *======================*
	//login -user=root -pass=123 -id=561A
	}else if strings.ToLower(parametros[0])=="login"{		
		if len(parametros)>1{	
			tmp,_ := USR.Login(parametros)		
			return tmp
		}else{
			 
			return "ERROR EN LOGIN, FALTAN PARAMETROS"
		}

	//logout
	}else if strings.ToLower(parametros[0])=="logout"{
		str, _ := USR.Logout()	
		return str
	
	//mkgrp -name=usuarios
	}else if strings.ToLower(parametros[0])=="mkgrp"{		
		if len(parametros)>1{			
			return USR.Mkgrp(parametros)
		}else{
			 
			return "ERROR EN MKGRP, FALTAN PARAMETROS"
		}
		
	//rmgrp -name=usuarios
	}else if strings.ToLower(parametros[0])=="rmgrp"{		
		if len(parametros)>1{			
			return USR.Rmgrp(parametros)
		}else{
			 
			return "ERROR EN RMGRP, FALTAN PARAMETROS"
		}
	
	}else if strings.ToLower(parametros[0])=="mkusr"{		
		if len(parametros)>1{			
			return USR.Mkusr(parametros)
		}else{
			 
			return "ERROR EN RMGRP, FALTAN PARAMETROS"
		}	

	}else if strings.ToLower(parametros[0])=="rmusr"{		
		if len(parametros)>1{			
			return USR.Rmusr(parametros)
		}else{
			 
			return "ERROR EN RMUSR, FALTAN PARAMETROS"
		}

	}else if strings.ToLower(parametros[0])=="chgrp"{		
		if len(parametros)>1{			
			return USR.Chgrp(parametros)
		}else{
			 
			return "ERROR EN CHGRP, FALTAN PARAMETROS"
		}
	
	// *=======================* PERMISOS DE CARPETAS Y ARHICVOS *============================*
	//mkfile -path=/home/archivos/docs/Tarea2.txt -size=75 -r
	}else if strings.ToLower(parametros[0])=="mkfile"{		
		if len(parametros)>1{			
			return AP.MKfile(parametros)
		}else{
			 
			return "ERROR EN MKFILE, FALTAN PARAMETROS"
		}
	
	//EJ: cat -file1=/home/user/docs/a.txt -file12=/home/user/docs/b.txt
	}else if strings.ToLower(parametros[0])=="cat"{		
		if len(parametros)>1{			
			return AP.Cat(parametros)
		}else{
			 
			return "ERROR EN CAT, FALTAN PARAMETROS"
		}

	//mkdir -r -path=/home/archivos/Fotos
	}else if strings.ToLower(parametros[0])=="mkdir"{		
		if len(parametros)>1{			
			return AP.Mkdir(parametros)
		}else{
			 
			return "ERROR EN MKDIR, FALTAN PARAMETROS"
		}
	
	}else if strings.ToLower(parametros[0])=="rename"{		
		if len(parametros)>1{			
			return AP.Rename(parametros)
		}else{
			 
			return "ERROR EN RENAME, FALTAN PARAMETROS"
		}
	
	}else if strings.ToLower(parametros[0])=="edit"{		
		if len(parametros)>1{			
			return AP.Edit(parametros)
		}else{
			 
			return "ERROR EN EDIT, FALTAN PARAMETROS"
		}
	
	}else if strings.ToLower(parametros[0])=="copy"{		
		if len(parametros)>1{			
			return AP.Copy(parametros)
		}else{
			 
			return "ERROR EN COPY, FALTAN PARAMETROS"
		}
	// *============================* OTROS *============================*
	} else if strings.ToLower(parametros[0]) == "rep" {
		//REP
		if len(parametros) > 1 {
			return REP.Rep(parametros)
		} else {
			 
			return "REP ERROR: parametros no encontrados"
		}
	} else if strings.ToLower(parametros[0]) == "" {
		//para agregar lineas con cada enter sin tomarlo como error
		return ""
	} else {
		 
		return "ERROR: COMANDO "+parametros[0]+" NO RECONOCIBLE"
	}
}

func getDiscos(w http.ResponseWriter, r *http.Request) {
	 
	 
	//struc de respuesta
	var disk RespuestaDisk
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")

	var entradaId string
	if err := json.NewDecoder(r.Body).Decode(&entradaId); err != nil {
		http.Error(w, "Error decodificion JSON", http.StatusBadRequest)
		return
	}
	 

	var pathDico string
	//BUsca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == entradaId{
			pathDico = montado.PathM			
		}
	}

	directorio := filepath.Dir(pathDico)
	tmp := strings.Split(pathDico, "/")
	nombre := tmp[len(tmp)-1]
	disk.DiskPart = nombre

	 

	//recorrer el directorio y buscar discos
	err := filepath.Walk(directorio, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			disk.Discos = append(disk.Discos, info.Name())
		}
		return nil
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al buscar archivos: %s", err), http.StatusInternalServerError)
	}

	respuestaJSON, err := json.Marshal(disk)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar datos a JSON: %s", err), http.StatusInternalServerError)
		return
	}
	w.Write(respuestaJSON)
}

func getParticiones(w http.ResponseWriter, r *http.Request) {
	 
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")

	var entrada string
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	var pathDico string
	//BUsca en struck de particiones montadas el id ingresado
	for _,montado := range Structs.Montadas{
		if montado.Id == entrada{
			pathDico = montado.PathM	
			break;		
		}
	}

	disco, err := Herramientas.OpenFile(pathDico)
	if err != nil {
		 
		return
	}

	//Se crea un mbr para cargar el mbr del disco
	var mbr Structs.MBR
	//Guardo el mbr leido
	if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
		return
	}

	// cerrar el archivo del disco
	defer disco.Close()

	//lista de discos encontrados	
	var resPartciones RespuestaPart

	for i := 0; i < 4; i++ {
		estado := string(mbr.Partitions[i].Status[:])
		id := string(mbr.Partitions[i].Id[:])
		if estado == "A" {
			resPartciones.Partciones = append(resPartciones.Partciones, string(mbr.Partitions[i].Id[:]))
		}
		if id == entrada{
			resPartciones.IdParticion = id
		}
	}

	 
	respuestaJSON, err := json.Marshal(resPartciones)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar datos a JSON: %s", err), http.StatusInternalServerError)
		return
	}
	w.Write(respuestaJSON)
}

func login(w http.ResponseWriter, r *http.Request) {	
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")

	var entrada Login
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	 
	 
	//Construir cadena para ejecutar el comando login
	//login -user=root -pass=123 -id=A148

	logear := [4]string{"login", "user=" + entrada.User, "pass=" + entrada.Pass, "id=" + entrada.Id}
	 
	_,tmp:= USR.Login(logear[:])

	respuestaJSON, err := json.Marshal(tmp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar datos a JSON: %s", err), http.StatusInternalServerError)
		return
	}
	w.Write(respuestaJSON)
}

func logout(w http.ResponseWriter, r *http.Request) {
	 
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")
	_,tmp :=USR.Logout()

	respuestaJSON, err := json.Marshal(tmp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar datos a JSON: %s", err), http.StatusInternalServerError)
		return
	}
	w.Write(respuestaJSON)
}

// variables para manejar los inodos (carpetas y archivos)
var idActual int32
var initSuperBloque int64

func getContenido(w http.ResponseWriter, r *http.Request) {
	 
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")

	//Abrimos el disco
	disco, err := Herramientas.OpenFile(Structs.UsuarioActual.PathD)
	if err != nil {
		return 
	}

	//Se crea un mbr para cargar el mbr del disco
	var mbr Structs.MBR
	//Guardo el mbr leido
	if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
		return
	}

	// cerrar el archivo del disco
	defer disco.Close()

	for i := 0; i < 4; i++ {
		identificador := Structs.GetId(string(mbr.Partitions[i].Id[:]))
		if identificador == Structs.UsuarioActual.IdPart {
			initSuperBloque = int64(mbr.Partitions[i].Start)
			break //para que ya no siga recorriendo si ya encontro la particion independientemente si se pudo o no reducir
		}
	}

	var superBloque Structs.Superblock
	Herramientas.ReadObject(disco, &superBloque, initSuperBloque)

	var Inode0 Structs.Inode
	Herramientas.ReadObject(disco, &Inode0, int64(superBloque.S_inode_start))

	//establezco valores de id (como es raiz ambos seran 0)
	idActual = 0

	//lista de discos encontrados
	var contenido []string

	var folderBlock Structs.Folderblock
	for i := 0; i < 12; i++ {
		idBloque := Inode0.I_block[i]
		if idBloque != -1 {
			Herramientas.ReadObject(disco, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))
			//Recorrer el bloque actual buscando la carpeta/archivo en la raiz
			for j := 2; j < 4; j++ {
				//apuntador es el apuntador del bloque al inodo (carpeta/archivo), si existe es distinto a -1
				apuntador := folderBlock.B_content[j].B_inodo
				if apuntador != -1 {
					pathActual := Structs.GetB_name(string(folderBlock.B_content[j].B_name[:]))
					contenido = append(contenido, pathActual)
				}
			}
		}
	}

	respuestaJSON, err := json.Marshal(contenido)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar datos a JSON: %s", err), http.StatusInternalServerError)
		return
	}
	w.Write(respuestaJSON)
}

// para manejar el anterior
var listaAnterior []int32

func getContenidoR(w http.ResponseWriter, r *http.Request) {
	 
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")

	var entrada string
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	 

	//Abrimos el disco
	disco, err := Herramientas.OpenFile(Structs.UsuarioActual.PathD)
	if err != nil {
		return 
	}
	// Close bin file
	defer disco.Close()

	var superBloque Structs.Superblock
	Herramientas.ReadObject(disco, &superBloque, initSuperBloque)

	//agrego el actual a la pila de anteriores (este sera el anterior)
	listaAnterior = append(listaAnterior, idActual)
	//busco en el actual
	idActual = ToolsInodos.BuscarInodo(idActual, "/"+entrada, superBloque, disco)

	//cargo el inodo actual
	var Inode Structs.Inode
	Herramientas.ReadObject(disco, &Inode, int64(superBloque.S_inode_start+(idActual*int32(binary.Size(Structs.Inode{})))))

	//lista de discos encontrados
	var contenido []string

	var folderBlock Structs.Folderblock
	for i := 0; i < 12; i++ {
		idBloque := Inode.I_block[i]
		if idBloque != -1 {
			Herramientas.ReadObject(disco, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))
			//Recorrer el bloque actual buscando la carpeta/archivo en la raiz
			for j := 2; j < 4; j++ {
				//apuntador es el apuntador del bloque al inodo (carpeta/archivo), si existe es distinto a -1
				apuntador := folderBlock.B_content[j].B_inodo
				if apuntador != -1 {
					pathActual := Structs.GetB_name(string(folderBlock.B_content[j].B_name[:]))
					contenido = append(contenido, pathActual)
				}
			}
		}
	}

	respuestaJSON, err := json.Marshal(contenido)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar datos a JSON: %s", err), http.StatusInternalServerError)
		return
	}
	w.Write(respuestaJSON)
}

func getFile(w http.ResponseWriter, r *http.Request) {
	 
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")

	var entrada string
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	 

	//Abrimos el disco
	disco, err := Herramientas.OpenFile(Structs.UsuarioActual.PathD)
	if err != nil {
		return 
	}
	// Close bin file
	defer disco.Close()

	var superBloque Structs.Superblock
	Herramientas.ReadObject(disco, &superBloque, initSuperBloque)

	//agrego el actual a la pila de anteriores (este sera el anterior)
	listaAnterior = append(listaAnterior, idActual)
	//busco en el actual
	idActual = ToolsInodos.BuscarInodo(idActual, "/"+entrada, superBloque, disco)

	//cargo el inodo actual
	var Inode Structs.Inode
	Herramientas.ReadObject(disco, &Inode, int64(superBloque.S_inode_start+(idActual*int32(binary.Size(Structs.Inode{})))))

	//lista de discos encontrados
	var contenido []string
	var textFile string
	var fileBlock Structs.Fileblock

	for _, idBlock := range Inode.I_block {
		if idBlock != -1 {
			Herramientas.ReadObject(disco, &fileBlock, int64(superBloque.S_block_start+(idBlock*int32(binary.Size(Structs.Fileblock{})))))
			textFile += Herramientas.EliminartIlegibles(string(fileBlock.B_content[:])) 
		}
	}

	contenido = append(contenido, textFile)
	 
	respuestaJSON, err := json.Marshal(contenido)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar datos a JSON: %s", err), http.StatusInternalServerError)
		return
	}
	w.Write(respuestaJSON)
}

func getBack(w http.ResponseWriter, r *http.Request) {
	 
	// Configurar la cabecera de respuesta
	w.Header().Set("Content-Type", "application/json")

	//Abrimos el disco
	disco, err := Herramientas.OpenFile(Structs.UsuarioActual.PathD)
	if err != nil {
		return 
	}

	var superBloque Structs.Superblock
	Herramientas.ReadObject(disco, &superBloque, initSuperBloque)

	//obtengo el ultimo elemento de la lista en el idActual global
	idActual = listaAnterior[len(listaAnterior)-1]
	//elimino el elemento de la lista
	listaAnterior = listaAnterior[:len(listaAnterior)-1]
	//cargo el inodo actual
	var Inode Structs.Inode
	Herramientas.ReadObject(disco, &Inode, int64(superBloque.S_inode_start+(idActual*int32(binary.Size(Structs.Inode{})))))

	//lista de discos encontrados
	var contenido []string

	var folderBlock Structs.Folderblock
	for i := 0; i < 12; i++ {
		idBloque := Inode.I_block[i]
		if idBloque != -1 {
			Herramientas.ReadObject(disco, &folderBlock, int64(superBloque.S_block_start+(idBloque*int32(binary.Size(Structs.Folderblock{})))))
			//Recorrer el bloque actual buscando la carpeta/archivo en la raiz
			for j := 2; j < 4; j++ {
				//apuntador es el apuntador del bloque al inodo (carpeta/archivo), si existe es distinto a -1
				apuntador := folderBlock.B_content[j].B_inodo
				if apuntador != -1 {
					pathActual := Structs.GetB_name(string(folderBlock.B_content[j].B_name[:]))					
					contenido = append(contenido, pathActual)
				}
			}
		}
	}

	respuestaJSON, err := json.Marshal(contenido)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al serializar datos a JSON: %s", err), http.StatusInternalServerError)
		return
	}
	w.Write(respuestaJSON)
}