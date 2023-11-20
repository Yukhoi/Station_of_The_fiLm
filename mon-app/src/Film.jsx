import React, { useState, useEffect } from "react";
import axios from "axios";
import {apiRoot} from './App.js'

export default function Film({ state , setState }) {
  const [details, setDetails] = useState(null);
  const [comments, setComments] = useState([]);
  const [commentcontent, setCommentcontent] = useState("");
  const [favorite, setFavorite] = useState(null);
  
  const userid = state.userid

  useEffect(() => {
    if(state.film){
      const id = state.film
      console.log("id :"+id)
      axios.get(`${apiRoot}/film/${id}`).then((response) => {
        setDetails(response.data);
      });
      
      axios.get(`${apiRoot}/comment/film/${id}`).then((response) => {
        setComments(response.data);
      });

      if(state.userid){
        const uid = state.userid
        axios.get(`${apiRoot}/favorite/${id}/${uid}`).then((response) => {
          setFavorite(response.data);
        });
      }
    }
  }, [state.film, state.userid]);

  const changeFavorite = () => {
    if(favorite === null) return;

    const film_id = favorite.film_id
    const user_id = favorite.user_id
    const fav = favorite.favorite
    
    console.log(state)
    console.log(favorite)

    const data = {film_id:film_id, user_id:user_id, favorite:!fav, auth:{userid:state.userid, session:state.session}}
    axios.post(`${apiRoot}/favorite/`, data).then((response) => {
      if (response.status == 200){
        setFavorite(response.data);
      }
    });
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

  const publishComment = () => {
    const comment = commentcontent
    console.log(comment)
    const film_id = state.film;
    const user_id = state.userid;
    const data = {film_id:film_id, user_id:user_id, contenu:comment.trim(), auth:{userid:state.userid, session:state.session}};
    axios.put(`${apiRoot}/comment/`, data).then((response) => {
      if (response.status == 200){
        setComments([response.data].concat(comments));
        setCommentcontent("")
      }
    });
  }

  return (
    <div className="beautify">
      <div>
        {details?(
          <div className="film-container" id={details.id} >
            <img className="film-poster" src={details.image} alt={details.title} />
            <div className="film-details">
              <h2 className="film-title">{details.title}</h2>
              <h4>Description:</h4>
              <p className="film-description">{details.description}</p>
            </div>
          </div>
        ):(<h1>Chargement</h1>)}
      </div>
      <div>
        <p>Favoris : 
      {favorite?(
        <button onClick={changeFavorite} title={favorite.favorite?"Ne plus aimer ce film":"Aimer ce film"}>{favorite.favorite?"Vous aimez ce film":"Vous n'aimez pas ce film"}</button>
      ):(favorite===null?(<span>Connectez vous</span>):(<span>Chargement</span>))}
        </p>
      </div>
        {(state.userid!==null)?(
          <div>
            <textarea name="comment" onChange={(e) => {setCommentcontent(e.target.value)}} value={commentcontent} placeholder="Commenter sur ce film" id="comment" cols="30" rows="4"></textarea>
            <br/>
            <button onClick={publishComment}>Publier</button>
          </div>)
        :(
          <div>
            <p>Connectez vous pour commenter sur ce film</p>
          </div>
        )} 
      {comments?(
        <div>
          {comments.map((com) => (
            <div className="film-container beautify centered comment">
              <p>
                <span className="link bold" onClick={() => setState({page:"user", uservisit:com.user.id, userid:userid, session:state.session})}>
                  {com.user.username}
                </span> a commenté
              </p>
              <div className="film-details">
                <p>{com.contenu}</p>
              
                {(com.user.id===state.userid)?(
                  <button onClick={()=>deleteComment(com)}>Supprimer</button>
                ):(<></>)}
              </div>
            </div>
          ))}
        </div>
      ):(<h1>Chargement</h1>)}
    </div>
  );
}