import React, { useState, useEffect } from "react";
import axios from "axios";

import {apiRoot} from './App.js'

export default function Homepage({ state , setState }) {
  const [films, setFilms] = useState([]);
  const [refilms, setReFilms] = useState([]);
  const userid = state.userid 

  useEffect(() => {
    if (!userid) {
      setState({page:"login"})
    } else {
      axios.get(`${apiRoot}/recommend/user/${userid}`).then((response) => {
        if (response.status == 200){
          setReFilms(response.data);
        }
      });
      axios.get(`${apiRoot}/favorites/user/${userid}`).then((response) => {
        if (response.status == 200){
          setFilms(response.data);
        }
      });
    }
  }, [state.userid]);
  return (
    <div className="beautify">
      <div>
        <h1>Vous avez aimé</h1>
        {films?(films.map((film) => (
          <div className="film-container" id={film.id} >
            <img onClick={()=>{setState({page:"film", film:film.id, userid:userid, session:state.session})}} src={film.image} alt={film.title} className="film-poster link" />
            <div className="film-details">
              <h2 onClick={()=>{setState({page:"film", film:film.id, userid:userid, session:state.session})}} className="film-title link">{film.title}</h2>
              <p className="film-description">{film.description}</p>
            </div>
          </div>
        ))):(<h2>Chargement</h2>)}
      </div>
      <div>
        <h1>Nous vous recommandons</h1>
        {refilms?(refilms.map((film) => (
          <div className="film-container recommend" id={film.id} >
            <img onClick={()=>{setState({page:"film", film:film.id, userid:userid, session:state.session})}} src={film.image} alt={film.title} className="film-poster link" />
            <div className="film-details">
              <h2 onClick={()=>{setState({page:"film", film:film.id, userid:userid, session:state.session})}} className="film-title link">{film.title}</h2>
              <p className="film-description">{film.description}</p>
            </div>
          </div>
        ))):(<h2>Chargement</h2>)}
      </div>
    </div>
  );
}