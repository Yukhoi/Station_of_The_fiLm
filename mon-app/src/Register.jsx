import React, { useState } from 'react';
import axios from 'axios';
import {apiRoot} from './App.js'

export default function Register ({state, setState}) {
  const [email, setDataEmail] = useState('');
  const [password, setDataPassword] = useState('');
  const [passVerif, setDataPassVerif] = useState('');
  const [username, setDataUsername] = useState('');

  const setPassword = (pass) => {
    setDataPassword(pass.trim())
  }

  const setEmail = (word) => {
    setDataEmail(word.trim())
  }
  
  const setPassVerif = (word) => {
    setDataPassVerif(word.trim())
  }
  
  const setUsername = (word) => {
    setDataUsername(word.trim())
  }

  const handleSubmit = async () => {
    if (email.includes("@") && password != '' && password == passVerif && username != ''){
      try {
        const data = {email:email, password:password, username:username}
        const response = await axios.put(`${apiRoot}/user/`, data);
        if (response.status === 200) {
            setState({page:"login"})
        }
      } catch (error) {
        console.error(error);
        window.alert(`Erreur : ${error}`)
      }
    }else {
      window.alert("Email mal formé, mots de passe non remplis ou ne correspondent pas")
    }
  };

  return (
    <div className='beautify center'>
      <h2>Créer un compte</h2>
      <label htmlFor="email">Email</label>
      <input
        className='in'
        type="text"
        placeholder="email"
        id='email'
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <label htmlFor="username">Non d'utilisateur</label>
      <input
        className='in'
        type="text"
        placeholder="username"
        id='username'
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <label htmlFor="password">Mot de passe</label>
      <input
        className='in'
        id='password'
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      />
      <label htmlFor="repass">Retapez le mot de passe</label>
      <input
        className='in'
        id='repass'
        type="password"
        placeholder="Password"
        value={passVerif}
        onChange={(e) => setPassVerif(e.target.value)}
        onKeyDown={(e) => {if (e.key == 'Enter') {handleSubmit()}}}
      />
      <button onClick={handleSubmit}>Register</button>
      <br/>
      <a className="link small" onClick={()=>{setState({page:"login"})}}>Se connecter</a>
    </div>
  );
};