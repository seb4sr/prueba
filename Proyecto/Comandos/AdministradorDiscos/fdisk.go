package administradordiscos

import (
	"MIA_2S_P2_201513656/Herramientas"
	"MIA_2S_P2_201513656/Structs"
	"encoding/binary"
	//"fmt"
	"os"
	"strconv"
	"strings"
)

func Fdisk(entrada []string) string{
	var respuesta string
	 
	 
	unit:=1024 	 
	tipe:="P"	 
	fit :="W"	 
	var size int			 
	var pathE string		 
	var name string			 
	
	var add int            
	var delete int		   
	var opcion int         
	
	Valido := true         
	var sizeValErr string  

	for _,parametro :=range entrada[1:]{
		tmp := strings.TrimRight(parametro," ")
		valores := strings.Split(tmp,"=")

		if len(valores)!=2{
			 
			Valido = false
			 
			return "ERROR FDISK, valor desconocido de parametros "+valores[1]
		}

		 
		if strings.ToLower(valores[0])=="size"{
			var err error
			 
			 
			size, err = strconv.Atoi(valores[1])  
			if err != nil {
				sizeValErr = valores[1]  
			}

		 
		} else if strings.ToLower(valores[0]) == "unit" {
			 
			if strings.ToLower(valores[1]) == "b" {
				unit = 1
				 
			} else if strings.ToLower(valores[1]) == "m" {
				unit = 1048576  
			} else if strings.ToLower(valores[1]) != "k" {
				Valido = false
				 
				return "ERROR FDISK en -unit. Valores aceptados: b, k, m. ingreso: "+ valores[1]				
			}

		 
		} else if strings.ToLower(valores[0]) == "path" {
			pathE = strings.ReplaceAll(valores[1],"\"","")

			_, err := os.Stat(pathE)
			if os.IsNotExist(err) {
				 
				return "ERROR FDISK: El disco no existe" 
			}
		
		 
		} else if strings.ToLower(valores[0]) == "type" {
			 
			if strings.ToLower(valores[1]) == "e" {
				tipe = "E"
			} else if strings.ToLower(valores[1]) == "l" {
				tipe = "L"
			} else if strings.ToLower(valores[1]) != "p" {
				 
				return "ERROR FDISK en -type. Valores aceptados: e, l, p. ingreso: "+ valores[1]
			}

		 
		}else if strings.ToLower(valores[0])=="fit"{
			if strings.ToLower(strings.TrimSpace(valores[1]))=="bf"{
				fit = "B"
			}else if strings.ToLower(valores[1])=="ff"{
				fit = "F"
			}else if strings.ToLower(valores[1])!="wf"{
				 
				return "EEROR: PARAMETRO FIT INCORRECTO. VALORES ACEPTADO: FF, BF,WF. SE INGRESO:"+valores[1]
			}
			
			
		 
		} else if strings.ToLower(valores[0]) == "name" {
			 
			name = strings.ReplaceAll(valores[1], "\"", "")
			 
			name = strings.TrimSpace(name)		
		
		 
		} else if strings.ToLower(valores[0]) == "delete" {	
			if strings.ToLower(valores[1]) == "full" {
				if opcion == 0 {
					opcion = 2  
					delete = 1
				}
			} else if strings.ToLower(valores[1]) == "fast" {
				if opcion == 0 {
					delete = 2
					opcion = 2  
				}
			} else {
				 
				Valido = false
				return "ERROR FDISK. Valor de delete desconocido"
			}
		 
		} else if strings.ToLower(valores[0]) == "add" {
			var err error
			add, err = strconv.Atoi(valores[1])  
			if err != nil {
				 
				Valido = false
				return "ERROR FDISK: El valor de \"add\" debe ser un valor numerico. se leyo "+ valores[1]
			} else {
				if opcion == 0 {
					opcion = 1
				}
			}
		 
		} else {
			 
			return "ERROR FDISK: Parametro desconocido: "+ valores[0]  
		}
	}

	 
	if size != 0{
		if sizeValErr == "" {  
			if size <= 0 {  
				 
				Valido = false
				return "ERROR FDISK: -size debe ser un valor positivo mayor a cero (0). se leyo " + string(size)
			}
		} else {  
			 
			Valido = false
			return "ERROR FDISK: -size debe ser un valor numerico. se leyo "+ sizeValErr
		}
	}else{
		 
		Valido =false
		return "ERROR FDISK: FALTO PARAMETRO SIZE"
	}

	if pathE == ""{
		 
		Valido = false
		return "ERROR FDISK: FALTA PARAMETRO PATH"
	}
	if name == ""{
		 
		Valido = false
		return "ERROR FDISK: FALTA PARAMETRO NAME"
	}

	if Valido{
		 
		disco, err := Herramientas.OpenFile(pathE)
		if err != nil {
			 
			return "ERROR FDISK: No se pudo leer el disco"+ "\n"
		}

		 
		var mbr Structs.MBR
		 
		if err := Herramientas.ReadObject(disco, &mbr, 0); err != nil {
			return "ERROR FDISK Read " + err.Error()+ "\n"
		}

		 
		if opcion == 0{
			 
			isPartExtend := false  
			isName := true         
			if tipe == "E" {
				for i := 0; i < 4; i++ {
					tipo := string(mbr.Partitions[i].Type[:])
					
					if tipo != "E" {
						isPartExtend = true
					} else {
						isPartExtend = false
						isName = false  
						 
						 
						return "ERROR FDISK. Ya existe una particion extendida \nFDISK Error. No se puede crear la nueva particion con nombre:  " + name+ "\n"
					}
				}
			}

			 
			if isName {
				for i := 0; i < 4; i++ {
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name {
						isName = false
						 
						 
						return "ERROR FDISK. Ya existe la particion : " + name + "\nFDISK Error. No se puede crear la nueva particion con nombre: " + name+ "\n"

					}
				}
			}

			if isName{
				 
				var partExtendida Structs.Partition
				 
				if string(mbr.Partitions[0].Type[:]) == "E" {
					partExtendida = mbr.Partitions[0]
				} else if string(mbr.Partitions[1].Type[:]) == "E" {
					partExtendida = mbr.Partitions[1]
				} else if string(mbr.Partitions[2].Type[:]) == "E" {
					partExtendida = mbr.Partitions[2]
				} else if string(mbr.Partitions[3].Type[:]) == "E" {
					partExtendida = mbr.Partitions[3]
				}

				if partExtendida.Size != 0{
					var actual Structs.EBR
					if err := Herramientas.ReadObject(disco, &actual, int64(partExtendida.Start)); err != nil {
						return "ERROR FDISK Read " + err.Error()+ "\n"
					}

					 
					if Structs.GetName(string(actual.Name[:])) == name {
						isName = false
					} else{
						for actual.Next != -1 {
							 
							if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
								return "ERROR FDISK Read " + err.Error()+ "\n"
							}
							if Structs.GetName(string(actual.Name[:])) == name {
								isName = false
								break
							}
						}
					}

					if !isName {
						 
						 
						respuesta += "ERROR FDISK. Ya existe la particion : " + name
						respuesta += "\nFDISK Error. No se puede crear la nueva particion con nombre: " + name+ "\n"
						return respuesta
						
					}
				}
			}

			 
			sizeNewPart := size * unit  
			guardar := false            
			var newPart Structs.Partition
			if (tipe == "P" || isPartExtend) && isName{ 
				sizeMBR := int32(binary.Size(mbr))  
				 
				 

				 
				var resTem string
				mbr, newPart, resTem = primerAjuste(mbr, tipe, sizeMBR, int32(sizeNewPart), name, fit)  
				respuesta += resTem
				guardar = newPart.Size != 0

				 
				if guardar{
					 
					if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {
						return "ERROR FDISK Write " +err.Error()+ "\n"
					}

					 
					if isPartExtend {
						var ebr Structs.EBR
						ebr.Start = newPart.Start
						ebr.Next = -1
						if err := Herramientas.WriteObject(disco, ebr, int64(ebr.Start)); err != nil {
							return "ERROR FDISK Write " +err.Error()+ "\n"
						}
					}
					 
					var TempMBR2 Structs.MBR
					 
					if err := Herramientas.ReadObject(disco, &TempMBR2, 0); err != nil {
						return "ERROR FDISK Read " + err.Error()+ "\n"
					}
					Structs.PrintMBR(TempMBR2)
					 
					respuesta += "\nParticion con nombre " + name + " creada exitosamente"+ "\n"					
				}else {
					 
					 
					return "ERROR FDISK. No se puede crear la nueva particion con nombre: "+ name
				}
			}else if tipe == "L" && isName{
				var partExtend Structs.Partition
				if string(mbr.Partitions[0].Type[:]) == "E" {
					partExtend = mbr.Partitions[0]
				} else if string(mbr.Partitions[1].Type[:]) == "E" {
					partExtend = mbr.Partitions[1]
				} else if string(mbr.Partitions[2].Type[:]) == "E" {
					partExtend = mbr.Partitions[2]
				} else if string(mbr.Partitions[3].Type[:]) == "E" {
					partExtend = mbr.Partitions[3]
				} else {
					 
					return"ERROR FDISK. No existe una particion extendida en la cual crear un particion logica"+ "\n"
				}

				 
				if partExtend.Size != 0 {
					 
					respuesta += primerAjusteLogicas(disco, partExtend, int32(sizeNewPart), name, fit) + "\n" 
					 
				}
				return respuesta
			}
		 
		 
		 
		}else if opcion == 1 {
			add = add * unit
			 
			 
			if add < 0 {
				 
				reducir := true  
				for i := 0; i < 4; i++{
					nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
					if nombre == name{
						reducir = false
						newSize := mbr.Partitions[i].Size + int32(add)
						if newSize > 0{
							mbr.Partitions[i].Size += int32(add)
							if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {  
								return "ERROR FDISK: " + err.Error()
							}
							 
							return "Particion con nombre "+ name+" se redujo correctamente"
						}else {
							 
							return "ERROR FDISK. El tamaño que intenta eliminar es demasiado grande"
						}
					}
				}

				 
				if reducir{
					var partExtendida Structs.Partition
					 
					if string(mbr.Partitions[0].Type[:]) == "E" {
						partExtendida = mbr.Partitions[0]
					} else if string(mbr.Partitions[1].Type[:]) == "E" {
						partExtendida = mbr.Partitions[1]
					} else if string(mbr.Partitions[2].Type[:]) == "E" {
						partExtendida = mbr.Partitions[2]
					} else if string(mbr.Partitions[3].Type[:]) == "E" {
						partExtendida = mbr.Partitions[3]
					}

					 
					if partExtendida.Size != 0{
						var actual Structs.EBR
						if err := Herramientas.ReadObject(disco, &actual, int64(partExtendida.Start)); err != nil {
							return "ERROR FDISK, READ "+ err.Error()
						}

						 
						if Structs.GetName(string(actual.Name[:])) == name {
							reducir = false
						} else {
							for actual.Next != -1 {
								 
								if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
									return "ERROR FDISK, READ "+ err.Error()
								}
								if Structs.GetName(string(actual.Name[:])) == name {
									reducir = false
									break
								}
							}
						}

						if !reducir {
							actual.Size += int32(add)
							if actual.Size > 0 {
								if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {  
									return "ERROR FDISK, write "+ err.Error()
								}
								 
								return "Particion con nombre "+ name+ " se redujo correctamente"
							} else {
								 
								return "ERROR FDISK. El tamaño que intenta eliminar es demasiado grande"
							}
						}
					}
				}

				if reducir {
					 
					return "ERROR FDISK. No existe la particion a reducir"
				}
			 
			}else if add > 0{
				 
				 
				evaluar := 0
				 
				if Structs.GetName(string(mbr.Partitions[0].Name[:])) == name{
					if mbr.Partitions[1].Start == 0 {
						if mbr.Partitions[2].Start == 0 {
							if mbr.Partitions[3].Start == 0 {
								evaluar = int(mbr.MbrSize - mbr.Partitions[0].GetEnd())
							} else {
								evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[0].GetEnd())
							}
						} else {
							evaluar = int(mbr.Partitions[2].Start - mbr.Partitions[0].GetEnd())
						}
					} else {
						evaluar = int(mbr.Partitions[1].Start - mbr.Partitions[0].GetEnd())
					}

					 
					if evaluar > 0 && add <= evaluar {
						 
						mbr.Partitions[0].Size += int32(add)
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {  
							return "ERROR FDISK, write "+ err.Error()
						}
						 
						return "Particion con nombre "+ name+ " aumento el espacio correctamente"
					} else {
						 
						return "ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion " + name
					}
				 
				}else if Structs.GetName(string(mbr.Partitions[1].Name[:])) == name{
					if mbr.Partitions[2].Start == 0 {
						if mbr.Partitions[3].Start == 0 {
							evaluar = int(mbr.MbrSize - mbr.Partitions[1].GetEnd())
						} else {
							evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[1].GetEnd())
						}
					} else {
						evaluar = int(mbr.Partitions[2].Start - mbr.Partitions[1].GetEnd())
					}
					 
					if evaluar > 0 && add <= evaluar {
						mbr.Partitions[1].Size += int32(add)
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {  
							return "ERROR FDISK WRITE "+err.Error()
						}
						 
						return "Particion con nombre "+ name+ " aumento el espacio correctamente"
					} else {
						 
						return "ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion " + name
					}
				 
				}else if Structs.GetName(string(mbr.Partitions[2].Name[:])) == name{
					if mbr.Partitions[3].Start == 0 {
						evaluar = int(mbr.MbrSize - mbr.Partitions[2].GetEnd())
					} else {
						evaluar = int(mbr.Partitions[3].Start - mbr.Partitions[2].GetEnd())
					}
					 
					if evaluar > 0 && add <= evaluar {
						mbr.Partitions[2].Size += int32(add)
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {  
							return "ERROR FDISK WRITE "+err.Error()
						}
						 
						return "Particion con nombre " + name+ " aumento el espacio correctamente"
					} else {
						 
						return "ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion "+name
					}
				 
				}else if Structs.GetName(string(mbr.Partitions[3].Name[:])) == name{
					evaluar = int(mbr.MbrSize - mbr.Partitions[3].GetEnd())
					 
					if evaluar > 0 && add <= evaluar {
						mbr.Partitions[3].Size += int32(add)
						if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {  
							return "ERROR FDISK, WRITE "+err.Error()
						}
						 
						return "Particion con nombre "+ name+ " aumento el espacio correctamente"
					} else {
						 
						return "ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion "+ name
					}		
				 
				}else{
					 
					var partExtendida Structs.Partition
					 
					if string(mbr.Partitions[0].Type[:]) == "E" {
						partExtendida = mbr.Partitions[0]
					} else if string(mbr.Partitions[1].Type[:]) == "E" {
						partExtendida = mbr.Partitions[1]
					} else if string(mbr.Partitions[2].Type[:]) == "E" {
						partExtendida = mbr.Partitions[2]
					} else if string(mbr.Partitions[3].Type[:]) == "E" {
						partExtendida = mbr.Partitions[3]
					}

					 
					if partExtendida.Size != 0{
						aumentar := false
						var actual Structs.EBR
						if err := Herramientas.ReadObject(disco, &actual, int64(partExtendida.Start)); err != nil {
							return "ERROR FDISK " +err.Error()
						}

						 
						if Structs.GetName(string(actual.Name[:])) == name {
							aumentar = true
						} else{
							for actual.Next != -1 {
								 
								if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
									return "ERROR FDISK " +err.Error()
								}
								if Structs.GetName(string(actual.Name[:])) == name {
									aumentar = true
									break
								}
							}
						}

						if aumentar {
							if actual.Next != -1 {
								if add <= int(actual.Next)-int(actual.GetEnd()) {
									actual.Size += int32(add)
									if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {  
										return "ERROR FDISK "+err.Error()
									}
									 
								} else {
									 
								}
							}else{
								if add <= int(partExtendida.GetEnd())-int(actual.GetEnd()) {
									actual.Size += int32(add)
									if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {  
										return "ERROR FDISK "+err.Error()
									}
									 
									return "Particion con nombre "+ name+ " aumento el espacio correctamente"
								} else {
									 
									return "ERROR FDISK. El tamaño que intenta aumentar es demasiado grande para la particion " + name
								}
							}
						}else {
							 
							return "ERROR FDISK. No existe la particion a aumentar"
						}
					}else{
						 
						return "ERROR FDISK. No existe particion extendida"
					}
				}
			} else {
				 
				return "ERROR FDISK. 0 no es un valor valido para aumentar o disminuir particiones"
			}


		 
		 
		 
		}else if opcion == 2 {
			 
			del := true  
			for i := 0; i < 4; i++ {
				nombre := Structs.GetName(string(mbr.Partitions[i].Name[:]))
				if nombre == name {
					if delete == 1{ 
						Herramientas.DeletePart(disco, int64(mbr.Partitions[i].Start), mbr.Partitions[i].Size)
					}
					var newPart Structs.Partition
					mbr.Partitions[i] = newPart
					if err := Herramientas.WriteObject(disco, mbr, 0); err != nil {  
						return "ERROR FDISK WRITE " + err.Error()
					}
					del = false
					 
					return "particion con nombre "+ name+ " eliminada"
				}
			}

			 
			if del{
				var partExtendida Structs.Partition
				 
				if string(mbr.Partitions[0].Type[:]) == "E" {
					partExtendida = mbr.Partitions[0]
				} else if string(mbr.Partitions[1].Type[:]) == "E" {
					partExtendida = mbr.Partitions[1]
				} else if string(mbr.Partitions[2].Type[:]) == "E" {
					partExtendida = mbr.Partitions[2]
				} else if string(mbr.Partitions[3].Type[:]) == "E" {
					partExtendida = mbr.Partitions[3]
				}
				
				 
				if partExtendida.Size != 0 {
					var actual Structs.EBR
					if err := Herramientas.ReadObject(disco, &actual, int64(partExtendida.Start)); err != nil {
						return "ERROR FDISK READ EBR "+err.Error()
					}
					var anterior Structs.EBR
					var eliminar Structs.EBR  

					 
					if Structs.GetName(string(actual.Name[:])) == name {
						del = false
					} else {
						for actual.Next != -1 {
							 
							if err := Herramientas.ReadObject(disco, &anterior, int64(actual.Start)); err != nil {
								return "ERROR FDISK READ EBR "+err.Error()
							}
							 
							if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
								return "ERROR FDISK READ EBR "+err.Error()
							}
							 
							if Structs.GetName(string(actual.Name[:])) == name {
								del = false
								break
							}
						}
					}

					 
					if !del {
						 
						sizeEBR := int32(binary.Size(actual))
						 
						if actual.Next != -1 {
							if anterior.Size == 0 {
								 
								if delete == 1{
									 
									if err := Herramientas.WriteObject(disco, Herramientas.DelPartL(actual.Size), int64(actual.Start+sizeEBR)); err != nil {
										return "ERROR FDISK ESCRITURA DE EBR "+err.Error()
									}
								}
								
								actual.Size = 0              
								actual.Name = eliminar.Name  
								
								if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
									return "ERROR FDISK ESCRITURA DE EBR "+err.Error()
								}
								 
																	
								 
								return "Particion con nombre "+ name+ " eliminada"
							}else{
								 
								if delete == 1 {
									 
									if err := Herramientas.WriteObject(disco, Herramientas.DelPartL(actual.Size+sizeEBR), int64(actual.Start)); err != nil {
										return "ERROR FDISK ESCRITURA DE EBR "+err.Error()
									}
								}
								 
								anterior.Next = actual.Next
								actual.Size = 0
								actual.Name = eliminar.Name 
								if err := Herramientas.WriteObject(disco, anterior, int64(anterior.Start)); err != nil {
									return "ERROR FDISK ESCRITURA DE EBR "+err.Error()
								}
								 
								return "Particion con nombre "+ name+ " eliminada"
							}
						}else{
							 
							if anterior.Size == 0 {
								 
								if delete == 1{
									 
									if err := Herramientas.WriteObject(disco, Herramientas.DelPartL(actual.Size), int64(actual.Start+sizeEBR)); err != nil {
										return "ERROR FDISK ESCRITURA DE EBR "+err.Error()
									}
								}
								actual.Size = 0
								actual.Name = eliminar.Name
								 
								if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
									return "ERROR FDISK ESCRITURA DE EBR "+err.Error()
								}
								 
								 
							}else{
								 
								anterior.Next = -1
								if err := Herramientas.WriteObject(disco, anterior, int64(anterior.Start)); err != nil {
									return "ERROR FDISK ESCRITURA DE EBR "+err.Error()
								}

								if delete == 1{
									 
									if err := Herramientas.WriteObject(disco, Herramientas.DelPartL(actual.Size+sizeEBR), int64(actual.Start)); err != nil {
										return "ERROR FDISK ESCRITURA DE EBR "+err.Error()
									}
								}
								
								 
								return "Particion con nombre "+ name+ " eliminada"
							}
						}
					}
				}else {
					 
					return "ERROR FDISK. No se encontro la particion de nombre "+name
				}
			}

			 
			if del {
				 
				return "ERROR FDISK. No se encontro la particion de nombre "+name
			}

		}else {
			 
			return "ERROR FDISK. Operación desconocida (operaciones aceptadas: crear, modificar o eliminar)"
		}

		 
		defer disco.Close()
	} 

	return respuesta
}

 
func primerAjuste(mbr Structs.MBR, typee string, sizeMBR int32, sizeNewPart int32, name string, fit string) (Structs.MBR, Structs.Partition, string) {
	var respuesta string
	var newPart Structs.Partition
	var noPart Structs.Partition  

	 
	if mbr.Partitions[0].Size == 0 {
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if mbr.Partitions[1].Size == 0 {
			if mbr.Partitions[2].Size == 0 {
				 
				if mbr.Partitions[3].Size == 0 {
					 
					if sizeNewPart <= mbr.MbrSize-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						newPart = noPart
						 
						return mbr, newPart, "ERROR FDISK. Espacio insuficiente \n"
					}
				} else {
					 
					 
					if sizeNewPart <= mbr.Partitions[3].Start-sizeMBR {
						mbr.Partitions[0] = newPart
					} else {
						 
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							 
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							 
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"+ "\n"
						}
					}
				}
				 
			} else {
				 
				 
				if sizeNewPart <= mbr.Partitions[2].Start-sizeMBR {
					mbr.Partitions[0] = newPart
				} else {
					 
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							 
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"+ "\n"
						}
					} else {
						 
						 
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							 
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3  
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							 
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							 
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							 
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"+ "\n"							
						}
					}  
				}  
			}  
		} else {
			 
			 
			if sizeNewPart <= mbr.Partitions[1].Start-sizeMBR {
				mbr.Partitions[0] = newPart
			} else {
				 
				 
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[2].Size == 0 {
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							newPart = noPart
							 
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente \n"
						}
					} else {
						 
						 
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[2] = newPart
						} else {
							 
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							if sizeNewPart <= mbr.MbrSize-newPart.Start {  
								mbr.Partitions[2] = mbr.Partitions[3]
								mbr.Partitions[3] = newPart
								 
								mbr.Partitions[2].Correlative = 3
							} else {
								newPart = noPart
								 
								return mbr, newPart,"ERROR FDISK. Espacio insuficiente \n"								
							}
						}  
					}  
				} else {
					 
					 
					if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
						mbr.Partitions[0] = mbr.Partitions[1]
						mbr.Partitions[1] = newPart
						 
						mbr.Partitions[0].Correlative = 1
						mbr.Partitions[1].Correlative = 2
					} else if mbr.Partitions[3].Size == 0 {
						 
						 
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
						} else {
							newPart = noPart
							 
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
						}
					} else {
						 
						 
						newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							 
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							 
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[0] = mbr.Partitions[1]
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							 
							mbr.Partitions[0].Correlative = 1
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							 
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
						}
					}  
				}  
			}  
		}  
		 

		 
	} else if mbr.Partitions[1].Size == 0 {
		 
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {  
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			 
			mbr.Partitions[1].Correlative = 2
		} else {
			 
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)  
			if mbr.Partitions[2].Size == 0 {
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[1] = newPart
					} else {
						newPart = noPart
						 
						return mbr, newPart,"ERROR FDISK. Espacio insuficiente"
					}
				} else {
					 
					 
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[1] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						 
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						 
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						 
						return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
					}
				}  
			} else {
				 
				 
				if sizeNewPart <= mbr.Partitions[2].Start-newPart.Start {
					mbr.Partitions[1] = newPart
				} else {
					 
					newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 3)
					if mbr.Partitions[3].Size == 0 {
						if sizeNewPart <= mbr.MbrSize-newPart.Start {
							mbr.Partitions[3] = newPart
							 
							mbr.Partitions[3].Correlative = 4
						} else {
							newPart = noPart
							 
							return mbr, newPart,"ERROR FDISK. Espacio insuficiente"
						}
					} else {
						 
						 
						if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = newPart
							 
							mbr.Partitions[1].Correlative = 2
						} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
							 
							newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
							mbr.Partitions[1] = mbr.Partitions[2]
							mbr.Partitions[2] = mbr.Partitions[3]
							mbr.Partitions[3] = newPart
							 
							mbr.Partitions[1].Correlative = 2
							mbr.Partitions[2].Correlative = 3
						} else {
							newPart = noPart
							 
							return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
						}
					}  
				}  
			}  
		}  
		 

		 
	} else if mbr.Partitions[2].Size == 0 {
		 
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			 
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			 
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				 
				mbr.Partitions[2].Correlative = 3
			} else {
				 
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				if mbr.Partitions[3].Size == 0 {
					if sizeNewPart <= mbr.MbrSize-newPart.Start {
						mbr.Partitions[2] = newPart
					} else {
						newPart = noPart
						 
						return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
					}
				} else {
					 
					 
					if sizeNewPart <= mbr.Partitions[3].Start-newPart.Start {
						mbr.Partitions[2] = newPart
					} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[3].GetEnd() {
						 
						newPart.SetInfo(typee, fit, mbr.Partitions[3].GetEnd(), sizeNewPart, name, 4)
						mbr.Partitions[2] = mbr.Partitions[3]
						mbr.Partitions[3] = newPart
						 
						mbr.Partitions[2].Correlative = 3
					} else {
						newPart = noPart
						 
						return mbr, newPart, "ERROR FDISK. Espacio insuficiente"
					}
				}  
			}  
		}  
		 

		 
	} else if mbr.Partitions[3].Size == 0 {
		 
		newPart.SetInfo(typee, fit, sizeMBR, sizeNewPart, name, 1)
		if sizeNewPart <= mbr.Partitions[0].Start-newPart.Start {
			mbr.Partitions[3] = mbr.Partitions[2]
			mbr.Partitions[2] = mbr.Partitions[1]
			mbr.Partitions[1] = mbr.Partitions[0]
			mbr.Partitions[0] = newPart
			 
			mbr.Partitions[3].Correlative = 4
			mbr.Partitions[2].Correlative = 3
			mbr.Partitions[1].Correlative = 2
		} else {
			 
			 
			newPart.SetInfo(typee, fit, mbr.Partitions[0].GetEnd(), sizeNewPart, name, 2)
			if sizeNewPart <= mbr.Partitions[1].Start-newPart.Start {
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = mbr.Partitions[1]
				mbr.Partitions[1] = newPart
				 
				mbr.Partitions[3].Correlative = 4
				mbr.Partitions[2].Correlative = 3
			} else if sizeNewPart <= mbr.Partitions[2].Start-mbr.Partitions[1].GetEnd() {
				 
				newPart.SetInfo(typee, fit, mbr.Partitions[1].GetEnd(), sizeNewPart, name, 3)
				mbr.Partitions[3] = mbr.Partitions[2]
				mbr.Partitions[2] = newPart
				 
				mbr.Partitions[3].Correlative = 4
			} else if sizeNewPart <= mbr.MbrSize-mbr.Partitions[2].GetEnd() {
				 
				newPart.SetInfo(typee, fit, mbr.Partitions[2].GetEnd(), sizeNewPart, name, 4)
				mbr.Partitions[3] = newPart
			} else {
				newPart = noPart
				 
				return mbr, newPart,"ERROR FDISK. Espacio insuficiente"
			}
		}  
		 
	} else {
		newPart = noPart
		 
		return mbr, newPart,"ERROR FDISK. Particiones primarias y/o extendidas ya no disponibles"
	}

	return mbr, newPart, respuesta
}

func primerAjusteLogicas(disco *os.File, partExtend Structs.Partition, sizeNewPart int32, name string, fit string) string{
	var respuesta string
	 
	save := true  
	var actual Structs.EBR
	sizeEBR := int32(binary.Size(actual))  
	 

	 
	if err := Herramientas.ReadObject(disco, &actual, int64(partExtend.Start)); err != nil {
		respuesta += "ERROR FDISK Read " + err.Error()+ "\n"
		return respuesta
	}

	 
	 
	 

	 
	if actual.Size == 0 {
		if actual.Next == -1 {
			 
			if sizeNewPart+sizeEBR <= partExtend.Size {
				actual.SetInfo(fit, partExtend.Start, sizeNewPart, name, -1)
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return"ERROR FDISK Write " +err.Error()+ "\n"
				}
				save = false  
				 
				respuesta += "Particion con nombre "+ name+ " creada correctamente"+ "\n"
			} else {
				 
				return "ERROR FDISK. Espacio insuficiente logicas"+ "\n"
			}
		} else {
			 
			 
			disponible := actual.Next - partExtend.Start  
			if sizeNewPart+sizeEBR <= disponible {
				actual.SetInfo(fit, partExtend.Start, sizeNewPart, name, actual.Next)
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return "ERROR FDISK Write " +err.Error()+ "\n"
				}
				save = false  
				 
				respuesta += "Particion con nombre " + name+ " creada correctamente"+ "\n"
			} else {
				 
				return "ERROR FDISK. Espacio insuficiente logicas"+ "\n"
			}
		}
		 
	}

	if save {
		 
		for actual.Next != -1 {
			 
			if sizeNewPart+sizeEBR <= actual.Next-actual.GetEnd() {
				break
			}
			 
			if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
				respuesta += "ERROR FDISK Read " + err.Error()+ "\n"
				return respuesta
			}

		}

		 
		if actual.Next == -1 {
			 
			if sizeNewPart+sizeEBR <= (partExtend.GetEnd() - actual.GetEnd()) {
				 
				actual.Next = actual.GetEnd()
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return"ERROR FDISK Write " +err.Error()+ "\n"
				}

				 
				newStart := actual.GetEnd()                           
				actual.SetInfo(fit, newStart, sizeNewPart, name, -1)  
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return "ERROR FDISK Write " +err.Error()+ "\n"
				}
				 
				respuesta += "Particion con nombre "+ name+" creada correctamente"+ "\n"
			} else {
				 
				return "ERROR FDISK. Espacio insuficiente logicas"+ "\n"
			}
		} else {
			 
			if sizeNewPart+sizeEBR <= (actual.Next - actual.GetEnd()) {
				siguiente := actual.Next  
				 
				actual.Next = actual.GetEnd()
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return "ERROR FDISK Write " +err.Error()+ "\n"
				}

				 
				newStart := actual.GetEnd()                                  
				actual.SetInfo(fit, newStart, sizeNewPart, name, siguiente)  
				if err := Herramientas.WriteObject(disco, actual, int64(actual.Start)); err != nil {
					return "ERROR FDISK Write " +err.Error()+ "\n"
				}

				 
				respuesta += "Particion con nombre "+ name +" creada correctamente"+ "\n"
			} else {
				 
				return "ERROR FDISK. Espacio insuficiente logicas "+ "\n"
			}
		}
	}
	return respuesta
}