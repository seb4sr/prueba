package Structs

import (
	"MIA_2S_P2_201513656/Herramientas"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

 
type MBR struct {
	MbrSize    int32         
	FechaC     [16]byte      
	Id         int32         
	Fit        [1]byte       
	Partitions [4]Partition  
}

 
func PrintMBR(data MBR) {
	 
	 
	for i := 0; i < 4; i++ {
		 
	}
}

func GetIdMBR(m MBR) int32{
	return m.Id
}

 

type Partition struct {
	Status      [1]byte   
	Type        [1]byte   
	Fit         [1]byte   
	Start       int32     
	Size        int32     
	Name        [16]byte  
	Correlative int32     
	Id          [4]byte   
}

func (p *Partition) GetEnd() int32 {
	return p.Start + p.Size
}

 
func GetName(nombre string) string {
	posicionNulo := strings.IndexByte(nombre, 0)
	 
	if posicionNulo != -1 {
		 
		nombre = nombre[:posicionNulo]
	}
	return nombre
}

 
func (p *Partition) SetInfo(newType string, fit string, newStart int32, newSize int32, name string, correlativo int32) {
	p.Size = newSize
	p.Start = newStart
	p.Correlative = 0
	copy(p.Name[:], name)
	copy(p.Fit[:], fit)
	copy(p.Status[:], "I")
	copy(p.Type[:], newType)
}

func GetId(nombre string) string {
	 
	posicionNulo := strings.IndexByte(nombre, 0)
	 
	if posicionNulo != -1 {
		nombre = "-"
	}
	return nombre
}


 
type EBR struct {
	Status [1]byte  
	Type   [1]byte
	Fit    [1]byte   
	Start  int32     
	Size   int32     
	Name   [16]byte  
	Next   int32     
}


func (e *EBR) SetInfo(fit string, newStart int32, newSize int32, name string, newNext int32) {
	e.Size = newSize
	e.Start = newStart
	e.Next = newNext
	copy(e.Name[:], name)
	copy(e.Fit[:], fit)
	copy(e.Status[:], "I")
	copy(e.Type[:], "L")
}

func (e *EBR) GetEnd() int32 {
	return e.Start + e.Size + int32(binary.Size(e))
}

/*func GetIdMount (data Mount) string{
	return data.MPath
}*/

/*===========================================================================================
 ======================================= REPORTE MBR =============================================
 ===============================================================================================*/

func RepGraphviz(data MBR, disco *os.File) string {
	disponible := int32(0)
	cad := ""
	inicioLibre := int32(binary.Size(data))  
	for i := 0; i < 4; i++ {
		if data.Partitions[i].Size > 0 {

			disponible = data.Partitions[i].Start - inicioLibre
			inicioLibre = data.Partitions[i].Start + data.Partitions[i].Size

			 
			if disponible > 0 {
				cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#808080' COLSPAN=\"2\"> ESPACIO LIBRE <br/> %d bytes </td> \n </tr> \n", disponible)
			}
			 
			cad += " <tr>\n  <td bgcolor='DeepSkyBlue' COLSPAN=\"2\"> PARTICION </td> \n </tr> \n"
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_status </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(data.Partitions[i].Status[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_type </td> \n  <td bgcolor='LightSkyBlue'> %s </td> \n </tr> \n", string(data.Partitions[i].Type[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_fit </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(data.Partitions[i].Fit[:]))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_start </td> \n  <td bgcolor='LightSkyBlue'> %d </td> \n </tr> \n", data.Partitions[i].Start)
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", data.Partitions[i].Size)
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='LightSkyBlue'> part_name </td> \n  <td bgcolor='LightSkyBlue'> %s </td> \n </tr> \n", GetName(string(data.Partitions[i].Name[:])))
			cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_id </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", GetId(string(data.Partitions[i].Id[:])))
			if string(data.Partitions[i].Type[:]) == "E" {
				cad += repLogicas(data.Partitions[i], disco)
			}
		}
	}

	 
	disponible = data.MbrSize - inicioLibre
	if disponible > 0 {
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#808080' COLSPAN=\"2\"> ESPACIO LIBRE <br/> %d bytes </td> \n </tr> \n", disponible)
	}

	return cad
}

func repLogicas(particion Partition, disco *os.File) string {
	cad := ""

	var actual EBR
	if err := Herramientas.ReadObject(disco, &actual, int64(particion.Start)); err != nil {
		 
		return ""
	}

	 
	if actual.Size != 0 {
		cad += " <tr>\n  <td bgcolor='SteelBlue' COLSPAN=\"2\"> PARTICION LOGICA </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_status </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Status[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_next </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Next)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_fit </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Fit[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_start </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Start)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", actual.Size)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_name </td> \n  <td bgcolor='SkyBlue'> %s </td> \n </tr> \n", GetName(string(actual.Name[:])))
	}

	 
	for actual.Next != -1 {
		if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
			 
			return ""
		}
		cad += " <tr>\n  <td bgcolor='SteelBlue' COLSPAN=\"2\"> PARTICION LOGICA </td> \n </tr> \n"
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_status </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Status[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_next </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Next)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_fit </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(actual.Fit[:]))
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_start </td> \n  <td bgcolor='SkyBlue'> %d </td> \n </tr> \n", actual.Start)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> part_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", actual.Size)
		cad += fmt.Sprintf(" <tr>\n  <td bgcolor='SkyBlue'> part_name </td> \n  <td bgcolor='SkyBlue'> %s </td> \n </tr> \n", GetName(string(actual.Name[:])))
	}
	return cad
}

/* ===============================================================================================
 ======================================= REPORTE DISK ============================================
 ===============================================================================================*/

 func RepDiskGraphviz(data MBR, disco *os.File) string {
	disponible := int32(0)
	cad := ""
	cadLogicas := ""
	cant := 0
	inicioLibre := int32(binary.Size(data))  
	for i := 0; i < 4; i++ {
		if data.Partitions[i].Size > 0 {
			disponible = data.Partitions[i].Start - inicioLibre
			inicioLibre = data.Partitions[i].Start + data.Partitions[i].Size
			 
			if disponible > 0 {
				porcentaje := float64(disponible) * 100 / float64(data.MbrSize)
				cad += fmt.Sprintf(" <td bgcolor='#808080'  ROWSPAN='3'> ESPACIO LIBRE <br/> %.2f %% </td> \n ", porcentaje)
			}
			porcentaje := float64(data.Partitions[i].Size) * 100 / float64(data.MbrSize)
			if string(data.Partitions[i].Type[:]) == "P" {
				cad += fmt.Sprintf(" <td bgcolor='LightSkyBlue' ROWSPAN='3'> PRIMARIA <br/> %.2f %% </td>\n", porcentaje)
			} else {
				cant, cadLogicas = repLogicasDisk(data.MbrSize, data.Partitions[i], disco)
				cad += fmt.Sprintf(" <td bgcolor='SteelBlue' COLSPAN='%d'> EXTENDIDA </td>\n", cant)
			}
		}
	}

	 
	disponible = data.MbrSize - inicioLibre
	if disponible > 0 {
		porcentaje := float64(disponible) * 100 / float64(data.MbrSize)
		cad += fmt.Sprintf(" <td bgcolor='#808080'  ROWSPAN='3'> ESPACIO LIBRE <br/> %.2f %% </td> \n", porcentaje)
	}
	cad += "</tr>"     
	cad += cadLogicas  
	return cad
}

func repLogicasDisk(MbrSize int32, particion Partition, disco *os.File) (int, string) {
	cant := 0
	cad := "\n\n<tr> \n"
	porcentaje := 0.0

	var actual EBR
	sizeEBR := int32(binary.Size(actual))

	 
	if err := Herramientas.ReadObject(disco, &actual, int64(particion.Start)); err != nil {
		 
		porcentaje = float64(particion.Size) * 100 / float64(MbrSize)
		return 1, fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
	}

	 
	if actual.Size != 0 {
		 
		porcentaje = float64(actual.Size+sizeEBR) * 100 / float64(MbrSize)
		cad += " <td bgcolor='royalblue3' ROWSPAN='2'> EBR </td>\n"
		cad += fmt.Sprintf(" <td bgcolor='darkgoldenrod2' ROWSPAN='2'> LOGICA <br/> %.2f %% </td>\n", porcentaje)
		cant += 2

		 
		if actual.Next != -1 {
			disponible := actual.Next - actual.GetEnd()
			if disponible > 0 {
				porcentaje = float64(disponible) * 100 / float64(MbrSize)
				cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
				cant++
			}
		} else {
			disponible := particion.GetEnd() - actual.GetEnd()
			if disponible > 0 {
				porcentaje = float64(disponible) * 100 / float64(MbrSize)
				cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
				cant++
			}
		}
	} else {
		 
		if actual.Next == -1 {
			 
			porcentaje = float64(particion.Size) * 100 / float64(MbrSize)
			cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
			cant++
		} else {
			 
			porcentaje = float64(actual.Next-particion.Start) * 100 / float64(MbrSize)
			cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
			cant++
		}
	}

	 
	for actual.Next != -1 {
		 
		if err := Herramientas.ReadObject(disco, &actual, int64(actual.Next)); err != nil {
			 
			porcentaje = float64(particion.Size) * 100 / float64(MbrSize)
			return 1, fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
		}

		 
		porcentaje = float64(actual.Size+sizeEBR) * 100 / float64(MbrSize)
		cad += " <td bgcolor='royalblue3' ROWSPAN='2'> EBR </td>\n"
		cad += fmt.Sprintf(" <td bgcolor='darkgoldenrod2' ROWSPAN='2'> LOGICA <br/> %.2f %% </td>\n", porcentaje)
		cant += 2

		 
		if actual.Next != -1 {
			disponible := actual.Next - actual.GetEnd()
			if disponible > 0 {
				porcentaje = float64(disponible) * 100 / float64(MbrSize)
				cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
				cant++
			}
		} else {
			 
			disponible := particion.GetEnd() - actual.GetEnd()
			if disponible > 0 {
				porcentaje = float64(disponible) * 100 / float64(MbrSize)
				cad += fmt.Sprintf(" <td bgcolor='#808080' ROWSPAN='2'> LIBRE <br/> %.2f %% </td>\n", porcentaje)
				cant++
			}
		}
	}

	cad += "</tr>\n"
	return cant, cad
}
