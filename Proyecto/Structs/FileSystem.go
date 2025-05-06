package Structs

import (
	"MIA_2S_P2_201513656/Herramientas"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

 
 
type Superblock struct {
	S_filesystem_type   int32     
	S_inodes_count      int32     
	S_blocks_count      int32     
	S_free_blocks_count int32     
	S_free_inodes_count int32     
	S_mtime             [16]byte  
	S_umtime            [16]byte  
	S_mnt_count         int32     
	S_magic             int32     
	S_inode_size        int32     
	S_block_size        int32     
	S_first_ino         int32     
	S_first_blo         int32     
	S_bm_inode_start    int32     
	S_bm_block_start    int32     
	S_inode_start       int32     
	S_block_start       int32     
}

 
type Inode struct {
	I_uid   int32      
	I_gid   int32      
	I_size  int32      
	I_atime [16]byte   
	I_ctime [16]byte   
	I_mtime [16]byte   
	I_block [15]int32  
	I_type  [1]byte    
	I_perm  [3]byte    
}

 
 
type Folderblock struct {
	B_content [4]Content  
}

type Content struct {
	B_name  [12]byte  
	B_inodo int32     
}

 
type Fileblock struct {
	B_content [64]byte  
}

 
type Pointerblock struct {
	B_pointers [16]int32  
}

type Journaling struct {
	Size      int32
	Ultimo    int32
	Contenido [50]Content_J
}

type Content_J struct {
	Operation [10]byte
	Path      [100]byte
	Content   [100]byte
	Date      [16]byte
}

 
func GetB_name(nombre string) string {
	posicionNulo := strings.IndexByte(nombre, 0)

	if posicionNulo != -1 {
		if posicionNulo != 0 {
			 
			nombre = nombre[:posicionNulo]
		} else {
			 
			nombre = "-"
		}

	}
	return nombre  
}

 
func GetB_content(nombre string) string {
	 
	nombre = strings.ReplaceAll(nombre, "\n", "<br/>")
	posicionNulo := strings.IndexByte(nombre, 0)

	if posicionNulo != -1 {
		if posicionNulo != 0 {
			 
			nombre = nombre[:posicionNulo]
		} else {
			 
			nombre = "-"
		}

	}
	 
	 
	return nombre  
}

 
func GetOperation(nombre string) string {
	posicionNulo := strings.IndexByte(nombre, 0)
	nombre = nombre[:posicionNulo]  
	return nombre
}

func GetPath(nombre string) string {
	posicionNulo := strings.IndexByte(nombre, 0)
	nombre = nombre[:posicionNulo]  
	return nombre
}

func GetContent(nombre string) string {
	posicionNulo := strings.IndexByte(nombre, 0)
	nombre = nombre[:posicionNulo]  
	return nombre
}


 
type Bite struct {
	Val [1]byte
}

 
func RepSB(particion Partition, disco *os.File) string {
	cad := ""
	 
	var SuperBloque Superblock
	 
	err := Herramientas.ReadObject(disco, &SuperBloque, int64(particion.Start))
	if err != nil {
		 
		return cad
	}
	 
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_filesystem_type </td> \n  <td bgcolor='Azure'> EXT%d </td> \n </tr> \n", SuperBloque.S_filesystem_type)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_inodes_count </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_inodes_count)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_blocks_count </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_blocks_count)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_free_inodes_count </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_free_inodes_count)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_free_blocks_count </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_free_blocks_count)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_mtime </td> \n  <td bgcolor='#7FC97F'> %s </td> \n </tr> \n", string(SuperBloque.S_mtime[:]))
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_umtime </td> \n  <td bgcolor='Azure'> %s </td> \n </tr> \n", string(SuperBloque.S_mtime[:]))
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_mnt_count </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_mnt_count)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_magic </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_magic)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_inode_size </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_inode_size)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_block_size </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_block_size)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_first_ino </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_first_ino)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_first_blo </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_first_blo)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_bm_inode_start </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_bm_inode_start)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_bm_block_start </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_bm_block_start)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='#7FC97F'> S_inode_start </td> \n  <td bgcolor='#7FC97F'> %d </td> \n </tr> \n", SuperBloque.S_inode_start)
	cad += fmt.Sprintf(" <tr>\n  <td bgcolor='Azure'> S_block_start </td> \n  <td bgcolor='Azure'> %d </td> \n </tr> \n", SuperBloque.S_block_start)
	
	return cad
}

 
func RepJournal(particion Partition, disco *os.File) string {
	cad := ""
	 
	var superBloque Superblock
	err := Herramientas.ReadObject(disco, &superBloque, int64(particion.Start))
	if err != nil {
		 
		cad += " <tr>\n  <td> Error No Journaling </td> \n </tr> \n"
		return cad
	}

	if superBloque.S_filesystem_type == 3 {
		 
		var journal Journaling
		 
		Herramientas.ReadObject(disco, &journal, int64(particion.Start+int32(binary.Size(Superblock{}))))
		for i := int32(0); i <= journal.Ultimo; i++ {
			dataJ := journal.Contenido[i]
			cad += " <tr>\n  <td> Operacion </td> \n  <td> Path </td> \n  <td> Contenido </td> \n  <td> Fecha </td> \n </tr> \n"
			cad += fmt.Sprintf(" <tr>\n  <td> %s </td> \n  <td> %s </td> \n  <td> %s </td> \n  <td> %s </td> \n </tr> \n", GetOperation(string(dataJ.Operation[:])), GetPath(string(dataJ.Path[:])), GetContent(string(dataJ.Content[:])), string(dataJ.Date[:]))
		}
	} else {
		cad += " <tr>\n  <td> Error No Journaling </td> \n </tr> \n"
		 
	}
	return cad
}