import { Routes, Route, HashRouter, Link } from 'react-router-dom'
import { useState } from 'react';

import hacker from '../iconos/data-entry.png';
import Consola from '../Pages/Consola';
import Discos from '../Pages/Discos';
import Partitions from '../Pages/Particiones';
import Login from '../Pages/Login';
import Explorer from '../Pages/SArchivos';
import "../Stylesheets/navbar.css"

export default function Navegador(){
    const [ ip, setIP ] = useState("localhost")
    

    const handleChange = (e) => {
        console.log(e.target.value)
        setIP(e.target.value)
    }
    
    const logOut = (e) => {
        e.preventDefault()
        
        fetch(`http://${ip}:8080/logout`)
        .then(Response => Response.json())
        .then(rawData => {
            console.log(rawData);  
            if (rawData === 0){
                alert('sesion cerrada')
                window.location.href = '#/Login';
            }else{
                alert('No hay sesion abierta')
            }
        }) 
        .catch(error => {
            console.error('Error en la solicitud Fetch:', error);
            
        });
    };

    const limpiar = (e) => {
        e.preventDefault()
        console.log("limpiando")
        fetch(`http://${ip}:8080/limpiar`)
        .then(Response => Response.json())
        .then(rawData => {
            console.log(rawData); 
            if (rawData === 1){
                alert('Discos y reportes eliminados')
                window.location.href = '#/Comandos';
            }else{
                alert('Error al eliminar archiovs')
            }
        }) 
    }

    return(
        <HashRouter>
            <nav className="navbar navbar-expand-lg navbar-dark bg-dark justify-content-center">
            <div className="container-fluid d-flex align-items-center justify-content-center w-100">
                <img id="imgIcon" src={hacker} alt="" width="64" height="64" className="d-inline-block align-text-top me-3" />

                <a className="navbar-brand me-3">
                    Proyecto2
                </a>

                <ul className="navbar-nav me-3">
                    <li className="nav-item">
                        <Link className="nav-link active" to="/Consola">Comandos</Link>
                    </li>
                    <li className="nav-item">
                        <Link className="nav-link" to="/Login">Navegador de Archivos</Link>
                    </li>
                    <li className="nav-item">
                        <button id="btnLogOut" onClick={logOut} className="nav-link btn btn-link">Cerrar Sesion</button>
                    </li>
                </ul>

                <input id="InIP" className="form-control" style={{ maxWidth: "180px" }} placeholder="IP" onChange={handleChange} />
            </div>
        </nav>

            
            <Routes>
                <Route path="/" element ={<Consola newIp={ip}/>}/> {/*home*/}
                <Route path="/Consola" element ={<Consola newIp={ip}/>}/> 
                <Route path="/Login" element ={<Login newIp={ip}/>}/>
                <Route path="/Discos/:id" element ={<Discos newIp={ip}/>}/> 
                <Route path="/Particiones/:id" element ={<Partitions newIp={ip}/>}/> 
                <Route path="/Explorador/:id" element ={<Explorer newIp={ip}/>}/>              
            </Routes>
        </HashRouter>
    );
}