import { useState } from "react";
import { useParams } from 'react-router-dom';
import { useNavigate } from "react-router-dom"
import partIMG from '../iconos/part.png';
import "../Stylesheets/Fondo.css"

export default function Partitions({newIp="localhost"}){
    const { id } = useParams()//Id que viene de Discos
    const [ particiones, setParticiones ] = useState([]);
    const [ namePart, setNamePart] = useState('');
    const navigate = useNavigate()
    
    useState(()=>{
        fetch(`http://${newIp}:8080/particiones`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json'},
            body: JSON.stringify(id)
        })
        .then(Response => Response.json())
        .then(rawData => {
            console.log(rawData); 
            setParticiones(rawData.Partciones);
            setNamePart(rawData.IdParticion);
        })
        .catch(error => {
            console.error('Error en la solicitud Fetch:', error);
            // Maneja el error aquí, como mostrar un mensaje al usuario
            //alert('Error en la solicitud Fetch. Por favor, inténtalo de nuevo más tarde.');
        });
    }, [])

    const onClick = (particion) => {
        console.log("click",particion)
        console.log("NamePar",namePart)
        console.log("id",id)
        if (particion === namePart){
            navigate(`/Explorador/${id}`) //navegar al objeto que hice click
        }else{
            alert('Error No hay sesion iniciada en esta Particion, intente con otra');
        }  
        
    }

    return(
        <div className='body'>
            <div>&nbsp;&nbsp;&nbsp;</div>
            <div style={{display:"flex", flexDirection:"row", justifyContent: "center"}}><h1>PARTICIONES DEL DISCO </h1></div>
            <div style={{display:"flex", flexDirection:"row", justifyContent: "center"}}>
                {particiones && particiones.length > 0 ? (
                    particiones.map((particion, index) => {
                        return (
                            <div key={index} style={{
                                display: "flex",
                                flexDirection: "column", // Alinea los elementos en columnas
                                alignItems: "center", // Centra verticalmente los elementos
                                maxWidth: "100px",
                                margin: "10px"
                                }}
                                onClick={() => onClick(particion)}
                            >
                                <img src={partIMG} alt="part" style={{width: "100px"}} />
                                <div class="circle">                                    
                                    {particion}
                                </div>                                
                            </div>
                        )
                    })
                ):(
                    <div>No hay particiones disponibles</div>
                )}
            </div> 
        </div> 
    );
}