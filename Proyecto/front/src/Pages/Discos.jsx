import { useState } from "react";
import { useParams, useNavigate } from "react-router-dom"

import diskIMG from '../iconos/disk.png';
import "../Stylesheets/Fondo.css"

export default function Discos({newIp="localhost"}){
    const { id } = useParams()//Id que viene del login
    const [discos, setDiscos] = useState([]);
    const navigate = useNavigate()
    const [nameDisk, setNameDisk] = useState('');

   
    useState(()=>{
        fetch(`http://${newIp}:8080/discos`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json'},
            body: JSON.stringify(id)
        })
        .then(Response => Response.json())
        .then(rawData => {
            console.log(rawData); 
            setDiscos(rawData.Discos); 
            setNameDisk(rawData.DiskPart);
        })
        .catch(error => {
            console.error('Error en la solicitud Fetch:', error);
            // Maneja el error aquí, como mostrar un mensaje al usuario
            //alert('Error No tiene permiso para acceder a este Disco');
        });
    }, [])

    const onClick = (disco) => {        
        if (disco === nameDisk){
            navigate(`/Particiones/${id}`) //navegar al objeto que hice click
        }else{
            alert('Error No hay sesion iniciada en este Disco, intente con otro');
        }       
    }

    return(
        <div className='body'>
            <div>&nbsp;&nbsp;&nbsp;</div>
            <div style={{display:"flex", flexDirection:"row", justifyContent: "center"}}><h1>DISCOS</h1></div>
            <div style={{display:"flex", flexDirection:"row", justifyContent: "center"}}> 
                {discos && discos.length > 0 ? (
                    discos.map((disco, index) => {
                        return (
                            
                            <div key={index} style={{
                                display: "flex",
                                flexDirection: "column", // Alinea los elementos en columnas
                                alignItems: "center", // Centra verticalmente los elementos
                                maxWidth: "100px",
                                margin: "10px"
                                }}
                                onClick={() => onClick(disco)}
                            >
                                <img src={diskIMG} alt="disk" style={{width: "100px"}} />  
                                <div class="circle">                      
                                    {disco}
                                </div>
                                
                            </div>
                        )
                    })
                ):(
                    <div>No hay discos disponibles</div>
                )}
            </div> 
        </div>
    );
}