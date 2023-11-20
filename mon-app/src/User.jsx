import React, { useState, useEffect } from "react";
import axios from "axios";
import {apiRoot} from './App.js'


export default function User({state, setState}) {
  const [details, setDetails] = useState(null);
  const [comments, setComments] = useState([]);
  const [favorites, setFavorites] = useState([]);
  
  useEffect(() => {
    if(state.uservisit){
      const id = state.uservisit
      axios.get(`${apiRoot}/user/${id}`).then((response) => {
        if (response.status == 200){
          setDetails(response.data);
        }
      });
      
      axios.get(`${apiRoot}/comment/user/${id}`).then((response) => {
        if (response.status == 200){
          setComments(response.data);
        }
      });
        
      axios.get(`${apiRoot}/favorites/user/${id}`).then((response) => {
        if (response.status == 200){
          setFavorites(response.data);
        }
      });
    }
  }, [state.uservisit]);

  const deleteAccount = () => {
    const del = window.confirm(`Supprimer le compte ${details.username} ?`);
    if (del) {
      const data = {auth:{userid:state.userid, session:state.session}}
      axios.delete(`${apiRoot}/user/${state.userid}`, {data:data}).then((response)=>{
        if (response.status == 200) {
          setState({page:"login"})
        } else {
          console.log(response.data)
          window.alert("Une erreur est survenue")
        }
      })
    }
  }
  
  const deleteComment = (comment) => {
    const data = {auth:{userid:state.userid, session:state.session}};
    console.log(data)
    axios.delete(`${apiRoot}/comment/${comment.id}`, {data}).then((response) => {
      if (response.status == 200){
        const new_comments = comments.filter((c) => (c.id!==comment.id));
        setComments(new_comments);
      }
    });
  }

  return (
    <div className="beautify centered">
      <div>
        {details?(
          <div className="beautify centered">
            <p className="big">Username : {details.username}</p>
            {state.userid === state.uservisit ? (
              <button onClick={deleteAccount}>Supprimer le compte</button>
            ):(
              <></>
            )}
          </div>
        ):(<h1>Chargement</h1>)}
      </div>
      <div>
        <h1>Favoris :</h1>
        {favorites?favorites.map((fav) => (
          <div className="film-container" >
            <img src={fav.image}  className="film-poster link" alt={fav.title} onClick={() => {setState({page:"film", userid:state.userid, film:fav.id, session:state.session})}}/>  
            <div className="film-details">
              <h2 className="film-title link" onClick={() => {setState({page:"film", userid:state.userid, film:fav.id, session:state.session})}}>{fav.title}</h2>
              <p className="film-description">{fav.description}</p>
            </div>
          </div>
          )):(<h2>Aucun film favoris</h2>)}
      </div>
      <div>
        <h1>Commentaires :</h1>
        {comments?comments.map((com) => (
          <div className="film-container comment">
            <img className="film-poster link" src={com.film.image} alt={com.film.title} onClick={() => {setState({page:"film", userid:state.userid, session:state.session, film:com.film.id})}}/>
            <div className="film-details">
              <h2 className="film-title link" onClick={() => {setState({page:"film", userid:state.userid, session:state.session, film:com.film.id})}}>{com.film.title}</h2>
              <p className="film-description">{com.contenu}</p>
              {(com.user_id===state.userid)?(
                <div><br/><button onClick={()=>deleteComment(com)}>Supprimer</button></div>
              ):(<></>)}
              </div>
          </div>
        )):(<h2>Aucun commentaires</h2>)}
      </div>
    </div>
  );
}