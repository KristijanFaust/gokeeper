import { Switch, Route, Redirect } from 'react-router-dom';
import {useEffect} from "react";
import {useState} from "react";

import Header from './components/header/header.component'
import SignInPage from "./pages/sign-in/sign-in-page.component";
import DashboardPage from "./pages/dashboard/dashboard-page.component";

import './App.scss';

function App() {
  const [authenticationToken, setAuthenticationToken] = useState(localStorage.getItem('authenticationToken'))

  useEffect(() => {
    setAuthenticationToken(localStorage.getItem('authenticationToken'));
  }, [authenticationToken]);

  let navigationLinks;
  authenticationToken ? navigationLinks = {'dashboard': ''} : navigationLinks = {'sign in': 'sign-in', 'sign up': 'sign-up'};

  return (
    <div className="App">
      <Header navigationLinks={navigationLinks} />
      <Switch>
        <Route exact path='/' render={() => authenticationToken ? (<DashboardPage />) : (<Redirect to='/sign-in' />)} />
        <Route exact path='/sign-in' render={() => authenticationToken ?
          (<Redirect to='/' />) : (<SignInPage signInCallback={setAuthenticationToken} />)}
        />
      </Switch>
    </div>
  );
}

export default App;
