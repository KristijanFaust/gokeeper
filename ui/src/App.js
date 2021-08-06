import { Switch, Route, Redirect } from 'react-router-dom';

import Header from './components/header/header.component'
import SignIn from "./pages/sign-in/sign-in-page.component";

import './App.scss';

function App() {
  return (
    <div className="App">
      <Header />
      <Switch>
        <Route exact path='/' render={() => <Redirect to='/sign-in' />} />
        <Route exact path='/sign-in' component={SignIn} />
      </Switch>
    </div>
  );
}

export default App;
