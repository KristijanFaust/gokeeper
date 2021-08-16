import React from 'react';
import ReactDOM from 'react-dom';
import {BrowserRouter} from 'react-router-dom';
import {ApolloClient, ApolloProvider, InMemoryCache} from '@apollo/react-hooks';

import App from './App';

import './index.scss';

const gqlClient = new ApolloClient({
  uri: 'http://localhost:8080/query',
  cache: new InMemoryCache(),
  headers: { Authentication: localStorage.getItem('authenticationToken') || '' }
});

ReactDOM.render(
  <React.StrictMode>
    <ApolloProvider client={gqlClient}>
      <BrowserRouter>
        <App/>
      </BrowserRouter>
    </ApolloProvider>
  </React.StrictMode>,
  document.getElementById('root')
);
