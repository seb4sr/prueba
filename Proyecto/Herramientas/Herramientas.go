package Herramientas

import (
	"encoding/binary"
	 
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"
)

func CrearDisco(path string) error {
	 
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		 
		return err
	}

	 
	if _, err := os.Stat(path); os.IsNotExist(err) {
		newFile, err := os.Create(path)
		if err != nil {
			 
			return err
		}
		defer newFile.Close()
	}
	return nil
}

func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		 
		return nil, err
	}
	return file, nil
}

 
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)  
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		 
		return err
	}
	return nil
}

 
func DeletePart(file *os.File, position int64, size int32)  error{
	zeros := make([]byte, size)
	file.Seek(position, 0)  
	err := binary.Write(file, binary.LittleEndian, zeros)
	if err != nil {
		 
		return err
	}
	return nil	
}

 
func DelPartL(size int32) []byte {
	datos := make([]byte, size)
	return datos
}

 
func EscribirPartL(file *os.File, data string, position int64) error {
	 
    dataBytes := []byte(data)
	file.Seek(position, 0)  
	err := binary.Write(file, binary.LittleEndian, dataBytes)	
	if err != nil {
		 
		return err
	}
	return nil
}

 
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		 
		return err
	}
	return nil
}


 
func EliminartIlegibles(entrada string) string{
	 
	transformFunc := func(r rune) rune {
		 
		 
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}

	 
	salida := strings.Map(transformFunc, entrada)
	return salida	
}


func Reporte(path string, contenido string) error {
	 
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		 
		return err
	}
	 
	file, err := os.Create(path)
	if err != nil {
		 
		return err
	}
	defer file.Close()

	 
	_, err = file.WriteString(contenido)
	if err != nil {
		 
		return err
	}

	return err
}


func RepGraphizMBR(path string, contenido string, nombre string) error {
	 
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		 
		return err
	}
	 
	file, err := os.Create(path)
	if err != nil {
		 
		return err
	}
	defer file.Close()

	 
	_, err = file.WriteString(contenido)
	if err != nil {
		 
		return err
	}

	rep2 := dir + "/" + nombre + ".png"
	cmd := exec.Command("dot", "-Tpng", path, "-o", rep2)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error al generar el reporte PNG: %v", err)
	}

	return err
}

