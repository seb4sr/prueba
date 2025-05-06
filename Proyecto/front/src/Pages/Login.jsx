import { useParams, useNavigate } from 'react-router-dom';
import { useState } from "react";

import "../Stylesheets/Login.css"
import user from '../iconos/profile.png';
import key from '../iconos/key.png';
import partDisk from '../iconos/IdPart.png';

export default function Login({newIp="localhost"}){
    const { disk, part } = useParams()
    const [ estado, setEstado ] = useState();
    const navigate = useNavigate()

    const handleSubmit = (e) => {
        e.preventDefault()
  
        const user = e.target.uname.value
        const pass = e.target.psw.value
        const id = e.target.particion.value
  
        console.log("user", user, pass, id)

        const data = {
            usuario: user,
            password: pass,
            id: id
        };

        fetch(`http://${newIp}:8080/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        })
        .then(Response => Response.json())
        .then(rawData => {
            console.log(rawData); 
            setEstado(rawData);
            if (rawData === -1){
                onClick(id)
            }
        })
    }

    const onClick = (id) => {
        console.log("nueva pagina",id)
        //navigate(`/explorador/${particion}`)
        navigate(`/Discos/${id}`)
    }

    return(
        <>
            <div className="container">
                <div className="d-flex justify-content-center">
                    <div className="card ">
                        <div className="card-header">
                            <h3>Inicio de Sesion</h3>
                        </div>
                        <div className="card-body">
                            <form onSubmit={handleSubmit}>
                                <div className="input-group form-group">
                                    <input type="text" className="form-control" placeholder="ID particion" name="particion" required/>
                                </div>
                                <div className="input-group form-group">
                                    <input type="text" className="form-control" placeholder="username" name="uname" required/>
                                </div>
                                <div className="input-group form-group">
                                    <input type="password" className="form-control" placeholder="password" name="psw" required/>
                                </div>
                                <div style={{textAlign:'center'}}>
                                    <button type="submit" className="btn btn-primary login_btn">Log In</button>
                                </div>
                            </form>

                            <div>&nbsp;&nbsp;&nbsp;</div>

                            <div className='estadoLogin'>
                                {estado === 0 ? (
                                    <div>Ya existe sesion activa</div>
                                ):estado === 2 ?(
                                    <div>Particion sin formato</div>
                                ):estado === 3 ?(
                                    <div>Contraseña incorrecta</div>
                                ):estado === 4 ?(
                                    <div>No se encontro el usuario</div>
                                ):estado === 5 ?(
                                    <div>Ocurrio un error inesperado</div>
                                ):(
                                    <div></div>
                                )}
                            </div>
                            
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}