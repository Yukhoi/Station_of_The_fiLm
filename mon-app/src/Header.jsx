import React, { useState, useEffect } from "react";
import axios from "axios";

import {apiRoot} from './App.js'

import LogoutIcon from '@mui/icons-material/Logout';
import SearchIcon from '@mui/icons-material/Search';
import HomeIcon from '@mui/icons-material/Home';
import AccountBoxIcon from '@mui/icons-material/AccountBox';

import './Header.css'

export default function Header({ state , setState }) {

  const logout = () => {
    const data = {auth: {userid:state.userid, session:state.session}}
    axios.post(`${apiRoot}/logout/`, data).then(() =>{
      setState({page:"login"})
    }).catch(() => {
      setState({page:"login"})
    });
    
  }

  return (
    <center className="beautify header">
      <nav style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ textAlign: 'center' }}>Station of The fiLms</h1>
        <div>
          <a className="link big" onClick={()=>{setState({page:"home", session:state.session, userid:state.userid})}}>
            <HomeIcon />
          </a>
          <a className="link big" onClick={()=>{setState({page:"profile", session:state.session, userid:state.userid})}}>
            <AccountBoxIcon />
          </a>
          <a className="link big" onClick={()=>{setState({page:"search", session:state.session, userid:state.userid})}}>
            <SearchIcon />
          </a>
          {state.userid!='' ? (
            <a className="link big" onClick={()=>{logout()}}>
              <LogoutIcon />
            </a>
          ) : (
            <a className="link big" onClick={()=>{setState({page:"login"})}}>Login</a>
          )}
        </div>
      </nav>
    </center>
  );
}
