import {apiRoot} from './App.js'
import React, { useState } from 'react';
import axios from 'axios';

export default function Search({state, setState}) {
    const [word, setWord] = useState('');
    const [type, setType] = useState('title');
    const [searchResults, setSearchResults] = useState([]);

    const handleSearch = async () => {
        try {
            let url = `${apiRoot}/search`;
            if (type === 'title') {
                if (!word) {
                    window.alert('Please enter a search term');
                    return;
                  }
                url += `/title/${word}`;
                const response = await axios.get(url);
                setSearchResults(response.data);
            } else if (type === 'genre') {
                if (!word) {
                    window.alert('Please enter a search term');
                    return;
                  }
                url += `/genre/${word}`;
                const response = await axios.get(url);
                setSearchResults(response.data);
            } 
            
        } catch (error) {
            console.log(error)
        }
    }

    return (
        <div className="beautify">
            <div className='beautify'>
            <fieldset>
            <legend>Rechercher les films par :</legend>

            <div>
            <input type="radio" id="title" name="drone" value="title" checked={type==="title"} onChange={()=>setType("title")}/>
            <label for="title">Titre</label>
            </div>

            <div>
            <input type="radio" id="genre" value="genre" checked={type==="genre"} onChange={()=>setType("genre")}/>
            <label for="genre">Genre</label>
            </div>
        {/*
            <div>
            <input type="radio" id="keyword" value="keyword" checked={type==="keyword"} onClick={setType}/>
            <label for="keyword">Keyword</label>
            </div>
    */}
        </fieldset>
            
                {type==="title"?<input 
                    className='in'
                    type="text" 
                    placeholder='Search text'
                    value={word} 
                    onChange={(e) => setWord(e.target.value)} 
                />:<select onChange={(e) => setWord(e.target.value)}>
                        <option value="Adventure">Adventure</option>
                        <option value="Family">Family</option>
                        <option value="Fantasy">Fantasy</option>
                        <option value="Crime">Crime</option>
                        <option value="Drama">Drama</option>
                        <option value="Comedy">Comedy</option>
                        <option value="Animation">Animation</option>
                        <option value="Sci-Fi">Sci-Fi</option>
                        <option value="Sport">Sport</option>
                        <option value="Action">Action</option>
                        <option value="Thriller">Thriller</option>
                        <option value="Mystery">Mystery</option>
                        <option value="Western">Western</option>
                        <option value="Romance">Romance</option>
                        <option value="Biography">Biography</option>
                        <option value="Horror">Horror</option>
                        <option value="War">War</option>
                        <option value="Musical">Musical</option>
                        <option value="History">History</option>
                        <option value="Music">Music</option>
                        <option value="Documentary">Documentary</option>

                        

                    </select>}
                <button onClick={() => handleSearch()}>Chercher</button>
            </div>
            <div>
                {searchResults?searchResults.map((film, index) => (
                    <div className="film-container recommend" id={film.id} >
                    <img onClick={()=>{setState({page:"film", film:film.id, userid:state.userid, session:state.session})}} src={film.image} alt={film.title} className="film-poster link" />
                    <div className="film-details">
                      <h2 onClick={()=>{setState({page:"film", film:film.id, userid:state.userid, session:state.session})}} className="film-title link">{film.title}</h2>
                      <p className="film-description">{film.description}</p>
                    </div>
                  </div>
                )):<></>}
            </div>
        </div>
    )
}