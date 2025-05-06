package Structs

/*Almacena la informacion de los Discos montados:
Se asigna una letra a cada disco montado y 
va sumando 1 cada vez que se monta otra particion en dicho disco
*/
var Pmontaje []DMontado
type DMontado struct{
	MPath  string	 
	Letter byte		 
	Cont   int 		 
}

 
func AddPathM (path string, L byte, cont int){
	Pmontaje = append(Pmontaje, DMontado{MPath: path ,Letter: L,Cont: cont})
}

 

 

var Montadas []mountAlready
type mountAlready struct{
	Id		 string	 
	PathM	 string	 
}

 
func AddMontadas(id string, path string){
	Montadas = append(Montadas, mountAlready{Id: id, PathM: path})
}