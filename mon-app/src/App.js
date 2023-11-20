//import logo from './logo.svg';
import './App.css';


import React, { useState } from 'react';
import Film from './Film.jsx';
import Homepage from './Homepage.jsx';
import Search from './Search.jsx';
import Login from './Login.jsx';
import User from './User.jsx';
import Header from './Header.jsx';
import Register from './Register.jsx';


const App = () => {
  const [state, setState] = useState({page:"home"})
  console.log(state)

  if (state.page === "register"){
    return (<Register setState={setState} />);
  }

  if (state.page === "home") {
    if(state.userid === null){
      setState({page:"login"})
    } else {
      return (<div>
          <Header state={state} setState={setState} />
          <Homepage state={state} setState={setState} />
        </div>);
    }
  }

  if (state.page === "profile") {
    if(state.userid === null){
      setState({page:"login"})
    } else {
      state.uservisit = state.userid
      return (<div>
          <Header state={state} setState={setState} />
          <User state={state} setState={setState} />
        </div>);
    }
  }

  if (state.page === "login") {
    return (<Login setState={setState} />);
  }

  if (state.page === "film" && state.film !== null) {
    return (<div>
        <Header state={state} setState={setState} />
        <Film state={state} setState={setState} />
      </div>);
  }

  if (state.page === "search") {
    return (<div>
        <Header state={state} setState={setState} />
        <Search state={state} setState={setState} />
      </div>);
  }

  if (state.page === "user" && state.uservisit !== null) {
    return (<div>
        <Header state={state} setState={setState} />
        <User state={state} setState={setState} />
      </div>);
  }


}

export default App;
export const apiRoot = "http://localhost:54321";