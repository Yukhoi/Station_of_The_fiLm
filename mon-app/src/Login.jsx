import React, { useState } from 'react';
import axios from 'axios';
import {apiRoot} from './App.js'

export default function Login ({state, setState}) {
  const [email, setDataEmail] = useState('');
  const [password, setDataPassword] = useState('');

  const setPassword = (pass) => {
    setDataPassword(pass.trim())
  }

  const setEmail = (word) => {
    setDataEmail(word.trim())
  }

  const handleSubmit = (e) => {
    if (email.includes("@") && password !== '') {
      axios.post(`${apiRoot}/login/`, { email:email, password:password }).then((response) => {
        if (response.status === 200) {
          console.log("next state : ", {page:"home", userid:response.data.userid, session:response.data.session})
          console.log(response)
          setState({page:"home", userid:response.data.userid, session:response.data.session})
        } else {
          console.log(response.response.data)
          window.alert("Erreur de connection : " + response.response.data.error)
        }
      }).catch((response) => {
        console.log(response.response.data)
        window.alert("Erreur de connection : " + response.response.data.error)
      })      
    } else {
      window.alert("Email mal formé ou mot de passe non rempli")
    }
  };

  return (
    <div className="beautify center">
      <center><h2>Se connecter</h2></center>
      <br/>
      <label htmlFor="email">Email</label>
      <input 
        className='in'
        type="text"
        placeholder="email@example.com"
        id="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <label htmlFor="password">Mot de passe</label>
      <input
        className='in'
        type="password"
        placeholder="Password"
        id='password'
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        onKeyDown={(e) => {if (e.key == 'Enter') {handleSubmit()}}}
      />
      
      <button onClick={handleSubmit}>Login</button>
      <br/>
      <a className="link small" onClick={()=>{setState({page:"register"})}}>Créer un compte</a>
    </div>
  );
};